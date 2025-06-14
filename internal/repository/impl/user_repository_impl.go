package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"
	"awesomeEcommerce/internal/repository/cache"

	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewUserRepository creates a new UserRepositoryImpl
func NewUserRepository(db *gorm.DB, cache *cache.RedisClient) repository.UserRepository {
	return &UserRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// FindByID retrieves a user by its ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)
	cachedUser, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var user domain.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
	}

	// Cache miss, get from database
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Store in cache for future requests (without sensitive data)
	userCopy := user
	userCopy.Password = "" // Don't cache password
	userJSON, err := json.Marshal(userCopy)
	if err == nil {
		r.cache.Set(ctx, cacheKey, userJSON, 30*time.Minute)
	}

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:email:%s", email)
	cachedUserID, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit for user ID
		var userID uint
		if err := json.Unmarshal([]byte(cachedUserID), &userID); err == nil {
			return r.FindByID(ctx, userID)
		}
	}

	// Cache miss, get from database
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	// Store user ID in cache for future requests
	userIDJSON, err := json.Marshal(user.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, userIDJSON, 30*time.Minute)
	}

	// Also cache the full user (without sensitive data)
	userCopy := user
	userCopy.Password = "" // Don't cache password
	userJSON, err := json.Marshal(userCopy)
	if err == nil {
		r.cache.Set(ctx, fmt.Sprintf("user:%d", user.ID), userJSON, 30*time.Minute)
	}

	return &user, nil
}

// Create creates a new user
func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}

	// Cache the user's email to ID mapping
	cacheKey := fmt.Sprintf("user:email:%s", user.Email)
	userIDJSON, err := json.Marshal(user.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, userIDJSON, 30*time.Minute)
	}

	return nil
}

// Update updates an existing user
func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	// Get the old user to check if email changed
	var oldUser domain.User
	if err := r.db.First(&oldUser, user.ID).Error; err != nil {
		return err
	}

	// Update in database
	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("user:%d", user.ID))
	
	// If email changed, invalidate old email cache and set new one
	if oldUser.Email != user.Email {
		r.cache.Delete(ctx, fmt.Sprintf("user:email:%s", oldUser.Email))
		
		// Cache the user's new email to ID mapping
		cacheKey := fmt.Sprintf("user:email:%s", user.Email)
		userIDJSON, err := json.Marshal(user.ID)
		if err == nil {
			r.cache.Set(ctx, cacheKey, userIDJSON, 30*time.Minute)
		}
	}

	return nil
}

// Delete deletes a user by its ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// Get the user first to get the email
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}

	// Delete the user
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
	r.cache.Delete(ctx, fmt.Sprintf("user:email:%s", user.Email))

	return nil
}

// FindAll retrieves all users with optional pagination
func (r *UserRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Count total records
	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Select("id, email, first_name, last_name, phone, address, role, created_at, updated_at").
		Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword updates a user's password
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	if err := r.db.Model(&domain.User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	// No need to invalidate cache since we don't cache passwords

	return nil
}

// FindByRole retrieves users by role with optional pagination
func (r *UserRepositoryImpl) FindByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Count total records for the role
	if err := r.db.Model(&domain.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Select("id, email, first_name, last_name, phone, address, role, created_at, updated_at").
		Where("role = ?", role).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserOrders retrieves all orders for a user
func (r *UserRepositoryImpl) GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	// Get items for each order
	for i := range orders {
		if err := r.db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			return nil, err
		}
	}

	return orders, nil
}

// GetUserCart retrieves the cart for a user
func (r *UserRepositoryImpl) GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("cart:user:%d", userID)
	cachedCartID, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit for cart ID
		var cartID uint
		if err := json.Unmarshal([]byte(cachedCartID), &cartID); err == nil {
			var cart domain.Cart
			if err := r.db.First(&cart, cartID).Error; err != nil {
				return nil, err
			}

			// Get cart items
			if err := r.db.Where("cart_id = ?", cart.ID).Find(&cart.Items).Error; err != nil {
				return nil, err
			}

			// Load product details for each cart item
			for i := range cart.Items {
				if err := r.db.First(&cart.Items[i].Product, cart.Items[i].ProductID).Error; err != nil {
					return nil, err
				}
			}

			return &cart, nil
		}
	}

	// Cache miss, get from database
	var cart domain.Cart
	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}

	// Get cart items
	if err := r.db.Where("cart_id = ?", cart.ID).Find(&cart.Items).Error; err != nil {
		return nil, err
	}

	// Load product details for each cart item
	for i := range cart.Items {
		if err := r.db.First(&cart.Items[i].Product, cart.Items[i].ProductID).Error; err != nil {
			return nil, err
		}
	}

	// Store cart ID in cache for future requests
	cartIDJSON, err := json.Marshal(cart.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, cartIDJSON, 15*time.Minute)
	}

	return &cart, nil
}