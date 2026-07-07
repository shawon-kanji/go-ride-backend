APP_NAME=go-ride-backend

.PHONY: run test tidy up down logs migrate-up migrate-down migrate-version

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f postgres

migrate-up:
	cd ../go-ride-db-schema && go run ./cmd/migrate up

migrate-down:
	cd ../go-ride-db-schema && go run ./cmd/migrate down

migrate-version:
	cd ../go-ride-db-schema && go run ./cmd/migrate version
