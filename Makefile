APP_NAME=go-ride-backend

.PHONY: run test tidy up down logs migrate-up migrate-down migrate-version

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

up:
	cd ../go-ride-infra/local && docker compose up -d

down:
	cd ../go-ride-infra/local && docker compose down

logs:
	cd ../go-ride-infra/local && docker compose logs -f postgres aistor

migrate-up:
	cd ../go-ride-db-schema && go run ./cmd/migrate up

migrate-down:
	cd ../go-ride-db-schema && go run ./cmd/migrate down

migrate-version:
	cd ../go-ride-db-schema && go run ./cmd/migrate version
