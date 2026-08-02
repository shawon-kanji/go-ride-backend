package driver

import (
	"strings"
	"time"

	domaindriver "go-ride-backend/domain/driver"

	"github.com/google/uuid"
)

func NewDriverFromSignup(req SignupRequest, hashedPassword string) *domaindriver.Driver {
	now := time.Now().UTC()

	return &domaindriver.Driver{
		ID:              uuid.New(),
		Email:           strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash:    hashedPassword,
		FirstName:       strings.TrimSpace(req.FirstName),
		LastName:        strings.TrimSpace(req.LastName),
		AccountStatus:   domaindriver.AccountStatusPending,
		IsEmailVerified: false,
		IsOnline:        false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func ToDriverResponse(d *domaindriver.Driver) DriverResponse {
	return DriverResponse{
		ID:              d.ID.String(),
		Email:           d.Email,
		FirstName:       d.FirstName,
		LastName:        d.LastName,
		AccountStatus:   d.AccountStatus,
		IsEmailVerified: d.IsEmailVerified,
		IsOnline:        d.IsOnline,
	}
}
