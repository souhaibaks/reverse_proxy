# reverse_proxy

A concurrent load-balancing reverse proxy built in Go, featuring Round-Robin and Least-Connections load balancing, periodic health checking, and a dynamic Admin API.

## Project Structure

```
reverse_proxy/
├── main.go           # Entry point — loads config, starts proxy + admin + health checker
├── Backend.go        # Backend struct with thread-safe alive status and health check
├── ServerPool.go     # Server pool with round-robin / least-conn selection
├── LoadBalancer.go   # LoadBalancer interface definition
├── ProxyHandler.go   # HTTP handler that forwards requests to backends
├── HealthChecker.go  # Periodic health check goroutine
├── AdminApi.go       # Admin API (port 8081) — add/remove/status backends
├── ProxyConfig.go    # Config struct and global config variable
├── config.json       # Configuration file
└── Backends/         # Sample backend servers for testing
```

## Configuration (`config.json`)

```json
{
  "port": 9000,
  "strategy": "round-robin",
  "health_check_frequency": 30000000000,
  "backends": [
    "http://localhost:9001",
    "http://localhost:9002",
    "http://localhost:9003"
  ]
}
```

| Field | Description |
|---|---|
| `port` | Port the proxy listens on |
| `strategy` | `"round-robin"` or `"least-conn"` |
| `health_check_frequency` | Health check interval in nanoseconds (30s = `30000000000`) |
| `backends` | Initial list of backend URLs |

## Running

### 1. Start backend servers (in separate terminals)

```bash
go run Backends/backend1.go   # listens on :9001
go run Backends/backend2.go   # listens on :9002
go run Backends/backend3.go   # listens on :9003
```

### 2. Start the proxy

```bash
go run *.go --config=config.json
```

The proxy listens on `:9000` and the Admin API on `:8081`.

## Admin API (port 8081)

### Get status summary

```bash
curl http://localhost:8081/status
```

Response:
```json
{
  "total_backends": 3,
  "active_backends": 3,
  "backends": [...]
}
```

### List all backends

```bash
curl http://localhost:8081/backends
```

### Add a backend

```bash
curl -X POST http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:9004"}'
```

### Remove a backend

```bash
curl -X DELETE http://localhost:8081/backends \
  -H "Content-Type: application/json" \
  -d '{"url": "http://localhost:9004"}'
```

## Testing the proxy

```bash
# Send a request through the proxy
curl http://localhost:9000/

# Send multiple requests to observe load balancing
for i in {1..6}; do curl -s http://localhost:9000/; done
```
