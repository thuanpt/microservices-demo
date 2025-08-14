package events

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    "github.com/go-redis/redis/v8"
)

type RedisClient struct {
    client *redis.Client
}

func NewRedisClient(host, port, password string) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: password,
        DB:       0,
    })

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    log.Println("✅ Connected to Redis")
    return &RedisClient{client: rdb}
}

// Publish event to Redis pub/sub
func (r *RedisClient) PublishEvent(ctx context.Context, channel string, event Event) error {
    eventJSON, err := ToJSON(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %v", err)
    }

    err = r.client.Publish(ctx, channel, eventJSON).Err()
    if err != nil {
        return fmt.Errorf("failed to publish event: %v", err)
    }

    log.Printf("📤 Published event %s to channel %s", event.GetType(), channel)
    return nil
}

// Subscribe to Redis channel
func (r *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
    pubsub := r.client.Subscribe(ctx, channels...)
    log.Printf("📥 Subscribed to channels: %v", channels)
    return pubsub
}

// Cache operations
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    valueJSON, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, valueJSON, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
    val, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}
