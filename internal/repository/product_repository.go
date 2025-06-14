package repository

import (
	"context"

	"awesomeEcommerce/internal/domain"
)

// ProductRepository defines the interface for product repository operations
type ProductRepository interface {
	// FindByID retrieves a product by its ID
	FindByID(ctx context.Context, id uint) (*domain.Product, error)

	// FindAll retrieves all products with optional pagination
	FindAll(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error)

	// FindByCategory retrieves products by category ID with optional pagination
	FindByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error)

	// Create creates a new product
	Create(ctx context.Context, product *domain.Product) error

	// Update updates an existing product
	Update(ctx context.Context, product *domain.Product) error

	// Delete deletes a product by its ID
	Delete(ctx context.Context, id uint) error

	// FindBySKU retrieves a product by its SKU
	FindBySKU(ctx context.Context, sku string) (*domain.Product, error)

	// UpdateStock updates the stock of a product
	UpdateStock(ctx context.Context, id uint, quantity int) error

	// FindCategories retrieves all product categories
	FindCategories(ctx context.Context) ([]domain.ProductCategory, error)

	// FindCategoryByID retrieves a product category by its ID
	FindCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error)

	// CreateCategory creates a new product category
	CreateCategory(ctx context.Context, category *domain.ProductCategory) error

	// UpdateCategory updates an existing product category
	UpdateCategory(ctx context.Context, category *domain.ProductCategory) error

	// DeleteCategory deletes a product category by its ID
	DeleteCategory(ctx context.Context, id uint) error
}