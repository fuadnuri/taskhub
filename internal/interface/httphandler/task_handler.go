package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service service.ITaskService
}

func NewTaskHandler(service service.ITaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
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

	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Create(c.Request.Context(), orgID, projectID, userID, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, task)
}

func (h *TaskHandler) GetByProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("project_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id format"})
		return
	}

	tasks, err := h.service.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, tasks)
}

func (h *TaskHandler) GetByOrg(c *gin.Context) {
	orgID, err := middleware.GetOrgID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	tasks, err := h.service.GetByOrganizationID(c.Request.Context(), orgID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "task deleted successfully"})
}

func (h *TaskHandler) AssignUser(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	var body struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignUser(c.Request.Context(), taskID, body.UserID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "user assigned to task"})
}

func (h *TaskHandler) UnassignUser(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	if err := h.service.UnassignUser(c.Request.Context(), taskID, userID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "user unassigned from task"})
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	var req dto.CreateTaskCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.AddComment(c.Request.Context(), taskID, userID, req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, comment)
}

func (h *TaskHandler) GetComments(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id format"})
		return
	}

	comments, err := h.service.GetComments(c.Request.Context(), taskID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, comments)
}

func (h *TaskHandler) DeleteComment(c *gin.Context) {
	commentID, err := uuid.Parse(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id format"})
		return
	}

	if err := h.service.DeleteComment(c.Request.Context(), commentID); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "comment deleted successfully"})
}
