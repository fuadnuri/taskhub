package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct { 
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	CreatedAt    time.Time
}

type OrganizationEntity struct {
	ID        uuid.UUID
	Name      string
	slug      string
	createdAt time.Time
}

type MembershipEntity struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	CreatedAt      time.Time
}

type ProjectsEntity struct {
	ID             uuid.UUID
	Name           string
	OrganizationID uuid.UUID
	Description    string
	CreatedAt      time.Time
	CreatedBy      uuid.UUID
}

type TasksEntity struct {
	ID              uuid.UUID
	Name            string
	ProjectID       uuid.UUID
	Description     string
	CreatedAt       time.Time
	CreatedBy       uuid.UUID
	Status          string
	DueDate         time.Time
	AssigneeID      uuid.UUID
	AssigneeName    string
	AssigneeEmail   string
	AssigneePhone   string
	AssigneeAddress string
	AssigneeCity    string
	AssigneeState   string
	AssigneeZip     string
	AssigneeCountry string
}
