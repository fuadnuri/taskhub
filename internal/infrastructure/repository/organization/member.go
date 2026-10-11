package orgrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrgMemberRepository implements irepository.IOrganizationMemberRepository.
type OrgMemberRepository struct {
	db *gorm.DB
}

func NewOrgMemberRepository(db *gorm.DB) *OrgMemberRepository {
	return &OrgMemberRepository{db: db}
}

func (r *OrgMemberRepository) Add(ctx context.Context, member *entities.OrganizationMemberEntity) error {
	row := toMemberSchema(member)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("OrgMemberRepository.Add: %w", err)
	}
	member.ID = row.ID
	return nil
}

func (r *OrgMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationMemberEntity, error) {
	var row schema.OrganizationMember
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("OrgMemberRepository.GetByID: %w", err)
	}
	e := toMemberEntity(&row)
	return &e, nil
}

func (r *OrgMemberRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.OrganizationMemberEntity, error) {
	var rows []schema.OrganizationMember
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("OrgMemberRepository.GetByOrganizationID: %w", err)
	}
	result := make([]entities.OrganizationMemberEntity, len(rows))
	for i, row := range rows {
		result[i] = toMemberEntity(&row)
	}
	return result, nil
}

func (r *OrgMemberRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.OrganizationMemberEntity, error) {
	var rows []schema.OrganizationMember
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("OrgMemberRepository.GetByUserID: %w", err)
	}
	result := make([]entities.OrganizationMemberEntity, len(rows))
	for i, row := range rows {
		result[i] = toMemberEntity(&row)
	}
	return result, nil
}

func (r *OrgMemberRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&schema.OrganizationMember{}).
		Where("id = ?", id).Update("role_id", roleID).Error; err != nil {
		return fmt.Errorf("OrgMemberRepository.UpdateRole: %w", err)
	}
	return nil
}

func (r *OrgMemberRepository) Remove(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.OrganizationMember{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("OrgMemberRepository.Remove: %w", err)
	}
	return nil
}

// ── OrgInvitationRepository ───────────────────────────────────────────────────

// OrgInvitationRepository implements irepository.IOrganizationInvitationRepository.
type OrgInvitationRepository struct {
	db *gorm.DB
}

func NewOrgInvitationRepository(db *gorm.DB) *OrgInvitationRepository {
	return &OrgInvitationRepository{db: db}
}

func (r *OrgInvitationRepository) Create(ctx context.Context, inv *entities.OrganizationInvitationEntity) error {
	row := toInvitationSchema(inv)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("OrgInvitationRepository.Create: %w", err)
	}
	inv.ID = row.ID
	return nil
}

func (r *OrgInvitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.OrganizationInvitationEntity, error) {
	var row schema.OrganizationInvitation
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("OrgInvitationRepository.GetByID: %w", err)
	}
	e := toInvitationEntity(&row)
	return &e, nil
}

func (r *OrgInvitationRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entities.OrganizationInvitationEntity, error) {
	var row schema.OrganizationInvitation
	if err := r.db.WithContext(ctx).First(&row, "token_hash = ?", tokenHash).Error; err != nil {
		return nil, fmt.Errorf("OrgInvitationRepository.GetByTokenHash: %w", err)
	}
	e := toInvitationEntity(&row)
	return &e, nil
}

func (r *OrgInvitationRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.OrganizationInvitationEntity, error) {
	var rows []schema.OrganizationInvitation
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("OrgInvitationRepository.GetByOrganizationID: %w", err)
	}
	result := make([]entities.OrganizationInvitationEntity, len(rows))
	for i, row := range rows {
		result[i] = toInvitationEntity(&row)
	}
	return result, nil
}

func (r *OrgInvitationRepository) Accept(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&schema.OrganizationInvitation{}).
		Where("id = ?", id).Update("accepted_at", &now).Error; err != nil {
		return fmt.Errorf("OrgInvitationRepository.Accept: %w", err)
	}
	return nil
}

func (r *OrgInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.OrganizationInvitation{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("OrgInvitationRepository.Delete: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toMemberSchema(e *entities.OrganizationMemberEntity) schema.OrganizationMember {
	return schema.OrganizationMember{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		UserID:         e.UserID,
		RoleID:         e.RoleID,
	}
}

func toMemberEntity(s *schema.OrganizationMember) entities.OrganizationMemberEntity {
	return entities.OrganizationMemberEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		UserID:         s.UserID,
		RoleID:         s.RoleID,
		JoinedAt:       s.JoinedAt,
	}
}

func toInvitationSchema(e *entities.OrganizationInvitationEntity) schema.OrganizationInvitation {
	return schema.OrganizationInvitation{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		Email:          e.Email,
		RoleID:         e.RoleID,
		TokenHash:      e.TokenHash,
		ExpiresAt:      e.ExpiresAt,
		AcceptedAt:     e.AcceptedAt,
	}
}

func toInvitationEntity(s *schema.OrganizationInvitation) entities.OrganizationInvitationEntity {
	return entities.OrganizationInvitationEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		Email:          s.Email,
		RoleID:         s.RoleID,
		TokenHash:      s.TokenHash,
		ExpiresAt:      s.ExpiresAt,
		AcceptedAt:     s.AcceptedAt,
		CreatedAt:      s.CreatedAt,
	}
}
