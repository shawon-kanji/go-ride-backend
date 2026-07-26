package vehicle

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type ListUseCase struct {
	repo domainvehicle.Repository
}

func NewListUseCase(repo domainvehicle.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(ctx context.Context, driverID string) (*ListVehiclesResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	vehicles, err := uc.repo.ListByDriver(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := ListVehiclesResponse{Vehicles: make([]VehicleResponse, 0, len(vehicles))}
	for _, v := range vehicles {
		resp.Vehicles = append(resp.Vehicles, ToVehicleResponse(v))
	}
	return &resp, nil
}
