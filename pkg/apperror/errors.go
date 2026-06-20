package apperror

import (
	"errors"
	"net/http"

	domainuser "go-ride-backend/domain/user"
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
	switch {
	case errors.Is(err, domainuser.ErrEmailAlreadyTaken):
		return HTTPError{Status: http.StatusConflict, Code: "EMAIL_ALREADY_TAKEN", Message: err.Error()}
	case errors.Is(err, domainuser.ErrInvalidCredential):
		return HTTPError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: err.Error()}
	default:
		return HTTPError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "internal server error"}
	}
}
