package api

import (
	"net/http"
	"strconv"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/middleware"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	productService service.ProductService
	userService    service.UserService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productService service.ProductService, userService service.UserService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		userService:    userService,
	}
}

// RegisterRoutes registers the routes for the ProductHandler
func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		// Public routes
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProductByID)
		products.GET("/categories", h.GetCategories)
		products.GET("/categories/:id", h.GetProductsByCategory)

		// Admin routes
		admin := products.Use(middleware.AuthMiddleware(h.userService), middleware.RoleMiddleware("admin"))
		{
			admin.POST("", h.CreateProduct)
			admin.PUT("/:id", h.UpdateProduct)
			admin.DELETE("/:id", h.DeleteProduct)
			admin.POST("/categories", h.CreateCategory)
			admin.PUT("/categories/:id", h.UpdateCategory)
			admin.DELETE("/categories/:id", h.DeleteCategory)
		}
	}
}

// GetProducts returns all products with pagination
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, total, err := h.productService.GetProducts(c, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	var productList []gin.H
	for _, product := range products {
		productList = append(productList, gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"sku":         product.SKU,
			"image_url":   product.ImageURL,
			"category_id": product.CategoryID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products": productList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProductByID returns a product by ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"sku":         product.SKU,
			"image_url":   product.ImageURL,
			"category_id": product.CategoryID,
		},
	})
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"required,gte=0"`
		SKU         string  `json:"sku" binding:"required"`
		ImageURL    string  `json:"image_url"`
		CategoryID  uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		SKU:         request.SKU,
		ImageURL:    request.ImageURL,
		CategoryID:  request.CategoryID,
	}

	if err := h.productService.CreateProduct(c, product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"sku":         product.SKU,
			"image_url":   product.ImageURL,
			"category_id": product.CategoryID,
		},
	})
}

// UpdateProduct updates an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"omitempty,gt=0"`
		Stock       int     `json:"stock" binding:"omitempty,gte=0"`
		SKU         string  `json:"sku"`
		ImageURL    string  `json:"image_url"`
		CategoryID  uint    `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current product
	product, err := h.productService.GetProductByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Update fields if provided
	if request.Name != "" {
		product.Name = request.Name
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price > 0 {
		product.Price = request.Price
	}
	if request.Stock >= 0 {
		product.Stock = request.Stock
	}
	if request.SKU != "" {
		product.SKU = request.SKU
	}
	if request.ImageURL != "" {
		product.ImageURL = request.ImageURL
	}
	if request.CategoryID > 0 {
		product.CategoryID = request.CategoryID
	}

	if err := h.productService.UpdateProduct(c, product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"sku":         product.SKU,
			"image_url":   product.ImageURL,
			"category_id": product.CategoryID,
		},
	})
}

// DeleteProduct deletes a product by ID
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetCategories returns all product categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get categories"})
		return
	}

	var categoryList []gin.H
	for _, category := range categories {
		categoryList = append(categoryList, gin.H{
			"id":        category.ID,
			"name":      category.Name,
			"parent_id": category.ParentID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"categories": categoryList})
}

// GetProductsByCategory returns products by category ID with pagination
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, total, err := h.productService.GetProductsByCategory(c, uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	var productList []gin.H
	for _, product := range products {
		productList = append(productList, gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"sku":         product.SKU,
			"image_url":   product.ImageURL,
			"category_id": product.CategoryID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products": productList,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateCategory creates a new product category
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var request struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &domain.ProductCategory{
		Name:     request.Name,
		ParentID: request.ParentID,
	}

	if err := h.productService.CreateCategory(c, category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"category": gin.H{
			"id":        category.ID,
			"name":      category.Name,
			"parent_id": category.ParentID,
		},
	})
}

// UpdateCategory updates an existing product category
func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var request struct {
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current category
	category, err := h.productService.GetCategoryByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Update fields if provided
	if request.Name != "" {
		category.Name = request.Name
	}
	if request.ParentID != nil {
		category.ParentID = request.ParentID
	}

	if err := h.productService.UpdateCategory(c, category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"category": gin.H{
			"id":        category.ID,
			"name":      category.Name,
			"parent_id": category.ParentID,
		},
	})
}

// DeleteCategory deletes a product category by ID
func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.productService.DeleteCategory(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
