package zabbixinteractions

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// NewZabbixConnectionJsonRPC создает объект соединения с Zabbix API с использование Json-RPC
func NewZabbixConnectionJsonRPC(ctx context.Context, settings SettingsZabbixConnectionJsonRPC) (*ZabbixConnectionJsonRPC, error) {
	var zc *ZabbixConnectionJsonRPC

	connTimeout := 30 * time.Second
	if settings.ConnectionTimeout > (1 * time.Second) {
		connTimeout = settings.ConnectionTimeout
	}

	if settings.Host == "" {
		return zc, supportingfunctions.CustomError(errors.New("the value 'host' should not be empty"))
	}

	client := &http.Client{Transport: &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     connTimeout,
		MaxIdleConnsPerHost: 10,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            x509.NewCertPool(),
		},
	}}

	zc = &ZabbixConnectionJsonRPC{
		url:             fmt.Sprintf("https://%s/api_jsonrpc.php", settings.Host),
		host:            settings.Host,
		login:           settings.Login,
		passwd:          settings.Passwd,
		applicationType: "application/json-rpc",
		connClient:      client,
	}

	return zc, authorizationZabbixAPI(ctx, settings.Login, settings.Passwd, zc)
}

// authorizationZabbixAPI делает запрос к Zabbix с целью получения хеша авторизации
// необходимого для дальнейшей работы с API
func authorizationZabbixAPI(ctx context.Context, login, passwd string, zc *ZabbixConnectionJsonRPC) error {
	ctxTimeout, close := context.WithTimeout(ctx, zc.connClient.Timeout)
	defer close()

	body := strings.NewReader(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"user.login\",\"params\":{\"username\":\"%s\",\"password\":\"%s\"},\"id\":1}", login, passwd))
	httpReq, err := http.NewRequestWithContext(ctxTimeout, "POST", zc.url, body)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	dataLen := body.Len()
	if dataLen > 0 {
		httpReq.ContentLength = int64(dataLen)
		httpReq.Body = io.NopCloser(body)
	}

	urlBase, err := url.Parse(zc.url)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	httpReq.URL = urlBase
	httpReq.URL.Path = "api_jsonrpc.php"
	/*
		httpReq.Header = http.Header{}
		httpReq.Header.Set("Authorization", client.AuthHash)
		httpReq.Header.Set("Content-type", "application/json")
		httpReq.Header.Set("Accept", "application/json")
	*/

	res, err := zc.connClient.Do(httpReq)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return supportingfunctions.CustomError(fmt.Errorf("error authorization, response status is %s", res.Status))
	}

	result := ZabbixAuthorizationData{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("error authorization, %w", err))
	}

	if len(result.Error) > 0 {
		var shortMsg, fullMsg string
		for k, v := range result.Error {
			if k == "message" {
				shortMsg = fmt.Sprint(v)
			}
			if k == "data" {
				fullMsg = fmt.Sprint(v)
			}
		}

		return supportingfunctions.CustomError(fmt.Errorf("error authorization, (%s %s)", shortMsg, fullMsg))
	}

	zc.authorizationHash = result.Result

	return nil
}

// GetAuthorizationData возвращает авторизационный хеш
func (zc *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return zc.authorizationHash
}

// SendPostRequest отправляет HTTP POST запрос с параметрами запроса вида JSON
func (zc *ZabbixConnectionJsonRPC) SendPostRequest(ctx context.Context, data *strings.Reader) (io.Reader, error) {
	var result io.Reader

	res, err := zc.connClient.Post(zc.url, zc.applicationType, data)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("error sending the request, %v %s:%d", err, f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)
		return result, fmt.Errorf("error sending the request, response status is %s %s:%d", res.Status, f, l-1)
	}

	return res.Body, nil
}
