package kyc

import (
	"fmt"
	"strings"

	domainkyc "go-ride-backend/domain/kyc"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func isValidDocumentType(documentType string) bool {
	if domainkyc.IsVehicleDocumentType(documentType) {
		return true
	}
	return domainkyc.IsIdentityDocumentType(documentType)
}

// validateVehicleIDPresence enforces that vehicle_id is supplied exactly
// when the document type is vehicle-scoped (vehicle_registration or any of
// the vehicle photo/plate types) and forbidden for every identity document
// type.
func validateVehicleIDPresence(documentType string, vehicleID string) error {
	if domainkyc.IsVehicleDocumentType(documentType) {
		if vehicleID == "" {
			return ValidationError{Message: "vehicle_id is required for this document type"}
		}
		if _, err := uuid.Parse(vehicleID); err != nil {
			return ValidationError{Message: "vehicle_id is invalid"}
		}
		return nil
	}
	if vehicleID != "" {
		return ValidationError{Message: "vehicle_id must not be set for this document type"}
	}
	return nil
}

func validateRequestUploadURLRequest(req RequestUploadURLRequest) error {
	if err := validate.Struct(req); err != nil {
		return structValidationError(err, "invalid upload url request")
	}
	if !isValidDocumentType(req.DocumentType) {
		return ValidationError{Message: "document_type is invalid"}
	}
	return validateVehicleIDPresence(req.DocumentType, req.VehicleID)
}

func validateConfirmUploadRequest(req ConfirmUploadRequest) error {
	if err := validate.Struct(req); err != nil {
		return structValidationError(err, "invalid confirm upload request")
	}
	if !isValidDocumentType(req.DocumentType) {
		return ValidationError{Message: "document_type is invalid"}
	}
	return validateVehicleIDPresence(req.DocumentType, req.VehicleID)
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
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
