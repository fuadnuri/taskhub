package service

import (
	"context"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/application/apperrors"
	"github.com/fuadnuri/taskhub/internal/core/application/dto"
	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/core/domain/irepository"
	"github.com/google/uuid"
)

// UserService implements IUserService.
type UserService struct {
	repo            irepository.IUserRepository
	passwordService IPasswordService
}

func NewUserService(repo irepository.IUserRepository, passwordService IPasswordService) *UserService {
	return &UserService{
		repo:            repo,
		passwordService: passwordService,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	existing, _ := s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	hash, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &entities.UserEntity{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  hash,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		IsActive:      true,
		EmailVerified: false,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	res := toUserResponse(user)
	return &res, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toUserResponse(user)
	return &res, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toUserResponse(user)
	return &res, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	result := make([]dto.UserResponse, len(users))
	for i := range users {
		result[i] = toUserResponse(&users[i])
	}
	return result, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, apperrors.ErrNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	res := toUserResponse(user)
	return &res, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil || user == nil {
		return apperrors.ErrNotFound
	}

	if err := s.passwordService.Verify(req.OldPassword, user.PasswordHash); err != nil {
		return apperrors.ErrInvalidCredentials
	}

	newHash, err := s.passwordService.Hash(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	user.PasswordHash = newHash
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}