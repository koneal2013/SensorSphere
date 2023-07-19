package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"

	"github.com/koneal2013/sensorsphere/internal/models"
)

const dbDriverName = "postgres"

type PgConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (c PgConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)
}

type Database interface {
	CreateSensor(ctx context.Context, newSensor *models.Sensor) (*models.Sensor, error)
	GetSensor(ctx context.Context, sensorName string) (*models.Sensor, error)
	UpdateSensor(ctx context.Context, updatedSensor *models.Sensor) (int64, error)
	GetNearestSensor(ctx context.Context, location *models.Location) (*models.Sensor, error)
	CreateSensorReading(ctx context.Context, reading *models.SensorReading) (*models.SensorReading, error)
	GetSensorReadingsForTimeRange(ctx context.Context,
		timeRange models.TimeRangeQuery) ([]*models.SensorReading, error)
	Close() error
	RunMigrations() error
}

type Db struct {
	*sql.DB
}

func New(config PgConfig) (*Db, error) {
	db, err := sql.Open(dbDriverName, config.ConnectionString())
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

func (d *Db) RunMigrations() error {
	err := goose.SetDialect(dbDriverName)
	if err != nil {
		return err
	}

	err = goose.Up(d.DB, "/migrations")
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) Close() error {
	return d.DB.Close()
}

func (d *Db) CreateSensor(ctx context.Context, newSensor *models.Sensor) (*models.Sensor, error) {
	sqlStatement := `
		INSERT INTO sensors (name, location, tags)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4);`

	_, err := d.ExecContext(ctx, sqlStatement, newSensor.Name, newSensor.Location.Longitude,
		newSensor.Location.Latitude, pq.Array(newSensor.Tags))
	if err != nil {
		return nil, err
	}

	return newSensor, nil
}

func (d *Db) GetSensor(ctx context.Context, sensorName string) (*models.Sensor, error) {
	sqlStatement := `
		SELECT name, ST_AsText(location), tags
		FROM sensors
		WHERE name = $1;`

	row := d.QueryRowContext(ctx, sqlStatement, sensorName)

	var sensor models.Sensor

	var location string

	err := row.Scan(&sensor.Name, &location, pq.Array(&sensor.Tags))
	if err != nil {
		return nil, err
	}

	// Parse location
	geometry, err := wkt.Unmarshal(location)
	if err != nil {
		return nil, err
	}

	point, ok := geometry.(*geom.Point)
	if !ok {
		return nil, fmt.Errorf("location is not a point")
	}

	sensor.Location = models.Location{
		Longitude: point.X(),
		Latitude:  point.Y(),
	}

	return &sensor, nil
}

func (d *Db) UpdateSensor(ctx context.Context, updatedSensor *models.Sensor) (int64, error) {
	sqlStatement := `
		UPDATE sensors
		SET  location = ST_SetSRID(ST_MakePoint($2, $3), 4326), tags = $4
		WHERE name = $1;`

	res, err := d.ExecContext(ctx, sqlStatement, updatedSensor.Name, updatedSensor.Location.Longitude,
		updatedSensor.Location.Latitude, pq.Array(updatedSensor.Tags))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, err
}

func (d *Db) GetNearestSensor(ctx context.Context, location *models.Location) (*models.Sensor, error) {
	sqlStatement := `
		SELECT name, ST_AsText(location), tags
		FROM sensors
		ORDER BY location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)
		LIMIT 1;`

	row := d.QueryRowContext(ctx, sqlStatement, location.Longitude, location.Latitude)

	var sensor models.Sensor

	var loc string

	err := row.Scan(&sensor.Name, &loc, pq.Array(&sensor.Tags))
	if err != nil {
		return nil, err
	}

	// Parse location
	geometry, err := wkt.Unmarshal(loc)
	if err != nil {
		return nil, err
	}

	point, ok := geometry.(*geom.Point)
	if !ok {
		return nil, fmt.Errorf("location is not a point")
	}

	sensor.Location = models.Location{
		Longitude: point.X(),
		Latitude:  point.Y(),
	}

	return &sensor, nil
}

func (d *Db) GetSensorReadingsForTimeRange(ctx context.Context,
	timeRange models.TimeRangeQuery) ([]*models.SensorReading, error) {
	sqlStatement := `
		SELECT name, value, time
		FROM sensor_readings
		WHERE name = $1 AND time BETWEEN $2 AND $3;`

	rows, err := d.QueryContext(ctx, sqlStatement, timeRange.SensorName, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sensorReadings := []*models.SensorReading{}

	for rows.Next() {
		var sensorReading models.SensorReading

		err = rows.Scan(&sensorReading.SensorName, &sensorReading.Value, &sensorReading.Time)
		if err != nil {
			return nil, err
		}

		sensorReadings = append(sensorReadings, &sensorReading)
	}

	return sensorReadings, nil
}

func (d *Db) CreateSensorReading(ctx context.Context, reading *models.SensorReading) (*models.SensorReading, error) {
	sqlStatement := `
		INSERT INTO sensor_readings (name, value, time)
		VALUES ($1, $2, NOW())
		RETURNING time;`

	row := d.QueryRowContext(ctx, sqlStatement, reading.SensorName, reading.Value)

	err := row.Scan(&reading.Time)
	if err != nil {
		return nil, err
	}

	return reading, nil
}
