package main

import (
    "database/sql"
    "log"
    "os"
    "user-service/handler"
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

    // Khởi tạo router Gin
    r := gin.Default()

    // Route kiểm tra service
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Route đăng ký user
    r.POST("/users", handler.RegisterUser(db))

    // Route đăng nhập user
    r.POST("/login", handler.LoginUser(db))

    // Route lấy thông tin user
    r.GET("/users/:id", handler.GetUser(db))

    // Chạy server
    r.Run(":" + os.Getenv("APP_PORT"))
}