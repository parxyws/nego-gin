package domain

type Auction struct {
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"column:user_id"`
}

func (a *Auction) TableName() string {
	return "auctions"
}
