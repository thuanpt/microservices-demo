# Microservices Demo

Dự án demo microservices được xây dựng bằng Go, sử dụng Gin framework và MySQL database.

## 📋 Mục lục

- [Tổng quan](#tổng-quan)
- [Kiến trúc](#kiến-trúc)
- [Công nghệ sử dụng](#công-nghệ-sử-dụng)
- [Cài đặt](#cài-đặt)
- [Cấu hình](#cấu-hình)
- [Migration](#migration)
- [Chạy ứng dụng](#chạy-ứng-dụng)
- [API Documentation](#api-documentation)
- [Cấu trúc thư mục](#cấu-trúc-thư-mục)

## 🎯 Tổng quan

Dự án này là một ví dụ về kiến trúc microservices sử dụng Go. Hiện tại bao gồm:

- **User Service**: Quản lý người dùng (đăng ký, đăng nhập, CRUD operations)
- **Product Service**: Quản lý sản phẩm (CRUD operations)
- **Order Service**: Quản lý đơn hàng (tạo, cập nhật, theo dõi đơn hàng)
- **Migration System**: Quản lý cấu trúc database với scripts tự động

## 🏗️ Kiến trúc

```
┌─────────────────┐
│   API Gateway   │
└─────────────────┘
         │
    ┌────┴────┬────────┬────────┐
    │         │        │        │
┌───▼───┐ ┌──▼────┐ ┌──▼────┐ ┌──▼──┐
│ User  │ │Product│ │ Order │ │ ... │
│Service│ │Service│ │Service│ │     │
└───────┘ └───────┘ └───────┘ └─────┘
```

## 🛠️ Công nghệ sử dụng

- **Go** 1.24.3
- **Gin** - HTTP web framework
- **MySQL** - Database
- **bcrypt** - Mã hóa mật khẩu
- **godotenv** - Quản lý biến môi trường

### Dependencies chính:

```go
github.com/gin-gonic/gin v1.10.1
github.com/go-sql-driver/mysql v1.9.3
github.com/joho/godotenv v1.5.1
golang.org/x/crypto v0.23.0
```

## 📥 Cài đặt

### Yêu cầu hệ thống

- Go 1.24.3 hoặc cao hơn
- MySQL 8.0 hoặc cao hơn
- Git

### Bước 1: Clone repository

```bash
git clone https://github.com/thuanpt/microservices-demo.git
cd microservices-demo
```

### Bước 2: Cài đặt dependencies

```bash
# User Service
cd user-service
go mod download
cd ..

# Product Service
cd product-service
go mod download
cd ..

# Order Service
cd order-service
go mod download
cd ..
```

## ⚙️ Cấu hình

### Database Setup

1. Tạo database MySQL:

```sql
CREATE DATABASE microservices_demo;
```

### Environment Variables

1. Copy và cấu hình file `.env` cho từng service:

**User Service:**
```bash
cd user-service
cp .env.example .env
```

**Product Service:**
```bash
cd product-service
cp .env.example .env
```

**Order Service:**
```bash
cd order-service
cp .env.example .env
```

2. Chỉnh sửa file `.env` với thông tin thực tế của bạn:

**User Service (.env):**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_username
DB_PASS=your_db_password
DB_NAME=microservices_demo

# Server Configuration
APP_PORT=8001
```

**Product Service (.env):**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_username
DB_PASS=your_db_password
DB_NAME=microservices_demo

# Server Configuration
APP_PORT=8002
```

**Order Service (.env):**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_username
DB_PASS=your_db_password
DB_NAME=microservices_demo

# Server Configuration
APP_PORT=8003

# External Services
USER_SERVICE_URL=http://localhost:8001
PRODUCT_SERVICE_URL=http://localhost:8002
```

**⚠️ Lưu ý**: File `.env` chứa thông tin nhạy cảm và đã được thêm vào `.gitignore`. Không bao giờ commit file này lên repository!

## 🗄️ Migration

Hệ thống migration cho phép tạo và cập nhật cấu trúc database một cách có tổ chức.

### User Service Migration

Migration sẽ tạo bảng `users` với cấu trúc:

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Product Service Migration

Migration sẽ tạo bảng `products` với cấu trúc:

```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Order Service Migration

Migration sẽ tạo bảng `orders` với cấu trúc:

```sql
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Chạy Migration

Để chạy migration, sử dụng các lệnh sau:

**User Service Migration:**
```bash
# Chạy migration up (tạo tables)
cd user-service/scripts
go run migrate.go up

# Chạy rollback migration (xóa tables) 
go run migrate.go down
```

**Product Service Migration:**
```bash
# Di chuyển vào thư mục scripts của product-service
cd product-service/scripts

# Chạy migration up (tạo database và bảng products)
go run migrate.go up

# Nếu cần rollback (xóa bảng products)
go run migrate.go down
```

**Order Service Migration:**
```bash
# Di chuyển vào thư mục scripts của order-service
cd order-service/scripts

# Chạy migration up (tạo database và bảng orders)
go run migrate.go up

# Nếu cần rollback (xóa bảng orders)
go run migrate.go down
```

## 🚀 Chạy ứng dụng

### Chạy tất cả services

```bash
# Terminal 1 - User Service
cd user-service
go run main.go

# Terminal 2 - Product Service
cd product-service
go run main.go

# Terminal 3 - Order Service
cd order-service
go run main.go
```

**Services sẽ chạy tại:**
- User Service: `http://localhost:8001`
- Product Service: `http://localhost:8002`
- Order Service: `http://localhost:8003`

### Chạy với development mode

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# User Service
cd user-service
air

# Product Service  
cd product-service
air

# Order Service
cd order-service
air
```

## 📚 API Documentation

### User Service Endpoints (Port 8001)

#### 1. Đăng ký người dùng
```
POST /register
Content-Type: application/json

{
    "name": "Nguyen Van A",
    "email": "user@example.com",
    "password": "password123"
}
```

#### 2. Đăng nhập
```
POST /login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

#### 3. Lấy thông tin người dùng
```
GET /user/:id
```

#### 4. Cập nhật thông tin người dùng
```
PUT /user/:id
Content-Type: application/json

{
    "name": "Updated Name",
    "email": "updated@example.com"
}
```

#### 5. Xóa người dùng
```
DELETE /user/:id
```

#### 6. Lấy danh sách tất cả người dùng
```
GET /users
```

### Product Service Endpoints (Port 8002)

#### 1. Tạo sản phẩm mới
```
POST /products
Content-Type: application/json

{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 100
}
```

#### 2. Lấy thông tin sản phẩm
```
GET /products/:id
```

#### 3. Cập nhật sản phẩm
```
PUT /products/:id
Content-Type: application/json

{
    "name": "Updated Product Name",
    "description": "Updated description",
    "price": 1199.99,
    "stock": 50
}
```

#### 4. Xóa sản phẩm
```
DELETE /products/:id
```

#### 5. Lấy danh sách tất cả sản phẩm
```
GET /products
```

#### 6. Tìm kiếm sản phẩm
```
GET /products/search?q=keyword
```

### Order Service Endpoints (Port 8003)

#### 1. Tạo đơn hàng mới
```
POST /orders
Content-Type: application/json

{
    "user_id": 1,
    "product_id": 1,
    "quantity": 2
}
```

#### 2. Lấy thông tin đơn hàng
```
GET /orders/:id
```

#### 3. Cập nhật trạng thái đơn hàng
```
PUT /orders/:id/status
Content-Type: application/json

{
    "status": "confirmed"
}
```

#### 4. Hủy đơn hàng
```
DELETE /orders/:id
```

#### 5. Lấy danh sách đơn hàng của user
```
GET /orders/user/:user_id
```

#### 6. Lấy tất cả đơn hàng
```
GET /orders
```

### Response Examples

**User Success Response:**
```json
{
    "id": 1,
    "name": "Nguyen Van A",
    "email": "user@example.com",
    "created_at": "2025-08-08T10:00:00Z"
}
```

**Product Success Response:**
```json
{
    "id": 1,
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 100,
    "created_at": "2025-08-08T10:00:00Z"
}
```

**Order Success Response:**
```json
{
    "id": 1,
    "user_id": 1,
    "product_id": 1,
    "quantity": 2,
    "total_amount": 1999.98,
    "status": "pending",
    "created_at": "2025-08-08T10:00:00Z"
}
```

**Error Response:**
```json
{
    "error": "Dữ liệu không hợp lệ"
}
```

## 📁 Cấu trúc thư mục

```
microservices-demo/
├── .gitignore
├── README.md
├── user-service/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── .env                 # Environment variables
│   ├── handler/
│   │   └── user.go         # HTTP handlers
│   ├── migrations/
│   │   ├── 001_create_users_table.down.sql
│   │   └── 001_create_users_table.up.sql
│   ├── model/
│   │   └── user.go         # Data models
│   ├── repository/
│   │   └── user.go         # Database operations
│   ├── scripts/
│   │   └── migrate.go      # Migration script
│   └── utils/
│       └── hash.go         # Utility functions
└── product-service/
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── .env                 # Environment variables
    ├── handler/
    │   └── product.go      # HTTP handlers
    ├── migrations/
    │   ├── 001_create_products_table.down.sql
    │   └── 001_create_products_table.up.sql
    ├── model/
    │   └── product.go      # Data models
    ├── repository/
    │   └── product.go      # Database operations
    └── scripts/
        └── migrate.go      # Migration script
└── order-service/
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── .env                 # Environment variables
    ├── handler/
    │   └── order.go        # HTTP handlers
    ├── migrations/
    │   ├── 001_create_orders_table.down.sql
    │   └── 001_create_orders_table.up.sql
    ├── model/
    │   └── order.go        # Data models
    ├── repository/
    │   └── order.go        # Database operations
    ├── scripts/
    │   └── migrate.go      # Migration script
    └── service/
        └── external.go     # External service calls
```

### Mô tả các thành phần:

- **`main.go`**: Entry point của ứng dụng
- **`handler/`**: Xử lý HTTP requests và responses
- **`model/`**: Định nghĩa cấu trúc dữ liệu
- **`repository/`**: Tương tác với database
- **`migrations/`**: SQL files để tạo/xóa database tables
- **`scripts/`**: Migration scripts để chạy database migrations
- **`utils/`**: Các hàm tiện ích (hash password, validation, ...)
- **`service/`**: External service calls (chỉ có trong order-service)

## 🧪 Testing

```bash
# Run tests for all services
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
cd user-service
go test ./... -v

cd product-service
go test ./... -v

cd order-service
go test ./... -v
```

## 🐳 Docker (Coming Soon)

```bash
# Build Docker images
docker build -t user-service ./user-service
docker build -t product-service ./product-service
docker build -t order-service ./order-service

# Run with Docker Compose
docker-compose up -d
```

## 🔧 Development

### Code Style

Sử dụng `gofmt` để format code:

```bash
go fmt ./...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run ./...
```

## 🤝 Contributing

1. Fork repository
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 👥 Authors

- **thuanpt** - *Initial work* - [thuanpt](https://github.com/thuanpt)

## 🆘 Support

Nếu bạn gặp vấn đề, vui lòng tạo issue tại [GitHub Issues](https://github.com/thuanpt/microservices-demo/issues)

## 🎯 Roadmap

- [x] User Service
- [x] Product Service
- [x] Order Service
- [x] Database Migration System
- [ ] Authentication Service
- [ ] API Gateway
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Monitoring và Logging
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Service Discovery
- [ ] Load Balancing
