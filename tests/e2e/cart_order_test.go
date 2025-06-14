package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartAndOrderOperations(t *testing.T) {
	// Setup test environment
	client := SetupTest(t)
	defer TeardownTest(t)

	var productID uint
	var categoryID uint
	var orderID uint

	// Create a test user
	t.Run("Register User", func(t *testing.T) {
		token := RegisterTestUser(t, client)
		client.SetAuthToken(token)
	})

	// Create a test category and product (as admin)
	t.Run("Create Test Product", func(t *testing.T) {
		// Login as admin
		token := CreateTestAdmin(t)
		client.SetAuthToken(token)

		// Create category
		categoryReqBody := map[string]interface{}{
			"name":        "Test Category",
			"description": "A category for testing",
		}

		resp, err := client.DoRequest(http.MethodPost, "/products/categories", categoryReqBody)
		require.NoError(t, err, "Failed to create category")

		var categoryRespBody struct {
			Category struct {
				ID uint `json:"id"`
			} `json:"category"`
		}

		err = ParseResponse(resp, &categoryRespBody)
		require.NoError(t, err, "Failed to parse category creation response")
		categoryID = categoryRespBody.Category.ID

		// Create product
		productReqBody := map[string]interface{}{
			"name":        "Test Product",
			"description": "A product for testing",
			"price":       19.99,
			"stock":       100,
			"sku":         "TEST-SKU-123",
			"image_url":   "https://example.com/test-product.jpg",
			"category_id": categoryID,
		}

		resp, err = client.DoRequest(http.MethodPost, "/products", productReqBody)
		require.NoError(t, err, "Failed to create product")

		var productRespBody struct {
			Product struct {
				ID uint `json:"id"`
			} `json:"product"`
		}

		err = ParseResponse(resp, &productRespBody)
		require.NoError(t, err, "Failed to parse product creation response")
		productID = productRespBody.Product.ID
	})

	// Login as regular user for cart operations
	t.Run("Login User", func(t *testing.T) {
		token := RegisterTestUser(t, client)
		client.SetAuthToken(token)
	})

	// Test adding an item to the cart
	t.Run("Add Item To Cart", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"product_id": productID,
			"quantity":   2,
		}

		resp, err := client.DoRequest(http.MethodPost, "/carts/items", reqBody)
		require.NoError(t, err, "Failed to add item to cart")

		var respBody struct {
			Message string `json:"message"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse add to cart response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Item added to cart successfully", respBody.Message)
	})

	// Test getting the cart
	t.Run("Get Cart", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/carts/me", nil)
		require.NoError(t, err, "Failed to get cart")

		var respBody struct {
			Cart struct {
				ID     uint `json:"id"`
				UserID uint `json:"user_id"`
				Items  []struct {
					ID        uint `json:"id"`
					ProductID uint `json:"product_id"`
					Quantity  int  `json:"quantity"`
					Product   struct {
						ID          uint    `json:"id"`
						Name        string  `json:"name"`
						Description string  `json:"description"`
						Price       float64 `json:"price"`
						ImageURL    string  `json:"image_url"`
					} `json:"product"`
				} `json:"items"`
			} `json:"cart"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse get cart response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.NotZero(t, respBody.Cart.ID, "Cart ID should not be zero")
		assert.NotZero(t, respBody.Cart.UserID, "User ID should not be zero")
		assert.Len(t, respBody.Cart.Items, 1, "Cart should have 1 item")
		assert.Equal(t, productID, respBody.Cart.Items[0].ProductID)
		assert.Equal(t, 2, respBody.Cart.Items[0].Quantity)
		assert.Equal(t, "Test Product", respBody.Cart.Items[0].Product.Name)
		assert.Equal(t, 19.99, respBody.Cart.Items[0].Product.Price)
	})

	// Test updating a cart item
	t.Run("Update Cart Item", func(t *testing.T) {
		// First get the cart to get the item ID
		resp, err := client.DoRequest(http.MethodGet, "/carts/me", nil)
		require.NoError(t, err, "Failed to get cart")

		var getCartResp struct {
			Cart struct {
				Items []struct {
					ID uint `json:"id"`
				} `json:"items"`
			} `json:"cart"`
		}

		err = ParseResponse(resp, &getCartResp)
		require.NoError(t, err, "Failed to parse get cart response")
		require.NotEmpty(t, getCartResp.Cart.Items, "Cart should have items")

		itemID := getCartResp.Cart.Items[0].ID

		// Update the cart item
		reqBody := map[string]interface{}{
			"quantity": 5,
		}

		resp, err = client.DoRequest(http.MethodPut, "/carts/items/"+strconv.FormatUint(uint64(itemID), 10), reqBody)
		require.NoError(t, err, "Failed to update cart item")

		var respBody struct {
			Message string `json:"message"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse update cart item response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Cart item updated successfully", respBody.Message)

		// Verify the update by getting the cart again
		resp, err = client.DoRequest(http.MethodGet, "/carts/me", nil)
		require.NoError(t, err, "Failed to get cart after update")

		var getCartAfterUpdateResp struct {
			Cart struct {
				Items []struct {
					ID       uint `json:"id"`
					Quantity int  `json:"quantity"`
				} `json:"items"`
			} `json:"cart"`
		}

		err = ParseResponse(resp, &getCartAfterUpdateResp)
		require.NoError(t, err, "Failed to parse get cart response after update")
		require.NotEmpty(t, getCartAfterUpdateResp.Cart.Items, "Cart should have items after update")
		assert.Equal(t, 5, getCartAfterUpdateResp.Cart.Items[0].Quantity, "Quantity should be updated to 5")
	})

	// Test getting the cart total
	t.Run("Get Cart Total", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/carts/total", nil)
		require.NoError(t, err, "Failed to get cart total")

		var respBody struct {
			Total float64 `json:"total"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse get cart total response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, 19.99*5, respBody.Total, "Total should be price * quantity")
	})

	// Test creating an order from the cart
	t.Run("Create Order", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"shipping_address": "123 Shipping St, Test City",
			"billing_address":  "123 Billing St, Test City",
		}

		resp, err := client.DoRequest(http.MethodPost, "/orders", reqBody)
		require.NoError(t, err, "Failed to create order")

		var respBody struct {
			Message string `json:"message"`
			Order   struct {
				ID              uint    `json:"id"`
				UserID          uint    `json:"user_id"`
				TotalAmount     float64 `json:"total_amount"`
				Status          string  `json:"status"`
				ShippingAddress string  `json:"shipping_address"`
				BillingAddress  string  `json:"billing_address"`
				Items           []struct {
					ProductID   uint    `json:"product_id"`
					ProductName string  `json:"product_name"`
					Price       float64 `json:"price"`
					Quantity    int     `json:"quantity"`
				} `json:"items"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse create order response")

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
		assert.Equal(t, "Order created successfully", respBody.Message)
		assert.NotZero(t, respBody.Order.ID, "Order ID should not be zero")
		assert.Equal(t, 19.99*5, respBody.Order.TotalAmount, "Total amount should be price * quantity")
		assert.Equal(t, "pending", respBody.Order.Status)
		assert.Equal(t, "123 Shipping St, Test City", respBody.Order.ShippingAddress)
		assert.Equal(t, "123 Billing St, Test City", respBody.Order.BillingAddress)
		assert.Len(t, respBody.Order.Items, 1, "Order should have 1 item")
		assert.Equal(t, productID, respBody.Order.Items[0].ProductID)
		assert.Equal(t, "Test Product", respBody.Order.Items[0].ProductName)
		assert.Equal(t, 19.99, respBody.Order.Items[0].Price)
		assert.Equal(t, 5, respBody.Order.Items[0].Quantity)

		// Save order ID for later tests
		orderID = respBody.Order.ID

		// Verify cart is empty after order creation
		resp, err = client.DoRequest(http.MethodGet, "/carts/me", nil)
		require.NoError(t, err, "Failed to get cart after order creation")

		var getCartAfterOrderResp struct {
			Cart struct {
				Items []struct{} `json:"items"`
			} `json:"cart"`
		}

		err = ParseResponse(resp, &getCartAfterOrderResp)
		require.NoError(t, err, "Failed to parse get cart response after order creation")
		assert.Empty(t, getCartAfterOrderResp.Cart.Items, "Cart should be empty after order creation")
	})

	// Test getting user's orders
	t.Run("Get My Orders", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/orders/me", nil)
		require.NoError(t, err, "Failed to get user orders")

		var respBody struct {
			Orders []struct {
				ID          uint    `json:"id"`
				TotalAmount float64 `json:"total_amount"`
				Status      string  `json:"status"`
			} `json:"orders"`
			Meta struct {
				Total    int `json:"total"`
				Page     int `json:"page"`
				PageSize int `json:"page_size"`
			} `json:"meta"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse get user orders response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.GreaterOrEqual(t, respBody.Meta.Total, 1, "User should have at least 1 order")

		// Find our test order
		var found bool
		for _, order := range respBody.Orders {
			if order.ID == orderID {
				found = true
				assert.Equal(t, 19.99*5, order.TotalAmount)
				assert.Equal(t, "pending", order.Status)
				break
			}
		}
		assert.True(t, found, "Test order not found in user orders")
	})

	// Test getting a specific order
	t.Run("Get My Order By ID", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/orders/me/"+strconv.FormatUint(uint64(orderID), 10), nil)
		require.NoError(t, err, "Failed to get order by ID")

		var respBody struct {
			Order struct {
				ID              uint    `json:"id"`
				UserID          uint    `json:"user_id"`
				TotalAmount     float64 `json:"total_amount"`
				Status          string  `json:"status"`
				ShippingAddress string  `json:"shipping_address"`
				BillingAddress  string  `json:"billing_address"`
				Items           []struct {
					ProductID   uint    `json:"product_id"`
					ProductName string  `json:"product_name"`
					Price       float64 `json:"price"`
					Quantity    int     `json:"quantity"`
				} `json:"items"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse get order by ID response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, orderID, respBody.Order.ID)
		assert.Equal(t, 19.99*5, respBody.Order.TotalAmount)
		assert.Equal(t, "pending", respBody.Order.Status)
		assert.Equal(t, "123 Shipping St, Test City", respBody.Order.ShippingAddress)
		assert.Equal(t, "123 Billing St, Test City", respBody.Order.BillingAddress)
		assert.Len(t, respBody.Order.Items, 1, "Order should have 1 item")
		assert.Equal(t, productID, respBody.Order.Items[0].ProductID)
		assert.Equal(t, "Test Product", respBody.Order.Items[0].ProductName)
		assert.Equal(t, 19.99, respBody.Order.Items[0].Price)
		assert.Equal(t, 5, respBody.Order.Items[0].Quantity)
	})

	// Test cancelling an order
	t.Run("Cancel Order", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodPost, "/orders/me/"+strconv.FormatUint(uint64(orderID), 10)+"/cancel", nil)
		require.NoError(t, err, "Failed to cancel order")

		var respBody struct {
			Message string `json:"message"`
			Order   struct {
				ID     uint   `json:"id"`
				Status string `json:"status"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse cancel order response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Order cancelled successfully", respBody.Message)
		assert.Equal(t, orderID, respBody.Order.ID)
		assert.Equal(t, "cancelled", respBody.Order.Status)

		// Verify order status by getting the order
		resp, err = client.DoRequest(http.MethodGet, "/orders/me/"+strconv.FormatUint(uint64(orderID), 10), nil)
		require.NoError(t, err, "Failed to get order after cancellation")

		var getOrderAfterCancelResp struct {
			Order struct {
				Status string `json:"status"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &getOrderAfterCancelResp)
		require.NoError(t, err, "Failed to parse get order response after cancellation")
		assert.Equal(t, "cancelled", getOrderAfterCancelResp.Order.Status, "Order status should be cancelled")
	})

	// Clean up - delete the test product and category
	t.Run("Cleanup", func(t *testing.T) {
		// Login as admin
		token := CreateTestAdmin(t)
		client.SetAuthToken(token)

		// Delete product
		_, err := client.DoRequest(http.MethodDelete, "/products/"+strconv.FormatUint(uint64(productID), 10), nil)
		require.NoError(t, err, "Failed to delete product during cleanup")

		// Delete category
		_, err = client.DoRequest(http.MethodDelete, "/products/categories/"+strconv.FormatUint(uint64(categoryID), 10), nil)
		require.NoError(t, err, "Failed to delete category during cleanup")
	})
}
