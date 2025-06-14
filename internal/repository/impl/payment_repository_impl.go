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

// PaymentRepositoryImpl implements the PaymentRepository interface
type PaymentRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewPaymentRepository creates a new PaymentRepositoryImpl
func NewPaymentRepository(db *gorm.DB, cache *cache.RedisClient) repository.PaymentRepository {
	return &PaymentRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// FindByID retrieves a payment by its ID
func (r *PaymentRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Payment, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("payment:%d", id)
	cachedPayment, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var payment domain.Payment
		if err := json.Unmarshal([]byte(cachedPayment), &payment); err == nil {
			return &payment, nil
		}
	}

	// Cache miss, get from database
	var payment domain.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return nil, err
	}

	// Store in cache for future requests
	paymentJSON, err := json.Marshal(payment)
	if err == nil {
		r.cache.Set(ctx, cacheKey, paymentJSON, 30*time.Minute)
	}

	return &payment, nil
}

// FindByOrderID retrieves a payment by order ID
func (r *PaymentRepositoryImpl) FindByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("payment:order:%d", orderID)
	cachedPaymentID, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit for payment ID
		var paymentID uint
		if err := json.Unmarshal([]byte(cachedPaymentID), &paymentID); err == nil {
			return r.FindByID(ctx, paymentID)
		}
	}

	// Cache miss, get from database
	var payment domain.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}

	// Store payment ID in cache for future requests
	paymentIDJSON, err := json.Marshal(payment.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, paymentIDJSON, 30*time.Minute)
	}

	// Also cache the full payment
	paymentJSON, err := json.Marshal(payment)
	if err == nil {
		r.cache.Set(ctx, fmt.Sprintf("payment:%d", payment.ID), paymentJSON, 30*time.Minute)
	}

	return &payment, nil
}

// Create creates a new payment
func (r *PaymentRepositoryImpl) Create(ctx context.Context, payment *domain.Payment) error {
	if err := r.db.Create(payment).Error; err != nil {
		return err
	}

	// Cache the payment by order ID
	cacheKey := fmt.Sprintf("payment:order:%d", payment.OrderID)
	paymentIDJSON, err := json.Marshal(payment.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, paymentIDJSON, 30*time.Minute)
	}

	return nil
}

// Update updates an existing payment
func (r *PaymentRepositoryImpl) Update(ctx context.Context, payment *domain.Payment) error {
	if err := r.db.Save(payment).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("payment:%d", payment.ID))
	r.cache.Delete(ctx, fmt.Sprintf("payment:order:%d", payment.OrderID))
	r.cache.Delete(ctx, fmt.Sprintf("payment:transaction:%s", payment.TransactionID))

	return nil
}

// UpdateStatus updates the status of a payment
func (r *PaymentRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status domain.PaymentStatus) error {
	// Get the payment first to get the order ID and transaction ID
	var payment domain.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return err
	}

	// Update the status
	if err := r.db.Model(&domain.Payment{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("payment:%d", id))
	r.cache.Delete(ctx, fmt.Sprintf("payment:order:%d", payment.OrderID))
	if payment.TransactionID != "" {
		r.cache.Delete(ctx, fmt.Sprintf("payment:transaction:%s", payment.TransactionID))
	}

	return nil
}

// Delete deletes a payment by its ID
func (r *PaymentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// Get the payment first to get the order ID and transaction ID
	var payment domain.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		return err
	}

	// Delete the payment
	if err := r.db.Delete(&domain.Payment{}, id).Error; err != nil {
		return err
	}

	// Invalidate caches
	r.cache.Delete(ctx, fmt.Sprintf("payment:%d", id))
	r.cache.Delete(ctx, fmt.Sprintf("payment:order:%d", payment.OrderID))
	if payment.TransactionID != "" {
		r.cache.Delete(ctx, fmt.Sprintf("payment:transaction:%s", payment.TransactionID))
	}

	return nil
}

// FindByTransactionID retrieves a payment by its transaction ID
func (r *PaymentRepositoryImpl) FindByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("payment:transaction:%s", transactionID)
	cachedPaymentID, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit for payment ID
		var paymentID uint
		if err := json.Unmarshal([]byte(cachedPaymentID), &paymentID); err == nil {
			return r.FindByID(ctx, paymentID)
		}
	}

	// Cache miss, get from database
	var payment domain.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}

	// Store payment ID in cache for future requests
	paymentIDJSON, err := json.Marshal(payment.ID)
	if err == nil {
		r.cache.Set(ctx, cacheKey, paymentIDJSON, 30*time.Minute)
	}

	// Also cache the full payment
	paymentJSON, err := json.Marshal(payment)
	if err == nil {
		r.cache.Set(ctx, fmt.Sprintf("payment:%d", payment.ID), paymentJSON, 30*time.Minute)
	}

	return &payment, nil
}

// FindAll retrieves all payments with optional pagination
func (r *PaymentRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// FindByStatus retrieves payments by status with optional pagination
func (r *PaymentRepositoryImpl) FindByStatus(ctx context.Context, status domain.PaymentStatus, page, pageSize int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	// Count total records for the status
	if err := r.db.Model(&domain.Payment{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// FindByDateRange retrieves payments created within a date range
func (r *PaymentRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	// Count total records in the date range
	if err := r.db.Model(&domain.Payment{}).
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
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// FindByMethod retrieves payments by payment method with optional pagination
func (r *PaymentRepositoryImpl) FindByMethod(ctx context.Context, method domain.PaymentMethod, page, pageSize int) ([]domain.Payment, int64, error) {
	var payments []domain.Payment
	var total int64

	// Count total records for the method
	if err := r.db.Model(&domain.Payment{}).Where("method = ?", method).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("method = ?", method).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}