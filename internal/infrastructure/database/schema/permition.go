package schema

import "github.com/google/uuid"

// Permission maps to the "permissions" table.
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `gorm:"type:text"`

	// Associations
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID"`
}

func (Permission) TableName() string {
	return "permissions"
}
