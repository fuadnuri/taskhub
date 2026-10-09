package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string      `json:"title" binding:"required"`
	Description string      `json:"description"`
	Status      string      `json:"status"`   // "todo", "in_progress", "done"
	Priority    string      `json:"priority"` // "low", "medium", "high", "urgent"
	DueDate     *time.Time  `json:"due_date,omitempty"`
	AssigneeIDs []uuid.UUID `json:"assignee_ids,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskResponse struct {
	ID             uuid.UUID   `json:"id"`
	OrganizationID uuid.UUID   `json:"organization_id"`
	ProjectID      uuid.UUID   `json:"project_id"`
	CreatedBy      uuid.UUID   `json:"created_by"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Status         string      `json:"status"`
	Priority       string      `json:"priority"`
	DueDate        *time.Time  `json:"due_date,omitempty"`
	AssigneeIDs    []uuid.UUID `json:"assignee_ids,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type CreateTaskCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateTaskCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type TaskCommentResponse struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
