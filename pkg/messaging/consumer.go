
package messaging

import (
    "context"
    "fmt"
)

func (r *RabbitMQ) ConsumeAndHandle(
    ctx context.Context,
    queueName string,
    handler func([]byte) error,
) error {
    if handler == nil {
        return fmt.Errorf("message handler is nil")
    }

    messages, err := r.Consume(queueName)
    if err != nil {
        return err
    }

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        case delivery, ok := <-messages:
            if !ok {
                return fmt.Errorf("message channel closed")
            }

            if err := handler(delivery.Body); err != nil {
                if nackErr := delivery.Nack(false, false); nackErr != nil {
                    return fmt.Errorf("nack message: %w", nackErr)
                }
                continue
            }

            if err := delivery.Ack(false); err != nil {
                return fmt.Errorf("ack message: %w", err)
            }
        }
    }
}