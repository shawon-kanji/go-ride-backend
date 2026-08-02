# Go Ride Backend Project Roadmap

Last Updated: 2026-06-29 (Phase 2 implementation progress: driver auth + lifecycle APIs)
Owner: Backend Team
Primary Tracking File: doc/PROJECT_ROADMAP.md

## How To Use This Roadmap

1. Keep phase status updated in the Phase Tracker table.
2. For each completed task, check the box in the corresponding phase section.
3. Add brief implementation notes in Progress Log after each meaningful change.
4. Do not remove old log entries; append new entries for auditability.

## Status Legend

- Not Started
- In Progress
- Blocked
- Completed

## Phase Tracker

| Phase | Name | Status | Target Outcome |
| --- | --- | --- | --- |
| 0 | Foundation and Architecture | Completed | DDD boilerplate, app bootstrapping, config, database wiring |
| 1 | Auth Module Stabilization | Completed | Production-ready signup/login with robust validation and error handling |
| 2 | Identity Expansion (Rider + Driver Auth) | In Progress | Minimal driver signup/login and role-aware auth without full driver management |
| 3 | Ride Booking Core | Not Started | End-to-end ride request lifecycle with assignment-ready contract |
| 4 | Driver and Fleet Domain | Not Started | Driver onboarding, availability, and assignment readiness |
| 5 | Pricing and Payments | Not Started | Fare calculation, payment intents, and transaction records |
| 6 | Notifications and Real-Time Updates | Not Started | Event-driven user and driver updates |
| 7 | Observability and Reliability | Not Started | Monitoring, tracing, resilience, and operability |
| 8 | Security and Compliance Hardening | Not Started | Hardened auth, data protection, and policy compliance |
| 9 | Testing, QA, and Release | Not Started | High-confidence release candidate with CI gates |
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

### Features

- [ ] Create ride request endpoint
- [ ] Ride status lifecycle: requested, accepted, started, completed, canceled
- [ ] Assignment-ready contract with nullable driver_id for early booking flow
- [ ] Ride cancellation with rules by actor and state
- [ ] Fare estimate endpoint before booking
- [ ] Ride history endpoint for users

### Technical Tasks

- [ ] Add ride aggregate with invariants in domain layer
- [ ] Add ride repository interfaces and PostgreSQL implementations
- [ ] Add migrations for rides, pickup/drop coordinates, status timeline
- [ ] Add initial assignment port interface for driver matching integration
- [ ] Add idempotency key handling for create ride endpoint

### Exit Criteria

- [ ] Ride state machine blocks invalid transitions
- [ ] Ride APIs meet latency and correctness targets in local tests

---

## Phase 4: Driver and Fleet Domain

Goal: Enable driver onboarding and assignment-ready availability tracking.

### Features

- [ ] Driver profile and verification status
- [ ] Driver vehicle registration and active vehicle selection
- [ ] Driver online/offline availability toggle
- [ ] Nearest-driver query for dispatch candidates

### Technical Tasks

- [ ] Introduce driver and vehicle bounded contexts
- [ ] Add geospatial indexing approach for nearby search
- [ ] Add migrations for drivers, vehicles, and availability records

### Exit Criteria

- [ ] Eligible drivers are discoverable for dispatch
- [ ] Driver availability changes are reflected in near real-time

---

## Phase 5: Pricing and Payments

Goal: Add transparent fare computation and payment processing support.

### Features

- [ ] Fare policy service with base, distance, and time components
- [ ] Surge pricing policy with configurable rules
- [ ] Payment intent creation and confirmation flow
- [ ] Transaction history for riders and drivers
- [ ] Failed payment recovery workflow

### Technical Tasks

- [ ] Add pricing module with versioned fare rules
- [ ] Add payment adapter interface and implementation stubs
- [ ] Add migrations for fares, invoices, and payment transactions

### Exit Criteria

- [ ] Fare is reproducible and auditable per ride
- [ ] Payment workflow handles retries and failure states safely

---

## Phase 6: Notifications and Real-Time Updates

Goal: Keep rider and driver apps synchronized with ride events.

### Features

- [ ] Push and email notification hooks
- [ ] Real-time ride status updates channel
- [ ] Driver assignment and ETA notifications
- [ ] Retry strategy for failed outbound notifications

### Technical Tasks

- [ ] Define event contracts for ride lifecycle events
- [ ] Add queue or pub/sub abstraction
- [ ] Implement notification workers with dead-letter handling

### Exit Criteria

- [ ] Critical ride events are delivered reliably
- [ ] Failed deliveries are visible and retryable

---

## Phase 7: Observability and Reliability

Goal: Make operations measurable and production-ready.

### Features

- [ ] Structured logs with correlation IDs
- [ ] Metrics for API latency, error rates, and DB performance
- [ ] Distributed tracing for request flows
- [ ] Health and readiness endpoints
- [ ] Graceful shutdown and startup checks

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
- [ ] Implement full vehicle and driver KYC with driver lincense upload , self image , govt id , car registration documents , etc

### Exit Criteria

- [ ] Security checklist completed with documented decisions
- [ ] No critical vulnerabilities in release candidate

---

## Phase 9: Testing, QA, and Release

Goal: Produce a stable release process with confidence gates.

### Features

- [ ] Unit tests for domain and use cases
- [ ] Integration tests against PostgreSQL test environment
- [ ] API contract tests and regression suite
- [ ] CI workflow for lint, test, and build
- [ ] Versioned release notes and changelog process

### Technical Tasks

- [ ] Add test fixtures and test data builders
- [ ] Add CI pipeline with required checks
- [ ] Add pre-release checklist automation

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

---

## Progress Log

| Date | Phase | Update | By |
| --- | --- | --- | --- |
| 2026-06-20 | 0 | Initial DDD scaffold, auth module skeleton, GORM PostgreSQL wiring, docker compose, and base docs created. | Copilot |
| 2026-06-20 | 2/3 | Roadmap revised: minimal driver signup/login moved earlier; ride booking core kept ahead of full driver/fleet management using assignment-ready contract. | Copilot |
| 2026-06-20 | 1/10 | Deferred refresh tokens, auth rate limiting, and brute-force protection to post-completion STP hardening phase. | Copilot |
| 2026-06-20 | 1 | Completed Phase 1 with strict request validation, unified auth error model, standardized JWT claims (issuer/audience/expiry), and auth test coverage for happy/failure/expired/malformed token paths. | Copilot |
| 2026-06-29 | 2 | Implemented driver signup/login, role-aware JWT + middleware role checks, protected driver endpoint, user profile lifecycle endpoints (`me`, profile update, change-password, deactivate), repository/model/migration updates, and end-to-end API+DB verification. | Copilot |

## Risks and Dependencies

- PostgreSQL availability and local environment consistency can block integration tests.
- Payment gateway selection impacts Phase 5 design and timeline.
- Real-time architecture choice (polling vs event streaming) impacts Phase 6 scope.
- Security hardening decisions can introduce breaking API behavior if delayed too long.

## Change Control

When scope changes:

1. Update relevant phase features and tasks.
2. Add a Progress Log entry describing the decision.
3. Mark affected phases as In Progress or Blocked accordingly.
