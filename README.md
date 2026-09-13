# StockLab Backend

> REST API backend for StockLab — portfolio management with real-time Yahoo Finance quotes.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21+ |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) + pgx driver |
| Database | PostgreSQL 16 |
| Decimals | [shopspring/decimal](https://github.com/shopspring/decimal) |
| Live Reload | [Air](https://github.com/cosmtrek/air) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |

---

## Project Structure

```
stocklab-backend/
├── cmd/
│   └── api/
│       └── main.go              # Entrypoint — wires dependencies, starts server
├── internal/
│   ├── config/
│   │   └── config.go            # Loads config from environment variables
│   ├── domain/
│   │   └── portfolio.go         # GORM entity models (Portfolio, Stock, Holding)
│   ├── portfolio/
│   │   ├── dto.go               # API response DTOs
│   │   ├── errors.go            # Sentinel errors (ErrPortfolioNotFound)
│   │   ├── handler.go           # HTTP handlers (Gin)
│   │   ├── repository.go        # Database access layer (GORM)
│   │   └── service.go           # Business logic
│   └── yahoofinance/
│       └── client.go            # Yahoo Finance HTTP client
├── migrations/
│   └── *.sql                    # SQL migration files
├── .air.toml                    # Air live-reload config
├── .env.example                 # Environment variable template
├── docker-compose.yml           # PostgreSQL via Docker
├── go.mod / go.sum              # Go module files
├── Makefile                     # Developer shortcuts
└── README.md
```

---

## Prerequisites

- [Go](https://golang.org/doc/install) 1.21+
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) — for database migrations
- [Air](https://github.com/cosmtrek/air) *(optional)* — for live reload during development

---

## Getting Started

### 1. Clone & configure environment

```bash
cp .env.example .env
# Edit .env with your database credentials if needed
```

### 2. Start the database

```bash
docker compose up -d
```

### 3. Run database migrations

```bash
make migrate-up
```

### 4. Run the server

```bash
# One-shot run
make run

# Live-reload (requires Air)
make dev
```

The server starts on `http://localhost:8080` by default. Override with the `PORT` env var.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable` | Full PostgreSQL connection string |
| `PORT` | `8080` | HTTP port the server listens on |

---

## Makefile Targets

```
make run           # Run server directly with go run
make dev           # Start Air live-reload watcher
make build         # Compile binary to ./tmp/main
make test          # Run all tests with race detector
make test/cover    # Run tests and open HTML coverage report
make tidy          # go mod tidy + verify
make vet           # go vet ./...
make migrate-up    # Apply all pending migrations
make migrate-down  # Roll back last migration
make migrate-create name=<migration_name>
make help          # Print all available targets
```

---

## API Reference

### `GET /api/portfolios/dashboard`

Returns the portfolio dashboard for the authenticated user.

**Headers**

| Header | Required | Description |
|---|---|---|
| `X-User-ID` | ✅ | Integer user ID |

**Example Response `200 OK`**

```json
{
  "summary": {
    "totalPortfolioValue": "15432.50",
    "cashBalance": "2000.00",
    "totalGainLoss": "1432.50",
    "totalGainLossPercentage": "10.23",
    "dayChange": "85.00",
    "dayChangePercentage": "0.63"
  },
  "holdings": [
    {
      "ticker": "BBCA",
      "companyName": "Bank Central Asia Tbk",
      "quantity": 100,
      "averageBuyPrice": "9000.00",
      "currentPrice": "9432.50",
      "marketValue": "943250.00",
      "gainLoss": "43250.00",
      "gainLossPercentage": "4.81"
    }
  ]
}
```

**Error Responses**

| Status | Condition |
|---|---|
| `401 Unauthorized` | Missing or invalid `X-User-ID` header |
| `404 Not Found` | No portfolio found for the given user |
| `500 Internal Server Error` | Unexpected server error |

---

## Database Migrations

```bash
# Create a new migration
make migrate-create name=add_column_to_users

# Apply all pending migrations
make migrate-up

# Roll back the last migration
make migrate-down
```

---

## Running Tests

```bash
make test

# With coverage
make test/cover
```
