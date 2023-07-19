-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE TABLE IF NOT EXISTS sensors (
                                       name TEXT PRIMARY KEY NOT NULL,
                                       location GEOGRAPHY(Point, 4326),
                                       tags TEXT[] NOT NULL
);
CREATE TABLE IF NOT EXISTS sensor_readings (
                                       name TEXT REFERENCES sensors NOT NULL,
                                       value DOUBLE PRECISION NOT NULL,
                                       time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ON sensor_readings (time, name);
SELECT create_hypertable('sensor_readings', 'time');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sensor_readings;
DROP TABLE IF EXISTS sensor;
DROP EXTENSION IF EXISTS "postgis";
DROP EXTENSION IF EXISTS timescaledb;
-- +goose StatementEnd
