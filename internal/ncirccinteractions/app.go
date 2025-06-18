package ncirccinteractions

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// NewClient клиент для доступа к НКЦКИ
func NewClient(url, token string, connTimeout time.Duration) (*ClientNICRCC, error) {
	settings := &ClientNICRCC{
		connectionTimeout: 30 * time.Second,
	}

	if url == "" {
		return settings, supportingfunctions.CustomError(errors.New("the 'url' parameter must not be empty"))
	}

	if token == "" {
		return settings, supportingfunctions.CustomError(errors.New("the 'token' parameter must not be empty"))
	}

	settings.url = url
	settings.token = token

	if connTimeout > (1 * time.Second) {
		settings.connectionTimeout = connTimeout
	}

	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     connTimeout,
		MaxIdleConnsPerHost: 10,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            x509.NewCertPool(),
		},
	}

	settings.client = &http.Client{Transport: transport}

	return settings, nil
}

// GetFullNameOrganizationByINN запрос для получения полной информации об организации
// по её ИНН
func (api *ClientNICRCC) GetFullNameOrganizationByINN(ctx context.Context, inn string) (Response, error) {
	resData := Response{}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	req, err := http.NewRequestWithContext(ctxTimeout, "GET", api.url, strings.NewReader(""))
	if err != nil {
		return resData, supportingfunctions.CustomError(err)
	}

	req.Header.Add("x-token", api.token)

	q := req.URL.Query()
	q.Add("fields", "[\"settings_name\",\"settings_sname\",\"settings_inn_of_subject\",\"settings_subject_type\"]")
	q.Add("filter", fmt.Sprintf("[{\"property\":\"settings_inn_of_subject\",\"operator\":\"eq\",\"value\":\"%s\"}]", inn))
	q.Add("limit", "10")
	q.Add("start", "0")
	req.URL.RawQuery = q.Encode()

	res, err := api.client.Do(req)
	if err != nil {
		return resData, supportingfunctions.CustomError(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return resData, supportingfunctions.CustomError(fmt.Errorf("error sending the request INN '%s', response status is %s", inn, res.Status))
	}

	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return resData, supportingfunctions.CustomError(err)
	}

	return resData, nil
}
