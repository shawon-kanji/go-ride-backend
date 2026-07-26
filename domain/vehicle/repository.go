package vehicle

import (
	"context"

	"github.com/google/uuid"
)

type CreateInput struct {
	DriverID    uuid.UUID
	PlateNumber string
	Color       string
	ModelName   string
	SeatCount   int
	Category    string
}

type UpdateInput struct {
	PlateNumber string
	Color       string
	ModelName   string
	SeatCount   int
	Category    string
}

// Repository defines persistence contracts for the vehicle aggregate.
type Repository interface {
	Create(ctx context.Context, input *CreateInput) (*Vehicle, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error)
	GetByPlateNumber(ctx context.Context, plate string) (*Vehicle, error)
	ListByDriver(ctx context.Context, driverID uuid.UUID) ([]*Vehicle, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateInput) (*Vehicle, error)
	Activate(ctx context.Context, driverID uuid.UUID, vehicleID uuid.UUID) (*Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
