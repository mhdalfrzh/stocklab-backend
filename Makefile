.PHONY: run dev build test migrate-up migrate-down tidy

# ─── Config ───────────────────────────────────────────────────────────────────
BINARY     := ./tmp/main
ENTRYPOINT := ./cmd/api
DB_URL     ?= postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable

# ─── Development ──────────────────────────────────────────────────────────────

## run: compile and run the server directly
run:
	go run $(ENTRYPOINT)

## dev: start Air live-reload watcher
dev:
	air

## build: compile the production binary
build:
	go build -o $(BINARY) $(ENTRYPOINT)

# ─── Testing ──────────────────────────────────────────────────────────────────

## test: run all tests with race detector
test:
	go test -race -count=1 ./...

## test/cover: run tests and open coverage report
test/cover:
	go test -race -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out

# ─── Code Quality ─────────────────────────────────────────────────────────────

## tidy: tidy and verify go modules
tidy:
	go mod tidy
	go mod verify

## vet: run go vet
vet:
	go vet ./...

# ─── Database ─────────────────────────────────────────────────────────────────

## migrate-up: apply all pending migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

## migrate-down: roll back the last migration
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

## migrate-create name=<name>: create a new migration file pair
migrate-create:
	migrate create -ext sql -dir migrations $(name)

# ─── Help ─────────────────────────────────────────────────────────────────────

## help: print this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
