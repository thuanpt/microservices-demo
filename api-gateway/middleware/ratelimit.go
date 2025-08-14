package middleware

import (
    "net/http"
    "sync"
    "time"
    "github.com/gin-gonic/gin"
)

// Rate limiter struct
type RateLimiter struct {
    mu      sync.RWMutex
    clients map[string]*ClientInfo
    limit   int           // Request limit per window
    window  time.Duration // Time window
}

// Client info để track request count
type ClientInfo struct {
    count     int
    resetTime time.Time
}

// Tạo rate limiter mới
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        clients: make(map[string]*ClientInfo),
        limit:   limit,
        window:  window,
    }

    // Cleanup goroutine để xóa expired clients
    go rl.cleanupExpiredClients()

    return rl
}

// Middleware rate limiting
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()
        
        rl.mu.Lock()
        defer rl.mu.Unlock()

        now := time.Now()
        client, exists := rl.clients[clientIP]

        if !exists {
            // Client mới
            rl.clients[clientIP] = &ClientInfo{
                count:     1,
                resetTime: now.Add(rl.window),
            }
            c.Next()
            return
        }

        // Reset counter nếu window đã hết hạn
        if now.After(client.resetTime) {
            client.count = 1
            client.resetTime = now.Add(rl.window)
            c.Next()
            return
        }

        // Kiểm tra limit
        if client.count >= rl.limit {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error":      "Rate limit exceeded",
                "limit":      rl.limit,
                "window":     rl.window.String(),
                "reset_time": client.resetTime.Unix(),
                "service":    "api-gateway",
            })
            c.Abort()
            return
        }

        // Increment counter
        client.count++
        c.Next()
    }
}

// Cleanup expired clients mỗi 5 phút
func (rl *RateLimiter) cleanupExpiredClients() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            rl.mu.Lock()
            now := time.Now()
            for clientIP, client := range rl.clients {
                if now.After(client.resetTime.Add(rl.window)) {
                    delete(rl.clients, clientIP)
                }
            }
            rl.mu.Unlock()
        }
    }
}