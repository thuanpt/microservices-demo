package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	shared "microservices.local/shared"
)

// EventService quản lý việc publish events
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

// PublishProductCreatedEvent - Publish khi tạo sản phẩm mới
func (es *EventService) PublishProductCreatedEvent(productID int, name string, price float64, stock int) error {
	event := shared.NewBaseEvent("product.created", map[string]interface{}{
		"product_id":   productID,
		"product_name": name,
		"price":        price,
		"stock":        stock,
	})

	// Publish qua Redis
	if err := es.redisClient.PublishEvent(context.Background(), "product.events", event); err != nil {
		log.Printf("Lỗi publish event qua Redis: %v", err)
	}

	// Publish qua RabbitMQ
	if err := es.rabbitmqClient.PublishEvent("product.exchange", "product.created", event); err != nil {
		log.Printf("Lỗi publish event qua RabbitMQ: %v", err)
		return err
	}

	log.Printf("✅ Published ProductCreated event: %d - %s", productID, name)
	return nil
}

// PublishProductUpdatedEvent - Publish khi cập nhật sản phẩm
func (es *EventService) PublishProductUpdatedEvent(productID int, name string, price float64, stock int) error {
	event := shared.NewBaseEvent("product.updated", map[string]interface{}{
		"product_id":   productID,
		"product_name": name,
		"price":        price,
		"stock":        stock,
	})

	// Publish qua Redis
	if err := es.redisClient.PublishEvent(context.Background(), "product.events", event); err != nil {
		log.Printf("Lỗi publish event qua Redis: %v", err)
	}

	// Publish qua RabbitMQ
	if err := es.rabbitmqClient.PublishEvent("product.exchange", "product.updated", event); err != nil {
		log.Printf("Lỗi publish event qua RabbitMQ: %v", err)
		return err
	}

	log.Printf("✅ Published ProductUpdated event: %d - %s", productID, name)
	return nil
}

// ListenUserEvents - Lắng nghe events từ User Service
func (es *EventService) ListenUserEvents() {
	go func() {
		log.Println("🎧 Product Service đang lắng nghe User events...")

		// Listen qua RabbitMQ
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
				// Logic xử lý khi có user mới đăng ký
				// Ví dụ: gửi email welcome với danh sách sản phẩm hot

			case "user.login":
				log.Printf("📨 Received user.login: %+v", event.Data)
				// Logic xử lý khi user login
				// Ví dụ: track user activity, recommend products
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
