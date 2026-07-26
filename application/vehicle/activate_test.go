package vehicle

import (
	"context"
	"errors"
	"testing"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

func TestActivateUseCaseExecuteDeactivatesPreviousVehicle(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverID := uuid.New()

	first, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register first vehicle: %v", err)
	}
	second, err := registerUC.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "A-2", Color: "Blue", ModelName: "Civic", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register second vehicle: %v", err)
	}

	activateUC := NewActivateUseCase(repo)
	if _, err := activateUC.Execute(context.Background(), driverID.String(), first.ID); err != nil {
		t.Fatalf("activate first vehicle: %v", err)
	}
	resp, err := activateUC.Execute(context.Background(), driverID.String(), second.ID)
	if err != nil {
		t.Fatalf("activate second vehicle: %v", err)
	}
	if !resp.IsActive {
		t.Fatalf("expected second vehicle to be active")
	}

	firstReloaded, err := repo.GetByID(context.Background(), uuid.MustParse(first.ID))
	if err != nil {
		t.Fatalf("reload first vehicle: %v", err)
	}
	if firstReloaded.IsActive {
		t.Fatalf("expected first vehicle to be deactivated after activating second")
	}
}

func TestActivateUseCaseExecuteForbiddenForOtherDriver(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	ownerID := uuid.New()

	vehicle, err := registerUC.Execute(context.Background(), ownerID.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	})
	if err != nil {
		t.Fatalf("register vehicle: %v", err)
	}

	activateUC := NewActivateUseCase(repo)
	_, err = activateUC.Execute(context.Background(), uuid.New().String(), vehicle.ID)
	if !errors.Is(err, domainvehicle.ErrVehicleForbidden) {
		t.Fatalf("expected ErrVehicleForbidden, got %v", err)
	}
}
