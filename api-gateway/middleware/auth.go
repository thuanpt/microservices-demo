package middleware

import (
    "errors"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// Claims struct (giống như User Service)
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

// JWT Middleware cho Gateway
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Lấy token từ Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "Authorization header is required",
                "service": "api-gateway",
            })
            c.Abort()
            return
        }

        // Kiểm tra format "Bearer <token>"
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "Invalid authorization format. Use: Bearer <token>",
                "service": "api-gateway",
            })
            c.Abort()
            return
        }

        tokenString := authHeader[7:] // Bỏ "Bearer " prefix

        // Verify token
        claims, err := verifyJWT(tokenString, jwtSecret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "Invalid or expired token: " + err.Error(),
                "service": "api-gateway",
            })
            c.Abort()
            return
        }

        // Set user info vào header để forward tới backend services
        c.Header("X-User-ID", string(rune(claims.UserID)))
        c.Header("X-Username", claims.Username)
        c.Header("X-Email", claims.Email)

        // Lưu thông tin user vào context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("email", claims.Email)

        c.Next()
    }
}

// Verify JWT token (copy từ User Service)
func verifyJWT(tokenString string, jwtSecret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}