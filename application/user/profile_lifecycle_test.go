package user

import (
	"context"
	"errors"
	"testing"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

func TestGetProfileUseCaseExecuteSuccess(t *testing.T) {
	userID := uuid.New()
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:            userID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domainuser.AccountStatusActive,
		},
	}}
	uc := NewGetProfileUseCase(repo)

	resp, err := uc.Execute(context.Background(), userID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.ID != userID.String() {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestUpdateProfileUseCaseExecuteSuccess(t *testing.T) {
	userID := uuid.New()
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:            userID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domainuser.AccountStatusActive,
		},
	}}
	uc := NewUpdateProfileUseCase(repo)

	resp, err := uc.Execute(context.Background(), userID.String(), UpdateProfileRequest{FirstName: "New", LastName: "Name"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.FirstName != "New" || resp.LastName != "Name" {
		t.Fatalf("profile not updated: %+v", resp)
	}
}

func TestChangePasswordUseCaseExecuteInvalidOldPassword(t *testing.T) {
	userID := uuid.New()
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:            userID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domainuser.AccountStatusActive,
		},
	}}
	uc := NewChangePasswordUseCase(repo, &fakeHasher{})

	_, err := uc.Execute(context.Background(), userID.String(), ChangePasswordRequest{OldPassword: "wrongpass", NewPassword: "newpassword123"})
	if !errors.Is(err, domainuser.ErrInvalidOldPassword) {
		t.Fatalf("expected ErrInvalidOldPassword, got %v", err)
	}
}

func TestDeactivateAccountUseCaseExecuteSuccess(t *testing.T) {
	userID := uuid.New()
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:            userID,
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domainuser.AccountStatusActive,
		},
	}}
	uc := NewDeactivateAccountUseCase(repo)

	resp, err := uc.Execute(context.Background(), userID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Message == "" {
		t.Fatalf("expected response message")
	}

	updated, err := repo.GetByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected updated user, got %v", err)
	}
	if updated.AccountStatus != domainuser.AccountStatusDeactivated {
		t.Fatalf("expected account status deactivated, got %s", updated.AccountStatus)
	}
}
