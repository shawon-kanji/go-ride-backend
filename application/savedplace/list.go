package savedplace

import (
	"context"

	domainsavedplace "go-ride-backend/domain/savedplace"

	"github.com/google/uuid"
)

type ListUseCase struct {
	repo domainsavedplace.Repository
}

func NewListUseCase(repo domainsavedplace.Repository) *ListUseCase {
	return &ListUseCase{repo: repo}
}

func (uc *ListUseCase) Execute(ctx context.Context, userIDStr string) (*ListSavedPlacesResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ValidationError{Message: "invalid user id"}
	}

	places, err := uc.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := ListSavedPlacesResponse{SavedPlaces: make([]SavedPlaceResponse, 0, len(places))}
	for _, p := range places {
		resp.SavedPlaces = append(resp.SavedPlaces, ToSavedPlaceResponse(p))
	}
	return &resp, nil
}
