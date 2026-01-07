package model

import "time"

type Room struct {
	ID        uint64
	CreatorID uint64
	Name      string
	ImgUrl    *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Creator  *User     `gorm:"foreignKey:CreatorID;references:ID"`
	Messages []Message `gorm:"foreignKey:RoomID;references:ID"`
}
