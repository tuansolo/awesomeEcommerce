package impl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockPaymentRepository is a mock implementation of the PaymentRepository interface
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
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, status, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error) {
	args := m.Called(ctx, method, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Payment), args.Get(1).(int64), args.Error(2)
}

// PaymentRepositoryTestSuite is a test suite for PaymentRepository
type PaymentRepositoryTestSuite struct {
	suite.Suite
	mockRepo repository.PaymentRepository
	ctx      context.Context
}

// SetupTest sets up the test suite
func (s *PaymentRepositoryTestSuite) SetupTest() {
	s.mockRepo = new(MockPaymentRepository)
	s.ctx = context.Background()
}

// TestFindByID tests the FindByID method
func (s *PaymentRepositoryTestSuite) TestFindByID() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a payment by ID
		paymentID := uint(1)
		paymentDate := time.Now()
		expectedPayment := &domain.Payment{
			ID:            paymentID,
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByID", s.ctx, paymentID).Return(expectedPayment, nil).Once()

		// Execute
		payment, err := s.mockRepo.FindByID(s.ctx, paymentID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayment, payment)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Payment not found
		paymentID := uint(999)
		expectedError := errors.New("payment not found")

		mockRepo.On("FindByID", s.ctx, paymentID).Return(nil, expectedError).Once()

		// Execute
		payment, err := s.mockRepo.FindByID(s.ctx, paymentID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payment)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByOrderID tests the FindByOrderID method
func (s *PaymentRepositoryTestSuite) TestFindByOrderID() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a payment by order ID
		orderID := uint(1)
		paymentID := uint(1)
		paymentDate := time.Now()
		expectedPayment := &domain.Payment{
			ID:            paymentID,
			OrderID:       orderID,
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByOrderID", s.ctx, orderID).Return(expectedPayment, nil).Once()

		// Execute
		payment, err := s.mockRepo.FindByOrderID(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayment, payment)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Payment not found
		orderID := uint(999)
		expectedError := errors.New("payment not found")

		mockRepo.On("FindByOrderID", s.ctx, orderID).Return(nil, expectedError).Once()

		// Execute
		payment, err := s.mockRepo.FindByOrderID(s.ctx, orderID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payment)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreate tests the Create method
func (s *PaymentRepositoryTestSuite) TestCreate() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully create a payment
		paymentDate := time.Now()
		payment := &domain.Payment{
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusPending,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
		}

		mockRepo.On("Create", s.ctx, payment).Return(nil).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, payment)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		paymentDate := time.Now()
		payment := &domain.Payment{
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusPending,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
		}
		expectedError := errors.New("database error")

		mockRepo.On("Create", s.ctx, payment).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, payment)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdate tests the Update method
func (s *PaymentRepositoryTestSuite) TestUpdate() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a payment
		paymentDate := time.Now()
		payment := &domain.Payment{
			ID:            uint(1),
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
		}

		mockRepo.On("Update", s.ctx, payment).Return(nil).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, payment)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		paymentDate := time.Now()
		payment := &domain.Payment{
			ID:            uint(1),
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: "txn_123456789",
			PaymentDate:   &paymentDate,
		}
		expectedError := errors.New("database error")

		mockRepo.On("Update", s.ctx, payment).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, payment)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdateStatus tests the UpdateStatus method
func (s *PaymentRepositoryTestSuite) TestUpdateStatus() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a payment status
		paymentID := uint(1)
		newStatus := domain.PaymentStatusCompleted

		mockRepo.On("UpdateStatus", s.ctx, paymentID, newStatus).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateStatus(s.ctx, paymentID, newStatus)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Payment Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Payment not found
		paymentID := uint(999)
		newStatus := domain.PaymentStatusCompleted
		expectedError := errors.New("payment not found")

		mockRepo.On("UpdateStatus", s.ctx, paymentID, newStatus).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateStatus(s.ctx, paymentID, newStatus)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDelete tests the Delete method
func (s *PaymentRepositoryTestSuite) TestDelete() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete a payment
		paymentID := uint(1)

		mockRepo.On("Delete", s.ctx, paymentID).Return(nil).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, paymentID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Payment Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Payment not found
		paymentID := uint(999)
		expectedError := errors.New("payment not found")

		mockRepo.On("Delete", s.ctx, paymentID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, paymentID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByTransactionID tests the FindByTransactionID method
func (s *PaymentRepositoryTestSuite) TestFindByTransactionID() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a payment by transaction ID
		transactionID := "txn_123456789"
		paymentDate := time.Now()
		expectedPayment := &domain.Payment{
			ID:            uint(1),
			OrderID:       uint(1),
			Amount:        100.50,
			Currency:      "USD",
			Method:        domain.PaymentMethodCreditCard,
			Status:        domain.PaymentStatusCompleted,
			TransactionID: transactionID,
			PaymentDate:   &paymentDate,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		mockRepo.On("FindByTransactionID", s.ctx, transactionID).Return(expectedPayment, nil).Once()

		// Execute
		payment, err := s.mockRepo.FindByTransactionID(s.ctx, transactionID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayment, payment)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Payment not found
		transactionID := "txn_nonexistent"
		expectedError := errors.New("payment not found")

		mockRepo.On("FindByTransactionID", s.ctx, transactionID).Return(nil, expectedError).Once()

		// Execute
		payment, err := s.mockRepo.FindByTransactionID(s.ctx, transactionID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payment)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindAll tests the FindAll method
func (s *PaymentRepositoryTestSuite) TestFindAll() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find all payments with pagination
		page := 1
		pageSize := 10
		paymentDate := time.Now()
		expectedPayments := []domain.Payment{
			{
				ID:            uint(1),
				OrderID:       uint(1),
				Amount:        100.50,
				Currency:      "USD",
				Method:        domain.PaymentMethodCreditCard,
				Status:        domain.PaymentStatusCompleted,
				TransactionID: "txn_123456789",
				PaymentDate:   &paymentDate,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uint(2),
				OrderID:       uint(2),
				Amount:        200.75,
				Currency:      "EUR",
				Method:        domain.PaymentMethodPayPal,
				Status:        domain.PaymentStatusPending,
				TransactionID: "txn_987654321",
				PaymentDate:   nil,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayments, payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: No payments found
		page := 1
		pageSize := 10
		var expectedPayments []domain.Payment
		expectedTotal := int64(0)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		payments, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payments)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByStatus tests the FindByStatus method
func (s *PaymentRepositoryTestSuite) TestFindByStatus() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find payments by status with pagination
		status := domain.PaymentStatusCompleted
		page := 1
		pageSize := 10
		paymentDate := time.Now()
		expectedPayments := []domain.Payment{
			{
				ID:            uint(1),
				OrderID:       uint(1),
				Amount:        100.50,
				Currency:      "USD",
				Method:        domain.PaymentMethodCreditCard,
				Status:        status,
				TransactionID: "txn_123456789",
				PaymentDate:   &paymentDate,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uint(3),
				OrderID:       uint(3),
				Amount:        300.25,
				Currency:      "USD",
				Method:        domain.PaymentMethodDebitCard,
				Status:        status,
				TransactionID: "txn_567891234",
				PaymentDate:   &paymentDate,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayments, payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: No payments found with the given status
		status := domain.PaymentStatusRefunded
		page := 1
		pageSize := 10
		var expectedPayments []domain.Payment
		expectedTotal := int64(0)

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		status := domain.PaymentStatusCompleted
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payments)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByDateRange tests the FindByDateRange method
func (s *PaymentRepositoryTestSuite) TestFindByDateRange() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find payments by date range with pagination
		startDate := "2023-01-01"
		endDate := "2023-12-31"
		page := 1
		pageSize := 10
		paymentDate := time.Now()
		expectedPayments := []domain.Payment{
			{
				ID:            uint(1),
				OrderID:       uint(1),
				Amount:        100.50,
				Currency:      "USD",
				Method:        domain.PaymentMethodCreditCard,
				Status:        domain.PaymentStatusCompleted,
				TransactionID: "txn_123456789",
				PaymentDate:   &paymentDate,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uint(2),
				OrderID:       uint(2),
				Amount:        200.75,
				Currency:      "EUR",
				Method:        domain.PaymentMethodPayPal,
				Status:        domain.PaymentStatusPending,
				TransactionID: "txn_987654321",
				PaymentDate:   nil,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayments, payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: No payments found in the given date range
		startDate := "2022-01-01"
		endDate := "2022-12-31"
		page := 1
		pageSize := 10
		var expectedPayments []domain.Payment
		expectedTotal := int64(0)

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		startDate := "2023-01-01"
		endDate := "2023-12-31"
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payments)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByMethod tests the FindByMethod method
func (s *PaymentRepositoryTestSuite) TestFindByMethod() {
	mockRepo := s.mockRepo.(*MockPaymentRepository)

	s.Run("Success", func() {
		// Test case: Successfully find payments by method with pagination
		method := domain.PaymentMethodCreditCard
		page := 1
		pageSize := 10
		paymentDate := time.Now()
		expectedPayments := []domain.Payment{
			{
				ID:            uint(1),
				OrderID:       uint(1),
				Amount:        100.50,
				Currency:      "USD",
				Method:        method,
				Status:        domain.PaymentStatusCompleted,
				TransactionID: "txn_123456789",
				PaymentDate:   &paymentDate,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uint(3),
				OrderID:       uint(3),
				Amount:        300.25,
				Currency:      "USD",
				Method:        method,
				Status:        domain.PaymentStatusPending,
				TransactionID: "txn_567891234",
				PaymentDate:   nil,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByMethod", s.ctx, method, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByMethod(s.ctx, method, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPayments, payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: No payments found with the given method
		method := domain.PaymentMethodBankTransfer
		page := 1
		pageSize := 10
		var expectedPayments []domain.Payment
		expectedTotal := int64(0)

		mockRepo.On("FindByMethod", s.ctx, method, page, pageSize).Return(expectedPayments, expectedTotal, nil).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByMethod(s.ctx, method, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), payments)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockPaymentRepository)

		// Test case: Database error
		method := domain.PaymentMethodCreditCard
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByMethod", s.ctx, method, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		payments, total, err := s.mockRepo.FindByMethod(s.ctx, method, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), payments)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestPaymentRepositorySuite runs the test suite
func TestPaymentRepositorySuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}
