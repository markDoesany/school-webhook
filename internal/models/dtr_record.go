package models

import (
	"fmt"
	"strings"
	"time"
)

// DTRRecord represents the school_cpeu_dtr_records_YYYY_MM table structure
type DTRRecord struct {
	ID             int       `gorm:"primaryKey;column:ID;autoIncrement" json:"id"`
	SchoolID       string    `gorm:"column:SchoolID;size:100;not null;index" json:"school_id"`
	StudentID      string    `gorm:"column:StudentID;size:100;not null;index" json:"student_id"`
	StudentName    string    `gorm:"column:StudentName;size:255;not null" json:"student_name"`
	RFIDCardNumber string    `gorm:"column:RFIDCardNumber;size:255;not null" json:"rfid_card_number"`
	Type           string    `gorm:"column:Type;size:100;not null" json:"type"`
	DateTimeIN     time.Time `gorm:"column:DateTimeIN;not null;index" json:"date_time_in"`
	IPAddress      string    `gorm:"column:IPAddress;size:100;not null;index" json:"ip_address"`
	BatchID        string    `gorm:"column:BatchID;size:100;not null;index" json:"batch_id"`
	Extra1         string    `gorm:"column:Extra1;size:100;not null;default:'."' json:"extra1"`
	Extra2         string    `gorm:"column:Extra2;size:100;not null;default:'."' json:"extra2"`
	Extra3         string    `gorm:"column:Extra3;size:100;not null;default:'."' json:"extra3"`
	Notes1         *string   `gorm:"column:Notes1;type:text" json:"notes1,omitempty"`
	Notes2         *string   `gorm:"column:Notes2;type:text" json:"notes2,omitempty"`
}

// TableName returns the dynamic table name based on school ID, year, and month
func (r DTRRecord) TableName(year int, month time.Month) string {
	return fmt.Sprintf("school_%s_dtr_records_%d_%02d", strings.ToLower(r.SchoolID), year, month)
}
