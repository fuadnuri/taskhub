package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims holds the JWT payload for access tokens.
type Claims struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTConfig holds tunable JWT parameters.
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// JWTService handles creation and validation of JWT access tokens
// and opaque refresh tokens.
type JWTService struct {
	cfg JWTConfig
}

// NewJWTService creates a JWTService with the provided config.
func NewJWTService(cfg JWTConfig) *JWTService {
	return &JWTService{cfg: cfg}
}

// GenerateAccessToken issues a signed JWT access token for a user.
func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string, orgID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:         userID,
		Email:          email,
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("signing access token: %w", err)
	}
	return signed, nil
}

// ValidateAccessToken parses and validates a JWT access token,
// returning a claims map with "user_id", "email", "organization_id", "expires_at".
func (s *JWTService) ValidateAccessToken(tokenString string) (map[string]any, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.cfg.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing access token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	return map[string]any{
		"user_id":         claims.UserID,
		"email":           claims.Email,
		"organization_id": claims.OrganizationID,
		"expires_at":      claims.ExpiresAt.Time,
	}, nil
}

// GenerateRefreshToken creates a cryptographically random opaque token
// and returns both the raw token (to send to the client) and its SHA-256
// hex-encoded hash (to store in the database).
func (s *JWTService) GenerateRefreshToken() (raw string, hash string, expiresAt time.Time, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generating refresh token bytes: %w", err)
	}

	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	expiresAt = time.Now().Add(s.cfg.RefreshTokenDuration)
	return raw, hash, expiresAt, nil
}

// HashRefreshToken returns the hex-encoded SHA-256 hash of a raw refresh token.
// Use this when looking up a token presented by the client.
func (s *JWTService) HashRefreshToken(raw string) string {
	return hashToken(raw)
}
