package vehicle

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestListUseCaseExecuteReturnsOnlyDriversVehicles(t *testing.T) {
	repo := newFakeVehicleRepo()
	registerUC := NewRegisterUseCase(repo)
	driverA := uuid.New()
	driverB := uuid.New()

	if _, err := registerUC.Execute(context.Background(), driverA.String(), RegisterVehicleRequest{
		PlateNumber: "A-1", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	}); err != nil {
		t.Fatalf("register driverA vehicle: %v", err)
	}
	if _, err := registerUC.Execute(context.Background(), driverA.String(), RegisterVehicleRequest{
		PlateNumber: "A-2", Color: "Blue", ModelName: "Camry", SeatCount: 6, Category: "luxury",
	}); err != nil {
		t.Fatalf("register second driverA vehicle: %v", err)
	}
	if _, err := registerUC.Execute(context.Background(), driverB.String(), RegisterVehicleRequest{
		PlateNumber: "B-1", Color: "Black", ModelName: "Civic", SeatCount: 4, Category: "normal",
	}); err != nil {
		t.Fatalf("register driverB vehicle: %v", err)
	}

	listUC := NewListUseCase(repo)
	resp, err := listUC.Execute(context.Background(), driverA.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Vehicles) != 2 {
		t.Fatalf("expected 2 vehicles for driverA, got %d", len(resp.Vehicles))
	}
}
