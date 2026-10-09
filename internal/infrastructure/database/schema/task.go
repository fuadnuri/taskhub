package schema

import (
	"time"

	"github.com/google/uuid"
)

// Task maps to the "tasks" table.
type Task struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;index"`
	ProjectID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title          string     `gorm:"type:varchar(255);not null"`
	Description    string     `gorm:"type:text"`
	Status         string     `gorm:"type:varchar(50);not null;default:'todo'"`
	Priority       string     `gorm:"type:varchar(50);not null;default:'medium'"`
	DueDate        *time.Time `gorm:"default:null"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`

	// Associations
	Organization Organization  `gorm:"foreignKey:OrganizationID"`
	Project      Project       `gorm:"foreignKey:ProjectID"`
	Creator      User          `gorm:"foreignKey:CreatedBy"`
	Assignees    []TaskAssignee `gorm:"foreignKey:TaskID"`
	Comments     []TaskComment  `gorm:"foreignKey:TaskID"`
}

func (Task) TableName() string {
	return "tasks"
}
