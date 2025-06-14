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

// OrderRepositoryImpl implements the OrderRepository interface
type OrderRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewOrderRepository creates a new OrderRepositoryImpl
func NewOrderRepository(db *gorm.DB, cache *cache.RedisClient) repository.OrderRepository {
	return &OrderRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// FindByID retrieves an order by its ID
func (r *OrderRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("order:%d", id)
	cachedOrder, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var order domain.Order
		if err := json.Unmarshal([]byte(cachedOrder), &order); err == nil {
			// Get order items separately
			if err := r.db.Where("order_id = ?", id).Find(&order.Items).Error; err != nil {
				return nil, err
			}
			return &order, nil
		}
	}

	// Cache miss, get from database
	var order domain.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}

	// Get order items
	if err := r.db.Where("order_id = ?", id).Find(&order.Items).Error; err != nil {
		return nil, err
	}

	// Store in cache for future requests (without items to avoid circular references)
	orderCopy := order
	orderCopy.Items = nil
	orderJSON, err := json.Marshal(orderCopy)
	if err == nil {
		r.cache.Set(ctx, cacheKey, orderJSON, 30*time.Minute)
	}

	return &order, nil
}

// FindByUserID retrieves orders by user ID with optional pagination
func (r *OrderRepositoryImpl) FindByUserID(ctx context.Context, userID uint, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	// Count total records for the user
	if err := r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// Get items for each order
	for i := range orders {
		if err := r.db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			return nil, 0, err
		}
	}

	return orders, total, nil
}

// Create creates a new order
func (r *OrderRepositoryImpl) Create(ctx context.Context, order *domain.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Create order items
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
			if err := tx.Create(&order.Items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update updates an existing order
func (r *OrderRepositoryImpl) Update(ctx context.Context, order *domain.Order) error {
	if err := r.db.Save(order).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("order:%d", order.ID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// UpdateStatus updates the status of an order
func (r *OrderRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status domain.OrderStatus) error {
	if err := r.db.Model(&domain.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("order:%d", id)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// Delete deletes an order by its ID
func (r *OrderRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete order items first
		if err := tx.Where("order_id = ?", id).Delete(&domain.OrderItem{}).Error; err != nil {
			return err
		}

		// Delete the order
		if err := tx.Delete(&domain.Order{}, id).Error; err != nil {
			return err
		}

		// Invalidate cache
		cacheKey := fmt.Sprintf("order:%d", id)
		r.cache.Delete(ctx, cacheKey)

		return nil
	})
}

// AddOrderItem adds an item to an order
func (r *OrderRepositoryImpl) AddOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	if err := r.db.Create(orderItem).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("order:%d", orderItem.OrderID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// GetOrderItems retrieves all items in an order
func (r *OrderRepositoryImpl) GetOrderItems(ctx context.Context, orderID uint) ([]domain.OrderItem, error) {
	var items []domain.OrderItem
	if err := r.db.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindAll retrieves all orders with optional pagination
func (r *OrderRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// Get items for each order
	for i := range orders {
		if err := r.db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			return nil, 0, err
		}
	}

	return orders, total, nil
}

// FindByStatus retrieves orders by status with optional pagination
func (r *OrderRepositoryImpl) FindByStatus(ctx context.Context, status domain.OrderStatus, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	// Count total records for the status
	if err := r.db.Model(&domain.Order{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// Get items for each order
	for i := range orders {
		if err := r.db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			return nil, 0, err
		}
	}

	return orders, total, nil
}

// GetOrderTotal calculates the total price of an order
func (r *OrderRepositoryImpl) GetOrderTotal(ctx context.Context, orderID uint) (float64, error) {
	var total float64
	err := r.db.Model(&domain.OrderItem{}).
		Select("SUM(price * quantity)").
		Where("order_id = ?", orderID).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// FindByDateRange retrieves orders created within a date range
func (r *OrderRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	// Count total records in the date range
	if err := r.db.Model(&domain.Order{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// Get items for each order
	for i := range orders {
		if err := r.db.Where("order_id = ?", orders[i].ID).Find(&orders[i].Items).Error; err != nil {
			return nil, 0, err
		}
	}

	return orders, total, nil
}