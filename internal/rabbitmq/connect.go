package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConnection(config *RabbitMQConfig) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	conn, err := amqp.Dial(connAddr)
	if err != nil {
		fmt.Println("failed to connect rabbitmq")
		return nil, err
	}

	return conn, nil
}
