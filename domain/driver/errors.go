package driver

import "errors"

var (
	ErrDriverNotFound    = errors.New("driver not found")
	ErrEmailAlreadyTaken = errors.New("email already registered")
	ErrInvalidCredential = errors.New("invalid credentials")
)
