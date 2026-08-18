package places

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

func validateAutocompleteRequest(req AutocompleteRequest) error {
	if err := validate.Struct(req); err != nil {
		return structValidationError(err, "invalid autocomplete request")
	}
	return nil
}

func validatePlaceDetailsRequest(req PlaceDetailsRequest) error {
	if err := validate.Struct(req); err != nil {
		return structValidationError(err, "invalid place details request")
	}
	return nil
}

func validateReverseGeocodeRequest(req ReverseGeocodeRequest) error {
	if err := validate.Struct(req); err != nil {
		return structValidationError(err, "invalid reverse geocode request")
	}
	return nil
}

func structValidationError(err error, fallback string) error {
	if verrs, ok := err.(validator.ValidationErrors); ok {
		first := verrs[0]
		return ValidationError{Message: formatFieldError(first)}
	}
	return ValidationError{Message: fallback}
}

func formatFieldError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
