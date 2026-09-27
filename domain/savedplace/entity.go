package savedplace

import (
	"time"

	"github.com/google/uuid"
)

const (
	PlaceTypeHome  = "home"
	PlaceTypeWork  = "work"
	PlaceTypeOther = "other"
)

// SavedPlace is a rider's saved location (Home, Work, or a custom place),
// owned by a user.
type SavedPlace struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Label            string
	PlaceType        string
	FormattedAddress string
	Lat              float64
	Lng              float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
