package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
)

type Handler struct {
	accountRepo *models.AccountRepository
	circleRepo  *models.CircleRepository
}

func NewHandler(accountRepo *models.AccountRepository, circleRepo *models.CircleRepository) *Handler {
	return &Handler{
		accountRepo: accountRepo,
		circleRepo:  circleRepo,
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

	html := `
		<div class="badge-list">
			<span class="badge bg-primary me-2 mb-2">Go Programming</span>
			<span class="badge bg-success me-2 mb-2">HTMX Expert</span>
			<span class="badge bg-info me-2 mb-2">Database Admin</span>
			<p class="text-muted mt-2">User ` + user.Username + ` has 3 badges.</p>
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
