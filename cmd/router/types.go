package router

import (
	"github.com/av-belyakov/enricher_sensor_information/interfaces"
)

type Router struct {
	counter       interfaces.Counter
	logger        interfaces.Logger
	commonInfo    interfaces.Searcher
	chFromNatsApi <-chan interfaces.Requester
	chToNatsApi   chan<- interfaces.Responser
}

type RouterSettings struct {
	SearchCommonInfo interfaces.Searcher
	ChanFromNatsApi  <-chan interfaces.Requester
	ChanToNatsApi    chan<- interfaces.Responser
}
