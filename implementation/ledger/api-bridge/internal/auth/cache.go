package auth

import (
	"sync"
	"time"
)

type AuthCache struct {
	cache    map[string]*CacheEntry
	mutex    sync.RWMutex
	duration time.Duration
}

type CacheEntry struct {
	KeyInfo   *KeyInfo
	ExpiresAt time.Time
}

type KeyInfo struct {
	UserID        int      `json:"user_id"`
	Username      string   `json:"username"`
	ComponentID   string   `json:"component_id"`
	DeviceType    string   `json:"device_type"`
	Permissions   []string `json:"permissions"`
	RateLimit     int      `json:"rate_limit"`
	MaxConcurrent int      `json:"max_concurrent"`
	Roles         []string `json:"roles"`
}

func NewAuthCache(duration time.Duration) *AuthCache {
	cache := &AuthCache{
		cache:    make(map[string]*CacheEntry),
		duration: duration,
	}

	// Start cleanup routine
	go cache.cleanup()

	return cache
}

func (a *AuthCache) Get(apiKey string) *KeyInfo {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	entry, exists := a.cache[apiKey]
	if !exists {
		return nil
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry expired
		delete(a.cache, apiKey)
		return nil
	}

	return entry.KeyInfo
}

func (a *AuthCache) Set(apiKey string, keyInfo *KeyInfo) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.cache[apiKey] = &CacheEntry{
		KeyInfo:   keyInfo,
		ExpiresAt: time.Now().Add(a.duration),
	}
}

func (a *AuthCache) Delete(apiKey string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	delete(a.cache, apiKey)
}

func (a *AuthCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.mutex.Lock()
			now := time.Now()
			for key, entry := range a.cache {
				if now.After(entry.ExpiresAt) {
					delete(a.cache, key)
				}
			}
			a.mutex.Unlock()
		}
	}
}
