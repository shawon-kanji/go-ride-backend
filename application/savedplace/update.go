package savedplace

import (
	"context"
	"errors"
	"strings"

	domainsavedplace "go-ride-backend/domain/savedplace"
)

type UpdateUseCase struct {
	repo domainsavedplace.Repository
}

func NewUpdateUseCase(repo domainsavedplace.Repository) *UpdateUseCase {
	return &UpdateUseCase{repo: repo}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, userID string, placeID string, req UpdateSavedPlaceRequest) (*SavedPlaceResponse, error) {
	if err := validateUpdateSavedPlaceRequest(req); err != nil {
		return nil, err
	}

	uID, pID, err := parseIDs(userID, placeID)
	if err != nil {
		return nil, err
	}

	existing, err := loadOwnedSavedPlace(ctx, uc.repo, uID, pID)
	if err != nil {
		return nil, err
	}

	placeType := strings.ToLower(strings.TrimSpace(req.PlaceType))
	if placeType != existing.PlaceType && (placeType == domainsavedplace.PlaceTypeHome || placeType == domainsavedplace.PlaceTypeWork) {
		other, err := uc.repo.GetByUserAndType(ctx, uID, placeType)
		if err == nil && other.ID != pID {
			return nil, domainsavedplace.ErrPlaceTypeTaken
		}
		if err != nil && !errors.Is(err, domainsavedplace.ErrSavedPlaceNotFound) {
			return nil, err
		}
	}

	updated, err := uc.repo.Update(ctx, pID, &domainsavedplace.UpdateInput{
		Label:            strings.TrimSpace(req.Label),
		PlaceType:        placeType,
		FormattedAddress: strings.TrimSpace(req.FormattedAddress),
		Lat:              req.Lat,
		Lng:              req.Lng,
	})
	if err != nil {
		return nil, err
	}

	resp := ToSavedPlaceResponse(updated)
	return &resp, nil
}
