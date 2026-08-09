package kyc

import (
	"context"
	"fmt"

	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

type RequestUploadURLUseCase struct {
	storage StoragePort
}

func NewRequestUploadURLUseCase(storage StoragePort) *RequestUploadURLUseCase {
	return &RequestUploadURLUseCase{storage: storage}
}

// Execute returns a presigned PUT URL and the object key the driver must
// echo back on confirm. It deliberately does not write a driver_documents
// row yet — that only happens once ConfirmUpload verifies the file was
// actually written, so a driver who never uploads never leaves an orphaned
// row behind.
func (uc *RequestUploadURLUseCase) Execute(ctx context.Context, driverID string, req RequestUploadURLRequest) (*RequestUploadURLResponse, error) {
	if err := validateRequestUploadURLRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	// Vehicle ownership of req.VehicleID isn't checked here — this step
	// doesn't persist anything yet, so a URL scoped to a vehicle the driver
	// doesn't own is harmless; ConfirmUpload is where ownership is enforced
	// before any driver_documents row is written.
	var key string
	if domainkyc.IsVehicleDocumentType(req.DocumentType) {
		key = fmt.Sprintf("drivers/%s/%s/%s/%s", id, req.DocumentType, req.VehicleID, uuid.New())
	} else {
		key = fmt.Sprintf("drivers/%s/%s/%s", id, req.DocumentType, uuid.New())
	}

	url, err := uc.storage.PresignUploadURL(ctx, key, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &RequestUploadURLResponse{UploadURL: url, Key: key}, nil
}
