package user

import (
	"strings"
	"time"

	domainuser "go-ride-backend/domain/user"

	"github.com/google/uuid"
)

func NewUserFromSignup(req SignupRequest, hashedPassword string) *domainuser.User {
	now := time.Now().UTC()

	return &domainuser.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hashedPassword,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func ToUserResponse(u *domainuser.User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
