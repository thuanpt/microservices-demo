FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod init email-service && go mod tidy && go build -o email-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/email-service .
CMD ["./email-service"]
