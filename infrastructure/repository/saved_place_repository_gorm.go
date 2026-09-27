package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainsavedplace "go-ride-backend/domain/savedplace"

	schemamodels "github.com/shawon-kanji/go-ride-db-schema/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavedPlaceRepositoryGorm struct {
	db *gorm.DB
}

func NewSavedPlaceRepositoryGorm(db *gorm.DB) *SavedPlaceRepositoryGorm {
	return &SavedPlaceRepositoryGorm{db: db}
}

func toDomainSavedPlace(model *schemamodels.SavedPlace) *domainsavedplace.SavedPlace {
	return &domainsavedplace.SavedPlace{
		ID:               model.ID,
		UserID:           model.UserID,
		Label:            model.Label,
		PlaceType:        model.PlaceType,
		FormattedAddress: model.FormattedAddress,
		Lat:              model.Lat,
		Lng:              model.Lng,
		CreatedAt:        model.CreatedAt,
		UpdatedAt:        model.UpdatedAt,
	}
}

func (r *SavedPlaceRepositoryGorm) Create(ctx context.Context, input *domainsavedplace.CreateInput) (*domainsavedplace.SavedPlace, error) {
	now := time.Now().UTC()
	model := &schemamodels.SavedPlace{
		ID:               uuid.New(),
		UserID:           input.UserID,
		Label:            input.Label,
		PlaceType:        input.PlaceType,
		FormattedAddress: input.FormattedAddress,
		Lat:              input.Lat,
		Lng:              input.Lng,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("create saved place: %w", err)
	}
	return toDomainSavedPlace(model), nil
}

func (r *SavedPlaceRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*domainsavedplace.SavedPlace, error) {
	var model schemamodels.SavedPlace
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainsavedplace.ErrSavedPlaceNotFound
		}
		return nil, fmt.Errorf("find saved place by id: %w", err)
	}
	return toDomainSavedPlace(&model), nil
}

func (r *SavedPlaceRepositoryGorm) GetByUserAndType(ctx context.Context, userID uuid.UUID, placeType string) (*domainsavedplace.SavedPlace, error) {
	var model schemamodels.SavedPlace
	if err := r.db.WithContext(ctx).Where("user_id = ? AND place_type = ?", userID, placeType).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainsavedplace.ErrSavedPlaceNotFound
		}
		return nil, fmt.Errorf("find saved place by user and type: %w", err)
	}
	return toDomainSavedPlace(&model), nil
}

func (r *SavedPlaceRepositoryGorm) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domainsavedplace.SavedPlace, error) {
	var rows []schemamodels.SavedPlace
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list saved places by user: %w", err)
	}
	places := make([]*domainsavedplace.SavedPlace, 0, len(rows))
	for i := range rows {
		places = append(places, toDomainSavedPlace(&rows[i]))
	}
	return places, nil
}

func (r *SavedPlaceRepositoryGorm) Update(ctx context.Context, id uuid.UUID, input *domainsavedplace.UpdateInput) (*domainsavedplace.SavedPlace, error) {
	updates := map[string]interface{}{
		"label":             input.Label,
		"place_type":        input.PlaceType,
		"formatted_address": input.FormattedAddress,
		"lat":               input.Lat,
		"lng":               input.Lng,
		"updated_at":        time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&schemamodels.SavedPlace{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update saved place: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domainsavedplace.ErrSavedPlaceNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *SavedPlaceRepositoryGorm) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&schemamodels.SavedPlace{})
	if result.Error != nil {
		return fmt.Errorf("delete saved place: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainsavedplace.ErrSavedPlaceNotFound
	}
	return nil
}
