package responses

type Response struct {
	FoundInformation []DetailedInformation `json:"found_information"`
	TaskId           string                `json:"task_id"`
	Source           string                `json:"source"`
	Error            string                `json:"error"`
}

type DetailedInformation struct {
	INN                      string `json:"inn"`                        //налоговый идентификатор
	HomeNet                  string `json:"home_net"`                   //перечень домашних сетей
	OrgName                  string `json:"organization_name"`          //наименование организации
	SensorId                 string `json:"sensor_id"`                  //идентификатор сенсора
	ObjectArea               string `json:"object_area"`                //сфера деятельности объекта
	FullOrgName              string `json:"full_organization_name"`     //полное наименование организации
	SpecialSensorId          string `json:"special_sensor_id"`          //идентификатор сенсора, специальный, для поиска информации в НКЦКИ
	SubjectRussianFederation string `json:"subject_russian_federation"` //субъект Российской Федерации
	Error                    string `json:"error"`
}
