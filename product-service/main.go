package main

import (
	"database/sql"
	"log"
	"os"
	"product-service/handler"
	"product-service/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load biến môi trường
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

	// Test kết nối DB
	if err := db.Ping(); err != nil {
		log.Fatalf("Không thể ping DB: %v", err)
	}

	// Khởi tạo EventService
	eventSvc := service.NewEventService()
	defer eventSvc.Close()
	eventSvc.ListenUserEvents()

	// Khởi tạo router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "product-service"})
	})

	// API routes cho products
	r.POST("/products", handler.CreateProduct(db, eventSvc))               // Tạo product
	r.GET("/products", handler.GetAllProducts(db))                         // Lấy tất cả products
	r.GET("/products/:id", handler.GetProduct(db))                         // Lấy product theo ID
	r.PUT("/products/:id/stock", handler.UpdateProductStock(db, eventSvc)) // Cập nhật stock

	// Chạy server trên port 8002
	log.Printf("Product Service đang chạy trên port %s", os.Getenv("APP_PORT"))
	r.Run(":" + os.Getenv("APP_PORT"))
}
