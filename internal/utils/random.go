package utils

import (
	"school-assistant-wh/internal/models"
	"time"

	"crypto/rand"
	"fmt"

	"gorm.io/gorm"
)

func GenerateUniqueCode(db *gorm.DB) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	for i := 0; i < 10; i++ {
		b := make([]byte, codeLength)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}

		for i := range b {
			b[i] = charset[int(b[i])%len(charset)]
		}

		code := "SA-" + string(b)

		var count int64
		if err := db.Model(&models.User{}).Where("Code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}

// Helper function to create string pointers
func StringPtr(s string) *string {
	return &s
}

func GenerateThreadID() string {
	now := time.Now()
	timestamp := now.Unix()
	seq := now.Nanosecond() % 10000
	if seq == 0 {
		seq = 1
	}

	return fmt.Sprintf("%d%04d", timestamp, seq)
}
