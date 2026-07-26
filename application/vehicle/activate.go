package vehicle

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"
)

type ActivateUseCase struct {
	repo domainvehicle.Repository
}

func NewActivateUseCase(repo domainvehicle.Repository) *ActivateUseCase {
	return &ActivateUseCase{repo: repo}
}

func (uc *ActivateUseCase) Execute(ctx context.Context, driverID string, vehicleID string) (*VehicleResponse, error) {
	dID, vID, err := parseIDs(driverID, vehicleID)
	if err != nil {
		return nil, err
	}

	// Ownership check up front so a forbidden/not-found driver never
	// reaches the transactional activate/deactivate flip.
	if _, err := loadOwnedVehicle(ctx, uc.repo, dID, vID); err != nil {
		return nil, err
	}

	activated, err := uc.repo.Activate(ctx, dID, vID)
	if err != nil {
		return nil, err
	}

	resp := ToVehicleResponse(activated)
	return &resp, nil
}
