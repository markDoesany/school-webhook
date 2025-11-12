package models

// StudentProfile represents the student information in the system
type StudentProfile struct {
	ID           int     `json:"ID"`
	StudentID    string  `json:"StudentID"`
	BorrowerID   string  `json:"BorrowerID"`
	FirstName    string  `json:"FirstName"`
	MiddleName   string  `json:"MiddleName"`
	LastName     string  `json:"LastName"`
	Course       string  `json:"Course"`
	YearLevel    string  `json:"YearLevel"`
	Status       string  `json:"Status"`
	MobileNumber string  `json:"MobileNumber"`
	EmailAddress string  `json:"EmailAddress"`
	Gender       string  `json:"Gender"`
	Birthdate    string  `json:"Birthdate"`
	School       *School `json:"School,omitempty" gorm:"-"` // School details will be populated separately
}
