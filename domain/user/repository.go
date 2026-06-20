package user

import "context"

// Repository defines persistence contracts for the user aggregate.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}
