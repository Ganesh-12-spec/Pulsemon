# Pulsemon

Pulsemon is a Go-based HTTP server health monitoring service.

It periodically checks multiple HTTP targets, tracks their health state, records response history, calculates uptime and average latency, exposes status and Prometheus-style metrics, and sends webhook alerts when targets fail repeatedly or recover.

## Features

* Periodic HTTP health checks
* Multiple monitored targets
* Concurrent-safe health state management
* Failure and recovery state tracking
* Alert after 3 consecutive failures
* Webhook notifications
* Health-check history
* Uptime percentage calculation
* Average response latency calculation
* JSON status API
* Prometheus-style metrics endpoint
* Structured logging with `log/slog`
* Graceful shutdown
* Docker support
* Unit tests
* Race detector testing

## Architecture

```text
                         ┌─────────────────┐
                         │     Pulsemon    │
                         └────────┬────────┘
                                  │
                           Periodic checks
                                  │
                    ┌─────────────┼─────────────┐
                    ↓             ↓             ↓
                Example        Google        Target N
                    │             │             │
                    └─────────────┼─────────────┘
                                  ↓
                           Health Checker
                                  │
                                  ↓
                           State Monitor
                         ┌────────┴────────┐
                         ↓                 ↓
                        UP                DOWN
                         │                 │
                         ↓                 ↓
                     History         Failure Counter
                         │                 │
                         ↓                 ↓
                      Metrics        3 consecutive
                         │              failures
                         │                 │
                         │                 ↓
                         │           Webhook Alert
                         │
                         ↓
                    HTTP API
                  /status /metrics
```

## Project Structure

```text
pulsemon/
├── cmd/
│   └── pulsemon/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── status.go
│   │   └── status_test.go
│   │
│   ├── checker/
│   │   └── checker.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── history/
│   │   ├── history.go
│   │   └── history_test.go
│   │
│   ├── metrics/
│   │   └── metrics.go
│   │
│   ├── state/
│   │   ├── state.go
│   │   └── state_test.go
│   │
│   └── webhook/
│       ├── webhook.go
│       └── webhook_test.go
│
├── Dockerfile
├── go.mod
└── README.md
```

## Monitoring Flow

```text
time.Ticker
     │
     ↓
HTTP Health Check
     │
     ↓
State Update
     │
     ├──────────────→ Failure / Recovery Detection
     │                         │
     │                         ↓
     │                    Webhook Alert
     │
     ↓
History Update
     │
     ↓
Uptime Calculation
     │
     ↓
Average Latency
     │
     ↓
Metrics Update
     │
     ↓
/status + /metrics
```

## Health Monitoring

Pulsemon periodically checks every configured target.

For every check it records:

* Target name
* UP/DOWN status
* Response latency
* Check time

A successful HTTP check produces an `UP` state.

A failed check produces a `DOWN` state.

## Failure Detection

Pulsemon tracks consecutive failures for every target.

A failure webhook is sent when the same target reaches:

```text
Failure
   ↓
Failure
   ↓
Failure
   ↓
ALERT
```

The threshold is **3 consecutive failures**.

If a successful check occurs, the consecutive failure counter is reset.

## Recovery Detection

When a target changes from:

```text
DOWN → UP
```

Pulsemon detects the recovery and sends a recovery webhook.

## Webhook

Pulsemon sends HTTP POST requests containing JSON data.

Example payload:

```json
{
  "target": "Google",
  "status": "DOWN",
  "event": "failure_threshold",
  "consecutive_failures": 3
}
```

Recovery example:

```json
{
  "target": "Google",
  "status": "UP",
  "event": "recovered"
}
```

## HTTP API

Pulsemon exposes its API on port `8080`.

### GET /status

Returns the current status of monitored targets.

```bash
curl http://localhost:8080/status
```

Example:

```json
{
  "targets": {
    "Example": "UP",
    "Google": "UP"
  }
}
```

### GET /metrics

Returns Prometheus-style metrics.

```bash
curl http://localhost:8080/metrics
```

Example:

```text
pulsemon_target_up{target="Example"} 1
pulsemon_target_up{target="Google"} 1
pulsemon_uptime_percent{target="Example"} 100.00
pulsemon_uptime_percent{target="Google"} 100.00
pulsemon_average_latency_ms{target="Example"} 120.50
pulsemon_average_latency_ms{target="Google"} 95.20
```

## Graceful Shutdown

Pulsemon handles operating-system shutdown signals such as:

* `Ctrl+C`
* `SIGTERM`

During shutdown:

```text
Shutdown Signal
      ↓
Cancel Monitoring Context
      ↓
Stop Monitoring Loop
      ↓
Gracefully Shutdown HTTP Server
      ↓
Finish Active Requests
      ↓
Pulsemon Stops Cleanly
```

The HTTP server receives a 5-second shutdown timeout for active requests.

## Running Locally

Clone the repository and enter the project:

```bash
cd ~/Projects/pulsemon
```

Run:

```bash
go run ./cmd/pulsemon
```

The server starts on:

```text
http://localhost:8080
```

Check the API:

```bash
curl http://localhost:8080/status
```

```bash
curl http://localhost:8080/metrics
```

## Testing

Run the normal test suite:

```bash
go test ./...
```

Run tests with the Go race detector:

```bash
go test -race ./...
```

Build the project:

```bash
go build ./...
```

The project uses tests for important components including:

* State transitions
* Failure threshold
* Recovery detection
* Failure counter reset
* Uptime calculation
* Average latency
* Status API
* Webhook requests

## Docker

Build the Docker image:

```bash
docker build -t pulsemon .
```

Run Pulsemon:

```bash
docker run --rm -p 8080:8080 pulsemon
```

Verify the API:

```bash
curl http://localhost:8080/status
```

```bash
curl http://localhost:8080/metrics
```

Stop the container:

```text
Ctrl+C
```

Pulsemon handles the shutdown signal and gracefully stops the HTTP server.

## Technologies

* Go
* `net/http`
* `time.Ticker`
* `context`
* `os/signal`
* `log/slog`
* HTTP
* JSON
* Prometheus-style metrics
* Docker
* Go testing
* Go race detector

## What This Project Demonstrates

Pulsemon was built to practice practical Go backend engineering concepts:

* HTTP clients and servers
* Goroutines
* Context cancellation
* Time-based processing
* Concurrent-safe state
* Mutexes
* State machines
* Error handling
* JSON serialization
* Webhooks
* Structured logging
* Metrics and observability
* HTTP testing
* Race detection
* Graceful shutdown
* Containerization

## Project Status

Pulsemon is a learning project focused on practical Go backend engineering, monitoring, observability, concurrency, testing, HTTP services, and containerization.
