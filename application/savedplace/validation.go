package savedplace

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

func validateCreateSavedPlaceRequest(req CreateSavedPlaceRequest) error {
	req.Label = strings.TrimSpace(req.Label)
	req.PlaceType = strings.ToLower(strings.TrimSpace(req.PlaceType))
	req.FormattedAddress = strings.TrimSpace(req.FormattedAddress)

	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid create saved place request"}
	}

	return nil
}

func validateUpdateSavedPlaceRequest(req UpdateSavedPlaceRequest) error {
	req.Label = strings.TrimSpace(req.Label)
	req.PlaceType = strings.ToLower(strings.TrimSpace(req.PlaceType))
	req.FormattedAddress = strings.TrimSpace(req.FormattedAddress)

	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid update saved place request"}
	}

	return nil
}

func formatFieldError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "latitude":
		return fmt.Sprintf("%s must be a valid latitude", field)
	case "longitude":
		return fmt.Sprintf("%s must be a valid longitude", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
