package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goforj/godump"

	"github.com/av-belyakov/enricher_sensor_information/interfaces"
	"github.com/av-belyakov/enricher_sensor_information/internal/natsapi"
	"github.com/av-belyakov/enricher_sensor_information/internal/requests"
	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
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
	obt := natsapi.ObjectBeingTransferred{Id: msg.GetId()}

	if err := json.Unmarshal(msg.GetData(), &req); err != nil {
		obt.Data = fmt.Appendf(nil, strJson, "the request received an incorrect json format")
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())

		r.chToNatsApi <- &obt

		return
	}

	r.logger.Send("info", fmt.Sprintf("we are starting to process task Id '%s', which came from source '%s' and contains a request %v", req.TaskId, req.Source, req.ListSensor))

	if len(req.ListSensor) == 0 {
		errMsg := "it is impossible to perform a search, an empty list of sensors is received"
		obt.Data = fmt.Appendf(nil, strJson, errMsg)
		r.logger.Send("error", supportingfunctions.CustomError(errors.New(errMsg)).Error())

		r.chToNatsApi <- &obt

		return
	}

	response := responses.Response{
		TaskId:           req.TaskId,
		Source:           req.Source,
		FoundInformation: make([]responses.DetailedInformation, 0),
	}

	//поиск подробной информации о сенсорах
	foundInfo, err := r.commonInfo.Search(ctx, req.ListSensor)
	if err != nil {
		println("0000000 Error", err)

		response.Error = err.Error()
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
	} else {
		r.logger.Send("info", fmt.Sprintf("task Id '%s', the search was completed successfully, sensors list:'%v'", req.TaskId, req.ListSensor))
		godump.Dump(foundInfo)
	}

	println("--------------------")
	fmt.Printf("foundInfo:'%+v'\n", foundInfo)
	println("--------------------")

	response.FoundInformation = foundInfo

	resByte, err := json.Marshal(response)
	if err != nil {
		obt.Data = fmt.Appendf(nil, strJson, err.Error())
		r.logger.Send("error", supportingfunctions.CustomError(err).Error())
	} else {
		obt.Data = resByte
		r.counter.SendMessage("update processed events", 1)
		r.logger.Send("info", fmt.Sprintf("the request for taskId '%s' from source '%s' has been processed", req.TaskId, req.Source))
	}

	r.chToNatsApi <- &obt
}
