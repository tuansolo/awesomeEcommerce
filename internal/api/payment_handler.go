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

// PaymentHandler handles HTTP requests related to payments
type PaymentHandler struct {
	paymentService service.PaymentService
	orderService   service.OrderService
	userService    service.UserService
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService service.PaymentService, orderService service.OrderService, userService service.UserService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		orderService:   orderService,
		userService:    userService,
	}
}

// RegisterRoutes registers the routes for the PaymentHandler
func (h *PaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	payments := router.Group("/payments")
	{
		// Customer routes (require authentication)
		auth := payments.Use(middleware.AuthMiddleware(h.userService))
		{
			auth.POST("/orders/:id", h.CreatePayment)
			auth.GET("/orders/:id", h.GetPaymentByOrderID)
		}

		// Admin routes
		admin := payments.Use(middleware.AuthMiddleware(h.userService), middleware.RoleMiddleware("admin"))
		{
			admin.GET("", h.GetAllPayments)
			admin.GET("/:id", h.GetPaymentByID)
			admin.PUT("/:id/status", h.UpdatePaymentStatus)
			admin.POST("/:id/refund", h.RefundPayment)
			admin.GET("/status/:status", h.GetPaymentsByStatus)
			admin.GET("/method/:method", h.GetPaymentsByMethod)
			admin.GET("/date-range", h.GetPaymentsByDateRange)
		}
	}
}

// CreatePayment creates a new payment for an order
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
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

	var request struct {
		Amount float64              `json:"amount" binding:"required,gt=0"`
		Method domain.PaymentMethod `json:"method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the order to check if it belongs to the user
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

	// Create the payment
	payment, err := h.paymentService.CreatePayment(c, uint(orderID), request.Amount, request.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real application, we would integrate with a payment gateway here
	// For this example, we'll simulate a successful payment
	transactionID := "txn_" + strconv.FormatInt(time.Now().Unix(), 10)
	err = h.paymentService.ProcessPayment(c, payment.ID, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated payment
	payment, err = h.paymentService.GetPaymentByID(c, payment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment processed successfully",
		"payment": gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		},
	})
}

// GetPaymentByOrderID returns the payment for a specific order
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
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

	// Get the order to check if it belongs to the user
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

	// Get the payment
	payment, err := h.paymentService.GetPaymentByOrderID(c, uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		},
	})
}

// GetAllPayments returns all payments (admin only)
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	payments, total, err := h.paymentService.GetAllPayments(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}

	var paymentList []gin.H
	for _, payment := range payments {
		paymentList = append(paymentList, gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": paymentList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPaymentByID returns a specific payment (admin only)
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.paymentService.GetPaymentByID(c, uint(paymentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		},
	})
}

// UpdatePaymentStatus updates the status of a payment (admin only)
func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var request struct {
		Status domain.PaymentStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the status
	if request.Status != domain.PaymentStatusPending &&
		request.Status != domain.PaymentStatusCompleted &&
		request.Status != domain.PaymentStatusFailed &&
		request.Status != domain.PaymentStatusRefunded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
		return
	}

	err = h.paymentService.UpdatePaymentStatus(c, uint(paymentID), request.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}

// RefundPayment refunds a payment (admin only)
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	err = h.paymentService.RefundPayment(c, uint(paymentID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment refunded successfully"})
}

// GetPaymentsByStatus returns payments by status (admin only)
func (h *PaymentHandler) GetPaymentsByStatus(c *gin.Context) {
	status := domain.PaymentStatus(c.Param("status"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate the status
	if status != domain.PaymentStatusPending &&
		status != domain.PaymentStatusCompleted &&
		status != domain.PaymentStatusFailed &&
		status != domain.PaymentStatusRefunded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
		return
	}

	payments, total, err := h.paymentService.GetPaymentsByStatus(c, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}

	var paymentList []gin.H
	for _, payment := range payments {
		paymentList = append(paymentList, gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": paymentList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPaymentsByMethod returns payments by method (admin only)
func (h *PaymentHandler) GetPaymentsByMethod(c *gin.Context) {
	method := domain.PaymentMethod(c.Param("method"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate the method
	if method != domain.PaymentMethodCreditCard &&
		method != domain.PaymentMethodDebitCard &&
		method != domain.PaymentMethodPayPal &&
		method != domain.PaymentMethodBankTransfer {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	payments, total, err := h.paymentService.GetPaymentsByMethod(c, method, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}

	var paymentList []gin.H
	for _, payment := range payments {
		paymentList = append(paymentList, gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": paymentList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetPaymentsByDateRange returns payments created within a date range (admin only)
func (h *PaymentHandler) GetPaymentsByDateRange(c *gin.Context) {
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

	payments, total, err := h.paymentService.GetPaymentsByDateRange(c, startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payments"})
		return
	}

	var paymentList []gin.H
	for _, payment := range payments {
		paymentList = append(paymentList, gin.H{
			"id":             payment.ID,
			"order_id":       payment.OrderID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"method":         payment.Method,
			"status":         payment.Status,
			"transaction_id": payment.TransactionID,
			"payment_date":   payment.PaymentDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": paymentList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
