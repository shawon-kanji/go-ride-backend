package driver

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	repo domaindriver.Repository
}

func NewGetProfileUseCase(repo domaindriver.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{repo: repo}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, driverID string) (*DriverResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	drv, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(drv)
	return &resp, nil
}
