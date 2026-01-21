package config

import (
	"go-chat/internal/database"
	"go-chat/internal/env"
	"go-chat/internal/jwt"
	"go-chat/internal/rabbitmq"
)

type Config struct {
	JWT      jwt.JwtConfig
	DBConfig database.DBConfig
	RabbitMQ rabbitmq.RabbitMQConfig
}

func NewAppConfig(env *env.Env) *Config {
	return &Config{
		JWT: jwt.JwtConfig{
			AccessExpire:  env.JWT.AccessExp,
			AccessSecret:  env.JWT.AccessSecret,
			RefreshExpire: env.JWT.RefreshExp,
			RefreshSecret: env.JWT.RefreshSecret,
		},
		DBConfig: database.DBConfig{
			User:     env.Database.User,
			Host:     env.Database.Host,
			Port:     env.Database.Port,
			DBName:   env.Database.DBName,
			Password: env.Database.Password,
			SslMode:  env.Database.SslMode,
		},
		RabbitMQ: rabbitmq.RabbitMQConfig{
			User:     env.MQ.User,
			Host:     env.MQ.Host,
			Port:     env.MQ.Port,
			Password: env.MQ.Password,
			VHost:    env.MQ.Vhost,
		},
	}
}
