package session

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	SessionTimeout = 24 * time.Hour // 24 hours
	CleanupInterval = 1 * time.Hour // Clean up every hour
)

// SessionData represents the data stored in a session
type SessionData struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
}

// SessionStore manages session storage and cleanup
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionData
	cleanup  *time.Ticker
	stopCh   chan struct{}
}

// NewSessionStore creates a new session store with automatic cleanup
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*SessionData),
		cleanup:  time.NewTicker(CleanupInterval),
		stopCh:   make(chan struct{}),
	}

	// Start background cleanup
	go store.cleanupLoop()

	return store
}

// Get retrieves session data for a session ID
func (s *SessionStore) Get(sessionID string) (*SessionData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data, exists := s.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session is expired
	if time.Since(data.LastActivity) > SessionTimeout {
		return nil, false
	}

	return data, true
}

// Set stores session data for a session ID
func (s *SessionStore) Set(sessionID string, data *SessionData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data.LastActivity = time.Now()
	s.sessions[sessionID] = data
}

// Delete removes a session
func (s *SessionStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.sessions, sessionID)
}

// UpdateActivity updates the last activity time for a session
func (s *SessionStore) UpdateActivity(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if data, exists := s.sessions[sessionID]; exists {
		data.LastActivity = time.Now()
		return true
	}
	return false
}

// CleanupExpired removes expired sessions
func (s *SessionStore) CleanupExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	count := 0
	
	for sessionID, data := range s.sessions {
		if now.Sub(data.LastActivity) > SessionTimeout {
			delete(s.sessions, sessionID)
			count++
		}
	}
	
	return count
}

// Count returns the number of active sessions
func (s *SessionStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return len(s.sessions)
}

// Stop stops the cleanup loop
func (s *SessionStore) Stop() {
	close(s.stopCh)
	s.cleanup.Stop()
}

// cleanupLoop runs in the background to clean up expired sessions
func (s *SessionStore) cleanupLoop() {
	for {
		select {
		case <-s.cleanup.C:
			count := s.CleanupExpired()
			if count > 0 {
				// Only log if we actually cleaned up sessions
				// log.Printf("Cleaned up %d expired sessions", count)
			}
		case <-s.stopCh:
			return
		}
	}
}

// SessionFromGin extracts session data from Gin session
func SessionFromGin(c *gin.Context) (*SessionData, error) {
	session := sessions.Default(c)
	
	userID := session.Get("user_id")
	username := session.Get("username")
	lastActivity := session.Get("last_activity")
	createdAt := session.Get("created_at")
	
	if userID == nil || username == nil {
		return nil, nil // No session data
	}
	
	data := &SessionData{
		UserID:   userID.(int),
		Username: username.(string),
	}
	
	if lastActivity != nil {
		if t, ok := lastActivity.(time.Time); ok {
			data.LastActivity = t
		} else if s, ok := lastActivity.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				data.LastActivity = t
			}
		}
	}
	
	if createdAt != nil {
		if t, ok := createdAt.(time.Time); ok {
			data.CreatedAt = t
		} else if s, ok := createdAt.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				data.CreatedAt = t
			}
		}
	}
	
	return data, nil
}

// SaveToGin saves session data to Gin session
func (data *SessionData) SaveToGin(c *gin.Context) error {
	session := sessions.Default(c)
	
	session.Set("user_id", data.UserID)
	session.Set("username", data.Username)
	session.Set("last_activity", data.LastActivity.Format(time.RFC3339))
	session.Set("created_at", data.CreatedAt.Format(time.RFC3339))
	
	return session.Save()
}

// ToJSON converts session data to JSON for debugging
func (data *SessionData) ToJSON() string {
	b, _ := json.Marshal(data)
	return string(b)
}