package model

import (
	"time"

	"gorm.io/datatypes"
)

type Message struct {
	ID       uint64 `gorm:"primaryKey"`
	RoomID   uint64
	SenderID uint64

	ReplyID      *uint64         `gorm:"index"`
	ReplyContent *datatypes.JSON `gorm:"type:jsonb"`

	Type      string
	Content   datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
