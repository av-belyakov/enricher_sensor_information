package sensorinformationapi

import (
	"context"
	"regexp"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// Search поиск информации о сенсоре
func (api *SensorInformationClient) Search(ctx context.Context, sensorId string) (responses.DetailedInformation, error) {
	//авторизуемся в Zabbix
	api.zabbixConn.Authorization(ctx)

	//fmt.Println("func 'SensorInformationClient.Search', search sensor with id:", sensorId)
	//fmt.Println("func 'SensorInformationClient.Search', поиск основной информации по сенсору в Zabbix")

	//поиск основной информации по сенсору в Zabbix
	commonInfo, err := zabbixinteractions.GetFullSensorInformation(ctx, sensorId, api.zabbixConn)
	if err != nil {
		return commonInfo, err
	}

	commonInfo.SensorId = sensorId

	reg, err := regexp.Compile(`^[0-9]+$`)
	if err != nil {
		return commonInfo, err
	}

	//поиск подробной информации об организации по её ИНН в НКЦКИ
	if reg.MatchString(commonInfo.INN) {
		innInfo, err := api.ncirccConn.GetFullNameOrganizationByINN(ctx, commonInfo.INN)
		if err != nil {
			//fmt.Println("func 'SensorInformationClient.Search', ERROR:", err)

			return commonInfo, err
		}

		if innInfo.Count == 0 {
			return commonInfo, err
		}

		commonInfo.OrgName = innInfo.Data[0].Name
		commonInfo.FullOrgName = innInfo.Data[0].Sname
	}

	//godump.Dump(commonInfo)

	return commonInfo, err
}
