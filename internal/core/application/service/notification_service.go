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

// NotificationService implements INotificationService.
type NotificationService struct {
	notificationRepo irepository.INotificationRepository
}

func NewNotificationService(notificationRepo irepository.INotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) Create(ctx context.Context, req dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	n := &entities.NotificationEntity{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		UserID:         req.UserID,
		Type:           req.Type,
		Title:          req.Title,
		Message:        req.Message,
		Data:           req.Data,
	}

	if err := s.notificationRepo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	res := toNotificationResponse(n)
	return &res, nil
}

func (s *NotificationService) GetByID(ctx context.Context, id uuid.UUID) (*dto.NotificationResponse, error) {
	n, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil || n == nil {
		return nil, apperrors.ErrNotFound
	}
	res := toNotificationResponse(n)
	return &res, nil
}

func (s *NotificationService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]dto.NotificationResponse, error) {
	notifications, err := s.notificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	result := make([]dto.NotificationResponse, len(notifications))
	for i := range notifications {
		result[i] = toNotificationResponse(&notifications[i])
	}
	return result, nil
}

func (s *NotificationService) GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]dto.NotificationResponse, error) {
	notifications, err := s.notificationRepo.GetUnreadByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}

	result := make([]dto.NotificationResponse, len(notifications))
	for i := range notifications {
		result[i] = toNotificationResponse(&notifications[i])
	}
	return result, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *NotificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.notificationRepo.Delete(ctx, id)
}

func toNotificationResponse(n *entities.NotificationEntity) dto.NotificationResponse {
	return dto.NotificationResponse{
		ID:             n.ID,
		OrganizationID: n.OrganizationID,
		UserID:         n.UserID,
		Type:           n.Type,
		Title:          n.Title,
		Message:        n.Message,
		Data:           n.Data,
		ReadAt:         n.ReadAt,
		CreatedAt:      n.CreatedAt,
	}
}
