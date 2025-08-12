package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	// For now, return a simple HTML response
	// Later we'll add proper template rendering with HTMX
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
        <h1>Welcome to P2K16</h1>
        <p class="lead">Hackerspace Management System</p>
        
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Quick Actions</h5>
                    </div>
                    <div class="card-body">
                        <a href="/login" class="btn btn-primary">Login</a>
                        <button class="btn btn-secondary" hx-get="/api/members/active" hx-target="#member-list">Show Active Members</button>
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

	// Successful login - in real implementation we'd set session/token
	html := `
		<div class="alert alert-success">
			Login successful! Welcome back, ` + account.Username + `
			<br><a href="/" class="btn btn-primary mt-2">Go to Dashboard</a>
		</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
