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

// ProductRepositoryImpl implements the ProductRepository interface
type ProductRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

// NewProductRepository creates a new ProductRepositoryImpl
func NewProductRepository(db *gorm.DB, cache *cache.RedisClient) repository.ProductRepository {
	return &ProductRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// FindByID retrieves a product by its ID
func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("product:%d", id)
	cachedProduct, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var product domain.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			return &product, nil
		}
	}

	// Cache miss, get from database
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	// Store in cache for future requests
	productJSON, err := json.Marshal(product)
	if err == nil {
		r.cache.Set(ctx, cacheKey, productJSON, 30*time.Minute)
	}

	return &product, nil
}

// FindAll retrieves all products with optional pagination
func (r *ProductRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	// Count total records
	if err := r.db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindByCategory retrieves products by category ID with optional pagination
func (r *ProductRepositoryImpl) FindByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	// Count total records for the category
	if err := r.db.Model(&domain.Product{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := r.db.Where("category_id = ?", categoryID).Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Create creates a new product
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *domain.Product) error {
	return r.db.Create(product).Error
}

// Update updates an existing product
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *domain.Product) error {
	// Update in database
	if err := r.db.Save(product).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// Delete deletes a product by its ID
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// Delete from database
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", id)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// FindBySKU retrieves a product by its SKU
func (r *ProductRepositoryImpl) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("product:sku:%s", sku)
	cachedProduct, err := r.cache.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit
		var product domain.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			return &product, nil
		}
	}

	// Cache miss, get from database
	var product domain.Product
	if err := r.db.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}

	// Store in cache for future requests
	productJSON, err := json.Marshal(product)
	if err == nil {
		r.cache.Set(ctx, cacheKey, productJSON, 30*time.Minute)
	}

	return &product, nil
}

// UpdateStock updates the stock of a product
func (r *ProductRepositoryImpl) UpdateStock(ctx context.Context, id uint, quantity int) error {
	// Update stock in database
	if err := r.db.Model(&domain.Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", id)
	r.cache.Delete(ctx, cacheKey)

	return nil
}

// FindCategories retrieves all product categories
func (r *ProductRepositoryImpl) FindCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	var categories []domain.ProductCategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// FindCategoryByID retrieves a product category by its ID
func (r *ProductRepositoryImpl) FindCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error) {
	var category domain.ProductCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateCategory creates a new product category
func (r *ProductRepositoryImpl) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	return r.db.Create(category).Error
}

// UpdateCategory updates an existing product category
func (r *ProductRepositoryImpl) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	return r.db.Save(category).Error
}

// DeleteCategory deletes a product category by its ID
func (r *ProductRepositoryImpl) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.Delete(&domain.ProductCategory{}, id).Error
}