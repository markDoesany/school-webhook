package repositories

import (
	"fmt"
	"school-assistant-wh/internal/models"

	"gorm.io/gorm"
)

type PaymentLogRepository struct {
	db *gorm.DB
}

func NewPaymentLogRepository(db *gorm.DB) *PaymentLogRepository {
	return &PaymentLogRepository{
		db: db,
	}
}

// GetPaymentLogsByStudentID retrieves payment logs for a specific student in a given year
func (r *PaymentLogRepository) GetPaymentLogsByStudentAndSchoolID(year int, studentID string, schoolID string) ([]models.PaymentLog, error) {
	if studentID == "" {
		return nil, fmt.Errorf("student ID cannot be empty")
	}

	var logs []models.PaymentLog
	err := r.db.Table(models.PaymentLog{}.TableName(year)).
		Where("StudentID = ? AND SchoolID = ?", studentID, schoolID).
		Order("DateTimePaid DESC").
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment logs: %w", err)
	}

	return logs, nil
}

// GetPaymentLogsBySOAID retrieves payment logs for a specific SOA ID in a given year
func (r *PaymentLogRepository) GetPaymentLogsBySOAID(year int, soaID string) ([]models.PaymentLog, error) {
	if soaID == "" {
		return nil, fmt.Errorf("SOA ID cannot be empty")
	}

	var logs []models.PaymentLog
	err := r.db.Table(models.PaymentLog{}.TableName(year)).
		Where("SOAID = ?", soaID).
		Order("DateTimePaid DESC").
		Find(&logs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment logs: %w", err)
	}

	return logs, nil
}

// GetPaymentLogByTxnID retrieves a specific payment log by its transaction ID in a given year
func (r *PaymentLogRepository) GetPaymentLogByTxnID(year int, txnID string) (*models.PaymentLog, error) {
	if txnID == "" {
		return nil, fmt.Errorf("transaction ID cannot be empty")
	}

	var log models.PaymentLog
	err := r.db.Table(models.PaymentLog{}.TableName(year)).
		Where("PaymentTxnID = ?", txnID).
		First(&log).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch payment log: %w", err)
	}

	return &log, nil
}

// CreatePaymentLog creates a new payment log entry in the specified year's table
func (r *PaymentLogRepository) CreatePaymentLog(year int, log *models.PaymentLog) error {
	if log == nil {
		return fmt.Errorf("payment log cannot be nil")
	}

	err := r.db.Table(models.PaymentLog{}.TableName(year)).Create(log).Error
	if err != nil {
		return fmt.Errorf("failed to create payment log: %w", err)
	}

	return nil
}

// UpdatePaymentLog updates an existing payment log in the specified year's table
func (r *PaymentLogRepository) UpdatePaymentLog(year int, log *models.PaymentLog) error {
	if log == nil || log.PaymentTxnID == "" {
		return fmt.Errorf("invalid payment log or transaction ID")
	}

	err := r.db.Table(models.PaymentLog{}.TableName(year)).Save(log).Error
	if err != nil {
		return fmt.Errorf("failed to update payment log: %w", err)
	}

	return nil
}
