package repositories

import (
	"fmt"
	"school-assistant-wh/internal/models"

	"gorm.io/gorm"
)

type StudentPayableRepository struct {
	db *gorm.DB
}

func NewStudentPayableRepository(db *gorm.DB) *StudentPayableRepository {
	return &StudentPayableRepository{
		db: db,
	}
}

// GetStudentPayables retrieves all payables for a specific student in a school
func (r *StudentPayableRepository) GetStudentPayables(schoolID, studentID string) ([]models.StudentPayable, error) {
	if schoolID == "" || studentID == "" {
		return nil, fmt.Errorf("invalid input")
	}

	var payables []models.StudentPayable
	result := r.db.Table("school_students_payables").
		Where("SchoolID = ? AND StudentID = ?", schoolID, studentID).
		Order("DateTimeIN DESC").
		Find(&payables)

	if result.Error != nil {
		return nil, result.Error
	}

	return payables, nil
}

// GetActiveStudentPayables retrieves active payables for a specific student in a school
func (r *StudentPayableRepository) GetActiveStudentPayables(schoolID, studentID string) ([]models.StudentPayable, error) {
	if schoolID == "" || studentID == "" {
		return nil, fmt.Errorf("invalid input")
	}

	var payables []models.StudentPayable
	result := r.db.Table("school_students_payables").
		Where("SchoolID = ? AND StudentID = ? AND Status = 'Active'", schoolID, studentID).
		Order("DateTimeIN DESC").
		Find(&payables)

	if result.Error != nil {
		return nil, result.Error
	}

	return payables, nil
}

// GetPayableBySOAID retrieves a specific payable by its SOA ID
func (r *StudentPayableRepository) GetPayableBySOAID(soaID string) (*models.StudentPayable, error) {
	if soaID == "" {
		return nil, fmt.Errorf("invalid input")
	}

	var payable models.StudentPayable
	result := r.db.Table("school_students_payables").
		Where("SOAID = ?", soaID).
		First(&payable)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &payable, nil
}
