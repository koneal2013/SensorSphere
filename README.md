# SensorSphere

SensorSphere is a project that provides a GRPC and HTTP server implementation for a sensor data management system. The system allows for the creation, retrieval, and updating of sensor data. It also provides functionality for retrieving sensor readings within a specific time range and finding the nearest sensor to a specific location.

The project uses Go as the primary language and leverages several libraries and frameworks such as Gorilla Mux for HTTP routing, OpenTelemetry for distributed tracing, and gRPC for efficient, high-performance communication.

## Getting Started

These instructions will help you get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Go (version 1.16 or later)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/koneal2013/SensorSphere.git
```

2. Navigate to the project directory:

```bash
cd SensorSphere
```

3. Build the Docker images:

```bash
docker-compose build
```

### Running the Application

1. Start the services:

```bash
docker-compose up
```

The application should now be running and accessible on the specified ports.

## Usage

The application provides several endpoints for managing sensor data:

- `POST /sensors`: Create a new sensor.
- `GET /sensors/{name}`: Get a sensor by its name.
- `GET /sensor_readings`: Get sensor readings for a specific time range.
- `PUT /sensors/{name}`: Update a sensor.
- `GET /sensors/nearest`: Get the nearest sensor to a specific location.
- `POST /sensor_readings`: Create a new sensor reading.

## Documentation

The project includes Swagger documentation for its HTTP API. You can access the Swagger UI at `http://localhost:8080/swagger/` when the application is running.

## Tests

The project includes unit tests for the database functions. You can run the tests using the following command:

```bash
go test ./...
```

## Contributing

Contributions are welcome. Please feel free to open an issue or submit a pull request.

## Questions

If you have any questions about the project, please feel free to [email me](mailto:koneal2013@gmail.com).

