package sensorinformationapi

import (
	"context"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
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

	conn, err := zabbixinteractions.NewZabbixConnectionJsonRPC(
		zabbixinteractions.SettingsZabbixConnectionJsonRPC{
			Host:              api.settings.host,
			Login:             api.settings.user,
			Passwd:            api.settings.passwd,
			ConnectionTimeout: time.Duration(api.settings.requestTimeout) * time.Second,
		})
	if err != nil {
		return api, err
	}

	api.zabbixConn = conn

	return api, nil
}

// SearchSensorInfo поиск информации о сенсоре
func (api *SensorInformationClient) SearchSensorInfo(ctx context.Context, sensorId string) (responses.DetailedInformation, error) {
	return zabbixinteractions.GetFullSensorInfo(ctx, sensorId, api.zabbixConn)
}
