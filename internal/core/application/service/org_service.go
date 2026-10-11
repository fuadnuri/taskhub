package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/core/domain/irepository"
	"github.com/google/uuid"
)

// OrganizationService implements IOrganizationService.
type OrganizationService struct {
	orgRepo        irepository.IOrganizationRepository
	memberRepo     irepository.IOrganizationMemberRepository
	invitationRepo irepository.IOrganizationInvitationRepository
	tokenService   ITokenService
}

func NewOrganizationService(
	orgRepo irepository.IOrganizationRepository,
	memberRepo irepository.IOrganizationMemberRepository,
	invitationRepo irepository.IOrganizationInvitationRepository,
	tokenService ITokenService,
) *OrganizationService {
	return &OrganizationService{
		orgRepo:        orgRepo,
		memberRepo:     memberRepo,
		invitationRepo: invitationRepo,
		tokenService:   tokenService,
	}
}

func (s *OrganizationService) Create(ctx context.Context, ownerID uuid.UUID, req dto.CreateOrgRequest) (*dto.OrgResponse, error) {
	existing, _ := s.orgRepo.GetBySlug(ctx, req.Slug)
	if existing != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	org := &entities.OrganizationEntity{
		ID:       uuid.New(),
		Name:     req.Name,
		Slug:     req.Slug,
		IsActive: true,
	}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add creator as member
	if s.memberRepo != nil && ownerID != uuid.Nil {
		_ = s.memberRepo.Add(ctx, &entities.OrganizationMemberEntity{
			ID:             uuid.New(),
			OrganizationID: org.ID,
			UserID:         ownerID,
			JoinedAt:       time.Now(),
		})
	}

	res := toOrgResponse(org)
	return &res, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*dto.OrgResponse, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil || org == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toOrgResponse(org)
	return &res, nil
}

func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*dto.OrgResponse, error) {
	org, err := s.orgRepo.GetBySlug(ctx, slug)
	if err != nil || org == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toOrgResponse(org)
	return &res, nil
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]dto.OrgResponse, error) {
	orgs, err := s.orgRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	result := make([]dto.OrgResponse, len(orgs))
	for i := range orgs {
		result[i] = toOrgResponse(&orgs[i])
	}
	return result, nil
}

func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateOrgRequest) (*dto.OrgResponse, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil || org == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Slug != nil {
		org.Slug = *req.Slug
	}
	if req.IsActive != nil {
		org.IsActive = *req.IsActive
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	res := toOrgResponse(org)
	return &res, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orgRepo.Delete(ctx, id)
}

func (s *OrganizationService) InviteMember(ctx context.Context, orgID uuid.UUID, req dto.InviteMemberRequest) (*dto.OrgInvitationResponse, error) {
	if s.invitationRepo == nil {
		return nil, apperrors.ErrInternal
	}

	// Generate random token and hash it
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}
	rawToken := hex.EncodeToString(b)

	var tokenHash string
	if s.tokenService != nil {
		tokenHash = s.tokenService.HashRefreshToken(rawToken)
	} else {
		tokenHash = rawToken
	}

	inv := &entities.OrganizationInvitationEntity{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          req.Email,
		RoleID:         req.RoleID,
		TokenHash:      tokenHash,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.invitationRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	return &dto.OrgInvitationResponse{
		ID:             inv.ID,
		OrganizationID: inv.OrganizationID,
		Email:          inv.Email,
		RoleID:         inv.RoleID,
		ExpiresAt:      inv.ExpiresAt,
		AcceptedAt:     inv.AcceptedAt,
		CreatedAt:      inv.CreatedAt,
	}, nil
}

func (s *OrganizationService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error {
	if s.invitationRepo == nil || s.memberRepo == nil {
		return apperrors.ErrInternal
	}

	var tokenHash string
	if s.tokenService != nil {
		tokenHash = s.tokenService.HashRefreshToken(token)
	} else {
		tokenHash = token
	}

	inv, err := s.invitationRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil || inv == nil {
		return apperrors.ErrNotFound
	}

	if inv.AcceptedAt != nil {
		return apperrors.ErrInvitationAccepted
	}

	if time.Now().After(inv.ExpiresAt) {
		return apperrors.ErrInvitationExpired
	}

	// Mark invitation accepted
	if err := s.invitationRepo.Accept(ctx, inv.ID); err != nil {
		return fmt.Errorf("failed to accept invitation: %w", err)
	}

	// Add user to members
	member := &entities.OrganizationMemberEntity{
		ID:             uuid.New(),
		OrganizationID: inv.OrganizationID,
		UserID:         userID,
		RoleID:         inv.RoleID,
		JoinedAt:       time.Now(),
	}
	return s.memberRepo.Add(ctx, member)
}

func (s *OrganizationService) GetMembers(ctx context.Context, orgID uuid.UUID) ([]dto.OrgMemberResponse, error) {
	if s.memberRepo == nil {
		return nil, nil
	}
	members, err := s.memberRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	result := make([]dto.OrgMemberResponse, len(members))
	for i, m := range members {
		result[i] = dto.OrgMemberResponse{
			ID:             m.ID,
			OrganizationID: m.OrganizationID,
			UserID:         m.UserID,
			RoleID:         m.RoleID,
			JoinedAt:       m.JoinedAt,
		}
	}
	return result, nil
}

func (s *OrganizationService) RemoveMember(ctx context.Context, memberID uuid.UUID) error {
	if s.memberRepo == nil {
		return apperrors.ErrInternal
	}
	return s.memberRepo.Remove(ctx, memberID)
}

func (s *OrganizationService) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, req dto.UpdateMemberRoleRequest) error {
	if s.memberRepo == nil {
		return apperrors.ErrInternal
	}
	return s.memberRepo.UpdateRole(ctx, memberID, req.RoleID)
}

func toOrgResponse(o *entities.OrganizationEntity) dto.OrgResponse {
	return dto.OrgResponse{
		ID:        o.ID,
		Name:      o.Name,
		Slug:      o.Slug,
		IsActive:  o.IsActive,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
