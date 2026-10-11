package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenRepository implements irepository.IRefreshTokenRepository using GORM.
type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entities.RefreshTokenEntity) error {
	row := toRefreshTokenSchema(token)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("RefreshTokenRepository.Create: %w", err)
	}
	token.ID = row.ID
	token.CreatedAt = row.CreatedAt
	return nil
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entities.RefreshTokenEntity, error) {
	var row schema.RefreshToken
	if err := r.db.WithContext(ctx).First(&row, "token_hash = ?", tokenHash).Error; err != nil {
		return nil, fmt.Errorf("RefreshTokenRepository.GetByTokenHash: %w", err)
	}
	e := toRefreshTokenEntity(&row)
	return &e, nil
}

func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.RefreshTokenEntity, error) {
	var rows []schema.RefreshToken
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("RefreshTokenRepository.GetByUserID: %w", err)
	}
	result := make([]entities.RefreshTokenEntity, len(rows))
	for i, row := range rows {
		result[i] = toRefreshTokenEntity(&row)
	}
	return result, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&schema.RefreshToken{}).
		Where("id = ?", id).Update("revoked_at", &now).Error; err != nil {
		return fmt.Errorf("RefreshTokenRepository.Revoke: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&schema.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &now).Error; err != nil {
		return fmt.Errorf("RefreshTokenRepository.RevokeAllByUserID: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", now).
		Delete(&schema.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("RefreshTokenRepository.DeleteExpired: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toRefreshTokenSchema(e *entities.RefreshTokenEntity) schema.RefreshToken {
	return schema.RefreshToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		RevokedAt: e.RevokedAt,
	}
}

func toRefreshTokenEntity(s *schema.RefreshToken) entities.RefreshTokenEntity {
	return entities.RefreshTokenEntity{
		ID:        s.ID,
		UserID:    s.UserID,
		TokenHash: s.TokenHash,
		ExpiresAt: s.ExpiresAt,
		RevokedAt: s.RevokedAt,
		CreatedAt: s.CreatedAt,
	}
}
