package model

import "time"

type Room struct {
	ID        uint64
	CreatorID uint64
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time

	Messages []Message `gorm:"foreignKey:RoomID"`
}
