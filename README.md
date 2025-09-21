# beto-shared

**`beto-shared`** is a shared Go library used across all microservices of the **Beto Cloud Platform**.  
It provides reusable building blocks (errors, middlewares, observability, OpenStack clients, etc.) to ensure consistency, speed up development, and enforce best practices.

---

##  Key Features

- **Centralized Error Handling ([pkg/errors](./pkg/errors))**
  - Standardized error codes (global, DB, OpenStack: Keystone, Nova, Neutron…).
  - Helper functions (`NewDBNotFound`, `NewKeystoneAuthFailed`, `NewNeutronIPConflict`…).
  - Unified conversion to HTTP and gRPC errors (`ToHTTPError`, `ToGRPCError`).

- **Messaging ([pkg/eventbus](./pkg/eventbus))**
  - RabbitMQ helpers: `Publish`, `Consume` with retries/ack handling.
  - Simplifies **event-driven architecture** between microservices.

- **Database Helpers ([pkg/gormhelpers](./pkg/gormhelpers))**
  - Generic GORM utilities: CRUD, pagination, transactions.
  - Functions like `Exists`, `FindByConditions`, `FirstOrCreate`… ready to use.

- **Middlewares ([pkg/middleware](./pkg/middleware))**
  - HTTP: logging, recovery, request ID, Keystone JWT auth, Redis-based rate limiting.
  - gRPC: interceptors for errors, rate limiting, tracing, and metrics.
  - Unified error responses → always JSON or gRPC status in a consistent format.

- **OpenStack Clients ([pkg/openstack-client](./pkg/openstack-client))**
  - Keystone: login, token validation, role/project fetch.
  - Nova: VM management (create/list/delete).
  - Neutron: networks, floating IPs, firewalls.
  - Swift/Cinder: object and block storage support.

- **Redis Utilities ([pkg/redis](./pkg/redis))**
  - Centralized Redis client initialization.
  - Simple key-value cache (`Set`, `Get`, TTL).
  - Distributed rate limiter (token bucket).

- **Observability**
  - **Metrics ([pkg/metrics](./pkg/metrics))**
    - Exposes `/metrics` endpoints for Prometheus (HTTP/gRPC).
    - Collects latency, status codes, request counts.
  - **Tracing ([pkg/tracing](./pkg/tracing))**
    - OpenTelemetry integration → export to Jaeger/Tempo.
    - HTTP middleware + gRPC interceptor → distributed spans across services.
  - **Logging ([pkg/middleware/logging.go](./pkg/middleware/logging.go))**
    - Request correlation with `X-Request-ID`.

- **AuthN/Z ([pkg/authz](./pkg/authz))**
  - Keystone token validation (`/v3/auth/tokens`).
  - Injects claims (user, project, roles) into `context.Context`.
  - HTTP middleware + gRPC interceptor.
  - Helpers: `RequireRole(ctx, "admin")`, `RequireProject(ctx, "tenant-id")`.

---

## Repository Structure
```
pkg/
├── errors/ # Centralized error handling (codes + helpers)
├── eventbus/ # RabbitMQ (publish/consume helpers)
├── gormhelpers/ # GORM helpers (CRUD, transactions)
├── middleware/ # HTTP/gRPC middlewares
├── openstack-client/ # OpenStack clients (Keystone, Nova, Neutron, Swift/Cinder)
├── redis/ # Redis utils (cache + rate limiter)
├── metrics/ # Prometheus metrics (HTTP/gRPC)
├── tracing/ # OpenTelemetry (Jaeger/Tempo)
└── authz/ # AuthN/Z via Keystone
```

---

## Usage Examples

### Errors
```go
if err != nil {
    return nil, errors.NewKeystoneAuthFailed(err)
}
```
HTTP Middleware
```go
    mux := http.NewServeMux()
    mux.Handle("/secure", middleware.AuthMiddleware(validator)(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, _ := authz.GetClaims(r.Context())
            w.Write([]byte("Hello " + claims.UserName))
        }),
    ))
```
gRPC with Observability
```go
    grpcServer := grpc.NewServer(
    middleware.GRPCErrorInterceptor(),
    middleware.GRPCRateLimitInterceptor(func(ctx context.Context, req interface{}) string {
        return "ratelimit:grpc:" // dynamic key
    }, 100, time.Minute),
    tracing.GRPCTracingOptions()..., // OpenTelemetry tracing
)
```

 ## Roadmap
   - Add full OpenStack Cinder client .
   - Add Loki for centralized logging.
   - Support multiple backends (Postgres, MySQL, MinIO, Kafka).
   - Auto-generate documentation for each package.