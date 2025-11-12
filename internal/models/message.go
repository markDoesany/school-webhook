package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SenderID string `json:"sender_id" gorm:"index"`
	Text     string `json:"text"`
}
