package dicontainer

import "github.com/av-belyakov/enricher_sensor_information/interfaces"

// NewDIContainer ленивая инициализация DI контейнера
func NewDIContainer(rootDir string, ch chan interfaces.Messager) *DiContainer {
	return &DiContainer{
		rootDir: rootDir,
		ch:      ch,
	}
}
