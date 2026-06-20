package user

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func validateSignupRequest(req SignupRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid signup request"}
	}

	return nil
}

func validateLoginRequest(req LoginRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid login request"}
	}

	return nil
}

func formatFieldError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
