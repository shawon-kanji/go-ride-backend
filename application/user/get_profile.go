package user

import (
	"context"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	repo domainuser.Repository
}

func NewGetProfileUseCase(repo domainuser.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{repo: repo}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ValidationError{Message: "invalid user id"}
	}

	usr, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:            usr.ID.String(),
		Email:         usr.Email,
		FirstName:     usr.FirstName,
		LastName:      usr.LastName,
		AccountStatus: usr.AccountStatus,
	}, nil
}
