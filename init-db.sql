-- Tạo các databases cho từng service
CREATE DATABASE IF NOT EXISTS user_service_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS product_service_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS order_service_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant quyền cho user microservices
GRANT ALL PRIVILEGES ON user_service_db.* TO 'microservices'@'%';
GRANT ALL PRIVILEGES ON product_service_db.* TO 'microservices'@'%';
GRANT ALL PRIVILEGES ON order_service_db.* TO 'microservices'@'%';

FLUSH PRIVILEGES;
