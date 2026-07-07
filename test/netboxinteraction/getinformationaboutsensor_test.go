package netboxinteraction

import (
	"fmt"
	"log"
	"maps"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/av-belyakov/enricher_sensor_information/internal/netboxinteractions"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	Host = "netbox.cloud.gcm"
	Port = 8005
)

func TestGetInformationAboutSensor(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error loading .env file: %v", err)
	}

	nbClient, err := netboxinteractions.New(
		os.Getenv("GO_ENRICHERSENSORINFO_NBTOKEN"),
		netboxinteractions.WithHost(Host),
		netboxinteractions.WithPort(Port),
		netboxinteractions.WithTimeout(10),
	)
	if err != nil {
		t.Fatalf("error creating Netbox client: %v", err)
	}

	t.Run("Тест 1. Получить информацию о сенсоре", func(t *testing.T) {
		var (
			totalDevices int
			countSteps   int = 1
			devicesLimit int = 350

			// список id устройств
			sensorIds map[string]int = map[string]int{
				"220052":  0, // Связь
				"220063":  0, // Образование
				"308047":  0, // Оборонная промышленность
				"310053":  0, // Ракетно-космическая промышленность
				"430019":  0, // Транспорт
				"570050":  0, // Образование
				"570057":  0, // Оборонная промышленность
				"8030154": 0, // Государственная/муниципальная власть
			}
			searchSensorsId []string
			//searchSensorsId []string = []string{"220065", "308051", "310067", "530013", "570027", "630019", "630062", "8030015"}
		)

		t.Run("Тест 1.1. Получить общее количество устройств", func(t *testing.T) {
			countDevices, statusCode, err := nbClient.GetCountDevices(t.Context())
			assert.NoError(t, err)
			assert.Equal(t, statusCode, http.StatusOK)
			assert.Greater(t, countDevices, 0)

			totalDevices = countDevices

			if totalDevices < devicesLimit {
				countSteps = 1
			} else {
				countSteps = int(math.Ceil(float64(totalDevices) / float64(devicesLimit)))
			}

			fmt.Println("Count devices:", countDevices)
			fmt.Printf("count steps=%d if devices limit=%d\n", countSteps, devicesLimit)
		})

		t.Run("Тест 1.2. Поиск сенсоров", func(t *testing.T) {
			if totalDevices == 0 {
				log.Fatal("the device list cannot be empty")
			}

			for sensorId := range maps.Keys(sensorIds) {
				searchSensorsId = append(searchSensorsId, sensorId)
			}

			var foundCountSteps, countIds int
			for step := range countSteps {
				foundCountSteps++

				devices, statusCode, err := nbClient.GetDevicesLimitInformation(t.Context(), devicesLimit, step*devicesLimit)
				if err != nil {
					log.Fatal(err)
				}
				assert.Equal(t, statusCode, http.StatusOK)

				//fmt.Println("Devices:", devices)

				/*
					220052 - id 965
					220063 - id 1600
				*/

				if statusCode == http.StatusOK {
					// было бы лучше класть результат в карту где ключем является name устройства, но к сожалению имена
					// устройств могут не точно соответствовать искомому сенсору, например '570027 (48832465)'
					// поэтому осуществляется поиск в срезе

					for key, device := range devices.Results {
						//if index
						_ = slices.IndexFunc(searchSensorsId, func(sensorId string) bool {
							if strings.Contains(device.Name, sensorId) {
								fmt.Println("device.Name", device.Name, " ==", sensorId, " sensorId")

								if _, ok := sensorIds[sensorId]; ok {
									sensorIds[sensorId] = devices.Results[key].Id

									countIds++
								}

								return true
							}

							return false
						})
						//; index != -1 {
						//	fmt.Println("found index:", index)
						//	fmt.Println("Key:", key)
						//}
					}

					if countIds == len(searchSensorsId) {
						break
					}
				}
			}

			fmt.Println("Все устройства найденны за", foundCountSteps, "попыток")
			fmt.Println("Список id устройств:", sensorIds)
		})

		t.Run("Тест 1.3. Получить группы арендаторов каждого устройства", func(t *testing.T) {
			for sensorId, deviceId := range sensorIds {
				tenantsGroup, statusCode, err := nbClient.GetTenantGroups(t.Context(), deviceId)
				assert.NoError(t, err)
				assert.Equal(t, statusCode, http.StatusOK)
				assert.NotEmpty(t, tenantsGroup)

				fmt.Printf("sensor id=%s device id=%d tenants group=%v\n", sensorId, deviceId, tenantsGroup)
			}
		})
		/*
			res, statusCode, err := nbClient.Get(t.Context(), "/api/dcim/devices/?fields=id,name&limit=150&offset=0")
			assert.NoError(t, err)
			assert.Equal(t, statusCode, http.StatusOK)

			fmt.Println("--- Result:")
			fmt.Println(string(res))
		*/
	})
}
