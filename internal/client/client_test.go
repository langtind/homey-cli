package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateDevice(t *testing.T) {
	// Create a test server
	var receivedPath string
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	// Create client with test server
	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	// Test UpdateDevice
	deviceID := "test-device-123"
	updates := map[string]interface{}{
		"name": "New Device Name",
	}

	err := client.UpdateDevice(deviceID, updates)
	if err != nil {
		t.Fatalf("UpdateDevice failed: %v", err)
	}

	// Verify the request
	expectedPath := "/api/manager/devices/device/test-device-123"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}

	if receivedBody["name"] != "New Device Name" {
		t.Errorf("expected name 'New Device Name', got %v", receivedBody["name"])
	}
}

func TestUpdateDevice_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "device not found"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.UpdateDevice("nonexistent", map[string]interface{}{"name": "Test"})
	if err == nil {
		t.Error("expected error for nonexistent device")
	}
}

func TestUpdateZone(t *testing.T) {
	var receivedPath string
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	zoneID := "test-zone-123"
	updates := map[string]interface{}{
		"name": "New Zone Name",
	}

	err := client.UpdateZone(zoneID, updates)
	if err != nil {
		t.Fatalf("UpdateZone failed: %v", err)
	}

	expectedPath := "/api/manager/zones/zone/test-zone-123"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}

	if receivedBody["name"] != "New Zone Name" {
		t.Errorf("expected name 'New Zone Name', got %v", receivedBody["name"])
	}
}

func TestUpdateZone_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "zone not found"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.UpdateZone("nonexistent", map[string]interface{}{"name": "Test"})
	if err == nil {
		t.Error("expected error for nonexistent zone")
	}
}

// Tests for Flow Folder API methods
func TestCreateFlowFolder(t *testing.T) {
	var receivedPath string
	var receivedBody map[string]interface{}
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "folder-123", "name": "Test Folder"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	folder := map[string]interface{}{
		"name": "Test Folder",
	}

	_, err := client.CreateFlowFolder(folder)
	if err != nil {
		t.Fatalf("CreateFlowFolder failed: %v", err)
	}

	if receivedMethod != "POST" {
		t.Errorf("expected POST method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/flow/flowfolder/"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}

	if receivedBody["name"] != "Test Folder" {
		t.Errorf("expected name 'Test Folder', got %v", receivedBody["name"])
	}
}

func TestDeleteFlowFolder(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.DeleteFlowFolder("folder-123")
	if err != nil {
		t.Fatalf("DeleteFlowFolder failed: %v", err)
	}

	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/flow/flowfolder/folder-123"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for Mood API methods
func TestGetMoods(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "mood-1", "name": "Relaxed"}]`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetMoods()
	if err != nil {
		t.Fatalf("GetMoods failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/moods/mood/"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestSetMood(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.SetMood("mood-123")
	if err != nil {
		t.Fatalf("SetMood failed: %v", err)
	}

	if receivedMethod != "POST" {
		t.Errorf("expected POST method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/moods/mood/mood-123/set"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for Dashboard API methods
func TestGetDashboards(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"dash-1": {"id": "dash-1", "name": "Main"}}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetDashboards()
	if err != nil {
		t.Fatalf("GetDashboards failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/dashboards/dashboard/"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for Presence API methods
func TestSetPresent(t *testing.T) {
	var receivedPath string
	var receivedMethod string
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.SetPresent("user-123", true)
	if err != nil {
		t.Fatalf("SetPresent failed: %v", err)
	}

	if receivedMethod != "PUT" {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/presence/user-123/present"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}

	if receivedBody["value"] != true {
		t.Errorf("expected value true, got %v", receivedBody["value"])
	}
}

func TestSetAsleep(t *testing.T) {
	var receivedPath string
	var receivedMethod string
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.SetAsleep("user-123", true)
	if err != nil {
		t.Fatalf("SetAsleep failed: %v", err)
	}

	if receivedMethod != "PUT" {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/presence/user-123/asleep"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}

	if receivedBody["value"] != true {
		t.Errorf("expected value true, got %v", receivedBody["value"])
	}
}

// Tests for Weather API methods
func TestGetWeather(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"state": "sunny", "temperature": 22.5}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetWeather()
	if err != nil {
		t.Fatalf("GetWeather failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/weather/weather"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestGetWeatherForecast(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"time": "2025-01-16T12:00:00", "state": "cloudy"}]`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetWeatherForecast()
	if err != nil {
		t.Fatalf("GetWeatherForecast failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/weather/forecast/hourly"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for Notification API methods
func TestDeleteNotification(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.DeleteNotification("notif-123")
	if err != nil {
		t.Fatalf("DeleteNotification failed: %v", err)
	}

	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/notifications/notification/notif-123"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestDeleteAllNotifications(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.DeleteAllNotifications()
	if err != nil {
		t.Fatalf("DeleteAllNotifications failed: %v", err)
	}

	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/notifications/notification/"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for System Name API methods
func TestGetSystemName(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"My Homey"`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetSystemName()
	if err != nil {
		t.Fatalf("GetSystemName failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/system/name"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestSetSystemName(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.SetSystemName("New Homey Name")
	if err != nil {
		t.Fatalf("SetSystemName failed: %v", err)
	}

	if receivedMethod != "PUT" {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/system/name"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for Energy API methods
func TestGetEnergyCurrency(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"NOK"`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	_, err := client.GetEnergyCurrency()
	if err != nil {
		t.Fatalf("GetEnergyCurrency failed: %v", err)
	}

	if receivedMethod != "GET" {
		t.Errorf("expected GET method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/energy/currency"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestDeleteEnergyReports(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.DeleteEnergyReports()
	if err != nil {
		t.Fatalf("DeleteEnergyReports failed: %v", err)
	}

	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/energy/reports"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

// Tests for App management API methods
func TestUninstallApp(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.UninstallApp("com.test.app")
	if err != nil {
		t.Fatalf("UninstallApp failed: %v", err)
	}

	if receivedMethod != "DELETE" {
		t.Errorf("expected DELETE method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/apps/app/com.test.app"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestEnableApp(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.EnableApp("com.test.app")
	if err != nil {
		t.Fatalf("EnableApp failed: %v", err)
	}

	if receivedMethod != "PUT" {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/apps/app/com.test.app/enable"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}

func TestDisableApp(t *testing.T) {
	var receivedPath string
	var receivedMethod string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	err := client.DisableApp("com.test.app")
	if err != nil {
		t.Fatalf("DisableApp failed: %v", err)
	}

	if receivedMethod != "PUT" {
		t.Errorf("expected PUT method, got %s", receivedMethod)
	}

	expectedPath := "/api/manager/apps/app/com.test.app/disable"
	if receivedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, receivedPath)
	}
}
