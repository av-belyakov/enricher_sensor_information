package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/natsapi/storagetemporary"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	counter              interfaces.Counter
	logger               interfaces.Logger
	natsConn             *nats.Conn
	storage              *storagetemporary.StorageTemporary
	subscriptionRequest  string
	subscriptionResponse string
	settings             apiNatsSettings
	chFromModule         chan interfaces.Requester
	chToModule           chan interfaces.Responser
}

type apiNatsSettings struct {
	nameRegionalObject string
	command            string
	host               string
	cachettl           int
	port               int
}

// NatsApiOptions функциональные опции
type NatsApiOptions func(*apiNatsModule) error

// ObjectBeingTransferred объект для передачи данных
type ObjectBeingTransferred struct {
	Data []byte
	Id   string
}
