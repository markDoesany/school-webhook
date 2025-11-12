package repositories

import (
	"fmt"
	"school-assistant-wh/internal/cache"
	"school-assistant-wh/internal/models"
	"time"

	"gorm.io/gorm"
)

type StudentProfileRepository struct {
	db    *gorm.DB
	cache *cache.StudentProfileCache
}

func NewStudentProfileRepository(db *gorm.DB) *StudentProfileRepository {
	// Cache student profiles for 12 hours
	profileCache := cache.NewStudentProfileCache(12 * time.Hour)
	return &StudentProfileRepository{
		db:    db,
		cache: profileCache,
	}
}

func (r *StudentProfileRepository) GetStudentProfile(schoolID, studentID string) (*models.StudentProfile, error) {
	if schoolID == "" || studentID == "" {
		return nil, fmt.Errorf("schoolID and studentID are required")
	}

	// Try to get from cache first
	if cachedProfile, found := r.cache.Get(schoolID, studentID); found {
		return cachedProfile, nil
	}

	// Not in cache, fetch from database
	studentTable := fmt.Sprintf("school_%s_students", schoolID)
	var student models.StudentProfile
	err := r.db.Table(studentTable).
		Select("ID, StudentID, BorrowerID, FirstName, MiddleName, LastName, Course, YearLevel, Status, MobileNumber, EmailAddress, Gender, Birthdate").
		Where("StudentID = ?", studentID).
		First(&student).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching student profile: %v", err)
	}

	// Fetch school information
	var school models.School
	err = r.db.Table("gk_miniapps.school").
		Where("SchoolID = ?", schoolID).
		First(&school).Error

	if err == nil {
		student.School = &school
	}

	// Cache the result
	r.cache.Set(schoolID, studentID, &student)

	return &student, nil
}

func (r *StudentProfileRepository) InvalidateCache(schoolID, studentID string) {
	r.cache.Invalidate(schoolID, studentID)
}

func (r *StudentProfileRepository) ClearCache() {
	r.cache.Clear()
}
