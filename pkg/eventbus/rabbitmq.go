package eventbus

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// initialise la connexion RabbitMQ
func InitRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to open RabbitMQ channel")
	}

	return &RabbitMQ{conn: conn, channel: ch}, nil
}

// envoie un message dans une exchange/topic
func (r *RabbitMQ) Publish(exchange, routingKey string, body []byte) error {
	if r.channel == nil {
		return errors.New(errors.CodeInternal, "RabbitMQ channel is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to publish RabbitMQ message")
	}
	log.Printf("[RabbitMQ] published message to %s:%s", exchange, routingKey)
	return nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
}
