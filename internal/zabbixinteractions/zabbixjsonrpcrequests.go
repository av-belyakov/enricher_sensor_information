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

// GetFullSensorInfo полная информация о сенсоре
func GetFullSensorInfo(ctx context.Context, sensorId string, zconn *ZabbixConnectionJsonRPC) (responses.DetailedInformation, error) {
	fullInfo := responses.DetailedInformation{SensorId: sensorId}
	reqSensorInfo, err := NewRequiestSensorInfo(sensorId, zconn)
	if err != nil {
		//при возникновении ошибки пытаемся авторизоватся повторно, так как при устаревании
		//авторизационного хеша возможно появления ошибки с сообщением: 'Invalid params. Session terminated, re-login, please.'.
		err = authorizationZabbixAPI(ctx, zconn.login, zconn.passwd, zconn)
		if err != nil {
			return fullInfo, err
		}

		//вторая попытка выполнения запроса на основе нового хеша
		reqSensorInfo, err = NewRequiestSensorInfo(sensorId, zconn)
		if err != nil {
			return fullInfo, err
		}
	}

	//получаем специальный код сенсора, вся информация в Zabbix хранится именно с ним
	fullInfo.SpecialSensorId = reqSensorInfo.GetSpecialId()

	//получаем географический код
	geoCode, err := reqSensorInfo.GetGeoCode()
	if err != nil {
		return fullInfo, err
	}
	fullInfo.GeoCode = geoCode

	//получаем сферу деятельности объекта
	objectArea, err := reqSensorInfo.GetObjectArea()
	if err != nil {
		return fullInfo, err
	}
	fullInfo.ObjectArea = objectArea

	//получаем наименование субъекта Российской Федерации
	subjectRF, err := reqSensorInfo.GetSubjectRF()
	if err != nil {
		return fullInfo, err
	}
	fullInfo.SubjectRussianFederation = subjectRF

	//получаем индивидуальный налоговый идентификатор
	inn, err := reqSensorInfo.GetINN()
	if err != nil {
		return fullInfo, err
	}
	fullInfo.INN = inn

	//получаем перечень домашних сетей контролируемого объекта
	homeNet, err := reqSensorInfo.GetHomeNet()
	if err != nil {
		return fullInfo, err
	}
	fullInfo.HomeNet = homeNet

	return fullInfo, nil
}

// NewRequiestSensorInfo объект запросов информации о сенсоре
func NewRequiestSensorInfo(sensorId string, zconn *ZabbixConnectionJsonRPC) (*RequiestSensorInfo, error) {
	req := RequiestSensorInfo{zabbixConnection: zconn}

	strReq := fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"host.get",
	  "params":{
	    "search":{"host":%s}
	  },
	  "auth":"%s",
	  "id":1
	}`, sensorId, zconn.GetAuthorizationData())

	if sensorId == "" {
		return &req, supportingfunctions.CustomError(errors.New("the sensor ID cannot be equal to 0"))
	}

	res, err := zconn.SendPostRequest(strings.NewReader(strReq))
	if err != nil {
		return &req, supportingfunctions.CustomError(fmt.Errorf("error send post request, %w", err))
	}

	rd := responseData{}
	err = json.NewDecoder(res).Decode(&rd)
	if err != nil {
		return &req, supportingfunctions.CustomError(fmt.Errorf("error decode request, %w", err))
	}

	if len(rd.Error) > 0 {
		var msg, data string

		for k, v := range rd.Error {
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
	for _, v := range rd.Result {
		for key, value := range v {
			if key == "hostid" {
				req.specialId = fmt.Sprint(value)

				break DONE
			}
		}
	}

	return &req, nil
}

// GetSpecialId специальный идентификатор под которым информация о сенсор хранится в БД
func (sid *RequiestSensorInfo) GetSpecialId() string {
	return sid.specialId
}

// SetSpecialId специальный идентификатор под которым информация о сенсор хранится в БД
func (sid *RequiestSensorInfo) SetSpecialId() string {
	return sid.specialId
}

// GetGeoCode географический код
func (sid *RequiestSensorInfo) GetGeoCode() (string, error) {
	return sid.sendRequest(fmt.Sprintf(`{
	  "jsonrpc": "2.0",
	  "method": "item.get",
	  "params": {
        "output": "extend",
		"hostids": "%s",
		"search": {"key_": "geo_code"},
		"sortfield": "name"
      },
	  "auth": "%s",
	  "id": 1
	}`, sid.specialId, sid.zabbixConnection.GetAuthorizationData()))
}

// GetObjectArea сфера деятельности объекта
func (sid *RequiestSensorInfo) GetObjectArea() (string, error) {
	return sid.sendRequest(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"item.get",
	  "params":{
        "output":"extend",
        "hostids":"%s",
		"search":{"key_": "object_area"},
		"sortfield":"name"
	  },
	  "auth":"%s",
	  "id": 1
	}`, sid.specialId, sid.zabbixConnection.GetAuthorizationData()))
}

// GetSubjectRF субъект Российской Федерации
func (sid *RequiestSensorInfo) GetSubjectRF() (string, error) {
	return sid.sendRequest(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"item.get",
	  "params":{
	    "output": "extend",
	    "hostids": "%s",
	    "search":{"key_":"subject_RF"},
	    "sortfield":"name"
	  },
	  "auth":"%s",
	  "id":1
	}`, sid.specialId, sid.zabbixConnection.GetAuthorizationData()))
}

// GetINN индивидуальный налоговый идентификатор
func (sid *RequiestSensorInfo) GetINN() (string, error) {
	return sid.sendRequest(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"item.get",
	  "params":{
	    "output":"extend",
	    "hostids":"%s",
	    "search":{"key_": "inn"},
	    "sortfield":"name"
	  },
	  "auth":"%s",
	  "id": 1
	}`, sid.specialId, sid.zabbixConnection.GetAuthorizationData()))
}

// GetHomeNet перечень домашних сетей
func (sid *RequiestSensorInfo) GetHomeNet() (string, error) {
	return sid.sendRequest(fmt.Sprintf(`{
      "jsonrpc":"2.0",
	  "method":"item.get",
	  "params":{
	    "output":"extend",
	    "hostids":"%s",
	    "search":{"key_":"home_net"},
	    "sortfield":"name"
	  },
	  "auth":"%s",
	  "id":1
	}`, sid.specialId, sid.zabbixConnection.GetAuthorizationData()))
}

// sendRequest передача запроса API
func (sid *RequiestSensorInfo) sendRequest(str string) (string, error) {
	res, err := sid.zabbixConnection.SendPostRequest(strings.NewReader(str))
	if err != nil {
		return "", supportingfunctions.CustomError(err)
	}

	rd := responseData{}
	err = json.NewDecoder(res).Decode(&rd)
	if err != nil {
		return "", supportingfunctions.CustomError(err)
	}

	if len(rd.Error) > 0 {
		var msg, data string

		for k, v := range rd.Error {
			if k == "message" {
				msg = fmt.Sprint(v)
			}

			if k == "data" {
				data = fmt.Sprint(v)
			}
		}

		return "", supportingfunctions.CustomError(fmt.Errorf("%s. %s", msg, data))
	}

	for _, v := range rd.Result {
		for key, value := range v {
			if key == "lastvalue" {
				return fmt.Sprint(value), nil
			}
		}
	}

	return "", nil
}
