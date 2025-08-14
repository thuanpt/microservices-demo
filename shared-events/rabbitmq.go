package events

import (
    "fmt"
    "log"
    "time"
    "github.com/streadway/amqp"
)

type RabbitMQClient struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

func NewRabbitMQClient(url string) *RabbitMQClient {
    conn, err := amqp.Dial(url)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open channel: %v", err)
    }

    log.Println("✅ Connected to RabbitMQ")
    return &RabbitMQClient{
        connection: conn,
        channel:    ch,
    }
}

// Declare exchange
func (r *RabbitMQClient) DeclareExchange(name, exchangeType string) error {
    return r.channel.ExchangeDeclare(
        name,         // name
        exchangeType, // type (direct, fanout, topic, headers)
        true,         // durable
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
}

// Declare queue
func (r *RabbitMQClient) DeclareQueue(name string) (amqp.Queue, error) {
    return r.channel.QueueDeclare(
        name,  // name
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
}

// Bind queue to exchange
func (r *RabbitMQClient) BindQueue(queueName, routingKey, exchangeName string) error {
    return r.channel.QueueBind(
        queueName,    // queue name
        routingKey,   // routing key
        exchangeName, // exchange
        false,
        nil,
    )
}

// Publish event to RabbitMQ
func (r *RabbitMQClient) PublishEvent(exchangeName, routingKey string, event Event) error {
    eventJSON, err := ToJSON(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %v", err)
    }

    err = r.channel.Publish(
        exchangeName, // exchange
        routingKey,   // routing key
        false,        // mandatory
        false,        // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         eventJSON,
            DeliveryMode: amqp.Persistent, // Make message persistent
            Timestamp:    time.Now(),
        },
    )

    if err != nil {
        return fmt.Errorf("failed to publish event: %v", err)
    }

    log.Printf("📤 Published event %s to exchange %s with key %s", event.GetType(), exchangeName, routingKey)
    return nil
}

// Consume messages
func (r *RabbitMQClient) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
    return r.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack (false để manual ack)
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
}

func (r *RabbitMQClient) Close() error {
    if r.channel != nil {
        r.channel.Close()
    }
    if r.connection != nil {
        return r.connection.Close()
    }
    return nil
}
