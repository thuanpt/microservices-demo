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
        // Hash password trước khi lưu
        hashedPassword, err := utils.HashPassword(user.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hash mật khẩu"})
            return
        }
        user.Password = hashedPassword
        id, err := repository.InsertUser(db, &user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo user"})
            return
        }
        user.ID = id
        c.JSON(http.StatusCreated, gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        })
    }
}

// Xử lý đăng nhập user
func LoginUser(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var loginData struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := c.ShouldBindJSON(&loginData); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
            return
        }
        user, err := repository.GetUserByUsername(db, loginData.Username)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
            return
        }
        if !utils.CheckPasswordHash(loginData.Password, user.Password) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
        })
    }
}

// Xử lý lấy thông tin user theo id
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