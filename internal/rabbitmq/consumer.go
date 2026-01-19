package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerDispatcher interface {
	Dispatch(ctx context.Context, msg amqp.Delivery)
}

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery

	exchangeName string
	kind         string
	queueName    string
	bindKey      string

	dispatcher ConsumerDispatcher
}

func NewConsumer(conn *amqp.Connection, dispatcher ConsumerDispatcher, exchangeName string, kind string, queueName string, bindKey string) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	args := amqp.Table{
		"x-queue-type": "quorum",
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		bindKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	deliveries, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:       conn,
		deliveries: deliveries,
		dispatcher: dispatcher,

		exchangeName: exchangeName,
		kind:         kind,
		queueName:    queueName,
		bindKey:      bindKey,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		defer c.channel.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-c.deliveries:
				if !ok {
					return // close the channel
				}
				c.dispatcher.Dispatch(ctx, delivery)
			case <-ctx.Done():
				c.channel.Close()
				return
			}
		}
	}()
}
