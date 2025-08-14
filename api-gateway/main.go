package main

import (
    "log"
    "time"
    "api-gateway/config"
    "api-gateway/middleware"
    "api-gateway/proxy"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: No .env file found, using default values")
    }

    // Load configuration
    cfg := config.LoadConfig()
    routes := config.GetRoutes(cfg)

    // Khởi tạo Gin router
    r := gin.Default()

    // Tạo rate limiter: 100 requests per minute per IP
    rateLimiter := middleware.NewRateLimiter(100, time.Minute)

    // Global middlewares
    r.Use(rateLimiter.RateLimit())
    
    // CORS middleware (cho web apps)
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Gateway health check
    services := map[string]string{
        "user-service":    cfg.UserServiceURL,
        "product-service": cfg.ProductServiceURL,
        "order-service":   cfg.OrderServiceURL,
    }
    r.GET("/health", proxy.HealthCheckProxy(services))

    // Gateway info endpoint
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "service":   "api-gateway",
            "version":   "1.0.0",
            "timestamp": time.Now().UTC().Format("2006-01-02 15:04:05"),
            "routes": []gin.H{
                {"prefix": "/api/users", "service": "user-service", "auth": false},
                {"prefix": "/api/profile", "service": "user-service", "auth": true},
                {"prefix": "/api/products", "service": "product-service", "auth": false},
                {"prefix": "/api/orders", "service": "order-service", "auth": true},
            },
        })
    })

    // Setup route groups
    for _, route := range routes {
        setupRouteGroup(r, route, cfg.JWTSecret)
    }

    // Start gateway
    log.Printf("🚀 API Gateway starting on port %s", cfg.GatewayPort)
    log.Printf("📊 Rate limit: 100 requests/minute per IP")
    log.Printf("🔗 Backend services:")
    log.Printf("   - User Service: %s", cfg.UserServiceURL)
    log.Printf("   - Product Service: %s", cfg.ProductServiceURL) 
    log.Printf("   - Order Service: %s", cfg.OrderServiceURL)
    
    if err := r.Run(":" + cfg.GatewayPort); err != nil {
        log.Fatalf("Failed to start gateway: %v", err)
    }
}

// Setup route group với authentication nếu cần
func setupRouteGroup(r *gin.Engine, route config.RouteConfig, jwtSecret string) {
    group := r.Group(route.Prefix)
    
    // Apply JWT middleware nếu route cần authentication
    if route.RequireAuth {
        group.Use(middleware.JWTMiddleware(jwtSecret))
    }
    
    // Proxy tất cả HTTP methods
    group.Any("/*path", proxy.ProxyRequest(route.ServiceURL, route.Prefix))
    
    log.Printf("✅ Route registered: %s -> %s (auth: %v)", route.Prefix, route.ServiceURL, route.RequireAuth)
}
