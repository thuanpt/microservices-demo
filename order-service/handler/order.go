package handler

import (
    "database/sql"
    "fmt"                                    // Thêm import này
    "net/http"
    "strconv"
    "order-service/model"
    "order-service/repository"
    "order-service/service"
    "github.com/gin-gonic/gin"
)

// Tạo đơn hàng mới
func CreateOrder(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req model.CreateOrderRequest
        
        // Bind và validate JSON request
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Bước 1: Kiểm tra user có tồn tại không
        user, err := service.CheckUserExists(req.UserID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Bước 2: Validate từng sản phẩm và tính tổng tiền
        var orderItems []model.OrderItem
        var totalAmount float64 = 0

        for _, item := range req.Items {
            // Lấy thông tin sản phẩm từ Product Service
            product, err := service.GetProduct(item.ProductID)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": fmt.Sprintf("Sản phẩm ID %d: %v", item.ProductID, err),
                })
                return
            }

            // Kiểm tra stock có đủ không
            if product.Stock < item.Quantity {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": fmt.Sprintf("Sản phẩm '%s' chỉ còn %d, không đủ %d", 
                             product.Name, product.Stock, item.Quantity),
                })
                return
            }

            // Tính subtotal cho item này
            subtotal := product.Price * float64(item.Quantity)
            totalAmount += subtotal

            // Tạo OrderItem
            orderItem := model.OrderItem{
                ProductID:   item.ProductID,
                ProductName: product.Name,
                Price:       product.Price,
                Quantity:    item.Quantity,
                Subtotal:    subtotal,
            }
            orderItems = append(orderItems, orderItem)
        }

        // Bước 3: Tạo Order object
        order := model.Order{
            UserID:      req.UserID,
            TotalAmount: totalAmount,
            Status:      "pending", // Trạng thái ban đầu
            Items:       orderItems,
        }

        // Bước 4: Lưu order vào database
        orderID, err := repository.CreateOrder(db, &order)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đơn hàng"})
            return
        }

        // Bước 5: Cập nhật stock các sản phẩm (gọi Product Service)
        for _, item := range req.Items {  // Bỏ biến i không dùng
            product, _ := service.GetProduct(item.ProductID) // Đã validate ở trên
            newStock := product.Stock - item.Quantity
            
            err := service.UpdateProductStock(item.ProductID, newStock)
            if err != nil {
                // Log lỗi nhưng không fail toàn bộ order (có thể xử lý async sau)
                // Trong production nên dùng message queue để đảm bảo consistency
                fmt.Printf("Warning: Không thể cập nhật stock cho product %d: %v\n", item.ProductID, err)
            }
        }

        // Trả về order đã tạo
        order.ID = orderID
        c.JSON(http.StatusCreated, gin.H{
            "message": "Đơn hàng được tạo thành công",
            "order":   order,
            "user":    user,
        })
    }
}

// Lấy đơn hàng theo ID
func GetOrder(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
            return
        }

        order, err := repository.GetOrderByID(db, id)
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi server"})
            return
        }

        c.JSON(http.StatusOK, order)
    }
}

// Lấy tất cả đơn hàng của một user
func GetOrdersByUser(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDStr := c.Param("user_id")
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User ID không hợp lệ"})
            return
        }

        // Kiểm tra user có tồn tại không
        _, err = service.CheckUserExists(userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        orders, err := repository.GetOrdersByUserID(db, userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách đơn hàng"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "user_id": userID,
            "orders":  orders,
        })
    }
}

// Cập nhật status đơn hàng
func UpdateOrderStatus(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
            return
        }

        var req struct {
            Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
            // oneof = chỉ chấp nhận các giá trị trong danh sách
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Kiểm tra order có tồn tại không
        _, err = repository.GetOrderByID(db, id)
        if err != nil {
            if err == sql.ErrNoRows {
                c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi server"})
            return
        }

        // Cập nhật status
        err = repository.UpdateOrderStatus(db, id, req.Status)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật status"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "message": "Cập nhật status thành công",
            "order_id": id,
            "status": req.Status,
        })
    }
}
