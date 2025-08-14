package handler

import (
    "database/sql"
    "net/http"
    "strconv"                          // Package convert string <-> number
    "product-service/model"
    "product-service/repository"
    "github.com/gin-gonic/gin"
)

// Tạo product mới
func CreateProduct(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req model.CreateProductRequest
        
        // ShouldBindJSON tự động validate theo tag binding trong struct
        if err := c.ShouldBindJSON(&req); err != nil {
            // c.JSON(status_code, data) để trả response JSON
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        // Tạo Product từ request
        product := model.Product{
            Name:        req.Name,
            Description: req.Description,
            Price:       req.Price,
            Stock:       req.Stock,
        }
        
        // Lưu vào DB
        id, err := repository.InsertProduct(db, &product)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo sản phẩm"})
            return
        }
        
        product.ID = id
        c.JSON(http.StatusCreated, product) // Trả về product vừa tạo
    }
}

// Lấy danh sách tất cả products
func GetAllProducts(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        products, err := repository.GetAllProducts(db)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách sản phẩm"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{"products": products})
    }
}

// Lấy product theo ID
func GetProduct(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // c.Param("id") lấy parameter từ URL path (ví dụ: /products/123)
        idStr := c.Param("id")
        
        // strconv.Atoi convert string thành int
        id, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
            return
        }
        
        product, err := repository.GetProductByID(db, id)
        if err != nil {
            if err == sql.ErrNoRows { // sql.ErrNoRows là lỗi khi không tìm thấy row
                c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi server"})
            return
        }
        
        c.JSON(http.StatusOK, product)
    }
}

// Cập nhật stock của product
func UpdateProductStock(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
            return
        }
        
        var req struct {
            Stock int `json:"stock" binding:"min=0"` // Struct inline cho request body
        }
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        err = repository.UpdateProductStock(db, id, req.Stock)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật stock"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Cập nhật stock thành công"})
    }
}