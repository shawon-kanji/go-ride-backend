package repository

import (
	"context"
	"errors"
	"fmt"

	domaindriver "go-ride-backend/domain/driver"
	"go-ride-backend/infrastructure/db/models"

	"gorm.io/gorm"
)

type DriverRepositoryGorm struct {
	db *gorm.DB
}

func NewDriverRepositoryGorm(db *gorm.DB) *DriverRepositoryGorm {
	return &DriverRepositoryGorm{db: db}
}

func (r *DriverRepositoryGorm) Create(ctx context.Context, driver *domaindriver.Driver) error {
	model := models.FromDriverDomain(driver)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create driver: %w", err)
	}
	return nil
}

func (r *DriverRepositoryGorm) GetByEmail(ctx context.Context, email string) (*domaindriver.Driver, error) {
	var model models.DriverModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domaindriver.ErrDriverNotFound
		}
		return nil, fmt.Errorf("find driver by email: %w", err)
	}
	return model.ToDomain(), nil
}
