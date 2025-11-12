package repositories

import (
	"fmt"
	"school-assistant-wh/internal/models"
	"time"

	"gorm.io/gorm"
)

type DTRRepository struct {
	db *gorm.DB
}

func NewDTRRepository(db *gorm.DB) *DTRRepository {
	return &DTRRepository{
		db: db,
	}
}

// GetDTRRecords retrieves paginated DTR records for a specific student and school
func (r *DTRRepository) GetDTRRecords(year int, month time.Month, schoolID, studentID string, offset, limit int) ([]models.DTRRecord, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school ID cannot be empty")
	}
	if studentID == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}

	record := models.DTRRecord{SchoolID: schoolID}
	var records []models.DTRRecord

	query := r.db.Table(record.TableName(year, month)).
		Where("StudentID = ?", studentID).
		Order("DateTimeIN DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DTR records: %w", err)
	}

	return records, nil
}

// GetDTRRecordsCount returns the total count of DTR records for a student
func (r *DTRRepository) GetDTRRecordsCount(year int, month time.Month, schoolID, studentID string) (int64, error) {
	if schoolID == "" {
		return 0, fmt.Errorf("school ID cannot be empty")
	}
	if studentID == "" {
		return 0, fmt.Errorf("student ID cannot be empty")
	}

	record := models.DTRRecord{SchoolID: schoolID}
	var count int64

	err := r.db.Table(record.TableName(year, month)).
		Where("StudentID = ?", studentID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count DTR records: %w", err)
	}

	return count, nil
}

// GetDTRRecordsByDateRange retrieves DTR records for a specific student and school within a date range
func (r *DTRRepository) GetDTRRecordsByDateRange(year int, month time.Month, schoolID, studentID string, startDate, endDate time.Time, offset, limit int) ([]models.DTRRecord, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school ID cannot be empty")
	}
	if studentID == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}

	record := models.DTRRecord{SchoolID: schoolID}
	var records []models.DTRRecord

	query := r.db.Table(record.TableName(year, month)).
		Where("StudentID = ? AND DateTimeIN BETWEEN ? AND ?",
			studentID, startDate, endDate).
		Order("DateTimeIN ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DTR records by date range: %w", err)
	}

	return records, nil
}
