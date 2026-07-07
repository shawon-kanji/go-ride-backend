package user

import (
	"context"
	"errors"
	"strings"

	domainuser "go-ride-backend/domain/user"
)

type SignupUseCase struct {
	repo   domainuser.Repository
	hasher PasswordHasher
}

func NewSignupUseCase(repo domainuser.Repository, hasher PasswordHasher) *SignupUseCase {
	return &SignupUseCase{repo: repo, hasher: hasher}
}

func (uc *SignupUseCase) Execute(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	if err := validateSignupRequest(req); err != nil {
		return nil, err
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domainuser.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domainuser.ErrEmailAlreadyTaken
	}

	hashed, err := uc.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	input := NewUserFromSignup(req, hashed)
	createdUser, err := uc.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return &SignupResponse{User: ToUserResponse(createdUser)}, nil
}
