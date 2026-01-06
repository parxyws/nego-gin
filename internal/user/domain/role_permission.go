package domain

import (
	"database/sql"
)

const TableNameRolePermission = "role_permissions"

type RolePermission struct {
	RoleID       int32        `gorm:"column:role_id;primaryKey;autoIncrement:true" json:"role_id"`
	PermissionID int32        `gorm:"column:permission_id;primaryKey;autoIncrement:true" json:"permission_id"`
	AssignedAt   sql.NullTime `gorm:"column:assigned_at;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	AssignedBy   string       `gorm:"column:assigned_by" json:"assigned_by"`
}

func (*RolePermission) TableName() string {
	return TableNameRolePermission
}
