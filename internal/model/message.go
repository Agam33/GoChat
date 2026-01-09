package model

import (
	"time"

	"gorm.io/datatypes"
)

type Message struct {
	ID       uint64 `gorm:"primaryKey"`
	RoomID   uint64
	SenderID uint64

	ReplyID *uint64 `gorm:"index"`

	Type    string         `gorm:"check:type IN ('text', 'image', 'system')"`
	Content datatypes.JSON `gorm:"type:jsonb"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Sender *User    `gorm:"foreignKey:SenderID;references:ID"`
	Room   *Room    `gorm:"foreignKey:RoomID;references:ID"`
	Reply  *Message `gorm:"foreignKey:ReplyID;references:ID"`
}
