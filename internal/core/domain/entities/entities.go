package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserEntity represents the core domain model for a user.
type UserEntity struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// OrganizationEntity represents a tenant/workspace.
type OrganizationEntity struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrganizationMemberEntity represents a user's membership in an organization.
type OrganizationMemberEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	RoleID         uuid.UUID
	JoinedAt       time.Time
}

// OrganizationInvitationEntity represents an invitation sent to join an organization.
type OrganizationInvitationEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	RoleID         uuid.UUID
	TokenHash      string
	ExpiresAt      time.Time
	AcceptedAt     *time.Time
	CreatedAt      time.Time
}

// RoleEntity represents an RBAC role scoped to an organization.
type RoleEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    string
	CreatedAt      time.Time
}

// PermissionEntity represents a named capability that can be assigned to roles.
type PermissionEntity struct {
	ID          uuid.UUID
	Name        string
	Description string
}

// RolePermissionEntity is the join between a role and a permission.
type RolePermissionEntity struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

// ProjectEntity represents a project within an organization.
type ProjectEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	CreatedBy      uuid.UUID
	Name           string
	Description    string
	Status         string
	StartDate      *time.Time
	DueDate        *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ProjectMemberEntity is the join between a project and a user.
type ProjectMemberEntity struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	JoinedAt  time.Time
}

// TaskEntity represents a task within a project.
type TaskEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ProjectID      uuid.UUID
	CreatedBy      uuid.UUID
	Title          string
	Description    string
	Status         string
	Priority       string
	DueDate        *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TaskAssigneeEntity is the join between a task and an assigned user.
type TaskAssigneeEntity struct {
	TaskID     uuid.UUID
	UserID     uuid.UUID
	AssignedAt time.Time
}

// TaskCommentEntity represents a comment left on a task.
type TaskCommentEntity struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NotificationEntity represents a notification sent to a user.
type NotificationEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Type           string
	Title          string
	Message        string
	Data           map[string]any
	ReadAt         *time.Time
	CreatedAt      time.Time
}

// RefreshTokenEntity represents an issued JWT refresh token.
type RefreshTokenEntity struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// AuditLogEntity represents an immutable record of a user action.
type AuditLogEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Action         string
	ResourceType   string
	ResourceID     uuid.UUID
	Metadata       map[string]any
	IPAddress      string
	CreatedAt      time.Time
}
