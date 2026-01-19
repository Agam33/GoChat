package rabbitmq

type RabbitMQConfig struct {
	Port     int
	User     string
	Host     string
	Password string
	VHost    string
}
