package model

import (
	"time"

	"gorm.io/datatypes"
)

type ReplyMessage struct {
	ID           uint64         `gorm:"primaryKey"`
	ContentType  string         `gorm:"check:content_type IN ('text', 'image', 'system')"`
	ReplyContent datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
