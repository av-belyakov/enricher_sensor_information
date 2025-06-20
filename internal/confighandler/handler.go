// Пакет confighandler формирует конфигурационные настройки приложения
package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

func New(rootDir string) (*ConfigApp, error) {
	conf := &ConfigApp{}

	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_ENRICHERSENSORINFO_MAIN": "",

			//Подключение к NATS
			"GO_ENRICHERSENSORINFO_NHOST":     "",
			"GO_ENRICHERSENSORINFO_NPORT":     "",
			"GO_ENRICHERSENSORINFO_NSUBSC":    "",
			"GO_ENRICHERSENSORINFO_NCACHETTL": "",

			//Подключение к БД с информацией о сенсорах
			"GO_ENRICHERSENSORINFO_SIHOST":        "",
			"GO_ENRICHERSENSORINFO_SIPORT":        "",
			"GO_ENRICHERSENSORINFO_SIUSER":        "",
			"GO_ENRICHERSENSORINFO_SIPASSWD":      "",
			"GO_ENRICHERSENSORINFO_SIRTIMEOUT":    "",
			"GO_ENRICHERSENSORINFO_SINCIRCCURL":   "",
			"GO_ENRICHERSENSORINFO_SINCIRCCTOKEN": "",

			//Настройки доступа к БД в которую будут записыватся логи
			"GO_ENRICHERSENSORINFO_DBWLOGHOST":        "",
			"GO_ENRICHERSENSORINFO_DBWLOGPORT":        "",
			"GO_ENRICHERSENSORINFO_DBWLOGNAME":        "",
			"GO_ENRICHERSENSORINFO_DBWLOGUSER":        "",
			"GO_ENRICHERSENSORINFO_DBWLOGPASSWD":      "",
			"GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME": "",
		}
	)

	getFileName := func(sf, confPath string, lfs []fs.DirEntry) (string, error) {
		for _, v := range lfs {
			if v.Name() == sf && !v.IsDir() {
				return filepath.Join(confPath, v.Name()), nil
			}
		}

		return "", fmt.Errorf("file '%s' is not found", sf)
	}

	setCommonSettings := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			conf.Common.Logs = ls.Logging
		}

		z := ZabbixSet{}
		if ok := viper.IsSet("ZABBIX"); ok {
			if err := viper.GetViper().Unmarshal(&z); err != nil {
				return err
			}

			np := 10051
			if z.Zabbix.NetworkPort != 0 && z.Zabbix.NetworkPort < 65536 {
				np = z.Zabbix.NetworkPort
			}

			conf.Common.Zabbix = ZabbixOptions{
				NetworkPort: np,
				NetworkHost: z.Zabbix.NetworkHost,
				ZabbixHost:  z.Zabbix.ZabbixHost,
				EventTypes:  z.Zabbix.EventTypes,
			}
		}

		return nil
	}

	setSpecial := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		//Настройки для модуля подключения к NATS
		if viper.IsSet("NATS.host") {
			conf.NATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			conf.NATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cache_ttl") {
			conf.NATS.CacheTTL = viper.GetInt("NATS.cache_ttl")
		}
		if viper.IsSet("NATS.subscription") {
			conf.NATS.Subscription = viper.GetString("NATS.subscription")
		}

		// Настройки доступа к БД GeoIP
		if viper.IsSet("SensorInformationDataBase.host") {
			conf.SensorInformationDB.Host = viper.GetString("SensorInformationDataBase.host")
		}
		if viper.IsSet("SensorInformationDataBase.port") {
			conf.SensorInformationDB.Port = viper.GetInt("SensorInformationDataBase.port")
		}
		if viper.IsSet("SensorInformationDataBase.user") {
			conf.SensorInformationDB.User = viper.GetString("SensorInformationDataBase.user")
		}
		if viper.IsSet("SensorInformationDataBase.ncircc_url") {
			conf.SensorInformationDB.NCIRCCURL = viper.GetString("SensorInformationDataBase.ncircc_url")
		}
		if viper.IsSet("SensorInformationDataBase.request_timeout") {
			conf.SensorInformationDB.RequestTimeout = viper.GetInt("SensorInformationDataBase.request_timeout")
		}

		// Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("WriteLogDataBase.host") {
			conf.LogDB.Host = viper.GetString("WriteLogDataBase.host")
		}
		if viper.IsSet("WriteLogDataBase.port") {
			conf.LogDB.Port = viper.GetInt("WriteLogDataBase.port")
		}
		if viper.IsSet("WriteLogDataBase.user") {
			conf.LogDB.User = viper.GetString("WriteLogDataBase.user")
		}
		if viper.IsSet("WriteLogDataBase.namedb") {
			conf.LogDB.NameDB = viper.GetString("WriteLogDataBase.namedb")
		}
		if viper.IsSet("WriteLogDataBase.storage_name_db") {
			conf.LogDB.StorageNameDB = viper.GetString("WriteLogDataBase.storage_name_db")
		}

		return nil
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return conf, err
	}

	confPath := filepath.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return conf, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return conf, err
	}

	var fn string
	if envList["GO_ENRICHERSENSORINFO_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	} else {
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return conf, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_ENRICHERSENSORINFO_NHOST"] != "" {
		conf.NATS.Host = envList["GO_ENRICHERSENSORINFO_NHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_NPORT"]); err == nil {
			conf.NATS.Port = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_NCACHETTL"] != "" {
		if ttl, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_NCACHETTL"]); err == nil {
			conf.NATS.CacheTTL = ttl
		}
	}
	if envList["GO_ENRICHERSENSORINFO_NSUBSC"] != "" {
		conf.NATS.Subscription = envList["GO_ENRICHERSENSORINFO_NSUBSC"]
	}

	//Подключение к БД с информацией о сенсорах
	if envList["GO_ENRICHERSENSORINFO_SIHOST"] != "" {
		conf.SensorInformationDB.Host = envList["GO_ENRICHERSENSORINFO_SIHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_SIPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_SIPORT"]); err == nil {
			conf.SensorInformationDB.Port = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_SIUSER"] != "" {
		conf.SensorInformationDB.User = envList["GO_ENRICHERSENSORINFO_SIUSER"]
	}

	if envList["GO_ENRICHERSENSORINFO_SIPASSWD"] != "" {
		conf.SensorInformationDB.Passwd = envList["GO_ENRICHERSENSORINFO_SIPASSWD"]
	}
	if envList["GO_ENRICHERSENSORINFO_SIRTIMEOUT"] != "" {
		if timeout, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_SIRTIMEOUT"]); err == nil {
			conf.SensorInformationDB.RequestTimeout = timeout
		}
	}
	if envList["GO_ENRICHERSENSORINFO_SINCIRCCURL"] != "" {
		conf.SensorInformationDB.NCIRCCURL = envList["GO_ENRICHERSENSORINFO_SINCIRCCURL"]
	}
	if envList["GO_ENRICHERSENSORINFO_SINCIRCCTOKEN"] != "" {
		conf.SensorInformationDB.NCIRCCToken = envList["GO_ENRICHERSENSORINFO_SINCIRCCTOKEN"]
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_ENRICHERSENSORINFO_DBWLOGHOST"] != "" {
		conf.LogDB.Host = envList["GO_ENRICHERSENSORINFO_DBWLOGHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_DBWLOGPORT"]); err == nil {
			conf.LogDB.Port = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGNAME"] != "" {
		conf.LogDB.NameDB = envList["GO_ENRICHERSENSORINFO_DBWLOGNAME"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGUSER"] != "" {
		conf.LogDB.User = envList["GO_ENRICHERSENSORINFO_DBWLOGUSER"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGPASSWD"] != "" {
		conf.LogDB.Passwd = envList["GO_ENRICHERSENSORINFO_DBWLOGPASSWD"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME"] != "" {
		conf.LogDB.StorageNameDB = envList["GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return conf, err
	}

	return conf, nil
}
