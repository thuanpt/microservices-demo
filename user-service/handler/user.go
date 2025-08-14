package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
	"user-service/model"
	"user-service/repository"
	"user-service/service"
	"user-service/utils"
	"github.com/gin-gonic/gin"
	shared "microservices.local/shared" // ✅ Local module
)

type UserHandler struct {
	db           *sql.DB
	eventService *service.EventService
}

func NewUserHandler(db *sql.DB, eventService *service.EventService) *UserHandler {
	return &UserHandler{
		db:           db,
		eventService: eventService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hash mật khẩu"})
		return
	}
	user.Password = hashedPassword

	id, err := repository.InsertUser(h.db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo user"})
		return
	}
	user.ID = id

	// 🎉 PUBLISH EVENT
	userRegisteredEvent := shared.NewUserRegisteredEvent(user.ID, user.Username, user.Email)

	ctx := context.Background()
	err = h.eventService.PublishEvent(ctx, userRegisteredEvent)
	if err != nil {
		log.Printf("⚠️ Failed to publish user registered event: %v", err)
	}

	// 📧 PUBLISH JOB
	welcomeEmailJob := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"template": "welcome",
	}

	err = h.eventService.PublishJob("email.send_welcome", welcomeEmailJob)
	if err != nil {
		log.Printf("⚠️ Failed to publish welcome email job: %v", err)
	}

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
		"token": token,
	})
}

// Các function khác giữ nguyên...
func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username và password là bắt buộc"})
		return
	}

	user, err := repository.GetUserByUsername(h.db, loginData.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
		return
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai tên đăng nhập hoặc mật khẩu"})
		return
	}

	// 📊 PUBLISH LOGIN EVENT
	userLoginEvent := shared.NewBaseEvent("user.login", map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"login_time": time.Now().UTC(),
		"ip":         c.ClientIP(),
	})

	ctx := context.Background()
	h.eventService.PublishEvent(ctx, userLoginEvent)

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
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy thông tin user"})
		return
	}

	user, err := repository.GetUserByID(h.db, userID.(int))
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

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}

	user, err := repository.GetUserByID(h.db, id)
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
