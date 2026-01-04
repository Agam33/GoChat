package model

import "time"

type User struct {
	ID        uint64
	Name      string
	Password  string
	ImgUrl    *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Rooms    []Room    `gorm:"foreignKey:CreatorID"`
	Messages []Message `gorm:"foreignKey:SenderID"`
}
