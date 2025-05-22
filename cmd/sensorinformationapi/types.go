package sensorinformationapi

// SensorInformationClient клиент для взаимодействия с API
type SensorInformationClient struct {
	Settings SensorInformationSettings
	ChInput  chan any
}

// SensorInformationSettings настройки модуля
type SensorInformationSettings struct {
	Host           string
	User           string
	Passwd         string
	NCIRCCURL      string
	NCIRCCToken    string
	Port           int
	RequestTimeout int
}

// sensorInformationClientOptions функциональные параметры
type sensorInformationClientOptions func(*SensorInformationClient) error
