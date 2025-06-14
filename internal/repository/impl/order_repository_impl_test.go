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

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uint, status domain.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) AddOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	args := m.Called(ctx, orderItem)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderItems(ctx context.Context, orderID uint) ([]domain.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.OrderItem), args.Error(1)
}

func (m *MockOrderRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, status, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetOrderTotal(ctx context.Context, orderID uint) (float64, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderRepository) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

// OrderRepositoryTestSuite is a test suite for OrderRepository
type OrderRepositoryTestSuite struct {
	suite.Suite
	mockRepo repository.OrderRepository
	ctx      context.Context
}

// SetupTest sets up the test suite
func (s *OrderRepositoryTestSuite) SetupTest() {
	s.mockRepo = new(MockOrderRepository)
	s.ctx = context.Background()
}

// TestFindByID tests the FindByID method
func (s *OrderRepositoryTestSuite) TestFindByID() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully find an order by ID
		orderID := uint(1)
		expectedOrder := &domain.Order{
			ID:              orderID,
			UserID:          uint(1),
			TotalAmount:     99.99,
			Status:          domain.OrderStatusPending,
			ShippingAddress: "123 Shipping St",
			BillingAddress:  "123 Billing St",
			Items: []domain.OrderItem{
				{
					ID:          uint(1),
					OrderID:     orderID,
					ProductID:   uint(2),
					ProductName: "Test Product",
					Price:       49.99,
					Quantity:    2,
				},
			},
		}

		mockRepo.On("FindByID", s.ctx, orderID).Return(expectedOrder, nil).Once()

		// Execute
		order, err := s.mockRepo.FindByID(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrder, order)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderID := uint(999)
		expectedError := errors.New("order not found")

		mockRepo.On("FindByID", s.ctx, orderID).Return(nil, expectedError).Once()

		// Execute
		order, err := s.mockRepo.FindByID(s.ctx, orderID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), order)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByUserID tests the FindByUserID method
func (s *OrderRepositoryTestSuite) TestFindByUserID() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully find orders by user ID
		userID := uint(1)
		page := 1
		pageSize := 10
		expectedOrders := []domain.Order{
			{
				ID:              uint(1),
				UserID:          userID,
				TotalAmount:     99.99,
				Status:          domain.OrderStatusPending,
				ShippingAddress: "123 Shipping St",
				BillingAddress:  "123 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(1),
						OrderID:     uint(1),
						ProductID:   uint(2),
						ProductName: "Test Product 1",
						Price:       49.99,
						Quantity:    2,
					},
				},
			},
			{
				ID:              uint(2),
				UserID:          userID,
				TotalAmount:     29.99,
				Status:          domain.OrderStatusDelivered,
				ShippingAddress: "123 Shipping St",
				BillingAddress:  "123 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(2),
						OrderID:     uint(2),
						ProductID:   uint(3),
						ProductName: "Test Product 2",
						Price:       29.99,
						Quantity:    1,
					},
				},
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByUserID", s.ctx, userID, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByUserID(s.ctx, userID, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrders, orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Orders", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: User has no orders
		userID := uint(2)
		page := 1
		pageSize := 10
		var expectedOrders []domain.Order
		expectedTotal := int64(0)

		mockRepo.On("FindByUserID", s.ctx, userID, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByUserID(s.ctx, userID, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Database error
		userID := uint(1)
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByUserID", s.ctx, userID, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByUserID(s.ctx, userID, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreate tests the Create method
func (s *OrderRepositoryTestSuite) TestCreate() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully create an order
		order := &domain.Order{
			UserID:          uint(1),
			TotalAmount:     99.99,
			Status:          domain.OrderStatusPending,
			ShippingAddress: "123 Shipping St",
			BillingAddress:  "123 Billing St",
			Items: []domain.OrderItem{
				{
					ProductID:   uint(2),
					ProductName: "Test Product",
					Price:       49.99,
					Quantity:    2,
				},
			},
		}

		mockRepo.On("Create", s.ctx, order).Return(nil).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, order)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Database error
		order := &domain.Order{
			UserID:          uint(1),
			TotalAmount:     99.99,
			Status:          domain.OrderStatusPending,
			ShippingAddress: "123 Shipping St",
			BillingAddress:  "123 Billing St",
			Items: []domain.OrderItem{
				{
					ProductID:   uint(2),
					ProductName: "Test Product",
					Price:       49.99,
					Quantity:    2,
				},
			},
		}
		expectedError := errors.New("database error")

		mockRepo.On("Create", s.ctx, order).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, order)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdate tests the Update method
func (s *OrderRepositoryTestSuite) TestUpdate() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully update an order
		order := &domain.Order{
			ID:              uint(1),
			UserID:          uint(1),
			TotalAmount:     129.99,
			Status:          domain.OrderStatusProcessing,
			ShippingAddress: "123 Shipping St",
			BillingAddress:  "123 Billing St",
		}

		mockRepo.On("Update", s.ctx, order).Return(nil).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, order)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		order := &domain.Order{
			ID:              uint(999),
			UserID:          uint(1),
			TotalAmount:     129.99,
			Status:          domain.OrderStatusProcessing,
			ShippingAddress: "123 Shipping St",
			BillingAddress:  "123 Billing St",
		}
		expectedError := errors.New("order not found")

		mockRepo.On("Update", s.ctx, order).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, order)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdateStatus tests the UpdateStatus method
func (s *OrderRepositoryTestSuite) TestUpdateStatus() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully update an order status
		orderID := uint(1)
		newStatus := domain.OrderStatusShipped

		mockRepo.On("UpdateStatus", s.ctx, orderID, newStatus).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateStatus(s.ctx, orderID, newStatus)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderID := uint(999)
		newStatus := domain.OrderStatusShipped
		expectedError := errors.New("order not found")

		mockRepo.On("UpdateStatus", s.ctx, orderID, newStatus).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateStatus(s.ctx, orderID, newStatus)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDelete tests the Delete method
func (s *OrderRepositoryTestSuite) TestDelete() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete an order
		orderID := uint(1)

		mockRepo.On("Delete", s.ctx, orderID).Return(nil).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderID := uint(999)
		expectedError := errors.New("order not found")

		mockRepo.On("Delete", s.ctx, orderID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, orderID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestAddOrderItem tests the AddOrderItem method
func (s *OrderRepositoryTestSuite) TestAddOrderItem() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully add an item to an order
		orderItem := &domain.OrderItem{
			OrderID:     uint(1),
			ProductID:   uint(3),
			ProductName: "New Product",
			Price:       19.99,
			Quantity:    1,
		}

		mockRepo.On("AddOrderItem", s.ctx, orderItem).Return(nil).Once()

		// Execute
		err := s.mockRepo.AddOrderItem(s.ctx, orderItem)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderItem := &domain.OrderItem{
			OrderID:     uint(999),
			ProductID:   uint(3),
			ProductName: "New Product",
			Price:       19.99,
			Quantity:    1,
		}
		expectedError := errors.New("order not found")

		mockRepo.On("AddOrderItem", s.ctx, orderItem).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.AddOrderItem(s.ctx, orderItem)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetOrderItems tests the GetOrderItems method
func (s *OrderRepositoryTestSuite) TestGetOrderItems() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully get items from an order
		orderID := uint(1)
		expectedItems := []domain.OrderItem{
			{
				ID:          uint(1),
				OrderID:     orderID,
				ProductID:   uint(2),
				ProductName: "Test Product 1",
				Price:       49.99,
				Quantity:    2,
			},
			{
				ID:          uint(2),
				OrderID:     orderID,
				ProductID:   uint(3),
				ProductName: "Test Product 2",
				Price:       29.99,
				Quantity:    1,
			},
		}

		mockRepo.On("GetOrderItems", s.ctx, orderID).Return(expectedItems, nil).Once()

		// Execute
		items, err := s.mockRepo.GetOrderItems(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedItems, items)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Order", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order has no items
		orderID := uint(2)
		var expectedItems []domain.OrderItem

		mockRepo.On("GetOrderItems", s.ctx, orderID).Return(expectedItems, nil).Once()

		// Execute
		items, err := s.mockRepo.GetOrderItems(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), items)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderID := uint(999)
		expectedError := errors.New("order not found")

		mockRepo.On("GetOrderItems", s.ctx, orderID).Return(nil, expectedError).Once()

		// Execute
		items, err := s.mockRepo.GetOrderItems(s.ctx, orderID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), items)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindAll tests the FindAll method
func (s *OrderRepositoryTestSuite) TestFindAll() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully find all orders
		page := 1
		pageSize := 10
		expectedOrders := []domain.Order{
			{
				ID:              uint(1),
				UserID:          uint(1),
				TotalAmount:     99.99,
				Status:          domain.OrderStatusPending,
				ShippingAddress: "123 Shipping St",
				BillingAddress:  "123 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(1),
						OrderID:     uint(1),
						ProductID:   uint(2),
						ProductName: "Test Product 1",
						Price:       49.99,
						Quantity:    2,
					},
				},
			},
			{
				ID:              uint(2),
				UserID:          uint(2),
				TotalAmount:     29.99,
				Status:          domain.OrderStatusDelivered,
				ShippingAddress: "456 Shipping St",
				BillingAddress:  "456 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(2),
						OrderID:     uint(2),
						ProductID:   uint(3),
						ProductName: "Test Product 2",
						Price:       29.99,
						Quantity:    1,
					},
				},
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrders, orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Orders", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: No orders in the system
		page := 1
		pageSize := 10
		var expectedOrders []domain.Order
		expectedTotal := int64(0)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Database error
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		orders, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByStatus tests the FindByStatus method
func (s *OrderRepositoryTestSuite) TestFindByStatus() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully find orders by status
		status := domain.OrderStatusPending
		page := 1
		pageSize := 10
		expectedOrders := []domain.Order{
			{
				ID:              uint(1),
				UserID:          uint(1),
				TotalAmount:     99.99,
				Status:          status,
				ShippingAddress: "123 Shipping St",
				BillingAddress:  "123 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(1),
						OrderID:     uint(1),
						ProductID:   uint(2),
						ProductName: "Test Product 1",
						Price:       49.99,
						Quantity:    2,
					},
				},
			},
			{
				ID:              uint(3),
				UserID:          uint(3),
				TotalAmount:     79.99,
				Status:          status,
				ShippingAddress: "789 Shipping St",
				BillingAddress:  "789 Billing St",
				Items: []domain.OrderItem{
					{
						ID:          uint(3),
						OrderID:     uint(3),
						ProductID:   uint(4),
						ProductName: "Test Product 3",
						Price:       79.99,
						Quantity:    1,
					},
				},
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrders, orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Orders With Status", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: No orders with the given status
		status := domain.OrderStatusCancelled
		page := 1
		pageSize := 10
		var expectedOrders []domain.Order
		expectedTotal := int64(0)

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Database error
		status := domain.OrderStatusPending
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByStatus", s.ctx, status, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByStatus(s.ctx, status, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetOrderTotal tests the GetOrderTotal method
func (s *OrderRepositoryTestSuite) TestGetOrderTotal() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully calculate the total price of an order
		orderID := uint(1)
		expectedTotal := 129.97 // (49.99 * 2) + (29.99 * 1) = 129.97

		mockRepo.On("GetOrderTotal", s.ctx, orderID).Return(expectedTotal, nil).Once()

		// Execute
		total, err := s.mockRepo.GetOrderTotal(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Order", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order has no items
		orderID := uint(2)
		expectedTotal := 0.0

		mockRepo.On("GetOrderTotal", s.ctx, orderID).Return(expectedTotal, nil).Once()

		// Execute
		total, err := s.mockRepo.GetOrderTotal(s.ctx, orderID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Order Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Order not found
		orderID := uint(999)
		expectedError := errors.New("order not found")

		mockRepo.On("GetOrderTotal", s.ctx, orderID).Return(0.0, expectedError).Once()

		// Execute
		total, err := s.mockRepo.GetOrderTotal(s.ctx, orderID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), 0.0, total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByDateRange tests the FindByDateRange method
func (s *OrderRepositoryTestSuite) TestFindByDateRange() {
	mockRepo := s.mockRepo.(*MockOrderRepository)

	s.Run("Success", func() {
		// Test case: Successfully find orders by date range
		startDate := "2023-01-01"
		endDate := "2023-01-31"
		page := 1
		pageSize := 10
		expectedOrders := []domain.Order{
			{
				ID:              uint(1),
				UserID:          uint(1),
				TotalAmount:     99.99,
				Status:          domain.OrderStatusPending,
				ShippingAddress: "123 Shipping St",
				BillingAddress:  "123 Billing St",
				CreatedAt:       time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC),
				Items: []domain.OrderItem{
					{
						ID:          uint(1),
						OrderID:     uint(1),
						ProductID:   uint(2),
						ProductName: "Test Product 1",
						Price:       49.99,
						Quantity:    2,
					},
				},
			},
			{
				ID:              uint(2),
				UserID:          uint(2),
				TotalAmount:     29.99,
				Status:          domain.OrderStatusDelivered,
				ShippingAddress: "456 Shipping St",
				BillingAddress:  "456 Billing St",
				CreatedAt:       time.Date(2023, 1, 20, 14, 30, 0, 0, time.UTC),
				Items: []domain.OrderItem{
					{
						ID:          uint(2),
						OrderID:     uint(2),
						ProductID:   uint(3),
						ProductName: "Test Product 2",
						Price:       29.99,
						Quantity:    1,
					},
				},
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrders, orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Orders In Date Range", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: No orders in the given date range
		startDate := "2022-01-01"
		endDate := "2022-01-31"
		page := 1
		pageSize := 10
		var expectedOrders []domain.Order
		expectedTotal := int64(0)

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(expectedOrders, expectedTotal, nil).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), orders)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Invalid Date Format", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Invalid date format
		startDate := "invalid-date"
		endDate := "2023-01-31"
		page := 1
		pageSize := 10
		expectedError := errors.New("invalid date format")

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockOrderRepository)

		// Test case: Database error
		startDate := "2023-01-01"
		endDate := "2023-01-31"
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByDateRange", s.ctx, startDate, endDate, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		orders, total, err := s.mockRepo.FindByDateRange(s.ctx, startDate, endDate, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestOrderRepositorySuite runs the test suite
func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}
