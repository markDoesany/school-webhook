package repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"school-assistant-wh/internal/cache"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	fbSvc *facebook.Service
	cache *cache.UserCache
}

func NewUserRepository(db *gorm.DB, fbSvc *facebook.Service) *UserRepository {
	// Cache users for 1 hour by default
	userCache := cache.NewUserCache(time.Hour)
	return &UserRepository{
		db:    db,
		fbSvc: fbSvc,
		cache: userCache,
	}
}

func (r *UserRepository) GetUserByPSID(psid string) (*models.User, error) {
	// Check cache first
	if cachedUser, found := r.cache.GetUser(psid); found {
		return cachedUser, nil
	}

	var user models.User
	now := time.Now()

	result := r.db.Where("PSID = ?", psid).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}

	// Update last seen and cache the user
	err := r.db.Model(&user).Update("UpdatedAt", now).Error
	if err == nil {
		r.cache.SetUser(&user)
	}
	return &user, err
}

func (r *UserRepository) RegisterUser(psid string) (*models.User, error) {
	existingUser, err := r.GetUserByPSID(psid)
	if err == nil {
		if existingUser.IsActive {
			return existingUser, nil
		}

		existingUser.IsActive = true
		timeNow := time.Now()
		existingUser.LastLoginAt = &timeNow
		if err := r.db.Save(existingUser).Error; err != nil {
			return nil, fmt.Errorf("failed to reactivate user: %v", err)
		}
		r.cache.Invalidate(psid)
		return existingUser, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking for existing user: %v", err)
	}

	now := time.Now()
	fbName := "USER"
	var fbImgURL *string
	var email *string

	profile, err := r.fbSvc.GetUserProfile(psid)
	if err != nil {
		log.Printf("Warning: Could not fetch Facebook profile for user %s: %v", psid, err)
	} else {
		fbName = profile.Name
		if profile.Picture.Data.URL != "" {
			fbImgURL = &profile.Picture.Data.URL
		}
		if profile.Email != "" {
			email = &profile.Email
		}
	}

	code, err := utils.GenerateUniqueCode(r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %v", err)
	}

	user := models.User{
		PSID:        psid,
		FBName:      fbName,
		FBImgURL:    fbImgURL,
		Email:       email,
		Code:        &code,
		IsActive:    true,
		LastLoginAt: &now,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	r.cache.SetUser(&user)

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	r.cache.Invalidate(user.PSID)
	return r.db.Save(user).Error
}

func (r *UserRepository) MarkUserAsRegistered(psid string) error {
	r.cache.Invalidate(psid)
	return r.db.Model(&models.User{}).
		Where("PSID = ?", psid).
		Update("IsRegistered", true).Error
}

func (r *UserRepository) IsUserRegistered(psid string) (bool, error) {
	var user models.User
	if err := r.db.Select("IsRegistered").Where("PSID = ?", psid).First(&user).Error; err != nil {
		return false, err
	}
	return user.IsActive, nil
}

func (r *UserRepository) InvalidateCache(psid string) {
	r.cache.Invalidate(psid)
}

func (r *UserRepository) ClearCache() {
	r.cache.Clear()
}

func (r *UserRepository) PreloadActiveUsers() error {
	var users []*models.User
	result := r.db.Where("IsActive = ?", true).Find(&users)
	if result.Error != nil {
		return fmt.Errorf("failed to preload active users: %v", result.Error)
	}

	for _, user := range users {
		r.cache.SetUser(user)
	}

	return nil
}

// UserExists checks if a user with the given PSID exists and is active
func (r *UserRepository) UserExists(psid string) (bool, error) {
	if cachedUser, found := r.cache.GetUser(psid); found {
		return cachedUser.IsActive, nil
	}

	var exists bool
	err := r.db.Model(&models.User{}).
		Select("1").
		Where("PSID = ? AND IsActive = ?", psid, true).
		First(&exists).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("error checking if user exists: %v", err)
	}

	return true, nil
}
