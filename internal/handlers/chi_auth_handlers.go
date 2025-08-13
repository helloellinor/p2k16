package handlers

import (
	"fmt"
	"net/http"

	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// ChiAuthLogin handles login form submission for chi
func (h *ChiHandler) ChiAuthLogin(w http.ResponseWriter, r *http.Request) {
	logging.LogHandlerAction("API REQUEST", "Login attempt received")
	
	if err := r.ParseForm(); err != nil {
		logging.LogError("LOGIN FAILED", "Failed to parse form data")
		h.writeHTML(w, http.StatusBadRequest, `<p>Invalid form data</p>`)
		return
	}
	
	username := r.FormValue("username")
	password := r.FormValue("password")

	logging.LogHandlerAction("LOGIN ATTEMPT", fmt.Sprintf("Username: %s", username))

	if username == "" || password == "" {
		logging.LogError("LOGIN FAILED", "Missing username or password")
		h.writeHTML(w, http.StatusBadRequest, `<p>Username and password are required</p>`)
		return
	}

	// Database authentication
	account, err := h.accountRepo.FindByUsername(username)
	if err != nil {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("User '%s' not found in database: %v", username, err))
		h.writeHTML(w, http.StatusUnauthorized, `<p>Invalid username or password</p>`)
		return
	}

	logging.LogHandlerAction("USER FOUND", fmt.Sprintf("User '%s' found in database, validating password", username))
	if !account.ValidatePassword(password) {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("Invalid password for user '%s'", username))
		h.writeHTML(w, http.StatusUnauthorized, `<p>Invalid username or password</p>`)
		return
	}

	logging.LogSuccess("LOGIN SUCCESS", fmt.Sprintf("User authenticated: %s", username))

	// Login user by setting session
	if err := h.sessionManager.LoginUser(w, r, account); err != nil {
		logging.LogError("SESSION ERROR", fmt.Sprintf("Failed to create session for user %s: %v", username, err))
		h.writeHTML(w, http.StatusInternalServerError, `<p>Failed to login. Please try again.</p>`)
		return
	}

	// Successful login - redirect via HTMX
	w.Header().Set("HX-Redirect", "/")
	h.writeHTML(w, http.StatusOK, `<p>Login successful! Redirecting...</p>`)
}

// ChiProfileCardFront returns the front of the membership card (requires auth) for chi
func (h *ChiHandler) ChiProfileCardFront(w http.ResponseWriter, r *http.Request) {
	user := ChiGetCurrentUser(r)
	html := h.renderChiProfileCardFrontHTML(user)
	h.writeHTML(w, http.StatusOK, html)
}

// ChiProfileCardBack returns the back (editing) of the membership card (requires auth) for chi
func (h *ChiHandler) ChiProfileCardBack(w http.ResponseWriter, r *http.Request) {
	user := ChiGetCurrentUser(r)
	html := h.renderChiProfileCardBackHTML(user)
	h.writeHTML(w, http.StatusOK, html)
}

// ChiChangePassword handles password change requests for chi
func (h *ChiHandler) ChiChangePassword(w http.ResponseWriter, r *http.Request) {
	user := ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, `<p>Authentication required</p>`)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.writeHTML(w, http.StatusBadRequest, `<p>Invalid form data</p>`)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if currentPassword == "" {
		h.writeHTML(w, http.StatusBadRequest, `<p>Current password is required</p>`)
		return
	}

	// Validate current password
	if !user.Account.ValidatePassword(currentPassword) {
		h.writeHTML(w, http.StatusUnauthorized, `<p>Current password is incorrect</p>`)
		return
	}

	// If new password is provided, validate and update
	if newPassword != "" {
		if newPassword != confirmPassword {
			h.writeHTML(w, http.StatusBadRequest, `<p>New password and confirmation do not match</p>`)
			return
		}

		if len(newPassword) < 8 {
			h.writeHTML(w, http.StatusBadRequest, `<p>New password must be at least 8 characters long</p>`)
			return
		}

		// Update password
		if err := h.accountRepo.UpdatePassword(user.Account.ID, newPassword); err != nil {
			logging.LogError("PASSWORD UPDATE", fmt.Sprintf("Failed to update password for user %s: %v", user.Username, err))
			h.writeHTML(w, http.StatusInternalServerError, `<p>Failed to update password</p>`)
			return
		}

		logging.LogSuccess("PASSWORD UPDATE", fmt.Sprintf("Password updated for user %s", user.Username))
	}

	// Return updated card front
	html := h.renderChiProfileCardFrontHTML(user)
	h.writeHTML(w, http.StatusOK, html)
}

// ChiUpdateProfile handles profile update requests for chi
func (h *ChiHandler) ChiUpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, `<p>Authentication required</p>`)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.writeHTML(w, http.StatusBadRequest, `<p>Invalid form data</p>`)
		return
	}

	currentPassword := r.FormValue("current_password")
	email := r.FormValue("email")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if currentPassword == "" {
		h.writeHTML(w, http.StatusBadRequest, `<p>Current password is required to save changes</p>`)
		return
	}

	// Validate current password
	if !user.Account.ValidatePassword(currentPassword) {
		h.writeHTML(w, http.StatusUnauthorized, `<p>Current password is incorrect</p>`)
		return
	}

	// Update email if provided and different
	if email != "" && email != user.Account.Email {
		if err := h.accountRepo.UpdateEmail(user.Account.ID, email); err != nil {
			logging.LogError("EMAIL UPDATE", fmt.Sprintf("Failed to update email for user %s: %v", user.Username, err))
			h.writeHTML(w, http.StatusInternalServerError, `<p>Failed to update email</p>`)
			return
		}
		user.Account.Email = email
		logging.LogSuccess("EMAIL UPDATE", fmt.Sprintf("Email updated for user %s", user.Username))
	}

	// Update password if provided
	if newPassword != "" {
		if newPassword != confirmPassword {
			h.writeHTML(w, http.StatusBadRequest, `<p>New password and confirmation do not match</p>`)
			return
		}

		if len(newPassword) < 8 {
			h.writeHTML(w, http.StatusBadRequest, `<p>New password must be at least 8 characters long</p>`)
			return
		}

		if err := h.accountRepo.UpdatePassword(user.Account.ID, newPassword); err != nil {
			logging.LogError("PASSWORD UPDATE", fmt.Sprintf("Failed to update password for user %s: %v", user.Username, err))
			h.writeHTML(w, http.StatusInternalServerError, `<p>Failed to update password</p>`)
			return
		}

		logging.LogSuccess("PASSWORD UPDATE", fmt.Sprintf("Password updated for user %s", user.Username))
	}

	// Return updated card front
	html := h.renderChiProfileCardFrontHTML(user)
	h.writeHTML(w, http.StatusOK, html)
}

// Helper function to get current user (fix import issue)
func ChiGetCurrentUser(r *http.Request) *middleware.ChiAuthenticatedUser {
	if user := r.Context().Value(middleware.UserContextKey); user != nil {
		return user.(*middleware.ChiAuthenticatedUser)
	}
	return nil
}