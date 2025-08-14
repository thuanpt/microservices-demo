package model

import "time"

// Struct Product định nghĩa cấu trúc dữ liệu sản phẩm
type Product struct {
    ID          int       `json:"id"`                    // ID tự động tăng
    Name        string    `json:"name"`                  // Tên sản phẩm
    Description string    `json:"description"`           // Mô tả sản phẩm
    Price       float64   `json:"price"`                // Giá sản phẩm (float64 = số thực)
    Stock       int       `json:"stock"`                // Số lượng tồn kho
    CreatedAt   time.Time `json:"created_at"`           // Thời gian tạo (time.Time là kiểu thời gian của Go)
}

// Struct để nhận dữ liệu tạo/cập nhật product (không cần ID, CreatedAt)
type CreateProductRequest struct {
    Name        string  `json:"name" binding:"required"`        // binding:"required" = bắt buộc phải có
    Description string  `json:"description"`                    // Không bắt buộc
    Price       float64 `json:"price" binding:"required,gt=0"`  // gt=0 nghĩa là phải > 0
    Stock       int     `json:"stock" binding:"min=0"`          // min=0 nghĩa là >= 0
}