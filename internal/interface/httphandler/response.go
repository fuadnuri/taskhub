package httpHandler

import (
	"errors"
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/gin-gonic/gin"
)

// RespondJSON sends a JSON response with status code and payload.
func RespondJSON(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

// RespondCreated sends a 201 Created JSON response.
func RespondCreated(c *gin.Context, data any) {
	RespondJSON(c, http.StatusCreated, data)
}

// RespondOK sends a 200 OK JSON response.
func RespondOK(c *gin.Context, data any) {
	RespondJSON(c, http.StatusOK, data)
}

// RespondError maps application errors to appropriate HTTP status codes.
func RespondError(c *gin.Context, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, apperrors.ErrAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		status = http.StatusUnauthorized
	case errors.Is(err, apperrors.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, apperrors.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, apperrors.ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, apperrors.ErrTokenExpired), errors.Is(err, apperrors.ErrInvalidToken):
		status = http.StatusUnauthorized
	case errors.Is(err, apperrors.ErrInvitationExpired), errors.Is(err, apperrors.ErrInvitationAccepted):
		status = http.StatusBadRequest
	case errors.Is(err, apperrors.ErrUserInactive), errors.Is(err, apperrors.ErrOrganizationInactive):
		status = http.StatusForbidden
	}

	c.JSON(status, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
