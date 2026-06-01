package sensorinformationapi

import (
	"context"
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

	var g errgroup.Group
	// поиск основной информации в Zabbix и НКЦКИ
	g.Go(func() error {
		// авторизуемся в Zabbix
		if err := api.zabbixConn.AuthorizationStart(ctx); err != nil {
			return err
		}

		result, err := api.SearchCommonInformation(ctx, sensorsId)
		for _, v := range result {
			storage.Add(v)
		}

		println("-0000--- method 'SearchCommonInformation', Error:", err)

		return err
	})
	// поиск дополнитольной информации в Netbox
	g.Go(func() error {
		result, err := api.SearchAdditionalInformation(ctx, sensorsId)
		for _, v := range result {
			storage.Add(v)
		}

		println("-1111--- method 'SearchAdditionalInformation', Error:", err)

		return err
	})

	if err := g.Wait(); err != nil {
		return storage.GetList(), err
	}

	return storage.GetList(), nil
}

// SearchCommonInformation поиск основной информации о сенсоре
func (api *SensorInformationClient) SearchCommonInformation(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	var (
		response []responses.DetailedInformation = make([]responses.DetailedInformation, 0, len(sensorsId))
	)

	for _, sensorId := range sensorsId {
		select {
		case <-ctx.Done():
			return response, ctx.Err()

		default:
			println("method 'SearchCommonInformation', sensorId:", sensorId, " get full sensor information")

			res := responses.DetailedInformation{
				SensorId: sensorId,
			}

			//поиск основной информации по сенсору в Zabbix
			info, err := zabbixinteractions.GetFullSensorInformation(ctx, sensorId, api.zabbixConn)
			if err != nil {
				res.Error = err.Error()
				response = append(response, res)

				continue
			}

			res.INN = info.INN
			res.GeoCode = info.GeoCode
			res.HomeNet = info.HomeNet
			res.SensorId = info.SensorId
			res.ObjectArea = info.ObjectArea
			res.SpecialSensorId = info.SpecialSensorId
			res.SubjectRussianFederation = info.SubjectRussianFederation

			reg, err := regexp.Compile(`^[0-9]+$`)
			if err != nil {
				res.Error = err.Error()
				response = append(response, res)

				continue
			}

			//поиск подробной информации об организации по её ИНН в НКЦКИ
			if reg.MatchString(res.INN) {
				innInfo, err := api.ncirccConn.GetFullNameOrganizationByINN(ctx, info.INN)
				if err != nil {
					res.Error = err.Error()
					response = append(response, res)

					continue
				}

				if innInfo.Count == 0 {
					res.Error = "inn information was not found"
					response = append(response, res)

					continue
				}

				res.OrgName = innInfo.Data[0].Name
				res.FullOrgName = innInfo.Data[0].Sname
			}

			response = append(response, res)
		}
	}

	return response, nil
}

// SearchAdditionalInformation поиск информации в Netbox
func (api *SensorInformationClient) SearchAdditionalInformation(ctx context.Context, sensorsId []string) ([]responses.DetailedInformation, error) {
	var (
		response   []responses.DetailedInformation = make([]responses.DetailedInformation, 0, len(sensorsId))
		sensors    map[string]int                  = make(map[string]int, len(sensorsId))
		countSteps int
	)

	println("method 'SearchAdditionalInformation', sensorsId:", sensorsId)

	countDevices, _, err := api.netboxConn.GetCountDevices(ctx)
	if err != nil {
		println("method 'SearchAdditionalInformation', 111 Error:", err)

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
			println("method 'SearchAdditionalInformation', 222 Error:", err)

			return response, ctx.Err()

		default:
			println("method 'SearchAdditionalInformation', step:", step)

			// получаем ограниченную информацию об устройствах, что бы получить внутренний id устройства
			// который понадобится для запроса дополнительной информации об группе арендаторов
			devices, statusCode, err := api.netboxConn.GetDevicesLimitInformation(ctx, constants.Devices_Limit, step*constants.Devices_Limit)
			if err != nil {
				println("method 'SearchAdditionalInformation', 333 Error:", err)

				return response, err
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
			println("method 'SearchAdditionalInformation', 444 Error:", err)

			return response, ctx.Err()

		default:
			println("method 'SearchAdditionalInformation', sensorId:", sensorId)

			tenantsGroup, _, err := api.netboxConn.GetTenantGroups(ctx, internalId)
			if err != nil {
				response = append(response, responses.DetailedInformation{
					SensorId: sensorId,
					Error:    err.Error(),
				})

				continue
			}

			response = append(response, responses.DetailedInformation{
				SensorId:          sensorId,
				NetboxTenantGroup: tenantsGroup,
			})
		}
	}

	println("method 'SearchAdditionalInformation', Final Error:", err)

	return response, nil
}
