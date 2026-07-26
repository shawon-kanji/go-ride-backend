package vehicle

import (
	"context"
	"errors"
	"testing"
	"time"

	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type fakeVehicleRepo struct {
	byID map[uuid.UUID]*domainvehicle.Vehicle
}

func newFakeVehicleRepo() *fakeVehicleRepo {
	return &fakeVehicleRepo{byID: map[uuid.UUID]*domainvehicle.Vehicle{}}
}

func (r *fakeVehicleRepo) Create(_ context.Context, input *domainvehicle.CreateInput) (*domainvehicle.Vehicle, error) {
	v := &domainvehicle.Vehicle{
		ID:          uuid.New(),
		DriverID:    input.DriverID,
		PlateNumber: input.PlateNumber,
		Color:       input.Color,
		ModelName:   input.ModelName,
		SeatCount:   input.SeatCount,
		Category:    input.Category,
		IsActive:    false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	r.byID[v.ID] = v
	return v, nil
}

func (r *fakeVehicleRepo) GetByID(_ context.Context, id uuid.UUID) (*domainvehicle.Vehicle, error) {
	if v, ok := r.byID[id]; ok {
		return v, nil
	}
	return nil, domainvehicle.ErrVehicleNotFound
}

func (r *fakeVehicleRepo) GetByPlateNumber(_ context.Context, plate string) (*domainvehicle.Vehicle, error) {
	for _, v := range r.byID {
		if v.PlateNumber == plate {
			return v, nil
		}
	}
	return nil, domainvehicle.ErrVehicleNotFound
}

func (r *fakeVehicleRepo) ListByDriver(_ context.Context, driverID uuid.UUID) ([]*domainvehicle.Vehicle, error) {
	var out []*domainvehicle.Vehicle
	for _, v := range r.byID {
		if v.DriverID == driverID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (r *fakeVehicleRepo) Update(_ context.Context, id uuid.UUID, input *domainvehicle.UpdateInput) (*domainvehicle.Vehicle, error) {
	v, ok := r.byID[id]
	if !ok {
		return nil, domainvehicle.ErrVehicleNotFound
	}
	v.PlateNumber = input.PlateNumber
	v.Color = input.Color
	v.ModelName = input.ModelName
	v.SeatCount = input.SeatCount
	v.Category = input.Category
	return v, nil
}

func (r *fakeVehicleRepo) Activate(_ context.Context, driverID uuid.UUID, vehicleID uuid.UUID) (*domainvehicle.Vehicle, error) {
	target, ok := r.byID[vehicleID]
	if !ok || target.DriverID != driverID {
		return nil, domainvehicle.ErrVehicleNotFound
	}
	for _, v := range r.byID {
		if v.DriverID == driverID && v.ID != vehicleID {
			v.IsActive = false
		}
	}
	target.IsActive = true
	return target, nil
}

func (r *fakeVehicleRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := r.byID[id]; !ok {
		return domainvehicle.ErrVehicleNotFound
	}
	delete(r.byID, id)
	return nil
}

func TestRegisterUseCaseExecuteSuccess(t *testing.T) {
	repo := newFakeVehicleRepo()
	uc := NewRegisterUseCase(repo)
	driverID := uuid.New()

	resp, err := uc.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "ABC-123",
		Color:       "Red",
		ModelName:   "Toyota Camry",
		SeatCount:   4,
		Category:    "normal",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.PlateNumber != "ABC-123" || resp.DriverID != driverID.String() {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestRegisterUseCaseExecuteDuplicatePlate(t *testing.T) {
	repo := newFakeVehicleRepo()
	uc := NewRegisterUseCase(repo)
	driverID := uuid.New()

	if _, err := uc.Execute(context.Background(), driverID.String(), RegisterVehicleRequest{
		PlateNumber: "ABC-123", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "normal",
	}); err != nil {
		t.Fatalf("expected no error on first register, got %v", err)
	}

	_, err := uc.Execute(context.Background(), uuid.New().String(), RegisterVehicleRequest{
		PlateNumber: "ABC-123", Color: "Blue", ModelName: "Civic", SeatCount: 4, Category: "normal",
	})
	if !errors.Is(err, domainvehicle.ErrPlateAlreadyRegistered) {
		t.Fatalf("expected ErrPlateAlreadyRegistered, got %v", err)
	}
}

func TestRegisterUseCaseExecuteValidation(t *testing.T) {
	repo := newFakeVehicleRepo()
	uc := NewRegisterUseCase(repo)

	_, err := uc.Execute(context.Background(), uuid.New().String(), RegisterVehicleRequest{
		PlateNumber: "A", Color: "Red", ModelName: "Camry", SeatCount: 4, Category: "sports",
	})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
