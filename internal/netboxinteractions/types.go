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
