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
- `migrations`: legacy local SQL migrations (shared schema is now owned by `go-ride-db-schema`)

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

4. Run shared schema migrations:

```bash
make migrate-up
```

5. Run API:

```bash
go run ./cmd/api
```

Server starts on `APP_PORT` (default `8080`).

JWT configuration supports issuer/audience claim validation through `JWT_ISSUER` and `JWT_AUDIENCE` in `.env`.

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
- Schema and migrations are managed by sibling package `go-ride-db-schema`.

## Deployment

Cluster/cloud provisioning (Terraform for VPC/EKS/RDS/ECR/IAM, plus local
kind tooling) lives in the sibling repo **`go-ride-infra`**, not here —
start there for the full picture:
[`docs/architecture.md`](https://github.com/shawon-kanji/go-ride-infra/blob/main/docs/architecture.md),
[`docs/runbook-cluster.md`](https://github.com/shawon-kanji/go-ride-infra/blob/main/docs/runbook-cluster.md),
[`docs/runbook-local.md`](https://github.com/shawon-kanji/go-ride-infra/blob/main/docs/runbook-local.md).

This repo owns:

- **Dockerfile**: [`Dockerfile`](Dockerfile) at repo root (single-module build context).
- **Helm chart**: [`deploy/helm/`](deploy/helm/) — `Chart.yaml`, `values.yaml` (defaults), and `values-local.yaml` / `values-staging.yaml` / `values-production.yaml` per-environment overrides.
- **Health check**: `GET /healthz` (added to [`interfaces/http/routes/routes.go`](interfaces/http/routes/routes.go) for the chart's liveness/readiness probes).
- **Secrets in staging/production**: [`internal/config/config.go`](internal/config/config.go)'s `Load(ctx)` checks for `DB_CREDENTIALS_SECRET_NAME` / `JWT_SECRET_NAME` env vars — if set, it fetches the named AWS Secrets Manager entry via [`go-ride-utils/awssecrets`](https://github.com/shawon-kanji/go-ride-utils/blob/main/awssecrets/awssecrets.go) (IRSA-authenticated) instead of reading `DB_USER`/`DB_PASSWORD`/`JWT_SECRET` from the environment directly. Locally those env vars are unset, so `Load()` behaves exactly as documented above (plain `.env` values).
- `go-ride-utils` is currently pulled in via a local `replace` directive in [`go.mod`](go.mod) pending a tagged release — see that file's comment.

Deploys happen via this repo's own CI/CD (not built yet): build image → push
to the ECR repo `go-ride-infra`'s Terraform creates → `helm upgrade
--install` using `deploy/helm/` against the `staging`/`production`
namespace. For manual local end-to-end testing (alongside the other
go-ride services) use `go-ride-infra/local/deploy-local.sh`.
