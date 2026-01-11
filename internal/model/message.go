package model

import (
	"time"

	"gorm.io/datatypes"
)

type Message struct {
	ID       uint64 `gorm:"primaryKey;not null"`
	RoomID   uint64 `gorm:"not null"`
	SenderID uint64 `gorm:"not null"`

	ReplyID *uint64 `gorm:"index"`

	ContentType string         `gorm:"check:content_type IN ('text', 'image', 'system')"`
	Content     datatypes.JSON `gorm:"type:jsonb"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Sender *User    `gorm:"foreignKey:SenderID;references:ID"`
	Room   *Room    `gorm:"foreignKey:RoomID;references:ID;constraint:OnDelete:CASCADE"`
	Reply  *Message `gorm:"foreignKey:ReplyID;references:ID"`
}
