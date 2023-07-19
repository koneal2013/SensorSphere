# SensorSphere

SensorSphere is a project that provides a GRPC and HTTP server implementation for a sensor data management system. The system allows for the creation, retrieval, and updating of sensor data. It also provides functionality for retrieving sensor readings within a specific time range and finding the nearest sensor to a specific location.

The project uses Go as the primary language and leverages several libraries and frameworks such as Gorilla Mux for HTTP routing, OpenTelemetry for distributed tracing, and gRPC for efficient, high-performance communication.

## Getting Started

These instructions will help you get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker
- Go (version 1.16 or later)

## Installation

1. Clone the repository:
   
   ```bash
   git clone https://github.com/koneal2013/SensorSphere.git
   ```
2. Navigate to the project directory:
   
   ```bash
   cd SensorSphere
   ```
3. Build and Start Docker containers:
   
   ```bash
   make docker-start
   ```
   
   - To stop the containers, run:
     ```bash
     make docker-down
     ```

Note: Make sure you have Docker and Docker Compose installed on your machine as the project runs in a Dockerized environment.

## Testing

To run the tests, use the following Makefile command:

```bash
make test


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

# Monitoring and Tracing with Prometheus, Jaeger, and Zipkin

SensorSphere integrates with Prometheus, Jaeger, and Zipkin to provide comprehensive monitoring and tracing capabilities.

## Prometheus

Prometheus is an open-source systems monitoring and alerting toolkit. SensorSphere exports metrics to Prometheus, which can then be visualized in Grafana. This allows you to monitor the performance and health of the SensorSphere service.

To access Prometheus:

1. Navigate to `http://localhost:9090`.
2. Use the expression browser for ad-hoc querying.

## Jaeger

Jaeger is an open-source, end-to-end distributed tracing system that helps developers monitor and troubleshoot complex, microservice-based architectures. It's used for monitoring and troubleshooting microservices-based distributed systems, including:

- Distributed context propagation
- Distributed transaction monitoring
- Root cause analysis
- Service dependency analysis
- Performance / latency optimization

To access Jaeger:

1. Navigate to `http://localhost:16686`.
2. Use the search functionality to find specific traces.

## Zipkin

Zipkin is a distributed tracing system that helps gather timing data needed to troubleshoot latency problems in service architectures. It manages both the collection and lookup of this data through a Collector and a Query service. SensorSphere uses Zipkin to trace requests as they pass through the various services, providing a detailed view of the request path and latency data.

To access Zipkin:

1. Navigate to `http://localhost:9411`.
2. Use the search functionality to find specific traces.

For more detailed information on these integrations, please refer to the respective official documentation:

- [Prometheus](https://prometheus.io/docs/introduction/overview/)
- [Jaeger](https://www.jaegertracing.io/docs/1.22/)
- [Zipkin](https://zipkin.io/pages/documentation.html)

# Future Improvements

SensorSphere provides a base solution for sensor data management, there are several areas where the project could be expanded or improved:

1. **Expanded API Functionality:** The current API provides a solid foundation, but there are many more operations that could be useful in a production environment. For example, more complex queries, bulk operations, or real-time updates via websockets or server-sent events.
2. **Improved Error Handling:** While the current system does include basic error handling, there is always room for improvement. More granular error messages, retries for certain types of failures, and better logging could all contribute to a more robust system.
3. **Performance Optimizations:** As with any system, performance can almost always be improved. This could include things like query optimization, caching, or parallel processing of requests.
4. **Security Enhancements:** While the system does include basic security measures, there are many additional security enhancements that could be added. This could include things like rate limiting, more advanced authentication mechanisms, or improved encryption.
5. **User Interface:** Currently, SensorSphere is a backend service with no user interface. A web or mobile UI could be developed to allow users to interact with the system more easily.
6. **Automated Testing:** While the project does include some tests, the coverage could be expanded.

## Contributing

Contributions are welcome. Please feel free to open an issue or submit a pull request.

## Questions

If you have any questions about the project, please feel free to [email me](mailto:koneal2013@gmail.com).

