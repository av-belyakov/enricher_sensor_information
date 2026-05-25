package dicontainer

import (
	"context"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
)

type Counter interface {
	SendMessage(msgType string, count int)
}

type Logger interface {
	GetChan() <-chan interfaces.Messager
	Send(msgType, message string)
	Close()
}

type SimpleLogger interface {
	SetDataBaseInteraction(dbi simplelogger.DataBaseInteractor)
	GetCountFileDescription() int
	GetListTypeFiles() []string
	Write(typeLog, msg string) bool
}

type Configer interface {
	GetCommon() *confighandler.CfgCommon
	GetNATS() *confighandler.CfgNats
	GetLogDB() *confighandler.CfgWriteLogDB
	GetListLogs() []*confighandler.LogSet
	GetSensorInformationDB() *confighandler.CfgSensorInformationDB
}

type NatsConnecter interface {
	GetChFromModule() chan interfaces.Requester
	GetChToModule() chan interfaces.Responser
}

type DbLogger interface {
	Write(msgType, msg string) error
}

type SensorInformationConnecter interface {
	Search(ctx context.Context, inn string) (responses.DetailedInformation, error)
}
