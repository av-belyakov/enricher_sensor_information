package netboxinteraction

import (
	"fmt"
	"log"
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
			listId          []int    = []int{}
			searchSensorsId []string = []string{
				"220052", // Связь
				"220063", // Образование
			}
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

			var foundCountSteps int
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
						if index := slices.IndexFunc(searchSensorsId, func(sensorId string) bool {
							//fmt.Println("searchSensorsId:", searchSensorsId)
							//fmt.Println("device.Name:", device.Name)

							if strings.Contains(device.Name, sensorId) {
								fmt.Println("device.Name", device.Name, " ==", sensorId, " sensorId")

								return true
							}

							return false

							//return strings.Contains(device.Name, sensorId)
						}); index != -1 {
							fmt.Println("found index:", index)
							fmt.Println("Key:", key)

							listId = append(listId, devices.Results[key].Id)
						}
					}

					if len(listId) == len(searchSensorsId) {
						break
					}
				}
			}

			fmt.Println("Все устройства найденны за", foundCountSteps, " попыток")
			fmt.Println("Список id устройств:", listId)
		})

		t.Run("Тест 1.3. Получить группы арендаторов каждого устройства", func(t *testing.T) {
			for _, id := range listId {
				tenantsGroup, statusCode, err := nbClient.GetTenantGroups(t.Context(), id)
				assert.NoError(t, err)
				assert.Equal(t, statusCode, http.StatusOK)
				assert.NotEmpty(t, tenantsGroup)

				fmt.Printf("device id=%d tenants group=%v\n", id, tenantsGroup)
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
