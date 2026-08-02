package driver

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccountStatusPending = "pending"
	AccountStatusActive  = "active"
	AccountStatusBlocked = "blocked"
)

// Driver is the core driver identity entity.
type Driver struct {
	ID              uuid.UUID
	Email           string
	PasswordHash    string
	FirstName       string
	LastName        string
	AccountStatus   string
	IsEmailVerified bool
	IsOnline        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
