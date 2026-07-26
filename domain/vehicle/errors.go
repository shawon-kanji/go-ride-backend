package vehicle

import "errors"

var (
	ErrVehicleNotFound        = errors.New("vehicle not found")
	ErrVehicleForbidden       = errors.New("vehicle does not belong to driver")
	ErrPlateAlreadyRegistered = errors.New("plate number already registered")
)
