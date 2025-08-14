package main

import (
    "database/sql"
    "log"
    "os"
    "user-service/handler"
    "user-service/utils"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // Load biến môi trường từ file .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Lỗi khi đọc file .env: %v", err)
    }

    // Tạo DSN kết nối MySQL
    dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
        "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" +
        os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

    // Kết nối DB
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Lỗi khi kết nối DB: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Không thể ping DB: %v", err)
    }

    // Khởi tạo router
    r := gin.Default()

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
    })

    // PUBLIC routes (không cần authentication)
    r.POST("/register", handler.RegisterUser(db))     // Đăng ký
    r.POST("/login", handler.LoginUser(db))          // Đăng nhập
    r.GET("/users/:id", handler.GetUser(db))         // Lấy user theo ID (public)

    // PROTECTED routes (cần authentication)
    protected := r.Group("/")
    protected.Use(utils.JWTMiddleware()) // Apply JWT middleware
    {
        protected.GET("/profile", handler.GetProfile(db))     // Lấy profile của chính mình
        // Có thể thêm các protected routes khác ở đây
    }

    log.Printf("User Service đang chạy trên port %s", os.Getenv("APP_PORT"))
    r.Run(":" + os.Getenv("APP_PORT"))
}
