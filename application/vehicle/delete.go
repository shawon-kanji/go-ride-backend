package vehicle

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"
)

type DeleteUseCase struct {
	repo domainvehicle.Repository
}

func NewDeleteUseCase(repo domainvehicle.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, driverID string, vehicleID string) error {
	dID, vID, err := parseIDs(driverID, vehicleID)
	if err != nil {
		return err
	}

	if _, err := loadOwnedVehicle(ctx, uc.repo, dID, vID); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, vID)
}
