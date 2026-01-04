package env

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	App struct {
		Port int `mapstructure:"APP_PORT"`
	}

	JWT struct {
		AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
		RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
		AccessExp     time.Duration `mapstructure:"JWT_ACCESS_EXP"`
		RefreshExp    time.Duration `mapstructure:"JWT_REFRESH_EXP"`
	}

	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		DBName   string `mapstructure:"DB_NAME"`
		Password string `mapstructure:"DB_PASSWORD"`
		SslMode  string `mapstructure:"DB_SSLMODE"`
	}
}

func NewEnv() (*Env, error) {
	v := viper.New()

	if os.Getenv("APP_ENV") != "prod" {
		v.SetConfigFile(".env")
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, err
			}
		}
	}

	v.AutomaticEnv()

	var env Env
	if err := v.Unmarshal(&env); err != nil {
		return nil, err
	}

	if err := env.validate(); err != nil {
		return nil, err
	}

	return &env, nil
}

func (e *Env) validate() error {
	if e.App.Port == 0 {
		return errors.New("APP_PORT is required")
	} else if e.JWT.AccessSecret == "" {
		return errors.New("JWT_ACCESS_SECRET is required")
	} else if e.JWT.AccessExp == 0 {
		return errors.New("JWT_ACCESS_EXP is required")
	} else if e.JWT.RefreshSecret == "" {
		return errors.New("JWT_REFRESH_SECRET is required")
	} else if e.JWT.RefreshExp == 0 {
		return errors.New("JWT_REFRESH_EXP is required")
	} else if e.Database.DBName == "" {
		return errors.New("DB is required")
	} else if e.Database.Host == "" {
		return errors.New("DB_HOST is required")
	} else if e.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	} else if e.Database.Port == 0 {
		return errors.New("DB_PORT is required")
	} else if e.Database.SslMode == "" {
		return errors.New("DB_SSLMODE is required (disable or enable)")
	}
	return nil
}
