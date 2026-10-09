package service

import (
	"context"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/google/uuid"
)

// ITokenService defines the token management contract.
type ITokenService interface {
	GenerateAccessToken(userID uuid.UUID, email string, orgID uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (map[string]any, error)
	GenerateRefreshToken() (raw string, hash string, expiresAt time.Time, err error)
	HashRefreshToken(raw string) string
}

// IPasswordService defines the password hashing contract.
type IPasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

// IAuthService defines authentication and authorization use cases.
type IAuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

// IUserService defines user management use cases.
type IUserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetAll(ctx context.Context) ([]dto.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IOrganizationService defines organization management use cases.
type IOrganizationService interface {
	Create(ctx context.Context, ownerID uuid.UUID, req dto.CreateOrgRequest) (*dto.OrgResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.OrgResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.OrgResponse, error)
	GetAll(ctx context.Context) ([]dto.OrgResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateOrgRequest) (*dto.OrgResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	InviteMember(ctx context.Context, orgID uuid.UUID, req dto.InviteMemberRequest) (*dto.OrgInvitationResponse, error)
	AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error
	GetMembers(ctx context.Context, orgID uuid.UUID) ([]dto.OrgMemberResponse, error)
	RemoveMember(ctx context.Context, memberID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, memberID uuid.UUID, req dto.UpdateMemberRoleRequest) error
}

// IProjectService defines project management use cases.
type IProjectService interface {
	Create(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID, req dto.CreateProjectRequest) (*dto.ProjectResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ProjectResponse, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]dto.ProjectResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	GetMembers(ctx context.Context, projectID uuid.UUID) ([]dto.ProjectMemberResponse, error)
}

// ITaskService defines task and comment management use cases.
type ITaskService interface {
	Create(ctx context.Context, orgID uuid.UUID, projectID uuid.UUID, createdBy uuid.UUID, req dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.TaskResponse, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]dto.TaskResponse, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]dto.TaskResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AssignUser(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error
	UnassignUser(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error
	GetAssigneeIDs(ctx context.Context, taskID uuid.UUID) ([]uuid.UUID, error)
	AddComment(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req dto.CreateTaskCommentRequest) (*dto.TaskCommentResponse, error)
	GetComments(ctx context.Context, taskID uuid.UUID) ([]dto.TaskCommentResponse, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}

// IRoleService defines RBAC role and permission use cases.
type IRoleService interface {
	CreateRole(ctx context.Context, orgID uuid.UUID, req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error)
	GetRolesByOrgID(ctx context.Context, orgID uuid.UUID) ([]dto.RoleResponse, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
	AssignPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
	RevokePermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error
	CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetAllPermissions(ctx context.Context) ([]dto.PermissionResponse, error)
}

// INotificationService defines user notification use cases.
type INotificationService interface {
	Create(ctx context.Context, req dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.NotificationResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]dto.NotificationResponse, error)
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]dto.NotificationResponse, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
