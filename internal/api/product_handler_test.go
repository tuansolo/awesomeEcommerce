package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"awesomeEcommerce/internal/api"
	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock implementation of service.ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProductByID(id uint) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) GetProducts(page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetProductsByCategory(categoryID uint, page, pageSize int) ([]domain.Product, int64, error) {
	args := m.Called(categoryID, page, pageSize)
	return args.Get(0).([]domain.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) CreateProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) UpdateProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductService) DeleteProduct(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) GetProductBySKU(sku string) (*domain.Product, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) UpdateProductStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockProductService) GetCategories() ([]domain.ProductCategory, error) {
	args := m.Called()
	return args.Get(0).([]domain.ProductCategory), args.Error(1)
}

func (m *MockProductService) GetCategoryByID(id uint) (*domain.ProductCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductCategory), args.Error(1)
}

func (m *MockProductService) CreateCategory(category *domain.ProductCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockProductService) UpdateCategory(category *domain.ProductCategory) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockProductService) DeleteCategory(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockUserService is a mock implementation of service.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Authenticate(email, password string) (*domain.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Setup test router
func setupTestRouter(productService service.ProductService, userService service.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	productHandler := api.NewProductHandler(productService, userService)

	// Set up routes
	v1 := router.Group("/api/v1")
	productHandler.RegisterRoutes(v1)

	return router
}

// Test cases
func TestGetProduct(t *testing.T) {
	// Setup
	mockProductService := new(MockProductService)
	mockUserService := new(MockUserService)

	router := setupTestRouter(mockProductService, mockUserService)

	// Test data
	productID := uint(1)
	product := &domain.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "This is a test product",
		Price:       99.99,
		Stock:       100,
		SKU:         "TEST-001",
	}

	// Expectations
	mockProductService.On("GetProductByID", productID).Return(product, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/products/1", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(productID), data["id"])
	assert.Equal(t, product.Name, data["name"])
	assert.Equal(t, product.Description, data["description"])
	assert.Equal(t, product.Price, data["price"])
	assert.Equal(t, float64(product.Stock), data["stock"])
	assert.Equal(t, product.SKU, data["sku"])

	mockProductService.AssertExpectations(t)
}

func TestGetProducts(t *testing.T) {
	// Setup
	mockProductService := new(MockProductService)
	mockUserService := new(MockUserService)

	router := setupTestRouter(mockProductService, mockUserService)

	// Test data
	products := []domain.Product{
		{
			ID:          1,
			Name:        "Test Product 1",
			Description: "This is test product 1",
			Price:       99.99,
			Stock:       100,
			SKU:         "TEST-001",
		},
		{
			ID:          2,
			Name:        "Test Product 2",
			Description: "This is test product 2",
			Price:       199.99,
			Stock:       50,
			SKU:         "TEST-002",
		},
	}
	total := int64(2)

	// Expectations
	mockProductService.On("GetProducts", 1, 10).Return(products, total, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/products?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)

	meta, ok := response["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(total), meta["total"])

	mockProductService.AssertExpectations(t)
}

func TestCreateProduct(t *testing.T) {
	// Setup
	mockProductService := new(MockProductService)
	mockUserService := new(MockUserService)

	router := setupTestRouter(mockProductService, mockUserService)

	// Test data
	product := domain.Product{
		Name:        "New Test Product",
		Description: "This is a new test product",
		Price:       149.99,
		Stock:       75,
		SKU:         "TEST-NEW",
		CategoryID:  1,
	}

	// Expectations
	mockProductService.On("CreateProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

	// Create request
	jsonData, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	message, ok := response["message"].(string)
	assert.True(t, ok)
	assert.Equal(t, "Product created successfully", message)

	mockProductService.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	// Setup
	mockProductService := new(MockProductService)
	mockUserService := new(MockUserService)

	router := setupTestRouter(mockProductService, mockUserService)

	// Test data
	productID := uint(1)
	product := domain.Product{
		ID:          productID,
		Name:        "Updated Test Product",
		Description: "This is an updated test product",
		Price:       199.99,
		Stock:       50,
		SKU:         "TEST-UPD",
		CategoryID:  1,
	}

	// Expectations
	mockProductService.On("UpdateProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

	// Create request
	jsonData, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/api/v1/products/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	message, ok := response["message"].(string)
	assert.True(t, ok)
	assert.Equal(t, "Product updated successfully", message)

	mockProductService.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	// Setup
	mockProductService := new(MockProductService)
	mockUserService := new(MockUserService)

	router := setupTestRouter(mockProductService, mockUserService)

	// Test data
	productID := uint(1)

	// Expectations
	mockProductService.On("DeleteProduct", productID).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/products/1", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	message, ok := response["message"].(string)
	assert.True(t, ok)
	assert.Equal(t, "Product deleted successfully", message)

	mockProductService.AssertExpectations(t)
}
