package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware gracefully catches panics and returns a 500 error response.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in HTTP handler", "error", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error occurred",
				})
			}
		}()
		c.Next()
	}
}
