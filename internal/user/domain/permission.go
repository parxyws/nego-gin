package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePermission = "permissions"

type Permission struct {
	PermissionID int32          `gorm:"column:permission_id;primaryKey;autoIncrement:true" json:"permission_id"`
	Name         string         `gorm:"column:name;not null" json:"name"`
	Description  string         `gorm:"column:description" json:"description"`
	Resource     string         `gorm:"column:resource;not null" json:"resource"`
	Action       string         `gorm:"column:action;not null" json:"action"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (*Permission) TableName() string {
	return TableNamePermission
}
