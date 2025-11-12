package repositories

import (
	"fmt"
	"time"

	"school-assistant-wh/internal/models"

	"gorm.io/gorm"
)

type BulletinRepository struct {
	db *gorm.DB
}

func NewBulletinRepository(db *gorm.DB) *BulletinRepository {
	return &BulletinRepository{
		db: db,
	}
}

func (r *BulletinRepository) GetBulletins(schoolID string, year *int, offset, limit int) ([]models.Bulletin, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school ID is required")
	}

	currentYear := time.Now().Year()
	if year == nil {
		year = &currentYear
	}

	bulletin := models.NewBulletin(schoolID)

	query := r.db.Table(bulletin.TableName()).
		Where("Status = ?", "active").
		Order("PeriodStart DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}

	var bulletins []models.Bulletin
	if err := query.Find(&bulletins).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch bulletins: %w", err)
	}

	return bulletins, nil
}

func (r *BulletinRepository) GetBulletinsCount(schoolID string) (int64, error) {
	if schoolID == "" {
		return 0, fmt.Errorf("school ID is required")
	}

	bulletin := models.NewBulletin(schoolID)

	var count int64
	err := r.db.Table(bulletin.TableName()).
		Where("Status = ?", "active").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count bulletins: %w", err)
	}

	return count, nil
}

func (r *BulletinRepository) GetBulletinByID(schoolID string, id int) (*models.Bulletin, error) {
	if schoolID == "" {
		return nil, fmt.Errorf("school ID is required")
	}

	bulletin := models.NewBulletin(schoolID)
	bulletin.ID = id

	if err := r.db.Table(bulletin.TableName()).First(bulletin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch bulletin: %w", err)
	}

	return bulletin, nil
}
