# CHIS - Webhook Delivery Service

A distributed webhook delivery system built with Go. It handles sending webhooks to customer endpoints with automatic retries, tracking, and monitoring.

## Architecture

```
┌─────────────────┐
│   API Service   │  REST API Gateway (port 8080)
└────────┬────────┘
         │ gRPC
         ▼
┌─────────────────┐
│Delivery Service │  Saves messages & queues for delivery
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌───────┐ ┌─────────┐
│Worker │ │Scheduler│
└───────┘ └─────────┘
    │         │
    └────┬────┘
         ▼
   PostgreSQL + Redis
```

**Services:**
- **API** - REST endpoints for webhook submission and user management
- **Delivery** - gRPC service that persists messages and queues them
- **Worker** - Pulls from queue and delivers webhooks via HTTP
- **Scheduler** - Retries failed messages with exponential backoff

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make

### Run Locally

1. Clone and setup environment:
```bash
cp .env.example .env
```

2. Start infrastructure:
```bash
docker-compose up -d
```

3. Run database migrations:
```bash
make migrate-up
```

4. Start services (in separate terminals):
```bash
make run-api
make run-grpc
make run-worker
make run-scheduler
```

## API Usage

### Authentication

Register a user:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123", "firstName": "John", "lastName": "Doe"}'
```

Login:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret123"}'
```

### API Keys

Create an API key (requires JWT token):
```bash
curl -X POST http://localhost:8080/api/api-key/create \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json"
```

### Send Webhook

```bash
curl -X POST http://localhost:8080/webhook/send \
  -H "X-API-Key: <api_key>" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/webhook",
    "payload": {"event": "user.created", "data": {"id": 123}}
  }'
```

## Configuration

Environment variables (see `.env.example`):

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `REDIS_URL` | Redis connection string |
| `JWT_SECRET` | Secret for JWT signing |
| `GRPC_ADDRESS` | Delivery service gRPC address |

## Retry Strategy

Failed webhooks are retried with exponential backoff:

| Attempt | Delay |
|---------|-------|
| 1 | 1 second |
| 2 | 2 seconds |
| 3 | 4 seconds |
| 4 | 8 seconds |
| 5 | 16 seconds |

After 5 failed attempts, the message is marked as dead-letter.

## Monitoring

- **Health checks**: Each service exposes `/health`
- **Metrics**: Prometheus metrics at port 9090 (`/metrics`)

Available metrics:
- `http_requests_total` - HTTP request count
- `http_request_duration_ms` - HTTP latency
- `webhooks_delivered_total` - Delivery count by status
- `webhook_delivery_duration_ms` - Delivery duration

## Kubernetes Deployment

Deploy to Kubernetes:
```bash
kubectl apply -f k8s/
```

This creates:
- Namespace `chis`
- Deployments for all services
- PostgreSQL and Redis
- ConfigMaps and Secrets
- LoadBalancer for API service

## Project Structure

```
├── cmd/                  # Service entry points
│   ├── api/             # REST API service
│   ├── delivery/        # gRPC delivery service
│   ├── worker/          # Webhook worker
│   └── scheduler/       # Retry scheduler
├── internal/            # Core application code
│   ├── handler/         # HTTP handlers
│   ├── repository/      # Data access
│   ├── queue/           # Redis queue
│   └── model/           # Data models
├── proto/               # gRPC definitions
├── k8s/                 # Kubernetes manifests
└── pkg/                 # Shared utilities
```

## Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Chi
- **Database**: PostgreSQL 16
- **Queue**: Redis 7.2
- **RPC**: gRPC
- **Monitoring**: Prometheus

## License

MIT
