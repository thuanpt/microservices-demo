package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	shared "microservices.local/shared"
)

// EventService quản lý việc publish và consume events cho Order Service
type EventService struct {
	redisClient    *shared.RedisClient
	rabbitmqClient *shared.RabbitMQClient
}

func NewEventService() *EventService {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}
	redisClient := shared.NewRedisClient(redisHost, redisPort, "")
	rabbitmqClient := shared.NewRabbitMQClient(rabbitmqURL)
	return &EventService{
		redisClient:    redisClient,
		rabbitmqClient: rabbitmqClient,
	}
}

// PublishOrderCreatedEvent - Publish khi tạo đơn hàng mới
func (es *EventService) PublishOrderCreatedEvent(orderID int, userID int, productID int, quantity int, total float64) error {
	event := shared.NewBaseEvent("order.created", map[string]interface{}{
		"order_id":   orderID,
		"user_id":    userID,
		"product_id": productID,
		"quantity":   quantity,
		"total":      total,
	})
	// Publish qua Redis
	if err := es.redisClient.PublishEvent(context.Background(), "order.events", event); err != nil {
		log.Printf("Lỗi publish event qua Redis: %v", err)
	}
	// Publish qua RabbitMQ
	if err := es.rabbitmqClient.PublishEvent("order.exchange", "order.created", event); err != nil {
		log.Printf("Lỗi publish event qua RabbitMQ: %v", err)
		return err
	}
	log.Printf("✅ Published OrderCreated event: %d", orderID)
	return nil
}

// ListenUserEvents - Lắng nghe events từ User Service
func (es *EventService) ListenUserEvents() {
	go func() {
		log.Println("🎧 Order Service đang lắng nghe User events...")
		msgs, err := es.rabbitmqClient.ConsumeMessages("user.events.queue")
		if err != nil {
			log.Printf("❌ Lỗi khi consume User events: %v", err)
			return
		}
		for msg := range msgs {
			var event shared.BaseEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Lỗi unmarshal event: %v", err)
				continue
			}
			switch event.Type {
			case "user.registered":
				log.Printf("📨 Received user.registered: %+v", event.Data)
			case "user.login":
				log.Printf("📨 Received user.login: %+v", event.Data)
			}
			msg.Ack(false)
		}
	}()
}

// ListenProductEvents - Lắng nghe events từ Product Service
func (es *EventService) ListenProductEvents() {
	go func() {
		log.Println("🎧 Order Service đang lắng nghe Product events...")
		msgs, err := es.rabbitmqClient.ConsumeMessages("product.events.queue")
		if err != nil {
			log.Printf("❌ Lỗi khi consume Product events: %v", err)
			return
		}
		for msg := range msgs {
			var event shared.BaseEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("❌ Lỗi unmarshal event: %v", err)
				continue
			}
			switch event.Type {
			case "product.created":
				log.Printf("📨 Received product.created: %+v", event.Data)
			case "product.updated":
				log.Printf("📨 Received product.updated: %+v", event.Data)
			}
			msg.Ack(false)
		}
	}()
}

func (es *EventService) Close() {
	if es.redisClient != nil {
		es.redisClient.Close()
	}
	if es.rabbitmqClient != nil {
		es.rabbitmqClient.Close()
	}
}
