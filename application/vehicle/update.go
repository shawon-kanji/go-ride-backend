package vehicle

import (
	"context"
	"errors"
	"strings"

	domainvehicle "go-ride-backend/domain/vehicle"
)

type UpdateUseCase struct {
	repo domainvehicle.Repository
}

func NewUpdateUseCase(repo domainvehicle.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, driverID string, vehicleID string, req UpdateVehicleRequest) (*VehicleResponse, error) {
	if err := validateUpdateVehicleRequest(req); err != nil {
		return nil, err
	}

	dID, vID, err := parseIDs(driverID, vehicleID)
	if err != nil {
		return nil, err
	}

	existing, err := loadOwnedVehicle(ctx, uc.repo, dID, vID)
	if err != nil {
		return nil, err
	}

	plateNumber := strings.TrimSpace(req.PlateNumber)
	if plateNumber != existing.PlateNumber {
		_, err := uc.repo.GetByPlateNumber(ctx, plateNumber)
		if err == nil {
			return nil, domainvehicle.ErrPlateAlreadyRegistered
		}
		if !errors.Is(err, domainvehicle.ErrVehicleNotFound) {
			return nil, err
		}
	}

	updated, err := uc.repo.Update(ctx, vID, &domainvehicle.UpdateInput{
		PlateNumber: plateNumber,
		Color:       strings.TrimSpace(req.Color),
		ModelName:   strings.TrimSpace(req.ModelName),
		SeatCount:   req.SeatCount,
		Category:    strings.ToLower(strings.TrimSpace(req.Category)),
	})
	if err != nil {
		return nil, err
	}

	resp := ToVehicleResponse(updated)
	return &resp, nil
}
