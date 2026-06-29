package driver

import (
	"context"
	"errors"
	"strings"

	domaindriver "go-ride-backend/domain/driver"
)

type LoginUseCase struct {
	repo   domaindriver.Repository
	hasher PasswordHasher
	tokens TokenManager
}

func NewLoginUseCase(repo domaindriver.Repository, hasher PasswordHasher, tokens TokenManager) *LoginUseCase {
	return &LoginUseCase{repo: repo, hasher: hasher, tokens: tokens}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	drv, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domaindriver.ErrDriverNotFound) {
			return nil, domaindriver.ErrInvalidCredential
		}
		return nil, err
	}

	if err := uc.hasher.Compare(drv.PasswordHash, req.Password); err != nil {
		return nil, domaindriver.ErrInvalidCredential
	}

	token, err := uc.tokens.Generate(drv.ID.String(), drv.Email, "driver")
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		Driver:      ToDriverResponse(drv),
	}, nil
}
