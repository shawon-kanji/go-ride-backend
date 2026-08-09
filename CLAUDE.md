# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Coding conventions — read this first

**[`CONSTITUTION.md`](CONSTITUTION.md) is the canonical, binding source for
how code in this repo is written**: DDD layering, the exact file shape
every `domain/<aggregate>`/`application/<aggregate>` package must follow,
error handling, repository/transaction conventions, testing philosophy,
and naming. Read it before any non-trivial change — new code should match
its Articles, and if a genuinely new situation doesn't fit one, amend the
Constitution rather than quietly deviating from it (see its Article XIII).
This file (`CLAUDE.md`) covers orientation and commands; `CONSTITUTION.md`
covers *how to write the code itself*.

## What this is

`go-ride-backend` is one of several sibling repos making up the **go-ride**
ride-hailing platform: user/rider and driver auth, driver vehicle
management, and driver KYC document verification (Gin + GORM + Postgres).
Ride booking/dispatch/realtime live in the sibling `go-ride-kafka-consumers`
repo instead — this repo does not implement the ride lifecycle itself. See
`doc/PROJECT_ROADMAP.md` for the full phase-by-phase status across both.

Sibling repos, checked out alongside this one for local dev:
`go-ride-db-schema` (owns the DB schema — see Constitution Article VIII),
`go-ride-kafka-consumers` (ride/dispatch/realtime microservices),
`go-ride-infra` (Terraform + local dev tooling, including the consolidated
`docker-compose.yml` this repo's `make up` delegates to),
`go-ride-utils` (shared Go library), `go-ride-driver-app` / `go-ride-user-app`
(the Expo/React Native clients this API serves).

## Repo layout

- `cmd/api`: app entrypoint (`main.go` → `config.Load` → `bootstrap.Build`)
- `domain/{driver,vehicle,user,kyc}`: entities, repository contracts, domain errors — see Constitution Article II
- `application/{driver,vehicle,user,kyc}`: DTOs, mappers, validation, use cases — see Constitution Article III
- `infrastructure/{repository,security,storage,db}`: GORM repository implementations, bcrypt/JWT, S3-compatible object storage client
- `interfaces/http/{handlers,routes,middleware}`: Gin handlers and per-domain route files — see Constitution Article X
- `internal/bootstrap`: the single dependency-injection composition root (`app.go`) — see Constitution Article IX
- `internal/config`: env-driven config, with AWS Secrets Manager fallback in staging/production
- `pkg/apperror`: the one place domain/validation errors become HTTP responses — see Constitution Article V
- `migrations`: legacy, unused — the real schema is in `go-ride-db-schema` (Constitution Article VIII)
- `doc/`: `PROJECT_ROADMAP.md` (phase-by-phase status across the whole platform) and `DRIVER_KYC_PLAN.md` (KYC feature detail)

## Commands

```bash
go run ./cmd/api        # run the API (reads .env — cp .env.example .env first)
go build ./...
go vet ./...
gofmt -l .               # must be empty
go test ./...
make up / make down / make logs   # backing services (Postgres, AIStor) — delegates to ../go-ride-infra/local
make migrate-up / migrate-down / migrate-version   # delegates to ../go-ride-db-schema
```

## Schema changes

Never edit this repo's own `migrations/`. See Constitution Article VIII —
schema changes happen in `go-ride-db-schema`, tagged, then pinned here via
`go get github.com/shawon-kanji/go-ride-db-schema@vX.Y.Z && go mod tidy` as
its own commit.
