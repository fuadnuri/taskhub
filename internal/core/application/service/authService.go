package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/core/domain/irepository"
	"github.com/google/uuid"
)

// AuthService implements IAuthService.
type AuthService struct {
	userRepo        irepository.IUserRepository
	tokenRepo       irepository.IRefreshTokenRepository
	orgRepo         irepository.IOrganizationRepository
	memberRepo      irepository.IOrganizationMemberRepository
	tokenService    ITokenService
	passwordService IPasswordService
}

func NewAuthService(
	userRepo irepository.IUserRepository,
	tokenRepo irepository.IRefreshTokenRepository,
	orgRepo irepository.IOrganizationRepository,
	memberRepo irepository.IOrganizationMemberRepository,
	tokenService ITokenService,
	passwordService IPasswordService,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		orgRepo:         orgRepo,
		memberRepo:      memberRepo,
		tokenService:    tokenService,
		passwordService: passwordService,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 1. Check if user already exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	// 2. Hash password
	pwdHash, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Create user
	user := &entities.UserEntity{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  pwdHash,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		IsActive:      true,
		EmailVerified: false,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 4. Create default organization if requested
	var orgID uuid.UUID
	if req.OrgName != "" && s.orgRepo != nil {
		slug := strings.ToLower(strings.ReplaceAll(req.OrgName, " ", "-"))
		org := &entities.OrganizationEntity{
			ID:       uuid.New(),
			Name:     req.OrgName,
			Slug:     slug,
			IsActive: true,
		}
		if err := s.orgRepo.Create(ctx, org); err == nil {
			orgID = org.ID
			if s.memberRepo != nil {
				_ = s.memberRepo.Add(ctx, &entities.OrganizationMemberEntity{
					ID:             uuid.New(),
					OrganizationID: org.ID,
					UserID:         user.ID,
					JoinedAt:       time.Now(),
				})
			}
		}
	}

	// 5. Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, refreshHash, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if s.tokenRepo != nil {
		_ = s.tokenRepo.Create(ctx, &entities.RefreshTokenEntity{
			ID:        uuid.New(),
			UserID:    user.ID,
			TokenHash: refreshHash,
			ExpiresAt: expiresAt,
		})
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		User:         toUserResponse(user),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// 1. Fetch user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, apperrors.ErrUserInactive
	}

	// 2. Verify password
	if err := s.passwordService.Verify(req.Password, user.PasswordHash); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 3. Determine active org
	var orgID uuid.UUID
	if req.OrganizationID != nil {
		orgID = *req.OrganizationID
	} else if s.memberRepo != nil {
		memberships, _ := s.memberRepo.GetByUserID(ctx, user.ID)
		if len(memberships) > 0 {
			orgID = memberships[0].OrganizationID
		}
	}

	// 4. Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, refreshHash, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if s.tokenRepo != nil {
		_ = s.tokenRepo.Create(ctx, &entities.RefreshTokenEntity{
			ID:        uuid.New(),
			UserID:    user.ID,
			TokenHash: refreshHash,
			ExpiresAt: expiresAt,
		})
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		User:         toUserResponse(user),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	if s.tokenRepo == nil {
		return nil, apperrors.ErrInternal
	}

	tokenHash := s.tokenService.HashRefreshToken(req.RefreshToken)
	tokenRecord, err := s.tokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil || tokenRecord == nil || tokenRecord.RevokedAt != nil {
		return nil, apperrors.ErrInvalidToken
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, apperrors.ErrTokenExpired
	}

	user, err := s.userRepo.GetByID(ctx, tokenRecord.UserID)
	if err != nil || user == nil || !user.IsActive {
		return nil, apperrors.ErrUnauthorized
	}

	// Revoke old refresh token
	_ = s.tokenRepo.Revoke(ctx, tokenRecord.ID)

	// Determine active org
	var orgID uuid.UUID
	if s.memberRepo != nil {
		memberships, _ := s.memberRepo.GetByUserID(ctx, user.ID)
		if len(memberships) > 0 {
			orgID = memberships[0].OrganizationID
		}
	}

	// Generate new tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRawRefresh, newRefreshHash, expiresAt, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	_ = s.tokenRepo.Create(ctx, &entities.RefreshTokenEntity{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: newRefreshHash,
		ExpiresAt: expiresAt,
	})

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRawRefresh,
		User:         toUserResponse(user),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if s.tokenRepo == nil || refreshToken == "" {
		return nil
	}
	tokenHash := s.tokenService.HashRefreshToken(refreshToken)
	tokenRecord, err := s.tokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil || tokenRecord == nil {
		return nil
	}
	return s.tokenRepo.Revoke(ctx, tokenRecord.ID)
}

func toUserResponse(u *entities.UserEntity) dto.UserResponse {
	return dto.UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}