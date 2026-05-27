package netboxinteractions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/av-belyakov/enricher_sensor_information/internal/supportingfunctions"
)

// GetSizeDevices общее количество устройств
func (c *Client) GetCountDevices(ctx context.Context) (int, int, error) {
	res, statusCode, err := c.Get(ctx, "/api/dcim/devices/?fields=id,name&limit=1&offset=0")
	if err != nil {
		return 0, statusCode, err
	}

	if statusCode != http.StatusOK {
		return 0, statusCode, fmt.Errorf("status code: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	var result DcimDivicesCountLimitedInformation
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, statusCode, err
	}

	return result.Count, statusCode, nil
}

// GetDevicesLimitInformation краткая информация о устройствах
func (c *Client) GetDevicesLimitInformation(ctx context.Context, limit, offset int) (DcimDivicesLimitedInformation, int, error) {
	limitedInformation := DcimDivicesLimitedInformation{}
	res, statusCode, err := c.Get(ctx, fmt.Sprintf("/api/dcim/devices/?fields=id,name&limit=%d&offset=%d", limit, offset))
	if err != nil {
		return limitedInformation, statusCode, err
	}

	if statusCode != http.StatusOK {
		return limitedInformation, statusCode, fmt.Errorf("status code: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	if err := json.Unmarshal(res, &limitedInformation); err != nil {
		return limitedInformation, statusCode, err
	}

	return limitedInformation, statusCode, nil
}

// GetTenantGroups получить группу арендаторов устроства
func (c *Client) GetTenantGroups(ctx context.Context, deviceId int) (string, int, error) {
	// получаем информацию по устройству
	data, statusCode, err := c.Get(ctx, fmt.Sprintf("/api/dcim/devices/%d/", deviceId))
	if err != nil {
		return "", statusCode, err
	}

	if statusCode != http.StatusOK {
		return "", statusCode, fmt.Errorf("status code: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	tenantId := struct {
		Tenant struct {
			Id int `json:"id"`
		} `json:"tenant"`
	}{}

	if err := json.Unmarshal(data, &tenantId); err != nil {
		return "", statusCode, err
	}

	// получаем информацию по арендатору
	data, statusCode, err = c.Get(ctx, fmt.Sprintf("/api/tenancy/tenants/%d/", tenantId.Tenant.Id))
	if err != nil {
		return "", statusCode, err
	}

	if statusCode != http.StatusOK {
		return "", statusCode, fmt.Errorf("status code: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	tenantGroupId := struct {
		Group struct {
			Id int `json:"id"`
		} `json:"group"`
	}{}

	if err := json.Unmarshal(data, &tenantGroupId); err != nil {
		return "", statusCode, err
	}

	// получаем информацию по группе арендаторов
	data, statusCode, err = c.Get(ctx, fmt.Sprintf("/api/tenancy/tenant-groups/%d/", tenantGroupId.Group.Id))
	if err != nil {
		return "", statusCode, err
	}

	if statusCode != http.StatusOK {
		return "", statusCode, fmt.Errorf("status code: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	tenantGroupName := struct {
		Name string `json:"name"`
	}{}

	if err := json.Unmarshal(data, &tenantGroupName); err != nil {
		return "", statusCode, err
	}

	return tenantGroupName.Name, statusCode, nil
}

// Get реализация HTTP GET запроса
func (c *Client) Get(ctx context.Context, query string) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.settings.timeout)*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s:%d%s", c.settings.host, c.settings.port, query)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Token "+c.settings.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
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
