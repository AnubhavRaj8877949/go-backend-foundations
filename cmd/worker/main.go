package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-backend-foundation/internal/config"
	"example.com/go-backend-foundation/pkg/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	rabbitMQ, err := messaging.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("connect RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	if err := rabbitMQ.DeclareQueue("default"); err != nil {
		log.Fatalf("declare queue: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	err = rabbitMQ.ConsumeAndHandle(
		ctx,
		"default",
		func(message []byte) error {
			log.Printf("received message: %s", message)
			return nil
		},
	)
	if err != nil && err != context.Canceled {
		log.Fatalf("consume messages: %v", err)
	}
}
