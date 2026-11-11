package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection    *amqp.Connection
	Channel       *amqp.Channel
	Confirmations <-chan amqp.Confirmation
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	if err := channel.Confirm(false); err != nil {
		_ = channel.Close()
		_ = connection.Close()

		return nil, fmt.Errorf("enable publisher confirms: %w", err)
	}

	confirmations := channel.NotifyPublish(
		make(chan amqp.Confirmation, 1),
	)
	return &RabbitMQ{
		Connection:    connection,
		Channel:       channel,
		Confirmations: confirmations,
	}, nil
}

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.Channel.QueueDeclare(
		name,  // queue का नाम
		true,  //  RabbitMQ restart के बाद queue रहेगी
		false, //unused होने पर खुद delete नहीं होगी
		false, //केवल एक connection के लिए नहीं
		false, // server confirmation का wait करेगा
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue %q: %w", name, err)
	}

	return nil
}

func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	err := r.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("publish message to queue %q: %w", queueName, err)
	}
	confirmation := <-r.Confirmations
	if !confirmation.Ack {
		return fmt.Errorf("RabbitMQ rejected published message")
	}

	return nil
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := r.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume queue %q: %w", queueName, err)
	}

	return messages, nil
}

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			return fmt.Errorf("close RabbitMQ channel: %w", err)
		}
	}

	if r.Connection != nil {
		if err := r.Connection.Close(); err != nil {
			return fmt.Errorf("close RabbitMQ connection: %w", err)
		}
	}

	return nil
}
