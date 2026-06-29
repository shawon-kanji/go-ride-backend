package user

import (
	"context"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type ChangePasswordUseCase struct {
	repo   domainuser.Repository
	hasher PasswordHasher
}

func NewChangePasswordUseCase(repo domainuser.Repository, hasher PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo, hasher: hasher}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID string, req ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if err := validateChangePasswordRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ValidationError{Message: "invalid user id"}
	}

	usr, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if usr.AccountStatus == domainuser.AccountStatusDeactivated {
		return nil, domainuser.ErrAccountDeactivated
	}

	if err := uc.hasher.Compare(usr.PasswordHash, req.OldPassword); err != nil {
		return nil, domainuser.ErrInvalidOldPassword
	}

	hashed, err := uc.hasher.Hash(req.NewPassword)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.UpdatePassword(ctx, id, hashed); err != nil {
		return nil, err
	}

	return &ChangePasswordResponse{Message: "password changed successfully"}, nil
}
