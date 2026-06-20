package user

import (
	"context"
	"errors"
	"strings"

	domainuser "go-ride-backend/domain/user"
)

type LoginUseCase struct {
	repo   domainuser.Repository
	hasher PasswordHasher
	tokens TokenManager
}

func NewLoginUseCase(repo domainuser.Repository, hasher PasswordHasher, tokens TokenManager) *LoginUseCase {
	return &LoginUseCase{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	usr, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, domainuser.ErrInvalidCredential
		}
		return nil, err
	}

	if err := uc.hasher.Compare(usr.PasswordHash, req.Password); err != nil {
		return nil, domainuser.ErrInvalidCredential
	}

	token, err := uc.tokens.Generate(usr.ID.String(), usr.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		User:        ToUserResponse(usr),
	}, nil
}
