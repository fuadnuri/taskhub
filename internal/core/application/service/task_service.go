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

// TaskService implements ITaskService.
type TaskService struct {
	taskRepo         irepository.ITaskRepository
	assigneeRepo     irepository.ITaskAssigneeRepository
	commentRepo      irepository.ITaskCommentRepository
	notificationRepo irepository.INotificationRepository
}

func NewTaskService(
	taskRepo irepository.ITaskRepository,
	assigneeRepo irepository.ITaskAssigneeRepository,
	commentRepo irepository.ITaskCommentRepository,
	notificationRepo irepository.INotificationRepository,
) *TaskService {
	return &TaskService{
		taskRepo:         taskRepo,
		assigneeRepo:     assigneeRepo,
		commentRepo:      commentRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, orgID uuid.UUID, projectID uuid.UUID, createdBy uuid.UUID, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	status := req.Status
	if status == "" {
		status = "todo"
	}
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	task := &entities.TaskEntity{
		ID:             uuid.New(),
		OrganizationID: orgID,
		ProjectID:      projectID,
		CreatedBy:      createdBy,
		Title:          req.Title,
		Description:    req.Description,
		Status:         status,
		Priority:       priority,
		DueDate:        req.DueDate,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Add assignees if provided
	if s.assigneeRepo != nil && len(req.AssigneeIDs) > 0 {
		for _, assigneeID := range req.AssigneeIDs {
			_ = s.assigneeRepo.Assign(ctx, &entities.TaskAssigneeEntity{
				TaskID:     task.ID,
				UserID:     assigneeID,
				AssignedAt: time.Now(),
			})

			// Optionally notify assignee
			if s.notificationRepo != nil && assigneeID != createdBy {
				_ = s.notificationRepo.Create(ctx, &entities.NotificationEntity{
					ID:             uuid.New(),
					OrganizationID: orgID,
					UserID:         assigneeID,
					Type:           "task_assigned",
					Title:          "New Task Assigned",
					Message:        fmt.Sprintf("You were assigned to task: %s", task.Title),
					Data: map[string]any{
						"task_id":    task.ID.String(),
						"project_id": projectID.String(),
					},
				})
			}
		}
	}

	res := s.toTaskResponse(task, req.AssigneeIDs)
	return &res, nil
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, apperrors.ErrNotFound
	}

	var assigneeIDs []uuid.UUID
	if s.assigneeRepo != nil {
		assignees, _ := s.assigneeRepo.GetByTaskID(ctx, id)
		for _, a := range assignees {
			assigneeIDs = append(assigneeIDs, a.UserID)
		}
	}

	res := s.toTaskResponse(task, assigneeIDs)
	return &res, nil
}

func (s *TaskService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]dto.TaskResponse, error) {
	tasks, err := s.taskRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	result := make([]dto.TaskResponse, len(tasks))
	for i := range tasks {
		var assigneeIDs []uuid.UUID
		if s.assigneeRepo != nil {
			assignees, _ := s.assigneeRepo.GetByTaskID(ctx, tasks[i].ID)
			for _, a := range assignees {
				assigneeIDs = append(assigneeIDs, a.UserID)
			}
		}
		result[i] = s.toTaskResponse(&tasks[i], assigneeIDs)
	}
	return result, nil
}

func (s *TaskService) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]dto.TaskResponse, error) {
	tasks, err := s.taskRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	result := make([]dto.TaskResponse, len(tasks))
	for i := range tasks {
		var assigneeIDs []uuid.UUID
		if s.assigneeRepo != nil {
			assignees, _ := s.assigneeRepo.GetByTaskID(ctx, tasks[i].ID)
			for _, a := range assignees {
				assigneeIDs = append(assigneeIDs, a.UserID)
			}
		}
		result[i] = s.toTaskResponse(&tasks[i], assigneeIDs)
	}
	return result, nil
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	assigneeIDs, _ := s.GetAssigneeIDs(ctx, id)
	res := s.toTaskResponse(task, assigneeIDs)
	return &res, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) AssignUser(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error {
	if s.assigneeRepo == nil {
		return apperrors.ErrInternal
	}
	return s.assigneeRepo.Assign(ctx, &entities.TaskAssigneeEntity{
		TaskID:     taskID,
		UserID:     userID,
		AssignedAt: time.Now(),
	})
}

func (s *TaskService) UnassignUser(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error {
	if s.assigneeRepo == nil {
		return apperrors.ErrInternal
	}
	return s.assigneeRepo.Unassign(ctx, taskID, userID)
}

func (s *TaskService) GetAssigneeIDs(ctx context.Context, taskID uuid.UUID) ([]uuid.UUID, error) {
	if s.assigneeRepo == nil {
		return nil, nil
	}
	assignees, err := s.assigneeRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task assignees: %w", err)
	}

	ids := make([]uuid.UUID, len(assignees))
	for i, a := range assignees {
		ids[i] = a.UserID
	}
	return ids, nil
}

func (s *TaskService) AddComment(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req dto.CreateTaskCommentRequest) (*dto.TaskCommentResponse, error) {
	if s.commentRepo == nil {
		return nil, apperrors.ErrInternal
	}

	comment := &entities.TaskCommentEntity{
		ID:      uuid.New(),
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return &dto.TaskCommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}, nil
}

func (s *TaskService) GetComments(ctx context.Context, taskID uuid.UUID) ([]dto.TaskCommentResponse, error) {
	if s.commentRepo == nil {
		return nil, nil
	}
	comments, err := s.commentRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	result := make([]dto.TaskCommentResponse, len(comments))
	for i, c := range comments {
		result[i] = dto.TaskCommentResponse{
			ID:        c.ID,
			TaskID:    c.TaskID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}
	return result, nil
}

func (s *TaskService) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	if s.commentRepo == nil {
		return apperrors.ErrInternal
	}
	return s.commentRepo.Delete(ctx, commentID)
}

func (s *TaskService) toTaskResponse(t *entities.TaskEntity, assigneeIDs []uuid.UUID) dto.TaskResponse {
	return dto.TaskResponse{
		ID:             t.ID,
		OrganizationID: t.OrganizationID,
		ProjectID:      t.ProjectID,
		CreatedBy:      t.CreatedBy,
		Title:          t.Title,
		Description:    t.Description,
		Status:         t.Status,
		Priority:       t.Priority,
		DueDate:        t.DueDate,
		AssigneeIDs:    assigneeIDs,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
