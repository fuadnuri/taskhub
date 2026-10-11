package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service service.IProjectService
}

func NewProjectHandler(service service.IProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	orgID, err := middleware.GetOrgID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.Create(c.Request.Context(), orgID, userID, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, project)
}

func (h *ProjectHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	project, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, project)
}

func (h *ProjectHandler) GetByOrg(c *gin.Context) {
	orgID, err := middleware.GetOrgID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	projects, err := h.service.GetByOrganizationID(c.Request.Context(), orgID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, projects)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "project deleted successfully"})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	var body struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddMember(c.Request.Context(), projectID, body.UserID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "member added to project"})
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), projectID, userID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "member removed from project"})
}

func (h *ProjectHandler) GetMembers(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	members, err := h.service.GetMembers(c.Request.Context(), projectID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, members)
}
