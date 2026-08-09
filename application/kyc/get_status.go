package kyc

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"

	"github.com/google/uuid"
)

// DriverKycReader abstracts fetching a driver's current kyc_status.
type DriverKycReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domaindriver.Driver, error)
}

type GetStatusUseCase struct {
	repo    domainkyc.Repository
	drivers DriverKycReader
}

func NewGetStatusUseCase(repo domainkyc.Repository, drivers DriverKycReader) *GetStatusUseCase {
	return &GetStatusUseCase{repo: repo, drivers: drivers}
}

func (uc *GetStatusUseCase) Execute(ctx context.Context, driverID string) (*KycStatusResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	driver, err := uc.drivers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	docs, err := uc.repo.ListCurrentByDriver(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := KycStatusResponse{
		KycStatus: driver.KycStatus,
		Documents: make([]DocumentResponse, 0, len(docs)),
	}
	for _, d := range docs {
		resp.Documents = append(resp.Documents, ToDocumentResponse(d))
	}
	return &resp, nil
}
