package service

import (
    "bytes"                    // Thêm import này
    "encoding/json"
    "fmt"
    "net/http"
    "order-service/model"
    "os"
    "time"
)

// HTTP Client với timeout để gọi các service khác
var httpClient = &http.Client{
    Timeout: 10 * time.Second, // Timeout 10 giây
}

// Kiểm tra user có tồn tại không bằng cách gọi User Service
func CheckUserExists(userID int) (*model.UserResponse, error) {
    url := fmt.Sprintf("%s/users/%d", os.Getenv("USER_SERVICE_URL"), userID)
    
    // Gửi GET request tới User Service
    resp, err := httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("không thể kết nối User Service: %v", err)
    }
    defer resp.Body.Close()

    // Kiểm tra status code
    if resp.StatusCode == 404 {
        return nil, fmt.Errorf("user không tồn tại")
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("User Service trả về lỗi: %d", resp.StatusCode)
    }

    // Parse JSON response
    var user model.UserResponse
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, fmt.Errorf("không thể parse response từ User Service: %v", err)
    }

    return &user, nil
}

// Lấy thông tin sản phẩm từ Product Service
func GetProduct(productID int) (*model.ProductResponse, error) {
    url := fmt.Sprintf("%s/products/%d", os.Getenv("PRODUCT_SERVICE_URL"), productID)
    
    // Gửi GET request tới Product Service
    resp, err := httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("không thể kết nối Product Service: %v", err)
    }
    defer resp.Body.Close()

    // Kiểm tra status code
    if resp.StatusCode == 404 {
        return nil, fmt.Errorf("sản phẩm không tồn tại")
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Product Service trả về lỗi: %d", resp.StatusCode)
    }

    // Parse JSON response
    var product model.ProductResponse
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, fmt.Errorf("không thể parse response từ Product Service: %v", err)
    }

    return &product, nil
}

// Cập nhật stock sản phẩm sau khi đặt hàng (gọi Product Service)
func UpdateProductStock(productID int, newStock int) error {
    url := fmt.Sprintf("%s/products/%d/stock", os.Getenv("PRODUCT_SERVICE_URL"), productID)
    
    // Tạo request body JSON
    requestBody := map[string]int{"stock": newStock}
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return fmt.Errorf("không thể tạo JSON request: %v", err)
    }

    // Tạo PUT request
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("không thể tạo request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Gửi request
    resp, err := httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("không thể gửi request tới Product Service: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("Product Service trả về lỗi khi cập nhật stock: %d", resp.StatusCode)
    }

    return nil
}