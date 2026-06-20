package repository

import (
	"context"
	"errors"
	"fmt"

	domainuser "go-ride-backend/domain/user"
	"go-ride-backend/infrastructure/db/models"

	"gorm.io/gorm"
)

type UserRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) Create(ctx context.Context, user *domainuser.User) error {
	model := models.FromDomain(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepositoryGorm) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainuser.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return model.ToDomain(), nil
}
