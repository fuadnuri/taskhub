package orgrepo

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationRepository implements irepository.IOrganizationRepository.
type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *entities.OrganizationEntity) error {
	row := toOrgSchema(org)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("OrganizationRepository.Create: %w", err)
	}
	org.ID = row.ID
	return nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationEntity, error) {
	var row schema.Organization
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("OrganizationRepository.GetByID: %w", err)
	}
	e := toOrgEntity(&row)
	return &e, nil
}

func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*entities.OrganizationEntity, error) {
	var row schema.Organization
	if err := r.db.WithContext(ctx).First(&row, "slug = ?", slug).Error; err != nil {
		return nil, fmt.Errorf("OrganizationRepository.GetBySlug: %w", err)
	}
	e := toOrgEntity(&row)
	return &e, nil
}

func (r *OrganizationRepository) GetAll(ctx context.Context) ([]entities.OrganizationEntity, error) {
	var rows []schema.Organization
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("OrganizationRepository.GetAll: %w", err)
	}
	result := make([]entities.OrganizationEntity, len(rows))
	for i, row := range rows {
		result[i] = toOrgEntity(&row)
	}
	return result, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *entities.OrganizationEntity) error {
	row := toOrgSchema(org)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("OrganizationRepository.Update: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Organization{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("OrganizationRepository.Delete: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toOrgSchema(e *entities.OrganizationEntity) schema.Organization {
	return schema.Organization{
		ID:       e.ID,
		Name:     e.Name,
		Slug:     e.Slug,
		IsActive: e.IsActive,
	}
}

func toOrgEntity(s *schema.Organization) entities.OrganizationEntity {
	return entities.OrganizationEntity{
		ID:        s.ID,
		Name:      s.Name,
		Slug:      s.Slug,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
