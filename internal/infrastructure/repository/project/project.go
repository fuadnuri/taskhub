package projectrepo

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectRepository implements irepository.IProjectRepository.
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *entities.ProjectEntity) error {
	row := toProjectSchema(project)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("ProjectRepository.Create: %w", err)
	}
	project.ID = row.ID
	return nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.ProjectEntity, error) {
	var row schema.Project
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("ProjectRepository.GetByID: %w", err)
	}
	e := toProjectEntity(&row)
	return &e, nil
}

func (r *ProjectRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.ProjectEntity, error) {
	var rows []schema.Project
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ProjectRepository.GetByOrganizationID: %w", err)
	}
	result := make([]entities.ProjectEntity, len(rows))
	for i, row := range rows {
		result[i] = toProjectEntity(&row)
	}
	return result, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *entities.ProjectEntity) error {
	row := toProjectSchema(project)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("ProjectRepository.Update: %w", err)
	}
	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Project{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("ProjectRepository.Delete: %w", err)
	}
	return nil
}

// ── ProjectMemberRepository ───────────────────────────────────────────────────

// ProjectMemberRepository implements irepository.IProjectMemberRepository.
type ProjectMemberRepository struct {
	db *gorm.DB
}

func NewProjectMemberRepository(db *gorm.DB) *ProjectMemberRepository {
	return &ProjectMemberRepository{db: db}
}

func (r *ProjectMemberRepository) Add(ctx context.Context, member *entities.ProjectMemberEntity) error {
	row := schema.ProjectMember{ProjectID: member.ProjectID, UserID: member.UserID}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("ProjectMemberRepository.Add: %w", err)
	}
	return nil
}

func (r *ProjectMemberRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]entities.ProjectMemberEntity, error) {
	var rows []schema.ProjectMember
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ProjectMemberRepository.GetByProjectID: %w", err)
	}
	result := make([]entities.ProjectMemberEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.ProjectMemberEntity{ProjectID: row.ProjectID, UserID: row.UserID, JoinedAt: row.JoinedAt}
	}
	return result, nil
}

func (r *ProjectMemberRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.ProjectMemberEntity, error) {
	var rows []schema.ProjectMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("ProjectMemberRepository.GetByUserID: %w", err)
	}
	result := make([]entities.ProjectMemberEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.ProjectMemberEntity{ProjectID: row.ProjectID, UserID: row.UserID, JoinedAt: row.JoinedAt}
	}
	return result, nil
}

func (r *ProjectMemberRepository) Remove(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&schema.ProjectMember{}).Error; err != nil {
		return fmt.Errorf("ProjectMemberRepository.Remove: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toProjectSchema(e *entities.ProjectEntity) schema.Project {
	return schema.Project{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		CreatedBy:      e.CreatedBy,
		Name:           e.Name,
		Description:    e.Description,
		Status:         e.Status,
		StartDate:      e.StartDate,
		DueDate:        e.DueDate,
	}
}

func toProjectEntity(s *schema.Project) entities.ProjectEntity {
	return entities.ProjectEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		CreatedBy:      s.CreatedBy,
		Name:           s.Name,
		Description:    s.Description,
		Status:         s.Status,
		StartDate:      s.StartDate,
		DueDate:        s.DueDate,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}
