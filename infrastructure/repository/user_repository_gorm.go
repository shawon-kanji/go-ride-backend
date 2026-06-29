package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domainuser "go-ride-backend/domain/user"
	"go-ride-backend/infrastructure/db/models"

	"github.com/google/uuid"
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

func (r *UserRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainuser.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return model.ToDomain(), nil
}

func (r *UserRepositoryGorm) UpdateProfile(ctx context.Context, id uuid.UUID, firstName string, lastName string) (*domainuser.User, error) {
	updates := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"updated_at": time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("update user profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, domainuser.ErrUserNotFound
	}

	return r.GetByID(ctx, id)
}

func (r *UserRepositoryGorm) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	updates := map[string]interface{}{
		"password_hash": passwordHash,
		"updated_at":    time.Now().UTC(),
	}

	result := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("update user password: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainuser.ErrUserNotFound
	}
	return nil
}

func (r *UserRepositoryGorm) Deactivate(ctx context.Context, id uuid.UUID, at time.Time) error {
	updates := map[string]interface{}{
		"account_status": domainuser.AccountStatusDeactivated,
		"deactivated_at": at,
		"updated_at":     at,
	}

	result := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("deactivate user account: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domainuser.ErrUserNotFound
	}
	return nil
}
