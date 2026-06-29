package driver

import (
	"context"
	"errors"
	"strings"

	domaindriver "go-ride-backend/domain/driver"
)

type SignupUseCase struct {
	repo   domaindriver.Repository
	hasher PasswordHasher
}

func NewSignupUseCase(repo domaindriver.Repository, hasher PasswordHasher) *SignupUseCase {
	return &SignupUseCase{repo: repo, hasher: hasher}
}

func (uc *SignupUseCase) Execute(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	if err := validateSignupRequest(req); err != nil {
		return nil, err
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domaindriver.ErrDriverNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domaindriver.ErrEmailAlreadyTaken
	}

	hashed, err := uc.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	newDriver := NewDriverFromSignup(req, hashed)
	if err := uc.repo.Create(ctx, newDriver); err != nil {
		return nil, err
	}

	return &SignupResponse{Driver: ToDriverResponse(newDriver)}, nil
}
