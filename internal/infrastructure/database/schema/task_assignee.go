package schema

import (
	"time"

	"github.com/google/uuid"
)

// TaskAssignee is the join table between Task and User.
// Composite PK: (task_id, user_id).
type TaskAssignee struct {
	TaskID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	AssignedAt time.Time `gorm:"autoCreateTime"`

	// Associations
	Task Task `gorm:"foreignKey:TaskID"`
	User User `gorm:"foreignKey:UserID"`
}

func (TaskAssignee) TableName() string {
	return "task_assignees"
}
