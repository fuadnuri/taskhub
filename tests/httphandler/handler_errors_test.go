package httphandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	httpHandler "github.com/fuadnuri/taskhub/internal/interface/httphandler"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRespondError_StatusMapping(t *testing.T) {
	tests := []struct {
		name           string
		appError       error
		expectedStatus int
	}{
		{
			name:           "ErrNotFound maps to 404",
			appError:       apperrors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ErrAlreadyExists maps to 409",
			appError:       apperrors.ErrAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ErrInvalidCredentials maps to 401",
			appError:       apperrors.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ErrUnauthorized maps to 401",
			appError:       apperrors.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ErrForbidden maps to 403",
			appError:       apperrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrBadRequest maps to 400",
			appError:       apperrors.ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrTokenExpired maps to 401",
			appError:       apperrors.ErrTokenExpired,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ErrInvalidToken maps to 401",
			appError:       apperrors.ErrInvalidToken,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "ErrInvitationExpired maps to 400",
			appError:       apperrors.ErrInvitationExpired,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrInvitationAccepted maps to 400",
			appError:       apperrors.ErrInvitationAccepted,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "ErrUserInactive maps to 403",
			appError:       apperrors.ErrUserInactive,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrOrganizationInactive maps to 403",
			appError:       apperrors.ErrOrganizationInactive,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "ErrInternal maps to 500",
			appError:       apperrors.ErrInternal,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Unknown generic error maps to 500",
			appError:       errors.New("something totally unexpected"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			httpHandler.RespondError(c, tt.appError)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected HTTP status %d, got %d", tt.expectedStatus, w.Code)
			}

			var body map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to decode response JSON: %v", err)
			}

			if success, ok := body["success"].(bool); !ok || success {
				t.Fatalf("expected success to be false, got %v", body["success"])
			}

			if errMsg, ok := body["error"].(string); !ok || errMsg != tt.appError.Error() {
				t.Fatalf("expected error message %q, got %v", tt.appError.Error(), body["error"])
			}
		})
	}
}

func TestHandlerInputValidationErrors(t *testing.T) {
	router := gin.New()

	router.POST("/test/register", func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	router.POST("/test/login", func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})

	t.Run("Register_InvalidEmail_Returns400", func(t *testing.T) {
		invalidPayload := []byte(`{"email":"not-an-email","password":"password123","first_name":"John","last_name":"Doe"}`)
		req := httptest.NewRequest(http.MethodPost, "/test/register", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for invalid email, got %d", w.Code)
		}
	})

	t.Run("Register_PasswordTooShort_Returns400", func(t *testing.T) {
		invalidPayload := []byte(`{"email":"test@example.com","password":"short","first_name":"John","last_name":"Doe"}`)
		req := httptest.NewRequest(http.MethodPost, "/test/register", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for short password, got %d", w.Code)
		}
	})

	t.Run("Login_MissingEmail_Returns400", func(t *testing.T) {
		invalidPayload := []byte(`{"password":"password123"}`)
		req := httptest.NewRequest(http.MethodPost, "/test/login", bytes.NewBuffer(invalidPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for missing email, got %d", w.Code)
		}
	})

	t.Run("MalformedJSON_Returns400", func(t *testing.T) {
		malformedJSON := []byte(`{"email": invalid-json}`)
		req := httptest.NewRequest(http.MethodPost, "/test/register", bytes.NewBuffer(malformedJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for malformed json, got %d", w.Code)
		}
	})
}
