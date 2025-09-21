package eventbus

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type MessageHandler func(ctx context.Context, msg amqp091.Delivery) error

// Consume démarre un consumer générique avec retry/backoff, lié au context
func (r *RabbitMQ) Consume(ctx context.Context, queue string, handler MessageHandler) error {
	if r.channel == nil {
		return errors.New(errors.CodeInternal, "RabbitMQ channel is nil")
	}

	// Context attaché au consumer  permet shutdown propre
	msgs, err := r.channel.ConsumeWithContext(
		ctx,
		queue,
		"",    // consumer tag
		false, // autoAck false → on ack manuellement
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to register RabbitMQ consumer")
	}

	// Goroutine  lit jusqu’à ce que ctx soit annulé
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("[RabbitMQ] consumer for queue=%s shutting down", queue)
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("[RabbitMQ] channel closed for queue=%s", queue)
					return
				}
				if err := processWithRetry(ctx, handler, msg, 3); err != nil {
					log.Printf("[RabbitMQ] failed after retries: %v", err)
					_ = msg.Nack(false, false) // rejet définitif
				} else {
					_ = msg.Ack(false)
				}
			}
		}
	}()
	return nil
}

// applique une stratégie retry simple
func processWithRetry(ctx context.Context, handler MessageHandler, msg amqp091.Delivery, maxRetries int) error {
	var err error
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		if err = handler(ctx, msg); err == nil {
			return nil
		}
		log.Printf("[RabbitMQ] handler error, retrying... (%d/%d): %v", i+1, maxRetries, err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			backoff *= 2 // backoff exponentiel
		}
	}
	return err
}
