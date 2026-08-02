package driver

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type UpdateOnlineStatusUseCase struct {
	repo domaindriver.Repository
}

func NewUpdateOnlineStatusUseCase(repo domaindriver.Repository) *UpdateOnlineStatusUseCase {
	return &UpdateOnlineStatusUseCase{repo: repo}
}

func (uc *UpdateOnlineStatusUseCase) Execute(ctx context.Context, driverID string, req UpdateOnlineStatusRequest) (*DriverResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	updated, err := uc.repo.UpdateOnlineStatus(ctx, id, req.IsOnline)
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(updated)
	return &resp, nil
}
