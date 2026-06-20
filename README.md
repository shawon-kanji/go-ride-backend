# Go Ride Backend

DDD-styled Golang backend boilerplate for a cab booking app with user authentication (signup/login), GORM, and PostgreSQL.

## Tech Stack

- Go
- Gin (HTTP)
- GORM (ORM)
- PostgreSQL
- JWT auth
- bcrypt password hashing

## Project Structure

- `cmd/api`: app entrypoint
- `domain/user`: domain entity, repository contract, domain errors
- `application/user`: DTOs, mappers, signup/login use cases
- `infrastructure`: db, repository implementations, security adapters
- `interfaces/http`: handlers, middleware, routes
- `migrations`: SQL migration files

## Setup

1. Copy env template:

```bash
cp .env.example .env
```

2. Start PostgreSQL:

```bash
docker compose up -d
```

3. Install deps:

```bash
go mod tidy
```

4. Run API:

```bash
go run ./cmd/api
```

Server starts on `APP_PORT` (default `8080`).

## API Endpoints

- `POST /api/v1/auth/signup`
- `POST /api/v1/auth/login`
- `GET /api/v1/me` (requires `Authorization: Bearer <token>`)

### Signup Body

```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "Jane",
  "last_name": "Doe"
}
```

### Login Body

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

## Notes

- `docker-compose.yml` provisions PostgreSQL with env-driven credentials.
- GORM `AutoMigrate` is enabled for bootstrap and SQL migrations are included for versioned migration flow.
