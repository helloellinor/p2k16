package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
)

type Handler struct {
	accountRepo *models.AccountRepository
	circleRepo  *models.CircleRepository
	badgeRepo   *models.BadgeRepository
}

func NewHandler(accountRepo *models.AccountRepository, circleRepo *models.CircleRepository, badgeRepo *models.BadgeRepository) *Handler {
	return &Handler{
		accountRepo: accountRepo,
		circleRepo:  circleRepo,
		badgeRepo:   badgeRepo,
	}
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
                    <div class="card-body">
                        <span class="badge bg-success">Online</span>
                        <p class="mt-2">All systems operational</p>
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
                    <div class="card-body">
                        <form hx-post="/api/auth/login" hx-target="#login-result">
                            <div class="mb-3">
                                <label for="username" class="form-label">Username</label>
                                <input type="text" class="form-control" id="username" name="username" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">Password</label>
                                <input type="password" class="form-control" id="password" name="password" required>
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
            <div class="col-md-4">
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
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        <h5>Recent Activity</h5>
                    </div>
                    <div class="card-body">
                        <p class="text-muted">No recent activity</p>
                    </div>
                </div>
            </div>
            <div class="col-md-4">
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
			[]byte(`<div class="alert alert-danger">Failed to award badge</div>`))
		return
	}

	html := `
<div class="alert alert-success">
Badge "` + badgeTitle + `" awarded! 
<button class="btn btn-sm btn-primary ms-2" hx-get="/api/user/badges" hx-target="#user-badges">
Refresh Badges
</button>
</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
