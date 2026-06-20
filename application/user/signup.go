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
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

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

	newUser := NewUserFromSignup(req, hashed)
	if err := uc.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return &SignupResponse{User: ToUserResponse(newUser)}, nil
}
