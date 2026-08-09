package kyc

type RequestUploadURLRequest struct {
	DocumentType string `json:"document_type" validate:"required"`
	ContentType  string `json:"content_type" validate:"required"`
	// VehicleID is required when DocumentType is vehicle_registration
	// (identifies which vehicle this registration document verifies) and
	// must be omitted for every other, driver-level document type.
	VehicleID string `json:"vehicle_id,omitempty"`
}

type RequestUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
}

type ConfirmUploadRequest struct {
	DocumentType string `json:"document_type" validate:"required"`
	Key          string `json:"key" validate:"required"`
	// VehicleID, same rule as RequestUploadURLRequest.VehicleID.
	VehicleID string `json:"vehicle_id,omitempty"`
}

type DocumentResponse struct {
	ID              string  `json:"id"`
	DocumentType    string  `json:"document_type"`
	VehicleID       *string `json:"vehicle_id,omitempty"`
	Status          string  `json:"status"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

type KycStatusResponse struct {
	KycStatus string             `json:"kyc_status"`
	Documents []DocumentResponse `json:"documents"`
}
