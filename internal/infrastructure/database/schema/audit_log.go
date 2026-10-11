package schema

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog maps to the "audit_logs" table.
type AuditLog struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index"`
	Action         string         `gorm:"type:varchar(100);not null"`
	ResourceType   string         `gorm:"type:varchar(100);not null"`
	ResourceID     uuid.UUID      `gorm:"type:uuid;not null"`
	Metadata       datatypes.JSON `gorm:"type:jsonb"`
	IPAddress      string         `gorm:"type:inet"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID"`
	User         User         `gorm:"foreignKey:UserID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
