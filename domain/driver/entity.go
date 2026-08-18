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

const (
	KycStatusNotStarted = "not_started"
	KycStatusInReview   = "in_review"
	KycStatusApproved   = "approved"
	KycStatusRejected   = "rejected"
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
	IsPaused        bool
	KycStatus       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
