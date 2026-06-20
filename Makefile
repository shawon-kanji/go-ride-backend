APP_NAME=go-ride-backend

.PHONY: run test tidy up down logs

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
