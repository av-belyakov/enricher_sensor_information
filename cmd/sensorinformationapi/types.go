package sensorinformationapi

import (
	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// SensorInformationClient клиент для взаимодействия с API
type SensorInformationClient struct {
	ncirccConn *ncirccinteractions.ClientNICRCC
	zabbixConn *zabbixinteractions.ZabbixConnectionJsonRPC
	settings   SensorInformationSettings
}

// SensorInformationSettings настройки модуля
type SensorInformationSettings struct {
	host           string
	user           string
	passwd         string
	ncirccURL      string
	ncirccToken    string
	port           int
	requestTimeout int
}

// sensorInformationClientOptions функциональные параметры
type sensorInformationClientOptions func(*SensorInformationClient) error
