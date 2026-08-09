package driver

import (
	"context"

	domaindriver "go-ride-backend/domain/driver"
	domainkyc "go-ride-backend/domain/kyc"
	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type UpdateOnlineStatusUseCase struct {
	repo             domaindriver.Repository
	vehicles         VehicleLister
	vehicleDocuments VehicleDocumentApprovalReader
}

func NewUpdateOnlineStatusUseCase(repo domaindriver.Repository, vehicles VehicleLister, vehicleDocuments VehicleDocumentApprovalReader) *UpdateOnlineStatusUseCase {
	return &UpdateOnlineStatusUseCase{repo: repo, vehicles: vehicles, vehicleDocuments: vehicleDocuments}
}

func (uc *UpdateOnlineStatusUseCase) Execute(ctx context.Context, driverID string, req UpdateOnlineStatusRequest) (*DriverResponse, error) {
	id, err := uuid.Parse(driverID)
	if err != nil {
		return nil, ValidationError{Message: "invalid driver id"}
	}

	if req.IsOnline {
		driver, err := uc.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if driver.KycStatus != domaindriver.KycStatusApproved {
			return nil, domainkyc.ErrKycNotApproved
		}

		vehicles, err := uc.vehicles.ListByDriver(ctx, id)
		if err != nil {
			return nil, err
		}
		activeVehicle := findActiveVehicle(vehicles)
		if activeVehicle == nil {
			return nil, domainvehicle.ErrNoActiveVehicle
		}

		vehicleApproved, err := uc.vehicleDocuments.IsVehicleDocumentApproved(ctx, activeVehicle.ID)
		if err != nil {
			return nil, err
		}
		if !vehicleApproved {
			return nil, domainkyc.ErrVehicleNotVerified
		}
	}

	updated, err := uc.repo.UpdateOnlineStatus(ctx, id, req.IsOnline)
	if err != nil {
		return nil, err
	}

	resp := ToDriverResponse(updated)
	return &resp, nil
}

func findActiveVehicle(vehicles []*domainvehicle.Vehicle) *domainvehicle.Vehicle {
	for _, v := range vehicles {
		if v.IsActive {
			return v
		}
	}
	return nil
}
