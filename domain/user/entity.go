package user

import (
	"time"

	"github.com/google/uuid"
)

// User is the core user domain entity.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
