package main

import (
	"email-service/config"
	"email-service/service"
	"log"
)

func main() {
	cfg := config.Load()
	worker := service.NewEmailWorker(cfg)
	log.Println("[EmailService] Worker started, waiting for jobs...")
	worker.Start()
}
