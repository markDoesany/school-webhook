package models

import "time"

type School struct {
	ID                     int        `gorm:"column:ID;primaryKey" json:"id"`
	SchoolID               string     `gorm:"column:SchoolID;unique;not null" json:"school_id"`
	SchoolName             string     `gorm:"column:SchoolName;not null" json:"school_name"`
	StreetAddress          string     `gorm:"column:StreetAddress;not null" json:"street_address"`
	City                   string     `gorm:"column:City;not null" json:"city"`
	Province               string     `gorm:"column:Province;not null" json:"province"`
	Country                string     `gorm:"column:Country;not null" json:"country"`
	AuthorizedPersonel     string     `gorm:"column:AuthorizedPersonel;not null" json:"authorized_personnel"`
	AuthorizedMobileNumber string     `gorm:"column:AuthorizedMobileNumber;not null" json:"authorized_mobile_number"`
	AuthorizedEmailAddress string     `gorm:"column:AuthorizedEmailAddress;not null" json:"authorized_email"`
	SchoolLogo             string     `gorm:"column:SchoolLogo;type:text;not null" json:"school_logo"`
	WithRFID               int        `gorm:"column:withRFID;not null;default:0" json:"with_rfid"`
	XMLDetails             *string    `gorm:"column:XMLDetails;type:text" json:"xml_details,omitempty"`
	DateTimeAdded          *time.Time `gorm:"column:DateTimeAdded" json:"date_time_added,omitempty"`
	AddedBy                string     `gorm:"column:AddedBy;not null" json:"added_by"`
	LastDateTimeUpdated    *time.Time `gorm:"column:LastDateTimeUpdated" json:"last_updated,omitempty"`
	Status                 string     `gorm:"column:Status;not null;index" json:"status"`
	Extra1                 string     `gorm:"column:Extra1;default:'.'" json:"extra1,omitempty"`
	Extra2                 string     `gorm:"column:Extra2;default:'.'" json:"extra2,omitempty"`
	Extra3                 string     `gorm:"column:Extra3;default:'.'" json:"extra3,omitempty"`
	Extra4                 string     `gorm:"column:Extra4;default:'.'" json:"extra4,omitempty"`
	Notes1                 *string    `gorm:"column:Notes1;type:text" json:"notes1,omitempty"`
	Notes2                 *string    `gorm:"column:Notes2;type:text" json:"notes2,omitempty"`
}

func (School) TableName() string {
	return "school"
}
