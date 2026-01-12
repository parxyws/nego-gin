package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePayment = "payments"

// Payment mapped from table <payments>
type Payment struct {
	PaymentID        string         `gorm:"column:payment_id;primaryKey" json:"payment_id"`
	OrderID          string         `gorm:"column:order_id;not null" json:"order_id"`
	PaymentMethod    string         `gorm:"column:payment_method;not null" json:"payment_method"`
	TransactionID    string         `gorm:"column:transaction_id" json:"transaction_id"`
	Amount           float64        `gorm:"column:amount;not null" json:"amount"`
	CurrencyCode     string         `gorm:"column:currency_code;not null;default:USD" json:"currency_code"`
	PaymentStatus    string         `gorm:"column:payment_status;not null;default:pending" json:"payment_status"`
	ProviderResponse string         `gorm:"column:provider_response" json:"provider_response"`
	FailureReason    string         `gorm:"column:failure_reason" json:"failure_reason"`
	ProcessedAt      time.Time      `gorm:"column:processed_at" json:"processed_at"`
	CreatedAt        time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Payment's table name
func (*Payment) TableName() string {
	return TableNamePayment
}
