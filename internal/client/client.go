package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/langtind/homeyctl/internal/config"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.BaseURL(),
		token:   cfg.EffectiveToken(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Devices

func (c *Client) GetDevices() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/devices/device/", nil)
}

func (c *Client) GetDevice(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/devices/device/"+id, nil)
}

func (c *Client) SetCapability(deviceID, capability string, value interface{}) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/devices/device/%s/capability/%s", deviceID, capability), body)
	return err
}

func (c *Client) DeleteDevice(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/devices/device/"+id, nil)
	return err
}

func (c *Client) UpdateDevice(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/devices/device/"+id, updates)
	return err
}

func (c *Client) GetDeviceSettings(id string) (json.RawMessage, error) {
	return c.doRequest("GET", fmt.Sprintf("/api/manager/devices/device/%s/settings_obj", id), nil)
}

func (c *Client) SetDeviceSetting(deviceID string, settings map[string]interface{}) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/devices/device/%s/settings", deviceID), settings)
	return err
}

// Device Groups

func (c *Client) CreateDeviceGroup(group map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/devices/group", group)
}

func (c *Client) UpdateDeviceGroup(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/devices/group/"+id, updates)
	return err
}

func (c *Client) RemoveDeviceFromGroup(groupID, deviceID string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/manager/devices/group/%s/device/%s", groupID, deviceID), nil)
	return err
}

// Flows

func (c *Client) GetFlows() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flow/", nil)
}

func (c *Client) GetAdvancedFlows() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/advancedflow/", nil)
}

func (c *Client) TriggerFlow(id string) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/manager/flow/flow/%s/trigger", id), nil)
	return err
}

func (c *Client) TriggerAdvancedFlow(id string) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/manager/flow/advancedflow/%s/trigger", id), nil)
	return err
}

func (c *Client) CreateFlow(flow map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/flow/flow/", flow)
}

func (c *Client) CreateAdvancedFlow(flow map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/flow/advancedflow/", flow)
}

func (c *Client) UpdateFlow(id string, flow map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("PUT", "/api/manager/flow/flow/"+id, flow)
}

func (c *Client) DeleteFlow(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/flow/flow/"+id, nil)
	return err
}

func (c *Client) UpdateAdvancedFlow(id string, flow map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("PUT", "/api/manager/flow/advancedflow/"+id, flow)
}

func (c *Client) DeleteAdvancedFlow(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/flow/advancedflow/"+id, nil)
	return err
}

// Flow cards

func (c *Client) GetFlowTriggers() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flowcardtrigger/", nil)
}

func (c *Client) GetFlowConditions() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flowcardcondition/", nil)
}

func (c *Client) GetFlowActions() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flowcardaction/", nil)
}

// Zones

func (c *Client) GetZones() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/zones/zone/", nil)
}

func (c *Client) GetZone(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/zones/zone/"+id, nil)
}

func (c *Client) CreateZone(zone map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/zones/zone/", zone)
}

func (c *Client) DeleteZone(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/zones/zone/"+id, nil)
	return err
}

func (c *Client) UpdateZone(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/zones/zone/"+id, updates)
	return err
}

// Apps

func (c *Client) GetApps() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/apps/app/", nil)
}

func (c *Client) GetApp(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/apps/app/"+id, nil)
}

func (c *Client) RestartApp(id string) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/manager/apps/app/%s/restart", id), nil)
	return err
}

// Notifications

func (c *Client) SendNotification(text string) error {
	// Use flow card action to create notification
	body := map[string]interface{}{
		"args": map[string]string{"text": text},
	}
	_, err := c.doRequest("POST", "/api/manager/flow/flowcardaction/homey:manager:notifications/homey:manager:notifications:create_notification/run", body)
	return err
}

func (c *Client) GetNotifications() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/notifications/notification/", nil)
}

// RunFlowCardAction runs any flow card action
func (c *Client) RunFlowCardAction(uri, id string, args map[string]interface{}) (json.RawMessage, error) {
	body := map[string]interface{}{
		"args": args,
	}
	return c.doRequest("POST", fmt.Sprintf("/api/manager/flow/flowcardaction/%s/%s/run", uri, id), body)
}

// Logic variables

func (c *Client) GetVariables() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/logic/variable/", nil)
}

func (c *Client) GetVariable(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/logic/variable/"+id, nil)
}

func (c *Client) SetVariable(id string, value interface{}) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", "/api/manager/logic/variable/"+id, body)
	return err
}

func (c *Client) CreateVariable(name string, varType string, value interface{}) (json.RawMessage, error) {
	body := map[string]interface{}{
		"name":  name,
		"type":  varType,
		"value": value,
	}
	return c.doRequest("POST", "/api/manager/logic/variable/", body)
}

func (c *Client) DeleteVariable(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/logic/variable/"+id, nil)
	return err
}

// System

func (c *Client) GetSystem() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/system/", nil)
}

func (c *Client) Reboot() error {
	_, err := c.doRequest("POST", "/api/manager/system/reboot/", nil)
	return err
}

// Users

func (c *Client) GetUsers() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/users/user/", nil)
}

// Insights (logs/history)

func (c *Client) GetInsights() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/insights/log/", nil)
}

func (c *Client) GetInsightEntries(uri, id, resolution string) (json.RawMessage, error) {
	// URL encode the URI and ID since they contain colons
	encodedURI := url.PathEscape(uri)
	encodedID := url.PathEscape(id)
	path := fmt.Sprintf("/api/manager/insights/log/%s/%s/entry", encodedURI, encodedID)
	if resolution != "" {
		path += "?resolution=" + resolution
	}
	return c.doRequest("GET", path, nil)
}

// Energy

func (c *Client) GetEnergyLive() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/live", nil)
}

func (c *Client) GetEnergyReportDay(date string) (json.RawMessage, error) {
	path := "/api/manager/energy/report/day"
	if date != "" {
		path += "?date=" + date
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) GetEnergyReportWeek(isoWeek string) (json.RawMessage, error) {
	path := "/api/manager/energy/report/week"
	if isoWeek != "" {
		path += "?isoWeek=" + isoWeek
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) GetEnergyReportMonth(yearMonth string) (json.RawMessage, error) {
	path := "/api/manager/energy/report/month"
	if yearMonth != "" {
		path += "?yearMonth=" + yearMonth
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) GetEnergyReportsAvailable() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/reports/available", nil)
}

func (c *Client) GetElectricityPrice(date string) (json.RawMessage, error) {
	path := "/api/manager/energy/price/electricity/dynamic"
	if date != "" {
		path += "?date=" + date
	}
	return c.doRequest("GET", path, nil)
}

func (c *Client) GetElectricityPriceType() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/price/electricity/type", nil)
}

func (c *Client) SetElectricityPriceType(priceType string) error {
	_, err := c.doRequest("PUT", "/api/manager/energy/price/electricity/"+priceType, nil)
	return err
}

func (c *Client) GetElectricityPriceFixed() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/option/electricityPriceFixed", nil)
}

func (c *Client) SetElectricityPriceFixed(price float64) error {
	body := map[string]interface{}{
		"value": map[string]interface{}{
			"costs": map[string]interface{}{
				"user_fixed_base": map[string]interface{}{
					"value": price,
				},
			},
		},
	}
	_, err := c.doRequest("PUT", "/api/manager/energy/option/electricityPriceFixed", body)
	return err
}

// Personal Access Tokens (PAT)

func (c *Client) ListPATs() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/users/pat", nil)
}

func (c *Client) CreatePAT(name string, scopes []string) (json.RawMessage, error) {
	body := map[string]interface{}{
		"name":   name,
		"scopes": scopes,
	}
	return c.doRequest("POST", "/api/manager/users/pat", body)
}

func (c *Client) DeletePAT(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/users/pat/"+id, nil)
	return err
}

// Flow Folders

func (c *Client) GetFlowFolders() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flowfolder/", nil)
}

func (c *Client) GetFlowFolder(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/flow/flowfolder/"+id, nil)
}

func (c *Client) CreateFlowFolder(folder map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/flow/flowfolder/", folder)
}

func (c *Client) UpdateFlowFolder(id string, folder map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/flow/flowfolder/"+id, folder)
	return err
}

func (c *Client) DeleteFlowFolder(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/flow/flowfolder/"+id, nil)
	return err
}

// Apps (extended)

func (c *Client) InstallApp(appID string, channel string) (json.RawMessage, error) {
	body := map[string]interface{}{
		"id":             appID,
		"waitForInstall": true,
	}
	if channel != "" {
		body["channel"] = channel
	}
	return c.doRequest("POST", "/api/manager/apps/store", body)
}

func (c *Client) UninstallApp(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/apps/app/"+id, nil)
	return err
}

func (c *Client) EnableApp(id string) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/apps/app/%s/enable", id), nil)
	return err
}

func (c *Client) DisableApp(id string) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/apps/app/%s/disable", id), nil)
	return err
}

func (c *Client) UpdateApp(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/apps/app/"+id, updates)
	return err
}

func (c *Client) GetAppSettings(id string) (json.RawMessage, error) {
	return c.doRequest("GET", fmt.Sprintf("/api/manager/apps/app/%s/setting", id), nil)
}

func (c *Client) SetAppSetting(appID, settingName string, value interface{}) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/apps/app/%s/setting/%s", appID, settingName), body)
	return err
}

func (c *Client) GetAppUsage(id string) (json.RawMessage, error) {
	return c.doRequest("GET", fmt.Sprintf("/api/manager/apps/app/%s/usage", id), nil)
}

// Users (extended)

func (c *Client) GetUser(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/users/user/"+id, nil)
}

func (c *Client) GetUserMe() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/users/user/me", nil)
}

func (c *Client) CreateUser(user map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/users/user/", user)
}

func (c *Client) UpdateUser(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/users/user/"+id, updates)
	return err
}

func (c *Client) DeleteUser(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/users/user/"+id, nil)
	return err
}

// Moods

func (c *Client) GetMoods() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/moods/mood/", nil)
}

func (c *Client) GetMood(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/moods/mood/"+id, nil)
}

func (c *Client) CreateMood(mood map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/moods/mood/", mood)
}

func (c *Client) UpdateMood(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/moods/mood/"+id, updates)
	return err
}

func (c *Client) DeleteMood(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/moods/mood/"+id, nil)
	return err
}

func (c *Client) SetMood(id string) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/api/manager/moods/mood/%s/set", id), nil)
	return err
}

// Dashboards

func (c *Client) GetDashboards() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/dashboards/dashboard/", nil)
}

func (c *Client) GetDashboard(id string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/dashboards/dashboard/"+id, nil)
}

func (c *Client) CreateDashboard(dashboard map[string]interface{}) (json.RawMessage, error) {
	return c.doRequest("POST", "/api/manager/dashboards/dashboard/", dashboard)
}

func (c *Client) UpdateDashboard(id string, updates map[string]interface{}) error {
	_, err := c.doRequest("PUT", "/api/manager/dashboards/dashboard/"+id, updates)
	return err
}

func (c *Client) DeleteDashboard(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/dashboards/dashboard/"+id, nil)
	return err
}

// Presence

func (c *Client) GetPresent(userID string) (json.RawMessage, error) {
	return c.doRequest("GET", fmt.Sprintf("/api/manager/presence/%s/present", userID), nil)
}

func (c *Client) SetPresent(userID string, value bool) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/presence/%s/present", userID), body)
	return err
}

func (c *Client) SetPresentMe(value bool) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", "/api/manager/presence/me/present", body)
	return err
}

func (c *Client) GetAsleep(userID string) (json.RawMessage, error) {
	return c.doRequest("GET", fmt.Sprintf("/api/manager/presence/%s/asleep", userID), nil)
}

func (c *Client) SetAsleep(userID string, value bool) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", fmt.Sprintf("/api/manager/presence/%s/asleep", userID), body)
	return err
}

func (c *Client) SetAsleepMe(value bool) error {
	body := map[string]interface{}{"value": value}
	_, err := c.doRequest("PUT", "/api/manager/presence/me/asleep", body)
	return err
}

// Notifications (extended)

func (c *Client) DeleteNotification(id string) error {
	_, err := c.doRequest("DELETE", "/api/manager/notifications/notification/"+id, nil)
	return err
}

func (c *Client) DeleteAllNotifications() error {
	_, err := c.doRequest("DELETE", "/api/manager/notifications/notification/", nil)
	return err
}

func (c *Client) GetNotificationOwners() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/notifications/owner/", nil)
}

// Insights (extended)

func (c *Client) DeleteInsightLog(uri, id string) error {
	encodedURI := url.PathEscape(uri)
	encodedID := url.PathEscape(id)
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/manager/insights/log/%s/%s", encodedURI, encodedID), nil)
	return err
}

func (c *Client) DeleteInsightLogEntries(uri, id string) error {
	encodedURI := url.PathEscape(uri)
	encodedID := url.PathEscape(id)
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/manager/insights/log/%s/%s/entry", encodedURI, encodedID), nil)
	return err
}

// System (extended)

func (c *Client) GetSystemName() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/system/name", nil)
}

func (c *Client) SetSystemName(name string) error {
	body := map[string]interface{}{"name": name}
	_, err := c.doRequest("PUT", "/api/manager/system/name", body)
	return err
}

// Weather

func (c *Client) GetWeather() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/weather/weather", nil)
}

func (c *Client) GetWeatherForecast() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/weather/forecast/hourly", nil)
}

// Energy (extended)

func (c *Client) GetEnergyReportYear(year string) (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/report/year?year="+year, nil)
}

func (c *Client) DeleteEnergyReports() error {
	_, err := c.doRequest("DELETE", "/api/manager/energy/reports", nil)
	return err
}

func (c *Client) GetEnergyCurrency() (json.RawMessage, error) {
	return c.doRequest("GET", "/api/manager/energy/currency", nil)
}
