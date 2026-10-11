package userrepo

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository implements irepository.IUserRepository using GORM.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.UserEntity) error {
	row := toUserSchema(user)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	user.ID = row.ID
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error) {
	var row schema.User
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.GetByID: %w", err)
	}
	e := toUserEntity(&row)
	return &e, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.UserEntity, error) {
	var row schema.User
	if err := r.db.WithContext(ctx).First(&row, "email = ?", email).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.GetByEmail: %w", err)
	}
	e := toUserEntity(&row)
	return &e, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]entities.UserEntity, error) {
	var rows []schema.User
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("UserRepository.GetAll: %w", err)
	}
	result := make([]entities.UserEntity, len(rows))
	for i, row := range rows {
		result[i] = toUserEntity(&row)
	}
	return result, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.UserEntity) error {
	row := toUserSchema(user)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("UserRepository.Update: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("UserRepository.Delete: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toUserSchema(e *entities.UserEntity) schema.User {
	return schema.User{
		ID:            e.ID,
		Email:         e.Email,
		PasswordHash:  e.PasswordHash,
		FirstName:     e.FirstName,
		LastName:      e.LastName,
		IsActive:      e.IsActive,
		EmailVerified: e.EmailVerified,
	}
}

func toUserEntity(s *schema.User) entities.UserEntity {
	return entities.UserEntity{
		ID:            s.ID,
		Email:         s.Email,
		PasswordHash:  s.PasswordHash,
		FirstName:     s.FirstName,
		LastName:      s.LastName,
		IsActive:      s.IsActive,
		EmailVerified: s.EmailVerified,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}