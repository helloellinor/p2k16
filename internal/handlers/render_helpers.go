package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// renderNavbar returns a minimal, classless navbar based on auth state
func (h *Handler) renderNavbar(c *gin.Context) string {
	user := middleware.GetCurrentUser(c)
	html := `<header><div><a href="/">P2K16</a>`
	if user != nil {
		html += `<form method="post" action="/logout" style="display:inline;margin-left:8px;"><button type="submit">Logout</button></form>`
	}
	html += `</div></header>`
	return html
}

// renderNavbarWithTrail is like renderNavbar but appends an inline
// breadcrumb trail (e.g., " / Profile") right after the brand.
func (h *Handler) renderNavbarWithTrail(c *gin.Context, trail string) string {
	user := middleware.GetCurrentUser(c)
	html := `<header><div><a href="/">P2K16</a>`
	if trail != "" {
		// If trail is 'Profile' or 'Admin', render as links
		switch trail {
		case "Profile":
			html += ` / <a href="/profile">Profile</a>`
		case "Admin":
			html += ` / <a href="/admin">Admin</a>`
		default:
			// For compound trails like 'Admin / Users', split and link first part if possible
			if trail == "Admin / Users" {
				html += ` / <a href="/admin">Admin</a> / <span>Users</span>`
			} else if trail == "Admin / Tools" {
				html += ` / <a href="/admin">Admin</a> / <span>Tools</span>`
			} else if trail == "Admin / Companies" {
				html += ` / <a href="/admin">Admin</a> / <span>Companies</span>`
			} else if trail == "Admin / Circles" {
				html += ` / <a href="/admin">Admin</a> / <span>Circles</span>`
			} else if trail == "Admin / Logs" {
				html += ` / <a href="/admin">Admin</a> / <span>Logs</span>`
			} else if trail == "Admin / Config" {
				html += ` / <a href="/admin">Admin</a> / <span>Config</span>`
			} else if trail == "Create Badge" {
				html += ` / <a href="/badges/new">Create Badge</a>`
			} else {
				html += ` / <span>` + trail + `</span>`
			}
		}
	}
	if user != nil {
		html += `<form method="post" action="/logout" style="display:inline;margin-left:8px;"><button type="submit">Logout</button></form>`
	}
	html += `<nav>`
	if user == nil {
		html += `<a href="/login">Login</a>`
	} else {
		html += `<a href="/profile">Profile</a><a href="/admin">Admin</a>`
	}
	html += `</nav></div></header>`
	return html
}

// renderUserBadgesSectionHTML builds the user badges section (classless, HTMX-friendly)
func (h *Handler) renderUserBadgesSectionHTML(accountID int) string {
	badges, _ := h.badgeRepo.GetBadgesForAccount(accountID)
	if len(badges) == 0 {
		return `
<section id="user-badges" aria-labelledby="user-badges-title">
	<h2 id="user-badges-title">Your Badges</h2>
	<p>No badges yet.</p>
	<button hx-get="/api/badges/available" hx-target="#available-badges" hx-swap="innerHTML">Browse Available Badges</button>
	<div id="available-badges"></div>
		  <div id="badge-feedback" aria-live="polite"></div>
</section>`
	}
	html := `
<section id="user-badges" aria-labelledby="user-badges-title">
	<h2 id="user-badges-title">Your Badges</h2>
		<ul>`
	for _, badge := range badges {
		html += `
			<li>
				<span>` + badge.BadgeDescription.Title + `</span>
				<button 
					hx-post="/api/badges/remove" 
					hx-vals='{"account_badge_id":"` + strconv.Itoa(badge.ID) + `"}'
					hx-target="#badge-feedback"
					hx-swap="innerHTML">Remove</button>
			</li>`
	}
	html += `
	</ul>
		<p>You have ` + fmt.Sprintf("%d", len(badges)) + ` badge(s).</p>
	<button hx-get="/api/badges/available" hx-target="#available-badges" hx-swap="innerHTML">Browse More Badges</button>
	<div id="available-badges"></div>
		<div id="badge-feedback" aria-live="polite"></div>
</section>`
	return html
}

// renderUserBadgesListReadOnly renders a simple list of badges without actions
func (h *Handler) renderUserBadgesListReadOnly(accountID int) string {
	badges, _ := h.badgeRepo.GetBadgesForAccount(accountID)
	html := `<section aria-labelledby="badges-title">
		<h2 id="badges-title">Badges</h2>`
	if len(badges) == 0 {
		html += `<p>No badges yet.</p>`
	} else {
		html += `<ul>`
		for _, badge := range badges {
			html += `<li><span>` + badge.BadgeDescription.Title + `</span></li>`
		}
		html += `</ul>`
		html += `<p>You have ` + fmt.Sprintf("%d", len(badges)) + ` badge(s).</p>`
	}
	html += `</section>`
	return html
}

// renderProfileCardFrontHTML composes the front of the membership card
func (h *Handler) renderProfileCardFrontHTML(user *middleware.AuthenticatedUser) string {
	info := `<section aria-labelledby="info-title"><h2 id="info-title">Member</h2>` +
		`<p><strong>Username:</strong> ` + user.Username + `</p>` +
		`<p><strong>Email:</strong> ` + user.Account.Email + `</p>`
	if user.Account.Name.Valid && user.Account.Name.String != "" {
		info += `<p><strong>Name:</strong> ` + user.Account.Name.String + `</p>`
	}
	if user.Account.Phone.Valid && user.Account.Phone.String != "" {
		info += `<p><strong>Phone:</strong> ` + user.Account.Phone.String + `</p>`
	}
	info += `</section>`

	badges := h.renderUserBadgesListReadOnly(user.ID)

	html := `<div>` +
		`<div><button hx-get="/api/profile/card/back" hx-target="#membership-card" hx-swap="innerHTML" aria-label="Edit membership card">Edit</button></div>` +
		info + badges + `</div>`
	return html
}

// renderProfileCardBackHTML composes the back of the membership card (editing)
func (h *Handler) renderProfileCardBackHTML(user *middleware.AuthenticatedUser) string {
	// Change Password form
	changePassword := `
<section>
	<header><h2>Change Password</h2></header>
	<div>
		<form hx-post="/api/profile/change-password" hx-target="#password-result">
			<div>
				<label for="oldPassword">Current Password</label>
				<input type="password" id="oldPassword" name="oldPassword" required>
			</div>
			<div>
				<label for="newPassword">New Password</label>
				<input type="password" id="newPassword" name="newPassword" required>
			</div>
			<div>
				<label for="confirmPassword">Confirm New Password</label>
				<input type="password" id="confirmPassword" name="confirmPassword" required>
			</div>
			<button type="submit">Change Password</button>
		</form>
		<div id="password-result"></div>
	</div>
</section>`

	// Profile details form
	details := `
<section>
	<header><h2>Profile Details</h2></header>
	<div>
		<form hx-post="/api/profile/update" hx-target="#profile-result">
			<div>
				<label for="name">Full Name</label>
				<input type="text" id="name" name="name" placeholder="Enter your full name">
				<div>This will be displayed on your public profile</div>
			</div>
			<div>
				<label for="email">Email Address</label>
				<input type="email" id="email" name="email" value="` + user.Account.Email + `" readonly>
				<div>Contact an administrator to change your email address</div>
			</div>
			<div>
				<label for="phone">Phone Number</label>
				<input type="tel" id="phone" name="phone" placeholder="Enter your phone number">
				<div>Used for emergency contact and door access notifications</div>
			</div>
			<button type="submit">Save Changes</button>
		</form>
		<div id="profile-result"></div>
	</div>
</section>`

	// Editable badges section
	badges := h.renderUserBadgesSectionHTML(user.ID)

	html := `<div>` +
		`<div><button hx-get="/api/profile/card/front" hx-target="#membership-card" hx-swap="innerHTML" aria-label="Done editing">Done</button></div>` +
		changePassword + details + badges + `</div>`
	return html
}