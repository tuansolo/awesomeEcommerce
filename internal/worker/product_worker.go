package worker

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/messaging"
	"awesomeEcommerce/internal/service"

	"github.com/segmentio/kafka-go"
)

// ProductWorker handles asynchronous product processing tasks
type ProductWorker struct {
	productService service.ProductService
	consumer       *messaging.KafkaConsumer
	producer       *messaging.KafkaProducer
}

// NewProductWorker creates a new ProductWorker
func NewProductWorker(
	productService service.ProductService,
	consumer *messaging.KafkaConsumer,
	producer *messaging.KafkaProducer,
) *ProductWorker {
	return &ProductWorker{
		productService: productService,
		consumer:       consumer,
		producer:       producer,
	}
}

// Start starts the product worker
func (w *ProductWorker) Start(ctx context.Context) {
	// Subscribe to product-related topics
	w.consumer.Subscribe(ctx, "product-sync", w.handleProductSync)
	w.consumer.Subscribe(ctx, "product-stock-update", w.handleProductStockUpdate)

	log.Println("Product worker started")
}

// handleProductSync processes product sync events
func (w *ProductWorker) handleProductSync(msg kafka.Message) error {
	log.Printf("Received product-sync event: %s", string(msg.Value))

	// Parse the product sync request from the message
	var syncRequest struct {
		Action string         `json:"action"` // "create", "update", "delete"
		Product domain.Product `json:"product"`
	}
	if err := json.Unmarshal(msg.Value, &syncRequest); err != nil {
		log.Printf("Error parsing product sync request: %v", err)
		return err
	}

	// Process the product sync request based on the action
	switch syncRequest.Action {
	case "create":
		// Create a new product
		log.Printf("Creating product: %s", syncRequest.Product.Name)
		err := w.productService.CreateProduct(context.Background(), &syncRequest.Product)
		if err != nil {
			log.Printf("Error creating product: %v", err)
			return err
		}
		log.Printf("Product created successfully: %d", syncRequest.Product.ID)

	case "update":
		// Update an existing product
		log.Printf("Updating product: %d", syncRequest.Product.ID)
		err := w.productService.UpdateProduct(context.Background(), &syncRequest.Product)
		if err != nil {
			log.Printf("Error updating product: %v", err)
			return err
		}
		log.Printf("Product updated successfully: %d", syncRequest.Product.ID)

	case "delete":
		// Delete a product
		log.Printf("Deleting product: %d", syncRequest.Product.ID)
		err := w.productService.DeleteProduct(context.Background(), syncRequest.Product.ID)
		if err != nil {
			log.Printf("Error deleting product: %v", err)
			return err
		}
		log.Printf("Product deleted successfully: %d", syncRequest.Product.ID)

	default:
		log.Printf("Unknown action: %s", syncRequest.Action)
		return nil
	}

	// Publish a product-sync-completed event
	syncCompletedJSON, _ := json.Marshal(map[string]interface{}{
		"action":      syncRequest.Action,
		"product_id":  syncRequest.Product.ID,
		"product_sku": syncRequest.Product.SKU,
		"timestamp":   time.Now(),
	})
	err := w.producer.Publish(context.Background(), "product-sync-completed", []byte(strconv.FormatUint(uint64(syncRequest.Product.ID), 10)), syncCompletedJSON)
	if err != nil {
		log.Printf("Error publishing product-sync-completed event: %v", err)
		return err
	}

	return nil
}

// handleProductStockUpdate processes product stock update events
func (w *ProductWorker) handleProductStockUpdate(msg kafka.Message) error {
	log.Printf("Received product-stock-update event: %s", string(msg.Value))

	// Parse the stock update request from the message
	var stockUpdate struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"` // Positive for increase, negative for decrease
	}
	if err := json.Unmarshal(msg.Value, &stockUpdate); err != nil {
		log.Printf("Error parsing stock update request: %v", err)
		return err
	}

	// Get the current product
	product, err := w.productService.GetProductByID(context.Background(), stockUpdate.ProductID)
	if err != nil {
		log.Printf("Error getting product: %v", err)
		return err
	}

	// Update the stock
	log.Printf("Updating stock for product %d (%s): %+d", product.ID, product.Name, stockUpdate.Quantity)
	err = w.productService.UpdateProductStock(context.Background(), product.ID, stockUpdate.Quantity)
	if err != nil {
		log.Printf("Error updating product stock: %v", err)
		return err
	}

	// Get the updated product
	updatedProduct, err := w.productService.GetProductByID(context.Background(), stockUpdate.ProductID)
	if err != nil {
		log.Printf("Error getting updated product: %v", err)
		return err
	}

	log.Printf("Stock updated successfully for product %d. New stock: %d", product.ID, updatedProduct.Stock)

	// Check if stock is low and needs replenishment
	if updatedProduct.Stock < 10 {
		log.Printf("Low stock alert for product %d (%s): %d units remaining", updatedProduct.ID, updatedProduct.Name, updatedProduct.Stock)
		
		// Publish a low-stock-alert event
		alertJSON, _ := json.Marshal(map[string]interface{}{
			"product_id":   updatedProduct.ID,
			"product_name": updatedProduct.Name,
			"product_sku":  updatedProduct.SKU,
			"current_stock": updatedProduct.Stock,
			"timestamp":    time.Now(),
		})
		err := w.producer.Publish(context.Background(), "low-stock-alert", []byte(strconv.FormatUint(uint64(updatedProduct.ID), 10)), alertJSON)
		if err != nil {
			log.Printf("Error publishing low-stock-alert event: %v", err)
			return err
		}
	}

	return nil
}

// SyncProductsFromExternalSource simulates syncing products from an external source
func (w *ProductWorker) SyncProductsFromExternalSource(ctx context.Context) error {
	log.Println("Starting product sync from external source")

	// In a real application, this would fetch products from an external API or database
	// For this example, we'll simulate it with a few sample products
	sampleProducts := []domain.Product{
		{
			Name:        "Sample Product 1",
			Description: "This is a sample product for testing",
			Price:       19.99,
			Stock:       100,
			SKU:         "SAMPLE-001",
			ImageURL:    "https://example.com/sample1.jpg",
			CategoryID:  1,
		},
		{
			Name:        "Sample Product 2",
			Description: "Another sample product for testing",
			Price:       29.99,
			Stock:       50,
			SKU:         "SAMPLE-002",
			ImageURL:    "https://example.com/sample2.jpg",
			CategoryID:  1,
		},
		{
			Name:        "Sample Product 3",
			Description: "Yet another sample product for testing",
			Price:       39.99,
			Stock:       25,
			SKU:         "SAMPLE-003",
			ImageURL:    "https://example.com/sample3.jpg",
			CategoryID:  2,
		},
	}

	// Process each sample product
	for _, product := range sampleProducts {
		// Check if the product already exists by SKU
		existingProduct, err := w.productService.GetProductBySKU(ctx, product.SKU)
		if err == nil {
			// Product exists, update it
			existingProduct.Name = product.Name
			existingProduct.Description = product.Description
			existingProduct.Price = product.Price
			existingProduct.Stock = product.Stock
			existingProduct.ImageURL = product.ImageURL
			existingProduct.CategoryID = product.CategoryID

			// Publish an update event
			syncRequest := map[string]interface{}{
				"action":  "update",
				"product": existingProduct,
			}
			syncRequestJSON, _ := json.Marshal(syncRequest)
			err = w.producer.Publish(ctx, "product-sync", []byte(existingProduct.SKU), syncRequestJSON)
			if err != nil {
				log.Printf("Error publishing product-sync event for update: %v", err)
				continue
			}
		} else {
			// Product doesn't exist, create it
			// Publish a create event
			syncRequest := map[string]interface{}{
				"action":  "create",
				"product": product,
			}
			syncRequestJSON, _ := json.Marshal(syncRequest)
			err = w.producer.Publish(ctx, "product-sync", []byte(product.SKU), syncRequestJSON)
			if err != nil {
				log.Printf("Error publishing product-sync event for create: %v", err)
				continue
			}
		}
	}

	log.Println("Product sync from external source completed")
	return nil
}