package messaging

import (
	"testing"
	"time"
)

func TestRabbitMQConnection(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"

	rabbitMQ, err := NewRabbitMQ(url)
	if err != nil {
		t.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	defer rabbitMQ.Close()

	if rabbitMQ.Connection == nil {
		t.Fatal("expected RabbitMQ connection")
	}

	if rabbitMQ.Channel == nil {
		t.Fatal("expected RabbitMQ channel")
	}
}

func TestDeclareQueue(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"

	rabbitMQ, err := NewRabbitMQ(url)
	if err != nil {
		t.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	defer rabbitMQ.Close()

	if err := rabbitMQ.DeclareQueue("test.queue"); err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}
}

func TestPublish(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	queueName := "test.queue"
	message := []byte(`{"message":"hello rabbitmq"}`)

	rabbitMQ, err := NewRabbitMQ(url)
	if err != nil {
		t.Fatalf("failed to connect RabbitMQ: %v", err)
	}

	defer rabbitMQ.Close()

	if err := rabbitMQ.DeclareQueue(queueName); err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}

	if err := rabbitMQ.Publish(queueName, message); err != nil {
		t.Fatalf("failed to publish message: %v", err)
	}
}

func TestPublishAndConsume(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	queueName := "test.consume.queue"
	message := []byte(`{"message":"hello consumer"}`)

	rabbitMQ, err := NewRabbitMQ(url)
	if err != nil {
		t.Fatalf("failed to connect RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	if err := rabbitMQ.DeclareQueue(queueName); err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}

	messages, err := rabbitMQ.Consume(queueName)
	if err != nil {
		t.Fatalf("failed to consume queue: %v", err)
	}

	if err := rabbitMQ.Publish(queueName, message); err != nil {
		t.Fatalf("failed to publish message: %v", err)
	}

	select {
	case delivery := <-messages:
		if string(delivery.Body) != string(message) {
			t.Fatalf("expected %q, got %q", message, delivery.Body)
		}

		if err := delivery.Ack(false); err != nil {
			t.Fatalf("failed to acknowledge message: %v", err)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for message")
	}
}
