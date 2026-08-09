package kyc

import (
	"context"

	"github.com/google/uuid"
)

// CreateInput is the data needed to record a confirmed upload as a new
// current document row.
type CreateInput struct {
	DriverID     uuid.UUID
	VehicleID    *uuid.UUID
	DocumentType string
	FileURL      string
}

// Repository defines persistence contracts for the driver-document aggregate.
type Repository interface {
	// Create marks any existing current row scoped the same way (same
	// driverID+documentType for identity documents, same vehicleID+
	// documentType for vehicle documents) as superseded, then inserts the
	// new row as current, in one transaction — this is how re-upload after
	// rejection preserves history without blocking a driver from having
	// current registration documents for two different vehicles at once.
	Create(ctx context.Context, input *CreateInput) (*DriverDocument, error)
	GetByID(ctx context.Context, id uuid.UUID) (*DriverDocument, error)
	ListCurrentByDriver(ctx context.Context, driverID uuid.UUID) ([]*DriverDocument, error)
	// IsVehicleDocumentApproved reports whether vehicleID has a current,
	// approved vehicle_registration document — the sole source of truth for
	// "is this vehicle verified", so there's no separate status column on
	// vehicles to keep in sync.
	IsVehicleDocumentApproved(ctx context.Context, vehicleID uuid.UUID) (bool, error)
}
