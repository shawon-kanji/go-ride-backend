package apperror

import (
	"errors"
	"net/http"

	appdriver "go-ride-backend/application/driver"
	appkyc "go-ride-backend/application/kyc"
	appplaces "go-ride-backend/application/places"
	appuser "go-ride-backend/application/user"
	appvehicle "go-ride-backend/application/vehicle"
	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"
	domainplaces "go-ride-backend/domain/places"
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
	var kycValidationErr appkyc.ValidationError
	var placesValidationErr appplaces.ValidationError

	switch {
	case errors.As(err, &driverValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: driverValidationErr.Message}
	case errors.As(err, &validationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: validationErr.Message}
	case errors.As(err, &vehicleValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: vehicleValidationErr.Message}
	case errors.As(err, &kycValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: kycValidationErr.Message}
	case errors.As(err, &placesValidationErr):
		return HTTPError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: placesValidationErr.Message}
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
	case errors.Is(err, domainvehicle.ErrNoActiveVehicle):
		return HTTPError{Status: http.StatusForbidden, Code: "NO_ACTIVE_VEHICLE", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrDocumentNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "DOCUMENT_NOT_FOUND", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrDocumentForbidden):
		return HTTPError{Status: http.StatusForbidden, Code: "DOCUMENT_FORBIDDEN", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrInvalidDocumentType):
		return HTTPError{Status: http.StatusBadRequest, Code: "INVALID_DOCUMENT_TYPE", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrUploadNotConfirmed):
		return HTTPError{Status: http.StatusBadRequest, Code: "UPLOAD_NOT_CONFIRMED", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrKycNotApproved):
		return HTTPError{Status: http.StatusForbidden, Code: "KYC_NOT_APPROVED", Message: err.Error()}
	case errors.Is(err, domainkyc.ErrVehicleNotVerified):
		return HTTPError{Status: http.StatusForbidden, Code: "VEHICLE_NOT_VERIFIED", Message: err.Error()}
	case errors.Is(err, domainplaces.ErrPlaceNotFound):
		return HTTPError{Status: http.StatusNotFound, Code: "PLACE_NOT_FOUND", Message: err.Error()}
	case errors.Is(err, domainplaces.ErrInvalidPlaceID):
		return HTTPError{Status: http.StatusBadRequest, Code: "INVALID_PLACE_ID", Message: err.Error()}
	case errors.Is(err, domainplaces.ErrProviderUnavailable):
		return HTTPError{Status: http.StatusBadGateway, Code: "PLACES_PROVIDER_UNAVAILABLE", Message: "places provider is temporarily unavailable"}
	default:
		return HTTPError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "internal server error"}
	}
}
