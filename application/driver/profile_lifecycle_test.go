package driver

import (
	"context"
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
