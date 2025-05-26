package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/enricher_sensor_information/cmd/natsapi"
	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/requests"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

func NewRouter(
	counter interfaces.Counter,
	logger interfaces.Logger,
	settings RouterSettings) *Router {
	return &Router{
		counter:       counter,
		logger:        logger,
		commonInfo:    settings.SearchCommonInfo,
		chFromNatsApi: settings.ChanFromNatsApi,
		chToNatsApi:   settings.ChanToNatsApi,
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
		//поиск подробной информации о сенсоре
		res, err := r.commonInfo.Search(ctx, sensor)
		if err != nil {
			res.Error = "error interacting with a remote database"
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		res.SensorId = sensor
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
