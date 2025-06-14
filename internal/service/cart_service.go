package service

import (
	"context"
	"errors"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"
)

// CartService defines the interface for cart-related business logic
type CartService interface {
	// GetCartByID retrieves a cart by its ID
	GetCartByID(ctx context.Context, id uint) (*domain.Cart, error)

	// GetCartByUserID retrieves a cart by user ID
	GetCartByUserID(ctx context.Context, userID uint) (*domain.Cart, error)

	// CreateCart creates a new cart
	CreateCart(ctx context.Context, cart *domain.Cart) error

	// AddItemToCart adds an item to a cart
	AddItemToCart(ctx context.Context, userID, productID uint, quantity int) error

	// UpdateCartItem updates an item in a cart
	UpdateCartItem(ctx context.Context, cartID, itemID uint, quantity int) error

	// RemoveItemFromCart removes an item from a cart
	RemoveItemFromCart(ctx context.Context, cartID, itemID uint) error

	// ClearCart removes all items from a cart
	ClearCart(ctx context.Context, cartID uint) error

	// GetCartItems retrieves all items in a cart
	GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error)

	// GetCartTotal calculates the total price of all items in a cart
	GetCartTotal(ctx context.Context, cartID uint) (float64, error)
}

// CartServiceImpl implements the CartService interface
type CartServiceImpl struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
}

// NewCartService creates a new CartServiceImpl
func NewCartService(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
) CartService {
	return &CartServiceImpl{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// GetCartByID retrieves a cart by its ID
func (s *CartServiceImpl) GetCartByID(ctx context.Context, id uint) (*domain.Cart, error) {
	return s.cartRepo.FindByID(ctx, id)
}

// GetCartByUserID retrieves a cart by user ID
func (s *CartServiceImpl) GetCartByUserID(ctx context.Context, userID uint) (*domain.Cart, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Try to find existing cart
	cart, err := s.cartRepo.FindByUserID(ctx, userID)
	if err == nil {
		return cart, nil
	}

	// Create new cart if not found
	newCart := &domain.Cart{
		UserID: userID,
		Items:  []domain.CartItem{},
	}
	if err := s.cartRepo.Create(ctx, newCart); err != nil {
		return nil, err
	}

	return newCart, nil
}

// CreateCart creates a new cart
func (s *CartServiceImpl) CreateCart(ctx context.Context, cart *domain.Cart) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, cart.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user already has a cart
	_, err = s.cartRepo.FindByUserID(ctx, cart.UserID)
	if err == nil {
		return errors.New("user already has a cart")
	}

	return s.cartRepo.Create(ctx, cart)
}

// AddItemToCart adds an item to a cart
func (s *CartServiceImpl) AddItemToCart(ctx context.Context, userID, productID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	// Get the user's cart
	cart, err := s.GetCartByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if product exists and has enough stock
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// Add item to cart
	cartItem := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
	}

	return s.cartRepo.AddItem(ctx, cartItem)
}

// UpdateCartItem updates an item in a cart
func (s *CartServiceImpl) UpdateCartItem(ctx context.Context, cartID, itemID uint, quantity int) error {
	if quantity <= 0 {
		return s.RemoveItemFromCart(ctx, cartID, itemID)
	}

	// Get the cart item
	items, err := s.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return err
	}

	var cartItem *domain.CartItem
	for i := range items {
		if items[i].ID == itemID {
			cartItem = &items[i]
			break
		}
	}

	if cartItem == nil {
		return errors.New("cart item not found")
	}

	// Check if product has enough stock
	product, err := s.productRepo.FindByID(ctx, cartItem.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// Update the cart item
	cartItem.Quantity = quantity
	return s.cartRepo.UpdateItem(ctx, cartItem)
}

// RemoveItemFromCart removes an item from a cart
func (s *CartServiceImpl) RemoveItemFromCart(ctx context.Context, cartID, itemID uint) error {
	return s.cartRepo.RemoveItem(ctx, cartID, itemID)
}

// ClearCart removes all items from a cart
func (s *CartServiceImpl) ClearCart(ctx context.Context, cartID uint) error {
	return s.cartRepo.ClearCart(ctx, cartID)
}

// GetCartItems retrieves all items in a cart
func (s *CartServiceImpl) GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	return s.cartRepo.GetCartItems(ctx, cartID)
}

// GetCartTotal calculates the total price of all items in a cart
func (s *CartServiceImpl) GetCartTotal(ctx context.Context, cartID uint) (float64, error) {
	return s.cartRepo.GetCartTotal(ctx, cartID)
}
