package vehicle

import (
	"context"
	"errors"
	"testing"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

func TestUpdateUseCaseExecuteForbiddenForOtherDriver(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	ownerID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), ownerID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	updateUC := NewUpdateUseCase(repo)
	_, err = updateUC.Execute(context.Background(), uuid.New().String(), vehicle.ID, UpdateVehicleRequest{
		PlateNumber: "A-1", Color: "Green", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if !errors.Is(err, domainvehicle.ErrVehicleForbidden) {
		t.Fatalf("expected ErrVehicleForbidden, got %v", err)
	}
}

func TestUpdateUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	updateUC := NewUpdateUseCase(repo)
	resp, err := updateUC.Execute(context.Background(), driverID.String(), vehicle.ID, UpdateVehicleRequest{
		PlateNumber: "A-1", Color: "Green", ModelName: "Camry LE", SeatCount: 6, Category: "luxury",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Color != "Green" || resp.SeatCount != 6 || resp.Category != "luxury" {
		t.Fatalf("vehicle not updated: %+v", resp)
	}
}

func TestDeleteUseCaseExecuteForbiddenForOtherDriver(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	ownerID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), ownerID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	deleteUC := NewDeleteUseCase(repo)
	err = deleteUC.Execute(context.Background(), uuid.New().String(), vehicle.ID)
	if !errors.Is(err, domainvehicle.ErrVehicleForbidden) {
		t.Fatalf("expected ErrVehicleForbidden, got %v", err)
	}
}

func TestDeleteUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	deleteUC := NewDeleteUseCase(repo)
	if err := deleteUC.Execute(context.Background(), driverID.String(), vehicle.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := repo.GetByID(context.Background(), uuid.MustParse(vehicle.ID)); !errors.Is(err, domainvehicle.ErrVehicleNotFound) {
		t.Fatalf("expected vehicle to be deleted, got %v", err)
	}
}
