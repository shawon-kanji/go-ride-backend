package kyc

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

func TestRequestUploadURLUseCaseExecuteSuccess(t *testing.T) {
	storage := newFakeStorage()
	uc := NewRequestUploadURLUseCase(storage)
	driverID := uuid.New()

	resp, err := uc.Execute(context.Background(), driverID.String(), RequestUploadURLRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		ContentType:  "image/jpeg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	wantPrefix := "drivers/" + driverID.String() + "/" + domainkyc.DocumentTypeSelfie + "/"
	if !strings.HasPrefix(resp.Key, wantPrefix) {
		t.Fatalf("expected key to start with %q, got %q", wantPrefix, resp.Key)
	}
	if resp.UploadURL == "" {
		t.Fatalf("expected a non-empty upload url")
	}
}

func TestRequestUploadURLUseCaseExecuteVehicleRegistrationScopesKeyByVehicle(t *testing.T) {
	storage := newFakeStorage()
	uc := NewRequestUploadURLUseCase(storage)
	driverID := uuid.New()
	vehicleID := uuid.New()

	resp, err := uc.Execute(context.Background(), driverID.String(), RequestUploadURLRequest{
		DocumentType: domainkyc.DocumentTypeVehicleRegistration,
		ContentType:  "image/jpeg",
		VehicleID:    vehicleID.String(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	wantPrefix := "drivers/" + driverID.String() + "/" + domainkyc.DocumentTypeVehicleRegistration + "/" + vehicleID.String() + "/"
	if !strings.HasPrefix(resp.Key, wantPrefix) {
		t.Fatalf("expected key to start with %q, got %q", wantPrefix, resp.Key)
	}
}

func TestRequestUploadURLUseCaseExecuteVehicleRegistrationRequiresVehicleID(t *testing.T) {
	storage := newFakeStorage()
	uc := NewRequestUploadURLUseCase(storage)

	_, err := uc.Execute(context.Background(), uuid.New().String(), RequestUploadURLRequest{
		DocumentType: domainkyc.DocumentTypeVehicleRegistration,
		ContentType:  "image/jpeg",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestRequestUploadURLUseCaseExecuteInvalidDocumentType(t *testing.T) {
	storage := newFakeStorage()
	uc := NewRequestUploadURLUseCase(storage)

	_, err := uc.Execute(context.Background(), uuid.New().String(), RequestUploadURLRequest{
		DocumentType: "passport",
		ContentType:  "image/jpeg",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
