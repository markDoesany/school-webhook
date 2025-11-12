package models

import (
	"fmt"
	"time"
)

// Bulletin represents a school bulletin/announcement in the system
type Bulletin struct {
	ID          int       `gorm:"primaryKey;column:ID;autoIncrement" json:"ID"`
	DateTimeIN  time.Time `gorm:"column:DateTimeIN;not null" json:"DateTimeIN"`
	Title       string    `gorm:"column:Title;not null;size:100;index" json:"Title"`
	Description *string   `gorm:"column:Description;type:text" json:"Description,omitempty"`
	ImageURL    string    `gorm:"column:ImageURL;type:text;not null" json:"ImageURL"`
	PeriodStart time.Time `gorm:"column:PeriodStart;not null" json:"PeriodStart"`
	PeriodEnd   time.Time `gorm:"column:PeriodEnd;not null" json:"PeriodEnd"`
	AddedBy     string    `gorm:"column:Addedby;not null;size:100" json:"AddedBy"`
	Status      string    `gorm:"column:Status;not null;size:100" json:"Status"`
	SchoolID    string    `gorm:"-" json:"-"` // Not stored in DB, used for table name generation
	Extra1      string    `gorm:"column:Extra1;not null;default:'';size:100" json:"Extra1"`
	Extra2      string    `gorm:"column:Extra2;not null;default:'';size:100" json:"Extra2"`
	Extra3      string    `gorm:"column:Extra3;not null;default:'';size:100" json:"Extra3"`
	Extra4      string    `gorm:"column:Extra4;not null;default:'';size:100" json:"Extra4"`
	Notes1      *string   `gorm:"column:Notes1;type:text" json:"Notes1,omitempty"`
	Notes2      *string   `gorm:"column:Notes2;type:text" json:"Notes2,omitempty"`
}

func (b Bulletin) TableName() string {
	year := time.Now().Year()
	if b.SchoolID == "" {
		b.SchoolID = "cpeu" // Default school ID if not specified
	}
	return fmt.Sprintf("school_%s_bulletin_%d", b.SchoolID, year)
}

func NewBulletin(schoolID string) *Bulletin {
	return &Bulletin{
		SchoolID: schoolID,
	}
}
