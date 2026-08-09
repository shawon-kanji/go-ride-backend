package kyc

import "errors"

var (
	ErrDocumentNotFound    = errors.New("document not found")
	ErrDocumentForbidden   = errors.New("document does not belong to driver")
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrUploadNotConfirmed  = errors.New("uploaded file not found in storage")
	ErrKycNotApproved      = errors.New("driver kyc is not approved")
	ErrVehicleNotVerified  = errors.New("vehicle registration document is not approved")
)
