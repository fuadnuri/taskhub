package httpHandler

import (
	"net/http"

	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.IAuthService
	userService service.IUserService
}

func NewAuthHandler(authService service.IAuthService, userService service.IUserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondCreated(c, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, res)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		RespondError(c, err)
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, err)
		return
	}

	RespondOK(c, user)
}
