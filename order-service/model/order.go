package model

import "time"

// Struct Order cho đơn hàng chính
type Order struct {
    ID          int         `json:"id"`
    UserID      int         `json:"user_id"`              // ID user từ User Service
    TotalAmount float64     `json:"total_amount"`         // Tổng tiền
    Status      string      `json:"status"`               // pending, confirmed, cancelled, completed
    Items       []OrderItem `json:"items,omitempty"`      // Danh sách sản phẩm trong đơn hàng
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// Struct OrderItem cho chi tiết sản phẩm trong đơn hàng
type OrderItem struct {
    ID          int     `json:"id"`
    OrderID     int     `json:"order_id"`
    ProductID   int     `json:"product_id"`              // ID product từ Product Service
    ProductName string  `json:"product_name"`            // Tên sản phẩm lưu lại
    Price       float64 `json:"price"`                   // Giá tại thời điểm đặt hàng
    Quantity    int     `json:"quantity"`                // Số lượng
    Subtotal    float64 `json:"subtotal"`                // Thành tiền
    CreatedAt   time.Time `json:"created_at"`
}

// Struct request tạo đơn hàng mới
type CreateOrderRequest struct {
    UserID int                    `json:"user_id" binding:"required"`      // Bắt buộc có user_id
    Items  []CreateOrderItemRequest `json:"items" binding:"required,dive"`   // dive = validate từng item trong slice
}

// Struct chi tiết item khi tạo đơn hàng
type CreateOrderItemRequest struct {
    ProductID int `json:"product_id" binding:"required"`     // ID sản phẩm
    Quantity  int `json:"quantity" binding:"required,min=1"` // Số lượng >= 1
}

// Struct response từ Product Service (để gọi API lấy thông tin sản phẩm)
type ProductResponse struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

// Struct response từ User Service (để gọi API kiểm tra user)
type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}