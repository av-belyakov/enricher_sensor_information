package netboxinteractions

import "net/http"

// Settings настройки для подключения к Netbox
type Settings struct {
	token   string
	host    string
	port    int
	timeout int
}

// Client клиент для работы с Netbox
type Client struct {
	client   *http.Client
	settings Settings
}

type Options func(*Client) error

// DcimDivicesCountLimitedInformation информация об общем количестве устройств
type DcimDivicesCountLimitedInformation struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Count    int    `json:"count"`
}

// DcimDivicesLimitedInformation информация о устройствах
type DcimDivicesLimitedInformation struct {
	Next     string                     `json:"next"`
	Previous string                     `json:"previous"`
	Count    int                        `json:"count"`
	Results  []DeviceLimitedInformation `json:"results"`
}

type DeviceLimitedInformation struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
