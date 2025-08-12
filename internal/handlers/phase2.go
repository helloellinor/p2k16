package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// ChangePassword handles password change requests (Phase 2)
func (h *Handler) ChangePassword(c *gin.Context) {
	logging.LogHandlerAction("PASSWORD CHANGE", "Password change request received")
	user := middleware.GetCurrentUser(c)
	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")
	confirmPassword := c.PostForm("confirmPassword")

	logging.LogHandlerAction("USER REQUEST", fmt.Sprintf("Password change for user ID: %d", user.ID))

	// Validation
	if oldPassword == "" || newPassword == "" || confirmPassword == "" {
		logging.LogHandlerAction("VALIDATION ERROR", "Missing required fields")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">All fields are required</div>`))
		return
	}

	if newPassword != confirmPassword {
		logging.LogHandlerAction("VALIDATION ERROR", "New passwords do not match")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">New passwords do not match</div>`))
		return
	}

	if len(newPassword) < 6 {
		logging.LogHandlerAction("VALIDATION ERROR", "Password too short")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Password must be at least 6 characters long</div>`))
		return
	}

	// Handle demo mode
	if h.demoMode || h.accountRepo == nil {
		logging.LogHandlerAction("DEMO MODE", "Password change simulated - no database update")
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--success">Password changed successfully (demo mode)</div>`))
		return
	}

	logging.LogHandlerAction("DATABASE OPERATION", "Fetching current account for password verification")
	// Get current account
	account, err := h.accountRepo.FindByID(user.ID)
	if err != nil {
		logging.LogError("DATABASE ERROR", fmt.Sprintf("Failed to load account: %v", err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Failed to load account</div>`))
		return
	}

	// Verify old password
	if !account.ValidatePassword(oldPassword) {
		logging.LogHandlerAction("AUTH ERROR", "Current password verification failed")
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Current password is incorrect</div>`))
		return
	}

	logging.LogHandlerAction("PASSWORD VALIDATION", "Current password verified successfully")

	// Update password
	if err := account.SetPassword(newPassword); err != nil {
		logging.LogError("HASH ERROR", fmt.Sprintf("Failed to hash new password: %v", err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Failed to hash new password</div>`))
		return
	}

	// Save to database
	if err := h.accountRepo.UpdatePassword(account.ID, account.Password); err != nil {
		logging.LogError("DATABASE ERROR", fmt.Sprintf("Failed to save new password: %v", err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Failed to save new password</div>`))
		return
	}

	logging.LogSuccess("PASSWORD CHANGE SUCCESS", fmt.Sprintf("Password updated successfully for user ID: %d", user.ID))
	c.Data(http.StatusOK, "text/html; charset=utf-8",
		[]byte(`<div class="p2k16-alert p2k16-alert--success">Password changed successfully!</div>`))
}

// UpdateProfile handles profile update requests (Phase 2)
func (h *Handler) UpdateProfile(c *gin.Context) {
	logging.LogHandlerAction("PROFILE UPDATE", "Profile update request received")
	user := middleware.GetCurrentUser(c)
	name := c.PostForm("name")
	phone := c.PostForm("phone")

	logging.LogHandlerAction("USER REQUEST", fmt.Sprintf("Profile update for user ID: %d - Name: %s, Phone: %s", user.ID, name, phone))

	// Handle demo mode
	if h.demoMode || h.accountRepo == nil {
		logging.LogHandlerAction("DEMO MODE", "Profile update simulated - no database update")
		message := "Profile updated successfully (demo mode)"
		if name != "" {
			message += " - Name: " + name
		}
		if phone != "" {
			message += " - Phone: " + phone
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--success">`+message+`</div>`))
		return
	}

	logging.LogHandlerAction("DATABASE OPERATION", "Fetching current account for profile update")
	// Get current account
	account, err := h.accountRepo.FindByID(user.ID)
	if err != nil {
		logging.LogError("DATABASE ERROR", fmt.Sprintf("Failed to load account: %v", err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Failed to load account</div>`))
		return
	}

	// Update fields
	updated := false
	if name != "" && (!account.Name.Valid || account.Name.String != name) {
		logging.LogHandlerAction("FIELD UPDATE", fmt.Sprintf("Updating name: %s -> %s", account.Name.String, name))
		account.Name.Valid = true
		account.Name.String = name
		updated = true
	}
	if phone != "" && (!account.Phone.Valid || account.Phone.String != phone) {
		logging.LogHandlerAction("FIELD UPDATE", fmt.Sprintf("Updating phone: %s -> %s", account.Phone.String, phone))
		account.Phone.Valid = true
		account.Phone.String = phone
		updated = true
	}

	if !updated {
		logging.LogHandlerAction("NO CHANGES", "No profile changes detected")
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--info">No changes made</div>`))
		return
	}

	// Save to database
	if err := h.accountRepo.UpdateProfile(account); err != nil {
		logging.LogError("DATABASE ERROR", fmt.Sprintf("Failed to save profile changes: %v", err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--danger">Failed to save profile changes</div>`))
		return
	}

	logging.LogSuccess("PROFILE UPDATE SUCCESS", fmt.Sprintf("Profile updated successfully for user ID: %d", user.ID))
	c.Data(http.StatusOK, "text/html; charset=utf-8",
		[]byte(`<div class="p2k16-alert p2k16-alert--success">Profile updated successfully!</div>`))
}