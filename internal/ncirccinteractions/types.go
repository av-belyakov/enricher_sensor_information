package ncirccinteractions

import (
	"net/http"
	"time"
)

// ClientNICRCC клиент НКЦКИ
type ClientNICRCC struct {
	client            *http.Client
	connectionTimeout time.Duration
	url               string
	token             string
}

// Response результат поиска информации об организации по ИНН
type Response struct {
	Data    []DetailedInformation `json:"data"`
	Count   int                   `json:"total"`
	Success bool                  `json:"success"`
}

// DetailedInformation
type DetailedInformation struct {
	Name       string `json:"settings_name"`
	Type       string `json:"settings_subject_type"`
	Sname      string `json:"settings_sname"`
	SubjectINN string `json:"settings_inn_of_subject"`
}
