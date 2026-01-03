package jwt

import (
	"time"
)

type JwtConfig struct {
	AccessExpire  time.Duration
	RefreshExpire time.Duration
	AccessSecret  string
	RefreshSecret string
}
