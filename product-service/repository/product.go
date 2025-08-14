package repository

import (
    "database/sql"
    "product-service/model"
)

// Thêm product mới vào DB
func InsertProduct(db *sql.DB, product *model.Product) (int, error) {
    // Câu SQL INSERT, dấu ? là placeholder để tránh SQL injection
    stmt := "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"
    
    // db.Exec thực hiện câu lệnh SQL không trả về dữ liệu (INSERT, UPDATE, DELETE)
    result, err := db.Exec(stmt, product.Name, product.Description, product.Price, product.Stock)
    if err != nil {
        return 0, err // Return 0 và error nếu có lỗi
    }
    
    // Lấy ID vừa được insert (LastInsertId là method của sql.Result)
    id, _ := result.LastInsertId()
    return int(id), nil // Convert int64 thành int và return
}

// Lấy tất cả products
func GetAllProducts(db *sql.DB) ([]model.Product, error) {
    // []model.Product là slice (mảng động) chứa các Product
    var products []model.Product
    
    // db.Query thực hiện câu SELECT và trả về nhiều rows
    rows, err := db.Query("SELECT id, name, description, price, stock, created_at FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close() // defer nghĩa là thực hiện cuối cùng trước khi return
    
    // Loop qua từng row
    for rows.Next() {
        var product model.Product
        // Scan data từ row vào struct Product
        err := rows.Scan(&product.ID, &product.Name, &product.Description, 
                        &product.Price, &product.Stock, &product.CreatedAt)
        if err != nil {
            return nil, err
        }
        // append thêm product vào slice products
        products = append(products, product)
    }
    
    return products, nil
}

// Lấy product theo ID
func GetProductByID(db *sql.DB, id int) (*model.Product, error) {
    var product model.Product
    
    // QueryRow trả về 1 row duy nhất
    err := db.QueryRow("SELECT id, name, description, price, stock, created_at FROM products WHERE id = ?", id).
        Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt)
    
    if err != nil {
        return nil, err // Return nil pointer nếu không tìm thấy
    }
    
    return &product, nil // Return pointer tới product
}

// Cập nhật stock của product
func UpdateProductStock(db *sql.DB, id int, stock int) error {
    stmt := "UPDATE products SET stock = ? WHERE id = ?"
    _, err := db.Exec(stmt, stock, id)
    return err // Chỉ return error, không cần return data
}