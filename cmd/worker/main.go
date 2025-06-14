package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeEcommerce/internal/config"
	"awesomeEcommerce/internal/messaging"
	"awesomeEcommerce/internal/repository"
	"awesomeEcommerce/internal/repository/cache"
	"awesomeEcommerce/internal/repository/db"
	"awesomeEcommerce/internal/repository/impl"
	"awesomeEcommerce/internal/service"
	"awesomeEcommerce/internal/worker"

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
			func(cfg *config.Config) *messaging.KafkaConsumer {
				return messaging.NewKafkaConsumer(cfg)
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
			func(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository, userRepo repository.UserRepository, producer *messaging.KafkaProducer) service.OrderService {
				return service.NewOrderService(orderRepo, cartRepo, productRepo, userRepo, producer)
			},
			func(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, producer *messaging.KafkaProducer) service.PaymentService {
				return service.NewPaymentService(paymentRepo, orderRepo, producer)
			},

			// Workers
			func(orderService service.OrderService, paymentService service.PaymentService, consumer *messaging.KafkaConsumer, producer *messaging.KafkaProducer) *worker.OrderWorker {
				return worker.NewOrderWorker(orderService, paymentService, consumer, producer)
			},
			func(productService service.ProductService, consumer *messaging.KafkaConsumer, producer *messaging.KafkaProducer) *worker.ProductWorker {
				return worker.NewProductWorker(productService, consumer, producer)
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

			// Start the workers
			func(lc fx.Lifecycle, orderWorker *worker.OrderWorker, productWorker *worker.ProductWorker, cfg *config.Config) {
				workerCtx, cancel := context.WithCancel(context.Background())

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						log.Println("Starting workers...")

						// Start the order worker
						orderWorker.Start(workerCtx)

						// Start the product worker
						productWorker.Start(workerCtx)

						// Run initial product sync (optional)
						go func() {
							time.Sleep(5 * time.Second) // Wait for everything to initialize
							if err := productWorker.SyncProductsFromExternalSource(workerCtx); err != nil {
								log.Printf("Error syncing products: %v", err)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Println("Stopping workers...")
						cancel() // Cancel the worker context
						return nil
					},
				})
			},

			// Set up graceful shutdown
			func(lc fx.Lifecycle, consumer *messaging.KafkaConsumer, producer *messaging.KafkaProducer, redisClient *cache.RedisClient) {
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
						// Close Kafka consumer
						if err := consumer.Close(); err != nil {
							log.Printf("Error closing Kafka consumer: %v", err)
						}

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
