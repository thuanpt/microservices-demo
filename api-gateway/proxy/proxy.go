package proxy

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    // Bỏ "net/url" vì không dùng
    "strings"
    "time"
    "github.com/gin-gonic/gin"
)

// HTTP Client với timeout cho proxy requests
var httpClient = &http.Client{
    Timeout: 30 * time.Second, // Timeout 30 giây
}

// Proxy request tới backend service
func ProxyRequest(targetServiceURL string, stripPrefix string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Tạo target URL
        targetURL := buildTargetURL(targetServiceURL, c.Request.URL.Path, stripPrefix, c.Request.URL.RawQuery)
        
        // Tạo request mới
        req, err := createProxyRequest(c, targetURL)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to create proxy request: " + err.Error(),
                "service": "api-gateway",
            })
            return
        }

        // Gửi request tới backend service
        resp, err := httpClient.Do(req)
        if err != nil {
            c.JSON(http.StatusBadGateway, gin.H{
                "error":   "Backend service unavailable: " + err.Error(),
                "service": "api-gateway",
                "target":  targetServiceURL,
            })
            return
        }
        defer resp.Body.Close()

        // Forward response từ backend về client
        forwardResponse(c, resp)
    }
}

// Tạo target URL cho backend service
func buildTargetURL(serviceURL, requestPath, stripPrefix, rawQuery string) string {
    // Bỏ prefix khỏi path
    targetPath := strings.TrimPrefix(requestPath, stripPrefix)
    if !strings.HasPrefix(targetPath, "/") {
        targetPath = "/" + targetPath
    }

    // Tạo full URL
    targetURL := serviceURL + targetPath
    if rawQuery != "" {
        targetURL += "?" + rawQuery
    }

    return targetURL
}

// Tạo proxy request từ original request
func createProxyRequest(c *gin.Context, targetURL string) (*http.Request, error) {
    // Đọc request body
    var bodyBytes []byte
    if c.Request.Body != nil {
        bodyBytes, _ = io.ReadAll(c.Request.Body)
        c.Request.Body.Close()
    }

    // Tạo request mới
    req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(bodyBytes))
    if err != nil {
        return nil, err
    }

    // Copy headers từ original request
    copyHeaders(c.Request.Header, req.Header)

    // Thêm headers từ authentication middleware (nếu có)
    if userID, exists := c.Get("user_id"); exists {
        req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
    }
    if username, exists := c.Get("username"); exists {
        req.Header.Set("X-Username", username.(string))
    }
    if email, exists := c.Get("email"); exists {
        req.Header.Set("X-Email", email.(string))
    }

    // Set client IP và other metadata
    req.Header.Set("X-Forwarded-For", c.ClientIP())
    req.Header.Set("X-Forwarded-Proto", getScheme(c))
    req.Header.Set("X-Forwarded-Host", c.Request.Host)
    req.Header.Set("X-Gateway", "api-gateway")

    return req, nil
}

// Copy headers từ source sang destination
func copyHeaders(src, dest http.Header) {
    for key, values := range src {
        // Bỏ qua một số headers không cần forward
        if shouldSkipHeader(key) {
            continue
        }
        
        for _, value := range values {
            dest.Add(key, value)
        }
    }
}

// Headers không nên forward
func shouldSkipHeader(key string) bool {
    skipHeaders := []string{
        "Connection",
        "Keep-Alive", 
        "Proxy-Authenticate",
        "Proxy-Authorization",
        "Te",
        "Trailers",
        "Transfer-Encoding",
        "Upgrade",
    }

    key = strings.ToLower(key)
    for _, skip := range skipHeaders {
        if strings.ToLower(skip) == key {
            return true
        }
    }
    return false
}

// Forward response từ backend về client
func forwardResponse(c *gin.Context, resp *http.Response) {
    // Copy response headers
    for key, values := range resp.Header {
        for _, value := range values {
            c.Header(key, value)
        }
    }

    // Set status code
    c.Status(resp.StatusCode)

    // Stream response body
    _, err := io.Copy(c.Writer, resp.Body)
    if err != nil {
        // Log error nhưng không thể thay đổi response nữa
        fmt.Printf("Error streaming response: %v\n", err)
    }
}

// Get scheme (http/https)
func getScheme(c *gin.Context) string {
    if c.Request.TLS != nil {
        return "https"
    }
    
    // Check X-Forwarded-Proto header (từ load balancer)
    if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
        return proto
    }
    
    return "http"
}

// Health check proxy - kiểm tra backend services
func HealthCheckProxy(services map[string]string) gin.HandlerFunc {
    return func(c *gin.Context) {
        results := make(map[string]interface{})
        
        // Kiểm tra từng service
        for serviceName, serviceURL := range services {
            healthURL := serviceURL + "/health"
            
            resp, err := httpClient.Get(healthURL)
            if err != nil {
                results[serviceName] = gin.H{
                    "status": "down",
                    "error":  err.Error(),
                }
                continue
            }
            resp.Body.Close()
            
            if resp.StatusCode == 200 {
                results[serviceName] = gin.H{
                    "status": "up",
                    "url":    serviceURL,
                }
            } else {
                results[serviceName] = gin.H{
                    "status": "unhealthy",
                    "code":   resp.StatusCode,
                }
            }
        }
        
        // Tính overall status
        allHealthy := true
        for _, result := range results {
            if result.(gin.H)["status"] != "up" {
                allHealthy = false
                break
            }
        }
        
        status := "healthy"
        httpStatus := http.StatusOK
        if !allHealthy {
            status = "unhealthy"
            httpStatus = http.StatusServiceUnavailable
        }
        
        c.JSON(httpStatus, gin.H{
            "gateway": gin.H{
                "status":    status,
                "timestamp": time.Now().UTC().Format("2006-01-02 15:04:05"),
                "version":   "1.0.0",
            },
            "services": results,
        })
    }
}