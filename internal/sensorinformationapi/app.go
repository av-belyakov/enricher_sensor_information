package sensorinformationapi

import (
	"errors"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// New настраивает новый модуль взаимодействия с API
func New(opts ...sensorInformationClientOptions) (*SensorInformationClient, error) {
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
	zConn, err := zabbixinteractions.NewZabbixConnectionJsonRPC(
		zabbixinteractions.SettingsZabbixConnectionJsonRPC{
			Host:              api.settings.host,
			Login:             api.settings.user,
			Passwd:            api.settings.passwd,
			ConnectionTimeout: time.Duration(api.settings.requestTimeout) * time.Second,
		})
	if err != nil {
		return api, err
	}
	api.zabbixConn = zConn

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

// WithHost имя или ip адрес сервера Zabbix
func WithHost(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		sic.settings.host = v

		return nil
	}
}

// WithPort порт сервера Zabbix
func WithPort(v int) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		sic.settings.port = v

		return nil
	}
}

// WithUser имя пользователя
func WithUser(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'user' cannot be empty")
		}

		sic.settings.user = v

		return nil
	}
}

// WithPasswd пароль пользователя пользователя
func WithPasswd(v string) sensorInformationClientOptions {
	return func(sic *SensorInformationClient) error {
		if v == "" {
			return errors.New("the value of 'passwd' cannot be empty")
		}

		sic.settings.passwd = v

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
