package config

import (
    "os"
)

// Cấu hình các service endpoints
type ServiceConfig struct {
    UserServiceURL    string
    ProductServiceURL string
    OrderServiceURL   string
    GatewayPort       string
    JWTSecret         string
}

// Load config từ environment variables
func LoadConfig() *ServiceConfig {
    return &ServiceConfig{
        UserServiceURL:    getEnv("USER_SERVICE_URL", "http://localhost:8001"),
        ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "http://localhost:8002"),
        OrderServiceURL:   getEnv("ORDER_SERVICE_URL", "http://localhost:8003"),
        GatewayPort:       getEnv("GATEWAY_PORT", "8000"),
        JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
    }
}

// Helper function để lấy env với default value
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// Định nghĩa route mapping
type RouteConfig struct {
    Prefix      string
    ServiceURL  string
    RequireAuth bool
}

// Cấu hình routes
func GetRoutes(config *ServiceConfig) []RouteConfig {
    return []RouteConfig{
        // User Service routes
        {
            Prefix:      "/api/users",
            ServiceURL:  config.UserServiceURL,
            RequireAuth: false, // Login/register không cần auth
        },
        {
            Prefix:      "/api/profile",
            ServiceURL:  config.UserServiceURL,
            RequireAuth: true, // Profile cần auth
        },
        
        // Product Service routes
        {
            Prefix:      "/api/products",
            ServiceURL:  config.ProductServiceURL,
            RequireAuth: false, // Xem sản phẩm không cần auth
        },
        
        // Order Service routes
        {
            Prefix:      "/api/orders",
            ServiceURL:  config.OrderServiceURL,
            RequireAuth: true, // Tạo/xem đơn hàng cần auth
        },
    }
}