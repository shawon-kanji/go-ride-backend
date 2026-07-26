package vehicle

import (
	"context"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

func parseIDs(driverIDStr, vehicleIDStr string) (driverID uuid.UUID, vehicleID uuid.UUID, err error) {
	driverID, err = uuid.Parse(driverIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, ValidationError{Message: "invalid driver id"}
	}
	vehicleID, err = uuid.Parse(vehicleIDStr)
	if err != nil {
		return uuid.Nil, uuid.Nil, ValidationError{Message: "invalid vehicle id"}
	}
	return driverID, vehicleID, nil
}

// loadOwnedVehicle fetches a vehicle and verifies it belongs to driverID,
// returning ErrVehicleForbidden rather than leaking existence of vehicles
// owned by other drivers.
func loadOwnedVehicle(ctx context.Context, repo domainvehicle.Repository, driverID, vehicleID uuid.UUID) (*domainvehicle.Vehicle, error) {
	v, err := repo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if v.DriverID != driverID {
		return nil, domainvehicle.ErrVehicleForbidden
	}
	return v, nil
}
