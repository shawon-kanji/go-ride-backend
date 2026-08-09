# Driver KYC & Document Verification — Implementation Plan

Status: **Driver-facing upload flow and gating implemented** (2026-08-09). Backoffice review (reviewer accounts, approve/reject endpoints, audit logging) is **not built** — document and `kyc_status` approval/rejection happens by editing the database directly for now. Referenced from `PROJECT_ROADMAP.md` Phase 8 ("Security and Compliance Hardening").

Goal: A driver cannot go online, and a vehicle cannot be activated, until the required identity and vehicle documents have been uploaded and approved.

AI verification remains explicitly **out of scope** — documented below only as a future enhancement and insertion point.

## Two independent verification tracks

KYC splits into two tracks, not one aggregate status, because "is this person who they say they are" and "is this specific vehicle verified" are different questions with different lifecycles:

- **Driver identity** — selfie, government ID (front/back), driving license (front/back). Aggregated on `drivers.kyc_status` (`not_started` / `in_review` / `approved` / `rejected`), independent of which vehicle the driver drives.
- **Vehicle verification** — five documents per vehicle (`domain/kyc.RequiredVehicleDocumentTypes`: registration, photo front/back/side, number plate), each linked via `driver_documents.vehicle_id`. There is **no** `vehicles.kyc_status` column; "is this vehicle verified" is derived live as "does it have a current, approved row for every required vehicle document type" (`domain/kyc.Repository.IsVehicleDocumentApproved`), so there's a single place to flip in the database rather than a column to keep in sync. A driver can have current documents for multiple vehicles at once.

## Documents Collected

- Driver selfie / live photo (identity)
- Government-issued ID, front + back (identity)
- Driving license, front + back (identity)
- Vehicle registration certificate (vehicle-scoped, requires `vehicle_id`)
- Vehicle photos: front, back, side (vehicle-scoped, requires `vehicle_id`)
- Vehicle number plate photo (vehicle-scoped, requires `vehicle_id`)
- Deferred to a later pass: insurance certificate, vehicle fitness/inspection certificate

## Data Model (`go-ride-db-schema`)

`driver_documents` table (migrations `000021`, `000023`, `000024`): `id`, `driver_id` (FK), `vehicle_id` (nullable FK, set for the five vehicle-scoped types), `document_type` (`selfie` / `govt_id_front` / `govt_id_back` / `driving_license_front` / `driving_license_back` / `vehicle_registration` / `vehicle_photo_front` / `vehicle_photo_back` / `vehicle_photo_side` / `vehicle_number_plate`), `file_url` (object storage key), `status` (`uploaded` / `approved` / `rejected`), `rejection_reason`, `is_current`, `created_at`/`updated_at`. `reviewed_by`/`reviewed_at` columns were deliberately **not added** — there's no reviewer identity yet (see Backoffice section below).

Reserved, currently unused columns for a future AI pass: `ai_verdict`, `ai_score`, `ai_checked_at`.

`is_current` + re-upload: on a new confirmed upload, any existing current row in the same scope is superseded (`is_current = false`) rather than overwritten, so rejection/re-upload history is preserved. "Same scope" means same `(driver_id, document_type)` for identity documents, same `(vehicle_id, document_type)` for vehicle documents — this is why the partial unique index is split in two (migration `000023`): one unique per `(driver_id, document_type) WHERE vehicle_id IS NULL`, one unique per `(vehicle_id, document_type) WHERE vehicle_id IS NOT NULL`.

`drivers.kyc_status` column (migration `000022`): `not_started` / `in_review` / `approved` / `rejected`, separate from the existing `account_status` — `account_status` governs auth lifecycle, `kyc_status` governs driver identity verification only.

## Storage

Object storage (S3-compatible via `infrastructure/storage/s3_storage.go`, targeting AIStor locally / real S3 in staging-production), not DB blobs. Backend issues presigned PUT URLs for driver app uploads; files never transit through the API server. Object keys: `drivers/{driverID}/{documentType}/{uuid}` for identity documents, `drivers/{driverID}/{documentType}/{vehicleID}/{uuid}` for any of the five vehicle-scoped types.

## Flow (implemented)

1. Driver signs up (existing flow) → `kyc_status = not_started`.
2. Driver requests a presigned upload URL: `POST /api/v1/driver/kyc/documents/upload-url` (`application/kyc.RequestUploadURLUseCase`). `vehicle_id` is required in the request when `document_type` is one of the five vehicle-scoped types, forbidden otherwise.
3. Driver uploads directly to object storage, then confirms: `POST /api/v1/driver/kyc/documents/confirm` (`application/kyc.ConfirmUploadUseCase`). This verifies the object actually exists in storage (extra `HEAD` call, guards against confirming a URL nothing was ever PUT to), verifies vehicle ownership for vehicle-scoped uploads (`ErrDocumentForbidden` if the vehicle isn't the driver's), supersedes the prior current row in scope, and records the new `driver_documents` row (`status = uploaded`).
4. Once all 5 identity document types have a current row → `drivers.kyc_status` auto-recomputes to `in_review` (`ConfirmUploadUseCase.recomputeKycStatus`). Vehicle-scoped uploads never touch `drivers.kyc_status`.
5. **Not automated yet**: moving a document to `approved`/`rejected`, and `drivers.kyc_status` to `approved`/`rejected`, is done directly against Postgres (`UPDATE driver_documents SET status = 'approved' WHERE ...`, `UPDATE drivers SET kyc_status = 'approved' WHERE ...`) — see Backoffice section.
6. Driver views own status: `GET /api/v1/driver/kyc/status` (`application/kyc.GetStatusUseCase`) — returns `kyc_status` plus every current document (identity and vehicle-scoped).
7. Gating (`application/driver/update_online_status.go`, `application/vehicle/activate.go`): going online requires driver `kyc_status = approved` **and** all 5 required vehicle document types approved for the currently active vehicle; activating/switching a vehicle requires driver `kyc_status = approved` **and** all 5 approved for *that specific* vehicle. Vehicle switching reuses the same activate use case, so one gate covers both "go online with vehicle X" and "switch to vehicle X while online."

## Backoffice/Admin Surface — not built

No admin/reviewer identity, login, or HTTP endpoints exist. There's no admin app anywhere in the workspace either. Approving/rejecting a document today means an operator runs SQL directly against `driver_documents.status` and `drivers.kyc_status`. When this is picked up:

- A minimal `backoffice_reviewers` table + login (issuing a JWT with `role = "backoffice"`, reusing the existing generic `middleware.AuthRequiredRole`) is the smallest addition that gives a verifiable "who" for `reviewed_by`.
- Review endpoints: list pending, view a document (presigned GET, `S3Storage.PresignDownloadURL` already exists and is unused for this), approve/reject with a reason.
- `reviewed_by`/`reviewed_at` columns need adding to `driver_documents`.
- Audit logging of review decisions ties into Phase 8's "audit logging for sensitive operations" item — build together, not bolted on later.
- **Re-upload resets review state deliberately**: re-uploading a document (identity or vehicle) after manual approval pulls the relevant status back (`kyc_status → in_review` for identity; the vehicle's derived approval flips false the moment its current row's `status` isn't `approved`) — the new document hasn't been reviewed yet, even if its siblings were. A future reviewer flow inherits this for free since it reads the same `is_current`/`status` state; no separate revocation logic needed.

## Future Enhancement — AI Pre-Verification (still out of scope)

Unchanged from the original plan: insert an async AI check between `uploaded` and `in_review` — document authenticity/tamper detection, OCR field cross-check, face-match between selfie and ID. Only sets `ai_verdict`/`ai_score` to prioritize the backoffice queue; never auto-approves/rejects. Likely fits the existing Kafka-consumer pattern (a worker subscribing to a `document.uploaded`-style event) once built.

## Key Files

- Schema: `go-ride-db-schema/migrations/000021_create_driver_documents.*`, `000022_add_driver_kyc_status.*`, `000023_split_driver_documents_current_index.*`, `000024_add_vehicle_photo_document_types.*`; `go-ride-db-schema/models/driver_document.go`, `models/driver.go`.
- Domain: `domain/kyc/` (entity, errors, repository contract), `domain/driver/entity.go` (`KycStatus`).
- Application: `application/kyc/` (upload-url, confirm-upload, get-status use cases, ports, validation), gating in `application/driver/update_online_status.go` and `application/vehicle/activate.go`.
- Infrastructure: `infrastructure/repository/kyc_repository_gorm.go`, `driver_repository_gorm.go` (`UpdateKycStatus`, `IsKycApproved`), `vehicle_repository_gorm.go` (`IsOwnedByDriver`).
- HTTP: `interfaces/http/handlers/kyc_handler.go`, routes under `/api/v1/driver/kyc/*` in `interfaces/http/routes/routes.go`.

## Exit Criteria

- [x] Driver cannot go online without driver `kyc_status = approved` **and** all 5 required document types approved for the active vehicle
- [x] Driver cannot activate/switch to a vehicle without driver `kyc_status = approved` **and** all 5 required document types approved for that vehicle
- [x] Rejected documents can be re-uploaded without losing review history (`is_current` supersede, scoped correctly per vehicle)
- [ ] Every document has a clear audit trail of who approved/rejected it and why — **not built**, no reviewer identity yet
