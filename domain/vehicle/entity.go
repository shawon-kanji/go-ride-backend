package vehicle

import (
	"time"

	"github.com/google/uuid"
)

const (
	CategoryNormal = "normal"
	CategoryLuxury = "luxury"
)

// Vehicle is the core vehicle domain entity, owned by a driver.
type Vehicle struct {
	ID          uuid.UUID
	DriverID    uuid.UUID
	PlateNumber string
	Color       string
	ModelName   string
	SeatCount   int
	Category    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
