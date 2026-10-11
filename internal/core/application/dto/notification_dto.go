package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateNotificationRequest struct {
	OrganizationID uuid.UUID      `json:"organization_id" binding:"required"`
	UserID         uuid.UUID      `json:"user_id" binding:"required"`
	Type           string         `json:"type" binding:"required"`
	Title          string         `json:"title" binding:"required"`
	Message        string         `json:"message" binding:"required"`
	Data           map[string]any `json:"data,omitempty"`
}

type NotificationResponse struct {
	ID             uuid.UUID      `json:"id"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	UserID         uuid.UUID      `json:"user_id"`
	Type           string         `json:"type"`
	Title          string         `json:"title"`
	Message        string         `json:"message"`
	Data           map[string]any `json:"data,omitempty"`
	ReadAt         *time.Time     `json:"read_at"`
	CreatedAt      time.Time      `json:"created_at"`
}
