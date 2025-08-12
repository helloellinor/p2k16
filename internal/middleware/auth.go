package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/models"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
	SessionName = "p2k16-session"
)

// AuthenticatedUser represents the currently logged-in user
type AuthenticatedUser struct {
	ID       int             `json:"id"`
	Username string          `json:"username"`
	Account  *models.Account `json:"account,omitempty"`
}

// RequireAuth middleware that requires authentication
func RequireAuth(accountRepo *models.AccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(UserIDKey)

		if userID == nil {
			// For HTMX requests, return HTML error
			if c.GetHeader("HX-Request") == "true" {
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-warning">Please <a href="/login">log in</a> to continue.</div>`))
				c.Abort()
				return
			}

			// For regular requests, redirect to login
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Load user account and add to context
		account, err := accountRepo.FindByID(userID.(int))
		if err != nil {
			// User ID in session but account doesn't exist - clear session
			session.Clear()
			session.Save()

			if c.GetHeader("HX-Request") == "true" {
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Session invalid. Please <a href="/login">log in again</a>.</div>`))
				c.Abort()
				return
			}

			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Add user to context
		user := &AuthenticatedUser{
			ID:       account.ID,
			Username: account.Username,
			Account:  account,
		}
		c.Set("user", user)
		c.Next()
	}
}

// OptionalAuth middleware that loads user if authenticated but doesn't require it
func OptionalAuth(accountRepo *models.AccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(UserIDKey)

		if userID != nil {
			// Try to load user account
			account, err := accountRepo.FindByID(userID.(int))
			if err == nil {
				user := &AuthenticatedUser{
					ID:       account.ID,
					Username: account.Username,
					Account:  account,
				}
				c.Set("user", user)
			}
		}

		c.Next()
	}
}

// GetCurrentUser retrieves the current user from the context
func GetCurrentUser(c *gin.Context) *AuthenticatedUser {
	if user, exists := c.Get("user"); exists {
		return user.(*AuthenticatedUser)
	}
	return nil
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	return GetCurrentUser(c) != nil
}

// LoginUser logs in a user by setting session variables
func LoginUser(c *gin.Context, account *models.Account) error {
	session := sessions.Default(c)
	session.Set(UserIDKey, account.ID)
	session.Set(UsernameKey, account.Username)
	return session.Save()
}

// LogoutUser logs out the current user by clearing the session
func LogoutUser(c *gin.Context) error {
	session := sessions.Default(c)
	session.Clear()
	return session.Save()
}
