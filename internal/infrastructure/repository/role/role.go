package rolerepo

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepository implements irepository.IRoleRepository.
type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *entities.RoleEntity) error {
	row := toRoleSchema(role)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("RoleRepository.Create: %w", err)
	}
	role.ID = row.ID
	return nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.RoleEntity, error) {
	var row schema.Role
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("RoleRepository.GetByID: %w", err)
	}
	e := toRoleEntity(&row)
	return &e, nil
}

func (r *RoleRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.RoleEntity, error) {
	var rows []schema.Role
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("RoleRepository.GetByOrganizationID: %w", err)
	}
	result := make([]entities.RoleEntity, len(rows))
	for i, row := range rows {
		result[i] = toRoleEntity(&row)
	}
	return result, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *entities.RoleEntity) error {
	row := toRoleSchema(role)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("RoleRepository.Update: %w", err)
	}
	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Role{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("RoleRepository.Delete: %w", err)
	}
	return nil
}

// ── PermissionRepository ──────────────────────────────────────────────────────

// PermissionRepository implements irepository.IPermissionRepository.
type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, perm *entities.PermissionEntity) error {
	row := toPermissionSchema(perm)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("PermissionRepository.Create: %w", err)
	}
	perm.ID = row.ID
	return nil
}

func (r *PermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.PermissionEntity, error) {
	var row schema.Permission
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("PermissionRepository.GetByID: %w", err)
	}
	e := toPermissionEntity(&row)
	return &e, nil
}

func (r *PermissionRepository) GetByName(ctx context.Context, name string) (*entities.PermissionEntity, error) {
	var row schema.Permission
	if err := r.db.WithContext(ctx).First(&row, "name = ?", name).Error; err != nil {
		return nil, fmt.Errorf("PermissionRepository.GetByName: %w", err)
	}
	e := toPermissionEntity(&row)
	return &e, nil
}

func (r *PermissionRepository) GetAll(ctx context.Context) ([]entities.PermissionEntity, error) {
	var rows []schema.Permission
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("PermissionRepository.GetAll: %w", err)
	}
	result := make([]entities.PermissionEntity, len(rows))
	for i, row := range rows {
		result[i] = toPermissionEntity(&row)
	}
	return result, nil
}

func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Permission{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("PermissionRepository.Delete: %w", err)
	}
	return nil
}

// ── RolePermissionRepository ──────────────────────────────────────────────────

// RolePermissionRepository implements irepository.IRolePermissionRepository.
type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{db: db}
}

func (r *RolePermissionRepository) Assign(ctx context.Context, rp *entities.RolePermissionEntity) error {
	row := schema.RolePermission{RoleID: rp.RoleID, PermissionID: rp.PermissionID}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("RolePermissionRepository.Assign: %w", err)
	}
	return nil
}

func (r *RolePermissionRepository) Revoke(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&schema.RolePermission{}).Error; err != nil {
		return fmt.Errorf("RolePermissionRepository.Revoke: %w", err)
	}
	return nil
}

func (r *RolePermissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entities.RolePermissionEntity, error) {
	var rows []schema.RolePermission
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("RolePermissionRepository.GetByRoleID: %w", err)
	}
	result := make([]entities.RolePermissionEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.RolePermissionEntity{RoleID: row.RoleID, PermissionID: row.PermissionID}
	}
	return result, nil
}

func (r *RolePermissionRepository) GetByPermissionID(ctx context.Context, permissionID uuid.UUID) ([]entities.RolePermissionEntity, error) {
	var rows []schema.RolePermission
	if err := r.db.WithContext(ctx).Where("permission_id = ?", permissionID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("RolePermissionRepository.GetByPermissionID: %w", err)
	}
	result := make([]entities.RolePermissionEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.RolePermissionEntity{RoleID: row.RoleID, PermissionID: row.PermissionID}
	}
	return result, nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toRoleSchema(e *entities.RoleEntity) schema.Role {
	return schema.Role{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		Name:           e.Name,
		Description:    e.Description,
	}
}

func toRoleEntity(s *schema.Role) entities.RoleEntity {
	return entities.RoleEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		Name:           s.Name,
		Description:    s.Description,
		CreatedAt:      s.CreatedAt,
	}
}

func toPermissionSchema(e *entities.PermissionEntity) schema.Permission {
	return schema.Permission{ID: e.ID, Name: e.Name, Description: e.Description}
}

func toPermissionEntity(s *schema.Permission) entities.PermissionEntity {
	return entities.PermissionEntity{ID: s.ID, Name: s.Name, Description: s.Description}
}
