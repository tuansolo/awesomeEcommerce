package api

import (
	"net/http"
	"strconv"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/middleware"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests related to orders
type OrderHandler struct {
	orderService service.OrderService
	userService  service.UserService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService service.OrderService, userService service.UserService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		userService:  userService,
	}
}

// RegisterRoutes registers the routes for the OrderHandler
func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		// Customer routes (require authentication)
		auth := orders.Use(middleware.AuthMiddleware(h.userService))
		{
			auth.POST("", h.CreateOrder)
			auth.GET("/me", h.GetMyOrders)
			auth.GET("/me/:id", h.GetMyOrderByID)
			auth.POST("/me/:id/cancel", h.CancelOrder)
		}

		// Admin routes
		admin := orders.Use(middleware.AuthMiddleware(h.userService), middleware.RoleMiddleware("admin"))
		{
			admin.GET("", h.GetAllOrders)
			admin.GET("/:id", h.GetOrderByID)
			admin.PUT("/:id/status", h.UpdateOrderStatus)
			admin.GET("/status/:status", h.GetOrdersByStatus)
			admin.GET("/date-range", h.GetOrdersByDateRange)
		}
	}
}

// CreateOrder creates a new order from the user's cart
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
		BillingAddress  string `json:"billing_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the order
	order, err := h.orderService.CreateOrder(c, userID.(uint), request.ShippingAddress, request.BillingAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Format the response
	var orderItems []gin.H
	for _, item := range order.Items {
		orderItems = append(orderItems, gin.H{
			"id":           item.ID,
			"product_id":   item.ProductID,
			"product_name": item.ProductName,
			"price":        item.Price,
			"quantity":     item.Quantity,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order": gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
			"items":            orderItems,
		},
	})
}

// GetMyOrders returns all orders for the authenticated user
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.GetOrdersByUserID(c, userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var orderList []gin.H
	for _, order := range orders {
		var orderItems []gin.H
		for _, item := range order.Items {
			orderItems = append(orderItems, gin.H{
				"id":           item.ID,
				"product_id":   item.ProductID,
				"product_name": item.ProductName,
				"price":        item.Price,
				"quantity":     item.Quantity,
			})
		}

		orderList = append(orderList, gin.H{
			"id":               order.ID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
			"items":            orderItems,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetMyOrderByID returns a specific order for the authenticated user
func (h *OrderHandler) GetMyOrderByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get the order
	order, err := h.orderService.GetOrderByID(c, uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if the order belongs to the user
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Format the response
	var orderItems []gin.H
	for _, item := range order.Items {
		orderItems = append(orderItems, gin.H{
			"id":           item.ID,
			"product_id":   item.ProductID,
			"product_name": item.ProductName,
			"price":        item.Price,
			"quantity":     item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
			"items":            orderItems,
		},
	})
}

// CancelOrder cancels an order for the authenticated user
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get the order
	order, err := h.orderService.GetOrderByID(c, uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if the order belongs to the user
	if order.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Cancel the order
	err = h.orderService.CancelOrder(c, uint(orderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// GetAllOrders returns all orders (admin only)
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.GetAllOrders(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var orderList []gin.H
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetOrderByID returns a specific order (admin only)
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Get the order
	order, err := h.orderService.GetOrderByID(c, uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Format the response
	var orderItems []gin.H
	for _, item := range order.Items {
		orderItems = append(orderItems, gin.H{
			"id":           item.ID,
			"product_id":   item.ProductID,
			"product_name": item.ProductName,
			"price":        item.Price,
			"quantity":     item.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order": gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
			"items":            orderItems,
		},
	})
}

// UpdateOrderStatus updates the status of an order (admin only)
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var request struct {
		Status domain.OrderStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the order status
	err = h.orderService.UpdateOrderStatus(c, uint(orderID), request.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// GetOrdersByStatus returns orders by status (admin only)
func (h *OrderHandler) GetOrdersByStatus(c *gin.Context) {
	status := domain.OrderStatus(c.Param("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate the status
	if status != domain.OrderStatusPending &&
		status != domain.OrderStatusProcessing &&
		status != domain.OrderStatusShipped &&
		status != domain.OrderStatusDelivered &&
		status != domain.OrderStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
		return
	}

	orders, total, err := h.orderService.GetOrdersByStatus(c, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var orderList []gin.H
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetOrdersByDateRange returns orders created within a date range (admin only)
func (h *OrderHandler) GetOrdersByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate the dates
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date and end date are required"})
		return
	}

	// Validate date format
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format, use YYYY-MM-DD"})
		return
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format, use YYYY-MM-DD"})
		return
	}

	orders, total, err := h.orderService.GetOrdersByDateRange(c, startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	var orderList []gin.H
	for _, order := range orders {
		orderList = append(orderList, gin.H{
			"id":               order.ID,
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
			"shipping_address": order.ShippingAddress,
			"billing_address":  order.BillingAddress,
			"created_at":       order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orderList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
