package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// renderNavbar returns a Bootstrap navbar based on auth state
func (h *Handler) renderNavbar(c *gin.Context) string {
	user := middleware.GetCurrentUser(c)

	html := `
<nav class="navbar navbar-expand-lg navbar-dark bg-primary">
	<div class="container">
		<a class="navbar-brand fw-bold" href="/">P2K16</a>
		
		<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav me-auto">
				<li class="nav-item">
					<a class="nav-link" href="/">Home</a>
				</li>`

	if user != nil {
		html += `
				<li class="nav-item">
					<a class="nav-link" href="/profile">Profile</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="/admin">Admin</a>
				</li>`
	}

	html += `
			</ul>
			
			<ul class="navbar-nav">
				<li class="nav-item">`

	if user == nil {
		html += `
					<a class="nav-link" href="/login">Login</a>`
	} else {
		html += `
					<span class="navbar-text me-3">Welcome, ` + user.Username + `</span>
				</li>
				<li class="nav-item">
					<form method="post" action="/logout" class="d-inline">
						<button type="submit" class="btn btn-outline-light btn-sm">Logout</button>
					</form>`
	}

	html += `
				</li>
			</ul>
		</div>
	</div>
</nav>`

	return html
}

// renderNavbarWithTrail is like renderNavbar but includes breadcrumb navigation
func (h *Handler) renderNavbarWithTrail(c *gin.Context, trail string) string {
	user := middleware.GetCurrentUser(c)

	html := `
<nav class="navbar navbar-expand-lg navbar-dark bg-primary">
	<div class="container">
		<a class="navbar-brand fw-bold" href="/">P2K16</a>
		
		<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav me-auto">
				<li class="nav-item">
					<a class="nav-link" href="/">Home</a>
				</li>`

	if user != nil {
		html += `
				<li class="nav-item">
					<a class="nav-link" href="/profile">Profile</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="/admin">Admin</a>
				</li>`
	}

	html += `
			</ul>
			
			<ul class="navbar-nav">
				<li class="nav-item">`

	if user == nil {
		html += `
					<a class="nav-link" href="/login">Login</a>`
	} else {
		html += `
					<span class="navbar-text me-3">Welcome, ` + user.Username + `</span>
				</li>
				<li class="nav-item">
					<form method="post" action="/logout" class="d-inline">
						<button type="submit" class="btn btn-outline-light btn-sm">Logout</button>
					</form>`
	}

	html += `
				</li>
			</ul>
		</div>
	</div>
</nav>`

	// Add breadcrumb if trail is provided
	if trail != "" {
		html += `
<nav aria-label="breadcrumb" class="bg-light border-bottom">
	<div class="container">
		<ol class="breadcrumb py-2 mb-0">
			<li class="breadcrumb-item"><a href="/">Home</a></li>`

		// Handle different trail types
		switch trail {
		case "Profile":
			html += `<li class="breadcrumb-item active" aria-current="page">Profile</li>`
		case "Admin":
			html += `<li class="breadcrumb-item active" aria-current="page">Admin</li>`
		default:
			// For compound trails like 'Admin / Users', split and link first part if possible
			if trail == "Admin / Users" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Users</li>`
			} else if trail == "Admin / Tools" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Tools</li>`
			} else if trail == "Admin / Companies" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Companies</li>`
			} else if trail == "Admin / Circles" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Circles</li>`
			} else if trail == "Admin / Logs" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Logs</li>`
			} else if trail == "Admin / Config" {
				html += `<li class="breadcrumb-item"><a href="/admin">Admin</a></li>
						<li class="breadcrumb-item active" aria-current="page">Config</li>`
			} else if trail == "Create Badge" {
				html += `<li class="breadcrumb-item active" aria-current="page">Create Badge</li>`
			} else {
				html += `<li class="breadcrumb-item active" aria-current="page">` + trail + `</li>`
			}
		}

		html += `
		</ol>
	</div>
</nav>`
	}

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
		<form hx-post="/profile/change-password" hx-target="#password-result">
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
		<form hx-post="/profile/update" hx-target="#profile-result">
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
