package savedplace

import (
	"context"
	"errors"
	"strings"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

type CreateUseCase struct {
	repo domainsavedplace.Repository
}

func NewCreateUseCase(repo domainsavedplace.Repository) *CreateUseCase {
	return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, userIDStr string, req CreateSavedPlaceRequest) (*SavedPlaceResponse, error) {
	if err := validateCreateSavedPlaceRequest(req); err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ValidationError{Message: "invalid user id"}
	}

	placeType := strings.ToLower(strings.TrimSpace(req.PlaceType))
	if placeType == domainsavedplace.PlaceTypeHome || placeType == domainsavedplace.PlaceTypeWork {
		_, err := uc.repo.GetByUserAndType(ctx, userID, placeType)
		if err == nil {
			return nil, domainsavedplace.ErrPlaceTypeTaken
		}
		if !errors.Is(err, domainsavedplace.ErrSavedPlaceNotFound) {
			return nil, err
		}
	}

	created, err := uc.repo.Create(ctx, &domainsavedplace.CreateInput{
		UserID:           userID,
		Label:            strings.TrimSpace(req.Label),
		PlaceType:        placeType,
		FormattedAddress: strings.TrimSpace(req.FormattedAddress),
		Lat:              req.Lat,
		Lng:              req.Lng,
	})
	if err != nil {
		return nil, err
	}

	resp := ToSavedPlaceResponse(created)
	return &resp, nil
}
