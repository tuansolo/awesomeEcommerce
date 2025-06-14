package repository

import (
	"context"

	"awesomeEcommerce/internal/domain"
)

// OrderRepository defines the interface for order repository operations
type OrderRepository interface {
	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, id uint) (*domain.Order, error)

	// FindByUserID retrieves orders by user ID with optional pagination
	FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]domain.Order, int64, error)

	// Create creates a new order
	Create(ctx context.Context, order *domain.Order) error

	// Update updates an existing order
	Update(ctx context.Context, order *domain.Order) error

	// UpdateStatus updates the status of an order
	UpdateStatus(ctx context.Context, id uint, status domain.OrderStatus) error

	// Delete deletes an order by its ID
	Delete(ctx context.Context, id uint) error

	// AddOrderItem adds an item to an order
	AddOrderItem(ctx context.Context, orderItem *domain.OrderItem) error

	// GetOrderItems retrieves all items in an order
	GetOrderItems(ctx context.Context, orderID uint) ([]domain.OrderItem, error)

	// FindAll retrieves all orders with optional pagination
	FindAll(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error)

	// FindByStatus retrieves orders by status with optional pagination
	FindByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error)

	// GetOrderTotal calculates the total price of an order
	GetOrderTotal(ctx context.Context, orderID uint) (float64, error)

	// FindByDateRange retrieves orders created within a date range
	FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error)
}