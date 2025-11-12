package repositories

import (
	"errors"
	"fmt"
	"log"

	"school-assistant-wh/internal/models"

	"gorm.io/gorm"
)

type UserLinkRepository struct {
	db                 *gorm.DB
	studentProfileRepo *StudentProfileRepository
}

func NewUserLinkRepository(db *gorm.DB, studentProfileRepo *StudentProfileRepository) *UserLinkRepository {
	return &UserLinkRepository{
		db:                 db,
		studentProfileRepo: studentProfileRepo,
	}
}

func (r *UserLinkRepository) GetUserLinks(userID int) ([]models.UserLinkWithStudent, error) {
	var links []models.UserLink
	err := r.db.Table("gk_miniapps.school_link_user").
		Where("UserID = ? AND IsActive = ?", userID, true).
		Find(&links).Error

	if err != nil {
		return nil, err
	}

	result := make([]models.UserLinkWithStudent, 0, len(links))

	for _, link := range links {
		linkWithStudent := models.UserLinkWithStudent{
			UserLink: link,
		}

		if link.SchoolID != "" && link.StudentID != "" {
			student, err := r.studentProfileRepo.GetStudentProfile(link.SchoolID, link.StudentID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Error fetching student profile: %v", err)
			}
			linkWithStudent.Student = student
		}

		result = append(result, linkWithStudent)
	}

	return result, nil
}

func (r *UserLinkRepository) GetPrimaryLink(userID int) (*models.UserLinkWithStudent, error) {
	var result models.UserLinkWithStudent

	err := r.db.Table("gk_miniapps.school_link_user").
		Where("UserID = ? AND IsPrimary = ? AND IsActive = ?", userID, true, true).
		First(&result.UserLink).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	profile, err := r.studentProfileRepo.GetStudentProfile(result.SchoolID, result.StudentID)

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to fetch student profile: %w", err)
	}

	if err == nil {
		result.Student = profile
	}

	return &result, nil
}

func (r *UserLinkRepository) UpdatePrimaryStatus(userID int, studentID string, schoolID string) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Table("gk_miniapps.school_link_user").
		Where("UserID = ?", userID).
		Update("IsPrimary", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("gk_miniapps.school_link_user").
		Where("UserID = ? AND StudentID = ? AND SchoolID = ?", userID, studentID, schoolID).
		Update("IsPrimary", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
