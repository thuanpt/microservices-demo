package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/gin-gonic/gin"  // Thêm import gin
)

// Struct chứa thông tin user trong JWT token
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

// Secret key để sign JWT (trong production nên lưu trong env)
var jwtSecret = []byte("your-secret-key-change-this-in-production")

// Tạo JWT token cho user
func GenerateJWT(userID int, username, email string) (string, error) {
    // Tạo claims với thời hạn 24 giờ
    claims := Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token hết hạn sau 24h
            IssuedAt:  jwt.NewNumericDate(time.Now()),                      // Thời gian tạo token
            Issuer:    "user-service",                                      // Service tạo token
        },
    }

    // Tạo token với signing method HS256
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // Sign token với secret key
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Verify JWT token và trả về claims
func VerifyJWT(tokenString string) (*Claims, error) {
    // Parse token
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Kiểm tra signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    // Kiểm tra token có valid không
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// Middleware để xác thực JWT token
func JWTMiddleware() gin.HandlerFunc {  // Sửa return type
    return func(c *gin.Context) {
        // Lấy token từ Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // Format: "Bearer <token>"
        if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
            c.JSON(401, gin.H{"error": "Invalid authorization format. Use: Bearer <token>"})
            c.Abort()
            return
        }

        tokenString := authHeader[7:] // Bỏ "Bearer " prefix

        // Verify token
        claims, err := VerifyJWT(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Lưu thông tin user vào context để sử dụng trong handler
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)

        c.Next() // Tiếp tục đến handler tiếp theo
    }
}