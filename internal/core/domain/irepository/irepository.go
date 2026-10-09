package irepository

import (
	"context"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/google/uuid"
)

// IUserRepository defines persistence operations for users.
type IUserRepository interface {
	Create(ctx context.Context, user *entities.UserEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error)
	GetByEmail(ctx context.Context, email string) (*entities.UserEntity, error)
	GetAll(ctx context.Context) ([]entities.UserEntity, error)
	Update(ctx context.Context, user *entities.UserEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IOrganizationRepository defines persistence operations for organizations.
type IOrganizationRepository interface {
	Create(ctx context.Context, org *entities.OrganizationEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationEntity, error)
	GetBySlug(ctx context.Context, slug string) (*entities.OrganizationEntity, error)
	GetAll(ctx context.Context) ([]entities.OrganizationEntity, error)
	Update(ctx context.Context, org *entities.OrganizationEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IOrganizationMemberRepository defines persistence operations for organization memberships.
type IOrganizationMemberRepository interface {
	Add(ctx context.Context, member *entities.OrganizationMemberEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationMemberEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.OrganizationMemberEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.OrganizationMemberEntity, error)
	UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) error
	Remove(ctx context.Context, id uuid.UUID) error
}

// IOrganizationInvitationRepository defines persistence operations for invitations.
type IOrganizationInvitationRepository interface {
	Create(ctx context.Context, invitation *entities.OrganizationInvitationEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationInvitationEntity, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entities.OrganizationInvitationEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.OrganizationInvitationEntity, error)
	Accept(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IRoleRepository defines persistence operations for RBAC roles.
type IRoleRepository interface {
	Create(ctx context.Context, role *entities.RoleEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.RoleEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.RoleEntity, error)
	Update(ctx context.Context, role *entities.RoleEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IPermissionRepository defines persistence operations for permissions.
type IPermissionRepository interface {
	Create(ctx context.Context, permission *entities.PermissionEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PermissionEntity, error)
	GetByName(ctx context.Context, name string) (*entities.PermissionEntity, error)
	GetAll(ctx context.Context) ([]entities.PermissionEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// IRolePermissionRepository defines persistence operations for role-permission assignments.
type IRolePermissionRepository interface {
	Assign(ctx context.Context, rp *entities.RolePermissionEntity) error
	Revoke(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entities.RolePermissionEntity, error)
	GetByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]entities.RolePermissionEntity, error)
}

// IProjectRepository defines persistence operations for projects.
type IProjectRepository interface {
	Create(ctx context.Context, project *entities.ProjectEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.ProjectEntity, error)
	Update(ctx context.Context, project *entities.ProjectEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IProjectMemberRepository defines persistence operations for project memberships.
type IProjectMemberRepository interface {
	Add(ctx context.Context, member *entities.ProjectMemberEntity) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]entities.ProjectMemberEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.ProjectMemberEntity, error)
	Remove(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
}

// ITaskRepository defines persistence operations for tasks.
type ITaskRepository interface {
	Create(ctx context.Context, task *entities.TaskEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TaskEntity, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]entities.TaskEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.TaskEntity, error)
	Update(ctx context.Context, task *entities.TaskEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ITaskAssigneeRepository defines persistence operations for task assignees.
type ITaskAssigneeRepository interface {
	Assign(ctx context.Context, assignee *entities.TaskAssigneeEntity) error
	Unassign(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error
	GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskAssigneeEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.TaskAssigneeEntity, error)
}

// ITaskCommentRepository defines persistence operations for task comments.
type ITaskCommentRepository interface {
	Create(ctx context.Context, comment *entities.TaskCommentEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TaskCommentEntity, error)
	GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskCommentEntity, error)
	Update(ctx context.Context, comment *entities.TaskCommentEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// INotificationRepository defines persistence operations for notifications.
type INotificationRepository interface {
	Create(ctx context.Context, notification *entities.NotificationEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.NotificationEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.NotificationEntity, error)
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]entities.NotificationEntity, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IRefreshTokenRepository defines persistence operations for refresh tokens.
type IRefreshTokenRepository interface {
	Create(ctx context.Context, token *entities.RefreshTokenEntity) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshTokenEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.RefreshTokenEntity, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// IAuditLogRepository defines persistence operations for audit logs.
type IAuditLogRepository interface {
	Create(ctx context.Context, log *entities.AuditLogEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLogEntity, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.AuditLogEntity, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.AuditLogEntity, error)
	GetByResourceID(ctx context.Context, resourceID uuid.UUID) ([]entities.AuditLogEntity, error)
}
