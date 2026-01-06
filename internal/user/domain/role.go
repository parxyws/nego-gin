package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameRole = "roles"

type Role struct {
	RoleID      int32          `gorm:"column:role_id;primaryKey;autoIncrement:true" json:"role_id"`
	RoleName    string         `gorm:"column:role_name;not null" json:"role_name"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (*Role) TableName() string {
	return TableNameRole
}
