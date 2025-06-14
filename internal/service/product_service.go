package service

import (
	"context"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"
)

// ProductService defines the interface for product-related business logic
type ProductService interface {
	// GetProductByID retrieves a product by its ID
	GetProductByID(ctx context.Context, id uint) (*domain.Product, error)

	// GetProducts retrieves all products with optional pagination
	GetProducts(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error)

	// GetProductsByCategory retrieves products by category ID with optional pagination
	GetProductsByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error)

	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, product *domain.Product) error

	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, product *domain.Product) error

	// DeleteProduct deletes a product by its ID
	DeleteProduct(ctx context.Context, id uint) error

	// GetProductBySKU retrieves a product by its SKU
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)

	// UpdateProductStock updates the stock of a product
	UpdateProductStock(ctx context.Context, id uint, quantity int) error

	// GetCategories retrieves all product categories
	GetCategories(ctx context.Context) ([]domain.ProductCategory, error)

	// GetCategoryByID retrieves a product category by its ID
	GetCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error)

	// CreateCategory creates a new product category
	CreateCategory(ctx context.Context, category *domain.ProductCategory) error

	// UpdateCategory updates an existing product category
	UpdateCategory(ctx context.Context, category *domain.ProductCategory) error

	// DeleteCategory deletes a product category by its ID
	DeleteCategory(ctx context.Context, id uint) error
}

// ProductServiceImpl implements the ProductService interface
type ProductServiceImpl struct {
	productRepo repository.ProductRepository
}

// NewProductService creates a new ProductServiceImpl
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

// GetProductByID retrieves a product by its ID
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, id)
}

// GetProducts retrieves all products with optional pagination
func (s *ProductServiceImpl) GetProducts(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	return s.productRepo.FindAll(ctx, page, pageSize)
}

// GetProductsByCategory retrieves products by category ID with optional pagination
func (s *ProductServiceImpl) GetProductsByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	return s.productRepo.FindByCategory(ctx, categoryID, page, pageSize)
}

// CreateProduct creates a new product
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Create(ctx, product)
}

// UpdateProduct updates an existing product
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

// DeleteProduct deletes a product by its ID
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id uint) error {
	return s.productRepo.Delete(ctx, id)
}

// GetProductBySKU retrieves a product by its SKU
func (s *ProductServiceImpl) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	return s.productRepo.FindBySKU(ctx, sku)
}

// UpdateProductStock updates the stock of a product
func (s *ProductServiceImpl) UpdateProductStock(ctx context.Context, id uint, quantity int) error {
	return s.productRepo.UpdateStock(ctx, id, quantity)
}

// GetCategories retrieves all product categories
func (s *ProductServiceImpl) GetCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	return s.productRepo.FindCategories(ctx)
}

// GetCategoryByID retrieves a product category by its ID
func (s *ProductServiceImpl) GetCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error) {
	return s.productRepo.FindCategoryByID(ctx, id)
}

// CreateCategory creates a new product category
func (s *ProductServiceImpl) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	return s.productRepo.CreateCategory(ctx, category)
}

// UpdateCategory updates an existing product category
func (s *ProductServiceImpl) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	return s.productRepo.UpdateCategory(ctx, category)
}

// DeleteCategory deletes a product category by its ID
func (s *ProductServiceImpl) DeleteCategory(ctx context.Context, id uint) error {
	return s.productRepo.DeleteCategory(ctx, id)
}