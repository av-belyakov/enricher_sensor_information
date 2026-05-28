package sensorinformationapi

import (
	"context"
	"regexp"

	"github.com/goforj/godump"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// Search поиск информации о сенсоре
func (api *SensorInformationClient) Search(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	//авторизуемся в Zabbix
	if err := api.zabbixConn.Authorization(ctx); err != nil {
		return nil, err
	}

	results := make([]responses.DetailedInformation, 0, len(sensorsId))
	for _, sensorId := range sensorsId {
		//поиск основной информации по сенсору в Zabbix
		info, err := zabbixinteractions.GetFullSensorInformation(ctx, sensorId, api.zabbixConn)
		if err != nil {
			info.Error = err.Error()

			continue
		}

		reg, err := regexp.Compile(`^[0-9]+$`)
		if err != nil {
			info.Error = err.Error()

			continue
		}

		//поиск подробной информации об организации по её ИНН в НКЦКИ
		if reg.MatchString(info.INN) {
			innInfo, err := api.ncirccConn.GetFullNameOrganizationByINN(ctx, info.INN)
			if err != nil {
				info.Error = err.Error()

				continue
			}

			if innInfo.Count == 0 {
				info.Error = "inn information was not found"

				continue
			}

			info.OrgName = innInfo.Data[0].Name
			info.FullOrgName = innInfo.Data[0].Sname
		}

		godump.Dump(info)

		results = append(results, info)
	}

	return results, nil
}
