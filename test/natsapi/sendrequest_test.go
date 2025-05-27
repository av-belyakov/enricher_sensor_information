package natsapi_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	NATS_HOST = "192.168.9.208"
	NATS_PORT = 4222

	CACHETTL     = 360
	SUBSCRIPTION = "object.sensor-info-request.test"
)

type ResponseData struct {
	Information []SensorInformation `json:"found_information"`
	Source      string              `json:"source"`
	TaskId      string              `json:"task_id"`
	Error       string              `json:"error"`
}

type SensorInformation struct {
	INN                      string `json:"inn"`
	GeoCode                  string `json:"geo_code"`
	HomeNet                  string `json:"home_net"`
	SensorId                 string `json:"sensor_id"`
	ObjectArea               string `json:"object_area"`
	SpecialSensorId          string `json:"special_sensor_id"`
	OrganizationName         string `json:"organization_name"`
	FullOrganizationName     string `json:"full_organization_name"`
	SubjectRussianFederation string `json:"subject_russian_federation"`
	Error                    string `json:"error"`
}

func CreateNatsConnect(host string, port int) (*nats.Conn, error) {
	var (
		nc  *nats.Conn
		err error
	)

	nc, err = nats.Connect(
		fmt.Sprintf("%s:%d", host, port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return nc, err
	}

	fmt.Println("func 'CreateNatsConnect', START")

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		if err != nil {
			fmt.Println(err)
			fmt.Printf("func 'CreateNatsConnect' the connection with NATS has been disconnected %s\n", err.Error())

			return
		}

		fmt.Println("func 'CreateNatsConnect' the connection with NATS has been disconnected")
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		if err != nil {
			fmt.Printf("func 'CreateNatsConnect' the connection to NATS has been re-established %s\n", err.Error())

			return
		}

		fmt.Println("func 'CreateNatsConnect' the connection to NATS has been re-established")
	})

	return nc, nil
}

func TestGetSensorCommonInfo(t *testing.T) {
	//инициализация с NATS
	nc, err := CreateNatsConnect(NATS_HOST, NATS_PORT)
	if err != nil {
		log.Fatalln(err)
	}

	//запрос информации через NATS
	msg, err := nc.RequestWithContext(t.Context(), SUBSCRIPTION, []byte(`{
			"source": "test_source",
	  		"task_id": "41af7c2b34",
	   		"list_sensors": ["8030073", "8030141", "8030017"]
		}`))

	assert.NotNil(t, msg)

	res := ResponseData{}
	if msg != nil {
		err = json.Unmarshal(msg.Data, &res)
		assert.NoError(t, err)
		assert.Empty(t, res.Error)
	}
	// обработка ответа
	t.Log("Response:")
	for k, v := range res.Information {
		t.Logf("%d. %s\n", k, v)
	}

	t.Cleanup(func() {
		nc.Close()
	})
}
