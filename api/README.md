# TimerWallet API

This is a Golang API server for a Telegram Mini App. It provides endpoints for Telegram authentication and data exchange.

## Features
- Basic REST API structure
- Telegram authentication handler (stub)
- Data exchange handler (stub)

## Getting Started

### Prerequisites
- Go 1.21 or newer

### Install dependencies
```
go mod tidy
```

### Run the server
```
go run ./cmd/main.go
```

The server will start on port 8080.

## Endpoints
- `POST /auth/telegram` — Telegram authentication
- `POST /data` — Data exchange

## Next Steps
- Implement Telegram authentication logic
- Implement data exchange logic
