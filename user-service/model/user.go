package model

import "time"

// Struct User với đầy đủ fields
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username" binding:"required"`
    Password  string    `json:"password" binding:"required"`
    Email     string    `json:"email" binding:"required,email"`
    CreatedAt time.Time `json:"created_at"`  // Thêm field này
    UpdatedAt time.Time `json:"updated_at"`  // Thêm luôn UpdatedAt cho đầy đủ
}

// Struct cho request đăng ký (không cần ID, CreatedAt, UpdatedAt)
type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=6"` // Tối thiểu 6 ký tự
    Email    string `json:"email" binding:"required,email"`
}

// Struct cho request đăng nhập
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}