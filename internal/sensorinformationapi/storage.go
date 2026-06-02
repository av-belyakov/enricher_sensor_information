package sensorinformationapi

import (
	"reflect"
	"sync"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
)

type InformationStorage struct {
	mtx     sync.RWMutex
	storage []responses.DetailedInformation
}

// NewInformationStorage
func NewInformationStorage(sensorsId []string) *InformationStorage {
	is := InformationStorage{
		storage: make([]responses.DetailedInformation, 0, len(sensorsId)),
	}

	for _, sensorId := range sensorsId {
		is.storage = append(is.storage, responses.DetailedInformation{SensorId: sensorId})
	}

	return &is
}

// Add добавление информации для существующего сенсора
func (s *InformationStorage) Add(sensorInfo responses.DetailedInformation) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for k, v := range s.storage {
		if v.SensorId != sensorInfo.SensorId {
			continue
		}

		s.storage[k] = concatInformation(v, sensorInfo)
	}
}

// GetList список информации
func (s *InformationStorage) GetList() []responses.DetailedInformation {
	return s.storage
}

func concatInformation(mainInfo, inputInfo responses.DetailedInformation) responses.DetailedInformation {
	if mainInfo.SensorId != inputInfo.SensorId {
		return mainInfo
	}

	listFields := []string{
		"INN",
		"GeoCode",
		"HomeNet",
		"OrgName",
		"SensorId",
		"ObjectArea",
		"FullOrgName",
		"SpecialSensorId",
		"SubjectRussianFederation",
		"NetboxTenantGroup",
		"Error",
	}

	valueInputInfo := reflect.ValueOf(inputInfo)
	valueMainInfo := reflect.ValueOf(&mainInfo).Elem()

	for _, field := range listFields {
		vMain := valueMainInfo.FieldByName(field)
		vInput := valueInputInfo.FieldByName(field)

		if !vMain.IsValid() || !vInput.IsValid() {
			continue
		}

		// Проверяем, что оба поля строковые
		if vMain.Kind() != reflect.String || vInput.Kind() != reflect.String {
			continue
		}

		if valueMainInfo.CanSet() && vInput.String() != "" && field != "SensorId" {
			if field == "Error" && vInput.String() != "" {
				vMain.SetString(valueMainInfo.FieldByName(field).String() + " " + vInput.String())

				continue
			}

			vMain.SetString(vInput.String())
		}
	}

	return mainInfo
}
