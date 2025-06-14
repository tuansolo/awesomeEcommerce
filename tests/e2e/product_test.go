package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductOperations(t *testing.T) {
	// Setup test environment
	client := SetupTest(t)
	defer TeardownTest(t)

	var productID uint
	var categoryID uint

	// Test creating a category (requires admin)
	t.Run("Create Category", func(t *testing.T) {
		// Register and login as admin (this is a placeholder - in a real test, you would need to create an admin user)
		token := CreateTestAdmin(t)
		client.SetAuthToken(token)

		reqBody := map[string]interface{}{
			"name":        "Test Category",
			"description": "A category for testing",
		}

		resp, err := client.DoRequest(http.MethodPost, "/products/categories", reqBody)
		require.NoError(t, err, "Failed to create category")

		var respBody struct {
			Message  string `json:"message"`
			Category struct {
				ID          uint   `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"category"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse category creation response")

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
		assert.Equal(t, "Category created successfully", respBody.Message)
		assert.Equal(t, "Test Category", respBody.Category.Name)
		assert.Equal(t, "A category for testing", respBody.Category.Description)

		// Save category ID for later tests
		categoryID = respBody.Category.ID
	})

	// Test creating a product (requires admin)
	t.Run("Create Product", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "Test Product",
			"description": "A product for testing",
			"price":       19.99,
			"stock":       100,
			"sku":         "TEST-SKU-123",
			"image_url":   "https://example.com/test-product.jpg",
			"category_id": categoryID,
		}

		resp, err := client.DoRequest(http.MethodPost, "/products", reqBody)
		require.NoError(t, err, "Failed to create product")

		var respBody struct {
			Message string `json:"message"`
			Product struct {
				ID          uint    `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				SKU         string  `json:"sku"`
				ImageURL    string  `json:"image_url"`
				CategoryID  uint    `json:"category_id"`
			} `json:"product"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse product creation response")

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
		assert.Equal(t, "Product created successfully", respBody.Message)
		assert.Equal(t, "Test Product", respBody.Product.Name)
		assert.Equal(t, "A product for testing", respBody.Product.Description)
		assert.Equal(t, 19.99, respBody.Product.Price)
		assert.Equal(t, 100, respBody.Product.Stock)
		assert.Equal(t, "TEST-SKU-123", respBody.Product.SKU)
		assert.Equal(t, "https://example.com/test-product.jpg", respBody.Product.ImageURL)
		assert.Equal(t, categoryID, respBody.Product.CategoryID)

		// Save product ID for later tests
		productID = respBody.Product.ID
	})

	// Test getting all products (public)
	t.Run("Get Products", func(t *testing.T) {
		// Clear auth token to test as anonymous user
		client.SetAuthToken("")

		resp, err := client.DoRequest(http.MethodGet, "/products", nil)
		require.NoError(t, err, "Failed to get products")

		var respBody struct {
			Products []struct {
				ID          uint    `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				SKU         string  `json:"sku"`
				ImageURL    string  `json:"image_url"`
				CategoryID  uint    `json:"category_id"`
			} `json:"products"`
			Meta struct {
				Total    int `json:"total"`
				Page     int `json:"page"`
				PageSize int `json:"page_size"`
			} `json:"meta"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse products response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.GreaterOrEqual(t, respBody.Meta.Total, 1, "Expected at least 1 product")

		// Find our test product
		var found bool
		for _, product := range respBody.Products {
			if product.ID == productID {
				found = true
				assert.Equal(t, "Test Product", product.Name)
				assert.Equal(t, "A product for testing", product.Description)
				assert.Equal(t, 19.99, product.Price)
				assert.Equal(t, 100, product.Stock)
				assert.Equal(t, "TEST-SKU-123", product.SKU)
				assert.Equal(t, "https://example.com/test-product.jpg", product.ImageURL)
				assert.Equal(t, categoryID, product.CategoryID)
				break
			}
		}
		assert.True(t, found, "Test product not found in products list")
	})

	// Test getting a specific product (public)
	t.Run("Get Product By ID", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/products/"+strconv.FormatUint(uint64(productID), 10), nil)
		require.NoError(t, err, "Failed to get product by ID")

		var respBody struct {
			Product struct {
				ID          uint    `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				SKU         string  `json:"sku"`
				ImageURL    string  `json:"image_url"`
				CategoryID  uint    `json:"category_id"`
			} `json:"product"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse product response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, productID, respBody.Product.ID)
		assert.Equal(t, "Test Product", respBody.Product.Name)
		assert.Equal(t, "A product for testing", respBody.Product.Description)
		assert.Equal(t, 19.99, respBody.Product.Price)
		assert.Equal(t, 100, respBody.Product.Stock)
		assert.Equal(t, "TEST-SKU-123", respBody.Product.SKU)
		assert.Equal(t, "https://example.com/test-product.jpg", respBody.Product.ImageURL)
		assert.Equal(t, categoryID, respBody.Product.CategoryID)
	})

	// Test getting all categories (public)
	t.Run("Get Categories", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/products/categories", nil)
		require.NoError(t, err, "Failed to get categories")

		var respBody struct {
			Categories []struct {
				ID          uint   `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"categories"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse categories response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.GreaterOrEqual(t, len(respBody.Categories), 1, "Expected at least 1 category")

		// Find our test category
		var found bool
		for _, category := range respBody.Categories {
			if category.ID == categoryID {
				found = true
				assert.Equal(t, "Test Category", category.Name)
				assert.Equal(t, "A category for testing", category.Description)
				break
			}
		}
		assert.True(t, found, "Test category not found in categories list")
	})

	// Test getting products by category (public)
	t.Run("Get Products By Category", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/products/categories/"+strconv.FormatUint(uint64(categoryID), 10), nil)
		require.NoError(t, err, "Failed to get products by category")

		var respBody struct {
			Products []struct {
				ID          uint    `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				SKU         string  `json:"sku"`
				ImageURL    string  `json:"image_url"`
				CategoryID  uint    `json:"category_id"`
			} `json:"products"`
			Category struct {
				ID          uint   `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"category"`
			Meta struct {
				Total    int `json:"total"`
				Page     int `json:"page"`
				PageSize int `json:"page_size"`
			} `json:"meta"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse products by category response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, categoryID, respBody.Category.ID)
		assert.Equal(t, "Test Category", respBody.Category.Name)
		assert.Equal(t, "A category for testing", respBody.Category.Description)
		assert.GreaterOrEqual(t, respBody.Meta.Total, 1, "Expected at least 1 product in category")

		// Find our test product
		var found bool
		for _, product := range respBody.Products {
			if product.ID == productID {
				found = true
				assert.Equal(t, "Test Product", product.Name)
				assert.Equal(t, categoryID, product.CategoryID)
				break
			}
		}
		assert.True(t, found, "Test product not found in category products list")
	})

	// Test updating a product (requires admin)
	t.Run("Update Product", func(t *testing.T) {
		// Login as admin again
		token := CreateTestAdmin(t)
		client.SetAuthToken(token)

		reqBody := map[string]interface{}{
			"name":        "Updated Test Product",
			"description": "An updated product for testing",
			"price":       29.99,
			"stock":       50,
		}

		resp, err := client.DoRequest(http.MethodPut, "/products/"+strconv.FormatUint(uint64(productID), 10), reqBody)
		require.NoError(t, err, "Failed to update product")

		var respBody struct {
			Message string `json:"message"`
			Product struct {
				ID          uint    `json:"id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				SKU         string  `json:"sku"`
				ImageURL    string  `json:"image_url"`
				CategoryID  uint    `json:"category_id"`
			} `json:"product"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse product update response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Product updated successfully", respBody.Message)
		assert.Equal(t, productID, respBody.Product.ID)
		assert.Equal(t, "Updated Test Product", respBody.Product.Name)
		assert.Equal(t, "An updated product for testing", respBody.Product.Description)
		assert.Equal(t, 29.99, respBody.Product.Price)
		assert.Equal(t, 50, respBody.Product.Stock)
		assert.Equal(t, "TEST-SKU-123", respBody.Product.SKU) // SKU should not change
		assert.Equal(t, "https://example.com/test-product.jpg", respBody.Product.ImageURL)
		assert.Equal(t, categoryID, respBody.Product.CategoryID)
	})

	// Test deleting a product (requires admin)
	t.Run("Delete Product", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodDelete, "/products/"+strconv.FormatUint(uint64(productID), 10), nil)
		require.NoError(t, err, "Failed to delete product")

		var respBody struct {
			Message string `json:"message"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse product deletion response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Product deleted successfully", respBody.Message)

		// Verify product is deleted by trying to get it
		getResp, _ := client.DoRequest(http.MethodGet, "/products/"+strconv.FormatUint(uint64(productID), 10), nil)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode, "Expected status code 404 after deletion")
	})

	// Test deleting a category (requires admin)
	t.Run("Delete Category", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodDelete, "/products/categories/"+strconv.FormatUint(uint64(categoryID), 10), nil)
		require.NoError(t, err, "Failed to delete category")

		var respBody struct {
			Message string `json:"message"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse category deletion response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Category deleted successfully", respBody.Message)

		// Verify category is deleted by trying to get products by category
		getResp, _ := client.DoRequest(http.MethodGet, "/products/categories/"+strconv.FormatUint(uint64(categoryID), 10), nil)
		assert.Equal(t, http.StatusNotFound, getResp.StatusCode, "Expected status code 404 after deletion")
	})
}
