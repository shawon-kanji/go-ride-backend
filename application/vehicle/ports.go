package vehicle

import (
	"context"

	"github.com/google/uuid"
)

// DriverKycStatusReader abstracts checking whether a driver is KYC-approved,
// needed to gate vehicle activation. Returns a bool rather than the raw
// kyc_status string so this package never needs to import domain/driver's
// status constants.
type DriverKycStatusReader interface {
	IsKycApproved(ctx context.Context, driverID uuid.UUID) (bool, error)
}

// VehicleDocumentApprovalReader abstracts checking whether a specific
// vehicle has an approved registration document, needed to gate activation
// (including switching the active vehicle, which reuses this same use
// case) on that vehicle's own verification, not just the driver's identity
// KYC. Narrow on purpose, so this package doesn't import domain/kyc's full
// Repository interface just to check one thing.
type VehicleDocumentApprovalReader interface {
	IsVehicleDocumentApproved(ctx context.Context, vehicleID uuid.UUID) (bool, error)
}
