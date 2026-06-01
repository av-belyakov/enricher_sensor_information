package sensorinformationapi

import (
	"errors"
	"time"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"

	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/netboxinteractions"
)

// New настраивает новый модуль взаимодействия с API
func New(opts ...sensorInformationClientOptions) (*SensorInformationClient, error) {
	var (
		zConn *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error
	)

	api := &SensorInformationClient{
		settings: SensorInformationSettings{
			requestTimeout: 10,
		},
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	//инициализация соединения с Zabbix
	if api.settings.zabbixUseTLS {
		zConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithTLS(),
			connectionjsonrpc.WithInsecureSkipVerify(),
			connectionjsonrpc.WithHost(api.settings.zabbixHost),
			connectionjsonrpc.WithPort(api.settings.zabbixPort),
			connectionjsonrpc.WithLogin(api.settings.zabbixUser),
			connectionjsonrpc.WithPasswd(api.settings.zabbixPasswd),
			connectionjsonrpc.WithConnectionTimeout(api.settings.requestTimeout),
		)
	} else {
		zConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost(api.settings.zabbixHost),
			connectionjsonrpc.WithPort(api.settings.zabbixPort),
			connectionjsonrpc.WithLogin(api.settings.zabbixUser),
			connectionjsonrpc.WithPasswd(api.settings.zabbixPasswd),
			connectionjsonrpc.WithConnectionTimeout(api.settings.requestTimeout),
		)
	}
	if err != nil {
		return api, err
	}
	api.zabbixConn = zConn

	//инициализация соединения с Netbox
	netboxConn, err := netboxinteractions.New(
		api.settings.netboxToken,
		netboxinteractions.WithHost(api.settings.netboxHost),
		netboxinteractions.WithPort(api.settings.netboxPort),
		netboxinteractions.WithTimeout(api.settings.requestTimeout),
	)
	if err != nil {
		return api, err
	}
	api.netboxConn = netboxConn

	//инициализация соединения с НКЦКИ
	ncirccConn, err := ncirccinteractions.NewClient(
		api.settings.ncirccURL,
		api.settings.ncirccToken,
		time.Duration(api.settings.requestTimeout)*time.Second,
	)
	if err != nil {
		return api, err
	}
	api.ncirccConn = ncirccConn

	return api, nil
}

// WithZabbixUseTLS использовать ли протокол TLS для доступа к серверу Zabbix
func WithZabbixUseTLS(v bool) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		sic.settings.zabbixUseTLS = v

		return nil
	}
}

// WithZabbixHost имя или ip адрес сервера Zabbix
func WithZabbixHost(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'host' for Zabbix cannot be empty")
		}

		sic.settings.zabbixHost = v

		return nil
	}
}

// WithZabbixPort порт сервера Zabbix
func WithZabbixPort(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		sic.settings.zabbixPort = v

		return nil
	}
}

// WithZabbixUser имя пользователя для сервера Zabbix
func WithZabbixUser(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'user' for Zabbix cannot be empty")
		}

		sic.settings.zabbixUser = v

		return nil
	}
}

// WithZabbixPasswd пароль пользователя для сервера Zabbix
func WithZabbixPasswd(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'passwd' for Zabbix cannot be empty")
		}

		sic.settings.zabbixPasswd = v

		return nil
	}
}

// WithNetboxHost имя или ip адрес сервера Netbox
func WithNetboxHost(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'host' for Netbox cannot be empty")
		}

		sic.settings.netboxHost = v

		return nil
	}
}

// WithNetboxPort порт сервера Netbox
func WithNetboxPort(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		sic.settings.netboxPort = v

		return nil
	}
}

// WithNetboxToken токен сервера Netbox
func WithNetboxToken(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'netboxToken' cannot be empty")
		}

		sic.settings.netboxToken = v

		return nil
	}
}

// WithNCIRCCURL URL API НКЦКИ
func WithNCIRCCURL(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'ncirccURL' cannot be empty")
		}

		sic.settings.ncirccURL = v

		return nil
	}
}

// WithNCIRCCToken токен API НКЦКИ
func WithNCIRCCToken(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'ncirccToken' cannot be empty")
		}

		sic.settings.ncirccToken = v

		return nil
	}
}

// WithRequestTimeout ограничение времени выполнения запроса от 1 до 60 сек.
func WithRequestTimeout(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 1 || v > 60 {
			return errors.New("the request execution time should be in the range from 1 to 60 seconds")
		}

		sic.settings.requestTimeout = v

		return nil
	}
}
