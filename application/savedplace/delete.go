package savedplace

import (
	"context"

	domainsavedplace "go-ride-backend/domain/savedplace"
)

type DeleteUseCase struct {
	repo domainsavedplace.Repository
}

func NewDeleteUseCase(repo domainsavedplace.Repository) *DeleteUseCase {
	return &DeleteUseCase{repo: repo}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, userID string, placeID string) error {
	uID, pID, err := parseIDs(userID, placeID)
	if err != nil {
		return err
	}

	if _, err := loadOwnedSavedPlace(ctx, uc.repo, uID, pID); err != nil {
		return err
	}

	return uc.repo.Delete(ctx, pID)
}
