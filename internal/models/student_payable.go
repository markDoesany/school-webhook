package models

import (
	"time"
)

// StudentPayable represents the school_students_payables table
// @gorm table:"school_students_payables"
type StudentPayable struct {
	ID                int       `gorm:"primaryKey;column:ID;autoIncrement" json:"id"`
	DateTimeIN        time.Time `gorm:"column:DateTimeIN;not null" json:"date_time_in"`
	DateTimeCompleted time.Time `gorm:"column:DateTimeCompleted;not null" json:"date_time_completed"`
	SOAID             string    `gorm:"column:SOAID;size:100;not null;uniqueIndex" json:"soa_id"`
	BillingID         string    `gorm:"column:BillingID;size:100;not null;index" json:"billing_id"`
	SchoolID          string    `gorm:"column:SchoolID;size:100;not null;index" json:"school_id"`
	SchoolName        string    `gorm:"column:SchoolName;size:100;not null" json:"school_name"`
	MerchantID        string    `gorm:"column:MerchantID;size:100;not null" json:"merchant_id"`
	MerchantName      string    `gorm:"column:MerchantName;size:100;not null" json:"merchant_name"`
	BorrowerID        string    `gorm:"column:BorrowerID;size:100;not null;index" json:"borrower_id"`
	StudentID         string    `gorm:"column:StudentID;size:100;not null;index" json:"student_id"`
	Course            string    `gorm:"column:Course;size:100;not null" json:"course"`
	YearLevel         string    `gorm:"column:YearLevel;size:100;not null" json:"year_level"`
	Semester          string    `gorm:"column:Semester;size:100;not null" json:"semester"`
	ExamTerm          string    `gorm:"column:ExamTerm;size:100;not null" json:"exam_term"`
	Particulars       string    `gorm:"column:Particulars;type:text;not null" json:"particulars"`
	TotalAmountToPay  float64   `gorm:"column:TotalAmountToPay;type:decimal(14,2);not null" json:"total_amount_to_pay"`
	Type              *string   `gorm:"column:Type;size:100;index" json:"type,omitempty"`
	SchoolYear        *string   `gorm:"column:SchoolYear;size:100" json:"school_year,omitempty"`
	Status            string    `gorm:"column:Status;size:100;not null;index" json:"status"`
	Extra1            string    `gorm:"column:Extra1;size:100;not null;default:'.'" json:"extra1"`
	Extra2            string    `gorm:"column:Extra2;size:100;not null;default:'.'" json:"extra2"`
	Extra3            string    `gorm:"column:Extra3;size:100;not null;default:'.'" json:"extra3"`
	Extra4            string    `gorm:"column:Extra4;size:100;not null;default:'.'" json:"extra4"`
	Notes1            *string   `gorm:"column:Notes1;type:text" json:"notes1,omitempty"`
	Notes2            *string   `gorm:"column:Notes2;type:text" json:"notes2,omitempty"`
}

// TableName specifies the table name for the StudentPayable model
func (StudentPayable) TableName() string {
	return "school_students_payables"
}
