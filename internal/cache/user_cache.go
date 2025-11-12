package cache

import (
	"school-assistant-wh/internal/models"
	"time"
)

type UserCache struct {
	cache *Cache
	ttl   time.Duration
}

func NewUserCache(ttl time.Duration) *UserCache {
	return &UserCache{
		cache: New(),
		ttl:   ttl,
	}
}

// GetUser retrieves a user from cache by PSID
func (uc *UserCache) GetUser(psid string) (*models.User, bool) {
	if val, found := uc.cache.Get(psid); found {
		if user, ok := val.(*models.User); ok {
			return user, true
		}
	}
	return nil, false
}

// SetUser adds a user to the cache with the configured TTL
func (uc *UserCache) SetUser(user *models.User) {
	if user == nil || user.PSID == "" {
		return
	}
	uc.cache.Set(user.PSID, user, uc.ttl)
}

// Invalidate removes a user from the cache
func (uc *UserCache) Invalidate(psid string) {
	uc.cache.Delete(psid)
}

// Clear removes all users from the cache
func (uc *UserCache) Clear() {
	uc.cache.Clear()
}
