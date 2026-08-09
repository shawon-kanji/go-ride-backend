package kyc

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

// StoragePort abstracts the object-storage operations KYC use cases need,
// so they depend on an interface rather than infrastructure/storage directly.
type StoragePort interface {
	PresignUploadURL(ctx context.Context, key string, contentType string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

// DriverKycUpdater abstracts the one driver-repository operation KYC use
// cases need, so this package doesn't import domain/driver's full
// Repository interface just to flip one column. Matches
// domaindriver.Repository.UpdateKycStatus's exact signature so
// *DriverRepositoryGorm satisfies it with no adapter needed.
type DriverKycUpdater interface {
	UpdateKycStatus(ctx context.Context, id uuid.UUID, kycStatus string) (*domaindriver.Driver, error)
}

// VehicleOwnerChecker abstracts verifying a vehicle_id belongs to the
// uploading driver, needed when confirming a vehicle_registration document.
// Narrow on purpose, so this package doesn't import domain/vehicle's full
// Repository interface just to check ownership.
type VehicleOwnerChecker interface {
	IsOwnedByDriver(ctx context.Context, vehicleID uuid.UUID, driverID uuid.UUID) (bool, error)
}
