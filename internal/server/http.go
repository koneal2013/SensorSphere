package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/koneal2013/sensorsphere/internal/db"
	"github.com/koneal2013/sensorsphere/internal/middleware/adaptor"
	"github.com/koneal2013/sensorsphere/internal/models"
)

type HttpConfig struct {
	Port            int
	MiddlewareFuncs []mux.MiddlewareFunc
	Db              db.Database
}

type SensorSphere struct {
	HttpTracer trace.Tracer
	database   db.Database
}

func NewHTTPServer(cfg *HttpConfig) (*http.Server, error) {
	s := &SensorSphere{
		HttpTracer: otel.GetTracerProvider().Tracer("httpTracer"),
		database:   cfg.Db,
	}
	r := mux.NewRouter()
	r.HandleFunc("/sensors", adaptor.GenericHttpAdaptor(s.HandleCreateSensor)).Methods(http.MethodPost)
	r.HandleFunc("/sensors/nearest", adaptor.GenericHttpAdaptor(s.HandleGetNearestSensor)).Methods(http.MethodGet)
	r.HandleFunc("/sensor_readings",
		adaptor.GenericHttpAdaptor(s.HandleGetSensorReadingsForTimeRange)).Methods(http.MethodGet)
	r.HandleFunc("/status", s.HandleStatus).Methods(http.MethodGet)
	r.HandleFunc("/sensor_readings",
		adaptor.GenericHttpAdaptor(s.HandleCreateSensorReading)).Methods(http.MethodPost)
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/",
		http.FileServer(http.Dir("./cmd/sensorsphere/docs"))))
	r.HandleFunc("/sensors/{name}", adaptor.GenericHttpAdaptor(s.HandleGetSensor)).Methods(http.MethodGet)
	r.HandleFunc("/sensors/{name}", adaptor.GenericHttpAdaptor(s.HandleUpdateSensor)).Methods(http.MethodPut)
	r.Use(cfg.MiddlewareFuncs...)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}, nil
}

// @Summary Get server status
// @Description Returns 200 OK if server is ready to accept requests
// @Tags status
// @Produce  text/plain
// @Success 200 {string} string "Server is running"
// @Router /status [get]
func (s *SensorSphere) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Server is running")
}

// @Summary Create a new sensor
// @Description Create a new sensor with the input payload
// @Tags sensors
// @Accept  json
// @Produce  json
// @Param sensor body models.Sensor true "Create sensor"
// @Success 200 {object} models.Sensor
// @Router /sensors [post]
func (s *SensorSphere) HandleCreateSensor(ctx context.Context, in models.Sensor) (*models.Sensor, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleCreateSensor")
	defer span.End()

	if in.Name == "" || in.Location.Latitude == 0.0 || in.Location.Longitude == 0.0 || in.Tags == nil {
		return nil, fmt.Errorf("missing required fields")
	}

	sensor, err := s.database.CreateSensor(ctx, &in)
	if err != nil {
		return &models.Sensor{}, err
	}

	return sensor, nil
}

// @Summary Get a sensor
// @Description Get a sensor by its name
// @Tags sensors
// @Accept  json
// @Produce  json
// @Param name path string true "Sensor name"
// @Success 200 {object} models.Sensor
// @Router /sensors/{name} [get]
func (s *SensorSphere) HandleGetSensor(ctx context.Context, in map[string]string) (*models.Sensor, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleGetSensor")
	defer span.End()

	sensorName, ok := in["name"]
	if !ok {
		return nil, fmt.Errorf("missing required fields")
	}

	sensor, err := s.database.GetSensor(ctx, sensorName)
	if err != nil {
		return nil, err
	}

	return sensor, nil
}

// @Summary Get sensor readings for a time range
// @Description Get sensor readings for a specific time range
// @Tags sensor_readings
// @Accept  json
// @Produce  json
// @Param timeRangeQuery body models.TimeRangeQuery true "Time range query"
// @Success 200 {array} models.SensorReading
// @Router /sensor_readings [get]
func (s *SensorSphere) HandleGetSensorReadingsForTimeRange(ctx context.Context,
	in models.TimeRangeQuery) ([]*models.SensorReading, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleGetSensorReadingsForTimeRange")
	defer span.End()

	if in.SensorName == "" || in.EndTime.IsZero() || in.StartTime.IsZero() {
		return nil, fmt.Errorf("missing required fields")
	}

	sensorReadings, err := s.database.GetSensorReadingsForTimeRange(ctx, in)
	if err != nil {
		return nil, err
	}

	return sensorReadings, nil
}

// @Summary Update a sensor
// @Description Update a sensor with the input payload
// @Tags sensors
// @Accept  json
// @Produce  json
// @Param sensor body models.Sensor true "Update sensor"
// @Success 200 {integer} int64
// @Router /sensors/{name} [put]
func (s *SensorSphere) HandleUpdateSensor(ctx context.Context, in models.Sensor) (int64, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleUpdateSensor")
	defer span.End()

	if in.Name == "" || in.Location.Latitude == 0.0 || in.Location.Longitude == 0.0 || in.Tags == nil {
		return 0, fmt.Errorf("missing required fields")
	}

	rows, err := s.database.UpdateSensor(ctx, &in)
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// @Summary Get the nearest sensor
// @Description Get the nearest sensor to a specific location
// @Tags sensors
// @Accept  json
// @Produce  json
// @Param location body models.Location true "Location"
// @Success 200 {object} models.Sensor
// @Router /sensors/nearest [get]
func (s *SensorSphere) HandleGetNearestSensor(ctx context.Context, in models.Location) (*models.Sensor, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleGetNearestSensor")
	defer span.End()

	if in.Latitude == 0.0 || in.Longitude == 0.0 {
		return nil, fmt.Errorf("missing required fields")
	}

	sensor, err := s.database.GetNearestSensor(ctx, &in)
	if err != nil {
		return &models.Sensor{}, err
	}

	return sensor, nil
}

// @Summary Create a new sensor reading
// @Description Create a new sensor reading with the input payload
// @Tags sensor_readings
// @Accept  json
// @Produce  json
// @Param sensorReading body models.SensorReading true "Create sensor reading"
// @Success 200 {object} models.SensorReading
// @Router /sensor_readings [post]
func (s *SensorSphere) HandleCreateSensorReading(ctx context.Context,
	reading models.SensorReading,
) (*models.SensorReading, error) {
	ctx, span := s.HttpTracer.Start(ctx, "HandleCreateSensorReading")
	defer span.End()

	if reading.SensorName == "" || reading.Value == 0.0 {
		return nil, fmt.Errorf("missing required fields")
	}

	sensorReading, err := s.database.CreateSensorReading(ctx, &reading)
	if err != nil {
		return nil, err
	}

	return sensorReading, nil
}
