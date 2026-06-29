package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyTaken  = errors.New("email already registered")
	ErrInvalidCredential  = errors.New("invalid credentials")
	ErrAccountDeactivated = errors.New("account is deactivated")
	ErrInvalidOldPassword = errors.New("old password is incorrect")
)
