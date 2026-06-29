package user

import (
	"context"
	"time"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type DeactivateAccountUseCase struct {
	repo domainuser.Repository
}

func NewDeactivateAccountUseCase(repo domainuser.Repository) *DeactivateAccountUseCase {
	return &DeactivateAccountUseCase{repo: repo}
}

func (uc *DeactivateAccountUseCase) Execute(ctx context.Context, userID string) (*DeactivateAccountResponse, error) {
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

	if err := uc.repo.Deactivate(ctx, id, time.Now().UTC()); err != nil {
		return nil, err
	}

	return &DeactivateAccountResponse{Message: "account deactivated successfully"}, nil
}
