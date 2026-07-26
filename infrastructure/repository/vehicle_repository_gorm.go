package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainvehicle "go-ride-backend/domain/vehicle"

	schemamodels "github.com/shawon-kanji/go-ride-db-schema/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleRepositoryGorm struct {
	db *gorm.DB
}

func NewVehicleRepositoryGorm(db *gorm.DB) *VehicleRepositoryGorm {
	return &VehicleRepositoryGorm{db: db}
}

func toDomainVehicle(model *schemamodels.Vehicle) *domainvehicle.Vehicle {
	return &domainvehicle.Vehicle{
		ID:          model.ID,
		DriverID:    model.DriverID,
		PlateNumber: model.PlateNumber,
		Color:       model.Color,
		ModelName:   model.ModelName,
		SeatCount:   model.SeatCount,
		Category:    model.Category,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (r *VehicleRepositoryGorm) Create(ctx context.Context, input *domainvehicle.CreateInput) (*domainvehicle.Vehicle, error) {
	now := time.Now().UTC()
	model := &schemamodels.Vehicle{
		ID:          uuid.New(),
		DriverID:    input.DriverID,
		PlateNumber: input.PlateNumber,
		Color:       input.Color,
		ModelName:   input.ModelName,
		SeatCount:   input.SeatCount,
		Category:    input.Category,
		IsActive:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("create vehicle: %w", err)
	}
	return toDomainVehicle(model), nil
}

func (r *VehicleRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*domainvehicle.Vehicle, error) {
	var model schemamodels.Vehicle
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainvehicle.ErrVehicleNotFound
		}
		return nil, fmt.Errorf("find vehicle by id: %w", err)
	}
	return toDomainVehicle(&model), nil
}

func (r *VehicleRepositoryGorm) GetByPlateNumber(ctx context.Context, plate string) (*domainvehicle.Vehicle, error) {
	var model schemamodels.Vehicle
	if err := r.db.WithContext(ctx).Where("plate_number = ?", plate).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainvehicle.ErrVehicleNotFound
		}
		return nil, fmt.Errorf("find vehicle by plate number: %w", err)
	}
	return toDomainVehicle(&model), nil
}

func (r *VehicleRepositoryGorm) ListByDriver(ctx context.Context, driverID uuid.UUID) ([]*domainvehicle.Vehicle, error) {
	var rows []schemamodels.Vehicle
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list vehicles by driver: %w", err)
	}
	vehicles := make([]*domainvehicle.Vehicle, 0, len(rows))
	for i := range rows {
		vehicles = append(vehicles, toDomainVehicle(&rows[i]))
	}
	return vehicles, nil
}

func (r *VehicleRepositoryGorm) Update(ctx context.Context, id uuid.UUID, input *domainvehicle.UpdateInput) (*domainvehicle.Vehicle, error) {
	updates := map[string]interface{}{
		"plate_number": input.PlateNumber,
		"color":        input.Color,
		"model_name":   input.ModelName,
		"seat_count":   input.SeatCount,
		"category":     input.Category,
		"updated_at":   time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&schemamodels.Vehicle{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domainvehicle.ErrVehicleNotFound
	}

	return r.GetByID(ctx, id)
}

// Activate deactivates any other active vehicle for driverID, then
// activates vehicleID, in one transaction. Worst-case under true
// concurrency is "last write wins" on which car is active — not a
// safety/financial-critical race — so a plain transaction (no SELECT FOR
// UPDATE) is enough; the partial unique index on (driver_id) WHERE
// is_active is the real backstop.
func (r *VehicleRepositoryGorm) Activate(ctx context.Context, driverID uuid.UUID, vehicleID uuid.UUID) (*domainvehicle.Vehicle, error) {
	now := time.Now().UTC()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&schemamodels.Vehicle{}).
			Where("driver_id = ? AND is_active = true AND id <> ?", driverID, vehicleID).
			Updates(map[string]interface{}{"is_active": false, "updated_at": now}).Error; err != nil {
			return fmt.Errorf("deactivate previous vehicle: %w", err)
		}

		result := tx.Model(&schemamodels.Vehicle{}).
			Where("id = ? AND driver_id = ?", vehicleID, driverID).
			Updates(map[string]interface{}{"is_active": true, "updated_at": now})
		if result.Error != nil {
			return fmt.Errorf("activate vehicle: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return domainvehicle.ErrVehicleNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, vehicleID)
}

func (r *VehicleRepositoryGorm) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&schemamodels.Vehicle{})
	if result.Error != nil {
		return fmt.Errorf("delete vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainvehicle.ErrVehicleNotFound
	}
	return nil
}
