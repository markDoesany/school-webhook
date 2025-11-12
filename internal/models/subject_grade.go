package models

import (
	"fmt"
	"time"
)

// SubjectGrade represents a student's grade for a specific subject
type SubjectGrade struct {
	ID                 int       `gorm:"primaryKey;column:ID;autoIncrement" json:"ID"`
	DateTimeIN         time.Time `gorm:"column:DateTimeIN;not null" json:"DateTimeIN"`
	StudentID          string    `gorm:"column:StudentID;not null;index" json:"StudentID"`
	Course             string    `gorm:"column:Course;not null;index" json:"Course"`
	YearLevel          string    `gorm:"column:YearLevel;not null;index" json:"YearLevel"`
	Semester           string    `gorm:"column:Semester;not null;index" json:"Semester"`
	ExamTerm           string    `gorm:"column:ExamTerm;not null" json:"ExamTerm"`
	SubjectID          string    `gorm:"column:SubjectID;not null" json:"SubjectID"`
	SubjectDescription string    `gorm:"column:SubjectDescription;not null;index" json:"SubjectDescription"`
	SubjectSchedule    string    `gorm:"column:SubjectSchedule;not null" json:"SubjectSchedule"`
	SubjectRoom        string    `gorm:"column:SubjectRoom;not null" json:"SubjectRoom"`
	SubjectUnit        string    `gorm:"column:SubjectUnit;not null" json:"SubjectUnit"`
	Type               string    `gorm:"column:Type;not null" json:"Type"`
	SchoolYear         string    `gorm:"column:SchoolYear;not null" json:"SchoolYear"`
	StudentGrade       string    `gorm:"column:StudentGrade;not null" json:"StudentGrade"`
	Status             string    `gorm:"column:Status;not null;index" json:"Status"`
	Extra1             string    `gorm:"column:Extra1;not null;default:''" json:"Extra1"`
	Extra2             string    `gorm:"column:Extra2;not null;default:''" json:"Extra2"`
	Extra3             string    `gorm:"column:Extra3;not null;default:''" json:"Extra3"`
	Extra4             string    `gorm:"column:Extra4;not null;default:''" json:"Extra4"`
	Notes1             *string   `gorm:"column:Notes1" json:"Notes1,omitempty"`
	Notes2             *string   `gorm:"column:Notes2" json:"Notes2,omitempty"`
	SchoolID           string    `gorm:"-" json:"-"`
}

func (s SubjectGrade) TableName() string {
	if s.SchoolID == "" {
		s.SchoolID = "cpeu" // Default school ID if not specified
	}
	return fmt.Sprintf("school_%s_students_subject_grades", s.SchoolID)
}

func NewSubjectGrade(schoolID string) *SubjectGrade {
	return &SubjectGrade{
		SchoolID: schoolID,
	}
}
