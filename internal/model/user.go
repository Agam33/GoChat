package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primaryKey;not null"`
	Username  string `gorm:"unique"`
	Password  string
	ImgUrl    *string
	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`

	Rooms    []Room    `gorm:"foreignKey:CreatorID;references:ID"`
	Messages []Message `gorm:"foreignKey:SenderID;references:ID"`
}
