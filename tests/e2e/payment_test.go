package e2e

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentOperations(t *testing.T) {
	// Setup test environment
	client := SetupTest(t)
	defer TeardownTest(t)

	var productID uint
	var categoryID uint
	var orderID uint
	var paymentID uint

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

	// Login as regular user for cart and order operations
	t.Run("Login User", func(t *testing.T) {
		token := RegisterTestUser(t, client)
		client.SetAuthToken(token)
	})

	// Create a cart and order for payment testing
	t.Run("Create Order for Payment", func(t *testing.T) {
		// Add item to cart
		addItemReqBody := map[string]interface{}{
			"product_id": productID,
			"quantity":   2,
		}

		resp, err := client.DoRequest(http.MethodPost, "/carts/items", addItemReqBody)
		require.NoError(t, err, "Failed to add item to cart")
		require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

		// Create order
		orderReqBody := map[string]interface{}{
			"shipping_address": "123 Shipping St, Test City",
			"billing_address":  "123 Billing St, Test City",
		}

		resp, err = client.DoRequest(http.MethodPost, "/orders", orderReqBody)
		require.NoError(t, err, "Failed to create order")

		var orderRespBody struct {
			Order struct {
				ID uint `json:"id"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &orderRespBody)
		require.NoError(t, err, "Failed to parse create order response")
		orderID = orderRespBody.Order.ID
	})

	// Test creating a payment for an order
	t.Run("Create Payment", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"amount": 39.98, // 19.99 * 2
			"method": "credit_card",
		}

		resp, err := client.DoRequest(http.MethodPost, "/payments/orders/"+strconv.FormatUint(uint64(orderID), 10), reqBody)
		require.NoError(t, err, "Failed to create payment")

		var respBody struct {
			Message string `json:"message"`
			Payment struct {
				ID            uint    `json:"id"`
				OrderID       uint    `json:"order_id"`
				Amount        float64 `json:"amount"`
				Currency      string  `json:"currency"`
				Method        string  `json:"method"`
				Status        string  `json:"status"`
				TransactionID string  `json:"transaction_id"`
				PaymentDate   string  `json:"payment_date"`
			} `json:"payment"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse payment creation response")

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
		assert.Equal(t, "Payment processed successfully", respBody.Message)
		assert.NotZero(t, respBody.Payment.ID, "Payment ID should not be zero")
		assert.Equal(t, orderID, respBody.Payment.OrderID)
		assert.Equal(t, 39.98, respBody.Payment.Amount)
		assert.Equal(t, "USD", respBody.Payment.Currency) // Assuming USD is the default
		assert.Equal(t, "credit_card", respBody.Payment.Method)
		assert.Equal(t, "completed", respBody.Payment.Status)
		assert.NotEmpty(t, respBody.Payment.TransactionID, "Transaction ID should not be empty")
		assert.NotEmpty(t, respBody.Payment.PaymentDate, "Payment date should not be empty")

		// Save payment ID for later tests
		paymentID = respBody.Payment.ID

		// Verify order status has been updated
		resp, err = client.DoRequest(http.MethodGet, "/orders/me/"+strconv.FormatUint(uint64(orderID), 10), nil)
		require.NoError(t, err, "Failed to get order after payment")

		var getOrderResp struct {
			Order struct {
				Status string `json:"status"`
			} `json:"order"`
		}

		err = ParseResponse(resp, &getOrderResp)
		require.NoError(t, err, "Failed to parse get order response after payment")
		assert.Equal(t, "paid", getOrderResp.Order.Status, "Order status should be 'paid' after payment")
	})

	// Test getting payment by order ID
	t.Run("Get Payment By Order ID", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/payments/orders/"+strconv.FormatUint(uint64(orderID), 10), nil)
		require.NoError(t, err, "Failed to get payment by order ID")

		var respBody struct {
			Payment struct {
				ID            uint    `json:"id"`
				OrderID       uint    `json:"order_id"`
				Amount        float64 `json:"amount"`
				Currency      string  `json:"currency"`
				Method        string  `json:"method"`
				Status        string  `json:"status"`
				TransactionID string  `json:"transaction_id"`
				PaymentDate   string  `json:"payment_date"`
			} `json:"payment"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse get payment response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, paymentID, respBody.Payment.ID)
		assert.Equal(t, orderID, respBody.Payment.OrderID)
		assert.Equal(t, 39.98, respBody.Payment.Amount)
		assert.Equal(t, "USD", respBody.Payment.Currency)
		assert.Equal(t, "credit_card", respBody.Payment.Method)
		assert.Equal(t, "completed", respBody.Payment.Status)
		assert.NotEmpty(t, respBody.Payment.TransactionID)
		assert.NotEmpty(t, respBody.Payment.PaymentDate)
	})

	// Test admin payment operations
	t.Run("Admin Payment Operations", func(t *testing.T) {
		// Login as admin
		token := CreateTestAdmin(t)
		client.SetAuthToken(token)

		// Get all payments
		t.Run("Get All Payments", func(t *testing.T) {
			resp, err := client.DoRequest(http.MethodGet, "/payments", nil)
			require.NoError(t, err, "Failed to get all payments")

			var respBody struct {
				Payments []struct {
					ID            uint    `json:"id"`
					OrderID       uint    `json:"order_id"`
					Amount        float64 `json:"amount"`
					Currency      string  `json:"currency"`
					Method        string  `json:"method"`
					Status        string  `json:"status"`
					TransactionID string  `json:"transaction_id"`
					PaymentDate   string  `json:"payment_date"`
				} `json:"payments"`
				Meta struct {
					Total    int `json:"total"`
					Page     int `json:"page"`
					PageSize int `json:"page_size"`
				} `json:"meta"`
			}

			err = ParseResponse(resp, &respBody)
			require.NoError(t, err, "Failed to parse get all payments response")

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
			assert.GreaterOrEqual(t, respBody.Meta.Total, 1, "Should have at least 1 payment")

			// Find our test payment
			var found bool
			for _, payment := range respBody.Payments {
				if payment.ID == paymentID {
					found = true
					assert.Equal(t, orderID, payment.OrderID)
					assert.Equal(t, 39.98, payment.Amount)
					assert.Equal(t, "USD", payment.Currency)
					assert.Equal(t, "credit_card", payment.Method)
					assert.Equal(t, "completed", payment.Status)
					break
				}
			}
			assert.True(t, found, "Test payment not found in payments list")
		})

		// Get payment by ID
		t.Run("Get Payment By ID", func(t *testing.T) {
			resp, err := client.DoRequest(http.MethodGet, "/payments/"+strconv.FormatUint(uint64(paymentID), 10), nil)
			require.NoError(t, err, "Failed to get payment by ID")

			var respBody struct {
				Payment struct {
					ID            uint    `json:"id"`
					OrderID       uint    `json:"order_id"`
					Amount        float64 `json:"amount"`
					Currency      string  `json:"currency"`
					Method        string  `json:"method"`
					Status        string  `json:"status"`
					TransactionID string  `json:"transaction_id"`
					PaymentDate   string  `json:"payment_date"`
				} `json:"payment"`
			}

			err = ParseResponse(resp, &respBody)
			require.NoError(t, err, "Failed to parse get payment by ID response")

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
			assert.Equal(t, paymentID, respBody.Payment.ID)
			assert.Equal(t, orderID, respBody.Payment.OrderID)
			assert.Equal(t, 39.98, respBody.Payment.Amount)
			assert.Equal(t, "USD", respBody.Payment.Currency)
			assert.Equal(t, "credit_card", respBody.Payment.Method)
			assert.Equal(t, "completed", respBody.Payment.Status)
		})

		// Test refunding a payment
		t.Run("Refund Payment", func(t *testing.T) {
			resp, err := client.DoRequest(http.MethodPost, "/payments/"+strconv.FormatUint(uint64(paymentID), 10)+"/refund", nil)
			require.NoError(t, err, "Failed to refund payment")

			var respBody struct {
				Message string `json:"message"`
				Payment struct {
					ID     uint   `json:"id"`
					Status string `json:"status"`
				} `json:"payment"`
			}

			err = ParseResponse(resp, &respBody)
			require.NoError(t, err, "Failed to parse refund payment response")

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
			assert.Equal(t, "Payment refunded successfully", respBody.Message)
			assert.Equal(t, paymentID, respBody.Payment.ID)
			assert.Equal(t, "refunded", respBody.Payment.Status)

			// Verify order status has been updated
			resp, err = client.DoRequest(http.MethodGet, "/orders/"+strconv.FormatUint(uint64(orderID), 10), nil)
			require.NoError(t, err, "Failed to get order after refund")

			var getOrderResp struct {
				Order struct {
					Status string `json:"status"`
				} `json:"order"`
			}

			err = ParseResponse(resp, &getOrderResp)
			require.NoError(t, err, "Failed to parse get order response after refund")
			assert.Equal(t, "refunded", getOrderResp.Order.Status, "Order status should be 'refunded' after payment refund")
		})
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
