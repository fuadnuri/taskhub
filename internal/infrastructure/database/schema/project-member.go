package schema

import (
	"time"

	"github.com/google/uuid"
)

// ProjectMember is the join table between Project and User.
// Composite PK: (project_id, user_id).
type ProjectMember struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	JoinedAt  time.Time `gorm:"autoCreateTime"`

	// Associations
	Project Project `gorm:"foreignKey:ProjectID"`
	User    User    `gorm:"foreignKey:UserID"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}
