package driver

import "context"

// Repository defines persistence contracts for the driver aggregate.
type Repository interface {
	Create(ctx context.Context, driver *Driver) error
	GetByEmail(ctx context.Context, email string) (*Driver, error)
}
