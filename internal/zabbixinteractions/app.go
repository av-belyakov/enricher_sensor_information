package zabbixinteractions

/*
// NewZabbixConnectionJsonRPC создает объект соединения с Zabbix API
func NewZabbixConnectionJsonRPC(settings SettingsZabbixConnectionJsonRPC) (*connectionjsonrpc.ZabbixConnectionJsonRPC, error) {
	var (
		zabbixConn *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error
	)

	connTimeout := 30 * time.Second
	if settings.ConnectionTimeout > (1 * time.Second) {
		connTimeout = settings.ConnectionTimeout
	}

	if settings.Host == "" {
		return zabbixConn, supportingfunctions.CustomError(errors.New("the value 'host' should not be empty"))
	}

	if settings.UseTLS {
		zabbixConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithTLS(),
			connectionjsonrpc.WithInsecureSkipVerify(),
			connectionjsonrpc.WithHost(settings.Host),
			connectionjsonrpc.WithPort(settings.Port),
			connectionjsonrpc.WithLogin(settings.Login),
			connectionjsonrpc.WithPasswd(settings.Passwd),
			connectionjsonrpc.WithConnectionTimeout(cmp.Or(settings.ConnectionTimeout, connTimeout)),
		)
	} else {
		zabbixConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost(settings.Host),
			connectionjsonrpc.WithPort(settings.Port),
			connectionjsonrpc.WithLogin(settings.Login),
			connectionjsonrpc.WithPasswd(settings.Passwd),
			connectionjsonrpc.WithConnectionTimeout(cmp.Or(cfg.GetZabbix().Timeout, connTimeout)),
		)
	}

	return zabbixConn, err
}

// Authorization запрос к Zabbix с целью получения хеша авторизации необходимого для
// дальнейшей работы с API
func (zc *ZabbixConnectionJsonRPC) Authorization(ctx context.Context) error {
	data := strings.NewReader(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"user.login",
	  "params": {
	    "username":"%s",
		"password":"%s"
	  },
	  "id":1
	}`, zc.login, zc.passwd))

	result := ZabbixAuthorizationData{}
	res, err := zc.PostRequest(ctx, data)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return supportingfunctions.CustomError(err)
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

// GetAuthorizationData хеш авторизации
func (zc *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return zc.authorizationHash
}

// PostRequest HTTP запрос типа POST
func (zc *ZabbixConnectionJsonRPC) PostRequest(ctx context.Context, data *strings.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", zc.url, data)
	if err != nil {
		return []byte{}, supportingfunctions.CustomError(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", zc.authorizationHash))
	req.Header.Set("Content-Type", "application/json-rpc")

	res, err := zc.connClient.Do(req)
	if err != nil {
		return []byte{}, supportingfunctions.CustomError(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, supportingfunctions.CustomError(fmt.Errorf("error sending the request, response status is %s", res.Status))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, supportingfunctions.CustomError(err)
	}

	return resBody, nil
}
*/
