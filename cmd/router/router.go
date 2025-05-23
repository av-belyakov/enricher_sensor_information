package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/enricher_sensor_information/cmd/natsapi"
	"github.com/av-belyakov/enricher_sensor_information/cmd/sensorinformationapi"
	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/requests"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

type Router struct {
	counter          interfaces.Counter
	logger           interfaces.Logger
	sensorInfoClient *sensorinformationapi.SensorInformationClient
	chFromNatsApi    <-chan interfaces.Requester
	chToNatsApi      chan<- interfaces.Responser
}

func NewRouter(
	counter interfaces.Counter,
	logger interfaces.Logger,
	sensorInfoClient *sensorinformationapi.SensorInformationClient,
	chFromNatsApi <-chan interfaces.Requester,
	chToNatsApi chan<- interfaces.Responser,
) *Router {
	return &Router{
		sensorInfoClient: sensorInfoClient,
		counter:          counter,
		logger:           logger,
		chFromNatsApi:    chFromNatsApi,
		chToNatsApi:      chToNatsApi,
	}
}

func (r *Router) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-r.chFromNatsApi:
				go r.handlerRequest(ctx, msg)
			}
		}
	}()
}

func (r *Router) handlerRequest(ctx context.Context, msg interfaces.Requester) {
	response := &natsapi.ObjectToNats{
		Id: msg.GetId(),
	}

	if ctx.Err() != nil {
		return
	}

	var req requests.Request
	if err := json.Unmarshal(msg.GetData(), &req); err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		response.Error = errors.New("the request received an incorrect json format")
		r.chToNatsApi <- response

		return
	}

	response.TaskId = req.TaskId
	response.Source = req.Source

	r.logger.Send("info", fmt.Sprintf("we are starting to process task Id '%s', which came from source '%s' and contains a request %v", req.TaskId, req.Source, req.ListSensor))

	results := make([]responses.DetailedInformation, 0, len(req.ListSensor))
	for _, sensor := range req.ListSensor {
		result := responses.DetailedInformation{SensorId: sensor}

		//поиск подробной информации о сенсоре
		res, err := r.sensorInfoClient.SearchSensorInfo(ctx, sensor)
		if err != nil {
			result.Error = "error interacting with a remote database"
			results = append(results, result)
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		results = append(results, res)
	}

	byteData, err := json.Marshal(results)
	if err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	response.Data = byteData

	r.counter.SendMessage("update processed events", 1)
	r.logger.Send("info", fmt.Sprintf("the request for taskId '%s' from source '%s' has been processed", req.TaskId, req.Source))

	r.chToNatsApi <- response
}
