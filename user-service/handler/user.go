package handler

import (
    "database/sql"
    "net/http"
    "strconv"
    "user-service/model"
    "user-service/repository"
    "user-service/utils"
    "github.com/gin-gonic/gin"
)

// Xử lý đăng ký user
func RegisterUser(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user model.User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
            return
        }

        // Hash password
        hashedPassword, err := utils.HashPassword(user.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hash mật khẩu"})
            return
        }
        user.Password = hashedPassword

        // Lưu user vào DB
        id, err := repository.InsertUser(db, &user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo user"})
            return
        }
        user.ID = id

        // Tạo JWT token cho user mới đăng ký
        token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "message": "Đăng ký thành công",
            "user": gin.H{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
            },
            "token": token, // Trả token để client có thể sử dụng ngay
        })
    }
}

// Xử lý đăng nhập user - QUAN TRỌNG: Trả JWT token
func LoginUser(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var loginData struct {
            Username string `json:"username" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        
        if err := c.ShouldBindJSON(&loginData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Username và password là bắt buộc"})
            return
        }

        // Tìm user trong DB
        user, err := repository.GetUserByUsername(db, loginData.Username)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
            return
        }

        // Kiểm tra password
        if !utils.CheckPasswordHash(loginData.Password, user.Password) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
            return
        }

        // Tạo JWT token
        token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "Đăng nhập thành công",
            "user": gin.H{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
            },
            "token": token, // Trả JWT token
        })
    }
}

// API lấy profile user hiện tại (cần authentication)
func GetProfile(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Lấy user_id từ JWT middleware đã set
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin user"})
            return
        }

        // Get user từ DB
        user, err := repository.GetUserByID(db, userID.(int))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User không tồn tại"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        })
    }
}

// API lấy user theo ID (public - không cần auth)
func GetUser(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
            return
        }

        user, err := repository.GetUserByID(db, id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy user"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        })
    }
}