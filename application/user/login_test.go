package user

import (
	"context"
	"errors"
	"testing"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

func TestLoginUseCaseExecuteSuccess(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:           uuid.New(),
			Email:        "foo@example.com",
			PasswordHash: "hashed:password123",
			FirstName:    "Foo",
			LastName:     "Bar",
		},
	}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	resp, err := uc.Execute(context.Background(), LoginRequest{Email: "foo@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatalf("expected access token to be set")
	}
}

func TestLoginUseCaseExecuteInvalidCredentials(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:           uuid.New(),
			Email:        "foo@example.com",
			PasswordHash: "hashed:password123",
			FirstName:    "Foo",
			LastName:     "Bar",
		},
	}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), LoginRequest{Email: "foo@example.com", Password: "wrongpass"})
	if !errors.Is(err, domainuser.ErrInvalidCredential) {
		t.Fatalf("expected ErrInvalidCredential, got %v", err)
	}
}

func TestLoginUseCaseExecuteValidation(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), LoginRequest{Email: "invalid", Password: "short"})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestLoginUseCaseExecuteDeactivatedAccount(t *testing.T) {
	repo := &fakeUserRepo{byEmail: map[string]*domainuser.User{
		"foo@example.com": {
			ID:            uuid.New(),
			Email:         "foo@example.com",
			PasswordHash:  "hashed:password123",
			FirstName:     "Foo",
			LastName:      "Bar",
			AccountStatus: domainuser.AccountStatusDeactivated,
		},
	}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), LoginRequest{Email: "foo@example.com", Password: "password123"})
	if !errors.Is(err, domainuser.ErrAccountDeactivated) {
		t.Fatalf("expected ErrAccountDeactivated, got %v", err)
	}
}
