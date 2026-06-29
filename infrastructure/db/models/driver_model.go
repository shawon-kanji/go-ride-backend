package models

import (
	"time"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

type DriverModel struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email           string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash    string    `gorm:"type:varchar(255);not null"`
	FirstName       string    `gorm:"type:varchar(100);not null"`
	LastName        string    `gorm:"type:varchar(100);not null"`
	AccountStatus   string    `gorm:"type:varchar(50);not null;default:pending"`
	IsEmailVerified bool      `gorm:"not null;default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (DriverModel) TableName() string {
	return "drivers"
}

func FromDriverDomain(d *domaindriver.Driver) *DriverModel {
	return &DriverModel{
		ID:              d.ID,
		Email:           d.Email,
		PasswordHash:    d.PasswordHash,
		FirstName:       d.FirstName,
		LastName:        d.LastName,
		AccountStatus:   d.AccountStatus,
		IsEmailVerified: d.IsEmailVerified,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func (m *DriverModel) ToDomain() *domaindriver.Driver {
	return &domaindriver.Driver{
		ID:              m.ID,
		Email:           m.Email,
		PasswordHash:    m.PasswordHash,
		FirstName:       m.FirstName,
		LastName:        m.LastName,
		AccountStatus:   m.AccountStatus,
		IsEmailVerified: m.IsEmailVerified,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
