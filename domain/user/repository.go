package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines persistence contracts for the user aggregate.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	Deactivate(ctx context.Context, id uuid.UUID, at time.Time) error
}
