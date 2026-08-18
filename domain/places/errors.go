package places

import "errors"

var (
	ErrPlaceNotFound       = errors.New("place not found")
	ErrInvalidPlaceID      = errors.New("invalid place id")
	ErrProviderUnavailable = errors.New("places provider unavailable")
)
