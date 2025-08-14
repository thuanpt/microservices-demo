package main

import (
	"database/sql"
	"log"
	"order-service/handler"
	"order-service/service"
	"os"

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

	// Test kết nối
	if err := db.Ping(); err != nil {
		log.Fatalf("Không thể ping DB: %v", err)
	}

	// Khởi tạo EventService
	eventSvc := service.NewEventService()
	defer eventSvc.Close()
	eventSvc.ListenUserEvents()
	eventSvc.ListenProductEvents()

	// Khởi tạo router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":          "ok",
			"service":         "order-service",
			"user_service":    os.Getenv("USER_SERVICE_URL"),
			"product_service": os.Getenv("PRODUCT_SERVICE_URL"),
		})
	})

	// API routes cho orders
	r.POST("/orders", handler.CreateOrder(db, eventSvc))         // Tạo đơn hàng mới
	r.GET("/orders/:id", handler.GetOrder(db))                   // Lấy đơn hàng theo ID
	r.GET("/users/:user_id/orders", handler.GetOrdersByUser(db)) // Lấy đơn hàng của user
	r.PUT("/orders/:id/status", handler.UpdateOrderStatus(db))   // Cập nhật status đơn hàng

	// Chạy server
	log.Printf("Order Service đang chạy trên port %s", os.Getenv("APP_PORT"))
	log.Printf("User Service URL: %s", os.Getenv("USER_SERVICE_URL"))
	log.Printf("Product Service URL: %s", os.Getenv("PRODUCT_SERVICE_URL"))

	r.Run(":" + os.Getenv("APP_PORT"))
}
