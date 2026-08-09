package driver

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

// PasswordHasher abstracts password hashing and verification.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

// TokenManager abstracts token generation for auth flows.
type TokenManager interface {
	Generate(userID string, email string, role string) (string, error)
}

// VehicleLister abstracts fetching a driver's vehicles, needed to gate going online
// on the driver having at least one active vehicle.
type VehicleLister interface {
	ListByDriver(ctx context.Context, driverID uuid.UUID) ([]*domainvehicle.Vehicle, error)
}

// VehicleDocumentApprovalReader abstracts checking whether a specific
// vehicle has an approved registration document, needed to gate going
// online on the driver's currently active vehicle being verified — not
// just the driver's own identity KYC. Narrow on purpose, so this package
// doesn't import domain/kyc's full Repository interface just to check one
// thing.
type VehicleDocumentApprovalReader interface {
	IsVehicleDocumentApproved(ctx context.Context, vehicleID uuid.UUID) (bool, error)
}
