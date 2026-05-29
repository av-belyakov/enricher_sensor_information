package sensorinformationapi

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/enricher_sensor_information/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
	"github.com/av-belyakov/enricher_sensor_information/internal/sensorinformationapi"
)

func TestGetFullSensorInformation(t *testing.T) {
	searchSensorsId := []string{"220065", "308051", "310067", "530013", "570027", "630019", "630062", "8030015"}

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("error loading .env file: %v", err)
	}

	cfg, err := confighandler.New(constants.Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	siClient, err := sensorinformationapi.New(
		sensorinformationapi.WithZabbixHost(cfg.GetSensorInformationDB().ZabbixHost),
		sensorinformationapi.WithZabbixUser(cfg.GetSensorInformationDB().ZabbixUser),
		sensorinformationapi.WithZabbixPasswd(cfg.GetSensorInformationDB().ZabbixPasswd),
		sensorinformationapi.WithNetboxHost(cfg.GetSensorInformationDB().NetboxHost),
		sensorinformationapi.WithNetboxPort(cfg.GetSensorInformationDB().NetboxPort),
		sensorinformationapi.WithNetboxToken(cfg.GetSensorInformationDB().NetboxToken),
		sensorinformationapi.WithNCIRCCURL(cfg.GetSensorInformationDB().NCIRCCURL),
		sensorinformationapi.WithNCIRCCToken(cfg.GetSensorInformationDB().NCIRCCToken),
		sensorinformationapi.WithRequestTimeout(cfg.GetSensorInformationDB().RequestTimeout),
	)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := siClient.Search(t.Context(), searchSensorsId)
	assert.NoError(t, err)
	assert.Greater(t, len(result), 0)

	t.Cleanup(func() {
		os.Unsetenv("GO_ENRICHERSENSORINFO_ZPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_NBTOKEN")
		os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_NCIRCCTOKEN")
	})
}
