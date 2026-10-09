package schema

import (
	"time"

	"github.com/google/uuid"
)

// Project maps to the "projects" table.
type Project struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;index"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null;index"`
	Name           string     `gorm:"type:varchar(255);not null"`
	Description    string     `gorm:"type:text"`
	Status         string     `gorm:"type:varchar(50);not null;default:'active'"`
	StartDate      *time.Time `gorm:"type:date;default:null"`
	DueDate        *time.Time `gorm:"type:date;default:null"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`

	// Associations
	Organization Organization    `gorm:"foreignKey:OrganizationID"`
	Creator      User            `gorm:"foreignKey:CreatedBy"`
	Members      []ProjectMember `gorm:"foreignKey:ProjectID"`
	Tasks        []Task          `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "projects"
}