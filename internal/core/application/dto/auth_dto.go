package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	OrgName   string `json:"org_name,omitempty"` // optional initial workspace
}

type LoginRequest struct {
	Email          string     `json:"email" binding:"required,email"`
	Password       string     `json:"password" binding:"required"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	User         UserResponse `json:"user"`
}

type TokenClaimsResponse struct {
	UserID         uuid.UUID  `json:"user_id"`
	Email          string     `json:"email"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	ExpiresAt      time.Time  `json:"expires_at"`
}
