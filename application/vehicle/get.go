package vehicle

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"
)

type GetUseCase struct {
	repo domainvehicle.Repository
}

func NewGetUseCase(repo domainvehicle.Repository) *GetUseCase {
	return &GetUseCase{repo: repo}
}

func (uc *GetUseCase) Execute(ctx context.Context, driverID string, vehicleID string) (*VehicleResponse, error) {
	dID, vID, err := parseIDs(driverID, vehicleID)
	if err != nil {
		return nil, err
	}

	v, err := loadOwnedVehicle(ctx, uc.repo, dID, vID)
	if err != nil {
		return nil, err
	}

	resp := ToVehicleResponse(v)
	return &resp, nil
}
