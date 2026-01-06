package domain

import (
	"database/sql"
)

const TableNameUserRole = "user_roles"

type UserRole struct {
	UserRoleID int32        `gorm:"column:user_role_id;primaryKey;autoIncrement:true" json:"user_role_id"`
	UserID     string       `gorm:"column:user_id;not null" json:"user_id"`
	RoleID     int32        `gorm:"column:role_id;not null" json:"role_id"`
	AssignedAt sql.NullTime `gorm:"column:assigned_at;not null;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	Role       Role         `gorm:"foreignKey:role_id;references:role_id"`
	User       User         `gorm:"foreignKey:user_id;references:user_id;"`
}

func (*UserRole) TableName() string {
	return TableNameUserRole
}
