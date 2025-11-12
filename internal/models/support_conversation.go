package models

import (
	"time"
)

// SupportConversation represents a single message in a support conversation
type SupportConversation struct {
	ID                 int       `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	DateTimeIN         time.Time `gorm:"column:DateTimeIN;not null;index" json:"date_time_in"`
	ThreadID           string    `gorm:"column:ThreadID;not null;index" json:"thread_id"`
	ReplySupportUserID string    `gorm:"column:ReplySupportUserID;not null;index" json:"reply_support_user_id"`
	ReplySupportName   string    `gorm:"column:ReplySupportName;not null;index" json:"reply_support_name"`
	ThreadType         string    `gorm:"column:ThreadType;not null;index" json:"thread_type"`
	Message            string    `gorm:"column:Message;type:text;not null" json:"message"`
	Extra1             *string   `gorm:"column:Extra1;size:100" json:"extra1,omitempty"`
	Extra2             *string   `gorm:"column:Extra2;size:100" json:"extra2,omitempty"`
	Extra3             *string   `gorm:"column:Extra3;size:100" json:"extra3,omitempty"`
	Notes1             *string   `gorm:"column:Notes1;type:text" json:"notes1,omitempty"`
	Notes2             *string   `gorm:"column:Notes2;type:text" json:"notes2,omitempty"`
}
