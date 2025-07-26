# dorsium-rpc-gateway

This project implements a microservice gateway for the Dorsium blockchain using a layered Onion Architecture. It exposes HTTP endpoints via the Fiber framework and is ready to host a large number of routes.

## Development

```
go run ./cmd
```

Run unit tests with:

```
go test ./...
```

## Configuration

The server reads its configuration from environment variables:

- `ADDRESS` sets the HTTP bind address (default `:8080`).
- `NODE_RPC` defines the node RPC endpoint (default `http://localhost:26657`).
- `APP_VERSION` overrides the build version string (default `dev`).
- `APP_MODE` sets the running mode (default `production`).
- `ADMIN_TOKEN` configures the admin authentication token (default `changeme`).
- `DISABLE_METRICS` turns off Prometheus metrics when set to `true`.
- `MAX_RESPONSE_SIZE` limits the allowed response size in bytes (default `1048576`).

