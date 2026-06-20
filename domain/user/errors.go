package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already registered")
	ErrInvalidCredential = errors.New("invalid credentials")
)
