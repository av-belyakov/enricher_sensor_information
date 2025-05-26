package zabbixinteractions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// GetSpecialId специальный идентификатор под которым информация о сенсор хранится в БД
func (sid *RequiestSensorInfo) GetSpecialId() string {
	return sid.specialId
}

// SetSpecialId специальный идентификатор под которым информация о сенсор хранится в БД
func (sid *RequiestSensorInfo) SetSpecialId() string {
	return sid.specialId
}

// GetGeoCode географический код
func (sid *RequiestSensorInfo) GetGeoCode(ctx context.Context) (string, error) {
	return sid.sendRequest(ctx, fmt.Sprintf(`{
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
func (sid *RequiestSensorInfo) GetObjectArea(ctx context.Context) (string, error) {
	return sid.sendRequest(ctx, fmt.Sprintf(`{
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
func (sid *RequiestSensorInfo) GetSubjectRF(ctx context.Context) (string, error) {
	return sid.sendRequest(ctx, fmt.Sprintf(`{
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
func (sid *RequiestSensorInfo) GetINN(ctx context.Context) (string, error) {
	return sid.sendRequest(ctx, fmt.Sprintf(`{
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
func (sid *RequiestSensorInfo) GetHomeNet(ctx context.Context) (string, error) {
	return sid.sendRequest(ctx, fmt.Sprintf(`{
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
func (sid *RequiestSensorInfo) sendRequest(ctx context.Context, str string) (string, error) {
	res, err := sid.zabbixConnection.PostRequest(ctx, strings.NewReader(str))
	if err != nil {
		return "", supportingfunctions.CustomError(err)
	}

	var resData ResponseData
	err = json.Unmarshal(res, &resData)
	if err != nil {
		return "", supportingfunctions.CustomError(err)
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

		return "", supportingfunctions.CustomError(fmt.Errorf("%s. %s", msg, data))
	}

	for _, v := range resData.Result {
		for key, value := range v {
			if key == "lastvalue" {
				return fmt.Sprint(value), nil
			}
		}
	}

	return "", nil
}
