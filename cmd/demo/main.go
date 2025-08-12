package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Demo mode - simple server without database for UI testing
func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Session middleware
	store := cookie.NewStore([]byte("demo-secret"))
	store.Options(sessions.Options{
		MaxAge:   86400 * 7,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
	r.Use(sessions.Sessions("p2k16-session", store))

	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")

		userInfo := ""
		if username != nil {
			userInfo = `
				<div class="alert alert-info">
					Welcome back, <strong>` + username.(string) + `</strong>! 
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
        <p class="lead">Hackerspace Management System (Demo Mode)</p>
        
        <div class="row">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Quick Actions</h5>
                    </div>
                    <div class="card-body">`

		if username == nil {
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
                        <span class="badge bg-warning">Demo Mode</span>
                        <p class="mt-2">Running without database connection</p>
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
	})

	r.GET("/login", func(c *gin.Context) {
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
                        <h4>Login to P2K16 (Demo)</h4>
                    </div>
                    <div class="card-body">
                        <div class="alert alert-info">
                            <strong>Demo Mode:</strong> Use username "demo" with any password
                        </div>
                        <form hx-post="/api/auth/login" hx-target="#login-result">
                            <div class="mb-3">
                                <label for="username" class="form-label">Username</label>
                                <input type="text" class="form-control" id="username" name="username" value="demo" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">Password</label>
                                <input type="password" class="form-control" id="password" name="password" value="password" required>
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
	})

	r.GET("/dashboard", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")

		if username == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

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
                <span class="me-3">Welcome, <strong>` + username.(string) + `</strong></span>
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
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusFound, "/")
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/members/active", func(c *gin.Context) {
			html := `
				<div class="list-group">
					<div class="list-group-item">
						<h6 class="mb-1">Demo Admin</h6>
						<p class="mb-1">System Administrator</p>
						<small>Last active: 2 hours ago</small>
					</div>
					<div class="list-group-item">
						<h6 class="mb-1">Demo User</h6>
						<p class="mb-1">Regular Member</p>
						<small>Last active: 1 day ago</small>
					</div>
				</div>`

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		})

		api.POST("/auth/login", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")

			if username == "" || password == "" {
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Username and password are required</div>`))
				return
			}

			// Demo login - accept "demo" with any password
			if username == "demo" {
				session := sessions.Default(c)
				session.Set("username", username)
				session.Save()

				html := `
					<div class="alert alert-success">
						Login successful! Welcome to demo mode, ` + username + `
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
				[]byte(`<div class="alert alert-danger">Invalid username or password. Use "demo" for demo mode.</div>`))
		})

		api.GET("/user/badges", func(c *gin.Context) {
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-warning">Please log in to view badges.</div>`))
				return
			}

			html := `
				<div class="badge-list">
					<span class="badge bg-primary me-2 mb-2">Go Programming</span>
					<span class="badge bg-success me-2 mb-2">HTMX Expert</span>
					<span class="badge bg-info me-2 mb-2">Database Admin</span>
					<span class="badge bg-warning me-2 mb-2">Demo User</span>
					<p class="text-muted mt-2">User ` + username.(string) + ` has 4 demo badges.</p>
				</div>`

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting demo server on port %s", port)
	log.Printf("Demo application available at http://localhost:%s", port)
	log.Printf("Login with username 'demo' and any password")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
