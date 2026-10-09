package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RouterConfig bundles all handlers required to configure the HTTP routes.
type RouterConfig struct {
	TokenService        service.ITokenService
	AuthHandler         *AuthHandler
	UserHandler         *UserHandler
	OrganizationHandler *OrganizationHandler
	ProjectHandler      *ProjectHandler
	TaskHandler         *TaskHandler
	RoleHandler         *RoleHandler
	NotificationHandler *NotificationHandler
}

// SetupRouter sets up and registers all middlewares and routes on a new Gin engine.
func SetupRouter(cfg RouterConfig) *gin.Engine {
	router := gin.New()

	// Global middlewares
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	api := router.Group("/api/v1")

	// ── Public Auth Routes ──
	if cfg.AuthHandler != nil {
		auth := api.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.POST("/refresh", cfg.AuthHandler.RefreshToken)
			auth.POST("/logout", cfg.AuthHandler.Logout)
		}
	}

	// ── Protected Routes ──
	var protected *gin.RouterGroup
	if cfg.TokenService != nil {
		protected = api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.TokenService))
	} else {
		protected = api.Group("")
	}

	// Auth (Protected)
	if cfg.AuthHandler != nil {
		protected.GET("/auth/me", cfg.AuthHandler.Me)
	}

	// Users
	if cfg.UserHandler != nil {
		users := protected.Group("/users")
		{
			users.POST("", cfg.UserHandler.Create)
			users.GET("", cfg.UserHandler.GetAllUsers)
			users.GET("/:id", cfg.UserHandler.GetByID)
			users.PUT("/:id", cfg.UserHandler.Update)
			users.PUT("/:id/password", cfg.UserHandler.ChangePassword)
			users.DELETE("/:id", cfg.UserHandler.Delete)
		}
	}

	// Organizations
	if cfg.OrganizationHandler != nil {
		orgs := protected.Group("/organizations")
		{
			orgs.POST("", cfg.OrganizationHandler.Create)
			orgs.GET("", cfg.OrganizationHandler.GetAll)
			orgs.GET("/:id", cfg.OrganizationHandler.GetByID)
			orgs.GET("/slug/:slug", cfg.OrganizationHandler.GetBySlug)
			orgs.PUT("/:id", cfg.OrganizationHandler.Update)
			orgs.DELETE("/:id", cfg.OrganizationHandler.Delete)
			orgs.POST("/:id/invitations", cfg.OrganizationHandler.InviteMember)
			orgs.POST("/invitations/accept", cfg.OrganizationHandler.AcceptInvitation)
			orgs.GET("/:id/members", cfg.OrganizationHandler.GetMembers)
			orgs.DELETE("/:id/members/:member_id", cfg.OrganizationHandler.RemoveMember)
			orgs.PUT("/:id/members/:member_id/role", cfg.OrganizationHandler.UpdateMemberRole)
		}
	}

	// Projects
	if cfg.ProjectHandler != nil {
		projects := protected.Group("/projects")
		{
			projects.POST("", cfg.ProjectHandler.Create)
			projects.GET("", cfg.ProjectHandler.GetByOrg)
			projects.GET("/:id", cfg.ProjectHandler.GetByID)
			projects.PUT("/:id", cfg.ProjectHandler.Update)
			projects.DELETE("/:id", cfg.ProjectHandler.Delete)
			projects.POST("/:id/members", cfg.ProjectHandler.AddMember)
			projects.DELETE("/:id/members/:user_id", cfg.ProjectHandler.RemoveMember)
			projects.GET("/:id/members", cfg.ProjectHandler.GetMembers)
		}
	}

	// Tasks
	if cfg.TaskHandler != nil {
		tasks := protected.Group("")
		{
			tasks.POST("/projects/:project_id/tasks", cfg.TaskHandler.Create)
			tasks.GET("/projects/:project_id/tasks", cfg.TaskHandler.GetByProject)
			tasks.GET("/tasks", cfg.TaskHandler.GetByOrg)
			tasks.GET("/tasks/:id", cfg.TaskHandler.GetByID)
			tasks.PUT("/tasks/:id", cfg.TaskHandler.Update)
			tasks.DELETE("/tasks/:id", cfg.TaskHandler.Delete)
			tasks.POST("/tasks/:id/assign", cfg.TaskHandler.AssignUser)
			tasks.DELETE("/tasks/:id/assign/:user_id", cfg.TaskHandler.UnassignUser)
			tasks.POST("/tasks/:id/comments", cfg.TaskHandler.AddComment)
			tasks.GET("/tasks/:id/comments", cfg.TaskHandler.GetComments)
			tasks.DELETE("/tasks/comments/:comment_id", cfg.TaskHandler.DeleteComment)
		}
	}

	// Roles & Permissions
	if cfg.RoleHandler != nil {
		roles := protected.Group("/roles")
		{
			roles.POST("", cfg.RoleHandler.CreateRole)
			roles.GET("", cfg.RoleHandler.GetRolesByOrg)
			roles.GET("/:id", cfg.RoleHandler.GetRole)
			roles.PUT("/:id", cfg.RoleHandler.UpdateRole)
			roles.DELETE("/:id", cfg.RoleHandler.DeleteRole)
			roles.POST("/:id/permissions", cfg.RoleHandler.AssignPermission)
			roles.DELETE("/:id/permissions/:perm_id", cfg.RoleHandler.RevokePermission)
		}

		perms := protected.Group("/permissions")
		{
			perms.POST("", cfg.RoleHandler.CreatePermission)
			perms.GET("", cfg.RoleHandler.GetAllPermissions)
		}
	}

	// Notifications
	if cfg.NotificationHandler != nil {
		notifs := protected.Group("/notifications")
		{
			notifs.GET("", cfg.NotificationHandler.GetUserNotifications)
			notifs.GET("/unread", cfg.NotificationHandler.GetUnreadNotifications)
			notifs.PUT("/:id/read", cfg.NotificationHandler.MarkAsRead)
			notifs.DELETE("/:id", cfg.NotificationHandler.DeleteNotification)
		}
	}

	return router
}
