package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/models"
	"github.com/helloellinor/p2k16/internal/session"
)

// ChiAuthenticatedUser represents the currently logged-in user for chi
type ChiAuthenticatedUser struct {
	ID       int             `json:"id"`
	Username string          `json:"username"`
	Account  *models.Account `json:"account,omitempty"`
}

// UserContextKey is used to store user in context
type contextKey string

const UserContextKey contextKey = "user"

// ChiRequireAuth middleware that requires authentication for chi
func ChiRequireAuth(sessionManager *session.ChiSessionManager, accountRepo *models.AccountRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sessionManager.IsAuthenticated(r) {
				// For HTMX requests, return HTML error
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`<div class="alert alert-warning">Please <a href="/login">log in</a> to continue.</div>`))
					return
				}

				// For regular requests, redirect to login
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Validate session is not expired
			if !sessionManager.ValidateSession(r) {
				sessionManager.LogoutUser(w, r)
				
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`<div class="alert alert-danger">Session expired. Please <a href="/login">log in again</a>.</div>`))
					return
				}

				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			userID := sessionManager.GetCurrentUserID(r)
			username := sessionManager.GetCurrentUsername(r)

			// In demo mode (no accountRepo), just use session info
			if accountRepo == nil {
				user := &ChiAuthenticatedUser{
					ID:       userID,
					Username: username,
					Account: &models.Account{
						ID:       userID,
						Username: username,
						Email:    username + "@demo.local",
					},
				}
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Load user account and add to context
			account, err := accountRepo.FindByID(userID)
			if err != nil {
				// User ID in session but account doesn't exist - clear session
				sessionManager.LogoutUser(w, r)

				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`<div class="alert alert-danger">Session invalid. Please <a href="/login">log in again</a>.</div>`))
					return
				}

				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			// Add user to context
			user := &ChiAuthenticatedUser{
				ID:       account.ID,
				Username: account.Username,
				Account:  account,
			}
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			
			// Update activity
			sessionManager.UpdateActivity(w, r)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ChiOptionalAuth middleware that loads user if authenticated but doesn't require it for chi
func ChiOptionalAuth(sessionManager *session.ChiSessionManager, accountRepo *models.AccountRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sessionManager.IsAuthenticated(r) && sessionManager.ValidateSession(r) {
				userID := sessionManager.GetCurrentUserID(r)
				username := sessionManager.GetCurrentUsername(r)

				// In demo mode (no accountRepo), just use session info
				if accountRepo == nil {
					user := &ChiAuthenticatedUser{
						ID:       userID,
						Username: username,
						Account: &models.Account{
							ID:       userID,
							Username: username,
							Email:    username + "@demo.local",
						},
					}
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					sessionManager.UpdateActivity(w, r)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				// Try to load user account
				account, err := accountRepo.FindByID(userID)
				if err == nil {
					user := &ChiAuthenticatedUser{
						ID:       account.ID,
						Username: account.Username,
						Account:  account,
					}
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					sessionManager.UpdateActivity(w, r)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ChiGetCurrentUser retrieves the current user from the context
func ChiGetCurrentUser(r *http.Request) *ChiAuthenticatedUser {
	if user := r.Context().Value(UserContextKey); user != nil {
		return user.(*ChiAuthenticatedUser)
	}
	return nil
}

// ChiIsAuthenticated checks if the current request is authenticated
func ChiIsAuthenticated(r *http.Request) bool {
	return ChiGetCurrentUser(r) != nil
}

// ChiLogger middleware for enhanced request logging
func ChiLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create a response writer wrapper to capture status code
			wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			next.ServeHTTP(wrapper, r)
			
			duration := time.Since(start)
			
			// Use our enhanced logger for HTTP requests
			logger := logging.ServerLogger
			
			// Log the request using our enhanced logger
			logger.LogRequest(
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				wrapper.statusCode,
				duration,
			)
		})
	}
}

// responseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ChiRecovery middleware for panic recovery
func ChiRecovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logging.LogError("PANIC", fmt.Sprintf("Panic recovered: %v", err))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ChiCORS middleware for handling cross-origin requests
func ChiCORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}