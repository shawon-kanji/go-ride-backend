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
	UpdateOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) (*Driver, error)
	UpdatePauseStatus(ctx context.Context, id uuid.UUID, isPaused bool) (*Driver, error)
	UpdateKycStatus(ctx context.Context, id uuid.UUID, kycStatus string) (*Driver, error)
}
