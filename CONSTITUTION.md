# go-ride-backend Constitution

This document governs how code is written in this repository: the layering
rules, the shape every domain/application package must follow, error
handling, testing, and naming conventions. It is **descriptive first,
prescriptive second** — every Article below is extracted from patterns
already consistently applied across `domain/{driver,vehicle,user,kyc}` and
`application/{driver,vehicle,user,kyc}` as of ratification, not invented
in the abstract. Follow it for new code; if a genuinely new situation
doesn't fit an Article, amend this document rather than silently
deviating (see Article XIII).

## Article I — Layered Architecture

Four layers, one-directional dependency flow:

```
interfaces/http  →  application  →  domain  ←  infrastructure
```

- **`domain/<aggregate>/`** has zero imports from any other layer in this
  repo. It defines the aggregate's shape and contracts only.
- **`application/<aggregate>/`** depends on its own `domain/<aggregate>`
  package, and on *other* aggregates only through a narrow port it defines
  itself (Article IV) — never on `infrastructure/*` directly, and never on
  another aggregate's full `domain.Repository`.
- **`infrastructure/*`** implements `domain` repository interfaces (and
  application-layer ports) against real dependencies (GORM/Postgres, S3,
  bcrypt, JWT). It may import `domain` types to satisfy those interfaces.
- **`interfaces/http/*`** depends on `application` use cases and
  `pkg/apperror`. It never touches `domain` or `infrastructure` directly.

`internal/bootstrap/app.go` is the one place all four layers are allowed
to meet (Article IX).

## Article II — Domain Layer Shape

Each aggregate gets its own package, `domain/<aggregate>/`, with exactly:

- **`entity.go`** — the struct (plain fields, no DB/JSON tags) plus any
  `const` enum blocks and small pure helper functions the entity's own
  invariants need (e.g. `domain/kyc/entity.go`'s `IsIdentityDocumentType`).
- **`errors.go`** — sentinel errors only: `var ErrX = errors.New("...")`.
  Compared with `errors.Is`, never a custom error type at this layer.
- **`repository.go`** — the `Repository` interface plus any `CreateInput`/
  `UpdateInput` structs its methods need. No implementation here.

Nothing else lives in `domain/`. If an aggregate seems to need more than
these three files, that's a signal the aggregate boundary is wrong, not
that this Article needs an exception.

## Article III — Application Layer Shape

Each aggregate gets `application/<aggregate>/` with:

- **`ports.go`** — narrow interfaces this package needs from outside
  itself (Article IV).
- **`dto.go`** — request/response structs, `validate:"..."` tags on
  request fields.
- **`mapper.go`** — `ToXResponse(entity) XResponse` functions, nothing
  else.
- **`validation.go`** — a local `ValidationError{Message string}` type,
  `validateXRequest(req) error` functions, and a local `formatFieldError`
  helper. **This is duplicated verbatim across every application package
  on purpose** (Article VII) — do not extract it to a shared `pkg/`.
- **One file per use case**, named for the action (`request_upload_url.go`,
  `update_online_status.go`): a `XUseCase` struct holding its dependencies,
  a `NewXUseCase(...)` constructor, and an `Execute(ctx, ...)` method. One
  use case, one exported responsibility.
- **One `_test.go` file per use case file**, same base name.

## Article IV — The Narrow-Port Rule

When a use case needs something from another aggregate, define the
smallest possible interface for exactly that need in the *consuming*
package's `ports.go` — never import another aggregate's full
`domain.Repository`.

```go
// application/vehicle/ports.go
type DriverKycStatusReader interface {
    IsKycApproved(ctx context.Context, driverID uuid.UUID) (bool, error)
}
```

Prefer a boolean or otherwise minimal return type over leaking the
producing package's raw enum/status string — `IsKycApproved(...) bool`,
not `GetKycStatus(...) string`, so the consumer never needs the
producer's status constants at all. An infrastructure struct (e.g.
`DriverRepositoryGorm`) implementing several such narrow ports from
different application packages, alongside its own `domain.Repository`, is
expected and correct — that's what makes the rule work without an
adapter layer.

## Article V — Error Handling

- **Domain errors**: sentinel `errors.New(...)` values in `errors.go`,
  checked with `errors.Is`.
- **Application validation errors**: a per-package `ValidationError`
  struct (Article III), checked with `errors.As`.
- **HTTP translation happens in exactly one place**:
  `pkg/apperror/errors.go`'s `Map(err) HTTPError`. One `switch`, every
  package's `ValidationError` type checked first via `errors.As`, then
  every domain sentinel via `errors.Is` in declaration order, default
  falls through to `500 INTERNAL_ERROR`. Every new domain error needs one
  new `case` here. **Never** map an error to an HTTP status anywhere else
  — no ad-hoc `c.JSON(4xx, ...)` for domain failures in a handler.

## Article VI — Repository Implementation

- `infrastructure/repository/<aggregate>_repository_gorm.go` implements
  the `domain.Repository` interface against `schemamodels` types from
  `go-ride-db-schema` — never against the domain entity directly.
- Every repo file has a private `toDomainX(model *schemamodels.X) *domainx.X`
  mapper. Domain entities never leak a GORM tag or a `schemamodels` type.
- Partial updates use `.Model(&schemamodels.X{}).Where(...).Updates(map[string]interface{}{...})`
  with an explicit `updated_at`, never a full-struct `Save()`.
- Every mutating method checks `result.RowsAffected == 0` and returns the
  matching domain `NotFound` error — a zero-rows update is not a silent
  no-op.
- Multi-statement invariants (supersede-then-insert in
  `KycRepositoryGorm.Create`, deactivate-then-activate in
  `VehicleRepositoryGorm.Activate`) are wrapped in `db.Transaction(...)`.
  Sequential calls are never assumed atomic.

## Article VII — Package Self-Containment

Application packages do not share code with each other via a common
`pkg/` for small, package-local concerns (`ValidationError`,
`formatFieldError`, mapper shapes). Each package is a complete, standalone
unit under `application/`. This trades a small amount of duplication for
zero coupling between aggregates at the application layer — consistent
with Article IV's narrow-port rule existing specifically to avoid
aggregates depending on each other's internals. `pkg/apperror` and
`pkg/events` are the only cross-cutting `pkg/` packages, and both exist
because they *must* see every aggregate (central HTTP error mapping,
shared event contracts) — that is the bar for adding a new one.

## Article VIII — Schema Ownership

This repo does **not** own its database schema. It lives in the sibling
`go-ride-db-schema` repo (`migrations/` + `models/`), consumed as a
versioned Go module dependency (`github.com/shawon-kanji/go-ride-db-schema`,
currently pinned in `go.mod` — check there for the exact version). This
repo's own `migrations/` directory is legacy and unused — never add to it.

Schema change workflow: add a migration + model change in
`go-ride-db-schema`, verify `migrate up`/`down` against real Postgres, tag
a release, then `go get github.com/shawon-kanji/go-ride-db-schema@vX.Y.Z && go mod tidy`
here **as its own commit**, separate from the feature code that uses the
new schema.

## Article IX — Dependency Injection

All construction happens in `internal/bootstrap/app.go`'s
`Build(ctx, cfg) (*App, error)`, in this order: repositories →
use cases (receiving repositories/ports as constructor args) → handlers
(receiving use cases) → router (receiving handlers). No package reaches
out and constructs its own concrete dependency when it could receive an
interface/constructor argument instead. `cmd/api/main.go` calls `Build`
once and does nothing else domain-related.

## Article X — HTTP Layer

- Handlers (`interfaces/http/handlers/*.go`) are thin:
  `ShouldBindJSON` → call exactly one use case's `Execute` → on error,
  `h.handleError(c, err)` (which only calls `apperror.Map`) → on success,
  `c.JSON`. No business logic, no repository/DB access, no error-to-status
  mapping inline in a handler.
- Routes live one file per domain under `interfaces/http/routes/`
  (`auth_routes.go`, `driver_routes.go`, `vehicle_routes.go`,
  `kyc_routes.go`), each exporting a single unexported
  `registerXRoutes(group *gin.RouterGroup, handler *handlers.XHandler, ...)`
  function that mounts that domain's routes on the passed-in group. A
  route file returns a `*gin.RouterGroup` only when a later file needs to
  mount under it (`registerDriverRoutes` returns the driver-protected
  group so `registerVehicleRoutes`/`registerKycRoutes` can nest under it
  without redeclaring the `AuthRequiredRole` middleware) — otherwise no
  return value. `routes.go` is the sole composition root — it builds the
  `gin.Engine`, attaches middleware, defines `/healthz`, and calls each
  `registerXRoutes` in dependency order. It defines no other route
  directly.

## Article XI — Testing

- Hand-rolled fakes implementing the relevant domain/port interface,
  never a mocking framework/library.
- A fake is defined once, in the `_test.go` file of the use case that
  first needs it (e.g. `fakeDriverRepo` in `signup_test.go`), and reused
  by every other test file in that package.
- Test names: `Test{UseCaseName}Execute{Scenario}` — one behavior per
  test function. Prefer several small, clearly-named test functions over
  one table-driven megatest for use-case tests.
- Every use case gets at minimum: one success-path test, one test per
  distinct error branch its `Execute` can return.
- `go build ./...`, `go vet ./...`, and `gofmt -l .` (empty output) must
  all be clean before considering any change done. `go test ./...` must
  pass.

## Article XII — Naming & Style

- Import aliases are consistent and predictable:
  `domain<aggregate>` for domain packages (`domaindriver`, `domainvehicle`,
  `domainkyc`, `domainuser`), `app<aggregate>` for application packages
  (`appdriver`, `appvehicle`, `appkyc`, `appuser`) when imported from
  outside that package.
- `NewXUseCase(...)` constructor, `Execute(ctx, ...)` method — every use
  case, no exceptions, no alternate verbs (`Run`, `Handle`, `Call`).
- No comments explaining *what* code does — only *why*, and only when
  genuinely non-obvious (a hidden constraint, a workaround for a specific
  backend behavior, a subtle invariant). If deleting the comment wouldn't
  confuse the next reader, delete it.
- No premature abstraction. Three similar lines across two files beats a
  shared helper introduced for its own sake — see Article VII for why
  this extends even to small validation helpers.
- Don't add error handling, validation, or a config knob for a scenario
  that can't actually happen given how this codebase calls the code in
  question. Trust internal callers; validate only at real boundaries
  (HTTP request bodies, external API responses).

## Article XIII — Amendments

This document is expected to change as the codebase grows past what it
currently describes. When new code would have to violate an Article to
be correct, that is a signal to amend the Article (with a one-line reason
in the log below), not to write code that quietly disagrees with it.
Vague aspirational principles that don't match what the code actually
does are worse than no constitution at all — keep every Article
falsifiable against the current codebase.

### Amendment Log

| Date | Change | Reason |
|------|--------|--------|
| 2026-08-09 | Initial ratification | Extracted from `domain/{driver,vehicle,user,kyc}` and `application/{driver,vehicle,user,kyc}`, written during the KYC feature's implementation as the first cross-cutting doc covering conventions rather than one feature. |

---
*Version 1.0.0 — Ratified 2026-08-09 — Last amended 2026-08-09*
