package requests

// Request запрос на поиск информации о сенсоре
type Request struct {
	ListSensor []string `json:"list_sensors"`
	TaskId     string   `json:"task_id"`
	Source     string   `json:"source"`
}
