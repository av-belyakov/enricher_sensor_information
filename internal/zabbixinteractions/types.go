package zabbixinteractions

import (
	"context"
	"net/http"
	"time"
)

// ZabbixAuthorizationData результат авторизации
type ZabbixAuthorizationData struct {
	Error   map[string]interface{} `json:"error"`
	JsonRPC string                 `json:"jsonrpc"`
	Result  string                 `json:"result"`
	Id      int                    `json:"id"`
}

type ZabbixAuthorizationErrorMessage struct {
	Data    string `json:"data"`
	Message string `json:"message"`
}

// SettingsZabbixConnection настройки Zabbix соединения
type SettingsZabbixConnection struct {
	Host              string         //ip адрес или доменное имя
	NetProto          string         //сетевой протокол (по умолчанию используется tcp)
	ZabbixHost        string         //имя Zabbix хоста
	ConnectionTimeout *time.Duration //время ожидания подключения (по умолчанию используется 5 сек)
	Port              int            //сетевой порт
}

// SettingsZabbixConnectionJsonRPC настройки Zabbix соединения
// Host - ip адрес или доменное имя
// ConnectionTimeout - время ожидания подключения (по умолчанию используется 5 сек)
type SettingsZabbixConnectionJsonRPC struct {
	ConnectionTimeout time.Duration
	Host              string
	Login             string
	Passwd            string
}

type HandlerZabbixConnection struct {
	ctx         context.Context
	connTimeout time.Duration
	host        string
	netProto    string
	zabbixHost  string
	chanErr     chan error
	port        int
}

type ZabbixConnectionJsonRPC struct {
	connClient        *http.Client
	url               string
	host              string
	login             string
	passwd            string
	applicationType   string
	authorizationHash string
}

type ZabbixOptions struct {
	EventTypes []EventType `yaml:"eventType"`
	ZabbixHost string      `yaml:"zabbixHost"`
}

type EventType struct {
	EventType  string    `yaml:"eventType"`
	ZabbixKey  string    `yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
	IsTransmit bool      `yaml:"isTransmit"`
}

type Handshake struct {
	Message      string `yaml:"message"`
	TimeInterval int    `yaml:"timeInterval"`
}

type MessageSettings struct {
	Message, EventType string
}

type PatternZabbix struct {
	Data    []DataZabbix `json:"data"`
	Request string       `json:"request"`
}

type DataZabbix struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequiestSensorInfo struct {
	zabbixConnection *ZabbixConnectionJsonRPC
	specialId        string
}

type responseData struct {
	Result []map[string]interface{} `json:"result"`
	Error  map[string]interface{}   `json:"error"`
}

type FullSensorInformationFromZabbixAPI struct {
	SensorId   string //id  сенсора
	HostId     string //id хоста
	GeoCode    string //геокод
	ObjectArea string //сфера деятельности
	SubjectRF  string //субъект РФ
	INN        string //ИНН
	HomeNet    string //список домашних сетей
}
