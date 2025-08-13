package session

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/helloellinor/p2k16/internal/models"
)

// ChiSessionManager wraps the SCS session manager for chi
type ChiSessionManager struct {
	*scs.SessionManager
}

// NewChiSessionManager creates a new session manager for chi
func NewChiSessionManager() *ChiSessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour // 24 hours
	sessionManager.Cookie.Name = "p2k16-session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	
	return &ChiSessionManager{
		SessionManager: sessionManager,
	}
}

// SessionFromChi extracts session data from chi request
func (csm *ChiSessionManager) SessionFromChi(r *http.Request) (*SessionData, error) {
	userID := csm.GetInt(r.Context(), "user_id")
	username := csm.GetString(r.Context(), "username")
	lastActivityStr := csm.GetString(r.Context(), "last_activity")
	createdAtStr := csm.GetString(r.Context(), "created_at")
	
	if userID == 0 || username == "" {
		return nil, nil // No session data
	}
	
	data := &SessionData{
		UserID:   userID,
		Username: username,
	}
	
	if lastActivityStr != "" {
		if t, err := time.Parse(time.RFC3339, lastActivityStr); err == nil {
			data.LastActivity = t
		}
	}
	
	if createdAtStr != "" {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			data.CreatedAt = t
		}
	}
	
	return data, nil
}

// SaveToChi saves session data to chi session
func (csm *ChiSessionManager) SaveToChi(w http.ResponseWriter, r *http.Request, data *SessionData) error {
	csm.Put(r.Context(), "user_id", data.UserID)
	csm.Put(r.Context(), "username", data.Username)
	csm.Put(r.Context(), "last_activity", data.LastActivity.Format(time.RFC3339))
	csm.Put(r.Context(), "created_at", data.CreatedAt.Format(time.RFC3339))
	return nil
}

// LoginUser logs in a user by setting session variables for chi
func (csm *ChiSessionManager) LoginUser(w http.ResponseWriter, r *http.Request, account *models.Account) error {
	now := time.Now()
	
	csm.Put(r.Context(), "user_id", account.ID)
	csm.Put(r.Context(), "username", account.Username)
	csm.Put(r.Context(), "last_activity", now.Format(time.RFC3339))
	csm.Put(r.Context(), "created_at", now.Format(time.RFC3339))
	
	return nil
}

// LogoutUser logs out the current user by clearing the session for chi
func (csm *ChiSessionManager) LogoutUser(w http.ResponseWriter, r *http.Request) error {
	return csm.Destroy(r.Context())
}

// GetCurrentUserID gets the current user ID from session
func (csm *ChiSessionManager) GetCurrentUserID(r *http.Request) int {
	return csm.GetInt(r.Context(), "user_id")
}

// GetCurrentUsername gets the current username from session
func (csm *ChiSessionManager) GetCurrentUsername(r *http.Request) string {
	return csm.GetString(r.Context(), "username")
}

// IsAuthenticated checks if the current request is authenticated
func (csm *ChiSessionManager) IsAuthenticated(r *http.Request) bool {
	return csm.GetCurrentUserID(r) != 0
}

// UpdateActivity updates the last activity timestamp
func (csm *ChiSessionManager) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	if csm.IsAuthenticated(r) {
		csm.Put(r.Context(), "last_activity", time.Now().Format(time.RFC3339))
	}
}

// ValidateSession checks if session is still valid (not expired)
func (csm *ChiSessionManager) ValidateSession(r *http.Request) bool {
	if !csm.IsAuthenticated(r) {
		return false
	}
	
	lastActivityStr := csm.GetString(r.Context(), "last_activity")
	if lastActivityStr == "" {
		return false
	}
	
	lastActivity, err := time.Parse(time.RFC3339, lastActivityStr)
	if err != nil {
		return false
	}
	
	// Check if session is expired (24 hours)
	return time.Since(lastActivity) <= 24*time.Hour
}

// ToJSON converts session data to JSON for debugging
func (data *SessionData) ToJSONChi() string {
	b, _ := json.Marshal(data)
	return string(b)
}