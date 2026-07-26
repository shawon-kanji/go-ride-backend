package vehicle

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

func validateRegisterVehicleRequest(req RegisterVehicleRequest) error {
	req.PlateNumber = strings.TrimSpace(req.PlateNumber)
	req.Color = strings.TrimSpace(req.Color)
	req.ModelName = strings.TrimSpace(req.ModelName)
	req.Category = strings.ToLower(strings.TrimSpace(req.Category))

	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid register vehicle request"}
	}

	return nil
}

func validateUpdateVehicleRequest(req UpdateVehicleRequest) error {
	req.PlateNumber = strings.TrimSpace(req.PlateNumber)
	req.Color = strings.TrimSpace(req.Color)
	req.ModelName = strings.TrimSpace(req.ModelName)
	req.Category = strings.ToLower(strings.TrimSpace(req.Category))

	if err := validate.Struct(req); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			first := verrs[0]
			return ValidationError{Message: formatFieldError(first)}
		}
		return ValidationError{Message: "invalid update vehicle request"}
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
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
