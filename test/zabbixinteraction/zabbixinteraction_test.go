package zabbixinteraction_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

const Sensor_ID = "8030174"

var (
	zc *zabbixinteractions.ZabbixConnectionJsonRPC

	err error
)

func TestMain(t *testing.M) {
	if err = gotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	zc, err = zabbixinteractions.NewZabbixConnectionJsonRPC(
		zabbixinteractions.SettingsZabbixConnectionJsonRPC{
			ConnectionTimeout: 10 * time.Second,
			Host:              "192.168.9.45",
			Login:             "803.p.vishnitsky@avz-center.ru",
			Passwd:            os.Getenv("GO_ENRICHERSENSORINFO_SIPASSWD"),
		})
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(t.Run())
}

func TestGetZabbixInfo(t *testing.T) {
	t.Run("Тест 1. Авторизация", func(t *testing.T) {
		assert.NoError(t, zc.Authorization(t.Context()))
	})

	t.Run("Тест 2. ", func(t *testing.T) {
		information, err := zabbixinteractions.GetFullSensorInformation(t.Context(), Sensor_ID, zc)
		assert.NoError(t, err)

		t.Logf("Sensor information:\n%+v\n", information)
	})

	t.Cleanup(func() {
		os.Unsetenv("GO_ENRICHERSENSORINFO_SIPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_SINCIRCCTOKEN")
	})
}
