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
	"github.com/av-belyakov/enricher_sensor_information/internal/information"
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
			SearchCommonInfo: app.diContainer.SensorInformationDB(),
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
			debugServerPort := 6263

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
	msg := information.GetInformationMessage(a.diContainer.Configer())
	a.diContainer.SimpleLogger(a.ctx).Write("info", strings.ToLower(msg))

	<-a.ctx.Done()
}
