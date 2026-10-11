package main

import (
	"log/slog"
	"os"

	"github.com/fuadnuri/taskhub/internal/config"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	infraauth "github.com/fuadnuri/taskhub/internal/infrastructure/auth"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	notificationrepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/notification"
	orgrepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/organization"
	projectrepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/project"
	rolerepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/role"
	taskrepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/task"
	userrepo "github.com/fuadnuri/taskhub/internal/infrastructure/repository/user"
	httpHandler "github.com/fuadnuri/taskhub/internal/interface/httphandler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	// 1. Connect to PostgreSQL
	db, err := database.NewPostgres(cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// 2. Auto-migrate schema tables
	err = db.AutoMigrate(
		&schema.User{},
		&schema.Organization{},
		&schema.OrganizationMember{},
		&schema.OrganizationInvitation{},
		&schema.Role{},
		&schema.Permission{},
		&schema.RolePermission{},
		&schema.Project{},
		&schema.ProjectMember{},
		&schema.Task{},
		&schema.TaskAssignee{},
		&schema.TaskComment{},
		&schema.Notification{},
		&schema.RefreshToken{},
		&schema.AuditLog{},
	)
	if err != nil {
		slog.Error("Auto-migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migration completed successfully")

	// 3. Infrastructure: Auth & Crypto
	jwtService := infraauth.NewJWTService(infraauth.JWTConfig{
		SecretKey:            cfg.JWTSecret,
		AccessTokenDuration:  cfg.AccessTokenDuration,
		RefreshTokenDuration: cfg.RefreshTokenDuration,
	})
	passwordHasher := infraauth.NewPasswordHasher(0)

	// 4. Infrastructure: Repositories
	userRepository := userrepo.NewUserRepository(db)
	refreshTokenRepository := userrepo.NewRefreshTokenRepository(db)
	organizationRepository := orgrepo.NewOrganizationRepository(db)
	orgMemberRepository := orgrepo.NewOrgMemberRepository(db)
	orgInvitationRepository := orgrepo.NewOrgInvitationRepository(db)
	projectRepository := projectrepo.NewProjectRepository(db)
	projectMemberRepository := projectrepo.NewProjectMemberRepository(db)
	taskRepository := taskrepo.NewTaskRepository(db)
	taskAssigneeRepository := taskrepo.NewTaskAssigneeRepository(db)
	taskCommentRepository := taskrepo.NewTaskCommentRepository(db)
	roleRepository := rolerepo.NewRoleRepository(db)
	permissionRepository := rolerepo.NewPermissionRepository(db)
	rolePermissionRepository := rolerepo.NewRolePermissionRepository(db)
	notificationRepository := notificationrepo.NewNotificationRepository(db)

	// 5. Application Services
	authService := service.NewAuthService(
		userRepository,
		refreshTokenRepository,
		organizationRepository,
		orgMemberRepository,
		jwtService,
		passwordHasher,
	)
	userService := service.NewUserService(userRepository, passwordHasher)
	organizationService := service.NewOrganizationService(
		organizationRepository,
		orgMemberRepository,
		orgInvitationRepository,
		jwtService,
	)
	projectService := service.NewProjectService(projectRepository, projectMemberRepository)
	taskService := service.NewTaskService(
		taskRepository,
		taskAssigneeRepository,
		taskCommentRepository,
		notificationRepository,
	)
	roleService := service.NewRoleService(
		roleRepository,
		permissionRepository,
		rolePermissionRepository,
	)
	notificationService := service.NewNotificationService(notificationRepository)

	// 6. HTTP Handlers
	authHandler := httpHandler.NewAuthHandler(authService, userService)
	userHandler := httpHandler.NewUserHandler(userService)
	orgHandler := httpHandler.NewOrganizationHandler(organizationService)
	projectHandler := httpHandler.NewProjectHandler(projectService)
	taskHandler := httpHandler.NewTaskHandler(taskService)
	roleHandler := httpHandler.NewRoleHandler(roleService)
	notificationHandler := httpHandler.NewNotificationHandler(notificationService)

	// 7. Router
	router := httpHandler.SetupRouter(httpHandler.RouterConfig{
		TokenService:        jwtService,
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		OrganizationHandler: orgHandler,
		ProjectHandler:      projectHandler,
		TaskHandler:         taskHandler,
		RoleHandler:         roleHandler,
		NotificationHandler: notificationHandler,
	})

	// 8. Start HTTP Server
	addr := ":" + cfg.Port
	slog.Info("Starting HTTP server", "address", addr)
	if err := router.Run(addr); err != nil {
		slog.Error("Server terminated unexpectedly", "error", err)
	}
}
