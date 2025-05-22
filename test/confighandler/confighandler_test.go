package confighandler_test

import (
	"log"
	"os"
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
	os.Unsetenv("GO_ENRICHERSENSORINFO_SIPORT")
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
			assert.Equal(t, conf.GetSensorInformationDB().Host, "192.168.9.45")
			assert.Equal(t, conf.GetSensorInformationDB().Port, 13013)
			assert.Equal(t, conf.GetSensorInformationDB().User, "Cherry")
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

		t.Run("Тест 2. Проверка настройки GeoIPDataBase", func(t *testing.T) {
			os.Setenv("GO_ENRICHERSENSORINFO_SIHOST", "127.0.0.1")
			os.Setenv("GO_ENRICHERSENSORINFO_SIPORT", "13813")
			os.Setenv("GO_ENRICHERSENSORINFO_SIUSER", "CherryTiggo")
			os.Setenv("GO_ENRICHERSENSORINFO_SIPASSWD", "SomE_oLd_pasSw")
			os.Setenv("GO_ENRICHERSENSORINFO_SIRTIMEOUT", "10")
			os.Setenv("GO_ENRICHERSENSORINFO_SINCIRCCURL", "https://example.io/api/v2/companies")
			os.Setenv("GO_ENRICHERSENSORINFO_SINCIRCCTOKEN", "fa932bca82")
			conf, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, conf.GetSensorInformationDB().Host, "127.0.0.1")
			assert.Equal(t, conf.GetSensorInformationDB().Port, 13813)
			assert.Equal(t, conf.GetSensorInformationDB().User, "CherryTiggo")
			assert.Equal(t, conf.GetSensorInformationDB().Passwd, "SomE_oLd_pasSw")
			assert.Equal(t, conf.GetSensorInformationDB().NCIRCCURL, "https://example.io/api/v2/companies")
			assert.Equal(t, conf.GetSensorInformationDB().NCIRCCToken, "fa932bca82")
			assert.Equal(t, conf.GetSensorInformationDB().RequestTimeout, 10)
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
