package driver

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines persistence contracts for the driver aggregate.
type Repository interface {
	Create(ctx context.Context, driver *Driver) error
	GetByEmail(ctx context.Context, email string) (*Driver, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Driver, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*Driver, error)
}
