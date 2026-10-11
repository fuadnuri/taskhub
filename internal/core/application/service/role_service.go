package service

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/core/domain/irepository"
	"github.com/google/uuid"
)

// RoleService implements IRoleService.
type RoleService struct {
	roleRepo     irepository.IRoleRepository
	permRepo     irepository.IPermissionRepository
	rolePermRepo irepository.IRolePermissionRepository
}

func NewRoleService(
	roleRepo irepository.IRoleRepository,
	permRepo irepository.IPermissionRepository,
	rolePermRepo irepository.IRolePermissionRepository,
) *RoleService {
	return &RoleService{
		roleRepo:     roleRepo,
		permRepo:     permRepo,
		rolePermRepo: rolePermRepo,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, orgID uuid.UUID, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &entities.RoleEntity{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Assign initial permissions if provided
	if s.rolePermRepo != nil && len(req.PermissionIDs) > 0 {
		for _, permID := range req.PermissionIDs {
			_ = s.rolePermRepo.Assign(ctx, &entities.RolePermissionEntity{
				RoleID:       role.ID,
				PermissionID: permID,
			})
		}
	}

	return s.getRoleResponse(ctx, role)
}

func (s *RoleService) GetRoleByID(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		return nil, apperrors.ErrNotFound
	}
	return s.getRoleResponse(ctx, role)
}

func (s *RoleService) GetRolesByOrgID(ctx context.Context, orgID uuid.UUID) ([]dto.RoleResponse, error) {
	roles, err := s.roleRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	result := make([]dto.RoleResponse, len(roles))
	for i := range roles {
		r, err := s.getRoleResponse(ctx, &roles[i])
		if err != nil {
			return nil, err
		}
		result[i] = *r
	}
	return result, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id uuid.UUID, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil || role == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return s.getRoleResponse(ctx, role)
}

func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) AssignPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error {
	if s.rolePermRepo == nil {
		return apperrors.ErrInternal
	}
	return s.rolePermRepo.Assign(ctx, &entities.RolePermissionEntity{
		RoleID:       roleID,
		PermissionID: permID,
	})
}

func (s *RoleService) RevokePermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error {
	if s.rolePermRepo == nil {
		return apperrors.ErrInternal
	}
	return s.rolePermRepo.Revoke(ctx, roleID, permID)
}

func (s *RoleService) CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	if s.permRepo == nil {
		return nil, apperrors.ErrInternal
	}

	perm := &entities.PermissionEntity{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.permRepo.Create(ctx, perm); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &dto.PermissionResponse{
		ID:          perm.ID,
		Name:        perm.Name,
		Description: perm.Description,
	}, nil
}

func (s *RoleService) GetAllPermissions(ctx context.Context) ([]dto.PermissionResponse, error) {
	if s.permRepo == nil {
		return nil, nil
	}

	perms, err := s.permRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	result := make([]dto.PermissionResponse, len(perms))
	for i, p := range perms {
		result[i] = dto.PermissionResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
		}
	}
	return result, nil
}

func (s *RoleService) getRoleResponse(ctx context.Context, role *entities.RoleEntity) (*dto.RoleResponse, error) {
	var permissions []dto.PermissionResponse
	if s.rolePermRepo != nil && s.permRepo != nil {
		rps, err := s.rolePermRepo.GetByRoleID(ctx, role.ID)
		if err == nil {
			for _, rp := range rps {
				if perm, err := s.permRepo.GetByID(ctx, rp.PermissionID); err == nil && perm != nil {
					permissions = append(permissions, dto.PermissionResponse{
						ID:          perm.ID,
						Name:        perm.Name,
						Description: perm.Description,
					})
				}
			}
		}
	}

	return &dto.RoleResponse{
		ID:             role.ID,
		OrganizationID: role.OrganizationID,
		Name:           role.Name,
		Description:    role.Description,
		CreatedAt:      role.CreatedAt,
		Permissions:    permissions,
	}, nil
}
