package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SignupInput struct {
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	AccountStatus string
}

// Repository defines persistence contracts for the user aggregate.
type Repository interface {
	Create(ctx context.Context, input *SignupInput) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	Deactivate(ctx context.Context, id uuid.UUID, at time.Time) error
}
