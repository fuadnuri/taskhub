package notificationrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fuadnuri/taskhub/internal/core/domain/entities"
	"github.com/fuadnuri/taskhub/internal/infrastructure/database/schema"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationRepository implements irepository.INotificationRepository.
type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *entities.NotificationEntity) error {
	row, err := toNotificationSchema(n)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return fmt.Errorf("NotificationRepository.Create: %w", err)
	}
	n.ID = row.ID
	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.NotificationEntity, error) {
	var row schema.Notification
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("NotificationRepository.GetByID: %w", err)
	}
	e, err := toNotificationEntity(&row)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entities.NotificationEntity, error) {
	var rows []schema.Notification
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("NotificationRepository.GetByUserID: %w", err)
	}
	return toNotificationEntitySlice(rows)
}

func (r *NotificationRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]entities.NotificationEntity, error) {
	var rows []schema.Notification
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND read_at IS NULL", userID).
		Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("NotificationRepository.GetUnreadByUserID: %w", err)
	}
	return toNotificationEntitySlice(rows)
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&schema.Notification{}).
		Where("id = ?", id).Update("read_at", "NOW()").Error; err != nil {
		return fmt.Errorf("NotificationRepository.MarkAsRead: %w", err)
	}
	return nil
}

func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&schema.Notification{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("NotificationRepository.Delete: %w", err)
	}
	return nil
}

// ── mappers ──────────────────────────────────────────────────────────────────

func toNotificationSchema(e *entities.NotificationEntity) (schema.Notification, error) {
	var data datatypes.JSON
	if e.Data != nil {
		b, err := json.Marshal(e.Data)
		if err != nil {
			return schema.Notification{}, fmt.Errorf("marshaling notification data: %w", err)
		}
		data = datatypes.JSON(b)
	}
	return schema.Notification{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		UserID:         e.UserID,
		Type:           e.Type,
		Title:          e.Title,
		Message:        e.Message,
		Data:           data,
		ReadAt:         e.ReadAt,
	}, nil
}

func toNotificationEntity(s *schema.Notification) (entities.NotificationEntity, error) {
	var data map[string]any
	if len(s.Data) > 0 {
		if err := json.Unmarshal(s.Data, &data); err != nil {
			return entities.NotificationEntity{}, fmt.Errorf("unmarshaling notification data: %w", err)
		}
	}
	return entities.NotificationEntity{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		UserID:         s.UserID,
		Type:           s.Type,
		Title:          s.Title,
		Message:        s.Message,
		Data:           data,
		ReadAt:         s.ReadAt,
		CreatedAt:      s.CreatedAt,
	}, nil
}

func toNotificationEntitySlice(rows []schema.Notification) ([]entities.NotificationEntity, error) {
	result := make([]entities.NotificationEntity, 0, len(rows))
	for _, row := range rows {
		e, err := toNotificationEntity(&row)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}
