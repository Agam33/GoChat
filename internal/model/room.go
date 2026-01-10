package model

import "time"

type Room struct {
	ID        uint64 `gorm:"primaryKey;not null"`
	CreatorID uint64 `gorm:"not null"`
	Name      string
	ImgUrl    *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Creator  *User     `gorm:"foreignKey:CreatorID;references:ID"`
	Messages []Message `gorm:"foreignKey:RoomID;references:ID"`
}
