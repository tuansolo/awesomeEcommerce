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

// OrderWorker handles asynchronous order processing tasks
type OrderWorker struct {
	orderService   service.OrderService
	paymentService service.PaymentService
	consumer       *messaging.KafkaConsumer
	producer       *messaging.KafkaProducer
}

// NewOrderWorker creates a new OrderWorker
func NewOrderWorker(
	orderService service.OrderService,
	paymentService service.PaymentService,
	consumer *messaging.KafkaConsumer,
	producer *messaging.KafkaProducer,
) *OrderWorker {
	return &OrderWorker{
		orderService:   orderService,
		paymentService: paymentService,
		consumer:       consumer,
		producer:       producer,
	}
}

// Start starts the order worker
func (w *OrderWorker) Start(ctx context.Context) {
	// Subscribe to order-related topics
	w.consumer.Subscribe(ctx, "order-created", w.handleOrderCreated)
	w.consumer.Subscribe(ctx, "order-updated", w.handleOrderUpdated)
	w.consumer.Subscribe(ctx, "payment-status", w.handlePaymentStatus)

	log.Println("Order worker started")
}

// handleOrderCreated processes order created events
func (w *OrderWorker) handleOrderCreated(msg kafka.Message) error {
	log.Printf("Received order-created event: %s", string(msg.Value))

	// Parse the order ID from the message
	orderID, err := strconv.ParseUint(string(msg.Value), 10, 32)
	if err != nil {
		log.Printf("Error parsing order ID: %v", err)
		return err
	}

	// Get the order details
	order, err := w.orderService.GetOrderByID(context.Background(), uint(orderID))
	if err != nil {
		log.Printf("Error getting order details: %v", err)
		return err
	}

	// Process the order (e.g., send confirmation email, notify inventory, etc.)
	log.Printf("Processing order %d for user %d with total amount %f", order.ID, order.UserID, order.TotalAmount)

	// In a real application, we would perform additional processing here
	// For this example, we'll just simulate a delay
	time.Sleep(500 * time.Millisecond)

	// Update the order status to processing
	err = w.orderService.UpdateOrderStatus(context.Background(), order.ID, domain.OrderStatusProcessing)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		return err
	}

	// Publish an order-updated event
	orderJSON, _ := json.Marshal(map[string]interface{}{
		"id":     order.ID,
		"status": domain.OrderStatusProcessing,
	})
	err = w.producer.Publish(context.Background(), "order-updated", []byte(strconv.FormatUint(uint64(order.ID), 10)), orderJSON)
	if err != nil {
		log.Printf("Error publishing order-updated event: %v", err)
		return err
	}

	log.Printf("Order %d processed successfully", order.ID)
	return nil
}

// handleOrderUpdated processes order updated events
func (w *OrderWorker) handleOrderUpdated(msg kafka.Message) error {
	log.Printf("Received order-updated event: %s", string(msg.Value))

	// Parse the order update from the message
	var orderUpdate struct {
		ID     uint               `json:"id"`
		Status domain.OrderStatus `json:"status"`
	}
	if err := json.Unmarshal(msg.Value, &orderUpdate); err != nil {
		log.Printf("Error parsing order update: %v", err)
		return err
	}

	// Process the order update based on the new status
	switch orderUpdate.Status {
	case domain.OrderStatusProcessing:
		// Order is being processed, notify the user
		log.Printf("Order %d is now being processed", orderUpdate.ID)
		// In a real application, we would send a notification to the user

	case domain.OrderStatusShipped:
		// Order has been shipped, update tracking information
		log.Printf("Order %d has been shipped", orderUpdate.ID)
		// In a real application, we would update tracking information

	case domain.OrderStatusDelivered:
		// Order has been delivered, send a follow-up email
		log.Printf("Order %d has been delivered", orderUpdate.ID)
		// In a real application, we would send a follow-up email

	case domain.OrderStatusCancelled:
		// Order has been cancelled, process refund if payment was made
		log.Printf("Order %d has been cancelled", orderUpdate.ID)
		// In a real application, we would process a refund if payment was made
	}

	return nil
}

// handlePaymentStatus processes payment status events
func (w *OrderWorker) handlePaymentStatus(msg kafka.Message) error {
	log.Printf("Received payment-status event: %s", string(msg.Value))

	// Parse the payment update from the message
	var paymentUpdate struct {
		ID        uint                 `json:"id"`
		OrderID   uint                 `json:"order_id"`
		Status    domain.PaymentStatus `json:"status"`
		Timestamp time.Time            `json:"timestamp"`
	}
	if err := json.Unmarshal(msg.Value, &paymentUpdate); err != nil {
		log.Printf("Error parsing payment update: %v", err)
		return err
	}

	// Process the payment update based on the new status
	switch paymentUpdate.Status {
	case domain.PaymentStatusCompleted:
		// Payment completed, update order status to processing
		log.Printf("Payment %d for order %d completed", paymentUpdate.ID, paymentUpdate.OrderID)
		err := w.orderService.UpdateOrderStatus(context.Background(), paymentUpdate.OrderID, domain.OrderStatusProcessing)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			return err
		}

	case domain.PaymentStatusFailed:
		// Payment failed, update order status to cancelled
		log.Printf("Payment %d for order %d failed", paymentUpdate.ID, paymentUpdate.OrderID)
		err := w.orderService.UpdateOrderStatus(context.Background(), paymentUpdate.OrderID, domain.OrderStatusCancelled)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			return err
		}

	case domain.PaymentStatusRefunded:
		// Payment refunded, update order status to cancelled
		log.Printf("Payment %d for order %d refunded", paymentUpdate.ID, paymentUpdate.OrderID)
		err := w.orderService.UpdateOrderStatus(context.Background(), paymentUpdate.OrderID, domain.OrderStatusCancelled)
		if err != nil {
			log.Printf("Error updating order status: %v", err)
			return err
		}
	}

	return nil
}