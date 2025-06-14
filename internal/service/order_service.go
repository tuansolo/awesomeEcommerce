package service

import (
	"context"
	"errors"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/messaging"
	"awesomeEcommerce/internal/repository"
)

// OrderService defines the interface for order-related business logic
type OrderService interface {
	// GetOrderByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, id uint) (*domain.Order, error)

	// GetOrdersByUserID retrieves orders by user ID with optional pagination
	GetOrdersByUserID(ctx context.Context, userID uint, page, pageSize int) ([]domain.Order, int64, error)

	// CreateOrder creates a new order from a cart
	CreateOrder(ctx context.Context, userID uint, shippingAddress, billingAddress string) (*domain.Order, error)

	// UpdateOrderStatus updates the status of an order
	UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error

	// CancelOrder cancels an order
	CancelOrder(ctx context.Context, id uint) error

	// GetAllOrders retrieves all orders with optional pagination
	GetAllOrders(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error)

	// GetOrdersByStatus retrieves orders by status with optional pagination
	GetOrdersByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error)

	// GetOrdersByDateRange retrieves orders created within a date range
	GetOrdersByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error)

	// GetOrderItems retrieves all items in an order
	GetOrderItems(ctx context.Context, orderID uint) ([]domain.OrderItem, error)

	// GetOrderTotal calculates the total price of an order
	GetOrderTotal(ctx context.Context, orderID uint) (float64, error)
}

// OrderServiceImpl implements the OrderService interface
type OrderServiceImpl struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	producer    *messaging.KafkaProducer
}

// NewOrderService creates a new OrderServiceImpl
func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	producer *messaging.KafkaProducer,
) OrderService {
	return &OrderServiceImpl{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		producer:    producer,
	}
}

// GetOrderByID retrieves an order by its ID
func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

// GetOrdersByUserID retrieves orders by user ID with optional pagination
func (s *OrderServiceImpl) GetOrdersByUserID(ctx context.Context, userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, 0, errors.New("user not found")
	}

	return s.orderRepo.FindByUserID(ctx, userID, page, pageSize)
}

// CreateOrder creates a new order from a cart
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, userID uint, shippingAddress, billingAddress string) (*domain.Order, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get the user's cart
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Check if cart has items
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Calculate total amount
	total, err := s.cartRepo.GetCartTotal(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	// Create order items from cart items
	var orderItems []domain.OrderItem
	for _, cartItem := range cart.Items {
		// Check if product exists and has enough stock
		product, err := s.productRepo.FindByID(ctx, cartItem.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < cartItem.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Create order item
		orderItem := domain.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    cartItem.Quantity,
		}
		orderItems = append(orderItems, orderItem)

		// Update product stock
		err = s.productRepo.UpdateStock(ctx, product.ID, -cartItem.Quantity)
		if err != nil {
			return nil, err
		}
	}

	// Create the order
	order := &domain.Order{
		UserID:          userID,
		Items:           orderItems,
		TotalAmount:     total,
		Status:          domain.OrderStatusPending,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
	}

	// Save the order
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// Clear the cart
	err = s.cartRepo.ClearCart(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	// Publish order created event
	// Note: In a real application, we would serialize the order to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "order-created", []byte(fmt.Sprintf("%d", order.ID)), []byte(orderJSON))

	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderServiceImpl) UpdateOrderStatus(ctx context.Context, id uint, status domain.OrderStatus) error {
	// Check if order exists
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("order not found")
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, status) {
		return errors.New("invalid status transition")
	}

	// Update the status
	err = s.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// Publish order updated event
	// Note: In a real application, we would serialize the order to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "order-updated", []byte(fmt.Sprintf("%d", id)), []byte(orderJSON))

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func (s *OrderServiceImpl) isValidStatusTransition(from, to domain.OrderStatus) bool {
	// Define valid transitions
	validTransitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderStatusPending: {
			domain.OrderStatusProcessing,
			domain.OrderStatusCancelled,
		},
		domain.OrderStatusProcessing: {
			domain.OrderStatusShipped,
			domain.OrderStatusCancelled,
		},
		domain.OrderStatusShipped: {
			domain.OrderStatusDelivered,
		},
		domain.OrderStatusDelivered: {},
		domain.OrderStatusCancelled: {},
	}

	// Check if the transition is valid
	for _, validTo := range validTransitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// CancelOrder cancels an order
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, id uint) error {
	// Check if order exists
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("order not found")
	}

	// Check if order can be cancelled
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusProcessing {
		return errors.New("order cannot be cancelled")
	}

	// Update the status
	err = s.orderRepo.UpdateStatus(ctx, id, domain.OrderStatusCancelled)
	if err != nil {
		return err
	}

	// Return items to inventory
	for _, item := range order.Items {
		err = s.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Publish order cancelled event
	// Note: In a real application, we would serialize the order to JSON
	// and publish it to Kafka. For simplicity, we're just logging here.
	// s.producer.Publish(ctx, "order-cancelled", []byte(fmt.Sprintf("%d", id)), []byte(orderJSON))

	return nil
}

// GetAllOrders retrieves all orders with optional pagination
func (s *OrderServiceImpl) GetAllOrders(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error) {
	return s.orderRepo.FindAll(ctx, page, pageSize)
}

// GetOrdersByStatus retrieves orders by status with optional pagination
func (s *OrderServiceImpl) GetOrdersByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error) {
	return s.orderRepo.FindByStatus(ctx, status, page, pageSize)
}

// GetOrdersByDateRange retrieves orders created within a date range
func (s *OrderServiceImpl) GetOrdersByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error) {
	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, errors.New("invalid start date format, use YYYY-MM-DD")
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, errors.New("invalid end date format, use YYYY-MM-DD")
	}

	return s.orderRepo.FindByDateRange(ctx, startDate, endDate, page, pageSize)
}

// GetOrderItems retrieves all items in an order
func (s *OrderServiceImpl) GetOrderItems(ctx context.Context, orderID uint) ([]domain.OrderItem, error) {
	return s.orderRepo.GetOrderItems(ctx, orderID)
}

// GetOrderTotal calculates the total price of an order
func (s *OrderServiceImpl) GetOrderTotal(ctx context.Context, orderID uint) (float64, error) {
	return s.orderRepo.GetOrderTotal(ctx, orderID)
}