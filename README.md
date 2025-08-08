# Microservices Demo

Dự án demo microservices được xây dựng bằng Go, sử dụng Gin framework và MySQL database.

## 📋 Mục lục

- [Tổng quan](#tổng-quan)
- [Kiến trúc](#kiến-trúc)
- [Công nghệ sử dụng](#công-nghệ-sử-dụng)
- [Cài đặt](#cài-đặt)
- [Cấu hình](#cấu-hình)
- [Chạy ứng dụng](#chạy-ứng-dụng)
- [API Documentation](#api-documentation)
- [Cấu trúc thư mục](#cấu-trúc-thư-mục)

## 🎯 Tổng quan

Dự án này là một ví dụ về kiến trúc microservices sử dụng Go. Hiện tại bao gồm:

- **User Service**: Quản lý người dùng (đăng ký, đăng nhập, CRUD operations)

## 🏗️ Kiến trúc

```
┌─────────────────┐
│   API Gateway   │
└─────────────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼──┐
│ User  │ │ ... │
│Service│ │     │
└───────┘ └─────┘
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

### Bước 2: Cài đặt dependencies cho User Service

```bash
cd user-service
go mod download
```

## ⚙️ Cấu hình

### Database Setup

1. Tạo database MySQL:

```sql
CREATE DATABASE microservices_demo;
USE microservices_demo;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Environment Variables

1. Copy file mẫu và cấu hình:

```bash
cd user-service
cp .env.example .env
```

2. Chỉnh sửa file `.env` với thông tin thực tế của bạn:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_username
DB_PASS=your_db_password
DB_NAME=your_db_name

# Server Configuration
APP_PORT=8001
```

**⚠️ Lưu ý**: File `.env` chứa thông tin nhạy cảm và đã được thêm vào `.gitignore`. Không bao giờ commit file này lên repository!

## 🚀 Chạy ứng dụng

### Chạy User Service

```bash
cd user-service
go run main.go
```

Server sẽ chạy tại: `http://localhost:8080`

### Chạy với development mode

```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

## 📚 API Documentation

### User Service Endpoints

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

### Response Examples

**Success Response:**
```json
{
    "id": 1,
    "name": "Nguyen Van A",
    "email": "user@example.com",
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
└── user-service/
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── .env                 # Environment variables
    ├── handler/
    │   └── user.go         # HTTP handlers
    ├── model/
    │   └── user.go         # Data models
    ├── repository/
    │   └── user.go         # Database operations
    └── utils/
        └── hash.go         # Utility functions
```

### Mô tả các thành phần:

- **`main.go`**: Entry point của ứng dụng
- **`handler/`**: Xử lý HTTP requests và responses
- **`model/`**: Định nghĩa cấu trúc dữ liệu
- **`repository/`**: Tương tác với database
- **`utils/`**: Các hàm tiện ích (hash password, validation, ...)

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./handler -v
```

## 🐳 Docker (Coming Soon)

```bash
# Build Docker image
docker build -t user-service .

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
golangci-lint run
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

- [ ] Authentication Service
- [ ] Product Service  
- [ ] Order Service
- [ ] API Gateway
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Monitoring và Logging
- [ ] Unit Tests
- [ ] Integration Tests
