package repositories

import (
	"fmt"
	"school-assistant-wh/internal/models"

	"gorm.io/gorm"
)

type GradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) *GradeRepository {
	return &GradeRepository{
		db: db,
	}
}

// GetStudentSubjectGrades retrieves all subject grades for a specific student
func (r *GradeRepository) GetStudentSubjectGrades(studentID, schoolID string) ([]models.SubjectGrade, error) {
	if studentID == "" || schoolID == "" {
		return nil, fmt.Errorf("studentID and schoolID are required")
	}

	table := fmt.Sprintf("school_%s_students_subject_grades", schoolID)

	var grades []models.SubjectGrade
	err := r.db.Table(table).
		Where("StudentID = ?", studentID).
		Order("Semester, SubjectDescription").
		Find(&grades).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching student subject grades: %v", err)
	}

	return grades, nil
}

// GetStudentGradesBySemester retrieves subject grades for a specific student and semester
func (r *GradeRepository) GetStudentGradesBySemester(studentID, schoolID, semester string) ([]models.SubjectGrade, error) {
	if studentID == "" || schoolID == "" || semester == "" {
		return nil, fmt.Errorf("studentID, schoolID, and semester are required")
	}

	table := fmt.Sprintf("school_%s_students_subject_grades", schoolID)

	var grades []models.SubjectGrade
	err := r.db.Table(table).
		Where("StudentID = ? AND Semester = ?", studentID, semester).
		Order("SubjectDescription").
		Find(&grades).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching student grades for semester %s: %v", semester, err)
	}

	return grades, nil
}

// GetEnrolledSubjects retrieves all enrolled subjects for a student
func (r *GradeRepository) GetEnrolledSubjects(studentID, schoolID string) ([]models.SubjectGrade, error) {
	if studentID == "" || schoolID == "" {
		return nil, fmt.Errorf("studentID and schoolID are required")
	}

	table := fmt.Sprintf("school_%s_students_subject_grades", schoolID)

	var subjects []models.SubjectGrade
	err := r.db.Table(table).
		Select("DISTINCT SubjectDescription, SubjectUnit, SubjectSchedule, SubjectRoom, SchoolYear").
		Where("StudentID = ?", studentID).
		Order("SubjectDescription").
		Find(&subjects).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching enrolled subjects: %v", err)
	}

	return subjects, nil
}

// GetEnrolledSubjectsByYear retrieves all enrolled subjects for a student in a specific school year
func (r *GradeRepository) GetEnrolledSubjectsByYear(studentID, schoolID, schoolYear string) ([]models.SubjectGrade, error) {
	if studentID == "" || schoolID == "" || schoolYear == "" {
		return nil, fmt.Errorf("studentID, schoolID, and schoolYear are required")
	}

	table := fmt.Sprintf("school_%s_students_subject_grades", schoolID)

	var subjects []models.SubjectGrade
	err := r.db.Table(table).
		Select("DISTINCT SubjectDescription, SubjectUnit, SubjectSchedule, SubjectRoom, SchoolYear").
		Where("StudentID = ? AND SchoolYear = ?", studentID, schoolYear).
		Order("SubjectDescription").
		Find(&subjects).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching enrolled subjects for school year %s: %v", schoolYear, err)
	}

	return subjects, nil
}
