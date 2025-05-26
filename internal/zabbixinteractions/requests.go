package zabbixinteractions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/av-belyakov/enricher_sensor_information/internal/responses"
	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// GetFullSensorInformation полная информация о сенсоре
func GetFullSensorInformation(ctx context.Context, sensorId string, zc *ZabbixConnectionJsonRPC) (responses.DetailedInformation, error) {
	fullInfo := responses.DetailedInformation{SensorId: sensorId}
	sensorInfo, err := GetSpecialId(ctx, sensorId, zc)
	if err != nil {
		//при возникновении ошибки пытаемся авторизоватся повторно, так как при устаревании
		//авторизационного хеша возможно появления ошибки с сообщением: 'Invalid params. Session terminated, re-login, please.'.
		err = zc.Authorization(ctx)
		if err != nil {
			return fullInfo, err
		}

		sensorInfo, err = GetSpecialId(ctx, sensorId, zc)
		if err != nil {
			return fullInfo, err
		}
	}

	//получаем специальный код сенсора, вся информация в Zabbix хранится именно с ним
	fullInfo.SpecialSensorId = sensorInfo.GetSpecialId()

	//получаем географический код
	geoCode, err := sensorInfo.GetGeoCode(ctx)
	if err != nil {
		return fullInfo, err
	}
	fullInfo.GeoCode = geoCode

	//получаем сферу деятельности объекта
	objectArea, err := sensorInfo.GetObjectArea(ctx)
	if err != nil {
		return fullInfo, err
	}
	fullInfo.ObjectArea = objectArea

	//получаем наименование субъекта Российской Федерации
	subjectRF, err := sensorInfo.GetSubjectRF(ctx)
	if err != nil {
		return fullInfo, err
	}
	fullInfo.SubjectRussianFederation = subjectRF

	//получаем индивидуальный налоговый идентификатор
	inn, err := sensorInfo.GetINN(ctx)
	if err != nil {
		return fullInfo, err
	}
	fullInfo.INN = inn

	//получаем перечень домашних сетей контролируемого объекта
	homeNet, err := sensorInfo.GetHomeNet(ctx)
	if err != nil {
		return fullInfo, err
	}
	fullInfo.HomeNet = homeNet

	return fullInfo, nil
}

// GetSpecialId объект со специальным id сенсора который нужен для подробных запросов
func GetSpecialId(ctx context.Context, sensorId string, zc *ZabbixConnectionJsonRPC) (*RequiestSensorInfo, error) {
	req := RequiestSensorInfo{zabbixConnection: zc}

	strReq := fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"host.get",
	  "params":{
	    "search":{"host":%s}
	  },
	  "auth":"%s",
	  "id":1
	}`, sensorId, zc.GetAuthorizationData())

	if sensorId == "" {
		return &req, supportingfunctions.CustomError(errors.New("the sensor ID cannot be equal to 0"))
	}

	res, err := zc.PostRequest(ctx, strings.NewReader(strReq))
	if err != nil {
		return &req, supportingfunctions.CustomError(fmt.Errorf("error send post request, %w", err))
	}

	resData := ResponseData{}
	err = json.Unmarshal(res, &resData)
	if err != nil {
		return &req, supportingfunctions.CustomError(fmt.Errorf("error decode request, %w", err))
	}

	if len(resData.Error) > 0 {
		var msg, data string

		for k, v := range resData.Error {
			if k == "message" {
				msg = fmt.Sprint(v)
			}

			if k == "data" {
				data = fmt.Sprint(v)
			}
		}

		return &req, supportingfunctions.CustomError(fmt.Errorf("error send post request, message:'%s', data:'%s'", msg, data))
	}

DONE:
	for _, v := range resData.Result {
		for key, value := range v {
			if key == "hostid" {
				req.specialId = fmt.Sprint(value)

				break DONE
			}
		}
	}

	return &req, nil
}
