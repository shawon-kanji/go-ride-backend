# Go Ride Backend Project Roadmap

Last Updated: 2026-08-02 (Full audit against actual code across go-ride-backend, go-ride-kafka-consumers, and go-ride-db-schema)
Owner: Backend Team
Primary Tracking File: doc/PROJECT_ROADMAP.md

## How To Use This Roadmap

1. Keep phase status updated in the Phase Tracker table.
2. For each completed task, check the box in the corresponding phase section.
3. Add brief implementation notes in Progress Log after each meaningful change.
4. Do not remove old log entries; append new entries for auditability.

## Architecture & Scope Note (added 2026-08-02)

This roadmap was originally written as if `go-ride-backend` (the DDD monolith) would host every domain: auth, rides, drivers, pricing, notifications. In practice, ride booking, dispatch, fare, and real-time notification logic were built as **separate services in `go-ride-kafka-consumers`**, backed by shared schema/models in **`go-ride-db-schema`**. That work is real and substantial (see Phases 3, 5, 6 below), but it is **not integrated into `go-ride-backend`** — there is no unified API/BFF in this repo that exposes ride booking to rider/driver apps through the DDD structure described here. Phase statuses below now reflect system-wide reality, with notes on which repo owns each piece.

## Status Legend

- Not Started
- In Progress
- Built Elsewhere (Not Integrated) — implemented in a sibling repo, not wired into `go-ride-backend`
- Blocked
- Completed

## Phase Tracker

| Phase | Name | Status | Target Outcome |
| --- | --- | --- | --- |
| 0 | Foundation and Architecture | Completed | DDD boilerplate, app bootstrapping, config, database wiring |
| 1 | Auth Module Stabilization | Completed | Production-ready signup/login with robust validation and error handling |
| 2 | Identity Expansion (Rider + Driver Auth) | In Progress | Minimal driver signup/login and role-aware auth without full driver management |
| 3 | Ride Booking Core | Built Elsewhere (Not Integrated) | Full ride request→assignment→completion lifecycle live in `go-ride-kafka-consumers`; no equivalent exists in `go-ride-backend` |
| 4 | Driver and Fleet Domain | In Progress | Vehicle registration + online toggle live in `go-ride-backend`; nearest-driver dispatch lives in `go-ride-kafka-consumers`; KYC/verification still missing everywhere |
| 5 | Pricing and Payments | Built Elsewhere (Partial) | Fare estimate live in `go-ride-kafka-consumers`; surge hardcoded to 1.0; payments are cash-collection tracking only, no gateway/transaction history |
| 6 | Notifications and Real-Time Updates | Built Elsewhere (Partial) | Websocket real-time updates fully live in `go-ride-kafka-consumers`; push/email notifications not implemented anywhere |
| 7 | Observability and Reliability | Not Started | Only liveness health endpoints and graceful shutdown exist; no structured logs, metrics, or tracing |
| 8 | Security and Compliance Hardening | Not Started | No secrets management, audit logging, CORS/security headers, or KYC document workflow anywhere |
| 9 | Testing, QA, and Release | In Progress | Unit tests exist in `go-ride-backend`; `go-ride-kafka-consumers` has zero tests and no CI; `go-ride-backend` also has no CI |
| 10 | Post-Completion STP Hardening | Not Started | Deferred auth hardening controls after core app completion |

---

## Phase 0: Foundation and Architecture

Goal: Establish a scalable DDD backend skeleton with local runtime support.

### Features

- [x] DDD structure with domain, application, infrastructure, and interfaces layers
- [x] Environment-based config loading
- [x] GORM database connection using PostgreSQL
- [x] Docker Compose for local PostgreSQL with healthcheck
- [x] Base HTTP server and route registration
- [x] Initial SQL migration for users table

### Exit Criteria

- [x] App builds and tests pass
- [x] Database container starts successfully
- [x] API boots with valid env configuration

---

## Phase 1: Auth Module Stabilization

Goal: Make authentication reliable, secure, and extensible.

### Features

- [x] Strict request validation for signup and login inputs
- [x] Consistent API error model for all auth endpoints
- [x] JWT claims standardization and expiry policy tuning
- [x] Access token guard for protected routes
- [x] Duplicate email handling with deterministic response codes

### Technical Tasks

- [x] Add validation package wiring and per-field error responses
- [x] Add centralized HTTP error mapper for all handlers
- [x] Add tests for signup/login happy path and failure path
- [x] Add security-focused tests for expired and malformed tokens

### Exit Criteria

- [x] Auth endpoints satisfy functional and security test suite
- [x] No plaintext secret leakage in logs or responses

---

## Phase 2: Identity Expansion (Rider + Driver Auth)

Goal: Add minimal driver authentication and role-aware identity while keeping full driver domain for later.

### Features

- [x] Driver signup endpoint
- [x] Driver login endpoint
- [x] Role-aware JWT claims (rider and driver)
- [x] Driver account status flag (pending, active, blocked)
- [x] Basic protected driver endpoint for auth flow validation
- [x] Get current user profile endpoint
- [x] Update profile endpoint with optimistic-safe updates
- [x] Change password endpoint with old-password verification
- [x] Soft delete or deactivate account flow
- [ ] Optional email verification workflow (deferred; `is_email_verified` defaults to false and is manually managed for now)
- [ ] Driver-side change-password endpoint (rider has `/change-password`; drivers have no equivalent yet — parity gap)

### Technical Tasks

- [x] Introduce identity role model and role checks in middleware
- [x] Extend auth DTOs and use cases for driver credentials
- [x] Add migrations for drivers table with auth-focused minimal fields
- [x] Add tests for driver signup/login and role-based route access
- [x] Add user service use cases and DTOs
- [x] Add repository methods for profile updates and status flags
- [x] Add migration for account status and audit columns
- [x] Add authorization checks for self-access operations

### Exit Criteria

- [x] Rider and driver can authenticate independently
- [x] Role-based route access is enforced by middleware
- [x] Profile lifecycle operations fully covered by tests
- [x] Account status state transitions are validated

---

## Phase 3: Ride Booking Core

Goal: Implement the rider journey from ride request to completion.

**Note:** All of this lives in `go-ride-kafka-consumers` (`cab-request-handler`, `driver-request-handler`, `trip-dispatch-worker`), backed by `go-ride-db-schema` migrations/models. None of it exists in `go-ride-backend`'s domain/application/interfaces layers — there is no `ride`/`trip` bounded context here, and no route in this repo exposes ride booking.

### Features

- [x] Create ride request endpoint — `cab-request-handler` `handleCreateCabRequest`
- [x] Ride status lifecycle: `search_started → offered → driver_accepted → assigned → in_progress → completed/cancelled` (superset of the originally planned states) — `schemamodels.TripRequest` / `OngoingTrip`
- [x] Assignment-ready contract with nullable driver_id — `trip_requests → driver_job_offers → ongoing_trips` pattern
- [x] Ride cancellation with rules by actor and state — `driver-request-handler` `CancelTrip`, `cab-request-handler` `cancel.go`
- [x] Fare estimate endpoint before booking — `handleFareEstimate`
- [x] Ride history endpoint for users — `trip_history` table/model

### Technical Tasks

- [ ] Add ride aggregate with invariants in domain layer — not applicable to `go-ride-backend`; logic exists as service code in `go-ride-kafka-consumers`, not a DDD aggregate
- [ ] Add ride repository interfaces and PostgreSQL implementations (in `go-ride-backend`) — none exist here
- [x] Add migrations for rides, pickup/drop coordinates, status timeline — done in `go-ride-db-schema` (`trip_requests`, `ongoing_trips`, etc.)
- [ ] Add initial assignment port interface for driver matching integration — services call each other directly via Kafka events, no port/adapter abstraction
- [x] Add idempotency key handling for create ride endpoint — `idempotency_key` header + DB uniqueness check in `createOrLoadCabRequest`

### Exit Criteria

- [ ] Ride state machine blocks invalid transitions — enforced ad hoc in handler code, not verified by tests (see gap: zero test files in `go-ride-kafka-consumers`)
- [ ] Ride APIs meet latency and correctness targets in local tests — no automated tests exist to confirm this
- [ ] **New:** Ride booking exposed through `go-ride-backend` (or an explicit decision that it stays a separate service) — undecided

---

## Phase 4: Driver and Fleet Domain

Goal: Enable driver onboarding and assignment-ready availability tracking.

**Note:** Split across repos — driver profile, vehicles, and online toggle live in `go-ride-backend`; nearest-driver dispatch lives in `go-ride-kafka-consumers`.

### Features

- [ ] Driver profile and verification status — profile exists (`domain/driver/entity.go`); "verification" is only `account_status` (pending/active/blocked), no real KYC/verification workflow (see Phase 8 gap)
- [x] Driver vehicle registration and active vehicle selection — `go-ride-backend` `application/vehicle/*` (register/list/get/update/activate/delete), fully tested
- [x] Driver online/offline availability toggle — `go-ride-backend` `PATCH /driver/online`, `application/driver/update_online_status.go`, tested
- [x] Nearest-driver query for dispatch candidates — `go-ride-kafka-consumers` `trip-dispatch-worker`, Haversine + S2-cell covering query over `driver_locations`

### Technical Tasks

- [x] Introduce driver and vehicle bounded contexts — done in `go-ride-backend`
- [x] Add geospatial indexing approach for nearby search — S2-cell indexing in `go-ride-db-schema`/`go-ride-kafka-consumers` (`driver_locations`, migration converting `s2_cell_id` to numeric for range queries)
- [x] Add migrations for drivers, vehicles, and availability records — drivers/vehicles migrations in `go-ride-backend`; `driver_locations` in `go-ride-db-schema`

### Exit Criteria

- [x] Eligible drivers are discoverable for dispatch — confirmed via `trip-dispatch-worker` nearest-driver query
- [x] Driver availability changes are reflected in near real-time — `driver-location-worker` + online/offline toggle
- [ ] **New:** Driver verification/KYC status is a real workflow, not just a manually-set enum — not started (tracked in Phase 8)

---

## Phase 5: Pricing and Payments

Goal: Add transparent fare computation and payment processing support.

**Note:** Fare estimation lives in `go-ride-kafka-consumers`. "Payments" so far means tracking cash collected by the driver at trip end, not a real payment gateway.

### Features

- [x] Fare policy service with base, distance, and time components — `buildFareEstimate` in `cab-request-handler`, backed by `fare_configs`/`fare_surcharges` tables
- [ ] Surge pricing policy with configurable rules — surge multiplier field exists but is hardcoded to `1.0`; no rules engine
- [ ] Payment intent creation and confirmation flow — not built; current flow is `handleCollectPayment` (cash-collection tracking via `ongoing_trips.payment_status`/`payment_collected_at`), not a payment-intent/gateway integration
- [ ] Transaction history for riders and drivers — no transaction table exists
- [ ] Failed payment recovery workflow — not started

### Technical Tasks

- [x] Add pricing module with versioned fare rules — `fare_configs`/`fare_surcharges` in `go-ride-db-schema`
- [ ] Add payment adapter interface and implementation stubs — no gateway adapter of any kind
- [ ] Add migrations for fares, invoices, and payment transactions — fare tables exist; no invoice/transaction tables

### Exit Criteria

- [ ] Fare is reproducible and auditable per ride — fare calc exists but not covered by tests
- [ ] Payment workflow handles retries and failure states safely — no gateway integration yet, so not applicable

---

## Phase 6: Notifications and Real-Time Updates

Goal: Keep rider and driver apps synchronized with ride events.

**Note:** Fully live in `go-ride-kafka-consumers` (`websocket-gateway`), Kafka-driven. Websocket-only — no push or email channel exists.

### Features

- [ ] Push and email notification hooks — not implemented; only websocket
- [x] Real-time ride status updates channel — `websocket-gateway` with per-event notifiers (`tripstart`, `tripend`, `tripcomplete`, `tripcancel`, `offers`, `tracking`), Redis-backed presence
- [x] Driver assignment and ETA notifications — `assignments`/`tracking` notifiers in `websocket-gateway`
- [ ] Retry strategy for failed outbound notifications — dead-letter/retry handling not confirmed

### Technical Tasks

- [x] Define event contracts for ride lifecycle events — `go-ride-utils/events` package
- [x] Add queue or pub/sub abstraction — Kafka, with a consumer per ride event
- [ ] Implement notification workers with dead-letter handling — consumers exist; explicit DLQ/retry handling not confirmed in code

### Exit Criteria

- [x] Critical ride events are delivered reliably — for websocket-connected clients, via Kafka consumers per event type
- [ ] Failed deliveries are visible and retryable — no DLQ/observability into delivery failures confirmed

---

## Phase 7: Observability and Reliability

Goal: Make operations measurable and production-ready.

### Features

- [ ] Structured logs with correlation IDs — a `correlation_id` field is propagated in events/DB rows for lineage, but actual logging is plain `log.Printf`/`gin.Logger()`, not structured
- [ ] Metrics for API latency, error rates, and DB performance — not started
- [ ] Distributed tracing for request flows — not started
- [x] Health endpoints — `/healthz` liveness check exists in `go-ride-backend` and every `go-ride-kafka-consumers` service (no separate readiness probe)
- [x] Graceful shutdown and startup checks — `signal.NotifyContext` + `http.Server.Shutdown` in `go-ride-backend/cmd/api` and all kafka-consumers services

### Technical Tasks

- [ ] Integrate OpenTelemetry-compatible instrumentation
- [ ] Add Prometheus metrics endpoint
- [ ] Add timeout and retry policies for downstream calls

### Exit Criteria

- [ ] SLO dashboards and basic alerts configured
- [ ] Runbook-ready operational visibility in place

---

## Phase 8: Security and Compliance Hardening

Goal: Reduce security risk and satisfy baseline compliance requirements.

### Features

- [ ] Secret management strategy beyond local env files
- [ ] Audit logging for sensitive operations
- [ ] Data retention and deletion policy controls
- [ ] CORS, headers, and transport hardening
- [ ] Role-based access model for internal/admin operations

### Technical Tasks

- [ ] Security headers middleware
- [ ] Threat model and abuse-case checklist
- [ ] Dependency vulnerability scanning in CI
- [ ] Implement full vehicle and driver KYC with driver lincense upload , self image , govt id , car registration documents , etc — confirmed zero hits across `go-ride-backend`, `go-ride-kafka-consumers`, and `go-ride-db-schema`; this is completely unbuilt

### Exit Criteria

- [ ] Security checklist completed with documented decisions
- [ ] No critical vulnerabilities in release candidate

### Driver KYC & Document Verification

- [ ] Document upload (selfie, govt ID, driving license, vehicle registration) with human backoffice approval gating driver online status and vehicle activation
- [ ] AI pre-screening of document authenticity — future enhancement only, not a final verifier; out of scope for now

Full implementation plan: [`doc/DRIVER_KYC_PLAN.md`](DRIVER_KYC_PLAN.md)

---

## Phase 9: Testing, QA, and Release

Goal: Produce a stable release process with confidence gates.

**Note:** Coverage is uneven across repos — `go-ride-backend` has unit tests but no CI; `go-ride-kafka-consumers` has no tests at all; `go-ride-db-schema` has both tests and CI (`.github/workflows/ci.yml`, `release.yml`).

### Features

- [x] Unit tests for domain and use cases — `go-ride-backend`: `application/user`, `application/driver`, `application/vehicle`, middleware all have `_test.go` coverage
- [ ] Integration tests against PostgreSQL test environment — not confirmed
- [ ] API contract tests and regression suite — not started
- [ ] CI workflow for lint, test, and build — **missing in `go-ride-backend` and `go-ride-kafka-consumers`**; exists only in `go-ride-db-schema`
- [ ] Versioned release notes and changelog process — only `go-ride-db-schema` has a `release.yml`

### Technical Tasks

- [ ] Add test fixtures and test data builders
- [ ] Add CI pipeline with required checks
- [ ] Add pre-release checklist automation
- [ ] **New:** Add any automated test coverage to `go-ride-kafka-consumers` — currently zero `_test.go` files across all services, including the ride lifecycle and dispatch logic from Phase 3

### Exit Criteria

- [ ] CI gates enforce quality baselines
- [ ] Release artifact and documentation are reproducible

---

## Phase 10: Post-Completion STP Hardening

Goal: Complete deferred auth security tasks after overall app delivery milestones are done.

### Features

- [ ] Optional refresh token model and rotation strategy
- [ ] Rate limiting for auth endpoints
- [ ] Brute-force protection policy

### Technical Tasks

- [ ] Add refresh-token storage strategy and token invalidation policy
- [ ] Add auth rate-limiter middleware and endpoint-specific thresholds
- [ ] Add login attempt tracking and temporary lockout rules
- [ ] Add tests for token rotation, throttling, and lockout behavior

### Exit Criteria

- [ ] Deferred auth hardening controls are implemented and tested
- [ ] Security review confirms acceptable auth abuse resistance

### Refresh Token Design Notes

Reference notes for implementing the "Optional refresh token model and rotation strategy" feature above.

**What it is:** A refresh token is a long-lived credential used to obtain new short-lived access tokens without forcing the user to log in again.

**Flow:**

1. **Login** issues two tokens: a short-lived access token (minutes to ~1 hour, sent on every API request, verified via JWT signature) and a long-lived refresh token (days to weeks, used only to mint new access tokens).
2. **Normal usage**: client sends the access token on API calls; server validates it statelessly.
3. **Access token expiry**: client calls a `/refresh` endpoint with the refresh token; server validates it against a store (DB lookup, since refresh tokens must be revocable) and issues a new access token.
4. **Logout/revocation**: refresh tokens are tracked server-side, so a session can be killed immediately, unlike a stateless access token which just has to expire on its own.

**Why split into two tokens:** short-lived access tokens limit blast radius if leaked; refresh tokens are transmitted less often and can be stored more securely (httpOnly cookie, secure storage), while still giving real revocation/session control.

**Implementation practices to apply here:**

- Store refresh tokens hashed in the DB (like passwords), not in plaintext.
- Use rotation: issue a new refresh token on every use and invalidate the old one. If an already-used (rotated-out) refresh token is replayed, treat it as theft and revoke the whole token family.
- Bind refresh tokens to device/session so users can view and revoke individual logged-in sessions.

---

## Progress Log

| Date | Phase | Update | By |
| --- | --- | --- | --- |
| 2026-06-20 | 0 | Initial DDD scaffold, auth module skeleton, GORM PostgreSQL wiring, docker compose, and base docs created. | Copilot |
| 2026-06-20 | 2/3 | Roadmap revised: minimal driver signup/login moved earlier; ride booking core kept ahead of full driver/fleet management using assignment-ready contract. | Copilot |
| 2026-06-20 | 1/10 | Deferred refresh tokens, auth rate limiting, and brute-force protection to post-completion STP hardening phase. | Copilot |
| 2026-06-20 | 1 | Completed Phase 1 with strict request validation, unified auth error model, standardized JWT claims (issuer/audience/expiry), and auth test coverage for happy/failure/expired/malformed token paths. | Copilot |
| 2026-06-29 | 2 | Implemented driver signup/login, role-aware JWT + middleware role checks, protected driver endpoint, user profile lifecycle endpoints (`me`, profile update, change-password, deactivate), repository/model/migration updates, and end-to-end API+DB verification. | Copilot |
| 2026-08-02 | 10 | Added refresh token design notes (flow, rationale, rotation/storage practices) to guide future implementation of the deferred refresh-token feature. | Claude |
| 2026-08-02 | 2-9 | Full audit of actual code vs. roadmap claims across `go-ride-backend`, `go-ride-kafka-consumers`, and `go-ride-db-schema`. Discovered ride booking, dispatch, fare estimation, and real-time notifications (Phases 3, 5, 6) are already substantially built, but as standalone services in sibling repos, not integrated into `go-ride-backend`. Updated phase statuses, checklists, and added an Identified Gaps section accordingly. | Claude |
| 2026-08-08 | 8 | Added Driver KYC & Document Verification as a bulleted feature under Phase 8, with the detailed implementation plan (document types, schema, storage, review flow, and the out-of-scope AI pre-screening insertion point) moved to `doc/DRIVER_KYC_PLAN.md` to keep the roadmap itself at feature-tracking altitude. | Claude |

## Identified Gaps

Gaps found during the 2026-08-02 audit that aren't fully captured by the phase checklists above, roughly in order of importance:

1. **No integration between `go-ride-backend` and the ride/dispatch services.** Ride booking, dispatch, and notifications work end-to-end as standalone services in `go-ride-kafka-consumers`, but `go-ride-backend` (the repo this roadmap describes) has no ride/trip domain, no routes, and no awareness of them. Either this roadmap should be rescoped to `go-ride-backend` only (with ride/dispatch tracked in a `go-ride-kafka-consumers` roadmap instead), or a decision is needed on building a unified API/BFF layer here.
2. **`go-ride-kafka-consumers` has zero automated tests and no CI**, despite containing the ride request, cancellation, dispatch, and fare-estimate logic — the highest-risk business logic in the system today has no regression safety net.
3. **`go-ride-backend` itself has no CI** (no `.github/workflows`), unlike `go-ride-db-schema` which already has one.
4. **Payments are not real payments** — the only implemented flow is marking cash as collected at trip end (`payment_status`/`payment_collected_at`). No payment gateway integration, no transaction/invoice history, no failed-payment recovery.
5. **Surge pricing is a hardcoded `1.0` multiplier**, not a configurable rules engine, despite Phase 5 listing it as a feature.
6. **No KYC / driver document verification workflow anywhere** (license, selfie, government ID, vehicle registration upload) — zero hits across all three repos. Driver "verification" today is just a manually-set `account_status` enum. Full plan: [`doc/DRIVER_KYC_PLAN.md`](DRIVER_KYC_PLAN.md).
7. **No push or email notifications** — only websocket delivery exists, so users not actively connected miss ride events.
8. **No structured logging, metrics, or distributed tracing**, despite a `correlation_id` field already being threaded through events/DB rows for this exact purpose — the plumbing for correlation exists but nothing consumes it yet.
9. **No email verification workflow** — `is_email_verified` column exists but nothing sets it to `true`.
10. **Driver-side change-password endpoint is missing** — rider has `/change-password`, driver does not (parity gap in Phase 2).
11. **No refresh tokens, rate limiting, or brute-force protection on auth** — already tracked in Phase 10, worth prioritizing sooner given auth is otherwise production-shaped.
12. **No secrets management, audit logging, or CORS/security-header hardening** anywhere (Phase 8, untouched).

## Risks and Dependencies

- PostgreSQL availability and local environment consistency can block integration tests.
- Payment gateway selection impacts Phase 5 design and timeline.
- Real-time architecture choice (polling vs event streaming) impacts Phase 6 scope — resolved in practice via Kafka + websocket, but only in `go-ride-kafka-consumers`, not reflected in this repo's architecture.
- Security hardening decisions can introduce breaking API behavior if delayed too long.
- **Architecture fragmentation risk:** as long as ride/dispatch logic lives outside `go-ride-backend` with no shared contract or integration tests between repos, changes to `go-ride-db-schema` models can silently break `go-ride-kafka-consumers` services without this roadmap or `go-ride-backend`'s CI (once it exists) catching it.

## Change Control

When scope changes:

1. Update relevant phase features and tasks.
2. Add a Progress Log entry describing the decision.
3. Mark affected phases as In Progress or Blocked accordingly.
