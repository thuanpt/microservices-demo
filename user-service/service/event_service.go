package service

import (
	"context"
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
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient := shared.NewRedisClient(redisHost, redisPort, redisPassword)

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	rabbitmqClient := shared.NewRabbitMQClient(rabbitmqURL)

	setupRabbitMQ(rabbitmqClient)

	return &EventService{
		redisClient:    redisClient,
		rabbitmqClient: rabbitmqClient,
	}
}

func setupRabbitMQ(client *shared.RabbitMQClient) {
	err := client.DeclareExchange("microservices.events", "fanout")
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	err = client.DeclareExchange("microservices.jobs", "direct")
	if err != nil {
		log.Fatalf("Failed to declare jobs exchange: %v", err)
	}

	log.Println("✅ RabbitMQ topology setup completed")
}

func (es *EventService) PublishEvent(ctx context.Context, event shared.Event) error {
	err := es.redisClient.PublishEvent(ctx, "events", event)
	if err != nil {
		log.Printf("❌ Failed to publish to Redis: %v", err)
	}

	err = es.rabbitmqClient.PublishEvent("microservices.events", event.GetType(), event)
	if err != nil {
		log.Printf("❌ Failed to publish to RabbitMQ: %v", err)
		return err
	}

	log.Printf("✅ Event published successfully: %s", event.GetType())
	return nil
}

func (es *EventService) PublishJob(jobType string, jobData interface{}) error {
	jobEvent := shared.NewBaseEvent(jobType, jobData)

	err := es.rabbitmqClient.PublishEvent("microservices.jobs", jobType, jobEvent)
	if err != nil {
		log.Printf("❌ Failed to publish job: %v", err)
		return err
	}

	log.Printf("✅ Job published: %s", jobType)
	return nil
}

func (es *EventService) Close() {
	if es.redisClient != nil {
		es.redisClient.Close()
	}
	if es.rabbitmqClient != nil {
		es.rabbitmqClient.Close()
	}
}
