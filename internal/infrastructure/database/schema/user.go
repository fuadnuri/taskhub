package schema

import (
	"time"

	"github.com/google/uuid"
)

// User maps to the "users" table.
type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email         string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string     `gorm:"type:text;not null"`
	FirstName     string     `gorm:"type:varchar(100);not null"`
	LastName      string     `gorm:"type:varchar(100);not null"`
	IsActive      bool       `gorm:"not null;default:true"`
	EmailVerified bool       `gorm:"not null;default:false"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`

	// Associations
	RefreshTokens       []RefreshToken       `gorm:"foreignKey:UserID"`
	OrganizationMembers []OrganizationMember `gorm:"foreignKey:UserID"`
	ProjectMembers      []ProjectMember      `gorm:"foreignKey:UserID"`
	TaskAssignees       []TaskAssignee       `gorm:"foreignKey:UserID"`
	TaskComments        []TaskComment        `gorm:"foreignKey:UserID"`
	Notifications       []Notification       `gorm:"foreignKey:UserID"`
	AuditLogs           []AuditLog           `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
