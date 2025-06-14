package api

import (
	"awesomeEcommerce/internal/middleware"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
)

// Router sets up all the routes for the API
type Router struct {
	userHandler    *UserHandler
	productHandler *ProductHandler
	cartHandler    *CartHandler
	orderHandler   *OrderHandler
	paymentHandler *PaymentHandler
}

// NewRouter creates a new Router
func NewRouter(
	userService service.UserService,
	productService service.ProductService,
	cartService service.CartService,
	orderService service.OrderService,
	paymentService service.PaymentService,
) *Router {
	return &Router{
		userHandler:    NewUserHandler(userService),
		productHandler: NewProductHandler(productService, userService),
		cartHandler:    NewCartHandler(cartService, userService, productService),
		orderHandler:   NewOrderHandler(orderService, userService),
		paymentHandler: NewPaymentHandler(paymentService, orderService, userService),
	}
}

// SetupRoutes sets up all the routes for the API
func (r *Router) SetupRoutes(engine *gin.Engine) {
	// Apply global middleware
	engine.Use(middleware.LoggingMiddleware())
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.SecureMiddleware())

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Register all handlers
		r.userHandler.RegisterRoutes(v1)
		r.productHandler.RegisterRoutes(v1)
		r.cartHandler.RegisterRoutes(v1)
		r.orderHandler.RegisterRoutes(v1)
		r.paymentHandler.RegisterRoutes(v1)
	}

	// No route found handler
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Endpoint not found"})
	})
}

// SetupSwagger sets up the Swagger documentation
// This is a placeholder for future Swagger integration
func (r *Router) SetupSwagger(engine *gin.Engine) {
	// TODO: Implement Swagger documentation
}
