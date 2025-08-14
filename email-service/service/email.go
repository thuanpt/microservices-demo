package service

import (
	"email-service/config"
	"email-service/utils"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type EmailWorker struct {
	cfg  *config.Config
	conn *amqp.Connection
	ch   *amqp.Channel
}

type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEmailWorker(cfg *config.Config) *EmailWorker {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("[EmailWorker] Lỗi kết nối RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[EmailWorker] Lỗi tạo channel RabbitMQ: %v", err)
	}
	return &EmailWorker{cfg: cfg, conn: conn, ch: ch}
}

func (w *EmailWorker) Start() {
	msgs, err := w.ch.Consume(
		"email.jobs.queue", // queue
		"",                 // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		log.Fatalf("[EmailWorker] Lỗi consume queue: %v", err)
	}
	for msg := range msgs {
		var job EmailJob
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("[EmailWorker] Lỗi decode job: %v", err)
			msg.Nack(false, false)
			continue
		}
		if err := utils.SendEmail(w.cfg, job.To, job.Subject, job.Body); err != nil {
			log.Printf("[EmailWorker] Lỗi gửi email: %v", err)
			msg.Nack(false, false)
		} else {
			log.Printf("[EmailWorker] Đã gửi email tới %s", job.To)
			msg.Ack(false)
		}
	}
}
