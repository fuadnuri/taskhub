package schema

import (
	"time"

	"github.com/google/uuid"
)

// Role maps to the "roles" table.
type Role struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(100);not null"`
	Description    string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	// Associations
	Organization        Organization         `gorm:"foreignKey:OrganizationID"`
	RolePermissions     []RolePermission     `gorm:"foreignKey:RoleID"`
	OrganizationMembers []OrganizationMember `gorm:"foreignKey:RoleID"`
	Invitations         []OrganizationInvitation `gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return "roles"
}
