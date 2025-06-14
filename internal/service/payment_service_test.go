package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPaymentRepository is a mock implementation of repository.PaymentRepository
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, id uint) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) UpdateStatus(ctx context.Context, id uint, status domain.PaymentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, status, page, pageSize)
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, method, page, pageSize)
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

// Test cases
func TestCreatePayment(t *testing.T) {
	// Setup
	mockPaymentRepo := new(MockPaymentRepository)
	mockOrderRepo := new(MockOrderRepository)
	paymentService := service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

	ctx := context.Background()
	orderID := uint(1)
	amount := 100.0
	method := domain.PaymentMethodCreditCard

	t.Run("Success", func(t *testing.T) {
		// Test data
		order := &domain.Order{
			ID:          orderID,
			UserID:      1,
			TotalAmount: amount,
			Status:      domain.OrderStatusPending,
		}

		// Expectations
		mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil).Once()
		mockPaymentRepo.On("FindByOrderID", ctx, orderID).Return(nil, errors.New("not found")).Once()
		mockOrderRepo.On("GetOrderTotal", ctx, orderID).Return(amount, nil).Once()
		mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()

		// Execute
		payment, err := paymentService.CreatePayment(ctx, orderID, amount, method)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, orderID, payment.OrderID)
		assert.Equal(t, amount, payment.Amount)
		assert.Equal(t, method, payment.Method)
		assert.Equal(t, domain.PaymentStatusPending, payment.Status)

		mockOrderRepo.AssertExpectations(t)
		mockPaymentRepo.AssertExpectations(t)
	})

	t.Run("Order Not Found", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Expectations
		mockOrderRepo.On("FindByID", ctx, orderID).Return(nil, errors.New("order not found")).Once()

		// Execute
		payment, err := paymentService.CreatePayment(ctx, orderID, amount, method)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Equal(t, "order not found", err.Error())

		mockOrderRepo.AssertExpectations(t)
		mockPaymentRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Payment Already Exists", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		order := &domain.Order{
			ID:          orderID,
			UserID:      1,
			TotalAmount: amount,
			Status:      domain.OrderStatusPending,
		}
		existingPayment := &domain.Payment{
			ID:      1,
			OrderID: orderID,
			Amount:  amount,
		}

		// Expectations
		mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil).Once()
		mockPaymentRepo.On("FindByOrderID", ctx, orderID).Return(existingPayment, nil).Once()

		// Execute
		payment, err := paymentService.CreatePayment(ctx, orderID, amount, method)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Contains(t, err.Error(), "payment already exists")

		mockOrderRepo.AssertExpectations(t)
		mockPaymentRepo.AssertExpectations(t)
	})

	t.Run("Amount Mismatch", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		order := &domain.Order{
			ID:          orderID,
			UserID:      1,
			TotalAmount: 200.0, // Different from payment amount
			Status:      domain.OrderStatusPending,
		}

		// Expectations
		mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil).Once()
		mockPaymentRepo.On("FindByOrderID", ctx, orderID).Return(nil, errors.New("not found")).Once()
		mockOrderRepo.On("GetOrderTotal", ctx, orderID).Return(200.0, nil).Once()

		// Execute
		payment, err := paymentService.CreatePayment(ctx, orderID, amount, method)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Contains(t, err.Error(), "payment amount does not match order total")

		mockOrderRepo.AssertExpectations(t)
		mockPaymentRepo.AssertExpectations(t)
	})

	t.Run("Payment Creation Error", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		order := &domain.Order{
			ID:          orderID,
			UserID:      1,
			TotalAmount: amount,
			Status:      domain.OrderStatusPending,
		}

		// Expectations
		mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil).Once()
		mockPaymentRepo.On("FindByOrderID", ctx, orderID).Return(nil, errors.New("not found")).Once()
		mockOrderRepo.On("GetOrderTotal", ctx, orderID).Return(amount, nil).Once()
		mockPaymentRepo.On("Create", ctx, mock.AnythingOfType("*domain.Payment")).Return(errors.New("failed to create payment")).Once()

		// Execute
		payment, err := paymentService.CreatePayment(ctx, orderID, amount, method)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Equal(t, "failed to create payment", err.Error())

		mockOrderRepo.AssertExpectations(t)
		mockPaymentRepo.AssertExpectations(t)
	})
}

func TestProcessPayment(t *testing.T) {
	// Setup
	mockPaymentRepo := new(MockPaymentRepository)
	mockOrderRepo := new(MockOrderRepository)
	paymentService := service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

	ctx := context.Background()
	paymentID := uint(1)
	transactionID := "txn_123456"

	t.Run("Success", func(t *testing.T) {
		// Test data
		payment := &domain.Payment{
			ID:      paymentID,
			OrderID: 1,
			Amount:  100.0,
			Status:  domain.PaymentStatusPending,
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(payment, nil).Once()
		mockPaymentRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()
		mockOrderRepo.On("UpdateStatus", ctx, payment.OrderID, domain.OrderStatusProcessing).Return(nil).Once()

		// Execute
		err := paymentService.ProcessPayment(ctx, paymentID, transactionID)

		// Assert
		assert.NoError(t, err)

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Payment Not Found", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(nil, errors.New("payment not found")).Once()

		// Execute
		err := paymentService.ProcessPayment(ctx, paymentID, transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "payment not found", err.Error())

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})

	t.Run("Payment Not Pending", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		payment := &domain.Payment{
			ID:      paymentID,
			OrderID: 1,
			Amount:  100.0,
			Status:  domain.PaymentStatusCompleted, // Already completed
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(payment, nil).Once()

		// Execute
		err := paymentService.ProcessPayment(ctx, paymentID, transactionID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in pending status")

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})

	t.Run("Update Error", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		payment := &domain.Payment{
			ID:      paymentID,
			OrderID: 1,
			Amount:  100.0,
			Status:  domain.PaymentStatusPending,
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(payment, nil).Once()
		mockPaymentRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(errors.New("update error")).Once()

		// Execute
		err := paymentService.ProcessPayment(ctx, paymentID, transactionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})
}

func TestGetPaymentByID(t *testing.T) {
	// Setup
	mockPaymentRepo := new(MockPaymentRepository)
	mockOrderRepo := new(MockOrderRepository)
	paymentService := service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

	ctx := context.Background()
	paymentID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		now := time.Now()
		expectedPayment := &domain.Payment{
			ID:            paymentID,
			OrderID:       1,
			Amount:        100.0,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456",
			PaymentDate:   &now,
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(expectedPayment, nil).Once()

		// Execute
		payment, err := paymentService.GetPaymentByID(ctx, paymentID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedPayment, payment)

		mockPaymentRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(nil, errors.New("payment not found")).Once()

		// Execute
		payment, err := paymentService.GetPaymentByID(ctx, paymentID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Equal(t, "payment not found", err.Error())

		mockPaymentRepo.AssertExpectations(t)
	})
}

func TestRefundPayment(t *testing.T) {
	// Setup
	mockPaymentRepo := new(MockPaymentRepository)
	mockOrderRepo := new(MockOrderRepository)
	paymentService := service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

	ctx := context.Background()
	paymentID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		now := time.Now()
		payment := &domain.Payment{
			ID:            paymentID,
			OrderID:       1,
			Amount:        100.0,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456",
			PaymentDate:   &now,
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(payment, nil).Once()
		mockPaymentRepo.On("UpdateStatus", ctx, paymentID, domain.PaymentStatusRefunded).Return(nil).Once()
		mockOrderRepo.On("UpdateStatus", ctx, payment.OrderID, domain.OrderStatusCancelled).Return(nil).Once()

		// Execute
		err := paymentService.RefundPayment(ctx, paymentID)

		// Assert
		assert.NoError(t, err)

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Payment Not Found", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(nil, errors.New("payment not found")).Once()

		// Execute
		err := paymentService.RefundPayment(ctx, paymentID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "payment not found", err.Error())

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})

	t.Run("Not Completed Payment", func(t *testing.T) {
		// Reset mocks
		mockPaymentRepo = new(MockPaymentRepository)
		mockOrderRepo = new(MockOrderRepository)
		paymentService = service.NewPaymentService(mockPaymentRepo, mockOrderRepo, nil)

		// Test data
		now := time.Now()
		payment := &domain.Payment{
			ID:            paymentID,
			OrderID:       1,
			Amount:        100.0,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusPending, // Not completed
			TransactionID: "txn_123456",
			PaymentDate:   &now,
		}

		// Expectations
		mockPaymentRepo.On("FindByID", ctx, paymentID).Return(payment, nil).Once()

		// Execute
		err := paymentService.RefundPayment(ctx, paymentID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only completed payments can be refunded")

		mockPaymentRepo.AssertExpectations(t)
		mockOrderRepo.AssertNotCalled(t, "UpdateStatus")
	})
}
