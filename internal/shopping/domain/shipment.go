package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameShipment = "shipments"

// Shipment mapped from table <shipments>
type Shipment struct {
	ShipmentID        int32          `gorm:"column:shipment_id;primaryKey;autoIncrement:true" json:"shipment_id"`
	OrderID           string         `gorm:"column:order_id;not null" json:"order_id"`
	Carrier           string         `gorm:"column:carrier" json:"carrier"`
	TrackingNumber    string         `gorm:"column:tracking_number" json:"tracking_number"`
	ShippingMethod    string         `gorm:"column:shipping_method" json:"shipping_method"`
	EstimatedDelivery time.Time      `gorm:"column:estimated_delivery" json:"estimated_delivery"`
	ActualDelivery    time.Time      `gorm:"column:actual_delivery" json:"actual_delivery"`
	ShipmentStatus    string         `gorm:"column:shipment_status;not null;default:pending" json:"shipment_status"`
	ShippedAt         time.Time      `gorm:"column:shipped_at" json:"shipped_at"`
	CreatedAt         time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Shipment's table name
func (*Shipment) TableName() string {
	return TableNameShipment
}
