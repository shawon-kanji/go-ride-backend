package vehicle

import (
	"context"

	domainkyc "go-ride-backend/domain/kyc"
	domainvehicle "go-ride-backend/domain/vehicle"
)

type ActivateUseCase struct {
	repo     domainvehicle.Repository
	drivers  DriverKycStatusReader
	vehicles VehicleDocumentApprovalReader
}

func NewActivateUseCase(repo domainvehicle.Repository, drivers DriverKycStatusReader, vehicles VehicleDocumentApprovalReader) *ActivateUseCase {
	return &ActivateUseCase{repo: repo, drivers: drivers, vehicles: vehicles}
}

// Execute activates vehicleID for driverID — this is also how a driver
// switches their active vehicle (activating a new one deactivates the
// previous one, see domainvehicle.Repository.Activate), so gating here
// covers both "go online with this vehicle" and "switch to this vehicle
// while online" with one check.
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

	approved, err := uc.drivers.IsKycApproved(ctx, dID)
	if err != nil {
		return nil, err
	}
	if !approved {
		return nil, domainkyc.ErrKycNotApproved
	}

	vehicleApproved, err := uc.vehicles.IsVehicleDocumentApproved(ctx, vID)
	if err != nil {
		return nil, err
	}
	if !vehicleApproved {
		return nil, domainkyc.ErrVehicleNotVerified
	}

	activated, err := uc.repo.Activate(ctx, dID, vID)
	if err != nil {
		return nil, err
	}

	resp := ToVehicleResponse(activated)
	return &resp, nil
}
