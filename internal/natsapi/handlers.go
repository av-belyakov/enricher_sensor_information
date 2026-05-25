package natsapi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// subscriptionHandler обработчик подписки приёма запросов
func (api *apiNatsModule) subscriptionRequestHandler() {
	_, err := api.natsConn.Subscribe(api.subscriptionRequest, func(m *nats.Msg) {
		id := uuid.NewString()

		api.storage.SetReq(id, m)
		api.chFromModule <- &ObjectBeingTransferred{
			Id:   id,
			Data: m.Data,
		}

		//счетчик принятых запросов
		api.counter.SendMessage("update accepted events", 1)
	})
	if err != nil {
		api.logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
}

// incomingInformationHandler обработчик информации полученной изнутри приложения
func (api *apiNatsModule) incomingInformationHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case incomingData := <-api.chToModule:
			m, ok := api.storage.GetReq(incomingData.GetId())
			if !ok {
				api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the responder for the request with id '%s' was not found", incomingData.GetId())).Error())

				continue
			}

			m.Respond(incomingData.GetData())
			api.storage.DelReq(incomingData.GetId())
		}
	}
}
