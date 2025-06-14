package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafkaMessaging(t *testing.T) {
	// Setup test environment
	client := SetupTest(t)
	defer TeardownTest(t)

	var productID uint
	var categoryID uint
	var orderID uint

	// Create Kafka consumers for the topics we want to test
	orderCreatedConsumer := NewKafkaConsumerTest("order-created")
	defer orderCreatedConsumer.Close()

	orderUpdatedConsumer := NewKafkaConsumerTest("order-updated")
	defer orderUpdatedConsumer.Close()

	paymentStatusConsumer := NewKafkaConsumerTest("payment-status")
	defer paymentStatusConsumer.Close()

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

	// Test order creation and verify Kafka message
	t.Run("Create Order and Verify Kafka Message", func(t *testing.T) {
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

		// Verify order-created Kafka message
		ctx := context.Background()
		msg, err := orderCreatedConsumer.ConsumeMessage(ctx, 5*time.Second)
		require.NoError(t, err, "Failed to consume order-created message")

		// Parse the message
		var orderCreatedMsg struct {
			OrderID uint `json:"order_id"`
			UserID  uint `json:"user_id"`
		}

		err = json.Unmarshal(msg.Value, &orderCreatedMsg)
		require.NoError(t, err, "Failed to parse order-created message")

		// Verify the message content
		assert.Equal(t, orderID, orderCreatedMsg.OrderID, "Order ID in Kafka message should match")
		assert.NotZero(t, orderCreatedMsg.UserID, "User ID in Kafka message should not be zero")
	})

	// Test order cancellation and verify Kafka message
	t.Run("Cancel Order and Verify Kafka Message", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodPost, "/orders/me/"+strconv.FormatUint(uint64(orderID), 10)+"/cancel", nil)
		require.NoError(t, err, "Failed to cancel order")
		require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")

		// Verify order-updated Kafka message
		ctx := context.Background()
		msg, err := orderUpdatedConsumer.ConsumeMessage(ctx, 5*time.Second)
		require.NoError(t, err, "Failed to consume order-updated message")

		// Parse the message
		var orderUpdatedMsg struct {
			OrderID uint   `json:"order_id"`
			Status  string `json:"status"`
		}

		err = json.Unmarshal(msg.Value, &orderUpdatedMsg)
		require.NoError(t, err, "Failed to parse order-updated message")

		// Verify the message content
		assert.Equal(t, orderID, orderUpdatedMsg.OrderID, "Order ID in Kafka message should match")
		assert.Equal(t, "cancelled", orderUpdatedMsg.Status, "Status in Kafka message should be 'cancelled'")
	})

	// Create a new order for payment testing
	t.Run("Create New Order for Payment", func(t *testing.T) {
		// Add item to cart
		addItemReqBody := map[string]interface{}{
			"product_id": productID,
			"quantity":   1,
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

		// Consume the order-created message to clear the queue
		ctx := context.Background()
		_, err = orderCreatedConsumer.ConsumeMessage(ctx, 5*time.Second)
		require.NoError(t, err, "Failed to consume order-created message")
	})

	// Test payment creation and verify Kafka message
	t.Run("Create Payment and Verify Kafka Message", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"amount": 19.99,
			"method": "credit_card",
		}

		resp, err := client.DoRequest(http.MethodPost, "/payments/orders/"+strconv.FormatUint(uint64(orderID), 10), reqBody)
		require.NoError(t, err, "Failed to create payment")
		require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")

		// Verify payment-status Kafka message
		ctx := context.Background()
		msg, err := paymentStatusConsumer.ConsumeMessage(ctx, 5*time.Second)
		require.NoError(t, err, "Failed to consume payment-status message")

		// Parse the message
		var paymentStatusMsg struct {
			OrderID   uint    `json:"order_id"`
			PaymentID uint    `json:"payment_id"`
			Status    string  `json:"status"`
			Amount    float64 `json:"amount"`
			Method    string  `json:"method"`
		}

		err = json.Unmarshal(msg.Value, &paymentStatusMsg)
		require.NoError(t, err, "Failed to parse payment-status message")

		// Verify the message content
		assert.Equal(t, orderID, paymentStatusMsg.OrderID, "Order ID in Kafka message should match")
		assert.NotZero(t, paymentStatusMsg.PaymentID, "Payment ID in Kafka message should not be zero")
		assert.Equal(t, "completed", paymentStatusMsg.Status, "Status in Kafka message should be 'completed'")
		assert.Equal(t, 19.99, paymentStatusMsg.Amount, "Amount in Kafka message should match")
		assert.Equal(t, "credit_card", paymentStatusMsg.Method, "Method in Kafka message should match")

		// Verify order-updated Kafka message (order status should be updated to 'paid')
		msg, err = orderUpdatedConsumer.ConsumeMessage(ctx, 5*time.Second)
		require.NoError(t, err, "Failed to consume order-updated message")

		// Parse the message
		var orderUpdatedMsg struct {
			OrderID uint   `json:"order_id"`
			Status  string `json:"status"`
		}

		err = json.Unmarshal(msg.Value, &orderUpdatedMsg)
		require.NoError(t, err, "Failed to parse order-updated message")

		// Verify the message content
		assert.Equal(t, orderID, orderUpdatedMsg.OrderID, "Order ID in Kafka message should match")
		assert.Equal(t, "paid", orderUpdatedMsg.Status, "Status in Kafka message should be 'paid'")
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
