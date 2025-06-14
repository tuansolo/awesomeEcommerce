package service_test

import (
	"context"
	"testing"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
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
	return args.Get(0).([]domain.OrderItem), args.Error(1)
}

func (m *MockOrderRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, status, page, pageSize)
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) GetOrderTotal(ctx context.Context, orderID uint) (float64, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockOrderRepository) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error) {
	args := m.Called(ctx, startDate, endDate, page, pageSize)
	return args.Get(0).([]domain.Order), args.Get(1).(int64), args.Error(2)
}

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FindByID(ctx context.Context, id uint) (*domain.Cart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) Update(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCartRepository) AddItem(ctx context.Context, cartItem *domain.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartRepository) UpdateItem(ctx context.Context, cartItem *domain.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartRepository) RemoveItem(ctx context.Context, cartID, itemID uint) error {
	args := m.Called(ctx, cartID, itemID)
	return args.Error(0)
}

func (m *MockCartRepository) ClearCart(ctx context.Context, cartID uint) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *MockCartRepository) GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	args := m.Called(ctx, cartID)
	return args.Get(0).([]domain.CartItem), args.Error(1)
}

func (m *MockCartRepository) GetCartTotal(ctx context.Context, cartID uint) (float64, error) {
	args := m.Called(ctx, cartID)
	return args.Get(0).(float64), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) FindByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) FindCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.ProductCategory), args.Error(1)
}

func (m *MockProductRepository) FindCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *MockProductRepository) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteCategory(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) FindByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error) {
	args := m.Called(ctx, role, page, pageSize)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockUserRepository) GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

// We don't need a MockKafkaProducer since we'll use nil for the producer parameter
// The Kafka publishing code is commented out in the service implementation

// Test cases
func TestGetOrderByID(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)

	// Use nil for the producer parameter since Kafka publishing is commented out in the service
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockProductRepo, mockUserRepo, nil)

	ctx := context.Background()
	orderID := uint(1)
	expectedOrder := &domain.Order{
		ID:              orderID,
		UserID:          uint(1),
		TotalAmount:     100.0,
		Status:          domain.OrderStatusPending,
		ShippingAddress: "123 Main St",
		BillingAddress:  "123 Main St",
	}

	// Expectations
	mockOrderRepo.On("FindByID", ctx, orderID).Return(expectedOrder, nil)

	// Execute
	order, err := orderService.GetOrderByID(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	mockOrderRepo.AssertExpectations(t)
}

func TestCreateOrder(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)

	// Use nil for the producer parameter since Kafka publishing is commented out in the service
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockProductRepo, mockUserRepo, nil)

	ctx := context.Background()
	userID := uint(1)
	shippingAddress := "123 Main St"
	billingAddress := "123 Main St"

	user := &domain.User{
		ID:    userID,
		Email: "user@example.com",
	}

	cart := &domain.Cart{
		ID:     uint(1),
		UserID: userID,
		Items: []domain.CartItem{
			{
				ID:        uint(1),
				CartID:    uint(1),
				ProductID: uint(1),
				Quantity:  2,
			},
		},
	}

	product := &domain.Product{
		ID:    uint(1),
		Name:  "Test Product",
		Price: 50.0,
		Stock: 10,
	}

	// Expectations
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockCartRepo.On("FindByUserID", ctx, userID).Return(cart, nil)
	mockCartRepo.On("GetCartTotal", ctx, cart.ID).Return(100.0, nil)
	mockProductRepo.On("FindByID", ctx, uint(1)).Return(product, nil)
	mockProductRepo.On("UpdateStock", ctx, uint(1), -2).Return(nil)
	mockOrderRepo.On("Create", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)
	mockCartRepo.On("ClearCart", ctx, cart.ID).Return(nil)

	// Execute
	order, err := orderService.CreateOrder(ctx, userID, shippingAddress, billingAddress)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	assert.Equal(t, shippingAddress, order.ShippingAddress)
	assert.Equal(t, billingAddress, order.BillingAddress)
	assert.Equal(t, 100.0, order.TotalAmount)
	assert.Len(t, order.Items, 1)

	mockUserRepo.AssertExpectations(t)
	mockCartRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)

	// Use nil for the producer parameter since Kafka publishing is commented out in the service
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockProductRepo, mockUserRepo, nil)

	ctx := context.Background()
	orderID := uint(1)
	currentStatus := domain.OrderStatusPending
	newStatus := domain.OrderStatusProcessing

	order := &domain.Order{
		ID:     orderID,
		Status: currentStatus,
	}

	// Expectations
	mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil)
	mockOrderRepo.On("UpdateStatus", ctx, orderID, newStatus).Return(nil)

	// Execute
	err := orderService.UpdateOrderStatus(ctx, orderID, newStatus)

	// Assert
	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
}

func TestCancelOrder(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)

	// Use nil for the producer parameter since Kafka publishing is commented out in the service
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockProductRepo, mockUserRepo, nil)

	ctx := context.Background()
	orderID := uint(1)

	order := &domain.Order{
		ID:     orderID,
		Status: domain.OrderStatusPending,
		Items: []domain.OrderItem{
			{
				ID:        uint(1),
				OrderID:   orderID,
				ProductID: uint(1),
				Quantity:  2,
			},
		},
	}

	// Expectations
	mockOrderRepo.On("FindByID", ctx, orderID).Return(order, nil)
	mockOrderRepo.On("UpdateStatus", ctx, orderID, domain.OrderStatusCancelled).Return(nil)
	mockProductRepo.On("UpdateStock", ctx, uint(1), 2).Return(nil)

	// Execute
	err := orderService.CancelOrder(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestGetOrdersByUserID(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)

	// Use nil for the producer parameter since Kafka publishing is commented out in the service
	orderService := service.NewOrderService(mockOrderRepo, mockCartRepo, mockProductRepo, mockUserRepo, nil)

	ctx := context.Background()
	userID := uint(1)
	page := 1
	pageSize := 10

	user := &domain.User{
		ID:    userID,
		Email: "user@example.com",
	}

	expectedOrders := []domain.Order{
		{
			ID:     uint(1),
			UserID: userID,
			Status: domain.OrderStatusPending,
		},
		{
			ID:     uint(2),
			UserID: userID,
			Status: domain.OrderStatusProcessing,
		},
	}
	expectedTotal := int64(2)

	// Expectations
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockOrderRepo.On("FindByUserID", ctx, userID, page, pageSize).Return(expectedOrders, expectedTotal, nil)

	// Execute
	orders, total, err := orderService.GetOrdersByUserID(ctx, userID, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	assert.Equal(t, expectedTotal, total)
	mockUserRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}
