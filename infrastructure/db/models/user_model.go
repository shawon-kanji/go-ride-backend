package models

import (
	"time"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

type UserModel struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string    `gorm:"type:varchar(255);not null"`
	FirstName     string    `gorm:"type:varchar(100);not null"`
	LastName      string    `gorm:"type:varchar(100);not null"`
	AccountStatus string    `gorm:"type:varchar(50);not null;default:active"`
	DeactivatedAt *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func FromDomain(u *domainuser.User) *UserModel {
	return &UserModel{
		ID:            u.ID,
		Email:         u.Email,
		PasswordHash:  u.PasswordHash,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AccountStatus: u.AccountStatus,
		DeactivatedAt: u.DeactivatedAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func (m *UserModel) ToDomain() *domainuser.User {
	return &domainuser.User{
		ID:            m.ID,
		Email:         m.Email,
		PasswordHash:  m.PasswordHash,
		FirstName:     m.FirstName,
		LastName:      m.LastName,
		AccountStatus: m.AccountStatus,
		DeactivatedAt: m.DeactivatedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
