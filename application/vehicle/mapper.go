package vehicle

import (
	domainvehicle "go-ride-backend/domain/vehicle"
)

func ToVehicleResponse(v *domainvehicle.Vehicle) VehicleResponse {
	return VehicleResponse{
		ID:          v.ID.String(),
		DriverID:    v.DriverID.String(),
		PlateNumber: v.PlateNumber,
		Color:       v.Color,
		ModelName:   v.ModelName,
		SeatCount:   v.SeatCount,
		Category:    v.Category,
		IsActive:    v.IsActive,
	}
}
