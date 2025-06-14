package api

import (
	"net/http"
	"strconv"

	"awesomeEcommerce/internal/middleware"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
)

// CartHandler handles HTTP requests related to shopping carts
type CartHandler struct {
	cartService    service.CartService
	userService    service.UserService
	productService service.ProductService
}

// NewCartHandler creates a new CartHandler
func NewCartHandler(cartService service.CartService, userService service.UserService, productService service.ProductService) *CartHandler {
	return &CartHandler{
		cartService:    cartService,
		userService:    userService,
		productService: productService,
	}
}

// RegisterRoutes registers the routes for the CartHandler
func (h *CartHandler) RegisterRoutes(router *gin.RouterGroup) {
	carts := router.Group("/carts")
	{
		// All cart routes require authentication
		auth := carts.Use(middleware.AuthMiddleware(h.userService))
		{
			auth.GET("/me", h.GetCart)
			auth.POST("/items", h.AddItemToCart)
			auth.PUT("/items/:id", h.UpdateCartItem)
			auth.DELETE("/items/:id", h.RemoveItemFromCart)
			auth.DELETE("/items", h.ClearCart)
			auth.GET("/total", h.GetCartTotal)
		}
	}
}

// GetCart returns the cart of the authenticated user
func (h *CartHandler) GetCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cart, err := h.cartService.GetCartByUserID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	var cartItems []gin.H
	for _, item := range cart.Items {
		cartItems = append(cartItems, gin.H{
			"id":         item.ID,
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
			"product": gin.H{
				"id":          item.Product.ID,
				"name":        item.Product.Name,
				"description": item.Product.Description,
				"price":       item.Product.Price,
				"image_url":   item.Product.ImageURL,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"cart": gin.H{
			"id":      cart.ID,
			"user_id": cart.UserID,
			"items":   cartItems,
		},
	})
}

// AddItemToCart adds an item to the cart of the authenticated user
func (h *CartHandler) AddItemToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	product, err := h.productService.GetProductByID(c, request.ProductID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if product has enough stock
	if product.Stock < request.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	// Add item to cart
	err = h.cartService.AddItemToCart(c, userID.(uint), request.ProductID, request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// UpdateCartItem updates an item in the cart of the authenticated user
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var request struct {
		Quantity int `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user's cart
	cart, err := h.cartService.GetCartByUserID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	// Update the cart item
	err = h.cartService.UpdateCartItem(c, cart.ID, uint(itemID), request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

// RemoveItemFromCart removes an item from the cart of the authenticated user
func (h *CartHandler) RemoveItemFromCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Get the user's cart
	cart, err := h.cartService.GetCartByUserID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	// Remove the item from the cart
	err = h.cartService.RemoveItemFromCart(c, cart.ID, uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

// ClearCart removes all items from the cart of the authenticated user
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the user's cart
	cart, err := h.cartService.GetCartByUserID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	// Clear the cart
	err = h.cartService.ClearCart(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// GetCartTotal returns the total price of all items in the cart of the authenticated user
func (h *CartHandler) GetCartTotal(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the user's cart
	cart, err := h.cartService.GetCartByUserID(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	// Get the cart total
	total, err := h.cartService.GetCartTotal(c, cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
