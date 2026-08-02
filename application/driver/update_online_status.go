package driver

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"
	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type UpdateOnlineStatusUseCase struct {
	repo     domaindriver.Repository
	vehicles VehicleLister
}

func NewUpdateOnlineStatusUseCase(repo domaindriver.Repository, vehicles VehicleLister) *UpdateOnlineStatusUseCase {
	return &UpdateOnlineStatusUseCase{repo: repo, vehicles: vehicles}
}

func (uc *UpdateOnlineStatusUseCase) Execute(ctx context.Context, driverID string, req UpdateOnlineStatusRequest) (*DriverResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	if req.IsOnline {
		vehicles, err := uc.vehicles.ListByDriver(ctx, id)
		if err != nil {
			return nil, err
		}
		if !hasActiveVehicle(vehicles) {
			return nil, domainvehicle.ErrNoActiveVehicle
		}
	}

	updated, err := uc.repo.UpdateOnlineStatus(ctx, id, req.IsOnline)
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(updated)
	return &resp, nil
}

func hasActiveVehicle(vehicles []*domainvehicle.Vehicle) bool {
	for _, v := range vehicles {
		if v.IsActive {
			return true
		}
	}
	return false
}
