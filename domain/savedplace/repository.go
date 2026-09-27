package savedplace

import (
	"context"

	"github.com/google/uuid"
)

type CreateInput struct {
	UserID           uuid.UUID
	Label            string
	PlaceType        string
	FormattedAddress string
	Lat              float64
	Lng              float64
}

type UpdateInput struct {
	Label            string
	PlaceType        string
	FormattedAddress string
	Lat              float64
	Lng              float64
}

// Repository defines persistence contracts for the saved-place aggregate.
type Repository interface {
	Create(ctx context.Context, input *CreateInput) (*SavedPlace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SavedPlace, error)
	// GetByUserAndType finds a user's existing 'home' or 'work' place, used
	// to pre-check the one-per-type limit before Create/Update rather than
	// parsing a Postgres unique-violation error code.
	GetByUserAndType(ctx context.Context, userID uuid.UUID, placeType string) (*SavedPlace, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*SavedPlace, error)
	Update(ctx context.Context, id uuid.UUID, input *UpdateInput) (*SavedPlace, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
