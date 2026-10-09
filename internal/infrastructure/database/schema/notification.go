package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Notification maps to the "notifications" table.
type Notification struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index"`
	Type           string         `gorm:"type:varchar(100);not null"`
	Title          string         `gorm:"type:varchar(255);not null"`
	Message        string         `gorm:"type:text;not null"`
	Data           datatypes.JSON `gorm:"type:jsonb"`
	ReadAt         *time.Time     `gorm:"default:null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID"`
	User         User         `gorm:"foreignKey:UserID"`
}

func (Notification) TableName() string {
	return "notifications"
}