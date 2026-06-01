package zabbixinteraction_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"

	"github.com/av-belyakov/enricher_sensor_information/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
	"github.com/av-belyakov/enricher_sensor_information/internal/zabbixinteractions"
)

const Sensor_ID = "8030174"

var (
	zConn *connectionjsonrpc.ZabbixConnectionJsonRPC
	err   error
)

func TestMain(t *testing.M) {
	os.Setenv("GO_ENRICHERSENSORINFO_MAIN", "development")

	if err = gotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	cfg, err := confighandler.New(constants.Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Config sensor information data base:", cfg.GetSensorInformationDB())

	if cfg.GetSensorInformationDB().ZabbixUseTLS {
		zConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithTLS(),
			connectionjsonrpc.WithInsecureSkipVerify(),
			connectionjsonrpc.WithHost(cfg.GetSensorInformationDB().ZabbixHost),
			connectionjsonrpc.WithPort(cfg.GetSensorInformationDB().ZabbixPort),
			connectionjsonrpc.WithLogin(cfg.GetSensorInformationDB().ZabbixUser),
			connectionjsonrpc.WithPasswd(cfg.GetSensorInformationDB().ZabbixPasswd),
			connectionjsonrpc.WithConnectionTimeout(cfg.GetSensorInformationDB().RequestTimeout),
		)
	} else {
		zConn, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost(cfg.GetSensorInformationDB().ZabbixHost),
			connectionjsonrpc.WithPort(cfg.GetSensorInformationDB().ZabbixPort),
			connectionjsonrpc.WithLogin(cfg.GetSensorInformationDB().ZabbixUser),
			connectionjsonrpc.WithPasswd(cfg.GetSensorInformationDB().ZabbixPasswd),
			connectionjsonrpc.WithConnectionTimeout(cfg.GetSensorInformationDB().RequestTimeout),
		)
	}
	if err != nil {
		log.Fatalln(err)
	}

	if err = zConn.AuthorizationStart(context.TODO()); err != nil {

	}

	os.Exit(t.Run())
}

func TestGetZabbixInfo(t *testing.T) {
	t.Run("Тест 1. Получение информации о сенсорах", func(t *testing.T) {
		information, err := zabbixinteractions.GetFullSensorInformation(t.Context(), Sensor_ID, zConn)
		assert.NoError(t, err)

		t.Logf("Sensor information:\n%+v\n", information)
	})

	t.Cleanup(func() {
		os.Unsetenv("GO_ENRICHERSENSORINFO_MAIN")
		os.Unsetenv("GO_ENRICHERSENSORINFO_ZPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_NBTOKEN")
		os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD")
		os.Unsetenv("GO_ENRICHERSENSORINFO_NCIRCCTOKEN")
	})
}
