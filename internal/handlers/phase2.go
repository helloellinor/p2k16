package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// ChangePassword handles password change requests (Phase 2)
func (h *Handler) ChangePassword(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	oldPassword := c.PostForm("oldPassword")
	newPassword := c.PostForm("newPassword")
	confirmPassword := c.PostForm("confirmPassword")

	// Validation
	if oldPassword == "" || newPassword == "" || confirmPassword == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">All fields are required</div>`))
		return
	}

	if newPassword != confirmPassword {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">New passwords do not match</div>`))
		return
	}

	if len(newPassword) < 6 {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Password must be at least 6 characters long</div>`))
		return
	}

	// Handle demo mode
	if h.demoMode || h.accountRepo == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--success">Password changed successfully (demo mode)</div>`))
		return
	}

	// Get current account
	account, err := h.accountRepo.FindByID(user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Failed to load account</div>`))
		return
	}

	// Verify old password
	if !account.ValidatePassword(oldPassword) {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Current password is incorrect</div>`))
		return
	}

	// Update password
	if err := account.SetPassword(newPassword); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Failed to hash new password</div>`))
		return
	}

	// Save to database
	if err := h.accountRepo.UpdatePassword(account.ID, account.Password); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Failed to save new password</div>`))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8",
		[]byte(`<div class="p2k16-alert p2k16-alert--success">Password changed successfully!</div>`))
}

// UpdateProfile handles profile update requests (Phase 2)
func (h *Handler) UpdateProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	name := c.PostForm("name")
	phone := c.PostForm("phone")

	// Handle demo mode
	if h.demoMode || h.accountRepo == nil {
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

	// Get current account
	account, err := h.accountRepo.FindByID(user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Failed to load account</div>`))
		return
	}

	// Update fields
	updated := false
	if name != "" && (!account.Name.Valid || account.Name.String != name) {
		account.Name.Valid = true
		account.Name.String = name
		updated = true
	}
	if phone != "" && (!account.Phone.Valid || account.Phone.String != phone) {
		account.Phone.Valid = true
		account.Phone.String = phone
		updated = true
	}

	if !updated {
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--info">No changes made</div>`))
		return
	}

	// Save to database
	if err := h.accountRepo.UpdateProfile(account); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="p2k16-alert p2k16-alert--error">Failed to save profile changes</div>`))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8",
		[]byte(`<div class="p2k16-alert p2k16-alert--success">Profile updated successfully!</div>`))
}