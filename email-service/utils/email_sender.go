package utils

import (
	"email-service/config"
	"fmt"
	"net/smtp"
)

func SendEmail(cfg *config.Config, to, subject, body string) error {
	// Thiết lập auth
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	// Tạo message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s", to, subject, body))
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	if err := smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, msg); err != nil {
		return fmt.Errorf("send mail error: %w", err)
	}
	return nil
}
