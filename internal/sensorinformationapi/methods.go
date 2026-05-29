package sensorinformationapi

import (
	"context"
	"errors"
	"math"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/enricher_sensor_information/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

// Search поиск информации о сенсоре
func (api *SensorInformationClient) Search(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	// инициализируем общее хранилище
	storage := NewInformationStorage(sensorsId)

	// авторизуемся в Zabbix
	if err := api.zabbixConn.Authorization(ctx); err != nil {
		return storage.GetList(), err
	}

	var g errgroup.Group
	// поиск основной информации в Zabbix и НКЦКИ
	g.Go(func() error {
		result, err := api.SearchCommonInformation(ctx, sensorsId)
		for _, v := range result {
			storage.Add(v)
		}

		return err
	})
	// поиск дополнитольной информации в Netbox
	g.Go(func() error {
		result, err := api.SearchAdditionalInformation(ctx, sensorsId)
		for _, v := range result {
			storage.Add(v)
		}

		return err
	})

	err := g.Wait()

	return storage.GetList(), err
}

// SearchCommonInformation поиск основной информации о сенсоре
func (api *SensorInformationClient) SearchCommonInformation(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	var (
		response  []responses.DetailedInformation = make([]responses.DetailedInformation, 0, len(sensorsId))
		errorList error
	)

	for _, sensorId := range sensorsId {
		select {
		case <-ctx.Done():
			return response, ctx.Err()

		default:
			//поиск основной информации по сенсору в Zabbix
			info, err := zabbixinteractions.GetFullSensorInformation(ctx, sensorId, api.zabbixConn)
			if err != nil {
				errorList = errors.Join(errorList, err)

				continue
			}

			reg, err := regexp.Compile(`^[0-9]+$`)
			if err != nil {
				errorList = errors.Join(errorList, err)

				continue
			}

			//поиск подробной информации об организации по её ИНН в НКЦКИ
			if reg.MatchString(info.INN) {
				innInfo, err := api.ncirccConn.GetFullNameOrganizationByINN(ctx, info.INN)
				if err != nil {
					errorList = errors.Join(errorList, err)

					continue
				}

				if innInfo.Count == 0 {
					errorList = errors.Join(errorList, errors.New("inn information was not found"))

					continue
				}

				info.OrgName = innInfo.Data[0].Name
				info.FullOrgName = innInfo.Data[0].Sname
			}

			response = append(response, info)
		}
	}

	return response, errorList
}

// SearchAdditionalInformation поиск информации в Netbox
func (api *SensorInformationClient) SearchAdditionalInformation(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	var (
		response   []responses.DetailedInformation = make([]responses.DetailedInformation, 0, len(sensorsId))
		sensors    map[string]int                  = make(map[string]int, len(sensorsId))
		countSteps int
		errorList  error
	)

	countDevices, _, err := api.netboxConn.GetCountDevices(ctx)
	if err != nil {
		return response, err
	}

	if countDevices < constants.Devices_Limit {
		countSteps = 1
	} else {
		countSteps = int(math.Ceil(float64(countDevices) / float64(constants.Devices_Limit)))
	}

	// поиск устройств с именем соответствующим имени сенсора
	for step := range countSteps {
		select {
		case <-ctx.Done():
			return response, ctx.Err()

		default:
			// получаем ограниченную информацию об устройствах, что бы получить внутренний id устройства
			// который понадобится для запроса дополнительной информации об группе арендаторов
			devices, statusCode, err := api.netboxConn.GetDevicesLimitInformation(ctx, constants.Devices_Limit, step*constants.Devices_Limit)
			if err != nil {
				errorList = errors.Join(errorList, err)

				continue
			}

			if statusCode == http.StatusOK {
				// было бы лучше класть результат в карту где ключем является name устройства, но к сожалению имена
				// устройств могут не точно соответствовать искомому сенсору, например '570027 (48832465)'
				// поэтому осуществляется поиск в срезе

				for _, device := range devices.Results {
					if index := slices.IndexFunc(sensorsId, func(sensorId string) bool {
						return strings.Contains(device.Name, sensorId)
					}); index != -1 {
						sensors[sensorsId[index]] = devices.Results[index].Id
					}
				}

				if len(sensors) == len(sensorsId) {
					break
				}
			}
		}
	}

	// поиск групп арендаторов по внутреннему id устройства
	for sensorId, internalId := range sensors {
		select {
		case <-ctx.Done():
			return response, ctx.Err()

		default:
			tenantsGroup, _, err := api.netboxConn.GetTenantGroups(ctx, internalId)
			if err != nil {
				errorList = errors.Join(errorList, err)

				continue
			}

			response = append(response, responses.DetailedInformation{
				SensorId:          sensorId,
				NetboxTenantGroup: tenantsGroup,
			})
		}
	}

	return response, errorList
}
