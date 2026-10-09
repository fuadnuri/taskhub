package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type RoleResponse struct {
	ID             uuid.UUID            `json:"id"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	CreatedAt      time.Time            `json:"created_at"`
	Permissions    []PermissionResponse `json:"permissions,omitempty"`
}

type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
