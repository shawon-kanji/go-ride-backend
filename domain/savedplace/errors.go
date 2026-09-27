package savedplace

import "errors"

var (
	ErrSavedPlaceNotFound  = errors.New("saved place not found")
	ErrSavedPlaceForbidden = errors.New("saved place does not belong to user")
	ErrPlaceTypeTaken      = errors.New("home or work is already saved")
)
