package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	service service.IRoleService
}

func NewRoleHandler(service service.IRoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	orgID, err := middleware.GetOrgID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), orgID, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, role)
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id format"})
		return
	}

	role, err := h.service.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, role)
}

func (h *RoleHandler) GetRolesByOrg(c *gin.Context) {
	orgID, err := middleware.GetOrgID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	roles, err := h.service.GetRolesByOrgID(c.Request.Context(), orgID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, roles)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id format"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, role)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id format"})
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "role deleted successfully"})
}

func (h *RoleHandler) AssignPermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id format"})
		return
	}

	var body struct {
		PermissionID uuid.UUID `json:"permission_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignPermission(c.Request.Context(), roleID, body.PermissionID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "permission assigned successfully"})
}

func (h *RoleHandler) RevokePermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id format"})
		return
	}

	permID, err := uuid.Parse(c.Param("perm_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id format"})
		return
	}

	if err := h.service.RevokePermission(c.Request.Context(), roleID, permID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "permission revoked successfully"})
}

func (h *RoleHandler) CreatePermission(c *gin.Context) {
	var req dto.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	perm, err := h.service.CreatePermission(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, perm)
}

func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	perms, err := h.service.GetAllPermissions(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, perms)
}
