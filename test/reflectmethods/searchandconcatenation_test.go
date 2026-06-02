package reflectmethods

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
)

func TestSearchAndConcatenation(t *testing.T) {
	sId := gofakeit.UUID()

	mainInfo := responses.DetailedInformation{
		INN:                      gofakeit.SSN(),
		GeoCode:                  gofakeit.Country(),
		HomeNet:                  gofakeit.IPv4Address(),
		OrgName:                  gofakeit.Company(),
		SensorId:                 sId,
		ObjectArea:               gofakeit.City(),
		FullOrgName:              gofakeit.Company(),
		SpecialSensorId:          "",
		SubjectRussianFederation: "",
		NetboxTenantGroup:        "",
		Error:                    gofakeit.JobDescriptor(),
	}

	inputInfo := responses.DetailedInformation{
		INN:                      "",
		GeoCode:                  "",
		HomeNet:                  "",
		OrgName:                  "",
		SensorId:                 sId,
		ObjectArea:               "",
		FullOrgName:              "",
		SpecialSensorId:          gofakeit.UUID(),
		SubjectRussianFederation: gofakeit.PronounObject(),
		NetboxTenantGroup:        gofakeit.AdjectiveQuantitative(),
		Error:                    gofakeit.JobDescriptor(),
	}

	result := concatInformation(mainInfo, inputInfo)
	fmt.Printf("Result:\n%+v\n", result)

	assert.NotEmpty(t, result.INN)
	assert.NotEmpty(t, result.GeoCode)
	assert.NotEmpty(t, result.HomeNet)
	assert.NotEmpty(t, result.OrgName)
	assert.NotEmpty(t, result.SensorId)
	assert.NotEmpty(t, result.ObjectArea)
	assert.NotEmpty(t, result.FullOrgName)
	assert.NotEmpty(t, result.SpecialSensorId)
	assert.NotEmpty(t, result.SubjectRussianFederation)
	assert.NotEmpty(t, result.NetboxTenantGroup)
	assert.NotEmpty(t, result.Error)
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
			if field == "Error" {
				vMain.SetString(valueMainInfo.FieldByName(field).String() + " " + vInput.String())

				continue
			}

			vMain.SetString(vInput.String())
		}
	}

	return mainInfo
}
