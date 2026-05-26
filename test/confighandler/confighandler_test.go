package confighandler_test

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/enricher_sensor_information/constants"
	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
)

var (
	conf *confighandler.ConfigApp

	err error
)

func TestMain(m *testing.M) {
	os.Unsetenv("GO_ENRICHERSENSORINFO_MAIN")

	//Подключение к NATS
	os.Unsetenv("GO_ENRICHERSENSORINFO_NHOST")
	os.Unsetenv("GO_ENRICHERSENSORINFO_NPORT")
	os.Unsetenv("GO_ENRICHERSENSORINFO_NSUBSC")
	os.Unsetenv("GO_ENRICHERSENSORINFO_NCACHETTL")

	//Подключение к БД с информацией о сенсорах
	os.Unsetenv("GO_ENRICHERSENSORINFO_SIHOST")
	os.Unsetenv("GO_ENRICHERSENSORINFO_SIUSER")
	os.Unsetenv("GO_ENRICHERSENSORINFO_SIPASSWD")
	os.Unsetenv("GO_ENRICHERSENSORINFO_SIRTIMEOUT")
	os.Unsetenv("GO_ENRICHERSENSORINFO_SINCIRCCURL")
	os.Unsetenv("GO_ENRICHERSENSORINFO_SINCIRCCTOKEN")

	//Настройки доступа к БД в которую будут записыватся логи
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGHOST")
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPORT")
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGNAME")
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGUSER")
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD")
	os.Unsetenv("GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME")

	//загружаем ключи и пароли
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	os.Setenv("GO_ENRICHERSENSORINFO_MAIN", "development")

	conf, err = confighandler.New(constants.Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestConfigHandler(t *testing.T) {
	t.Run("Тест чтения конфигурационного файла", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки NATS из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetNATS().Host, "192.168.9.208")
			assert.Equal(t, conf.GetNATS().Port, 4222)
			assert.Equal(t, conf.GetNATS().CacheTTL, 3600)
			assert.Equal(t, conf.GetNATS().Subscription, "object.sensor-info-request.test")
		})

		t.Run("Тест 2. Проверка настройки SensorInformationDataBase из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetSensorInformationDB().ZabbixHost, "192.168.9.45")
			assert.Equal(t, conf.GetSensorInformationDB().ZabbixUser, "803.p.vishnitsky@avz-center.ru")
			assert.Equal(t, conf.GetSensorInformationDB().NetboxHost, "netbox.cloud.gcm")
			assert.Equal(t, conf.GetSensorInformationDB().NetboxPort, 8005)
			assert.Equal(t, conf.GetSensorInformationDB().NCIRCCURL, "https://10.0.227.10/api/v2/companies")
			assert.Equal(t, conf.GetSensorInformationDB().RequestTimeout, 7)
		})

		t.Run("Тест 3. Проверка настройки WriteLogDataBase из файла config_dev.yml", func(t *testing.T) {
			assert.Equal(t, conf.GetLogDB().Host, "datahook.cloud.gcm")
			assert.Equal(t, conf.GetLogDB().Port, 9200)
			assert.Equal(t, conf.GetLogDB().User, "log_writer")
			assert.Equal(t, conf.GetLogDB().Passwd, os.Getenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD"))
			assert.Equal(t, conf.GetLogDB().NameDB, "")
			assert.Equal(t, conf.GetLogDB().StorageNameDB, "enricher_geoip")
		})

		t.Run("Тест 4. Проверка настройки сервера отладки", func(t *testing.T) {
			assert.True(t, conf.GetDebugServer().Enable)
			assert.Equal(t, conf.GetDebugServer().Host, "localhost")
			assert.Equal(t, conf.GetDebugServer().Port, 6262)
		})
	})

	t.Run("Тест чтения переменных окружения", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки NATS", func(t *testing.T) {
			os.Setenv("GO_ENRICHERSENSORINFO_NHOST", "127.0.0.1")
			os.Setenv("GO_ENRICHERSENSORINFO_NPORT", "4242")
			os.Setenv("GO_ENRICHERSENSORINFO_NCACHETTL", "650")
			os.Setenv("GO_ENRICHERSENSORINFO_NSUBSC", "obj.subscript.test_request")

			conf, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetNATS().Host, "127.0.0.1")
			assert.Equal(t, conf.GetNATS().Port, 4242)
			assert.Equal(t, conf.GetNATS().CacheTTL, 650)
			assert.Equal(t, conf.GetNATS().Subscription, "obj.subscript.test_request")
		})

		t.Run("Тест 2. Проверка настройки базы данных с информацией о сенсорах", func(t *testing.T) {
			zhost := "127.0.0.1"
			zuser := "CherryTiggo"
			nhost := "netbox.cloudhost"
			nport := "7562"
			ncirccurl := "https://example.io/api/v2/companies"
			nbtoken := "yydooosmmmskfkfkflddlfgj"
			zpasswd := "SomE_oLd_pasSw"
			ncircctoken := "fa932bca82"
			rtimeout := "10"

			os.Setenv("GO_ENRICHERSENSORINFO_ZHOST", zhost)
			os.Setenv("GO_ENRICHERSENSORINFO_ZUSER", zuser)
			os.Setenv("GO_ENRICHERSENSORINFO_NBHOST", nhost)
			os.Setenv("GO_ENRICHERSENSORINFO_NBPORT", nport)
			os.Setenv("GO_ENRICHERSENSORINFO_NCIRCCURL", ncirccurl)
			os.Setenv("GO_ENRICHERSENSORINFO_NBTOKEN", nbtoken)
			os.Setenv("GO_ENRICHERSENSORINFO_ZPASSWD", zpasswd)
			os.Setenv("GO_ENRICHERSENSORINFO_NCIRCCTOKEN", ncircctoken)
			os.Setenv("GO_ENRICHERSENSORINFO_RTIMEOUT", rtimeout)
			conf, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			netboxPort, _ := strconv.Atoi(nport)
			requestTimeout, _ := strconv.Atoi(rtimeout)

			assert.Equal(t, conf.GetSensorInformationDB().ZabbixHost, zhost)
			assert.Equal(t, conf.GetSensorInformationDB().ZabbixUser, zuser)
			assert.Equal(t, conf.GetSensorInformationDB().ZabbixPasswd, zpasswd)
			assert.Equal(t, conf.GetSensorInformationDB().NCIRCCURL, ncirccurl)
			assert.Equal(t, conf.GetSensorInformationDB().NCIRCCToken, ncircctoken)
			assert.Equal(t, conf.GetSensorInformationDB().NetboxHost, nhost)
			assert.Equal(t, conf.GetSensorInformationDB().NetboxPort, netboxPort)
			assert.Equal(t, conf.GetSensorInformationDB().NetboxToken, nbtoken)
			assert.Equal(t, conf.GetSensorInformationDB().RequestTimeout, requestTimeout)
		})

		t.Run("Тест 3. Проверка настройки WriteLogDataBase", func(t *testing.T) {
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGHOST", "domaniname.database.cm")
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGPORT", "8989")
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGUSER", "somebody_user")
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGNAME", "any_name_db")
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD", "your_passwd")
			os.Setenv("GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME", "log_storage")

			conf, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetLogDB().Host, "domaniname.database.cm")
			assert.Equal(t, conf.GetLogDB().Port, 8989)
			assert.Equal(t, conf.GetLogDB().User, "somebody_user")
			assert.Equal(t, conf.GetLogDB().Passwd, os.Getenv("GO_ENRICHERSENSORINFO_DBWLOGPASSWD"))
			assert.Equal(t, conf.GetLogDB().NameDB, "any_name_db")
			assert.Equal(t, conf.GetLogDB().StorageNameDB, "log_storage")
		})
	})
}
