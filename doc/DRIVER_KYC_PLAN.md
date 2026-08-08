# Driver KYC & Document Verification — Implementation Plan

Status: Not started. Referenced from `PROJECT_ROADMAP.md` Phase 8 ("Security and Compliance Hardening").

Goal: A driver cannot go online (and a vehicle cannot be activated) until their identity and vehicle documents have been uploaded and approved by a human backoffice reviewer.

AI verification is explicitly **out of scope for this phase** — it is documented below only as a future enhancement and insertion point, not something to build now.

## Documents Collected

- Driver selfie / live photo
- Government-issued ID (front + back)
- Driving license (front + back)
- Vehicle registration certificate
- Deferred to a later pass: insurance certificate, vehicle fitness/inspection certificate

## Data Model (`go-ride-db-schema`)

- New `driver_documents` table: `id`, `driver_id` (FK), `document_type` (enum: `selfie`, `govt_id_front`, `govt_id_back`, `driving_license_front`, `driving_license_back`, `vehicle_registration`, ...), `file_url` (object storage key, not a DB blob), `status` (`uploaded` / `approved` / `rejected`), `rejection_reason`, `reviewed_by` (admin id), `reviewed_at`, `created_at`, `updated_at`.
  - Reserve nullable `ai_verdict` / `ai_score` / `ai_checked_at` columns now so the future AI pass (see below) doesn't require a breaking schema change, but leave them unused for this phase.
  - Append-only / `is_current` flag on re-upload after rejection, so review history is auditable instead of overwritten.
- New `drivers.kyc_status` column (`not_started` / `in_review` / `approved` / `rejected`), separate from the existing `account_status` — `account_status` governs auth lifecycle, `kyc_status` governs document verification. A vehicle-level KYC link (registration doc tied to `vehicle_id`) follows the same pattern.

## Storage

Object storage (S3-compatible), not DB blobs. Backend issues presigned PUT URLs for driver app uploads and presigned GET URLs for backoffice review, so files never transit through the API server.

## Flow

1. Driver signs up (existing flow) → `kyc_status = not_started`.
2. Driver app requests a presigned upload URL per document type, uploads directly to object storage, then confirms the upload to the backend, which records a `driver_documents` row with `status = uploaded`.
3. Once all required documents are uploaded → `kyc_status = in_review`.
4. Backoffice reviewer approves/rejects each document individually with a reason. All required docs approved → `kyc_status = approved`; any rejection → `kyc_status = rejected` and the driver is prompted to re-upload just the rejected document(s).
5. `is_online` toggle (`application/driver/update_online_status.go`) and vehicle activation are gated on `kyc_status = approved`.

## Backoffice/Admin Surface

List pending submissions, view documents, approve/reject with a reason. Every decision should be audit-logged (who, when, before/after status) — ties directly into Phase 8's existing "Audit logging for sensitive operations" and "Role-based access model for internal/admin operations" items, so build them together rather than bolting audit logging on afterward.

## Future Enhancement — AI Pre-Verification (explicitly out of scope now)

Insert an async AI check between `uploaded` and `in_review` — document authenticity/tamper detection, OCR field cross-check against the driver's profile data, and a face-match between the selfie and ID photo. This would only set `ai_verdict`/`ai_score` to prioritize or flag submissions for the backoffice queue; it must never auto-approve or auto-reject. Human backoffice review stays the sole final authority.

When this is picked up, it likely fits the existing Kafka-consumer pattern: a worker subscribing to a `document.uploaded` event, calling a vision/OCR provider, and writing results back before the item reaches backoffice review.

## Technical Tasks (when this is scheduled)

- [ ] Migration: `driver_documents` table + `drivers.kyc_status` column
- [ ] Object storage integration (S3/MinIO) + presigned upload/download URL endpoints
- [ ] Driver-facing endpoints: request upload URL, confirm upload, view own KYC status/documents
- [ ] Backoffice-facing endpoints: list submissions, view document, approve/reject with reason
- [ ] Gate `is_online` toggle and vehicle activation on `kyc_status = approved`
- [ ] Audit log entry on every review decision
- [ ] Driver notification on approval/rejection (depends on Phase 6 push/email gap being closed, or at minimum a websocket/in-app notice)
- [ ] **Future phase, not now:** AI pre-screening worker (authenticity/tamper/OCR/face-match) feeding `ai_verdict` for reviewer triage only

## Exit Criteria (when this is built)

- [ ] Driver cannot go online without `kyc_status = approved`
- [ ] Every document has a clear audit trail of who approved/rejected it and why
- [ ] Rejected documents can be re-uploaded without losing review history
