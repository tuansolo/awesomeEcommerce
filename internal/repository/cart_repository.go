package repository

import (
	"context"

	"awesomeEcommerce/internal/domain"
)

// CartRepository defines the interface for cart repository operations
type CartRepository interface {
	// FindByID retrieves a cart by its ID
	FindByID(ctx context.Context, id uint) (*domain.Cart, error)

	// FindByUserID retrieves a cart by user ID
	FindByUserID(ctx context.Context, userID uint) (*domain.Cart, error)

	// Create creates a new cart
	Create(ctx context.Context, cart *domain.Cart) error

	// Update updates an existing cart
	Update(ctx context.Context, cart *domain.Cart) error

	// Delete deletes a cart by its ID
	Delete(ctx context.Context, id uint) error

	// AddItem adds an item to a cart
	AddItem(ctx context.Context, cartItem *domain.CartItem) error

	// UpdateItem updates an item in a cart
	UpdateItem(ctx context.Context, cartItem *domain.CartItem) error

	// RemoveItem removes an item from a cart
	RemoveItem(ctx context.Context, cartID, itemID uint) error

	// ClearCart removes all items from a cart
	ClearCart(ctx context.Context, cartID uint) error

	// GetCartItems retrieves all items in a cart
	GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error)

	// GetCartTotal calculates the total price of all items in a cart
	GetCartTotal(ctx context.Context, cartID uint) (float64, error)
}