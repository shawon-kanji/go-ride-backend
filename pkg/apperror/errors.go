package apperror

import (
	"errors"
	"net/http"

	appdriver "go-ride-backend/application/driver"
	appuser "go-ride-backend/application/user"
	appvehicle "go-ride-backend/application/vehicle"
	domaindriver "go-ride-backend/domain/driver"
	domainuser "go-ride-backend/domain/user"
	domainvehicle "go-ride-backend/domain/vehicle"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func Map(err error) HTTPError {
	var driverValidationErr appdriver.ValidationError
	var validationErr appuser.ValidationError
	var vehicleValidationErr appvehicle.ValidationError

	switch {
	case errors.As(err, &driverValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: driverValidationErr.Message}
	case errors.As(err, &validationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: validationErr.Message}
	case errors.As(err, &vehicleValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: vehicleValidationErr.Message}
	case errors.Is(err, domaindriver.ErrEmailAlreadyTaken):
		return HTTPError{Status: http.StatusConflict, Code: "EMAIL_ALREADY_TAKEN", Message: err.Error()}
	case errors.Is(err, domaindriver.ErrInvalidCredential):
		return HTTPError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: err.Error()}
	case errors.Is(err, domaindriver.ErrDriverNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "DRIVER_NOT_FOUND", Message: err.Error()}
	case errors.Is(err, domainuser.ErrEmailAlreadyTaken):
		return HTTPError{Status: http.StatusConflict, Code: "EMAIL_ALREADY_TAKEN", Message: err.Error()}
	case errors.Is(err, domainuser.ErrInvalidCredential):
		return HTTPError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: err.Error()}
	case errors.Is(err, domainuser.ErrUserNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: err.Error()}
	case errors.Is(err, domainuser.ErrAccountDeactivated):
		return HTTPError{Status: http.StatusForbidden, Code: "ACCOUNT_DEACTIVATED", Message: err.Error()}
	case errors.Is(err, domainuser.ErrInvalidOldPassword):
		return HTTPError{Status: http.StatusUnauthorized, Code: "INVALID_OLD_PASSWORD", Message: err.Error()}
	case errors.Is(err, domainvehicle.ErrPlateAlreadyRegistered):
		return HTTPError{Status: http.StatusConflict, Code: "PLATE_ALREADY_REGISTERED", Message: err.Error()}
	case errors.Is(err, domainvehicle.ErrVehicleForbidden):
		return HTTPError{Status: http.StatusForbidden, Code: "VEHICLE_FORBIDDEN", Message: err.Error()}
	case errors.Is(err, domainvehicle.ErrVehicleNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "VEHICLE_NOT_FOUND", Message: err.Error()}
	default:
		return HTTPError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "internal server error"}
	}
}
