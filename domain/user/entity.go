package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	AccountStatusActive      = "active"
	AccountStatusDeactivated = "deactivated"
)

// User is the core user domain entity.
type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	AccountStatus string
	DeactivatedAt *time.Time
	// RatingAverage is nil until the rider's first rating from a driver —
	// written only by driver-request-handler (go-ride-kafka-consumers);
	// this repo never sets it, only reads it.
	RatingAverage *float64
	RatingCount   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
