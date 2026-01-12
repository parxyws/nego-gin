package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameNotification = "notifications"

// Notification mapped from table <notifications>
type Notification struct {
	NotificationID    string         `gorm:"column:notification_id;primaryKey" json:"notification_id"`
	UserID            string         `gorm:"column:user_id;not null" json:"user_id"`
	NotificationType  string         `gorm:"column:notification_type;not null" json:"notification_type"`
	Title             string         `gorm:"column:title;not null" json:"title"`
	Message           string         `gorm:"column:message;not null" json:"message"`
	RelatedEntityType string         `gorm:"column:related_entity_type" json:"related_entity_type"`
	RelatedEntityID   string         `gorm:"column:related_entity_id" json:"related_entity_id"`
	ActionURL         string         `gorm:"column:action_url" json:"action_url"`
	IsRead            bool           `gorm:"column:is_read;not null" json:"is_read"`
	ReadAt            time.Time      `gorm:"column:read_at" json:"read_at"`
	CreatedAt         time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Notification's table name
func (*Notification) TableName() string {
	return TableNameNotification
}
