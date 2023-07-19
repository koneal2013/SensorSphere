-- +goose Up
-- +goose StatementBegin
-- Generate sensors
DO
$DO$
    DECLARE
        counter INTEGER := 1;
        sensor_name TEXT;
        sensor_location GEOMETRY;
        sensor_tags TEXT[];
        reading_time TIMESTAMPTZ;
        reading_value DOUBLE PRECISION;
    BEGIN
        WHILE counter <= 100 LOOP
                sensor_name := 'sensor' || counter;
                sensor_location := 'POINT(' || (random()*360 - 180)::text || ' ' || (random()*180 - 90)::text || ')';
                sensor_tags := ARRAY['tag' || counter, 'value' || counter];

                -- Insert into sensors table
                INSERT INTO sensors (name, location, tags) VALUES (sensor_name, sensor_location, sensor_tags);

                -- Insert into sensor_readings table
                FOR i IN 1..100 LOOP
                        reading_time := NOW() - ((random() * 86400)::int || ' seconds')::interval - ((random() * 1000 + i)::int || ' days')::interval;
                        reading_value := random()*100;

                        BEGIN
                            INSERT INTO sensor_readings (name, time, value) VALUES (sensor_name, reading_time, reading_value);
                        END;
                    END LOOP;

                counter := counter + 1;
        END LOOP;
    END;
$DO$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete sensor readings associated with generated sensors
DELETE FROM sensor_readings
WHERE name IN (
    SELECT name FROM sensors WHERE name LIKE 'sensor_%'
);

-- Delete generated sensors
DELETE FROM sensors
WHERE name LIKE 'sensor_%';
-- +goose StatementEnd
