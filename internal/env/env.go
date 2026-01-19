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
	} `mapstructure:",squash"`

	JWT struct {
		AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
		RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
		AccessExp     time.Duration `mapstructure:"JWT_ACCESS_EXP"`
		RefreshExp    time.Duration `mapstructure:"JWT_REFRESH_EXP"`
	} `mapstructure:",squash"`

	Database struct {
		Host     string `mapstructure:"DB_HOST"`
		User     string `mapstructure:"DB_USER"`
		Port     int    `mapstructure:"DB_PORT"`
		DBName   string `mapstructure:"DB_NAME"`
		Password string `mapstructure:"DB_PASSWORD"`
		SslMode  string `mapstructure:"DB_SSLMODE"`
	} `mapstructure:",squash"`

	MQ struct {
		User     string `mapstructure:"MQ_USER"`
		Password string `mapstructure:"MQ_PASSWORD"`
		Port     int    `mapstructure:"MQ_PORT"`
		Host     string `mapstructure:"MQ_HOST"`
		Vhost    string `mapstructure:"MQ_VHOST"`
	} `mapstructure:",squash"`
}

func NewEnv() (*Env, error) {
	v := viper.New()

	if os.Getenv("APP_ENV") != "prod" {
		v.SetConfigName(".env")
		v.SetConfigType("env")
		v.AddConfigPath(".")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, err
			}
		}
	}

	keys := getKeyBind()
	for _, key := range keys {
		_ = v.BindEnv(key)
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

func getKeyBind() []string {
	return []string{
		"APP_PORT",
		"JWT_ACCESS_SECRET", "JWT_ACCESS_EXP", "JWT_REFRESH_SECRET", "JWT_REFRESH_EXP",
		"DB_NAME", "DB_USER", "DB_PORT", "DB_HOST", "DB_SSLMODE", "DB_PASSWORD",
		"MQ_USER", "MQ_PASSWORD", "MQ_PORT", "MQ_VHOST",
	}
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
		return errors.New("DB_NAME is required")
	} else if e.Database.Host == "" {
		return errors.New("DB_HOST is required")
	} else if e.Database.User == "" {
		return errors.New("DB_USER is required")
	} else if e.Database.Password == "" {
		return errors.New("DB_PASSWORD is required")
	} else if e.Database.Port == 0 {
		return errors.New("DB_PORT is required")
	} else if e.Database.SslMode == "" {
		return errors.New("DB_SSLMODE is required (disable or enable)")
	} else if e.MQ.User == "" {
		return errors.New("MQ_USER is required")
	} else if e.MQ.Password == "" {
		return errors.New("MQ_PASSWORD is required")
	} else if e.MQ.Vhost == "" {
		return errors.New("MQ_VHOST is required")
	} else if e.MQ.Port == 0 {
		return errors.New("MQ_PORT is required")
	} else {
		return nil
	}
}
