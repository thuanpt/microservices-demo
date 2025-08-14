-- Migration tạo bảng orders và order_items
CREATE DATABASE IF NOT EXISTS order_service_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE order_service_db;

-- Bảng orders (đơn hàng chính)
CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,                    -- ID của user từ User Service
    total_amount DECIMAL(10,2) NOT NULL,     -- Tổng tiền đơn hàng
    status ENUM('pending', 'confirmed', 'cancelled', 'completed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Bảng order_items (chi tiết sản phẩm trong đơn hàng)
CREATE TABLE IF NOT EXISTS order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,                   -- Foreign key tới bảng orders
    product_id INT NOT NULL,                 -- ID của product từ Product Service
    product_name VARCHAR(255) NOT NULL,      -- Lưu tên sản phẩm tại thời điểm đặt hàng
    price DECIMAL(10,2) NOT NULL,           -- Giá sản phẩm tại thời điểm đặt hàng
    quantity INT NOT NULL,                   -- Số lượng sản phẩm
    subtotal DECIMAL(10,2) NOT NULL,        -- Thành tiền (price * quantity)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;