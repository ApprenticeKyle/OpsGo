# OpsGo

Standalone DevOps service decoupled from FlowGo.

## Features
- Real-time deployment logs via Server-Sent Events (SSE).
- Multi-service deployment support.
- Decoupled process management: OpsGo can restart other services without being terminated.

## API Endpoints
- `GET /api/v1/devops/summary`: Overall status and history.
- `POST /api/v1/devops/config`: Configure a new repository.
- `POST /api/v1/devops/deploy`: Trigger a deployment.
- `GET /api/v1/devops/events`: SSE endpoint for real-time logs.

## Setup
1. Configure environment/database in `internal/infrastructure/config`.
2. Run `go mod tidy`.
3. Build: `go build -o opsgo-server ./cmd/server/main.go`.
4. Run: `./opsgo-server`.
