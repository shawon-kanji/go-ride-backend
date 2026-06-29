package user

import (
	"context"
	"strings"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type UpdateProfileUseCase struct {
	repo domainuser.Repository
}

func NewUpdateProfileUseCase(repo domainuser.Repository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: repo}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req UpdateProfileRequest) (*UserResponse, error) {
	if err := validateUpdateProfileRequest(req); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ValidationError{Message: "invalid user id"}
	}

	updated, err := uc.repo.UpdateProfile(ctx, id, strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName))
	if err != nil {
		return nil, err
	}

	resp := ToUserResponse(updated)
	return &resp, nil
}
