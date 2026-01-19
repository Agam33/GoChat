package rabbitmq

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, exchangeName, kind, routingKey string, data []byte) error
}

type publisher struct {
	conn *amqp.Connection
}

func NewPublisher(conn *amqp.Connection) Publisher {
	return &publisher{
		conn: conn,
	}
}

func (p *publisher) Publish(ctx context.Context, exchangeName, kind, routingKey string, data []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		log.Printf("fail create a channel:  %v", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // exchangeName
		kind,         // kind
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("fail declare exchange:  %v", err)
		return err
	}

	pub := amqp.Publishing{
		Body:         data,
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.NewString(),
		Timestamp:    time.Now(),
	}

	return ch.PublishWithContext(ctx, exchangeName, routingKey, false, false, pub)
}
