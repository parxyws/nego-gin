package domain

import "time"

const TableNameProductMedia = "product_media"

type ProductMedia struct {
	MediaID      int32     `gorm:"column:media_id;primaryKey;autoIncrement:true" json:"media_id"`
	ProductID    string    `gorm:"column:product_id;not null" json:"product_id"`
	MediaType    string    `gorm:"column:media_type;not null" json:"media_type"`
	MediaURL     string    `gorm:"column:media_url;not null" json:"media_url"`
	ThumbnailURL string    `gorm:"column:thumbnail_url" json:"thumbnail_url"`
	AltText      string    `gorm:"column:alt_text" json:"alt_text"`
	SortOrder    int32     `gorm:"column:sort_order" json:"sort_order"`
	IsPrimary    bool      `gorm:"column:is_primary;not null" json:"is_primary"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*ProductMedia) TableName() string {
	return TableNameProductMedia
}
