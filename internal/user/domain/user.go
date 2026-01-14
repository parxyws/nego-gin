package domain

import (
	"database/sql"
	"time"

	auction "github.com/parxyws/nego-gin/internal/auction/domain"
	product "github.com/parxyws/nego-gin/internal/product/domain"
	"gorm.io/gorm"
)

const TableNameUser = "users"

const (
	Active    = "active"
	Suspended = "suspended"
	Banned    = "banned"
	Pending   = "pending"
)

type User struct {
	UserID        string         `gorm:"column:user_id;primaryKey" json:"user_id"`
	Email         string         `gorm:"column:email;not null" json:"email"`
	PasswordHash  string         `gorm:"column:password_hash;not null" json:"password_hash"`
	Username      string         `gorm:"column:username;not null" json:"username"`
	FirstName     string         `gorm:"column:first_name" json:"first_name"`
	LastName      string         `gorm:"column:last_name" json:"last_name"`
	Phone         string         `gorm:"column:phone" json:"phone"`
	AvatarURL     string         `gorm:"column:avatar_url" json:"avatar_url"`
	AccountStatus string         `gorm:"column:account_status;not null;default:active" json:"account_status"`
	IsVerified    sql.NullTime   `gorm:"column:is_verified" json:"is_verified"`
	LastLogin     sql.NullTime   `gorm:"column:last_login" json:"last_login"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`

	Products      []product.Product `gorm:"foreignKey:seller_id;references:user_id"`
	Auctions      []auction.Auction `gorm:"foreignKey:user_id;references:user_id"`
	UserRoles     []UserRole        `gorm:"foreignKey:user_id;references:user_id"`
	UserAddresses []UserAddress     `gorm:"foreignKey:user_id;references:user_id"`
}

func (u *User) TableName() string {
	return TableNameUser
}
