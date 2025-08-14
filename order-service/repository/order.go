package repository

import (
    "database/sql"
    "order-service/model"
)

// Tạo đơn hàng mới (sử dụng transaction để đảm bảo data consistency)
func CreateOrder(db *sql.DB, order *model.Order) (int, error) {
    // Bắt đầu transaction
    tx, err := db.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback() // Rollback nếu có lỗi

    // Insert order chính
    stmt := "INSERT INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)"
    result, err := tx.Exec(stmt, order.UserID, order.TotalAmount, order.Status)
    if err != nil {
        return 0, err
    }

    // Lấy order_id vừa tạo
    orderID, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    // Insert từng order_item
    itemStmt := "INSERT INTO order_items (order_id, product_id, product_name, price, quantity, subtotal) VALUES (?, ?, ?, ?, ?, ?)"
    for _, item := range order.Items {
        _, err := tx.Exec(itemStmt, orderID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.Subtotal)
        if err != nil {
            return 0, err
        }
    }

    // Commit transaction nếu tất cả thành công
    if err := tx.Commit(); err != nil {
        return 0, err
    }

    return int(orderID), nil
}

// Lấy đơn hàng theo ID (kèm danh sách items)
func GetOrderByID(db *sql.DB, id int) (*model.Order, error) {
    var order model.Order
    
    // Lấy thông tin order chính
    err := db.QueryRow(`
        SELECT id, user_id, total_amount, status, created_at, updated_at 
        FROM orders WHERE id = ?`, id).
        Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
    
    if err != nil {
        return nil, err
    }

    // Lấy danh sách order_items
    rows, err := db.Query(`
        SELECT id, order_id, product_id, product_name, price, quantity, subtotal, created_at 
        FROM order_items WHERE order_id = ?`, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Loop qua từng item và add vào order.Items
    var items []model.OrderItem
    for rows.Next() {
        var item model.OrderItem
        err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, 
                        &item.Price, &item.Quantity, &item.Subtotal, &item.CreatedAt)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    order.Items = items

    return &order, nil
}

// Lấy tất cả đơn hàng của một user
func GetOrdersByUserID(db *sql.DB, userID int) ([]model.Order, error) {
    var orders []model.Order
    
    rows, err := db.Query(`
        SELECT id, user_id, total_amount, status, created_at, updated_at 
        FROM orders WHERE user_id = ? ORDER BY created_at DESC`, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var order model.Order
        err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }

    return orders, nil
}

// Cập nhật status đơn hàng
func UpdateOrderStatus(db *sql.DB, orderID int, status string) error {
    stmt := "UPDATE orders SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
    _, err := db.Exec(stmt, status, orderID)
    return err
}
