package sensorinformationapi

import (
	"context"
	"regexp"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/ncirccinteractions"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// New настраивает новый модуль взаимодействия с API
func New(opts ...sensorInformationClientOptions) (*SensorInformationClient, error) {
	api := &SensorInformationClient{
		settings: SensorInformationSettings{
			requestTimeout: 10,
		},
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	//инициализация соединения с Zabbix
	zConn, err := zabbixinteractions.NewZabbixConnectionJsonRPC(
		zabbixinteractions.SettingsZabbixConnectionJsonRPC{
			Host:              api.settings.host,
			Login:             api.settings.user,
			Passwd:            api.settings.passwd,
			ConnectionTimeout: time.Duration(api.settings.requestTimeout) * time.Second,
		})
	if err != nil {
		return api, err
	}
	api.zabbixConn = zConn

	//инициализация соединения с НКЦКИ
	ncirccConn, err := ncirccinteractions.NewClient(
		api.settings.ncirccURL,
		api.settings.ncirccToken,
		time.Duration(api.settings.requestTimeout)*time.Second,
	)
	if err != nil {
		return api, err
	}
	api.ncirccConn = ncirccConn

	return api, nil
}

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
