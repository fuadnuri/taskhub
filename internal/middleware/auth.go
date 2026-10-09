package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextUserID = "user_id"
	ContextEmail  = "email"
	ContextOrgID  = "organization_id"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidOrgID = errors.New("organization context missing or invalid")
)

// AuthMiddleware creates a Gin middleware that validates JWT access tokens.
func AuthMiddleware(tokenService service.ITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Format: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]
		claims, err := tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Store user id
		if uid, ok := claims["user_id"].(uuid.UUID); ok {
			c.Set(ContextUserID, uid)
		} else if uidStr, ok := claims["user_id"].(string); ok {
			if parsed, err := uuid.Parse(uidStr); err == nil {
				c.Set(ContextUserID, parsed)
			}
		}

		// Store email
		if email, ok := claims["email"].(string); ok {
			c.Set(ContextEmail, email)
		}

		// Store org id if present
		if orgID, ok := claims["organization_id"].(uuid.UUID); ok && orgID != uuid.Nil {
			c.Set(ContextOrgID, orgID)
		} else if orgIDStr, ok := claims["organization_id"].(string); ok && orgIDStr != "" {
			if parsed, err := uuid.Parse(orgIDStr); err == nil && parsed != uuid.Nil {
				c.Set(ContextOrgID, parsed)
			}
		}

		c.Next()
	}
}

// GetUserID extracts the authenticated user ID from Gin context.
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return uuid.Nil, ErrUnauthorized
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUnauthorized
	}
	return id, nil
}

// GetOrgID extracts the organization ID from Gin context.
func GetOrgID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get(ContextOrgID)
	if !exists {
		// Also allow override from header or query param if needed
		if headerOrg := c.GetHeader("X-Organization-ID"); headerOrg != "" {
			if parsed, err := uuid.Parse(headerOrg); err == nil {
				return parsed, nil
			}
		}
		return uuid.Nil, ErrInvalidOrgID
	}
	id, ok := val.(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, ErrInvalidOrgID
	}
	return id, nil
}

// GetUserEmail extracts the authenticated user's email.
func GetUserEmail(c *gin.Context) (string, error) {
	val, exists := c.Get(ContextEmail)
	if !exists {
		return "", ErrUnauthorized
	}
	email, ok := val.(string)
	if !ok {
		return "", ErrUnauthorized
	}
	return email, nil
}