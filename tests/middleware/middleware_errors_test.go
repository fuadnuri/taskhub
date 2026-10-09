package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fuadnuri/taskhub/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockTokenService implements service.ITokenService for middleware tests
type mockTokenService struct {
	validateFn func(token string) (map[string]any, error)
}

func (m *mockTokenService) GenerateAccessToken(userID uuid.UUID, email string, orgID uuid.UUID) (string, error) {
	return "mock-token", nil
}

func (m *mockTokenService) ValidateAccessToken(tokenString string) (map[string]any, error) {
	if m.validateFn != nil {
		return m.validateFn(tokenString)
	}
	return nil, errors.New("invalid token")
}

func (m *mockTokenService) GenerateRefreshToken() (string, string, time.Time, error) {
	return "raw", "hash", time.Now().Add(time.Hour), nil
}

func (m *mockTokenService) HashRefreshToken(raw string) string {
	return "hashed-" + raw
}

func TestAuthMiddleware_ErrorScenarios(t *testing.T) {
	t.Run("MissingAuthorizationHeader_Returns401", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AuthMiddleware(&mockTokenService{}))
		router.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("MalformedHeader_NoBearerPrefix_Returns401", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.AuthMiddleware(&mockTokenService{}))
		router.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("InvalidOrExpiredToken_Returns401", func(t *testing.T) {
		mockService := &mockTokenService{
			validateFn: func(token string) (map[string]any, error) {
				return nil, errors.New("token expired")
			},
		}

		router := gin.New()
		router.Use(middleware.AuthMiddleware(mockService))
		router.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-or-expired-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("ValidToken_PassesAndSetsClaims", func(t *testing.T) {
		expectedUserID := uuid.New()
		expectedOrgID := uuid.New()
		mockService := &mockTokenService{
			validateFn: func(token string) (map[string]any, error) {
				return map[string]any{
					"user_id":         expectedUserID,
					"email":           "test@example.com",
					"organization_id": expectedOrgID,
				}, nil
			},
		}

		router := gin.New()
		router.Use(middleware.AuthMiddleware(mockService))
		router.GET("/protected", func(c *gin.Context) {
			uid, err := middleware.GetUserID(c)
			if err != nil || uid != expectedUserID {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			orgID, err := middleware.GetOrgID(c)
			if err != nil || orgID != expectedOrgID {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			email, err := middleware.GetUserEmail(c)
			if err != nil || email != "test@example.com" {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
	})
}

func TestContextHelpers_ErrorScenarios(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("GetUserID_Unset_ReturnsErrUnauthorized", func(t *testing.T) {
		_, err := middleware.GetUserID(c)
		if !errors.Is(err, middleware.ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("GetOrgID_UnsetWithoutHeader_ReturnsErrInvalidOrgID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		c.Request = req

		_, err := middleware.GetOrgID(c)
		if !errors.Is(err, middleware.ErrInvalidOrgID) {
			t.Fatalf("expected ErrInvalidOrgID, got %v", err)
		}
	})

	t.Run("GetOrgID_WithValidHeaderFallback_Succeeds", func(t *testing.T) {
		orgID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Organization-ID", orgID.String())
		c.Request = req

		id, err := middleware.GetOrgID(c)
		if err != nil || id != orgID {
			t.Fatalf("expected fallback to header orgID %v, got %v (err: %v)", orgID, id, err)
		}
	})

	t.Run("GetUserEmail_Unset_ReturnsErrUnauthorized", func(t *testing.T) {
		_, err := middleware.GetUserEmail(c)
		if !errors.Is(err, middleware.ErrUnauthorized) {
			t.Fatalf("expected ErrUnauthorized, got %v", err)
		}
	})
}

func TestRecoveryMiddleware_PanicScenarios(t *testing.T) {
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	router.GET("/panic-nil", func(c *gin.Context) {
		var ptr *string
		_ = *ptr // triggers nil pointer dereference panic
	})

	router.GET("/panic-string", func(c *gin.Context) {
		panic("database connection collapsed")
	})

	t.Run("HandlerPanicsWithNilPointer_RecoversAndReturns500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic-nil", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500 on recovered panic, got %d", w.Code)
		}
	})

	t.Run("HandlerPanicsWithStringMessage_RecoversAndReturns500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic-string", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500 on recovered panic, got %d", w.Code)
		}
	})
}
