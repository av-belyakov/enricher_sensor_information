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
	"github.com/goforj/godump"
)

// NewRouter инициализация маршрутизатора
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

// Start
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
	if ctx.Err() != nil {
		return
	}

	strJson := `{
			      "found_information": [],
			      "task_id": "",
			      "source": "",
			      "error": "%s"
			    }`

	var req requests.Request
	if err := json.Unmarshal(msg.GetData(), &req); err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		r.chToNatsApi <- &natsapi.ObjectBeingTransferred{
			Id:   msg.GetId(),
			Data: fmt.Appendf(nil, strJson, "the request received an incorrect json format"),
		}

		return
	}

	r.logger.Send("info", fmt.Sprintf("we are starting to process task Id '%s', which came from source '%s' and contains a request %v", req.TaskId, req.Source, req.ListSensor))

	if len(req.ListSensor) == 0 {
		errMsg := "it is impossible to perform a search, an empty list of sensors is received"
		r.logger.Send("error", supportingfunctions.CustomError(errors.New(errMsg)).Error())
		r.chToNatsApi <- &natsapi.ObjectBeingTransferred{
			Id:   msg.GetId(),
			Data: fmt.Appendf(nil, strJson, errMsg),
		}

		return
	}

	results := make([]responses.DetailedInformation, 0, len(req.ListSensor))
	for _, sensor := range req.ListSensor {
		//поиск подробной информации о сенсоре
		res, err := r.commonInfo.Search(ctx, sensor)
		if err != nil {
			//fmt.Println("func 'Router.handlerRequest', ERROR:", err)

			res.Error = "error interacting with a remote database"
			r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		res.SensorId = sensor
		results = append(results, res)
	}

	//fmt.Println("func 'Router.handlerRequest', result:")
	//godump.Dump(results)
	r.logger.Send("info", fmt.Sprintf("task Id '%s', result:'%s'", req.TaskId, godump.DumpStr(results)))

	resByte, err := json.Marshal(responses.Response{
		TaskId:           req.TaskId,
		Source:           req.Source,
		FoundInformation: results,
	})
	if err != nil {
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
		r.chToNatsApi <- &natsapi.ObjectBeingTransferred{
			Id:   msg.GetId(),
			Data: fmt.Appendf(nil, strJson, err.Error()),
		}

		return
	}

	r.counter.SendMessage("update processed events", 1)
	r.logger.Send("info", fmt.Sprintf("the request for taskId '%s' from source '%s' has been processed", req.TaskId, req.Source))

	r.chToNatsApi <- &natsapi.ObjectBeingTransferred{
		Id:   msg.GetId(),
		Data: resByte,
	}
}
