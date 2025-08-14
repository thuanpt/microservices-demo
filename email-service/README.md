# Email Service (Background Worker)

Service này lắng nghe các email jobs từ RabbitMQ và gửi email (welcome, order confirmation...)

## Cấu trúc
- `main.go`: Khởi động worker
- `config/config.go`: Đọc config từ biến môi trường
- `service/email.go`: Worker lắng nghe queue và gửi email
- `utils/email_sender.go`: Hàm gửi email sử dụng gomail
- `Dockerfile`: Build image

## Biến môi trường
- `RABBITMQ_URL`: Kết nối RabbitMQ
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `FROM_EMAIL`: Thông tin SMTP

## Chạy local
```sh
go mod tidy
go run main.go
```

## Chạy bằng Docker
```sh
docker build -t email-service .
docker run --env-file .env email-service
```

## Queue
- Lắng nghe queue: `email.jobs.queue`
- Định dạng job:
```json
{
  "to": "user@example.com",
  "subject": "Welcome!",
  "body": "<b>Chào mừng bạn!</b>"
}
```
