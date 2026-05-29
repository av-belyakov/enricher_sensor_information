package sensorinformationapi

import (
	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/netboxinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// SensorInformationClient клиент для взаимодействия с API
type SensorInformationClient struct {
	ncirccConn *ncirccinteractions.ClientNICRCC
	zabbixConn *zabbixinteractions.ZabbixConnectionJsonRPC
	netboxConn *netboxinteractions.Client
	settings   SensorInformationSettings
}

// SensorInformationSettings настройки модуля
type SensorInformationSettings struct {
	zabbixPasswd   string
	zabbixHost     string
	zabbixUser     string
	ncirccToken    string
	ncirccURL      string
	netboxToken    string
	netboxHost     string
	netboxPort     int
	requestTimeout int
}

// sensorInformationClientOptions функциональные параметры
type sensorInformationClientOptions func(*SensorInformationClient) error

type TenantGroupsInformation struct {
	SiensorId, Display, Name string
}
