package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/koneal2013/sensorsphere/internal/models"
)

// MockDb is a mock type for db.Db
type MockDb struct {
	mock.Mock
}

// CreateSensor is a mock implementation of db.Db.CreateSensor
func (m *MockDb) CreateSensor(ctx context.Context, sensor *models.Sensor) (*models.Sensor, error) {
	args := m.Called(ctx, sensor)

	return args.Get(0).(*models.Sensor), args.Error(1)
}

// GetSensor is a mock implementation of db.Db.GetSensor
func (m *MockDb) GetSensor(ctx context.Context, name string) ([]*models.Sensor, error) {
	args := m.Called(ctx, name)

	return args.Get(0).([]*models.Sensor), args.Error(1)
}

// UpdateSensor is a mock implementation of db.Db.UpdateSensor
func (m *MockDb) UpdateSensor(ctx context.Context, sensor *models.Sensor) (int64, error) {
	args := m.Called(ctx, sensor)

	return args.Get(0).(int64), args.Error(1)
}

// GetNearestSensor is a mock implementation of db.Db.GetNearestSensor
func (m *MockDb) GetNearestSensor(ctx context.Context, location *models.Location) (*models.Sensor, error) {
	args := m.Called(ctx, location)

	return args.Get(0).(*models.Sensor), args.Error(1)
}

// GetSensorReadingsForTimeRange is a mock implementation of db.Db.GetSensorReadingsForTimeRange
func (m *MockDb) GetSensorReadingsForTimeRange(ctx context.Context,
	timeRange models.TimeRangeQuery) ([]*models.SensorReading, error) {
	args := m.Called(ctx, timeRange)

	return args.Get(0).([]*models.SensorReading), args.Error(1)
}

// CreateSensorReading is a mock implementation of db.Db.CreateSensorReading
func (m *MockDb) CreateSensorReading(ctx context.Context,
	reading *models.SensorReading) (*models.SensorReading, error) {
	args := m.Called(ctx, reading)

	return args.Get(0).(*models.SensorReading), args.Error(1)
}

func TestHandleCreateSensor(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new sensor
	sensor := models.Sensor{Name: "Test Sensor"}

	// Setup expectations
	mockDB.On("CreateSensor", mock.Anything, &sensor).Return(&sensor, nil)

	// Convert the sensor to JSON
	jsonSensor, _ := json.Marshal(sensor)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/sensors", bytes.NewBuffer(jsonSensor))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `{"name":"Test Sensor","location":{"longitude":0,"latitude":0},"tags":null}
`
	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestHandleGetSensor(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new sensor
	sensor := models.Sensor{Name: "Test Sensor"}

	// Setup expectations
	mockDB.On("GetSensor", mock.Anything, sensor.Name).Return([]*models.Sensor{&sensor}, nil)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodGet, "/sensors/"+sensor.Name, io.NopCloser(bytes.NewReader(nil)))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `[{"name":"Test Sensor","location":{"longitude":0,"latitude":0},"tags":null}]
`
	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestHandleUpdateSensor(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new sensor
	sensor := models.Sensor{Name: "Test Sensor"}

	// Setup expectations
	mockDB.On("UpdateSensor", mock.Anything, &sensor).Return(int64(1), nil)

	// Convert the sensor to JSON
	jsonSensor, _ := json.Marshal(sensor)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPut, "/sensors/"+sensor.Name, bytes.NewBuffer(jsonSensor))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `1
`
	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestHandleGetNearestSensor(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new location
	location := models.Location{Longitude: 0, Latitude: 0}

	// Create a new sensor
	sensor := models.Sensor{Name: "Test Sensor"}

	// Setup expectations
	mockDB.On("GetNearestSensor", mock.Anything, &location).Return(&sensor, nil)

	// Convert the location to JSON
	jsonLocation, _ := json.Marshal(location)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodGet, "/sensors/nearest", bytes.NewBuffer(jsonLocation))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `{"name":"Test Sensor","location":{"longitude":0,"latitude":0},"tags":null}
`
	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestHandleCreateSensorReading(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new sensor reading
	reading := models.SensorReading{SensorName: "Test Sensor", Value: 1.0}

	// Setup expectations
	mockDB.On("CreateSensorReading", mock.Anything, &reading).Return(&reading, nil)

	// Convert the reading to JSON
	jsonReading, _ := json.Marshal(reading)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPost, "/sensor_readings", bytes.NewBuffer(jsonReading))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `{"sensorName":"Test Sensor","time":"0001-01-01T00:00:00Z","value":1}
`
	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestHandleGetSensorReadingsForTimeRange(t *testing.T) {
	// Create a new instance of our mock Db
	mockDB := new(MockDb)

	// Create a new HTTP server with the mock database
	server := NewHTTPServer(&HttpConfig{Port: 8080}, mockDB)

	// Create a new time range query
	timeRangeQuery := models.TimeRangeQuery{
		SensorName: "Test Sensor",
		StartTime:  time.Now().Add(-1 * time.Hour).Truncate(time.Second),
		EndTime:    time.Now().Truncate(time.Second),
	}

	// Create a slice of sensor readings
	sensorReadings := []*models.SensorReading{
		{SensorName: "Test Sensor", Value: 1.0, Time: time.Now().Add(-30 * time.Minute)},
		{SensorName: "Test Sensor", Value: 2.0, Time: time.Now().Add(-15 * time.Minute)},
	}

	// Setup expectations
	mockDB.On("GetSensorReadingsForTimeRange", mock.Anything, timeRangeQuery).Return(sensorReadings, nil)

	// Convert the time range query to JSON
	jsonTimeRangeQuery, _ := json.Marshal(timeRangeQuery)

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodGet, "/sensor_readings", bytes.NewBuffer(jsonTimeRangeQuery))

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	server.Handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := fmt.Sprintf(
		"[{\"sensorName\":\"Test Sensor\",\"time\":\"%s\",\"value\":1},{\"sensorName\":\"Test Sensor\",\"time\":\"%s\",\"value\":2}]\n",
		sensorReadings[0].Time.Format(time.RFC3339Nano),
		sensorReadings[1].Time.Format(time.RFC3339Nano),
	)

	assert.Equal(t, expected, rr.Body.String())

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}
