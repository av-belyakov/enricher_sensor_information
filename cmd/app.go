package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/enricher_sensor_information/cmd/router"
	"github.com/av-belyakov/enricher_sensor_information/constants"
	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/dicontainer"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
	"github.com/av-belyakov/enricher_sensor_information/internal/wrappers"
)

type App struct {
	diContainer *dicontainer.DiContainer
	router      *router.Router
	ctx         context.Context
}

func NewApp(ctx context.Context) *App {
	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	ch := make(chan interfaces.Messager)
	app := &App{
		ctx:         ctx,
		diContainer: dicontainer.NewDIContainer(rootPath, ch),
	}

	// настройка обёртки для взаимодействия с Zabbix
	zabbixSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: app.diContainer.Configer().GetCommon().Zabbix.NetworkPort,
		NetworkHost: app.diContainer.Configer().GetCommon().Zabbix.NetworkHost,
		ZabbixHost:  app.diContainer.Configer().GetCommon().Zabbix.ZabbixHost,
		EventTypes:  make([]wrappers.EventType, len(app.diContainer.Configer().GetCommon().Zabbix.EventTypes)),
	}
	for _, v := range app.diContainer.Configer().GetCommon().Zabbix.EventTypes {
		zabbixSettings.EventTypes = append(zabbixSettings.EventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	// обертка для взаимодействия с Zabbix
	wrappers.WrappersZabbixInteraction(ctx, zabbixSettings, app.diContainer.SimpleLogger(ctx), ch)

	app.router = router.NewRouter(
		app.diContainer.Counter(ctx),
		app.diContainer.Logger(ctx),
		router.RouterSettings{
			SearchCommonInfo: app.diContainer.SensoeInformationDB(),
			ChanFromNatsApi:  app.diContainer.NatsConnecter(ctx).GetChFromModule(),
			ChanToNatsApi:    app.diContainer.NatsConnecter(ctx).GetChToModule(),
		})

	return app
}

func (a *App) Start() {
	// сервер для отладки
	if os.Getenv("GO_ENRICHERSENSORINFO_MAIN") == "test" || os.Getenv("GO_ENRICHERSENSORINFO_MAIN") == "development" {
		go func() {
			debugServerHost := "localhost"
			debugServerPort := 6161

			httpServer := &http.Server{
				Addr: fmt.Sprintf("%s:%d", debugServerHost, debugServerPort),
				BaseContext: func(_ net.Listener) context.Context {
					return a.ctx
				},
			}

			g, gCtx := errgroup.WithContext(a.ctx)
			g.Go(func() error {
				return httpServer.ListenAndServe()
			})
			g.Go(func() error {
				<-gCtx.Done()

				return httpServer.Shutdown(context.Background())
			})

			log.Printf("%vdebug server %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, debugServerHost, debugServerPort, constants.Ansi_Reset)

			if err := g.Wait(); err != nil {
				log.Fatal("error debugging server:", err)
			}
		}()
	}

	// старт приложения
	a.router.Start(a.ctx)

	// вывод информационного сообщения при старте приложения
	msg := getInformationMessage(a.diContainer.Configer().GetNATS())
	a.diContainer.SimpleLogger(a.ctx).Write("info", strings.ToLower(msg))

	<-a.ctx.Done()
}

/*
func app(ctx context.Context) {
	var nameRegionalObject string
	if os.Getenv("GO_ENRICHERSENSORINFO_MAIN") == "development" {
		nameRegionalObject = "enricher_sensor_information-dev"
	} else {
		nameRegionalObject = "enricher_sensor_information"
	}

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// ****************************************************************************
	// *********** инициализируем модуль чтения конфигурационного файла ***********
	conf, err := confighandler.New(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	var listLog []simplelogger.OptionsManager
	for _, v := range conf.GetListLogs() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	confDB := conf.GetLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: nameRegionalObject,
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	// ************************************************************************
	// ************* инициализация модуля взаимодействия с Zabbix *************
	chZabbix := make(chan interfaces.Messager)
	confZabbix := conf.GetZabbix()
	wziSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confZabbix.NetworkPort,
		NetworkHost: confZabbix.NetworkHost,
		ZabbixHost:  confZabbix.ZabbixHost,
	}
	eventTypes := []wrappers.EventType(nil)
	for _, v := range confZabbix.EventTypes {
		eventTypes = append(eventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	wziSettings.EventTypes = eventTypes
	wrappers.WrappersZabbixInteraction(ctx, wziSettings, simpleLogger, chZabbix)

	//***************************************************************************
	//************** инициализация обработчика логирования данных ***************
	//фактически это мост между simpleLogger и пакетом соединения с Zabbix
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	// ***************************************************************************
	// *********** инициализируем модуль счётчика для подсчёта сообщений *********
	counting := countermessage.New(chZabbix)
	counting.Start(ctx)

	// ***********************************************************************
	// ************** инициализация модуля взаимодействия с NATS *************
	confNats := conf.NATS
	apiNats, err := natsapi.New(
		logging,
		counting,
		natsapi.WithHost(confNats.Host),
		natsapi.WithPort(confNats.Port),
		natsapi.WithCacheTTL(confNats.CacheTTL),
		natsapi.WithSubscription(confNats.Subscription))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}
	//--- старт модуля
	if err = apiNats.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

	// ***************************************************************
	// ************ инициализация модуля поиска информации ***********
	confSIDB := conf.GetSensorInformationDB()
	sensorInformationClient, err := sensorinformationapi.New(
		sensorinformationapi.WithHost(confSIDB.Host),
		sensorinformationapi.WithUser(confSIDB.User),
		sensorinformationapi.WithPasswd(confSIDB.Passwd),
		sensorinformationapi.WithNCIRCCURL(confSIDB.NCIRCCURL),
		sensorinformationapi.WithNCIRCCToken(confSIDB.NCIRCCToken),
		sensorinformationapi.WithRequestTimeout(confSIDB.RequestTimeout))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

	router := router.NewRouter(
		counting,
		logging,
		router.RouterSettings{
			SearchCommonInfo: sensorInformationClient,
			ChanFromNatsApi:  apiNats.GetChFromModule(),
			ChanToNatsApi:    apiNats.GetChToModule()})
	router.Start(ctx)

	//информационное сообщение
	msg := getInformationMessage(conf.GetNATS())
	_ = simpleLogger.Write("info", msg)

	<-ctx.Done()
}
*/
