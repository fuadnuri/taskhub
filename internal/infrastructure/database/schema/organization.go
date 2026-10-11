package schema

import (
	"time"

	"github.com/google/uuid"
)

// Organization maps to the "organizations" table.
type Organization struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Slug      string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	IsActive  bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Associations
	Members     []OrganizationMember     `gorm:"foreignKey:OrganizationID"`
	Roles       []Role                   `gorm:"foreignKey:OrganizationID"`
	Invitations []OrganizationInvitation `gorm:"foreignKey:OrganizationID"`
	Projects    []Project                `gorm:"foreignKey:OrganizationID"`
	Tasks       []Task                   `gorm:"foreignKey:OrganizationID"`
	Notifications []Notification         `gorm:"foreignKey:OrganizationID"`
	AuditLogs   []AuditLog               `gorm:"foreignKey:OrganizationID"`
}

func (Organization) TableName() string {
	return "organizations"
}

// OrganizationMember maps to the "organization_members" join table.
type OrganizationMember struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID         uuid.UUID `gorm:"type:uuid;not null;index"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID"`
	User         User         `gorm:"foreignKey:UserID"`
	Role         Role         `gorm:"foreignKey:RoleID"`
}

func (OrganizationMember) TableName() string {
	return "organization_members"
}

// OrganizationInvitation maps to the "organization_invitations" table.
type OrganizationInvitation struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Email          string     `gorm:"type:varchar(255);not null"`
	RoleID         uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash      string     `gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt      time.Time  `gorm:"not null"`
	AcceptedAt     *time.Time `gorm:"default:null"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`

	// Associations
	Organization Organization `gorm:"foreignKey:OrganizationID"`
	Role         Role         `gorm:"foreignKey:RoleID"`
}

func (OrganizationInvitation) TableName() string {
	return "organization_invitations"
}
