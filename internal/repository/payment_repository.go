package repository

import (
	"context"

	"awesomeEcommerce/internal/domain"
)

// PaymentRepository defines the interface for payment repository operations
type PaymentRepository interface {
	// FindByID retrieves a payment by its ID
	FindByID(ctx context.Context, id uint) (*domain.Payment, error)

	// FindByOrderID retrieves a payment by order ID
	FindByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error)

	// Create creates a new payment
	Create(ctx context.Context, payment *domain.Payment) error

	// Update updates an existing payment
	Update(ctx context.Context, payment *domain.Payment) error

	// UpdateStatus updates the status of a payment
	UpdateStatus(ctx context.Context, id uint, status domain.PaymentStatus) error

	// Delete deletes a payment by its ID
	Delete(ctx context.Context, id uint) error

	// FindByTransactionID retrieves a payment by its transaction ID
	FindByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error)

	// FindAll retrieves all payments with optional pagination
	FindAll(ctx context.Context, page, pageSize int) ([]domain.Payment, int64, error)

	// FindByStatus retrieves payments by status with optional pagination
	FindByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error)

	// FindByDateRange retrieves payments created within a date range
	FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error)

	// FindByMethod retrieves payments by payment method with optional pagination
	FindByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error)
}