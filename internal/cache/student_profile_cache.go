package cache

import (
	"fmt"
	"school-assistant-wh/internal/models"
	"time"
)

type StudentProfileCache struct {
	cache *Cache
	ttl   time.Duration
}

func NewStudentProfileCache(ttl time.Duration) *StudentProfileCache {
	return &StudentProfileCache{
		cache: New(),
		ttl:   ttl,
	}
}

func (c *StudentProfileCache) getCacheKey(schoolID, studentID string) string {
	return fmt.Sprintf("student_profile:%s:%s", schoolID, studentID)
}

// Get retrieves a student profile from the cache
func (c *StudentProfileCache) Get(schoolID, studentID string) (*models.StudentProfile, bool) {
	key := c.getCacheKey(schoolID, studentID)
	if val, found := c.cache.Get(key); found {
		if profile, ok := val.(*models.StudentProfile); ok {
			return profile, true
		}
	}
	return nil, false
}

// Set adds a student profile to the cache
func (c *StudentProfileCache) Set(schoolID, studentID string, profile *models.StudentProfile) {
	if profile == nil || schoolID == "" || studentID == "" {
		return
	}
	key := c.getCacheKey(schoolID, studentID)
	c.cache.Set(key, profile, c.ttl)
}

// Invalidate removes a student profile from the cache
func (c *StudentProfileCache) Invalidate(schoolID, studentID string) {
	key := c.getCacheKey(schoolID, studentID)
	c.cache.Delete(key)
}

// Clear removes all student profiles from the cache
func (c *StudentProfileCache) Clear() {
	c.cache.Clear()
}
