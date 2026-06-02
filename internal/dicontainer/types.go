package dicontainer

import "github.com/av-belyakov/enricher_sensor_information/interfaces"

// DiContainer DI контейнер
type DiContainer struct {
	logger       Logger
	counter      Counter
	configer     Configer
	simpleLogger SimpleLogger

	sensorInformationDB SensorInformationConnecter
	dbLogger            DbLogger
	nats                NatsConnecter

	ch      chan interfaces.Messager
	rootDir string
}
