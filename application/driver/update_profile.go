package driver

import (
	"context"
	"strings"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type UpdateProfileUseCase struct {
	repo domaindriver.Repository
}

func NewUpdateProfileUseCase(repo domaindriver.Repository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: repo}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, driverID string, req UpdateProfileRequest) (*DriverResponse, error) {
	if err := validateUpdateProfileRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	updated, err := uc.repo.UpdateProfile(ctx, id, strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName))
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(updated)
	return &resp, nil
}
