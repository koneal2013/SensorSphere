package models

import (
	"time"
)

type Sensor struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Tags     []string `json:"tags"`
}

type SensorReading struct {
	SensorName string    `json:"sensorName"`
	Time       time.Time `json:"time"`
	Value      float64   `json:"value"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type TimeRangeQuery struct {
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	SensorName string    `json:"sensorName"`
}
