package interfaces

import (
	"context"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
)

//**************** счётчик *****************

type Counter interface {
	SendMessage(string, int)
}

//************** логирование ***************

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

type WriterLoggingData interface {
	Write(typeLogFile, str string) bool
}

// ************** запрос ***************
type Requester interface {
	CommonTransmitter
	GetData() []byte
	SetData([]byte)
}

//************** ответ ***************

type Responser interface {
	CommonTransmitter
	GetData() []byte
	SetData([]byte)
}

type CommonTransmitter interface {
	GetId() string
	SetId(string)
}

//********** поиск информации ***********

type Searcher interface {
	Search(context.Context, string) (responses.DetailedInformation, error)
}
