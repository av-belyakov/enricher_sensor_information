package countermessage

import (
	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/informationcountingstorage"
)

// CounterMessage счетчик сообщений
type CounterMessage struct {
	storage  *informationcountingstorage.InformationCountingStorage
	chInput  chan DataCounterMessage
	chOutput chan<- interfaces.Messager
}

// SomeMessage некоторое сообщение
type SomeMessage struct {
	Type, Message string
}

// DataCounterMessage содержит информацию для подсчета
type DataCounterMessage struct {
	msgType string
	count   int
}
