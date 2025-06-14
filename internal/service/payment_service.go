package service

import (
	"context"
	"errors"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/messaging"
	"awesomeEcommerce/internal/repository"
)

// PaymentService defines the interface for payment-related business logic
type PaymentService interface {
	// GetPaymentByID retrieves a payment by its ID
	GetPaymentByID(ctx context.Context, id uint) (*domain.Payment, error)

	// GetPaymentByOrderID retrieves a payment by order ID
	GetPaymentByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error)

	// CreatePayment creates a new payment for an order
	CreatePayment(ctx context.Context, orderID uint, amount float64, method domain.PaymentMethod) (*domain.Payment, error)

	// ProcessPayment processes a payment (simulated)
	ProcessPayment(ctx context.Context, paymentID uint, transactionID string) error

	// UpdatePaymentStatus updates the status of a payment
	UpdatePaymentStatus(ctx context.Context, id uint, status domain.PaymentStatus) error

	// RefundPayment refunds a payment
	RefundPayment(ctx context.Context, id uint) error

	// GetAllPayments retrieves all payments with optional pagination
	GetAllPayments(ctx context.Context, page, pageSize int) ([]domain.Payment, int64, error)

	// GetPaymentsByStatus retrieves payments by status with optional pagination
	GetPaymentsByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error)

	// GetPaymentsByDateRange retrieves payments created within a date range
	GetPaymentsByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error)

	// GetPaymentsByMethod retrieves payments by payment method with optional pagination
	GetPaymentsByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error)
}

// PaymentServiceImpl implements the PaymentService interface
type PaymentServiceImpl struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	producer    *messaging.KafkaProducer
}

// NewPaymentService creates a new PaymentServiceImpl
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	producer *messaging.KafkaProducer,
) PaymentService {
	return &PaymentServiceImpl{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		producer:    producer,
	}
}

// GetPaymentByID retrieves a payment by its ID
func (s *PaymentServiceImpl) GetPaymentByID(ctx context.Context, id uint) (*domain.Payment, error) {
	return s.paymentRepo.FindByID(ctx, id)
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentServiceImpl) GetPaymentByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error) {
	return s.paymentRepo.FindByOrderID(ctx, orderID)
}

// CreatePayment creates a new payment for an order
func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, orderID uint, amount float64, method domain.PaymentMethod) (*domain.Payment, error) {
	// Check if order exists
	_, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if payment already exists for this order
	_, err = s.paymentRepo.FindByOrderID(ctx, orderID)
	if err == nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Validate amount
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Validate that amount matches order total
	orderTotal, err := s.orderRepo.GetOrderTotal(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if amount != orderTotal {
		return nil, errors.New("payment amount does not match order total")
	}

	// Create the payment
	payment := &domain.Payment{
		OrderID:  orderID,
		Amount:   amount,
		Currency: "USD", // Default currency
		Method:   method,
		Status:   domain.PaymentStatusPending,
	}

	// Save the payment
	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Publish payment created event
	// Note: In a real application, we would serialize the payment to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "payment-created", []byte(fmt.Sprintf("%d", payment.ID)), []byte(paymentJSON))

	return payment, nil
}

// ProcessPayment processes a payment (simulated)
func (s *PaymentServiceImpl) ProcessPayment(ctx context.Context, paymentID uint, transactionID string) error {
	// Check if payment exists
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Check if payment is in pending status
	if payment.Status != domain.PaymentStatusPending {
		return errors.New("payment is not in pending status")
	}

	// Update payment with transaction ID and completed status
	payment.TransactionID = transactionID
	payment.Status = domain.PaymentStatusCompleted
	payment.PaymentDate = func() *time.Time { t := time.Now(); return &t }()

	// Save the payment
	err = s.paymentRepo.Update(ctx, payment)
	if err != nil {
		return err
	}

	// Update order status to processing
	err = s.orderRepo.UpdateStatus(ctx, payment.OrderID, domain.OrderStatusProcessing)
	if err != nil {
		return err
	}

	// Publish payment processed event
	// Note: In a real application, we would serialize the payment to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "payment-processed", []byte(fmt.Sprintf("%d", payment.ID)), []byte(paymentJSON))

	return nil
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentServiceImpl) UpdatePaymentStatus(ctx context.Context, id uint, status domain.PaymentStatus) error {
	// Check if payment exists
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("payment not found")
	}

	// Validate status transition
	if !s.isValidStatusTransition(payment.Status, status) {
		return errors.New("invalid status transition")
	}

	// Update the status
	err = s.paymentRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// If payment is completed, update order status to processing
	if status == domain.PaymentStatusCompleted {
		err = s.orderRepo.UpdateStatus(ctx, payment.OrderID, domain.OrderStatusProcessing)
		if err != nil {
			return err
		}
	}

	// If payment is failed, update order status to cancelled
	if status == domain.PaymentStatusFailed {
		err = s.orderRepo.UpdateStatus(ctx, payment.OrderID, domain.OrderStatusCancelled)
		if err != nil {
			return err
		}
	}

	// Publish payment status updated event
	// Note: In a real application, we would serialize the payment to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "payment-status-updated", []byte(fmt.Sprintf("%d", id)), []byte(paymentJSON))

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func (s *PaymentServiceImpl) isValidStatusTransition(from, to domain.PaymentStatus) bool {
	// Define valid transitions
	validTransitions := map[domain.PaymentStatus][]domain.PaymentStatus{
		domain.PaymentStatusPending: {
			domain.PaymentStatusCompleted,
			domain.PaymentStatusFailed,
		},
		domain.PaymentStatusCompleted: {
			domain.PaymentStatusRefunded,
		},
		domain.PaymentStatusFailed:    {},
		domain.PaymentStatusRefunded:  {},
	}

	// Check if the transition is valid
	for _, validTo := range validTransitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// RefundPayment refunds a payment
func (s *PaymentServiceImpl) RefundPayment(ctx context.Context, id uint) error {
	// Check if payment exists
	payment, err := s.paymentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("payment not found")
	}

	// Check if payment can be refunded
	if payment.Status != domain.PaymentStatusCompleted {
		return errors.New("only completed payments can be refunded")
	}

	// Update payment status to refunded
	err = s.paymentRepo.UpdateStatus(ctx, id, domain.PaymentStatusRefunded)
	if err != nil {
		return err
	}

	// Update order status to cancelled
	err = s.orderRepo.UpdateStatus(ctx, payment.OrderID, domain.OrderStatusCancelled)
	if err != nil {
		return err
	}

	// Get order items to return to inventory
	// Note: In a real implementation, we would get the order items and return them to inventory
	// For simplicity, we're skipping this step

	// Publish payment refunded event
	// Note: In a real application, we would serialize the payment to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "payment-refunded", []byte(fmt.Sprintf("%d", id)), []byte(paymentJSON))

	return nil
}

// GetAllPayments retrieves all payments with optional pagination
func (s *PaymentServiceImpl) GetAllPayments(ctx context.Context, page, pageSize int) ([]domain.Payment, int64, error) {
	return s.paymentRepo.FindAll(ctx, page, pageSize)
}

// GetPaymentsByStatus retrieves payments by status with optional pagination
func (s *PaymentServiceImpl) GetPaymentsByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error) {
	return s.paymentRepo.FindByStatus(ctx, status, page, pageSize)
}

// GetPaymentsByDateRange retrieves payments created within a date range
func (s *PaymentServiceImpl) GetPaymentsByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error) {
	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, errors.New("invalid start date format, use YYYY-MM-DD")
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, errors.New("invalid end date format, use YYYY-MM-DD")
	}

	return s.paymentRepo.FindByDateRange(ctx, startDate, endDate, page, pageSize)
}

// GetPaymentsByMethod retrieves payments by payment method with optional pagination
func (s *PaymentServiceImpl) GetPaymentsByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error) {
	return s.paymentRepo.FindByMethod(ctx, method, page, pageSize)
}
