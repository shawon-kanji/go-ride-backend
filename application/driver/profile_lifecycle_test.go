package driver

import (
	"context"
	"errors"
	"testing"

	domaindriver "go-ride-backend/domain/driver"
	domainvehicle "go-ride-backend/domain/vehicle"

	"github.com/google/uuid"
)

type fakeVehicleLister struct {
	byDriver map[uuid.UUID][]*domainvehicle.Vehicle
}

func (f *fakeVehicleLister) ListByDriver(_ context.Context, driverID uuid.UUID) ([]*domainvehicle.Vehicle, error) {
	return f.byDriver[driverID], nil
}

func TestGetProfileUseCaseExecuteSuccess(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:            driverID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domaindriver.AccountStatusActive,
		},
	}}
	uc := NewGetProfileUseCase(repo)

	resp, err := uc.Execute(context.Background(), driverID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.ID != driverID.String() {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUpdateProfileUseCaseExecuteSuccess(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:            driverID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domaindriver.AccountStatusActive,
		},
	}}
	uc := NewUpdateProfileUseCase(repo)

	resp, err := uc.Execute(context.Background(), driverID.String(), UpdateProfileRequest{FirstName: "New", LastName: "Name"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.FirstName != "New" || resp.LastName != "Name" {
		t.Fatalf("profile not updated: %+v", resp)
	}
}

func TestUpdateOnlineStatusUseCaseExecuteSuccess(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:            driverID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domaindriver.AccountStatusActive,
			IsOnline:      false,
		},
	}}
	vehicles := &fakeVehicleLister{byDriver: map[uuid.UUID][]*domainvehicle.Vehicle{
		driverID: {{ID: uuid.New(), DriverID: driverID, IsActive: true}},
	}}
	uc := NewUpdateOnlineStatusUseCase(repo, vehicles)

	resp, err := uc.Execute(context.Background(), driverID.String(), UpdateOnlineStatusRequest{IsOnline: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.IsOnline {
		t.Fatalf("online status not updated: %+v", resp)
	}
}

func TestUpdateOnlineStatusUseCaseExecuteDriverNotFound(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	uc := NewUpdateOnlineStatusUseCase(repo, &fakeVehicleLister{})

	_, err := uc.Execute(context.Background(), uuid.New().String(), UpdateOnlineStatusRequest{IsOnline: false})
	if !errors.Is(err, domaindriver.ErrDriverNotFound) {
		t.Fatalf("expected ErrDriverNotFound, got %v", err)
	}
}

func TestUpdateOnlineStatusUseCaseExecuteBlockedNoActiveVehicle(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:            driverID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domaindriver.AccountStatusActive,
			IsOnline:      false,
		},
	}}
	vehicles := &fakeVehicleLister{byDriver: map[uuid.UUID][]*domainvehicle.Vehicle{
		driverID: {{ID: uuid.New(), DriverID: driverID, IsActive: false}},
	}}
	uc := NewUpdateOnlineStatusUseCase(repo, vehicles)

	_, err := uc.Execute(context.Background(), driverID.String(), UpdateOnlineStatusRequest{IsOnline: true})
	if !errors.Is(err, domainvehicle.ErrNoActiveVehicle) {
		t.Fatalf("expected ErrNoActiveVehicle, got %v", err)
	}
}

func TestUpdateOnlineStatusUseCaseExecuteOfflineAlwaysAllowed(t *testing.T) {
	driverID := uuid.New()
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:            driverID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domaindriver.AccountStatusActive,
			IsOnline:      true,
		},
	}}
	uc := NewUpdateOnlineStatusUseCase(repo, &fakeVehicleLister{})

	resp, err := uc.Execute(context.Background(), driverID.String(), UpdateOnlineStatusRequest{IsOnline: false})
	if err != nil {
		t.Fatalf("expected no error going offline, got %v", err)
	}
	if resp.IsOnline {
		t.Fatalf("expected driver offline, got %+v", resp)
	}
}
