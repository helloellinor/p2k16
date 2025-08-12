package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	demoMode       bool
}

func NewHandler(accountRepo *models.AccountRepository, circleRepo *models.CircleRepository, badgeRepo *models.BadgeRepository, toolRepo *models.ToolRepository, eventRepo *models.EventRepository, membershipRepo *models.MembershipRepository) *Handler {
	return &Handler{
		accountRepo:    accountRepo,
		circleRepo:     circleRepo,
		badgeRepo:      badgeRepo,
		toolRepo:       toolRepo,
		eventRepo:      eventRepo,
		membershipRepo: membershipRepo,
		demoMode:       false,
	}
}

// SetDemoMode sets whether the handler is running in demo mode
func (h *Handler) SetDemoMode(demoMode bool) {
	h.demoMode = demoMode
}

// GetAccountRepo returns the account repository (may be nil in demo mode)
func (h *Handler) GetAccountRepo() *models.AccountRepository {
	return h.accountRepo
}

// Home renders the front page
func (h *Handler) Home(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	userInfo := ""

	if user != nil {
		userInfo = `
			<div class="alert alert-info">
				Welcome back, <strong>` + user.Username + `</strong>! 
				<a href="/logout" class="btn btn-sm btn-outline-secondary ms-2">Logout</a>
			</div>`
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>P2K16 - Hackerspace Management System</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        ` + userInfo + `
        <h1>Welcome to P2K16</h1>
        <p class="lead">Hackerspace Management System</p>
        
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Quick Actions</h5>
                    </div>
                    <div class="card-body">`

	if user == nil {
		html += `
                        <a href="/login" class="btn btn-primary">Login</a>`
	} else {
		html += `
                        <a href="/dashboard" class="btn btn-primary">Dashboard</a>
                        <a href="/profile" class="btn btn-secondary">Profile</a>`
	}

	html += `
                        <button class="btn btn-info" hx-get="/api/members/active" hx-target="#member-list">Show Active Members</button>
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>System Status</h5>
                    </div>
                    <div class="card-body">`

	if h.demoMode {
		html += `
                        <span class="badge bg-warning">Demo Mode</span>
                        <p class="mt-2">Running without database connection</p>
                        <small class="text-muted">Use username "demo", "super", or "foo" with any password to login</small>`
	} else {
		html += `
                        <span class="badge bg-success">Online</span>
                        <p class="mt-2">Database connected - all systems operational</p>`
	}

	html += `
                    </div>
                </div>
            </div>
        </div>
        
        <div class="mt-4">
            <h4>Active Members</h4>
            <div id="member-list">
                <p class="text-muted">Click "Show Active Members" to load...</p>
            </div>
        </div>
    </div>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	// If already logged in, redirect to home
	if middleware.IsAuthenticated(c) {
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
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-5">
        <div class="row justify-content-center">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h4>Login to P2K16</h4>
                    </div>
                    <div class="card-body">`

	if h.demoMode {
		html += `
                        <div class="alert alert-info">
                            <strong>Demo Mode:</strong> Use username "demo", "super", or "foo" with any password
                        </div>`
	}

	html += `
                        <form hx-post="/api/auth/login" hx-target="#login-result">
                            <div class="mb-3">
                                <label for="username" class="form-label">Username</label>`

	if h.demoMode {
		html += `
                                <input type="text" class="form-control" id="username" name="username" value="demo" required>`
	} else {
		html += `
                                <input type="text" class="form-control" id="username" name="username" required>`
	}

	html += `
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">Password</label>`

	if h.demoMode {
		html += `
                                <input type="password" class="form-control" id="password" name="password" value="password" required>`
	} else {
		html += `
                                <input type="password" class="form-control" id="password" name="password" required>`
	}

	html += `
                            </div>
                            <button type="submit" class="btn btn-primary">Login</button>
                            <a href="/" class="btn btn-secondary">Back to Home</a>
                        </form>
                        <div id="login-result" class="mt-3"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	if err := middleware.LogoutUser(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// Dashboard shows the user dashboard (requires authentication)
func (h *Handler) Dashboard(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Dashboard - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <div class="d-flex justify-content-between align-items-center">
            <h1>Dashboard</h1>
            <div>
                <span class="me-3">Welcome, <strong>` + user.Username + `</strong></span>
                <a href="/logout" class="btn btn-outline-secondary">Logout</a>
            </div>
        </div>
        
        <div class="row mt-4">
            <div class="col-md-3">
                <div class="card">
                    <div class="card-header">
                        <h5>Your Badges</h5>
                    </div>
                    <div class="card-body">
                        <div id="user-badges">
                            <button class="btn btn-primary" hx-get="/api/user/badges" hx-target="#user-badges">Load Badges</button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-header">
                        <h5>Tool Management</h5>
                    </div>
                    <div class="card-body">
                        <div id="tool-section">
                            <button class="btn btn-success mb-2" hx-get="/api/tools" hx-target="#tool-section">Browse Tools</button>
                            <button class="btn btn-warning" hx-get="/api/tools/checkouts" hx-target="#tool-section">Active Checkouts</button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-header">
                        <h5>Membership</h5>
                    </div>
                    <div class="card-body">
                        <div id="membership-section">
                            <button class="btn btn-info mb-2" hx-get="/api/membership/status" hx-target="#membership-section">My Status</button>
                            <button class="btn btn-secondary" hx-get="/api/membership/active" hx-target="#membership-section">Active Members</button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-header">
                        <h5>Quick Actions</h5>
                    </div>
                    <div class="card-body">
                        <a href="/profile" class="btn btn-primary d-block mb-2">Edit Profile</a>
                        <a href="/" class="btn btn-secondary d-block">Back to Home</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveMembers returns a list of active members (for HTMX)
func (h *Handler) GetActiveMembers(c *gin.Context) {
	// This is a placeholder - in real implementation we'd fetch from database
	html := `
		<div class="list-group">
			<div class="list-group-item">
				<h6 class="mb-1">Super Admin</h6>
				<p class="mb-1">System Administrator</p>
				<small>Last active: 2 hours ago</small>
			</div>
			<div class="list-group-item">
				<h6 class="mb-1">Foo User</h6>
				<p class="mb-1">Regular Member</p>
				<small>Last active: 1 day ago</small>
			</div>
		</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetUserBadges returns user badges (for HTMX, requires authentication)
func (h *Handler) GetUserBadges(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	
	badges, err := h.badgeRepo.GetBadgesForAccount(user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", 
			[]byte(`<div class="alert alert-danger">Failed to load badges</div>`))
		return
	}
	
	if len(badges) == 0 {
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
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Username and password are required</div>`))
		return
	}

	// Handle demo mode authentication
	if h.demoMode || h.accountRepo == nil {
		if username == "demo" || username == "super" || username == "foo" {
			// Create a demo account for session
			account := &models.Account{
				ID:       1,
				Username: username,
				Email:    username + "@demo.local",
			}

			// Login user by setting session
			if err := middleware.LoginUser(c, account); err != nil {
				c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Failed to login. Please try again.</div>`))
				return
			}

			// Successful login - redirect via HTMX
			html := `
				<div class="alert alert-success">
					Login successful! Welcome to demo mode, ` + account.Username + `
				</div>
				<script>
					setTimeout(function() {
						window.location.href = '/dashboard';
					}, 1000);
				</script>`

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}

		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Invalid username or password. In demo mode, use "demo", "super", or "foo" with any password.</div>`))
		return
	}

	// Normal database authentication
	account, err := h.accountRepo.FindByUsername(username)
	if err != nil {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Invalid username or password</div>`))
		return
	}

	if !account.ValidatePassword(password) {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Invalid username or password</div>`))
		return
	}

	// Login user by setting session
	if err := middleware.LoginUser(c, account); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Failed to login. Please try again.</div>`))
		return
	}

	// Successful login - redirect via HTMX
	html := `
		<div class="alert alert-success">
			Login successful! Welcome back, ` + account.Username + `
		</div>
		<script>
			setTimeout(function() {
				window.location.href = '/dashboard';
			}, 1000);
		</script>`

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
			"<p class=\"card-text\">Type: " + tool.Type + "</p>" +
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
				"<h6 class=\"mb-1\">" + checkout.Tool.Name + " (" + checkout.Tool.Type + ")</h6>" +
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
