package model

import "time"

type UserRoom struct {
	UserID    uint64 `gorm:"primaryKey;not null"`
	RoomID    uint64 `gorm:"primaryKey;not null"`
	Role      string `gorm:"check:role IN ('admin', 'member')"`
	CreatedAt time.Time

	Room Room `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
