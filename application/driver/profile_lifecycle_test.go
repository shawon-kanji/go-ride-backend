package driver

import (
	"context"
	"errors"
	"testing"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

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
	uc := NewUpdateOnlineStatusUseCase(repo)

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
	uc := NewUpdateOnlineStatusUseCase(repo)

	_, err := uc.Execute(context.Background(), uuid.New().String(), UpdateOnlineStatusRequest{IsOnline: true})
	if !errors.Is(err, domaindriver.ErrDriverNotFound) {
		t.Fatalf("expected ErrDriverNotFound, got %v", err)
	}
}
