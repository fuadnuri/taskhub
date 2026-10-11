package service

import (
	"context"
	"fmt"
	"time"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/core/domain/irepository"
	"github.com/google/uuid"
)

// ProjectService implements IProjectService.
type ProjectService struct {
	projectRepo   irepository.IProjectRepository
	projectMember irepository.IProjectMemberRepository
}

func NewProjectService(
	projectRepo irepository.IProjectRepository,
	projectMember irepository.IProjectMemberRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo:   projectRepo,
		projectMember: projectMember,
	}
}

func (s *ProjectService) Create(ctx context.Context, orgID uuid.UUID, createdBy uuid.UUID, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	status := req.Status
	if status == "" {
		status = "planning"
	}

	project := &entities.ProjectEntity{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CreatedBy:      createdBy,
		Name:           req.Name,
		Description:    req.Description,
		Status:         status,
		StartDate:      req.StartDate,
		DueDate:        req.DueDate,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add creator as initial member
	if s.projectMember != nil && createdBy != uuid.Nil {
		_ = s.projectMember.Add(ctx, &entities.ProjectMemberEntity{
			ProjectID: project.ID,
			UserID:    createdBy,
			JoinedAt:  time.Now(),
		})
	}

	res := toProjectResponse(project)
	return &res, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toProjectResponse(project)
	return &res, nil
}

func (s *ProjectService) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]dto.ProjectResponse, error) {
	projects, err := s.projectRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	result := make([]dto.ProjectResponse, len(projects))
	for i := range projects {
		result[i] = toProjectResponse(&projects[i])
	}
	return result, nil
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil || project == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.StartDate != nil {
		project.StartDate = req.StartDate
	}
	if req.DueDate != nil {
		project.DueDate = req.DueDate
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	res := toProjectResponse(project)
	return &res, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) AddMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	if s.projectMember == nil {
		return apperrors.ErrInternal
	}
	return s.projectMember.Add(ctx, &entities.ProjectMemberEntity{
		ProjectID: projectID,
		UserID:    userID,
		JoinedAt:  time.Now(),
	})
}

func (s *ProjectService) RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	if s.projectMember == nil {
		return apperrors.ErrInternal
	}
	return s.projectMember.Remove(ctx, projectID, userID)
}

func (s *ProjectService) GetMembers(ctx context.Context, projectID uuid.UUID) ([]dto.ProjectMemberResponse, error) {
	if s.projectMember == nil {
		return nil, nil
	}
	members, err := s.projectMember.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project members: %w", err)
	}

	result := make([]dto.ProjectMemberResponse, len(members))
	for i, m := range members {
		result[i] = dto.ProjectMemberResponse{
			ProjectID: m.ProjectID,
			UserID:    m.UserID,
			JoinedAt:  m.JoinedAt,
		}
	}
	return result, nil
}

func toProjectResponse(p *entities.ProjectEntity) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:             p.ID,
		OrganizationID: p.OrganizationID,
		CreatedBy:      p.CreatedBy,
		Name:           p.Name,
		Description:    p.Description,
		Status:         p.Status,
		StartDate:      p.StartDate,
		DueDate:        p.DueDate,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
