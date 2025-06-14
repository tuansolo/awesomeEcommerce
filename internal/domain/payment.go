package domain

import (
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodDebitCard    PaymentMethod = "debit_card"
	PaymentMethodPayPal       PaymentMethod = "paypal"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// Payment represents a payment transaction
type Payment struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	OrderID       uint          `json:"order_id" gorm:"not null"`
	Amount        float64       `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency      string        `json:"currency" gorm:"size:3;not null;default:'USD'"`
	Method        PaymentMethod `json:"method" gorm:"type:varchar(20);not null"`
	Status        PaymentStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TransactionID string        `json:"transaction_id" gorm:"size:100"`
	PaymentDate   *time.Time    `json:"payment_date"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Payment
func (Payment) TableName() string {
	return "payments"
}
