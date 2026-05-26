package netboxinteractions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// Get реализация HTTP GET запроса
func (api *Client) Get(ctx context.Context, query string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(api.settings.timeout)*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s:%d%s", api.settings.host, api.settings.port, query)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Token "+api.settings.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := api.client.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("status code: %d (%s)", res.StatusCode, res.Status)
	}
	defer supportingfunctions.CloseHTTPResponse(res)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 500, err
	}

	return resBody, res.StatusCode, nil
}
