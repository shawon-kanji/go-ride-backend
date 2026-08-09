package kyc

import (
	"context"
	"fmt"
	"strings"

	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

type ConfirmUploadUseCase struct {
	repo     domainkyc.Repository
	storage  StoragePort
	drivers  DriverKycUpdater
	vehicles VehicleOwnerChecker
}

func NewConfirmUploadUseCase(repo domainkyc.Repository, storage StoragePort, drivers DriverKycUpdater, vehicles VehicleOwnerChecker) *ConfirmUploadUseCase {
	return &ConfirmUploadUseCase{repo: repo, storage: storage, drivers: drivers, vehicles: vehicles}
}

func (uc *ConfirmUploadUseCase) Execute(ctx context.Context, driverID string, req ConfirmUploadRequest) (*DocumentResponse, error) {
	if err := validateConfirmUploadRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	var vehicleID *uuid.UUID
	if domainkyc.IsVehicleDocumentType(req.DocumentType) {
		vID, err := uuid.Parse(req.VehicleID)
		if err != nil {
			return nil, ValidationError{Message: "invalid vehicle id"}
		}
		owned, err := uc.vehicles.IsOwnedByDriver(ctx, vID, id)
		if err != nil {
			return nil, err
		}
		if !owned {
			return nil, domainkyc.ErrDocumentForbidden
		}
		vehicleID = &vID
	}

	// The key must actually belong to this driver+document type (and, for
	// vehicle-scoped types, this vehicle), so a driver can't confirm an
	// upload URL scoped to someone else's key.
	expectedPrefix := fmt.Sprintf("drivers/%s/%s/", id, req.DocumentType)
	if vehicleID != nil {
		expectedPrefix = fmt.Sprintf("drivers/%s/%s/%s/", id, req.DocumentType, vehicleID)
	}
	if !strings.HasPrefix(req.Key, expectedPrefix) {
		return nil, ValidationError{Message: "key does not match driver id, document type, and vehicle id"}
	}

	exists, err := uc.storage.Exists(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainkyc.ErrUploadNotConfirmed
	}

	created, err := uc.repo.Create(ctx, &domainkyc.CreateInput{
		DriverID:     id,
		VehicleID:    vehicleID,
		DocumentType: req.DocumentType,
		FileURL:      req.Key,
	})
	if err != nil {
		return nil, err
	}

	// Only identity documents feed the driver's aggregate kyc_status;
	// vehicle-scoped documents are tracked per vehicle_id instead (see
	// domain/kyc.Repository.IsVehicleDocumentApproved) and never touch
	// drivers.kyc_status.
	if !domainkyc.IsVehicleDocumentType(req.DocumentType) {
		if err := uc.recomputeKycStatus(ctx, id); err != nil {
			return nil, err
		}
	}

	resp := ToDocumentResponse(created)
	return &resp, nil
}

// recomputeKycStatus moves the driver to in_review whenever every required
// identity document type has a current row — including when re-uploading
// one document after a manual approval, which deliberately pulls
// kyc_status back from approved to in_review: the newly uploaded document
// hasn't been reviewed yet, even though its siblings were. It only ever
// sets in_review; it never sets rejected/approved, since those are applied
// outside the API for now per current scope.
func (uc *ConfirmUploadUseCase) recomputeKycStatus(ctx context.Context, driverID uuid.UUID) error {
	current, err := uc.repo.ListCurrentByDriver(ctx, driverID)
	if err != nil {
		return err
	}

	present := make(map[string]bool, len(current))
	for _, d := range current {
		present[d.DocumentType] = true
	}

	for _, t := range domainkyc.RequiredIdentityDocumentTypes {
		if !present[t] {
			return nil
		}
	}

	_, err = uc.drivers.UpdateKycStatus(ctx, driverID, domaindriver.KycStatusInReview)
	return err
}
