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
	cfg := &ConfigApp{}

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
			"GO_ENRICHERSENSORINFO_ZUSETLS":     "",
			"GO_ENRICHERSENSORINFO_ZHOST":       "",
			"GO_ENRICHERSENSORINFO_ZPORT":       "",
			"GO_ENRICHERSENSORINFO_ZUSER":       "",
			"GO_ENRICHERSENSORINFO_NBHOST":      "",
			"GO_ENRICHERSENSORINFO_NBPORT":      "",
			"GO_ENRICHERSENSORINFO_NCIRCCURL":   "",
			"GO_ENRICHERSENSORINFO_ZPASSWD":     "",
			"GO_ENRICHERSENSORINFO_NBTOKEN":     "",
			"GO_ENRICHERSENSORINFO_NCIRCCTOKEN": "",
			"GO_ENRICHERSENSORINFO_RTIMEOUT":    "",

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

			cfg.Common.Logs = ls.Logging
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

			cfg.Common.Zabbix = ZabbixOptions{
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
			cfg.NATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			cfg.NATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cache_ttl") {
			cfg.NATS.CacheTTL = viper.GetInt("NATS.cache_ttl")
		}
		if viper.IsSet("NATS.subscription") {
			cfg.NATS.Subscription = viper.GetString("NATS.subscription")
		}

		// Настройки доступа к базам обогащения дополнительной информацией
		if viper.IsSet("SensorInformationDataBase.zabbix_use_tls") {
			cfg.SensorInformationDB.ZabbixUseTLS = viper.GetBool("SensorInformationDataBase.zabbix_use_tls")
		}
		if viper.IsSet("SensorInformationDataBase.zabbix_host") {
			cfg.SensorInformationDB.ZabbixHost = viper.GetString("SensorInformationDataBase.zabbix_host")
		}
		if viper.IsSet("SensorInformationDataBase.zabbix_port") {
			cfg.SensorInformationDB.ZabbixPort = viper.GetInt("SensorInformationDataBase.zabbix_port")
		}
		if viper.IsSet("SensorInformationDataBase.zabbix_user") {
			cfg.SensorInformationDB.ZabbixUser = viper.GetString("SensorInformationDataBase.zabbix_user")
		}
		if viper.IsSet("SensorInformationDataBase.netbox_host") {
			cfg.SensorInformationDB.NetboxHost = viper.GetString("SensorInformationDataBase.netbox_host")
		}
		if viper.IsSet("SensorInformationDataBase.netbox_port") {
			cfg.SensorInformationDB.NetboxPort = viper.GetInt("SensorInformationDataBase.netbox_port")
		}
		if viper.IsSet("SensorInformationDataBase.ncircc_url") {
			cfg.SensorInformationDB.NCIRCCURL = viper.GetString("SensorInformationDataBase.ncircc_url")
		}
		if viper.IsSet("SensorInformationDataBase.request_timeout") {
			cfg.SensorInformationDB.RequestTimeout = viper.GetInt("SensorInformationDataBase.request_timeout")
		}

		// Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("WriteLogDataBase.host") {
			cfg.LogDB.Host = viper.GetString("WriteLogDataBase.host")
		}
		if viper.IsSet("WriteLogDataBase.port") {
			cfg.LogDB.Port = viper.GetInt("WriteLogDataBase.port")
		}
		if viper.IsSet("WriteLogDataBase.user") {
			cfg.LogDB.User = viper.GetString("WriteLogDataBase.user")
		}
		if viper.IsSet("WriteLogDataBase.namedb") {
			cfg.LogDB.NameDB = viper.GetString("WriteLogDataBase.namedb")
		}
		if viper.IsSet("WriteLogDataBase.storage_name_db") {
			cfg.LogDB.StorageNameDB = viper.GetString("WriteLogDataBase.storage_name_db")
		}

		// Настройки для отладочного сервера
		if viper.IsSet("DebugServer.enable") {
			cfg.DebugServer.Enable = viper.GetBool("DebugServer.enable")
		}
		if viper.IsSet("DebugServer.host") {
			cfg.DebugServer.Host = viper.GetString("DebugServer.host")
		}
		if viper.IsSet("DebugServer.port") {
			cfg.DebugServer.Port = viper.GetInt("DebugServer.port")
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
		return cfg, err
	}

	confPath := filepath.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return cfg, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return cfg, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return cfg, err
	}

	var fn string
	switch envList["GO_ENRICHERSENSORINFO_MAIN"] {
	case "development":
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return cfg, err
		}

	case "test":
		fn, err = getFileName("config_test.yml", confPath, list)
		if err != nil {
			return cfg, err
		}

	default:
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return cfg, err
		}

	}

	if err := setSpecial(fn); err != nil {
		return cfg, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_ENRICHERSENSORINFO_NHOST"] != "" {
		cfg.NATS.Host = envList["GO_ENRICHERSENSORINFO_NHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_NPORT"]); err == nil {
			cfg.NATS.Port = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_NCACHETTL"] != "" {
		if ttl, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_NCACHETTL"]); err == nil {
			cfg.NATS.CacheTTL = ttl
		}
	}
	if envList["GO_ENRICHERSENSORINFO_NSUBSC"] != "" {
		cfg.NATS.Subscription = envList["GO_ENRICHERSENSORINFO_NSUBSC"]
	}

	//Подключение к БД с информацией о сенсорах
	if envList["GO_ENRICHERSENSORINFO_ZHOST"] != "" {
		cfg.SensorInformationDB.ZabbixHost = envList["GO_ENRICHERSENSORINFO_ZHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_ZPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_ZPORT"]); err == nil {
			cfg.SensorInformationDB.ZabbixPort = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_ZUSER"] != "" {
		cfg.SensorInformationDB.ZabbixUser = envList["GO_ENRICHERSENSORINFO_ZUSER"]
	}
	if envList["GO_ENRICHERSENSORINFO_ZUSETLS"] != "" {
		if envList["GO_ENRICHERSENSORINFO_ZUSETLS"] == "true" {
			cfg.SensorInformationDB.ZabbixUseTLS = true
		} else {
			cfg.SensorInformationDB.ZabbixUseTLS = false
		}
	}

	if envList["GO_ENRICHERSENSORINFO_NBHOST"] != "" {
		cfg.SensorInformationDB.NetboxHost = envList["GO_ENRICHERSENSORINFO_NBHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_NBPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_NBPORT"]); err == nil {
			cfg.SensorInformationDB.NetboxPort = p
		}
	}

	if envList["GO_ENRICHERSENSORINFO_NCIRCCURL"] != "" {
		cfg.SensorInformationDB.NCIRCCURL = envList["GO_ENRICHERSENSORINFO_NCIRCCURL"]
	}
	if envList["GO_ENRICHERSENSORINFO_ZPASSWD"] != "" {
		cfg.SensorInformationDB.ZabbixPasswd = envList["GO_ENRICHERSENSORINFO_ZPASSWD"]
	}
	if envList["GO_ENRICHERSENSORINFO_NBTOKEN"] != "" {
		cfg.SensorInformationDB.NetboxToken = envList["GO_ENRICHERSENSORINFO_NBTOKEN"]
	}
	if envList["GO_ENRICHERSENSORINFO_NCIRCCTOKEN"] != "" {
		cfg.SensorInformationDB.NCIRCCToken = envList["GO_ENRICHERSENSORINFO_NCIRCCTOKEN"]
	}
	if envList["GO_ENRICHERSENSORINFO_RTIMEOUT"] != "" {
		if timeout, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_RTIMEOUT"]); err == nil {
			cfg.SensorInformationDB.RequestTimeout = timeout
		}
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_ENRICHERSENSORINFO_DBWLOGHOST"] != "" {
		cfg.LogDB.Host = envList["GO_ENRICHERSENSORINFO_DBWLOGHOST"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_ENRICHERSENSORINFO_DBWLOGPORT"]); err == nil {
			cfg.LogDB.Port = p
		}
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGNAME"] != "" {
		cfg.LogDB.NameDB = envList["GO_ENRICHERSENSORINFO_DBWLOGNAME"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGUSER"] != "" {
		cfg.LogDB.User = envList["GO_ENRICHERSENSORINFO_DBWLOGUSER"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGPASSWD"] != "" {
		cfg.LogDB.Passwd = envList["GO_ENRICHERSENSORINFO_DBWLOGPASSWD"]
	}
	if envList["GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME"] != "" {
		cfg.LogDB.StorageNameDB = envList["GO_ENRICHERSENSORINFO_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
