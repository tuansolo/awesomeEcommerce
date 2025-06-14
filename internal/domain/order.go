package domain

import (
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents a customer order
type Order struct {
	ID              uint        `json:"id" gorm:"primaryKey"`
	UserID          uint        `json:"user_id" gorm:"not null"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount     float64     `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status          OrderStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	ShippingAddress string      `json:"shipping_address" gorm:"type:text;not null"`
	BillingAddress  string      `json:"billing_address" gorm:"type:text;not null"`
	PaymentID       *uint       `json:"payment_id"`
	CreatedAt       time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// OrderItem represents an item in a customer order
type OrderItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrderID     uint      `json:"order_id" gorm:"not null"`
	ProductID   uint      `json:"product_id" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"size:255;not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Order
func (Order) TableName() string {
	return "orders"
}

// TableName specifies the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}
