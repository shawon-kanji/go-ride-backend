package driver

import (
	"context"
	"errors"
	"testing"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

func TestLoginUseCaseExecuteSuccess(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:              uuid.New(),
			Email:           "foo@example.com",
			PasswordHash:    "hashed:password123",
			FirstName:       "Foo",
			LastName:        "Bar",
			AccountStatus:   domaindriver.AccountStatusPending,
			IsEmailVerified: false,
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
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{
		"foo@example.com": {
			ID:              uuid.New(),
			Email:           "foo@example.com",
			PasswordHash:    "hashed:password123",
			FirstName:       "Foo",
			LastName:        "Bar",
			AccountStatus:   domaindriver.AccountStatusPending,
			IsEmailVerified: false,
		},
	}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), LoginRequest{Email: "foo@example.com", Password: "wrongpass"})
	if !errors.Is(err, domaindriver.ErrInvalidCredential) {
		t.Fatalf("expected ErrInvalidCredential, got %v", err)
	}
}

func TestLoginUseCaseExecuteValidation(t *testing.T) {
	repo := &fakeDriverRepo{byEmail: map[string]*domaindriver.Driver{}}
	hasher := &fakeHasher{}
	tokens := &fakeTokens{}
	uc := NewLoginUseCase(repo, hasher, tokens)

	_, err := uc.Execute(context.Background(), LoginRequest{Email: "invalid", Password: "short"})
	var validationErr ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}
