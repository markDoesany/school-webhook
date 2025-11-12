package models

import (
	"time"
)

// SupportThread represents a support thread in the system
type SupportThread struct {
	ID                           int        `gorm:"column:ID;primaryKey;autoIncrement" json:"id"`
	DateTimeIN                   *time.Time `gorm:"column:DateTimeIN" json:"date_time_in,omitempty"`
	DateTimeCompleted            *time.Time `gorm:"column:DateTimeCompleted" json:"date_time_completed,omitempty"`
	ThreadID                     string     `gorm:"column:ThreadID;not null;default:'.';index" json:"thread_id"`
	HelpTopic                    string     `gorm:"column:HelpTopic;not null;default:'.';index" json:"help_topic"`
	Subject                      *string    `gorm:"column:Subject;size:100" json:"subject,omitempty"`
	SupportUserID                *string    `gorm:"column:SupportUserID;size:100" json:"support_user_id,omitempty"`
	MobileNo                     string     `gorm:"column:MobileNo;not null;index:MobileNo;index:Mobile" json:"mobile_no"`
	GKBorrowerID                 string     `gorm:"column:GKBorrowerID;not null;index" json:"borrower_id"`
	GKBorrowerName               string     `gorm:"column:GKBorrowerName;not null;default:'.';index" json:"borrower_name"`
	GKIMEI                       *string    `gorm:"column:GKIMEI;size:100;index" json:"imei,omitempty"`
	GKEmailAddress               *string    `gorm:"column:GKEmailAddress;size:100;index" json:"email,omitempty"`
	SubscriberNotificationStatus *string    `gorm:"column:SubscriberNotificationStatus;size:100;index" json:"subscriber_notification_status,omitempty"` // UnRead-0 Read-1
	SupportNotificationStatus    *string    `gorm:"column:SupportNotificationStatus;size:100;index" json:"support_notification_status,omitempty"`       // READ -1, UNREAD-0
	Status                       *string    `gorm:"column:Status;size:100;index" json:"status,omitempty"`
	Extra1                       *string    `gorm:"column:Extra1;size:100" json:"extra1,omitempty"`
	Extra2                       *string    `gorm:"column:Extra2;size:100" json:"extra2,omitempty"`
	Extra3                       *string    `gorm:"column:Extra3;size:100" json:"extra3,omitempty"`
	Notes1                       *string    `gorm:"column:Notes1;size:100" json:"notes1,omitempty"`
	Notes2                       *string    `gorm:"column:Notes2;size:100" json:"notes2,omitempty"`
}
