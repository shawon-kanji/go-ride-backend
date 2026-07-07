package repository

import (
	"context"
	"errors"
	"fmt"

	domaindriver "go-ride-backend/domain/driver"

	schemamodels "github.com/shawon-kanji/go-ride-db-schema/models"

	"gorm.io/gorm"
)

type DriverRepositoryGorm struct {
	db *gorm.DB
}

func NewDriverRepositoryGorm(db *gorm.DB) *DriverRepositoryGorm {
	return &DriverRepositoryGorm{db: db}
}

func (r *DriverRepositoryGorm) Create(ctx context.Context, driver *domaindriver.Driver) error {
	model := &schemamodels.Driver{
		ID:              driver.ID,
		Email:           driver.Email,
		PasswordHash:    driver.PasswordHash,
		FirstName:       driver.FirstName,
		LastName:        driver.LastName,
		AccountStatus:   driver.AccountStatus,
		IsEmailVerified: driver.IsEmailVerified,
		CreatedAt:       driver.CreatedAt,
		UpdatedAt:       driver.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create driver: %w", err)
	}
	return nil
}

func (r *DriverRepositoryGorm) GetByEmail(ctx context.Context, email string) (*domaindriver.Driver, error) {
	var model schemamodels.Driver
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindriver.ErrDriverNotFound
		}
		return nil, fmt.Errorf("find driver by email: %w", err)
	}
	return &domaindriver.Driver{
		ID:              model.ID,
		Email:           model.Email,
		PasswordHash:    model.PasswordHash,
		FirstName:       model.FirstName,
		LastName:        model.LastName,
		AccountStatus:   model.AccountStatus,
		IsEmailVerified: model.IsEmailVerified,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}, nil
}
