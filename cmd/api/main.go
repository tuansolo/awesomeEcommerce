package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeEcommerce/internal/api"
	"awesomeEcommerce/internal/config"
	"awesomeEcommerce/internal/messaging"
	"awesomeEcommerce/internal/repository"
	"awesomeEcommerce/internal/repository/cache"
	"awesomeEcommerce/internal/repository/db"
	"awesomeEcommerce/internal/repository/impl"
	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	app := fx.New(
		// Provide all the constructors
		fx.Provide(
			// Config
			config.LoadConfig,

			// Database
			func(cfg *config.Config) (*gorm.DB, error) {
				return db.NewDatabase(cfg)
			},

			// Redis cache
			func(cfg *config.Config) (*cache.RedisClient, error) {
				redisClient, err := cache.NewRedisClient(cfg)
				if err != nil {
					return nil, err
				}
				return redisClient, nil
			},

			// Kafka
			func(cfg *config.Config) *messaging.KafkaProducer {
				return messaging.NewKafkaProducer(cfg)
			},

			// Repositories
			func(database *gorm.DB, redisClient *cache.RedisClient) repository.UserRepository {
				return impl.NewUserRepository(database, redisClient)
			},
			func(database *gorm.DB, redisClient *cache.RedisClient) repository.ProductRepository {
				return impl.NewProductRepository(database, redisClient)
			},
			func(database *gorm.DB, redisClient *cache.RedisClient) repository.CartRepository {
				return impl.NewCartRepository(database, redisClient)
			},
			func(database *gorm.DB, redisClient *cache.RedisClient) repository.OrderRepository {
				return impl.NewOrderRepository(database, redisClient)
			},
			func(database *gorm.DB, redisClient *cache.RedisClient) repository.PaymentRepository {
				return impl.NewPaymentRepository(database, redisClient)
			},

			// Services
			func(repo repository.UserRepository) service.UserService {
				return service.NewUserService(repo, nil, nil)
			},
			func(repo repository.ProductRepository) service.ProductService {
				return service.NewProductService(repo)
			},
			func(cartRepo repository.CartRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository) service.CartService {
				return service.NewCartService(cartRepo, productRepo, userRepo)
			},
			func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, producer *messaging.KafkaProducer) service.OrderService {
				return service.NewOrderService(orderRepo, cartRepo, productRepo, nil, producer)
			},
			func(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, producer *messaging.KafkaProducer) service.PaymentService {
				return service.NewPaymentService(paymentRepo, orderRepo, producer)
			},

			// API Router
			func(userService service.UserService, productService service.ProductService, cartService service.CartService, orderService service.OrderService, paymentService service.PaymentService) *api.Router {
				return api.NewRouter(userService, productService, cartService, orderService, paymentService)
			},

			// Gin Engine
			func() *gin.Engine {
				return gin.Default()
			},
		),

		// Register lifecycle hooks
		fx.Invoke(
			// Initialize database
			func(database *gorm.DB) {
				if err := db.AutoMigrate(database); err != nil {
					log.Fatalf("Failed to migrate database: %v", err)
				}
			},

			// Set up API routes and Swagger
			func(router *api.Router, engine *gin.Engine) {
				router.SetupRoutes(engine)
				router.SetupSwagger(engine)
			},

			// Start the HTTP server
			func(lc fx.Lifecycle, engine *gin.Engine, cfg *config.Config) {
				server := &http.Server{
					Addr:         ":" + cfg.Server.Port,
					Handler:      engine,
					ReadTimeout:  cfg.Server.ReadTimeout,
					WriteTimeout: cfg.Server.WriteTimeout,
					IdleTimeout:  cfg.Server.IdleTimeout,
				}

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							log.Printf("Starting API server on port %s", cfg.Server.Port)
							if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
								log.Fatalf("Failed to start server: %v", err)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Println("Shutting down API server...")
						ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
						defer cancel()
						return server.Shutdown(ctx)
					},
				})
			},

			// Set up graceful shutdown
			func(lc fx.Lifecycle, producer *messaging.KafkaProducer, redisClient *cache.RedisClient) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						// Set up signal handling
						go func() {
							sigCh := make(chan os.Signal, 1)
							signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
							sig := <-sigCh
							log.Printf("Received signal: %v", sig)
							os.Exit(0)
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						// Close Kafka producer
						if err := producer.Close(); err != nil {
							log.Printf("Error closing Kafka producer: %v", err)
						}

						// Close Redis client
						if err := redisClient.Close(); err != nil {
							log.Printf("Error closing Redis client: %v", err)
						}

						return nil
					},
				})
			},
		),
	)

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Block until the application is stopped
	<-app.Done()
}
