package vehicle

import (
	"context"
	"errors"
	"strings"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type RegisterUseCase struct {
	repo domainvehicle.Repository
}

func NewRegisterUseCase(repo domainvehicle.Repository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, driverID string, req RegisterVehicleRequest) (*VehicleResponse, error) {
	if err := validateRegisterVehicleRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	plateNumber := strings.TrimSpace(req.PlateNumber)

	_, err = uc.repo.GetByPlateNumber(ctx, plateNumber)
	if err == nil {
		return nil, domainvehicle.ErrPlateAlreadyRegistered
	}
	if !errors.Is(err, domainvehicle.ErrVehicleNotFound) {
		return nil, err
	}

	created, err := uc.repo.Create(ctx, &domainvehicle.CreateInput{
		DriverID:    id,
		PlateNumber: plateNumber,
		Color:       strings.TrimSpace(req.Color),
		ModelName:   strings.TrimSpace(req.ModelName),
		SeatCount:   req.SeatCount,
		Category:    strings.ToLower(strings.TrimSpace(req.Category)),
	})
	if err != nil {
		return nil, err
	}

	resp := ToVehicleResponse(created)
	return &resp, nil
}
