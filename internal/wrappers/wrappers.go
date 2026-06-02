package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionzabbixagent"
	zainterfaces "github.com/av-belyakov/zabbixapicommunicator/v2/interfaces"

	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// WrappersZabbixInteraction обёртка для взаимодействия с модулем zabbixapi
func WrappersZabbixInteraction(
	ctx context.Context,
	settings WrappersZabbixInteractionSettings,
	logging interfaces.WriterLoggingData,
	channelZabbix <-chan interfaces.Messager) {

	connTimeout := time.Duration(3 * time.Second)
	zc, err := connectionzabbixagent.New(connectionzabbixagent.SettingsZabbixConnection{
		Port:              settings.NetworkPort,
		Host:              settings.NetworkHost,
		NetProto:          "tcp",
		ZabbixHost:        settings.ZabbixHost,
		ConnectionTimeout: &connTimeout,
	})
	if err != nil {
		logging.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %w", err)).Error())

		return
	}

	et := make([]connectionzabbixagent.EventType, len(settings.EventTypes))
	for _, v := range settings.EventTypes {
		et = append(et, connectionzabbixagent.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: connectionzabbixagent.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}

	recipient := make(chan zainterfaces.Messager)
	if err = zc.Start(ctx, et, recipient); err != nil {
		logging.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %w", err)).Error())

		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-channelZabbix:
				newMessageSettings := &connectionzabbixagent.MessageSettings{}
				newMessageSettings.SetType(msg.GetType())
				newMessageSettings.SetMessage(msg.GetMessage())

				recipient <- newMessageSettings

			case errMsg := <-zc.GetChanErr():
				logging.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %W", errMsg)).Error())

			}
		}
	}()
}
