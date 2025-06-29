# shipment-coordinator

## What is it?

**shipment-coordinator** is a service for managing orders, carriers, and freight contracts. It allows you to create orders, register carriers with freight policies, get shipping quotes, and contract carriers for deliveries, applying region and status rules.

## Dependencies

- Go 1.24.4+
- PostgreSQL (Docker)
- [Go modules](https://golang.org/doc/go1.11#modules)
- [Docker Compose](https://docs.docker.com/compose/) (for local environment)

## How to Run

### Start the database

With Docker Compose:

```sh
make up
```

### Run the application

```sh
make build
./bin/shipment-coordinator
```
Or directly:
```sh
go run ./cmd/main.go
```

The API will be available on the configured port (default: 8080).

#### Links 

| Name                   | URL                                                            | Description                        |
|------------------------|----------------------------------------------------------------| ---------------------------------- |
| **API**                | [http://localhost:8080](http://localhost:8080)                 | Main API server                    |
| **API Documentation**  | [http://localhost:8080/docs](http://localhost:8080/docs)       | API documentation using Swagger    |
| **Prometheus UI**      | [http://localhost:9090](http://localhost:9090)                 | Prometheus UI for querying metrics |
| **Prometheus Metrics** | [http://localhost:8889/metrics](http://localhost:8889/metrics) | OTEL metrics endpoint (exposed)    |
| **Grafana**            | [http://localhost:3000/dashboards](http://localhost:3000)         | Visualization dashboard (Grafana)  |

## How to Test

Run all automated tests:

```sh
make test
```
or
```sh
go test ./... -v
```

## Project Structure

- `cmd/` — server startup and HTTP handlers
- `internal/order/` — order domain
- `internal/carrier/` — carrier domain
- `internal/shipping/` — contracts and quotes domain
- `internal/platform/` — integrations (DB, logger, config)
- `pkg/states/` — state and region utilities
