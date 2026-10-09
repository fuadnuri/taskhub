package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service service.INotificationService
}

func NewNotificationHandler(service service.INotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	notifications, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, notifications)
}

func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	notifications, err := h.service.GetUnreadByUserID(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id format"})
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "notification marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "notification deleted"})
}
