package models

import (
	"fmt"
	"time"
)

// PaymentLog represents the school_payment_logs_YYYY table structure
type PaymentLog struct {
	ID                     int       `gorm:"primaryKey;column:ID;autoIncrement" json:"id"`
	DateTimeIN             time.Time `gorm:"column:DateTimeIN;not null" json:"date_time_in"`
	DateTimeCompleted      time.Time `gorm:"column:DateTimeCompleted;not null" json:"date_time_completed"`
	PaymentTxnID           string    `gorm:"column:PaymentTxnID;size:100;not null;uniqueIndex" json:"payment_txn_id"`
	SOAID                  string    `gorm:"column:SOAID;size:100;not null;index" json:"soa_id"`
	BillingID              string    `gorm:"column:BillingID;size:100;not null;index" json:"billing_id"`
	SchoolID               string    `gorm:"column:SchoolID;size:100;not null;index" json:"school_id"`
	MerchantID             string    `gorm:"column:MerchantID;size:100;not null;index" json:"merchant_id"`
	MerchantName           string    `gorm:"column:MerchantName;size:100;not null;index" json:"merchant_name"`
	BorrowerID             string    `gorm:"column:BorrowerID;size:100;not null;index" json:"borrower_id"`
	StudentID              string    `gorm:"column:StudentID;size:100;not null;index" json:"student_id"`
	StudentMobileNumber    string    `gorm:"column:StudentMobileNumber;size:100;not null;index" json:"student_mobile_number"`
	StudentFirstName       string    `gorm:"column:StudentFirstName;size:100;not null" json:"student_first_name"`
	StudentLastName        string    `gorm:"column:StudentLastName;size:100;not null" json:"student_last_name"`
	PaymentDetails         string    `gorm:"column:PaymentDetails;type:text;not null" json:"payment_details"`
	Amount                 float64   `gorm:"column:Amount;type:decimal(14,2);not null" json:"amount"`
	CustomerServiceCharge  float64   `gorm:"column:CustomerServiceCharge;type:decimal(14,2);default:0.00" json:"customer_service_charge"`
	MerchantServiceCharge  float64   `gorm:"column:MerchantServiceCharge;type:decimal(14,2);default:0.00" json:"merchant_service_charge"`
	ResellerDiscount       float64   `gorm:"column:ResellerDiscount;type:decimal(14,2);default:0.00" json:"reseller_discount"`
	TotalAmount            float64   `gorm:"column:TotalAmount;type:decimal(14,2);default:0.00" json:"total_amount"`
	TransactionMedium      string    `gorm:"column:TransactionMedium;size:100;not null" json:"transaction_medium"`
	ProcessID              string    `gorm:"column:ProcessID;size:100;not null;default:'.'" json:"process_id"`
	PaymentType            string    `gorm:"column:PaymentType;size:100;not null;index" json:"payment_type"`
	DateTimePaid           time.Time `gorm:"column:DateTimePaid;not null;index" json:"date_time_paid"`
	PartnerNetworkID       string    `gorm:"column:PartnerNetworkID;size:100;not null" json:"partner_network_id"`
	PartnerNetworkName     string    `gorm:"column:PartnerNetworkName;size:100;not null" json:"partner_network_name"`
	PartnerOutletID        string    `gorm:"column:PartnerOutletID;size:100;not null" json:"partner_outlet_id"`
	PartnerOutletName      string    `gorm:"column:PartnerOutletName;size:100;not null" json:"partner_outlet_name"`
	PreConsummationSession string    `gorm:"column:PreConsummationSession;size:100;not null;default:'.'" json:"pre_consummation_session"`
	Status                 string    `gorm:"column:Status;size:100;not null;index" json:"status"`
	Extra1                 string    `gorm:"column:Extra1;size:100;not null;default:'.'" json:"extra1"`
	Extra2                 string    `gorm:"column:Extra2;size:100;not null;default:'.'" json:"extra2"`
	Extra3                 string    `gorm:"column:Extra3;size:100;not null;default:'.'" json:"extra3"`
	Extra4                 string    `gorm:"column:Extra4;size:100;not null;default:'.'" json:"extra4"`
	Notes1                 *string   `gorm:"column:Notes1;type:text" json:"notes1,omitempty"`
	Notes2                 *string   `gorm:"column:Notes2;type:text" json:"notes2,omitempty"`
}

// TableName returns the dynamic table name based on the year
func (PaymentLog) TableName(year int) string {
	return fmt.Sprintf("school_payment_logs_%d", year)
}
