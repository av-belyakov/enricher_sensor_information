package sensorinformationapi

import (
	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"

	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/netboxinteractions"
)

// SensorInformationClient клиент для взаимодействия с API
type SensorInformationClient struct {
	ncirccConn *ncirccinteractions.ClientNICRCC
	zabbixConn *connectionjsonrpc.ZabbixConnectionJsonRPC
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
	zabbixPort     int
	requestTimeout int
	zabbixUseTLS   bool
}

// sensorInformationClientOptions функциональные параметры
type sensorInformationClientOptions func(*SensorInformationClient) error

type TenantGroupsInformation struct {
	SiensorId, Display, Name string
}
