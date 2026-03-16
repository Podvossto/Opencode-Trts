.PHONY: up down build logs migrate tidy dev-main dev-portal dev-api

## Start all services
up:
	docker compose up -d

## Stop all services
down:
	docker compose down

## Build all Docker images
build:
	docker compose build

## Follow logs for all services
logs:
	docker compose logs -f

## Run latest migration against local postgres
migrate:
	@PGPASSWORD=$(POSTGRES_PASSWORD) psql -h localhost -U $(POSTGRES_USER) -d $(POSTGRES_DB) \
		-f $(shell ls migrations/*.sql | sort | tail -1)

## Run go mod tidy inside go-api container
tidy:
	cd services/go-api && go mod tidy

## Dev mode — apps/main (port 3000)
dev-main:
	cd apps/main && npm run dev

## Dev mode — apps/portal (port 3001)
dev-portal:
	cd apps/portal && npm run dev

## Dev mode — Go API
dev-api:
	cd services/go-api && go run ./cmd/api

## Dev mode — AI service
dev-ai:
	cd ai-service && uvicorn main:app --reload --port 8000

## Type-check all TypeScript projects
type-check:
	cd packages/shared && npm run type-check
	cd apps/main && npm run type-check
	cd apps/portal && npm run type-check
