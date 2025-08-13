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
	userInfo := ""

	if user != nil {
		logging.LogHandlerAction("USER STATUS", fmt.Sprintf("Authenticated user: %s", user.Username))
		userInfo = `
			<div class="p2k16-header__user">
				<span class="p2k16-text--secondary">Welcome, <strong>` + user.Username + `</strong></span>
				<a href="/logout" class="p2k16-button p2k16-button--secondary p2k16-button--sm">Logout</a>
			</div>`
	} else {
		logging.LogHandlerAction("USER STATUS", "Anonymous user")
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>P2K16 - Hackerspace Management System</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="/styles/p2k16-design-system.css" rel="stylesheet">
</head>
<body>
    <header class="p2k16-header">
        <div class="p2k16-container p2k16-header__container">
            <a href="/" class="p2k16-header__brand">P2K16</a>
            <nav class="p2k16-header__nav">` + userInfo + `</nav>
        </div>
    </header>
    
    <main class="p2k16-container p2k16-mt-8">
        <div class="p2k16-text--center p2k16-mb-8">
            <h1>Welcome to P2K16</h1>
            <p class="p2k16-text--secondary">Hackerspace Management System</p>
        </div>
        
        <div class="p2k16-grid p2k16-grid--2-col">
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Quick Actions</h5>
                </div>
                <div class="p2k16-card__body">`

	if user == nil {
		html += `
                        <a href="/login" class="p2k16-button p2k16-button--primary">Login</a>`
	} else {
		html += `
                        <a href="/dashboard" class="p2k16-button p2k16-button--primary">Dashboard</a>
                        <a href="/profile" class="p2k16-button p2k16-button--secondary p2k16-mt-4">Profile</a>`
	}

	html += `
                        <button class="p2k16-button p2k16-button--secondary p2k16-mt-4" 
                                hx-get="/api/members/active" 
                                hx-target="#member-list">Show Active Members</button>
                    </div>
                </div>
            </div>
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">System Status</h5>
                </div>
                <div class="p2k16-card__body">
                        <div class="p2k16-badge p2k16-badge--success">Online</div>
                        <p class="p2k16-mt-4">Database connected - all systems operational</p>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="p2k16-mt-8">
            <h4>Active Members</h4>
            <div id="member-list" class="p2k16-mt-4">
                <p class="p2k16-text--muted">Click "Show Active Members" to load...</p>
            </div>
        </div>
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
    <link href="/styles/p2k16-design-system.css" rel="stylesheet">
</head>
<body>
    <header class="p2k16-header">
        <div class="p2k16-container p2k16-header__container">
            <a href="/" class="p2k16-header__brand">P2K16</a>
        </div>
    </header>
    
    <main class="p2k16-container p2k16-container--narrow p2k16-mt-8">
        <div class="p2k16-card">
            <div class="p2k16-card__header">
                <h4 class="p2k16-card__title p2k16-text--center">Login to P2K16</h4>
            </div>
            <div class="p2k16-card__body">`

	html += `
                        <form class="p2k16-form" hx-post="/api/auth/login" hx-target="#login-result" method="post" action="/api/auth/login">
                            <div class="p2k16-field">
                                <label for="username" class="p2k16-field__label">Username</label>
                                <input type="text" class="p2k16-field__input" id="username" name="username" required>
                            </div>
                            <div class="p2k16-field">
                                <label for="password" class="p2k16-field__label">Password</label>
                                <input type="password" class="p2k16-field__input" id="password" name="password" required>
                            </div>
                            <div class="p2k16-flex p2k16-flex--between">
                                <button type="submit" class="p2k16-button p2k16-button--primary">Login</button>
                                <a href="/" class="p2k16-button p2k16-button--secondary">Back to Home</a>
                            </div>
                        </form>
                        <div id="login-result" class="p2k16-mt-6"></div>
                    </div>
                </div>
            </div>
        </div>
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

// Dashboard shows the user dashboard (requires authentication)
func (h *Handler) Dashboard(c *gin.Context) {
	logging.LogHandlerAction("PAGE REQUEST", "Dashboard page visited")
	user := middleware.GetCurrentUser(c)
	if user != nil {
		logging.LogHandlerAction("USER ACCESS", fmt.Sprintf("Dashboard accessed by user: %s", user.Username))
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Dashboard - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="/styles/p2k16-design-system.css" rel="stylesheet">
</head>
<body>
    <header class="p2k16-header">
        <div class="p2k16-container p2k16-header__container">
            <a href="/" class="p2k16-header__brand">P2K16</a>
            <nav class="p2k16-header__nav">
                <div class="p2k16-header__user">
                    <span class="p2k16-text--secondary">Welcome, <strong>` + user.Username + `</strong></span>
                    <a href="/logout" class="p2k16-button p2k16-button--secondary p2k16-button--sm">Logout</a>
                </div>
            </nav>
        </div>
    </header>
    
    <main class="p2k16-container p2k16-mt-8">
        <div class="p2k16-mb-8">
            <h1>Dashboard</h1>
        </div>
        
        <div class="p2k16-grid p2k16-grid--4-col">
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Your Badges</h5>
                </div>
                <div class="p2k16-card__body">
                    <div id="user-badges">
                        <button class="p2k16-button p2k16-button--primary p2k16-button--full" 
                                hx-get="/api/user/badges" 
                                hx-target="#user-badges">Load Badges</button>
                    </div>
                </div>
            </div>
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Tool Management</h5>
                </div>
                <div class="p2k16-card__body">
                    <div id="tool-section">
                        <button class="p2k16-button p2k16-button--success p2k16-button--full p2k16-mb-4" 
                                hx-get="/api/tools" 
                                hx-target="#tool-section">Browse Tools</button>
                        <button class="p2k16-button p2k16-button--warning p2k16-button--full" 
                                hx-get="/api/tools/checkouts" 
                                hx-target="#tool-section">Active Checkouts</button>
                    </div>
                </div>
            </div>
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Membership</h5>
                </div>
                <div class="p2k16-card__body">
                    <div id="membership-section">
                        <button class="p2k16-button p2k16-button--secondary p2k16-button--full p2k16-mb-4" 
                                hx-get="/api/membership/status" 
                                hx-target="#membership-section">My Status</button>
                        <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                                hx-get="/api/membership/active" 
                                hx-target="#membership-section">Active Members</button>
                    </div>
                </div>
            </div>
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Quick Actions</h5>
                </div>
                <div class="p2k16-card__body">
                    <a href="/profile" class="p2k16-button p2k16-button--primary p2k16-button--full p2k16-mb-4">Edit Profile</a>
                    <a href="/admin" class="p2k16-button p2k16-button--warning p2k16-button--full p2k16-mb-4">Administration</a>
                    <a href="/" class="p2k16-button p2k16-button--secondary p2k16-button--full">Back to Home</a>
                </div>
            </div>
        </div>
    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveMembers returns a list of active members (for HTMX)
func (h *Handler) GetActiveMembers(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "Active members list requested")
	// This is a placeholder - in real implementation we'd fetch from database
	logging.LogWarning("UNIMPLEMENTED", "GetActiveMembers not yet connected to database - returning mock data")
	html := `
		<div class="p2k16-grid p2k16-grid--2-col">
			<div class="p2k16-card">
				<div class="p2k16-card__body">
					<h6 class="p2k16-mb-4">Super Admin</h6>
					<p class="p2k16-text--secondary p2k16-mb-4">System Administrator</p>
					<small class="p2k16-text--muted">Last active: 2 hours ago</small>
				</div>
			</div>
			<div class="p2k16-card">
				<div class="p2k16-card__body">
					<h6 class="p2k16-mb-4">Foo User</h6>
					<p class="p2k16-text--secondary p2k16-mb-4">Regular Member</p>
					<small class="p2k16-text--muted">Last active: 1 day ago</small>
				</div>
			</div>
		</div>`

	logging.LogHandlerAction("API RESPONSE", "Active members list returned (mock data)")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetUserBadges returns user badges (for HTMX, requires authentication)
func (h *Handler) GetUserBadges(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "User badges requested")
	user := middleware.GetCurrentUser(c)
	
	badges, err := h.badgeRepo.GetBadgesForAccount(user.ID)
	if err != nil {
		logging.LogError("DATABASE ERROR", fmt.Sprintf("Failed to fetch badges for user %s: %v", user.Username, err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", 
			[]byte(`<div class="alert alert-danger">Failed to load badges</div>`))
		return
	}
	
	if len(badges) == 0 {
		logging.LogHandlerAction("BADGES RESULT", fmt.Sprintf("No badges found for user: %s", user.Username))
		html := `
			<div class="text-center">
				<p class="text-muted">No badges yet!</p>
				<button class="btn btn-primary" hx-get="/api/badges/available" hx-target="#available-badges" hx-swap="innerHTML">
					Browse Available Badges
				</button>
				<div id="available-badges" class="mt-3"></div>
			</div>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		return
	}
	
	html := `<div class="badge-list">`
	for _, badge := range badges {
		color := "primary"
		if badge.BadgeDescription.Color.Valid && badge.BadgeDescription.Color.String != "" {
			color = badge.BadgeDescription.Color.String
		}
		
		html += `<span class="badge bg-` + color + ` me-2 mb-2">` + badge.BadgeDescription.Title + `</span>`
	}
	html += `<p class="text-muted mt-2">You have ` + fmt.Sprintf("%d", len(badges)) + ` badges.</p>`
	html += `<button class="btn btn-sm btn-outline-primary" hx-get="/api/badges/available" hx-target="#available-badges" hx-swap="innerHTML">
				Browse More Badges
			</button>
			<div id="available-badges" class="mt-3"></div>
		</div>`
	
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
			[]byte(`<div class="alert alert-danger">Username and password are required</div>`))
		return
	}

	// Database authentication
	account, err := h.accountRepo.FindByUsername(username)
	if err != nil {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("User '%s' not found in database: %v", username, err))
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Invalid username or password</div>`))
		return
	}
	
	logging.LogHandlerAction("USER FOUND", fmt.Sprintf("User '%s' found in database, validating password", username))
	if !account.ValidatePassword(password) {
		logging.LogError("LOGIN FAILED", fmt.Sprintf("Invalid password for user '%s'", username))
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Invalid username or password</div>`))
		return
	}
	
	logging.LogSuccess("LOGIN SUCCESS", fmt.Sprintf("User authenticated: %s", username))

	// Login user by setting session
	if err := middleware.LoginUser(c, account); err != nil {
		logging.LogError("SESSION ERROR", fmt.Sprintf("Failed to create session for user %s: %v", username, err))
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Failed to login. Please try again.</div>`))
		return
	}

	// Successful login - redirect via HTMX
	html := `
		<div class="alert alert-success">
			Login successful! Welcome, ` + account.Username + `
		</div>
		<script>
			setTimeout(function() {
				window.location.href = '/dashboard';
			}, 1000);
		</script>`

	logging.LogSuccess("SESSION CREATED", fmt.Sprintf("Session created for user: %s", username))
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetAvailableBadges returns a list of available badge descriptions
func (h *Handler) GetAvailableBadges(c *gin.Context) {
	descriptions, err := h.badgeRepo.GetAllDescriptions()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Failed to load available badges</div>`))
		return
	}

	html := `
<div class="card">
<div class="card-header">
<h6>Available Badges</h6>
</div>
<div class="card-body">
<div class="row">`

	for _, desc := range descriptions {
		color := "outline-primary"
		if desc.Color.Valid && desc.Color.String != "" {
			color = "outline-" + desc.Color.String
		}

		html += `
<div class="col-md-6 mb-2">
<div class="d-flex justify-content-between align-items-center">
<span class="badge ` + color + `">` + desc.Title + `</span>
<button class="btn btn-sm btn-success" 
        hx-post="/api/badges/award" 
        hx-vals='{"badge_title":"` + desc.Title + `"}'
        hx-target="#badge-result"
        hx-swap="innerHTML">
Award to Self
</button>
</div>
</div>`
	}

	html += `
</div>
<div id="badge-result" class="mt-3"></div>
<div class="mt-3">
<h6>Create New Badge</h6>
<form hx-post="/api/badges/create" hx-target="#badge-result">
<div class="input-group">
<input type="text" class="form-control" name="title" placeholder="Badge title" required>
<button type="submit" class="btn btn-primary">Create & Award</button>
</div>
</form>
</div>
</div>
</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CreateBadge creates a new badge description
func (h *Handler) CreateBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	title := c.PostForm("title")

	if title == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Badge title is required</div>`))
		return
	}

	// Check if badge already exists
	existing, _ := h.badgeRepo.FindBadgeDescriptionByTitle(title)
	if existing != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-warning">Badge "`+title+`" already exists</div>`))
		return
	}

	// Create badge description
	desc, err := h.badgeRepo.CreateBadgeDescription(title, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Failed to create badge</div>`))
		return
	}

	// Award to self
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Badge created but failed to award</div>`))
		return
	}

	html := `
<div class="alert alert-success">
Badge "` + title + `" created and awarded! 
<button class="btn btn-sm btn-primary ms-2" hx-get="/api/user/badges" hx-target="#user-badges">
Refresh Badges
</button>
</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AwardBadge awards an existing badge to the current user
func (h *Handler) AwardBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	badgeTitle := c.PostForm("badge_title")

	if badgeTitle == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Badge title is required</div>`))
		return
	}

	// Find badge description
	desc, err := h.badgeRepo.FindBadgeDescriptionByTitle(badgeTitle)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Badge "`+badgeTitle+`" not found</div>`))
		return
	}

	// Award badge
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to award badge</div>"))
		return
	}

	html := "<div class=\"alert alert-success\">" +
		"Badge \"" + badgeTitle + "\" awarded! " +
		"<button class=\"btn btn-sm btn-primary ms-2\" hx-get=\"/api/user/badges\" hx-target=\"#user-badges\">" +
		"Refresh Badges" +
		"</button>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetTools returns a list of all tools
func (h *Handler) GetTools(c *gin.Context) {
	tools, err := h.toolRepo.GetAllTools()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load tools</div>"))
		return
	}

	html := "<div class=\"card\">" +
		"<div class=\"card-header\">" +
		"<h6>Available Tools</h6>" +
		"</div>" +
		"<div class=\"card-body\">" +
		"<div class=\"row\">"

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
		"<div id=\"tool-result\" class=\"mt-3\"></div>" +
		"</div>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveCheckouts returns currently checked out tools
func (h *Handler) GetActiveCheckouts(c *gin.Context) {
	checkouts, err := h.toolRepo.GetActiveCheckouts()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load active checkouts</div>"))
		return
	}

	html := "<div class=\"card\">" +
		"<div class=\"card-header\">" +
		"<h6>Currently Checked Out Tools</h6>" +
		"</div>" +
		"<div class=\"card-body\">"

	if len(checkouts) == 0 {
		html += "<p class=\"text-muted\">No tools currently checked out.</p>"
	} else {
		html += "<div class=\"list-group\">"
		for _, checkout := range checkouts {
			html += "<div class=\"list-group-item d-flex justify-content-between align-items-center\">" +
				"<div>" +
				"<h6 class=\"mb-1\">" + checkout.Tool.Name + " (" + checkout.Tool.Description.String + ")</h6>" +
				"<p class=\"mb-1\">Checked out by: " + checkout.Account.Username + "</p>" +
				"<small>Since: " + checkout.CheckoutAt.Format("2006-01-02 15:04") + "</small>" +
				"</div>" +
				"<button class=\"btn btn-warning btn-sm\" " +
				"hx-post=\"/api/tools/checkin\" " +
				"hx-vals='{\"checkout_id\":\"" + strconv.Itoa(checkout.ID) + "\"}' " +
				"hx-target=\"#tool-result\" " +
				"hx-swap=\"innerHTML\">" +
				"Check In" +
				"</button>" +
				"</div>"
		}
		html += "</div>"
	}

	html += "<div id=\"tool-result\" class=\"mt-3\"></div>" +
		"</div>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckoutTool handles tool checkout
func (h *Handler) CheckoutTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	toolIDStr := c.PostForm("tool_id")

	toolID, err := strconv.Atoi(toolIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Invalid tool ID</div>"))
		return
	}

	// Check if tool exists
	tool, err := h.toolRepo.FindToolByID(toolID)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Tool not found</div>"))
		return
	}

	// Create checkout record
	_, err = h.toolRepo.CheckoutTool(toolID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to checkout tool</div>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkout", user.ID)

	html := "<div class=\"alert alert-success\">" +
		"Successfully checked out \"" + tool.Name + "\"! " +
		"<button class=\"btn btn-sm btn-primary ms-2\" hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckinTool handles tool checkin
func (h *Handler) CheckinTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	checkoutIDStr := c.PostForm("checkout_id")

	checkoutID, err := strconv.Atoi(checkoutIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Invalid checkout ID</div>"))
		return
	}

	// Check in tool
	err = h.toolRepo.CheckinTool(checkoutID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to check in tool: "+err.Error()+"</div>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkin", user.ID)

	html := "<div class=\"alert alert-success\">" +
		"Tool checked in successfully! " +
		"<button class=\"btn btn-sm btn-primary ms-2\" hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</div>"

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

	html := "<div class=\"card\">" +
		"<div class=\"card-header\">" +
		"<h6>Membership Status</h6>" +
		"</div>" +
		"<div class=\"card-body\">"

	if isActive {
		html += "<div class=\"alert alert-success\">" +
			"<strong>Active Member</strong>"
		if isPaying {
			html += " (Paying Member)"
		}
		if isEmployee {
			html += " (Company Employee)"
		}
		html += "</div>"
	} else {
		html += "<div class=\"alert alert-warning\">" +
			"<strong>Inactive Member</strong>" +
			"</div>"
	}

	if membership != nil {
		html += "<div class=\"mt-3\">" +
			"<h6>Membership Details</h6>" +
			"<p><strong>Member since:</strong> " + membership.FirstMembership.Format("2006-01-02") + "</p>" +
			"<p><strong>Current membership start:</strong> " + membership.StartMembership.Format("2006-01-02") + "</p>" +
			"<p><strong>Monthly fee:</strong> " + fmt.Sprintf("%.2f NOK", float64(membership.Fee)/100) + "</p>"
		if membership.MembershipNumber.Valid {
			html += "<p><strong>Membership number:</strong> " + fmt.Sprintf("%d", membership.MembershipNumber.Int64) + "</p>"
		}
		html += "</div>"
	}

	html += "</div>" +
		"</div>"

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
    <link href="/styles/p2k16-design-system.css" rel="stylesheet">
</head>
<body>
    <header class="p2k16-header">
        <div class="p2k16-container p2k16-header__container">
            <a href="/" class="p2k16-header__brand">P2K16</a>
            <nav class="p2k16-header__nav">
                <div class="p2k16-header__user">
                    <span class="p2k16-text--secondary">Welcome, <strong>` + user.Username + `</strong></span>
                    <a href="/logout" class="p2k16-button p2k16-button--secondary p2k16-button--sm">Logout</a>
                </div>
            </nav>
        </div>
    </header>
    
    <main class="p2k16-container p2k16-container--narrow p2k16-mt-8">
        <div class="p2k16-mb-8">
            <h1>User Profile</h1>
            <p class="p2k16-text--secondary">Manage your account settings and view your membership status</p>
        </div>

        <!-- Change Password Section -->
        <div class="p2k16-card p2k16-mb-8">
            <div class="p2k16-card__header">
                <h5 class="p2k16-card__title">Change Password</h5>
            </div>
            <div class="p2k16-card__body">
                <form class="p2k16-form" hx-post="/api/profile/change-password" hx-target="#password-result">
                    <div class="p2k16-field">
                        <label for="oldPassword" class="p2k16-field__label">Current Password</label>
                        <input type="password" class="p2k16-field__input" id="oldPassword" name="oldPassword" required>
                    </div>
                    <div class="p2k16-field">
                        <label for="newPassword" class="p2k16-field__label">New Password</label>
                        <input type="password" class="p2k16-field__input" id="newPassword" name="newPassword" required>
                    </div>
                    <div class="p2k16-field">
                        <label for="confirmPassword" class="p2k16-field__label">Confirm New Password</label>
                        <input type="password" class="p2k16-field__input" id="confirmPassword" name="confirmPassword" required>
                    </div>
                    <button type="submit" class="p2k16-button p2k16-button--primary">Change Password</button>
                </form>
                <div id="password-result" class="p2k16-mt-6"></div>
            </div>
        </div>

        <!-- Profile Details Section -->
        <div class="p2k16-card p2k16-mb-8">
            <div class="p2k16-card__header">
                <h5 class="p2k16-card__title">Profile Details</h5>
            </div>
            <div class="p2k16-card__body">
                <form class="p2k16-form" hx-post="/api/profile/update" hx-target="#profile-result">
                    <div class="p2k16-field">
                        <label for="name" class="p2k16-field__label">Full Name</label>
                        <input type="text" class="p2k16-field__input" id="name" name="name" placeholder="Enter your full name">
                        <div class="p2k16-field__help">This will be displayed on your public profile</div>
                    </div>
                    <div class="p2k16-field">
                        <label for="email" class="p2k16-field__label">Email Address</label>
                        <input type="email" class="p2k16-field__input" id="email" name="email" value="` + user.Account.Email + `" readonly>
                        <div class="p2k16-field__help">Contact an administrator to change your email address</div>
                    </div>
                    <div class="p2k16-field">
                        <label for="phone" class="p2k16-field__label">Phone Number</label>
                        <input type="tel" class="p2k16-field__input" id="phone" name="phone" placeholder="Enter your phone number">
                        <div class="p2k16-field__help">Used for emergency contact and door access notifications</div>
                    </div>
                    <button type="submit" class="p2k16-button p2k16-button--primary">Save Changes</button>
                </form>
                <div id="profile-result" class="p2k16-mt-6"></div>
            </div>
        </div>

        <!-- Badges Section -->
        <div class="p2k16-card p2k16-mb-8">
            <div class="p2k16-card__header">
                <h5 class="p2k16-card__title">Your Badges</h5>
            </div>
            <div class="p2k16-card__body">
                <div id="profile-badges">
                    <button class="p2k16-button p2k16-button--secondary" 
                            hx-get="/api/user/badges" 
                            hx-target="#profile-badges">Load Your Badges</button>
                </div>
            </div>
        </div>

        <!-- Navigation -->
        <div class="p2k16-text--center">
            <a href="/dashboard" class="p2k16-button p2k16-button--secondary">Back to Dashboard</a>
            <a href="/" class="p2k16-button p2k16-button--secondary">Home</a>
        </div>
    </main>
</body>
</html>`

c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Admin shows the admin interface (requires authentication)
func (h *Handler) Admin(c *gin.Context) {
user := middleware.GetCurrentUser(c)

html := `
<!DOCTYPE html>
<html>
<head>
    <title>Admin - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="/styles/p2k16-design-system.css" rel="stylesheet">
</head>
<body>
    <header class="p2k16-header">
        <div class="p2k16-container p2k16-header__container">
            <a href="/" class="p2k16-header__brand">P2K16</a>
            <nav class="p2k16-header__nav">
                <div class="p2k16-header__user">
                    <span class="p2k16-text--secondary">Welcome, <strong>` + user.Username + `</strong></span>
                    <a href="/logout" class="p2k16-button p2k16-button--secondary p2k16-button--sm">Logout</a>
                </div>
            </nav>
        </div>
    </header>
    
    <main class="p2k16-container p2k16-mt-8">
        <div class="p2k16-mb-8">
            <h1>Administration</h1>
            <p class="p2k16-text--secondary">Manage users, tools, badges, and system settings</p>
        </div>

        <div class="p2k16-grid p2k16-grid--3-col">
            <!-- User Management -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">User Management</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Manage user accounts and permissions</p>
                    <button class="p2k16-button p2k16-button--primary p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/users" 
                            hx-target="#admin-content">View All Users</button>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                            hx-get="/api/admin/users/new" 
                            hx-target="#admin-content">Create New User</button>
                </div>
            </div>

            <!-- Badge Management -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Badge System</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Manage badges and certifications</p>
                    <button class="p2k16-button p2k16-button--success p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/badges" 
                            hx-target="#admin-content">Manage Badges</button>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                            hx-get="/api/admin/badges/award" 
                            hx-target="#admin-content">Award Badge</button>
                </div>
            </div>

            <!-- Tool Management -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Tool Management</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Manage tools and equipment</p>
                    <button class="p2k16-button p2k16-button--warning p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/tools" 
                            hx-target="#admin-content">Manage Tools</button>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                            hx-get="/api/admin/tools/new" 
                            hx-target="#admin-content">Add New Tool</button>
                </div>
            </div>

            <!-- Company Management -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Companies</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Manage corporate memberships</p>
                    <button class="p2k16-button p2k16-button--primary p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/companies" 
                            hx-target="#admin-content">View Companies</button>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                            hx-get="/api/admin/companies/new" 
                            hx-target="#admin-content">Add Company</button>
                </div>
            </div>

            <!-- Circle Management -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">Circles</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Manage user groups and permissions</p>
                    <button class="p2k16-button p2k16-button--primary p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/circles" 
                            hx-target="#admin-content">View Circles</button>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full" 
                            hx-get="/api/admin/circles/new" 
                            hx-target="#admin-content">Create Circle</button>
                </div>
            </div>

            <!-- System Settings -->
            <div class="p2k16-card">
                <div class="p2k16-card__header">
                    <h5 class="p2k16-card__title">System Settings</h5>
                </div>
                <div class="p2k16-card__body">
                    <p class="p2k16-text--secondary p2k16-mb-6">Configure system parameters</p>
                    <button class="p2k16-button p2k16-button--secondary p2k16-button--full p2k16-mb-4" 
                            hx-get="/api/admin/settings" 
                            hx-target="#admin-content">System Settings</button>
                    <button class="p2k16-button p2k16-button--danger p2k16-button--full" 
                            hx-get="/api/admin/logs" 
                            hx-target="#admin-content">View Logs</button>
                </div>
            </div>
        </div>

        <!-- Dynamic Content Area -->
        <div class="p2k16-mt-8">
            <div id="admin-content">
                <div class="p2k16-card">
                    <div class="p2k16-card__body">
                        <p class="p2k16-text--center p2k16-text--muted">Select an action from above to get started</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Navigation -->
        <div class="p2k16-text--center p2k16-mt-8">
            <a href="/dashboard" class="p2k16-button p2k16-button--secondary">Back to Dashboard</a>
            <a href="/" class="p2k16-button p2k16-button--secondary">Home</a>
        </div>
    </main>
</body>
</html>`

c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
