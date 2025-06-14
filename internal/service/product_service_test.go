package service_test

import (
	"context"
	"errors"
	"testing"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestGetProductByID(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	productID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedProduct := &domain.Product{
			ID:          productID,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			SKU:         "TEST-SKU-123",
			Stock:       100,
			CategoryID:  1,
		}

		// Expectations
		mockRepo.On("FindByID", ctx, productID).Return(expectedProduct, nil).Once()

		// Execute
		product, err := productService.GetProductByID(ctx, productID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindByID", ctx, productID).Return(nil, errors.New("product not found")).Once()

		// Execute
		product, err := productService.GetProductByID(ctx, productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "product not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProducts(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	page := 1
	pageSize := 10

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedProducts := []domain.Product{
			{
				ID:          1,
				Name:        "Product 1",
				Description: "Description 1",
				Price:       99.99,
				SKU:         "SKU-1",
				Stock:       100,
				CategoryID:  1,
			},
			{
				ID:          2,
				Name:        "Product 2",
				Description: "Description 2",
				Price:       199.99,
				SKU:         "SKU-2",
				Stock:       50,
				CategoryID:  2,
			},
		}
		expectedTotal := int64(2)

		// Expectations
		mockRepo.On("FindAll", ctx, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := productService.GetProducts(ctx, page, pageSize)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindAll", ctx, page, pageSize).Return([]domain.Product{}, int64(0), errors.New("database error")).Once()

		// Execute
		products, total, err := productService.GetProducts(ctx, page, pageSize)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, products)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProductsByCategory(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	categoryID := uint(1)
	page := 1
	pageSize := 10

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedProducts := []domain.Product{
			{
				ID:          1,
				Name:        "Product 1",
				Description: "Description 1",
				Price:       99.99,
				SKU:         "SKU-1",
				Stock:       100,
				CategoryID:  categoryID,
			},
			{
				ID:          3,
				Name:        "Product 3",
				Description: "Description 3",
				Price:       299.99,
				SKU:         "SKU-3",
				Stock:       75,
				CategoryID:  categoryID,
			},
		}
		expectedTotal := int64(2)

		// Expectations
		mockRepo.On("FindByCategory", ctx, categoryID, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := productService.GetProductsByCategory(ctx, categoryID, page, pageSize)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, products)
		assert.Equal(t, expectedTotal, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindByCategory", ctx, categoryID, page, pageSize).Return([]domain.Product{}, int64(0), errors.New("database error")).Once()

		// Execute
		products, total, err := productService.GetProductsByCategory(ctx, categoryID, page, pageSize)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, products)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateProduct(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		product := &domain.Product{
			Name:        "New Product",
			Description: "New Description",
			Price:       149.99,
			SKU:         "NEW-SKU-123",
			Stock:       200,
			CategoryID:  1,
		}

		// Expectations
		mockRepo.On("Create", ctx, product).Return(nil).Once()

		// Execute
		err := productService.CreateProduct(ctx, product)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Test data
		product := &domain.Product{
			Name:        "New Product",
			Description: "New Description",
			Price:       149.99,
			SKU:         "NEW-SKU-123",
			Stock:       200,
			CategoryID:  1,
		}

		// Expectations
		mockRepo.On("Create", ctx, product).Return(errors.New("database error")).Once()

		// Execute
		err := productService.CreateProduct(ctx, product)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProduct(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		product := &domain.Product{
			ID:          1,
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       199.99,
			SKU:         "UPD-SKU-123",
			Stock:       150,
			CategoryID:  2,
		}

		// Expectations
		mockRepo.On("Update", ctx, product).Return(nil).Once()

		// Execute
		err := productService.UpdateProduct(ctx, product)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Test data
		product := &domain.Product{
			ID:          1,
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       199.99,
			SKU:         "UPD-SKU-123",
			Stock:       150,
			CategoryID:  2,
		}

		// Expectations
		mockRepo.On("Update", ctx, product).Return(errors.New("database error")).Once()

		// Execute
		err := productService.UpdateProduct(ctx, product)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	productID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Expectations
		mockRepo.On("Delete", ctx, productID).Return(nil).Once()

		// Execute
		err := productService.DeleteProduct(ctx, productID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("Delete", ctx, productID).Return(errors.New("database error")).Once()

		// Execute
		err := productService.DeleteProduct(ctx, productID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetProductBySKU(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	sku := "TEST-SKU-123"

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedProduct := &domain.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			SKU:         sku,
			Stock:       100,
			CategoryID:  1,
		}

		// Expectations
		mockRepo.On("FindBySKU", ctx, sku).Return(expectedProduct, nil).Once()

		// Execute
		product, err := productService.GetProductBySKU(ctx, sku)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindBySKU", ctx, sku).Return(nil, errors.New("product not found")).Once()

		// Execute
		product, err := productService.GetProductBySKU(ctx, sku)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "product not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProductStock(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	productID := uint(1)
	quantity := 50

	t.Run("Success", func(t *testing.T) {
		// Expectations
		mockRepo.On("UpdateStock", ctx, productID, quantity).Return(nil).Once()

		// Execute
		err := productService.UpdateProductStock(ctx, productID, quantity)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("UpdateStock", ctx, productID, quantity).Return(errors.New("database error")).Once()

		// Execute
		err := productService.UpdateProductStock(ctx, productID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCategories(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedCategories := []domain.ProductCategory{
			{
				ID:   1,
				Name: "Category 1",
			},
			{
				ID:   2,
				Name: "Category 2",
			},
		}

		// Expectations
		mockRepo.On("FindCategories", ctx).Return(expectedCategories, nil).Once()

		// Execute
		categories, err := productService.GetCategories(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindCategories", ctx).Return([]domain.ProductCategory{}, errors.New("database error")).Once()

		// Execute
		categories, err := productService.GetCategories(ctx)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, categories)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetCategoryByID(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	categoryID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedCategory := &domain.ProductCategory{
			ID:   categoryID,
			Name: "Test Category",
		}

		// Expectations
		mockRepo.On("FindCategoryByID", ctx, categoryID).Return(expectedCategory, nil).Once()

		// Execute
		category, err := productService.GetCategoryByID(ctx, categoryID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("FindCategoryByID", ctx, categoryID).Return(nil, errors.New("category not found")).Once()

		// Execute
		category, err := productService.GetCategoryByID(ctx, categoryID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Contains(t, err.Error(), "category not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateCategory(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		category := &domain.ProductCategory{
			Name: "New Category",
		}

		// Expectations
		mockRepo.On("CreateCategory", ctx, category).Return(nil).Once()

		// Execute
		err := productService.CreateCategory(ctx, category)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Test data
		category := &domain.ProductCategory{
			Name: "New Category",
		}

		// Expectations
		mockRepo.On("CreateCategory", ctx, category).Return(errors.New("database error")).Once()

		// Execute
		err := productService.CreateCategory(ctx, category)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateCategory(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		category := &domain.ProductCategory{
			ID:   1,
			Name: "Updated Category",
		}

		// Expectations
		mockRepo.On("UpdateCategory", ctx, category).Return(nil).Once()

		// Execute
		err := productService.UpdateCategory(ctx, category)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Test data
		category := &domain.ProductCategory{
			ID:   1,
			Name: "Updated Category",
		}

		// Expectations
		mockRepo.On("UpdateCategory", ctx, category).Return(errors.New("database error")).Once()

		// Execute
		err := productService.UpdateCategory(ctx, category)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteCategory(t *testing.T) {
	// Setup
	mockRepo := new(MockProductRepository)
	productService := service.NewProductService(mockRepo)
	ctx := context.Background()
	categoryID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Expectations
		mockRepo.On("DeleteCategory", ctx, categoryID).Return(nil).Once()

		// Execute
		err := productService.DeleteCategory(ctx, categoryID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mock
		mockRepo = new(MockProductRepository)
		productService = service.NewProductService(mockRepo)

		// Expectations
		mockRepo.On("DeleteCategory", ctx, categoryID).Return(errors.New("database error")).Once()

		// Execute
		err := productService.DeleteCategory(ctx, categoryID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}
