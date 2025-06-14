package repository

import (
	"context"

	"awesomeEcommerce/internal/domain"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// FindByID retrieves a user by its ID
	FindByID(ctx context.Context, id uint) (*domain.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by its ID
	Delete(ctx context.Context, id uint) error

	// FindAll retrieves all users with optional pagination
	FindAll(ctx context.Context, page, pageSize int) ([]domain.User, int64, error)

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, id uint, hashedPassword string) error

	// FindByRole retrieves users by role with optional pagination
	FindByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error)

	// GetUserOrders retrieves all orders for a user
	GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error)

	// GetUserCart retrieves the cart for a user
	GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error)
}