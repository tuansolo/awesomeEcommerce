package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"awesomeEcommerce/internal/config"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer handles producing messages to Kafka
type KafkaProducer struct {
	writers map[string]*kafka.Writer
}

// KafkaConsumer handles consuming messages from Kafka
type KafkaConsumer struct {
	readers map[string]*kafka.Reader
	groupID string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(cfg *config.Config) *KafkaProducer {
	writers := make(map[string]*kafka.Writer)

	// Create writers for each topic
	topics := []string{
		cfg.Kafka.Topics.OrderCreated,
		cfg.Kafka.Topics.OrderUpdated,
		cfg.Kafka.Topics.ProductSync,
		cfg.Kafka.Topics.PaymentStatus,
	}

	for _, topic := range topics {
		writers[topic] = &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 100 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		}
	}

	return &KafkaProducer{
		writers: writers,
	}
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(cfg *config.Config) *KafkaConsumer {
	readers := make(map[string]*kafka.Reader)

	// Create readers for each topic
	topics := []string{
		cfg.Kafka.Topics.OrderCreated,
		cfg.Kafka.Topics.OrderUpdated,
		cfg.Kafka.Topics.ProductSync,
		cfg.Kafka.Topics.PaymentStatus,
	}

	for _, topic := range topics {
		readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.Kafka.Brokers,
			Topic:          topic,
			GroupID:        cfg.Kafka.GroupID,
			MinBytes:       10e3,    // 10KB
			MaxBytes:       10e6,    // 10MB
			MaxWait:        1 * time.Second,
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 1 * time.Second,
		})
	}

	return &KafkaConsumer{
		readers: readers,
		groupID: cfg.Kafka.GroupID,
	}
}

// Publish publishes a message to the specified topic
func (p *KafkaProducer) Publish(ctx context.Context, topic string, key, value []byte) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("no writer found for topic: %s", topic)
	}

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	return nil
}

// Subscribe consumes messages from the specified topic and calls the handler function
func (c *KafkaConsumer) Subscribe(ctx context.Context, topic string, handler func(kafka.Message) error) {
	reader, ok := c.readers[topic]
	if !ok {
		log.Printf("No reader found for topic: %s", topic)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled, stopping consumer for topic: %s", topic)
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message from topic %s: %v", topic, err)
					continue
				}

				if err := handler(msg); err != nil {
					log.Printf("Error handling message from topic %s: %v", topic, err)
				}
			}
		}
	}()
}

// Close closes all Kafka writers
func (p *KafkaProducer) Close() error {
	var lastErr error
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close writer for topic %s: %w", topic, err)
			log.Println(lastErr)
		}
	}
	return lastErr
}

// Close closes all Kafka readers
func (c *KafkaConsumer) Close() error {
	var lastErr error
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close reader for topic %s: %w", topic, err)
			log.Println(lastErr)
		}
	}
	return lastErr
}