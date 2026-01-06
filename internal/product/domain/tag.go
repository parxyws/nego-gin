package domain

import "time"

const TableNameTag = "tags"

type Tag struct {
	TagID     int32     `gorm:"column:tag_id;primaryKey;autoIncrement:true" json:"tag_id"`
	TagName   string    `gorm:"column:tag_name;not null" json:"tag_name"`
	Slug      string    `gorm:"column:slug;not null" json:"slug"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*Tag) TableName() string {
	return TableNameTag
}
