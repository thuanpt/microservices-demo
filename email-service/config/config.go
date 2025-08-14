package config

import (
	"os"
)

type Config struct {
	RabbitMQURL string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	FromEmail   string
}

func Load() *Config {
	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		SMTPHost:    getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
		FromEmail:   getEnv("FROM_EMAIL", "noreply@example.com"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
