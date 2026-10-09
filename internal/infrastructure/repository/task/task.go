package taskrepo

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskRepository implements irepository.ITaskRepository.
type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *entities.TaskEntity) error {
	row := toTaskSchema(task)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("TaskRepository.Create: %w", err)
	}
	task.ID = row.ID
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.TaskEntity, error) {
	var row schema.Task
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("TaskRepository.GetByID: %w", err)
	}
	e := toTaskEntity(&row)
	return &e, nil
}

func (r *TaskRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]entities.TaskEntity, error) {
	var rows []schema.Task
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("TaskRepository.GetByProjectID: %w", err)
	}
	return toTaskEntitySlice(rows), nil
}

func (r *TaskRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entities.TaskEntity, error) {
	var rows []schema.Task
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("TaskRepository.GetByOrganizationID: %w", err)
	}
	return toTaskEntitySlice(rows), nil
}

func (r *TaskRepository) Update(ctx context.Context, task *entities.TaskEntity) error {
	row := toTaskSchema(task)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("TaskRepository.Update: %w", err)
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Task{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("TaskRepository.Delete: %w", err)
	}
	return nil
}

// ── TaskAssigneeRepository ────────────────────────────────────────────────────

// TaskAssigneeRepository implements irepository.ITaskAssigneeRepository.
type TaskAssigneeRepository struct {
	db *gorm.DB
}

func NewTaskAssigneeRepository(db *gorm.DB) *TaskAssigneeRepository {
	return &TaskAssigneeRepository{db: db}
}

func (r *TaskAssigneeRepository) Assign(ctx context.Context, assignee *entities.TaskAssigneeEntity) error {
	row := schema.TaskAssignee{TaskID: assignee.TaskID, UserID: assignee.UserID}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("TaskAssigneeRepository.Assign: %w", err)
	}
	return nil
}

func (r *TaskAssigneeRepository) Unassign(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("task_id = ? AND user_id = ?", taskID, userID).
		Delete(&schema.TaskAssignee{}).Error; err != nil {
		return fmt.Errorf("TaskAssigneeRepository.Unassign: %w", err)
	}
	return nil
}

func (r *TaskAssigneeRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskAssigneeEntity, error) {
	var rows []schema.TaskAssignee
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("TaskAssigneeRepository.GetByTaskID: %w", err)
	}
	result := make([]entities.TaskAssigneeEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.TaskAssigneeEntity{TaskID: row.TaskID, UserID: row.UserID, AssignedAt: row.AssignedAt}
	}
	return result, nil
}

func (r *TaskAssigneeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.TaskAssigneeEntity, error) {
	var rows []schema.TaskAssignee
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("TaskAssigneeRepository.GetByUserID: %w", err)
	}
	result := make([]entities.TaskAssigneeEntity, len(rows))
	for i, row := range rows {
		result[i] = entities.TaskAssigneeEntity{TaskID: row.TaskID, UserID: row.UserID, AssignedAt: row.AssignedAt}
	}
	return result, nil
}

// ── TaskCommentRepository ─────────────────────────────────────────────────────

// TaskCommentRepository implements irepository.ITaskCommentRepository.
type TaskCommentRepository struct {
	db *gorm.DB
}

func NewTaskCommentRepository(db *gorm.DB) *TaskCommentRepository {
	return &TaskCommentRepository{db: db}
}

func (r *TaskCommentRepository) Create(ctx context.Context, comment *entities.TaskCommentEntity) error {
	row := toCommentSchema(comment)
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("TaskCommentRepository.Create: %w", err)
	}
	comment.ID = row.ID
	return nil
}

func (r *TaskCommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.TaskCommentEntity, error) {
	var row schema.TaskComment
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("TaskCommentRepository.GetByID: %w", err)
	}
	e := toCommentEntity(&row)
	return &e, nil
}

func (r *TaskCommentRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]entities.TaskCommentEntity, error) {
	var rows []schema.TaskComment
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at asc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("TaskCommentRepository.GetByTaskID: %w", err)
	}
	result := make([]entities.TaskCommentEntity, len(rows))
	for i, row := range rows {
		result[i] = toCommentEntity(&row)
	}
	return result, nil
}

func (r *TaskCommentRepository) Update(ctx context.Context, comment *entities.TaskCommentEntity) error {
	row := toCommentSchema(comment)
	if err := r.db.WithContext(ctx).Save(&row).Error; err != nil {
		return fmt.Errorf("TaskCommentRepository.Update: %w", err)
	}
	return nil
}

func (r *TaskCommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.TaskComment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("TaskCommentRepository.Delete: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toTaskSchema(e *entities.TaskEntity) schema.Task {
	return schema.Task{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		ProjectID:      e.ProjectID,
		CreatedBy:      e.CreatedBy,
		Title:          e.Title,
		Description:    e.Description,
		Status:         e.Status,
		Priority:       e.Priority,
		DueDate:        e.DueDate,
	}
}

func toTaskEntity(s *schema.Task) entities.TaskEntity {
	return entities.TaskEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		ProjectID:      s.ProjectID,
		CreatedBy:      s.CreatedBy,
		Title:          s.Title,
		Description:    s.Description,
		Status:         s.Status,
		Priority:       s.Priority,
		DueDate:        s.DueDate,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func toTaskEntitySlice(rows []schema.Task) []entities.TaskEntity {
	result := make([]entities.TaskEntity, len(rows))
	for i, row := range rows {
		result[i] = toTaskEntity(&row)
	}
	return result
}

func toCommentSchema(e *entities.TaskCommentEntity) schema.TaskComment {
	return schema.TaskComment{ID: e.ID, TaskID: e.TaskID, UserID: e.UserID, Content: e.Content}
}

func toCommentEntity(s *schema.TaskComment) entities.TaskCommentEntity {
	return entities.TaskCommentEntity{
		ID:        s.ID,
		TaskID:    s.TaskID,
		UserID:    s.UserID,
		Content:   s.Content,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
