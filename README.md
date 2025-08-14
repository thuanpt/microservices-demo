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

- **API Gateway**: Cổng chính xử lý routing, authentication, và rate limiting
- **User Service**: Quản lý người dùng (đăng ký, đăng nhập, CRUD operations)├── product-service/
│   ├── .dockerignore
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── .env                 # Environment variables
│   ├── handler/
│   │   └── product.go      # HTTP handlers
│   ├── migrations/
│   │   ├── 001_create_products_table.down.sql
│   │   └── 001_create_products_table.up.sql
│   ├── model/
│   │   └── product.go      # Data models
│   ├── repository/
│   │   └── product.go      # Database operations
│   └── scripts/
│       └── migrate.go      # Migration script
└── order-service/ipts/
│       └── migrate.go      # Migration script
└── order-service/uct Service**: Quản lý sản phẩm (CRUD operations)
- **Order Service**: Quản lý đơn hàng (tạo, cập nhật, theo dõi đơn hàng)
- **JWT Authentication**: Hệ thống xác thực và phân quyền người dùng
- **Migration System**: Quản lý cấu trúc database với scripts tự động

## 🏗️ Kiến trúc

```
                    Client Requests
                          │
                          ▼
┌─────────────────────────────────────────┐
│              API Gateway                │
│  ┌─────────────────┐ ┌─────────────────┐│
│  │  Authentication │ │  Rate Limiting  ││
│  │   & JWT Auth    │ │   & Routing     ││
│  └─────────────────┘ └─────────────────┘│
└─────────────────────────────────────────┘
                          │
       ┌──────────────────┼──────────────────┐
       │                  │                  │
┌─────▼──────┐ ┌─────────▼────┐ ┌──────────▼──┐
│    User    │ │   Product    │ │    Order    │
│  Service   │ │   Service    │ │   Service   │
│ (Port 8001)│ │ (Port 8002)  │ │(Port 8003)  │
└────────────┘ └──────────────┘ └─────────────┘
       │                  │                  │
       └──────────────────┼──────────────────┘
                          │
                   ┌─────▼──────┐
                   │   MySQL    │
                   │ Database   │
                   └────────────┘
```

## 🛠️ Công nghệ sử dụng

- **Go** 1.24.3
- **Gin** - HTTP web framework
- **MySQL** - Database
- **JWT** - JSON Web Token cho authentication
- **bcrypt** - Mã hóa mật khẩu
- **godotenv** - Quản lý biến môi trường
- **Rate Limiting** - Giới hạn số lượng request

### Dependencies chính:

```go
github.com/gin-gonic/gin v1.10.1
github.com/go-sql-driver/mysql v1.9.3
github.com/joho/godotenv v1.5.1
github.com/golang-jwt/jwt/v4 v4.5.0
golang.org/x/crypto v0.23.0
golang.org/x/time v0.5.0
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

# API Gateway
cd api-gateway
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

**API Gateway:**
```bash
cd api-gateway
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

**API Gateway (.env):**
```env
# Server Configuration
APP_PORT=8000

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# External Services
USER_SERVICE_URL=http://localhost:8001
PRODUCT_SERVICE_URL=http://localhost:8002
ORDER_SERVICE_URL=http://localhost:8003

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m
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

### 🐳 Chạy với Docker (Khuyến nghị)

**Cách nhanh nhất để chạy toàn bộ hệ thống:**

```bash
# Clone và chạy ngay
git clone https://github.com/thuanpt/microservices-demo.git
cd microservices-demo

# Build và chạy tất cả services
docker-compose up --build -d

# Kiểm tra logs
docker-compose logs -f
```

**API sẽ có sẵn tại `http://localhost:8000`**

### 💻 Chạy Manual (Development)

**Yêu cầu trước:** MySQL đang chạy và đã cấu hình .env files

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

# Terminal 4 - API Gateway
cd api-gateway
go run main.go
```

**Services sẽ chạy tại:**
- API Gateway: `http://localhost:8000` (Main entry point)
- User Service: `http://localhost:8001` (Internal)
- Product Service: `http://localhost:8002` (Internal)
- Order Service: `http://localhost:8003` (Internal)

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

# API Gateway
cd api-gateway
air
```

## 📚 API Documentation

**⚡ Tất cả API requests đều được gửi thông qua API Gateway tại `http://localhost:8000`**

### Authentication

#### 1. Đăng ký người dùng
```
POST /api/v1/auth/register
Content-Type: application/json

{
    "name": "Nguyen Van A",
    "email": "user@example.com",
    "password": "password123"
}
```

#### 2. Đăng nhập
```
POST /api/v1/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}

Response:
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "name": "Nguyen Van A",
        "email": "user@example.com"
    }
}
```

### User Service Endpoints

**🔒 Yêu cầu Authentication: Thêm header `Authorization: Bearer <token>`**

#### 3. Lấy thông tin người dùng hiện tại
```
GET /api/v1/users/me
Authorization: Bearer <token>
```

#### 4. Cập nhật thông tin người dùng
```
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Updated Name",
    "email": "updated@example.com"
}
```

#### 5. Lấy danh sách tất cả người dùng (Admin only)
```
GET /api/v1/users
Authorization: Bearer <admin_token>
```

### Product Service Endpoints

**🔒 Yêu cầu Authentication: Thêm header `Authorization: Bearer <token>`**

#### 6. Tạo sản phẩm mới (Admin only)
```
POST /api/v1/products
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 100
}
```

#### 7. Lấy thông tin sản phẩm
```
GET /api/v1/products/:id
Authorization: Bearer <token>
```

#### 8. Cập nhật sản phẩm (Admin only)
```
PUT /api/v1/products/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "name": "Updated Product Name",
    "description": "Updated description",
    "price": 1199.99,
    "stock": 50
}
```

#### 9. Xóa sản phẩm (Admin only)
```
DELETE /api/v1/products/:id
Authorization: Bearer <admin_token>
```

#### 10. Lấy danh sách tất cả sản phẩm
```
GET /api/v1/products
Authorization: Bearer <token>
```

#### 11. Tìm kiếm sản phẩm
```
GET /api/v1/products/search?q=keyword
Authorization: Bearer <token>
```

### Order Service Endpoints

**🔒 Yêu cầu Authentication: Thêm header `Authorization: Bearer <token>`**

#### 12. Tạo đơn hàng mới
```
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
    "product_id": 1,
    "quantity": 2
}
```

#### 13. Lấy thông tin đơn hàng
```
GET /api/v1/orders/:id
Authorization: Bearer <token>
```

#### 14. Cập nhật trạng thái đơn hàng (Admin only)
```
PUT /api/v1/orders/:id/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
    "status": "confirmed"
}
```

#### 15. Hủy đơn hàng
```
DELETE /api/v1/orders/:id
Authorization: Bearer <token>
```

#### 16. Lấy danh sách đơn hàng của user hiện tại
```
GET /api/v1/orders/my-orders
Authorization: Bearer <token>
```

#### 17. Lấy tất cả đơn hàng (Admin only)
```
GET /api/v1/orders
Authorization: Bearer <admin_token>
```

### API Gateway Features

#### Rate Limiting
- **Giới hạn**: 100 requests per minute per IP
- **Response khi vượt giới hạn**: HTTP 429 Too Many Requests

#### Authentication Middleware
- **JWT Token Validation**: Tự động xác thực token cho các protected routes
- **User Context**: Tự động inject thông tin user vào request headers cho các microservices

#### Request/Response Logging
- Log tất cả requests đi qua gateway
- Performance monitoring và error tracking

### Response Examples

**Login Success Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "expires_in": "24h",
    "user": {
        "id": 1,
        "name": "Nguyen Van A",
        "email": "user@example.com",
        "created_at": "2025-08-08T10:00:00Z"
    }
}
```

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
    "error": "Unauthorized",
    "message": "Token không hợp lệ hoặc đã hết hạn"
}
```

**Rate Limit Error:**
```json
{
    "error": "Too Many Requests",
    "message": "Vượt quá giới hạn request. Vui lòng thử lại sau."
}
```

## 📁 Cấu trúc thư mục

```
microservices-demo/
├── .gitignore
├── README.md
├── docker-compose.yml       # Docker Compose configuration
├── init-db.sql             # Database initialization script
├── api-gateway/
│   ├── .dockerignore
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── .env                 # Environment variables
│   ├── config/
│   │   └── config.go       # Configuration management
│   ├── middleware/
│   │   ├── auth.go         # JWT authentication
│   │   └── ratelimit.go    # Rate limiting
│   └── proxy/
│       └── proxy.go        # Request routing and proxying
├── user-service/
│   ├── .dockerignore
│   ├── Dockerfile
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
│       ├── hash.go         # Utility functions
│       └── jwt.go          # JWT token utilities
├── product-service/
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
    ├── .dockerignore
    ├── Dockerfile
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

**Docker Files:**
- **`docker-compose.yml`**: Orchestration cho tất cả services và database
- **`init-db.sql`**: Script khởi tạo database và tables
- **`Dockerfile`**: Container build instructions cho từng service  
- **`.dockerignore`**: Files/folders bị ignore khi build image

**API Gateway:**
- **`config/`**: Quản lý cấu hình ứng dụng
- **`middleware/`**: Authentication, rate limiting, logging
- **`proxy/`**: Request routing và proxying tới các microservices

**Microservices:**
- **`main.go`**: Entry point của ứng dụng
- **`handler/`**: Xử lý HTTP requests và responses
- **`model/`**: Định nghĩa cấu trúc dữ liệu
- **`repository/`**: Tương tác với database
- **`migrations/`**: SQL files để tạo/xóa database tables
- **`scripts/`**: Migration scripts để chạy database migrations
- **`utils/`**: Các hàm tiện ích (hash password, JWT, validation, ...)
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

cd api-gateway
go test ./... -v
```

## 🐳 Docker

Dự án đã được containerized với Docker và Docker Compose để dễ dàng deploy và development.

### Quick Start với Docker

```bash
# Build và chạy tất cả services
docker-compose up --build

# Hoặc chạy background
docker-compose up --build -d

# Xem logs
docker-compose logs -f

# Stop tất cả services
docker-compose down

# Stop và xóa volumes
docker-compose down -v
```

### Services sẽ chạy tại:
- **API Gateway**: `http://localhost:8000`
- **User Service**: `http://localhost:8001` (Internal)
- **Product Service**: `http://localhost:8002` (Internal) 
- **Order Service**: `http://localhost:8003` (Internal)
- **MySQL Database**: `localhost:3306`

### Useful Docker Commands

```bash
# Build specific service
docker-compose build api-gateway

# Restart specific service
docker-compose restart user-service

# View logs of specific service
docker-compose logs -f api-gateway

# Execute command in running container
docker-compose exec api-gateway sh

# View running services
docker-compose ps

# Remove stopped containers and networks
docker-compose down --remove-orphans
```

### Development với Docker

```bash
# Development mode với hot reload (nếu có Air setup)
docker-compose -f docker-compose.dev.yml up

# Chỉ chạy database cho local development
docker-compose up mysql -d

# Build lại service sau khi thay đổi code
docker-compose build api-gateway
docker-compose up api-gateway -d

# Reset everything (containers, networks, volumes)
docker-compose down -v --remove-orphans
docker system prune -a
```

### Environment Variables trong Docker

Environment variables được quản lý thông qua:
1. **File `.env`** trong mỗi service directory
2. **docker-compose.yml** environment section
3. **Dockerfile ENV** statements

**Lưu ý**: Trong production, sử dụng Docker Secrets hoặc external config management.

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
- [x] JWT Authentication
- [x] API Gateway
- [x] Rate Limiting
- [x] Docker containerization
- [ ] Docker Compose for development
- [ ] Kubernetes deployment
- [ ] Monitoring và Logging
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Service Discovery
- [ ] Load Balancing
- [ ] Circuit Breaker Pattern
- [ ] Distributed Tracing
- [ ] Health Checks
- [ ] Graceful Shutdown
