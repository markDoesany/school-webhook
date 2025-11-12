package repositories

import (
	"fmt"
	"log"
	"strings"
	"time"

	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/utils"

	"gorm.io/gorm"
)

type SupportRepository struct {
	db *gorm.DB
}

func NewSupportRepository(db *gorm.DB) *SupportRepository {
	return &SupportRepository{db: db}
}

// CreateThread creates a new support thread and returns the thread ID
func (r *SupportRepository) CreateThread(thread *models.SupportThread, schoolID string) (string, error) {
	year := time.Now().Year()
	tableName := fmt.Sprintf("gk_support.school_%s_support_thread_%d", strings.ToLower(schoolID), year)

	var tableExists bool
	err := r.db.Raw(
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = 'gk_support' AND table_name = ?",
		tableName,
	).Scan(&tableExists).Error

	if err != nil {
		return "", fmt.Errorf("failed to check if table exists: %w", err)
	}

	if !tableExists {
		createTableSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				ID bigint AUTO_INCREMENT,
				DateTimeIN datetime(3) NULL,
				DateTimeCompleted datetime(3) NULL,
				ThreadID varchar(191) NOT NULL DEFAULT '.',
				HelpTopic varchar(191) NOT NULL DEFAULT '.',
				Subject varchar(100),
				SupportUserID varchar(100),
				MobileNo varchar(191) NOT NULL,
				GKBorrowerID varchar(191) NOT NULL,
				GKBorrowerName varchar(191) NOT NULL DEFAULT '.',
				GKIMEI varchar(100),
				GKEmailAddress varchar(100),
				SubscriberNotificationStatus varchar(100),
				SupportNotificationStatus varchar(100),
				Status varchar(100),
				Extra1 varchar(100),
				Extra2 varchar(100),
				Extra3 varchar(100),
				Notes1 text,
				Notes2 text,
				PRIMARY KEY (ID),
				INDEX idx_thread_id (ThreadID),
				INDEX idx_help_topic (HelpTopic),
				INDEX idx_mobile (MobileNo)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`, tableName)

		if err := r.db.Exec(createTableSQL).Error; err != nil {
			return "", fmt.Errorf("failed to create support thread table: %w", err)
		}
	}

	if thread.ThreadID == "" {
		thread.ThreadID = utils.GenerateThreadID()
	}

	now := time.Now()
	thread.DateTimeIN = &now

	// Use the full table name with schema for the insert
	if err := r.db.Table(tableName).Create(thread).Error; err != nil {
		return "", fmt.Errorf("failed to create support thread: %w", err)
	}

	return thread.ThreadID, nil
}

// CreateMessage adds a new message to an existing support thread
func (r *SupportRepository) CreateMessage(message *models.SupportConversation, schoolID string) error {
	if message.ThreadID == "" {
		return fmt.Errorf("thread ID is required")
	}

	year := time.Now().Year()
	tableName := fmt.Sprintf("gk_support.school_%s_support_conversation_%d", strings.ToLower(schoolID), year)

	var tableExists bool
	err := r.db.Raw(
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = 'gk_support' AND table_name = ?",
		tableName,
	).Scan(&tableExists).Error

	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}

	if !tableExists {
		createTableSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				ID bigint AUTO_INCREMENT,
				DateTimeIN datetime(3) NOT NULL,
				ThreadID varchar(191) NOT NULL,
				ReplySupportUserID varchar(100) NOT NULL,
				ReplySupportName varchar(191) NOT NULL,
				ThreadType varchar(50) NOT NULL,
				Message text NOT NULL,
				Extra1 varchar(100),
				Extra2 varchar(100),
				Extra3 varchar(100),
				Notes1 text,
				Notes2 text,
				PRIMARY KEY (ID),
				INDEX idx_thread_id (ThreadID),
				INDEX idx_datetime (DateTimeIN)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
		`, tableName)

		if err := r.db.Exec(createTableSQL).Error; err != nil {
			return fmt.Errorf("failed to create support conversation table: %w", err)
		}
	}

	message.DateTimeIN = time.Now()

	if err := r.db.Table(tableName).Create(message).Error; err != nil {
		return fmt.Errorf("failed to create support message: %w", err)
	}

	return nil
}

func (r *SupportRepository) GetThreadsByBorrowerID(borrowerID, schoolID string) ([]*models.SupportThread, error) {
	if borrowerID == "" || schoolID == "" {
		return nil, fmt.Errorf("borrower ID and school ID are required")
	}

	year := time.Now().Year()
	tableName := fmt.Sprintf("school_%s_support_thread_%d", strings.ToLower(schoolID), year)
	fullTableName := fmt.Sprintf("`gk_support`.`%s`", tableName)

	// Check if table exists in gk_support database
	var tableExists bool
	err := r.db.Raw(
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = 'gk_support' AND table_name = ?",
		tableName,
	).Scan(&tableExists).Error

	if err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}

	if !tableExists {
		log.Printf("Table %s does not exist in gk_support database", tableName)
		return []*models.SupportThread{}, nil
	}

	var threads []*models.SupportThread

	// First, try exact match
	sql := fmt.Sprintf("SELECT * FROM %s WHERE GKBorrowerID = ? ORDER BY DateTimeIN DESC", fullTableName)
	err = r.db.Raw(sql, borrowerID).Scan(&threads).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get support threads: %w", err)
	}

	// If no results, try trimming any whitespace from borrowerID
	if len(threads) == 0 {
		trimmedBorrowerID := strings.TrimSpace(borrowerID)
		if trimmedBorrowerID != borrowerID {
			log.Printf("No results with exact match, trying with trimmed borrower ID: %s", trimmedBorrowerID)
			err = r.db.Raw(sql, trimmedBorrowerID).Scan(&threads).Error

			if err != nil {
				return nil, fmt.Errorf("failed to get support threads with trimmed borrower ID: %w", err)
			}
		}
	}

	log.Printf("Found %d threads for borrower ID %s", len(threads), borrowerID)
	return threads, nil
}

func (r *SupportRepository) GetThread(threadID, schoolID string) (*models.SupportThread, error) {
	if threadID == "" || schoolID == "" {
		return nil, fmt.Errorf("thread ID and school ID are required")
	}

	year := time.Now().Year()
	tableName := fmt.Sprintf("school_%s_support_thread_%d", strings.ToLower(schoolID), year)
	fullTableName := fmt.Sprintf("`gk_support`.`%s`", tableName)

	var tableExists bool
	err := r.db.Raw(
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = 'gk_support' AND table_name = ?",
		tableName,
	).Scan(&tableExists).Error

	if err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}

	if !tableExists {
		return nil, fmt.Errorf("support thread table not found in gk_support database")
	}

	var thread models.SupportThread
	sql := fmt.Sprintf("SELECT * FROM %s WHERE ThreadID = ? LIMIT 1", fullTableName)
	err = r.db.Raw(sql, threadID).Scan(&thread).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("support thread not found")
		}
		return nil, fmt.Errorf("failed to get support thread: %w", err)
	}

	return &thread, nil
}

// GetMessages retrieves all messages for a specific support thread
func (r *SupportRepository) GetMessages(threadID, schoolID string) ([]*models.SupportConversation, error) {
	if threadID == "" || schoolID == "" {
		return nil, fmt.Errorf("thread ID and school ID are required")
	}

	year := time.Now().Year()
	tableName := fmt.Sprintf("school_%s_support_conversation_%d", strings.ToLower(schoolID), year)
	fullTableName := fmt.Sprintf("`gk_support`.`%s`", tableName)

	// Check if table exists in gk_support database
	var tableExists bool
	err := r.db.Raw(
		"SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = 'gk_support' AND table_name = ?",
		tableName,
	).Scan(&tableExists).Error

	if err != nil {
		return nil, fmt.Errorf("failed to check if table exists: %w", err)
	}

	if !tableExists {
		log.Printf("Table %s does not exist in gk_support database", tableName)
		return []*models.SupportConversation{}, nil
	}

	// Use raw SQL to query messages
	var messages []*models.SupportConversation
	sql := fmt.Sprintf("SELECT * FROM %s WHERE ThreadID = ? ORDER BY DateTimeIN ASC", fullTableName)
	err = r.db.Raw(sql, threadID).Scan(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get support messages: %w", err)
	}

	return messages, nil
}
