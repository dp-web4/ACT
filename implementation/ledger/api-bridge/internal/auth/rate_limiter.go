package auth

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type RateLimiter struct {
	clients map[int]*ClientLimit
	mutex   sync.RWMutex
	logger  zerolog.Logger
}

type ClientLimit struct {
	requests     []time.Time
	lastActivity time.Time
}

func NewRateLimiter(logger zerolog.Logger) *RateLimiter {
	limiter := &RateLimiter{
		clients: make(map[int]*ClientLimit),
		logger:  logger,
	}

	// Start cleanup routine
	go limiter.cleanup()

	return limiter
}

func (r *RateLimiter) IsAllowed(userID int, limit int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	client, exists := r.clients[userID]

	if !exists {
		client = &ClientLimit{
			requests:     []time.Time{now},
			lastActivity: now,
		}
		r.clients[userID] = client
		return true
	}

	// Remove old requests (older than 1 minute)
	validRequests := []time.Time{}
	for _, req := range client.requests {
		if now.Sub(req) < time.Minute {
			validRequests = append(validRequests, req)
		}
	}

	if len(validRequests) >= limit {
		return false
	}

	client.requests = append(validRequests, now)
	client.lastActivity = now
	return true
}

func (r *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mutex.Lock()
			now := time.Now()
			for userID, client := range r.clients {
				if now.Sub(client.lastActivity) > 24*time.Hour {
					delete(r.clients, userID)
				}
			}
			r.mutex.Unlock()
		}
	}
}
