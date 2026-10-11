package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/application/service"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/google/uuid"
)

// ── In-Memory Mocks ─────────────────────────────────────────────────────────

type mockUserRepo struct {
	users map[uuid.UUID]*entities.UserEntity
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uuid.UUID]*entities.UserEntity)}
}

func (m *mockUserRepo) Create(ctx context.Context, u *entities.UserEntity) error {
	for _, existing := range m.users {
		if existing.Email == u.Email {
			return errors.New("duplicate key")
		}
	}
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetAll(ctx context.Context) ([]entities.UserEntity, error) {
	var list []entities.UserEntity
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *entities.UserEntity) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.users, id)
	return nil
}

type mockTokenRepo struct {
	tokens map[string]*entities.RefreshTokenEntity
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*entities.RefreshTokenEntity)}
}

func (m *mockTokenRepo) Create(ctx context.Context, t *entities.RefreshTokenEntity) error {
	m.tokens[t.TokenHash] = t
	return nil
}

func (m *mockTokenRepo) GetByTokenHash(ctx context.Context, hash string) (*entities.RefreshTokenEntity, error) {
	t, ok := m.tokens[hash]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockTokenRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.RefreshTokenEntity, error) {
	var res []entities.RefreshTokenEntity
	for _, t := range m.tokens {
		if t.UserID == userID {
			res = append(res, *t)
		}
	}
	return res, nil
}

func (m *mockTokenRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	for _, t := range m.tokens {
		if t.ID == id {
			now := time.Now()
			t.RevokedAt = &now
			return nil
		}
	}
	return nil
}

func (m *mockTokenRepo) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	for _, t := range m.tokens {
		if t.UserID == userID {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (m *mockTokenRepo) DeleteExpired(ctx context.Context) error {
	return nil
}

type mockTokenService struct{}

func (m *mockTokenService) GenerateAccessToken(userID uuid.UUID, email string, orgID uuid.UUID) (string, error) {
	return "mock-access-token", nil
}

func (m *mockTokenService) ValidateAccessToken(tokenString string) (map[string]any, error) {
	return map[string]any{"user_id": uuid.New()}, nil
}

func (m *mockTokenService) GenerateRefreshToken() (string, string, time.Time, error) {
	return "raw-refresh-token", "hashed-refresh-token", time.Now().Add(24 * time.Hour), nil
}

func (m *mockTokenService) HashRefreshToken(raw string) string {
	return "hashed-" + raw
}

type mockPasswordService struct{}

func (m *mockPasswordService) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *mockPasswordService) Verify(password, hash string) error {
	if hash == "hashed-"+password {
		return nil
	}
	return errors.New("password mismatch")
}

// ── Tests: AuthService Error Flows ──────────────────────────────────────────

func TestAuthService_ErrorFlows(t *testing.T) {
	ctx := context.Background()
	userRepo := newMockUserRepo()
	tokenRepo := newMockTokenRepo()
	tokenSvc := &mockTokenService{}
	pwdSvc := &mockPasswordService{}

	authService := service.NewAuthService(userRepo, tokenRepo, nil, nil, tokenSvc, pwdSvc)

	// Preload existing user
	existingID := uuid.New()
	userRepo.users[existingID] = &entities.UserEntity{
		ID:           existingID,
		Email:        "existing@example.com",
		PasswordHash: "hashed-correctpassword",
		IsActive:     true,
	}

	// Inactive user
	inactiveID := uuid.New()
	userRepo.users[inactiveID] = &entities.UserEntity{
		ID:           inactiveID,
		Email:        "inactive@example.com",
		PasswordHash: "hashed-correctpassword",
		IsActive:     false,
	}

	t.Run("Register_DuplicateEmail_ReturnsErrAlreadyExists", func(t *testing.T) {
		req := dto.RegisterRequest{
			Email:     "existing@example.com",
			Password:  "newpassword",
			FirstName: "Jane",
			LastName:  "Doe",
		}
		_, err := authService.Register(ctx, req)
		if !errors.Is(err, apperrors.ErrAlreadyExists) {
			t.Fatalf("expected ErrAlreadyExists, got %v", err)
		}
	})

	t.Run("Login_UnknownEmail_ReturnsErrInvalidCredentials", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "unknown@example.com",
			Password: "correctpassword",
		}
		_, err := authService.Login(ctx, req)
		if !errors.Is(err, apperrors.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("Login_WrongPassword_ReturnsErrInvalidCredentials", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "existing@example.com",
			Password: "wrongpassword",
		}
		_, err := authService.Login(ctx, req)
		if !errors.Is(err, apperrors.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("Login_InactiveUser_ReturnsErrUserInactive", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "inactive@example.com",
			Password: "correctpassword",
		}
		_, err := authService.Login(ctx, req)
		if !errors.Is(err, apperrors.ErrUserInactive) {
			t.Fatalf("expected ErrUserInactive, got %v", err)
		}
	})

	t.Run("RefreshToken_NonExistentToken_ReturnsErrInvalidToken", func(t *testing.T) {
		req := dto.RefreshTokenRequest{
			RefreshToken: "non-existent-token",
		}
		_, err := authService.RefreshToken(ctx, req)
		if !errors.Is(err, apperrors.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("RefreshToken_ExpiredToken_ReturnsErrTokenExpired", func(t *testing.T) {
		expiredHash := tokenSvc.HashRefreshToken("expired-token")
		tokenRepo.tokens[expiredHash] = &entities.RefreshTokenEntity{
			ID:        uuid.New(),
			UserID:    existingID,
			TokenHash: expiredHash,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}

		req := dto.RefreshTokenRequest{
			RefreshToken: "expired-token",
		}
		_, err := authService.RefreshToken(ctx, req)
		if !errors.Is(err, apperrors.ErrTokenExpired) {
			t.Fatalf("expected ErrTokenExpired, got %v", err)
		}
	})

	t.Run("RefreshToken_RevokedToken_ReturnsErrInvalidToken", func(t *testing.T) {
		revokedHash := tokenSvc.HashRefreshToken("revoked-token")
		now := time.Now()
		tokenRepo.tokens[revokedHash] = &entities.RefreshTokenEntity{
			ID:        uuid.New(),
			UserID:    existingID,
			TokenHash: revokedHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			RevokedAt: &now,
		}

		req := dto.RefreshTokenRequest{
			RefreshToken: "revoked-token",
		}
		_, err := authService.RefreshToken(ctx, req)
		if !errors.Is(err, apperrors.ErrInvalidToken) {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})
}

// ── Tests: UserService Error Flows ──────────────────────────────────────────

func TestUserService_ErrorFlows(t *testing.T) {
	ctx := context.Background()
	userRepo := newMockUserRepo()
	pwdSvc := &mockPasswordService{}
	userService := service.NewUserService(userRepo, pwdSvc)

	existingID := uuid.New()
	userRepo.users[existingID] = &entities.UserEntity{
		ID:           existingID,
		Email:        "existing@example.com",
		PasswordHash: "hashed-oldpass",
		IsActive:     true,
	}

	t.Run("CreateUser_DuplicateEmail_ReturnsErrAlreadyExists", func(t *testing.T) {
		req := dto.CreateUserRequest{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "Alice",
			LastName:  "Smith",
		}
		_, err := userService.CreateUser(ctx, req)
		if !errors.Is(err, apperrors.ErrAlreadyExists) {
			t.Fatalf("expected ErrAlreadyExists, got %v", err)
		}
	})

	t.Run("GetByID_NotFound_ReturnsErrNotFound", func(t *testing.T) {
		_, err := userService.GetByID(ctx, uuid.New())
		if !errors.Is(err, apperrors.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("GetByEmail_NotFound_ReturnsErrNotFound", func(t *testing.T) {
		_, err := userService.GetByEmail(ctx, "does-not-exist@example.com")
		if !errors.Is(err, apperrors.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Update_NotFound_ReturnsErrNotFound", func(t *testing.T) {
		name := "NewName"
		_, err := userService.Update(ctx, uuid.New(), dto.UpdateUserRequest{FirstName: &name})
		if !errors.Is(err, apperrors.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("ChangePassword_UserNotFound_ReturnsErrNotFound", func(t *testing.T) {
		err := userService.ChangePassword(ctx, uuid.New(), dto.ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "newpassword123",
		})
		if !errors.Is(err, apperrors.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("ChangePassword_WrongOldPassword_ReturnsErrInvalidCredentials", func(t *testing.T) {
		err := userService.ChangePassword(ctx, existingID, dto.ChangePasswordRequest{
			OldPassword: "wrongoldpassword",
			NewPassword: "newpassword123",
		})
		if !errors.Is(err, apperrors.ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}
