package model

import "time"

type User struct {
	ID        uint64
	Username  string `gorm:"unique"`
	Password  string
	ImgUrl    *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Rooms    []Room    `gorm:"foreignKey:CreatorID;references:ID"`
	Messages []Message `gorm:"foreignKey:SenderID;references:ID"`
}
