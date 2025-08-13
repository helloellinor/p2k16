package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
)

type Handler struct {
	accountRepo    *models.AccountRepository
	circleRepo     *models.CircleRepository
	badgeRepo      *models.BadgeRepository
	toolRepo       *models.ToolRepository
	eventRepo      *models.EventRepository
	membershipRepo *models.MembershipRepository
}

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

func NewHandler(accountRepo *models.AccountRepository, circleRepo *models.CircleRepository, badgeRepo *models.BadgeRepository, toolRepo *models.ToolRepository, eventRepo *models.EventRepository, membershipRepo *models.MembershipRepository) *Handler {
	return &Handler{
		accountRepo:    accountRepo,
		circleRepo:     circleRepo,
		badgeRepo:      badgeRepo,
		toolRepo:       toolRepo,
		eventRepo:      eventRepo,
		membershipRepo: membershipRepo,
	}
}

// GetAccountRepo returns the account repository
func (h *Handler) GetAccountRepo() *models.AccountRepository {
	return h.accountRepo
}

// Home renders the front page
func (h *Handler) Home(c *gin.Context) {
	logging.LogHandlerAction("PAGE REQUEST", "Home page visited")
	user := middleware.GetCurrentUser(c)

	if user != nil {
		logging.LogHandlerAction("USER STATUS", fmt.Sprintf("Authenticated user: %s", user.Username))
	} else {
		logging.LogHandlerAction("USER STATUS", "Anonymous user")
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<title>P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbar(c) + `
    
	<main>
		<div>
			<h1>Welcome to P2K16</h1>
			<p>Hackerspace Management System</p>
		</div>
`
	if user == nil {
		html += `
		<section>
			<header>
				<h2>Login</h2>
			</header>
			<div>
				<form hx-post="/api/auth/login" hx-target="#login-result" method="post" action="/api/auth/login">
					<div>
						<label for="username">Username</label>
						<input type="text" id="username" name="username" required>
					</div>
					<div>
						<label for="password">Password</label>
						<input type="password" id="password" name="password" required>
					</div>
					<div>
						<button type="submit">Login</button>
					</div>
				</form>
				<div id="login-result"></div>
			</div>
		</section>
		`
	}
	html += `
		<section>
			<article>
				<header>
					<h2>System Status</h2>
				</header>
				<div>
					<p>Online</p>
					<p>Database connected - all systems operational</p>
				</div>
			</article>
		</section>

	</main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	logging.LogHandlerAction("PAGE REQUEST", "Login page visited")
	// If already logged in, redirect to home
	if middleware.IsAuthenticated(c) {
		logging.LogHandlerAction("LOGIN REDIRECT", "User already authenticated, redirecting to home")
		c.Redirect(http.StatusFound, "/")
		return
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Login - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
		` + h.renderNavbar(c) + `
    
	<main>
		<section>
			<header>
				<h2>Login to P2K16</h2>
			</header>
			<div>`

	html += `
						<form hx-post="/api/auth/login" hx-target="#login-result" method="post" action="/api/auth/login">
							<div>
								<label for="username">Username</label>
								<input type="text" id="username" name="username" required>
							</div>
							<div>
								<label for="password">Password</label>
								<input type="password" id="password" name="password" required>
							</div>
							<div>
								<button type="submit">Login</button>
							</div>
						</form>
						<div id="login-result"></div>
                    </div>
                
		</section>
    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	logging.LogHandlerAction("USER ACTION", "User logout requested")
	if err := middleware.LogoutUser(c); err != nil {
		logging.LogError("LOGOUT ERROR", fmt.Sprintf("Failed to logout user: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	logging.LogSuccess("USER ACTION", "User successfully logged out - redirecting to home")
	c.Redirect(http.StatusFound, "/")
}

// Active members listing is intentionally not exposed on the landing page.

// GetUserBadges returns user badges (for HTMX, requires authentication)
func (h *Handler) GetUserBadges(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "User badges requested")
	user := middleware.GetCurrentUser(c)
	// Render centralized section
	html := h.renderUserBadgesSectionHTML(user.ID)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ProfileCardFront returns the front of the membership card (requires auth)
func (h *Handler) ProfileCardFront(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	html := h.renderProfileCardFrontHTML(user)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ProfileCardBack returns the back (editing) of the membership card (requires auth)
func (h *Handler) ProfileCardBack(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	html := h.renderProfileCardBackHTML(user)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AuthLogin handles login form submission
func (h *Handler) AuthLogin(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "Login attempt received")
	username := c.PostForm("username")
	password := c.PostForm("password")

	logging.LogHandlerAction("LOGIN ATTEMPT", fmt.Sprintf("Username: %s", username))

	if username == "" || password == "" {
		logging.LogError("LOGIN FAILED", "Missing username or password")
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Username and password are required</p>`))
		return
	}

	// Database authentication
	account, err := h.accountRepo.FindByUsername(username)
	if err != nil {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("User '%s' not found in database: %v", username, err))
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<p>Invalid username or password</p>`))
		return
	}

	logging.LogHandlerAction("USER FOUND", fmt.Sprintf("User '%s' found in database, validating password", username))
	if !account.ValidatePassword(password) {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("Invalid password for user '%s'", username))
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<p>Invalid username or password</p>`))
		return
	}

	logging.LogSuccess("LOGIN SUCCESS", fmt.Sprintf("User authenticated: %s", username))

	// Login user by setting session
	if err := middleware.LoginUser(c, account); err != nil {
		logging.LogError("SESSION ERROR", fmt.Sprintf("Failed to create session for user %s: %v", username, err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to login. Please try again.</p>`))
		return
	}

	// Successful login - redirect via HTMX
	html := `
		<section aria-live="polite">
			<p>Login successful! Welcome, ` + account.Username + `</p>
		</section>
		<script>
			window.location.href = '/';
		</script>`

	logging.LogSuccess("SESSION CREATED", fmt.Sprintf("Session created for user: %s", username))
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetAvailableBadges returns a list of available badge descriptions
func (h *Handler) GetAvailableBadges(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	descriptions, err := h.badgeRepo.GetAllDescriptions()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to load available badges</p>`))
		return
	}

	html := `
<section aria-labelledby="available-badges-title">
    <h2 id="available-badges-title">Explore Badges</h2>
    <ul>`

	for _, desc := range descriptions {
		// Check if user already has this badge
		has, _ := h.badgeRepo.AccountHasBadge(user.ID, desc.ID)
		if has {
			html += `
		<li>
			<span>` + desc.Title + `</span>
			<span>(already added)</span>
		</li>`
		} else {
			html += `
		<li>
			<span>` + desc.Title + `</span>
			<button 
				hx-post="/api/badges/award" 
				hx-vals='{"badge_title":"` + desc.Title + `"}'
				hx-target="#badge-feedback"
				hx-swap="innerHTML">Add to My Badges</button>
		</li>`
		}
	}

	html += `
	</ul>
	<div id="badge-feedback" aria-live="polite"></div>
	<div>
		<p>Want something new? Use the dedicated page to create a badge.</p>
		<p><a href="/badges/new">Create a new badge</a></p>
	</div>
</section>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CreateBadge creates a new badge description
func (h *Handler) CreateBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	title := c.PostForm("title")

	if title == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge title is required</p>`))
		return
	}

	// Check if badge already exists
	existing, _ := h.badgeRepo.FindBadgeDescriptionByTitle(title)
	if existing != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge "`+title+`" already exists</p>`))
		return
	}

	// Create badge description
	desc, err := h.badgeRepo.CreateBadgeDescription(title, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to create badge</p>`))
		return
	}

	// Award to self
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Badge created but failed to award</p>`))
		return
	}

	// Build success feedback and update the user badges via OOB swap
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := `
<section aria-live="polite">
	<p>Badge "` + title + `" created and added to your badges.</p>
</section>
` + `<div id="user-badges" hx-swap-oob="true">` + updated + `</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AwardBadge awards an existing badge to the current user
func (h *Handler) AwardBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	badgeTitle := c.PostForm("badge_title")

	if badgeTitle == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge title is required</p>`))
		return
	}

	// Find badge description
	desc, err := h.badgeRepo.FindBadgeDescriptionByTitle(badgeTitle)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte(`<p>Badge "`+badgeTitle+`" not found</p>`))
		return
	}

	// Prevent duplicates
	has, _ := h.badgeRepo.AccountHasBadge(user.ID, desc.ID)
	if has {
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<section aria-live="polite"><p>You already have '`+badgeTitle+`'.</p></section>`))
		return
	}

	// Award badge
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to award badge</p>"))
		return
	}

	// Return feedback and out-of-band update of the user badges section
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := "<section aria-live=\"polite\"><p>Added '" + badgeTitle + "' to your badges.</p></section>" +
		"<div id=\"user-badges\" hx-swap-oob=\"true\">" + updated + "</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RemoveBadge removes an awarded badge from the current user
func (h *Handler) RemoveBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	idStr := c.PostForm("account_badge_id")
	accountBadgeID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Invalid badge id</p>`))
		return
	}
	if err := h.badgeRepo.DeleteAccountBadge(accountBadgeID, user.ID); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to remove badge</p>`))
		return
	}
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := `<section aria-live="polite"><p>Badge removed.</p></section>` +
		`<div id="user-badges" hx-swap-oob="true">` + updated + `</div>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetTools returns a list of all tools
func (h *Handler) GetTools(c *gin.Context) {
	tools, err := h.toolRepo.GetAllTools()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to load tools</p>"))
		return
	}

	html := "<section aria-labelledby=\"tools-title\">" +
		"<h2 id=\"tools-title\">Available Tools</h2>" +
		"<div>"

	for _, tool := range tools {
		html += "<div class=\"col-md-6 mb-3\">" +
			"<div class=\"card border-primary\">" +
			"<div class=\"card-body\">" +
			"<h6 class=\"card-title\">" + tool.Name + "</h6>" +
			"<p class=\"card-text\">Description: " + tool.Description.String + "</p>" +
			"<button class=\"btn btn-success btn-sm\" " +
			"hx-post=\"/api/tools/checkout\" " +
			"hx-vals='{\"tool_id\":\"" + strconv.Itoa(tool.ID) + "\"}' " +
			"hx-target=\"#tool-result\" " +
			"hx-swap=\"innerHTML\">" +
			"Checkout" +
			"</button>" +
			"</div>" +
			"</div>" +
			"</div>"
	}

	html += "</div>" +
		"<div id=\"tool-result\"></div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveCheckouts returns currently checked out tools
func (h *Handler) GetActiveCheckouts(c *gin.Context) {
	checkouts, err := h.toolRepo.GetActiveCheckouts()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to load active checkouts</p>"))
		return
	}

	html := "<section aria-labelledby=\"checkouts-title\">" +
		"<h2 id=\"checkouts-title\">Currently Checked Out Tools</h2>" +
		"<div>"

	if len(checkouts) == 0 {
		html += "<p>No tools currently checked out.</p>"
	} else {
		html += "<ul>"
		for _, checkout := range checkouts {
			html += "<li>" +
				"<div>" + checkout.Tool.Name + " (" + checkout.Tool.Description.String + ") - Checked out by: " + checkout.Account.Username + " - Since: " + checkout.CheckoutAt.Format("2006-01-02 15:04") + "</div>" +
				"<button " +
				"hx-post=\"/api/tools/checkin\" " +
				"hx-vals='{\"checkout_id\":\"" + strconv.Itoa(checkout.ID) + "\"}' " +
				"hx-target=\"#tool-result\" " +
				"hx-swap=\"innerHTML\">" +
				"Check In" +
				"</button>" +
				"</li>"
		}
		html += "</ul>"
	}

	html += "<div id=\"tool-result\"></div>" +
		"</div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckoutTool handles tool checkout
func (h *Handler) CheckoutTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	toolIDStr := c.PostForm("tool_id")

	toolID, err := strconv.Atoi(toolIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<p>Invalid tool ID</p>"))
		return
	}

	// Check if tool exists
	tool, err := h.toolRepo.FindToolByID(toolID)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte("<p>Tool not found</p>"))
		return
	}

	// Create checkout record
	_, err = h.toolRepo.CheckoutTool(toolID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to checkout tool</p>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkout", user.ID)

	html := "<section aria-live=\"polite\">" +
		"<p>Successfully checked out \"" + tool.Name + "\"!</p>" +
		"<button hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckinTool handles tool checkin
func (h *Handler) CheckinTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	checkoutIDStr := c.PostForm("checkout_id")

	checkoutID, err := strconv.Atoi(checkoutIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<p>Invalid checkout ID</p>"))
		return
	}

	// Check in tool
	err = h.toolRepo.CheckinTool(checkoutID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to check in tool: "+err.Error()+"</p>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkin", user.ID)

	html := "<section aria-live=\"polite\">" +
		"<p>Tool checked in successfully!</p>" +
		"<button hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetMembershipStatus returns the membership status for current user
func (h *Handler) GetMembershipStatus(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	// Check if user is an active member
	isActive, err := h.membershipRepo.IsActiveMember(user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to check membership status</div>"))
		return
	}

	// Check payment status
	isPaying, _ := h.membershipRepo.IsAccountPayingMember(user.ID)
	isEmployee, _ := h.membershipRepo.IsAccountCompanyEmployee(user.ID)

	// Get membership details
	membership, _ := h.membershipRepo.GetMembershipByAccount(user.ID)

	html := "<section aria-labelledby=\"membership-title\">" +
		"<h2 id=\"membership-title\">Membership Status</h2>" +
		"<div>"

	if isActive {
		html += "<p>" +
			"Active Member"
		if isPaying {
			html += " (Paying Member)"
		}
		if isEmployee {
			html += " (Company Employee)"
		}
		html += "</p>"
	} else {
		html += "<p>Inactive Member</p>"
	}

	if membership != nil {
		html += "<section>" +
			"<h3>Membership Details</h3>" +
			"<p>Member since: " + membership.FirstMembership.Format("2006-01-02") + "</p>" +
			"<p>Current membership start: " + membership.StartMembership.Format("2006-01-02") + "</p>" +
			"<p>Monthly fee: " + fmt.Sprintf("%.2f NOK", float64(membership.Fee)/100) + "</p>"
		if membership.MembershipNumber.Valid {
			html += "<p><strong>Membership number:</strong> " + fmt.Sprintf("%d", membership.MembershipNumber.Int64) + "</p>"
		}
		html += "</section>"
	}

	html += "</div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveMembers returns a list of active members
func (h *Handler) GetActiveMembersDetailed(c *gin.Context) {
	payingMembers, err := h.membershipRepo.GetActivePayingMembers()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load active members</div>"))
		return
	}

	activeCompanies, err := h.membershipRepo.GetActiveCompanies()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load companies</div>"))
		return
	}

	html := "<div class=\"card\">" +
		"<div class=\"card-header\">" +
		"<h6>Active Members</h6>" +
		"</div>" +
		"<div class=\"card-body\">"

	// Show paying members
	if len(payingMembers) > 0 {
		html += "<h6 class=\"text-success\">Paying Members (" + fmt.Sprintf("%d", len(payingMembers)) + ")</h6>" +
			"<div class=\"row\">"
		for _, member := range payingMembers {
			displayName := member.Username
			if member.Name.Valid && member.Name.String != "" {
				displayName = member.Name.String + " (" + member.Username + ")"
			}
			html += "<div class=\"col-md-6 mb-2\">" +
				"<span class=\"badge bg-success\">" + displayName + "</span>" +
				"</div>"
		}
		html += "</div>"
	}

	// Show companies
	if len(activeCompanies) > 0 {
		html += "<h6 class=\"text-primary mt-3\">Active Companies (" + fmt.Sprintf("%d", len(activeCompanies)) + ")</h6>" +
			"<div class=\"list-group\">"
		for _, company := range activeCompanies {
			contactName := company.Contact.Username
			if company.Contact.Name.Valid && company.Contact.Name.String != "" {
				contactName = company.Contact.Name.String
			}
			html += "<div class=\"list-group-item\">" +
				"<h6 class=\"mb-1\">" + company.Name + "</h6>" +
				"<p class=\"mb-1\">Contact: " + contactName + "</p>" +
				"</div>"
		}
		html += "</div>"
	}

	if len(payingMembers) == 0 && len(activeCompanies) == 0 {
		html += "<p class=\"text-muted\">No active members found.</p>"
	}

	html += "</div>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Profile shows the user profile page (requires authentication)
func (h *Handler) Profile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Profile - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Profile") + `
    
	<main>
		<div>
			<h1>Membership Card</h1>
			<p>Front shows your info and badges. Click Edit to flip to the back and manage details.</p>
		</div>
		<section>
			<div id="membership-card">` + h.renderProfileCardFrontHTML(user) + `</div>
		</section>


    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Admin shows the admin interface (requires authentication)
func (h *Handler) Admin(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Admin - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin") + `
    
	<main>
		<div>
			<h1>Admin Console</h1>
			<p>Workspace for privileged tasks: manage users, tools, companies, circles, logs, and configuration.</p>
		</div>
		<section>
			<nav aria-label="Admin sections">
				<ul>
					<li><a href="/admin/users">Users</a></li>
					<li><a href="/admin/tools">Tools</a></li>
					<li><a href="/admin/companies">Companies</a></li>
					<li><a href="/admin/circles">Circles</a></li>
					<li><a href="/admin/logs">Logs</a></li>
					<li><a href="/admin/config">Config</a></li>
				</ul>
			</nav>
		</section>

    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Admin subpages — minimal pages to "paginate" admin options
func (h *Handler) AdminUsers(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Users - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Users") + `
	<main>
		<h1>Users</h1>
		<p>Manage user accounts and permissions.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) AdminTools(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Tools - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Tools") + `
	<main>
		<h1>Tools</h1>
		<p>Manage tools and equipment.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) AdminCompanies(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Companies - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Companies") + `
	<main>
		<h1>Companies</h1>
		<p>Manage corporate memberships.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) AdminCircles(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Circles - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Circles") + `
	<main>
		<h1>Circles</h1>
		<p>Manage user groups and permissions.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) AdminLogs(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Logs - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Logs") + `
	<main>
		<h1>Logs</h1>
		<p>View system logs.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (h *Handler) AdminConfig(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Config - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Config") + `
	<main>
		<h1>Config</h1>
		<p>Configure system parameters.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// BadgeCreatePage shows a dedicated page to create a new badge (requires authentication)
func (h *Handler) BadgeCreatePage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Create Badge - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Create Badge") + `
    
	<main>
		<section aria-labelledby="create-badge-title">
			<header>
				<h1 id="create-badge-title">Create a New Badge</h1>
			</header>
			<div>
				<form hx-post="/api/badges/create" hx-target="#create-result">
					<div>
						<label for="title">Title</label>
						<input id="title" type="text" name="title" required>
						<div>Choose a short, clear name (e.g., Laser Cutter Trained)</div>
					</div>
					<div>
						<button type="submit">Create Badge</button>
					</div>
				</form>
				<div id="create-result" aria-live="polite"></div>
			</div>
		</section>

		<section>
			<p>After creating, you can award it to yourself or others from the Admin Console.</p>
		</section>
	</main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
