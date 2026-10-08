package irepository

import (
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"

	"context"

	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(user *entities.UserEntity) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.UserEntity, error)
	GetAll()([]entities.UserEntity,error)
}
