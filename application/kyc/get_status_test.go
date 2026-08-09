package kyc

import (
	"context"
	"testing"

	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

type fakeDriverKycReader struct {
	byID map[uuid.UUID]*domaindriver.Driver
}

func (f *fakeDriverKycReader) GetByID(_ context.Context, id uuid.UUID) (*domaindriver.Driver, error) {
	if d, ok := f.byID[id]; ok {
		return d, nil
	}
	return nil, domaindriver.ErrDriverNotFound
}

func TestGetStatusUseCaseExecuteReturnsCurrentDocumentsAndStatus(t *testing.T) {
	repo := newFakeKycRepo()
	storage := newFakeStorage()
	drivers := newFakeDriverKycUpdater()
	confirmUC := NewConfirmUploadUseCase(repo, storage, drivers, newFakeVehicleOwnerChecker())
	driverID := uuid.New()

	key := "drivers/" + driverID.String() + "/selfie/" + uuid.New().String()
	storage.existingKeys[key] = true
	if _, err := confirmUC.Execute(context.Background(), driverID.String(), ConfirmUploadRequest{
		DocumentType: domainkyc.DocumentTypeSelfie,
		Key:          key,
	}); err != nil {
		t.Fatalf("confirm upload: %v", err)
	}

	readers := &fakeDriverKycReader{byID: map[uuid.UUID]*domaindriver.Driver{
		driverID: {ID: driverID, KycStatus: domaindriver.KycStatusNotStarted},
	}}
	getStatusUC := NewGetStatusUseCase(repo, readers)

	resp, err := getStatusUC.Execute(context.Background(), driverID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.KycStatus != domaindriver.KycStatusNotStarted {
		t.Fatalf("unexpected kyc status: %s", resp.KycStatus)
	}
	if len(resp.Documents) != 1 || resp.Documents[0].DocumentType != domainkyc.DocumentTypeSelfie {
		t.Fatalf("unexpected documents: %+v", resp.Documents)
	}
}
