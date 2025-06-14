package impl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) FindByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) FindCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ProductCategory), args.Error(1)
}

func (m *MockProductRepository) FindCategoryByID(ctx context.Context, id uint) (*domain.ProductCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *MockProductRepository) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteCategory(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRedisClient is a mock implementation of the RedisClient
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// ProductRepositoryTestSuite is a test suite for ProductRepository
type ProductRepositoryTestSuite struct {
	suite.Suite
	mockRepo  repository.ProductRepository
	mockCache *MockRedisClient
	ctx       context.Context
}

// SetupTest sets up the test suite
func (s *ProductRepositoryTestSuite) SetupTest() {
	s.mockRepo = new(MockProductRepository)
	s.mockCache = new(MockRedisClient)
	s.ctx = context.Background()
}

// TestFindByID tests the FindByID method
func (s *ProductRepositoryTestSuite) TestFindByID() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a product by ID
		productID := uint(1)
		expectedProduct := &domain.Product{
			ID:          productID,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			SKU:         "TEST-SKU-123",
			ImageURL:    "http://example.com/image.jpg",
			CategoryID:  uint(1),
		}

		mockRepo.On("FindByID", s.ctx, productID).Return(expectedProduct, nil).Once()

		// Execute
		product, err := s.mockRepo.FindByID(s.ctx, productID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedProduct, product)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Product not found
		productID := uint(999)
		expectedError := errors.New("product not found")

		mockRepo.On("FindByID", s.ctx, productID).Return(nil, expectedError).Once()

		// Execute
		product, err := s.mockRepo.FindByID(s.ctx, productID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), product)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindAll tests the FindAll method
func (s *ProductRepositoryTestSuite) TestFindAll() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find all products
		page := 1
		pageSize := 10
		expectedProducts := []domain.Product{
			{
				ID:          uint(1),
				Name:        "Test Product 1",
				Description: "This is test product 1",
				Price:       99.99,
				Stock:       100,
				SKU:         "TEST-SKU-1",
				ImageURL:    "http://example.com/image1.jpg",
				CategoryID:  uint(1),
			},
			{
				ID:          uint(2),
				Name:        "Test Product 2",
				Description: "This is test product 2",
				Price:       49.99,
				Stock:       50,
				SKU:         "TEST-SKU-2",
				ImageURL:    "http://example.com/image2.jpg",
				CategoryID:  uint(2),
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedProducts, products)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Products", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: No products in the system
		page := 1
		pageSize := 10
		var expectedProducts []domain.Product
		expectedTotal := int64(0)

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), products)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Database error
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		products, total, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), products)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByCategory tests the FindByCategory method
func (s *ProductRepositoryTestSuite) TestFindByCategory() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find products by category
		categoryID := uint(1)
		page := 1
		pageSize := 10
		expectedProducts := []domain.Product{
			{
				ID:          uint(1),
				Name:        "Test Product 1",
				Description: "This is test product 1",
				Price:       99.99,
				Stock:       100,
				SKU:         "TEST-SKU-1",
				ImageURL:    "http://example.com/image1.jpg",
				CategoryID:  categoryID,
			},
			{
				ID:          uint(3),
				Name:        "Test Product 3",
				Description: "This is test product 3",
				Price:       79.99,
				Stock:       75,
				SKU:         "TEST-SKU-3",
				ImageURL:    "http://example.com/image3.jpg",
				CategoryID:  categoryID,
			},
		}
		expectedTotal := int64(2)

		mockRepo.On("FindByCategory", s.ctx, categoryID, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := s.mockRepo.FindByCategory(s.ctx, categoryID, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedProducts, products)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Products In Category", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: No products in the category
		categoryID := uint(3)
		page := 1
		pageSize := 10
		var expectedProducts []domain.Product
		expectedTotal := int64(0)

		mockRepo.On("FindByCategory", s.ctx, categoryID, page, pageSize).Return(expectedProducts, expectedTotal, nil).Once()

		// Execute
		products, total, err := s.mockRepo.FindByCategory(s.ctx, categoryID, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), products)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Database error
		categoryID := uint(1)
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByCategory", s.ctx, categoryID, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		products, total, err := s.mockRepo.FindByCategory(s.ctx, categoryID, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), products)
		assert.Equal(s.T(), int64(0), total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreate tests the Create method
func (s *ProductRepositoryTestSuite) TestCreate() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully create a product
		product := &domain.Product{
			Name:        "New Test Product",
			Description: "This is a new test product",
			Price:       129.99,
			Stock:       200,
			SKU:         "NEW-TEST-SKU",
			ImageURL:    "http://example.com/new-image.jpg",
			CategoryID:  uint(1),
		}

		mockRepo.On("Create", s.ctx, product).Return(nil).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, product)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Database error
		product := &domain.Product{
			Name:        "New Test Product",
			Description: "This is a new test product",
			Price:       129.99,
			Stock:       200,
			SKU:         "NEW-TEST-SKU",
			ImageURL:    "http://example.com/new-image.jpg",
			CategoryID:  uint(1),
		}
		expectedError := errors.New("database error")

		mockRepo.On("Create", s.ctx, product).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, product)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdate tests the Update method
func (s *ProductRepositoryTestSuite) TestUpdate() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a product
		product := &domain.Product{
			ID:          uint(1),
			Name:        "Updated Test Product",
			Description: "This is an updated test product",
			Price:       149.99,
			Stock:       150,
			SKU:         "TEST-SKU-1",
			ImageURL:    "http://example.com/updated-image.jpg",
			CategoryID:  uint(2),
		}

		mockRepo.On("Update", s.ctx, product).Return(nil).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, product)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Product Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Product not found
		product := &domain.Product{
			ID:          uint(999),
			Name:        "Updated Test Product",
			Description: "This is an updated test product",
			Price:       149.99,
			Stock:       150,
			SKU:         "TEST-SKU-999",
			ImageURL:    "http://example.com/updated-image.jpg",
			CategoryID:  uint(2),
		}
		expectedError := errors.New("product not found")

		mockRepo.On("Update", s.ctx, product).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, product)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDelete tests the Delete method
func (s *ProductRepositoryTestSuite) TestDelete() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete a product
		productID := uint(1)

		mockRepo.On("Delete", s.ctx, productID).Return(nil).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, productID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Product Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Product not found
		productID := uint(999)
		expectedError := errors.New("product not found")

		mockRepo.On("Delete", s.ctx, productID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, productID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindBySKU tests the FindBySKU method
func (s *ProductRepositoryTestSuite) TestFindBySKU() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a product by SKU
		sku := "TEST-SKU-123"
		expectedProduct := &domain.Product{
			ID:          uint(1),
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			SKU:         sku,
			ImageURL:    "http://example.com/image.jpg",
			CategoryID:  uint(1),
		}

		mockRepo.On("FindBySKU", s.ctx, sku).Return(expectedProduct, nil).Once()

		// Execute
		product, err := s.mockRepo.FindBySKU(s.ctx, sku)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedProduct, product)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Product not found
		sku := "NONEXISTENT-SKU"
		expectedError := errors.New("product not found")

		mockRepo.On("FindBySKU", s.ctx, sku).Return(nil, expectedError).Once()

		// Execute
		product, err := s.mockRepo.FindBySKU(s.ctx, sku)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), product)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdateStock tests the UpdateStock method
func (s *ProductRepositoryTestSuite) TestUpdateStock() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success - Increase Stock", func() {
		// Test case: Successfully increase product stock
		productID := uint(1)
		quantity := 50

		mockRepo.On("UpdateStock", s.ctx, productID, quantity).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateStock(s.ctx, productID, quantity)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Decrease Stock", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Successfully decrease product stock
		productID := uint(1)
		quantity := -20

		mockRepo.On("UpdateStock", s.ctx, productID, quantity).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateStock(s.ctx, productID, quantity)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Product Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Product not found
		productID := uint(999)
		quantity := 10
		expectedError := errors.New("product not found")

		mockRepo.On("UpdateStock", s.ctx, productID, quantity).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateStock(s.ctx, productID, quantity)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Insufficient Stock", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Insufficient stock for decrease
		productID := uint(1)
		quantity := -200
		expectedError := errors.New("insufficient stock")

		mockRepo.On("UpdateStock", s.ctx, productID, quantity).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateStock(s.ctx, productID, quantity)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindCategories tests the FindCategories method
func (s *ProductRepositoryTestSuite) TestFindCategories() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find all product categories
		expectedCategories := []domain.ProductCategory{
			{
				ID:       uint(1),
				Name:     "Category 1",
				ParentID: nil,
			},
			{
				ID:       uint(2),
				Name:     "Category 2",
				ParentID: nil,
			},
			{
				ID:       uint(3),
				Name:     "Subcategory 1",
				ParentID: func() *uint { id := uint(1); return &id }(),
			},
		}

		mockRepo.On("FindCategories", s.ctx).Return(expectedCategories, nil).Once()

		// Execute
		categories, err := s.mockRepo.FindCategories(s.ctx)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedCategories, categories)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Categories", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: No categories in the system
		var expectedCategories []domain.ProductCategory

		mockRepo.On("FindCategories", s.ctx).Return(expectedCategories, nil).Once()

		// Execute
		categories, err := s.mockRepo.FindCategories(s.ctx)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), categories)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Database error
		expectedError := errors.New("database error")

		mockRepo.On("FindCategories", s.ctx).Return(nil, expectedError).Once()

		// Execute
		categories, err := s.mockRepo.FindCategories(s.ctx)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), categories)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindCategoryByID tests the FindCategoryByID method
func (s *ProductRepositoryTestSuite) TestFindCategoryByID() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a category by ID
		categoryID := uint(1)
		expectedCategory := &domain.ProductCategory{
			ID:       categoryID,
			Name:     "Category 1",
			ParentID: nil,
		}

		mockRepo.On("FindCategoryByID", s.ctx, categoryID).Return(expectedCategory, nil).Once()

		// Execute
		category, err := s.mockRepo.FindCategoryByID(s.ctx, categoryID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedCategory, category)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Category not found
		categoryID := uint(999)
		expectedError := errors.New("category not found")

		mockRepo.On("FindCategoryByID", s.ctx, categoryID).Return(nil, expectedError).Once()

		// Execute
		category, err := s.mockRepo.FindCategoryByID(s.ctx, categoryID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), category)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreateCategory tests the CreateCategory method
func (s *ProductRepositoryTestSuite) TestCreateCategory() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully create a category
		category := &domain.ProductCategory{
			Name:     "New Category",
			ParentID: nil,
		}

		mockRepo.On("CreateCategory", s.ctx, category).Return(nil).Once()

		// Execute
		err := s.mockRepo.CreateCategory(s.ctx, category)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - With Parent", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Successfully create a subcategory
		parentID := uint(1)
		category := &domain.ProductCategory{
			Name:     "New Subcategory",
			ParentID: &parentID,
		}

		mockRepo.On("CreateCategory", s.ctx, category).Return(nil).Once()

		// Execute
		err := s.mockRepo.CreateCategory(s.ctx, category)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Database error
		category := &domain.ProductCategory{
			Name:     "New Category",
			ParentID: nil,
		}
		expectedError := errors.New("database error")

		mockRepo.On("CreateCategory", s.ctx, category).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.CreateCategory(s.ctx, category)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdateCategory tests the UpdateCategory method
func (s *ProductRepositoryTestSuite) TestUpdateCategory() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a category
		category := &domain.ProductCategory{
			ID:       uint(1),
			Name:     "Updated Category",
			ParentID: nil,
		}

		mockRepo.On("UpdateCategory", s.ctx, category).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateCategory(s.ctx, category)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Category Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Category not found
		category := &domain.ProductCategory{
			ID:       uint(999),
			Name:     "Updated Category",
			ParentID: nil,
		}
		expectedError := errors.New("category not found")

		mockRepo.On("UpdateCategory", s.ctx, category).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateCategory(s.ctx, category)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDeleteCategory tests the DeleteCategory method
func (s *ProductRepositoryTestSuite) TestDeleteCategory() {
	mockRepo := s.mockRepo.(*MockProductRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete a category
		categoryID := uint(1)

		mockRepo.On("DeleteCategory", s.ctx, categoryID).Return(nil).Once()

		// Execute
		err := s.mockRepo.DeleteCategory(s.ctx, categoryID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Category Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Category not found
		categoryID := uint(999)
		expectedError := errors.New("category not found")

		mockRepo.On("DeleteCategory", s.ctx, categoryID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.DeleteCategory(s.ctx, categoryID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Category Has Products", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockProductRepository)

		// Test case: Category has associated products
		categoryID := uint(2)
		expectedError := errors.New("cannot delete category with associated products")

		mockRepo.On("DeleteCategory", s.ctx, categoryID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.DeleteCategory(s.ctx, categoryID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestProductRepositorySuite runs the test suite
func TestProductRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProductRepositoryTestSuite))
}
