package requests

// Request запрос на поиск информации о сенсоре
type Request struct {
	ListSensor []string `json:"list_ip_addresses"`
	TaskId     string   `json:"task_id"`
	Source     string   `json:"source"`
}
