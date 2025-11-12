package models

import (
	"time"
)

type User struct {
	ID          uint       `gorm:"primaryKey;column:ID"`
	CreatedAt   time.Time  `gorm:"column:CreatedAt"`
	UpdatedAt   time.Time  `gorm:"column:UpdatedAt"`
	IsActive    bool       `gorm:"column:IsActive;default:true"`
	Code        *string    `gorm:"column:Code;size:10;unique"`
	PSID        string     `gorm:"column:PSID;size:100;unique;not null"`
	FBName      string     `gorm:"column:FBName;size:100"`
	FBImgURL    *string    `gorm:"column:FBImgURL;type:text"`
	Email       *string    `gorm:"column:Email;size:50"`
	LastLoginAt *time.Time `gorm:"column:LastLoginAt"`
	Notes1      *string    `gorm:"column:Notes1;type:text"`
}

func (User) TableName() string {
	return "school_messenger_users"
}
