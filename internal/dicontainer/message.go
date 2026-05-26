package dicontainer

import (
	"context"
	"log"
	"os"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/enricher_sensor_information/internal/confighandler"
	"github.com/av-belyakov/enricher_sensor_information/internal/countermessage"
	"github.com/av-belyakov/enricher_sensor_information/internal/elasticsearchapi"
	"github.com/av-belyakov/enricher_sensor_information/internal/logginghandler"
	"github.com/av-belyakov/enricher_sensor_information/internal/natsapi"
	"github.com/av-belyakov/enricher_sensor_information/internal/sensorinformationapi"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// Configer чтения конфигурационного файла
func (d *DiContainer) Configer() Configer {
	if d.configer == nil {
		rootPath, err := supportingfunctions.GetRootPath(d.rootDir)
		if err != nil {
			log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
		}

		cfg, err := confighandler.New(rootPath)
		if err != nil {
			log.Fatal("error module 'confighandler':", err)
		}

		d.configer = cfg
	}

	return d.configer
}

// SimpleLogger простое логирование с помощью стороннего пакета
func (d *DiContainer) SimpleLogger(ctx context.Context) SimpleLogger {
	if d.simpleLogger == nil {
		listLog := make([]simplelogger.OptionsManager, 0, len(d.Configer().GetListLogs()))
		for _, v := range d.Configer().GetListLogs() {
			listLog = append(listLog, v)
		}

		opts := simplelogger.CreateOptions(listLog...)
		simpleLogger, err := simplelogger.NewSimpleLogger(ctx, d.rootDir, opts)
		if err != nil {
			log.Fatal("error module 'simplelogger':", err)
		}

		d.simpleLogger = simpleLogger

		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(d.DbLogger())
	}

	return d.simpleLogger
}

// Logger основное логирование
func (d *DiContainer) Logger(ctx context.Context) Logger {
	if d.logger == nil {
		logger := logginghandler.New(d.SimpleLogger(ctx), d.ch)
		logger.Start(ctx)

		d.logger = logger
	}

	return d.logger
}

// Counter счетчик сообщений
func (d *DiContainer) Counter(ctx context.Context) Counter {
	if d.counter == nil {
		counter := countermessage.New(d.ch)
		counter.Start(ctx)

		d.counter = counter
	}

	return d.counter
}

// DbLogger запись логов в БД
func (d *DiContainer) DbLogger() DbLogger {
	if d.dbLogger == nil {
		var nameRegionalObject = "gcm"
		if os.Getenv("GO_PHMISP_MAIN") == "development" {
			nameRegionalObject = "gcm-test"
		}

		conn, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
			Port:               d.Configer().GetLogDB().Port,
			Host:               d.Configer().GetLogDB().Host,
			User:               d.Configer().GetLogDB().User,
			Passwd:             d.Configer().GetLogDB().Passwd,
			IndexDB:            d.Configer().GetLogDB().StorageNameDB,
			NameRegionalObject: nameRegionalObject,
		})
		if err != nil {
			log.Fatal("error module 'elasticsearchapi':", err)
		}

		d.dbLogger = conn
	}

	return d.dbLogger
}

// NatsConnecter подключение к NATS
func (d *DiContainer) NatsConnecter(ctx context.Context) NatsConnecter {
	if d.nats == nil {
		apiNats, err := natsapi.New(
			d.Logger(ctx),
			d.Counter(ctx),
			natsapi.WithHost(d.Configer().GetNATS().Host),
			natsapi.WithPort(d.Configer().GetNATS().Port),
			natsapi.WithCacheTTL(d.Configer().GetNATS().CacheTTL),
			natsapi.WithSubscription(d.Configer().GetNATS().Subscription))
		if err != nil {
			log.Fatal("error initialization module 'natsapi':", err)
		}

		if err = apiNats.Start(ctx); err != nil {
			log.Fatal("error start module 'natsapi':", err)
		}

		d.nats = apiNats
	}

	return d.nats
}

// SensorInformationDB подключение к БД с данными о сенсорах
func (d *DiContainer) SensorInformationDB() SensorInformationConnecter {
	if d.sensorInformationDB == nil {
		client, err := sensorinformationapi.New(
			sensorinformationapi.WithHost(d.Configer().GetSensorInformationDB().ZabbixHost),
			sensorinformationapi.WithUser(d.Configer().GetSensorInformationDB().ZabbixUser),
			sensorinformationapi.WithPasswd(d.Configer().GetSensorInformationDB().ZabbixPasswd),
			sensorinformationapi.WithNCIRCCURL(d.Configer().GetSensorInformationDB().NCIRCCURL),
			sensorinformationapi.WithNCIRCCToken(d.Configer().GetSensorInformationDB().NCIRCCToken),
			sensorinformationapi.WithRequestTimeout(d.Configer().GetSensorInformationDB().RequestTimeout))
		if err != nil {
			log.Fatal("error initialization module 'sensorinformationapi':", err)
		}

		d.sensorInformationDB = client
	}

	return d.sensorInformationDB
}
