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

// CartRepositoryImpl implements the CartRepository interface
type CartRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewCartRepository creates a new CartRepositoryImpl
func NewCartRepository(db *gorm.DB, cache *cache.RedisClient) repository.CartRepository {
	return &CartRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// FindByID retrieves a cart by its ID
func (r *CartRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Cart, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("cart:%d", id)
	cachedCart, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var cart domain.Cart
		if err := json.Unmarshal([]byte(cachedCart), &cart); err == nil {
			// Get cart items separately
			if err := r.db.Where("cart_id = ?", id).Find(&cart.Items).Error; err != nil {
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
	if err := r.db.First(&cart, id).Error; err != nil {
		return nil, err
	}

	// Get cart items
	if err := r.db.Where("cart_id = ?", id).Find(&cart.Items).Error; err != nil {
		return nil, err
	}

	// Load product details for each cart item
	for i := range cart.Items {
		if err := r.db.First(&cart.Items[i].Product, cart.Items[i].ProductID).Error; err != nil {
			return nil, err
		}
	}

	// Store in cache for future requests (without items to avoid circular references)
	cartCopy := cart
	cartCopy.Items = nil
	cartJSON, err := json.Marshal(cartCopy)
	if err == nil {
		r.cache.Set(ctx, cacheKey, cartJSON, 15*time.Minute)
	}

	return &cart, nil
}

// FindByUserID retrieves a cart by user ID
func (r *CartRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*domain.Cart, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("cart:user:%d", userID)
	cachedCartID, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit for cart ID
		var cartID uint
		if err := json.Unmarshal([]byte(cachedCartID), &cartID); err == nil {
			return r.FindByID(ctx, cartID)
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

// Create creates a new cart
func (r *CartRepositoryImpl) Create(ctx context.Context, cart *domain.Cart) error {
	if err := r.db.Create(cart).Error; err != nil {
		return err
	}

	// Cache the user's cart ID
	cacheKey := fmt.Sprintf("cart:user:%d", cart.UserID)
	cartIDJSON, err := json.Marshal(cart.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, cartIDJSON, 15*time.Minute)
	}

	return nil
}

// Update updates an existing cart
func (r *CartRepositoryImpl) Update(ctx context.Context, cart *domain.Cart) error {
	if err := r.db.Save(cart).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("cart:%d", cart.ID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// Delete deletes a cart by its ID
func (r *CartRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// Get the cart to find the user ID
	var cart domain.Cart
	if err := r.db.First(&cart, id).Error; err != nil {
		return err
	}

	// Delete cart items first
	if err := r.db.Where("cart_id = ?", id).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}

	// Delete the cart
	if err := r.db.Delete(&domain.Cart{}, id).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("cart:%d", id))
	r.cache.Delete(ctx, fmt.Sprintf("cart:user:%d", cart.UserID))

	return nil
}

// AddItem adds an item to a cart
func (r *CartRepositoryImpl) AddItem(ctx context.Context, cartItem *domain.CartItem) error {
	// Check if the item already exists in the cart
	var existingItem domain.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartItem.CartID, cartItem.ProductID).First(&existingItem).Error
	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += cartItem.Quantity
		return r.UpdateItem(ctx, &existingItem)
	}

	// Item doesn't exist, create new
	if err := r.db.Create(cartItem).Error; err != nil {
		return err
	}

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("cart:%d", cartItem.CartID))

	return nil
}

// UpdateItem updates an item in a cart
func (r *CartRepositoryImpl) UpdateItem(ctx context.Context, cartItem *domain.CartItem) error {
	if err := r.db.Save(cartItem).Error; err != nil {
		return err
	}

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("cart:%d", cartItem.CartID))

	return nil
}

// RemoveItem removes an item from a cart
func (r *CartRepositoryImpl) RemoveItem(ctx context.Context, cartID, itemID uint) error {
	if err := r.db.Where("id = ? AND cart_id = ?", itemID, cartID).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("cart:%d", cartID))

	return nil
}

// ClearCart removes all items from a cart
func (r *CartRepositoryImpl) ClearCart(ctx context.Context, cartID uint) error {
	if err := r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("cart:%d", cartID))

	return nil
}

// GetCartItems retrieves all items in a cart
func (r *CartRepositoryImpl) GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	var items []domain.CartItem
	if err := r.db.Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
		return nil, err
	}

	// Load product details for each cart item
	for i := range items {
		if err := r.db.First(&items[i].Product, items[i].ProductID).Error; err != nil {
			return nil, err
		}
	}

	return items, nil
}

// GetCartTotal calculates the total price of all items in a cart
func (r *CartRepositoryImpl) GetCartTotal(ctx context.Context, cartID uint) (float64, error) {
	items, err := r.GetCartItems(ctx, cartID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range items {
		total += item.Product.Price * float64(item.Quantity)
	}

	return total, nil
}