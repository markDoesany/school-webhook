package models

// UserLink represents the relationship between a user and a student in the system
type UserLink struct {
	ID          int    `gorm:"primaryKey;column:ID"`
	CreatedAt   string `gorm:"column:CreatedAt"`
	UpdatedAt   string `gorm:"column:UpdatedAt"`
	IsActive    bool   `gorm:"column:IsActive"`
	UserID      int    `gorm:"column:UserID"`
	StudentID   string `gorm:"column:StudentID"`
	SchoolID    string `gorm:"column:SchoolID"`
	IsNewlyLink bool   `gorm:"column:IsNewlyLink"`
	IsPrimary   bool   `gorm:"column:IsPrimary"`
}
