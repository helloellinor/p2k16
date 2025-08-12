package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// logAction logs demo mode actions with timestamp and clear formatting
func logAction(action, details string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("\n🎯 [DEMO MODE] %s | %s | %s\n", timestamp, action, details)
}

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
		logAction("PAGE REQUEST", "Home page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		userInfo := ""
		if username != nil {
			logAction("USER STATUS", fmt.Sprintf("Authenticated user: %s", username.(string)))
			userInfo = `
				<div class="alert alert-info">
					Welcome back, <strong>` + username.(string) + `</strong>! 
					<a href="/logout" class="btn btn-sm btn-outline-secondary ms-2">Logout</a>
				</div>`
		} else {
			logAction("USER STATUS", "Anonymous user")
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
		logAction("PAGE REQUEST", "Login page visited")
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
		logAction("PAGE REQUEST", "Dashboard page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		if username == nil {
			logAction("ACCESS DENIED", "Dashboard access attempted without authentication")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		logAction("USER ACCESS", fmt.Sprintf("Dashboard accessed by user: %s", username.(string)))

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

	r.GET("/profile", func(c *gin.Context) {
		logAction("PAGE REQUEST", "Profile page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		if username == nil {
			logAction("ACCESS DENIED", "Profile access attempted without authentication")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		logAction("USER ACCESS", fmt.Sprintf("Profile page accessed by user: %s", username.(string)))

		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Profile - P2K16 (Demo)</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <div class="d-flex justify-content-between align-items-center">
            <h1>User Profile</h1>
            <div>
                <span class="me-3">Welcome, <strong>` + username.(string) + `</strong></span>
                <a href="/logout" class="btn btn-outline-secondary">Logout</a>
            </div>
        </div>
        
        <div class="alert alert-warning mt-3">
            <strong>Demo Mode:</strong> Changes will be simulated but not saved to database
        </div>
        
        <div class="row mt-4">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Change Password</h5>
                    </div>
                    <div class="card-body">
                        <form hx-post="/api/profile/change-password" hx-target="#password-result">
                            <div class="mb-3">
                                <label for="oldPassword" class="form-label">Current Password</label>
                                <input type="password" class="form-control" id="oldPassword" name="oldPassword" required>
                            </div>
                            <div class="mb-3">
                                <label for="newPassword" class="form-label">New Password</label>
                                <input type="password" class="form-control" id="newPassword" name="newPassword" required>
                            </div>
                            <div class="mb-3">
                                <label for="confirmPassword" class="form-label">Confirm New Password</label>
                                <input type="password" class="form-control" id="confirmPassword" name="confirmPassword" required>
                            </div>
                            <button type="submit" class="btn btn-primary">Change Password</button>
                        </form>
                        <div id="password-result" class="mt-3"></div>
                    </div>
                </div>
            </div>
            
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Profile Details</h5>
                    </div>
                    <div class="card-body">
                        <form hx-post="/api/profile/update" hx-target="#profile-result">
                            <div class="mb-3">
                                <label for="name" class="form-label">Full Name</label>
                                <input type="text" class="form-control" id="name" name="name" placeholder="Enter your full name">
                                <div class="form-text">This will be displayed on your public profile</div>
                            </div>
                            <div class="mb-3">
                                <label for="email" class="form-label">Email Address</label>
                                <input type="email" class="form-control" id="email" name="email" value="` + username.(string) + `@demo.p2k16.io" readonly>
                                <div class="form-text">Contact an administrator to change your email address</div>
                            </div>
                            <div class="mb-3">
                                <label for="phone" class="form-label">Phone Number</label>
                                <input type="tel" class="form-control" id="phone" name="phone" placeholder="Enter your phone number">
                                <div class="form-text">Used for emergency contact and door access notifications</div>
                            </div>
                            <button type="submit" class="btn btn-primary">Save Changes</button>
                        </form>
                        <div id="profile-result" class="mt-3"></div>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="mt-4">
            <h4>Your Badges</h4>
            <div id="user-badges">
                <button class="btn btn-primary" hx-get="/api/user/badges" hx-target="#user-badges">Load Your Badges</button>
            </div>
        </div>
        
        <div class="mt-4 text-center">
            <a href="/dashboard" class="btn btn-secondary me-2">Back to Dashboard</a>
            <a href="/" class="btn btn-secondary">Home</a>
        </div>
    </div>
</body>
</html>`

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/logout", func(c *gin.Context) {
		logAction("USER ACTION", "User logout requested")
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		logAction("USER ACTION", "User session cleared - redirecting to home")
		c.Redirect(http.StatusFound, "/")
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/members/active", func(c *gin.Context) {
			logAction("API REQUEST", "Active members list requested")
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

			logAction("API RESPONSE", "Active members list returned (2 demo users)")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		})

		api.POST("/auth/login", func(c *gin.Context) {
			logAction("API REQUEST", "Login attempt received")
			username := c.PostForm("username")
			password := c.PostForm("password")

			logAction("LOGIN ATTEMPT", fmt.Sprintf("Username: %s", username))

			if username == "" || password == "" {
				logAction("LOGIN FAILED", "Missing username or password")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Username and password are required</div>`))
				return
			}

			// Demo login - accept "demo" with any password
			if username == "demo" {
				logAction("LOGIN SUCCESS", fmt.Sprintf("Demo user authenticated: %s", username))
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

				logAction("SESSION CREATED", fmt.Sprintf("Session created for user: %s", username))
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
				return
			}

			logAction("LOGIN FAILED", fmt.Sprintf("Invalid credentials for username: %s", username))
			c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
				[]byte(`<div class="alert alert-danger">Invalid username or password. Use "demo" for demo mode.</div>`))
		})

		api.GET("/user/badges", func(c *gin.Context) {
			logAction("API REQUEST", "User badges requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logAction("API ERROR", "Badges request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-warning">Please log in to view badges.</div>`))
				return
			}

			logAction("BADGES LOADED", fmt.Sprintf("Badges loaded for user: %s", username.(string)))
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

		// Profile management API endpoints (Phase 2)
		api.POST("/profile/change-password", func(c *gin.Context) {
			logAction("API REQUEST", "Password change requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logAction("API ERROR", "Password change request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Please log in to change password.</div>`))
				return
			}

			oldPassword := c.PostForm("oldPassword")
			newPassword := c.PostForm("newPassword")
			confirmPassword := c.PostForm("confirmPassword")

			logAction("PASSWORD CHANGE", fmt.Sprintf("User %s attempting password change", username.(string)))

			// Validation
			if oldPassword == "" || newPassword == "" || confirmPassword == "" {
				logAction("PASSWORD CHANGE FAILED", "Missing required fields")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">All fields are required</div>`))
				return
			}

			if newPassword != confirmPassword {
				logAction("PASSWORD CHANGE FAILED", "New passwords do not match")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">New passwords do not match</div>`))
				return
			}

			if len(newPassword) < 6 {
				logAction("PASSWORD CHANGE FAILED", "Password too short")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Password must be at least 6 characters long</div>`))
				return
			}

			// Demo mode - simulate password change
			logAction("PASSWORD CHANGE SUCCESS", fmt.Sprintf("Password changed successfully for user: %s (demo mode)", username.(string)))
			c.Data(http.StatusOK, "text/html; charset=utf-8",
				[]byte(`<div class="alert alert-success">Password changed successfully (demo mode)</div>`))
		})

		api.POST("/profile/update", func(c *gin.Context) {
			logAction("API REQUEST", "Profile update requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logAction("API ERROR", "Profile update request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="alert alert-danger">Please log in to update profile.</div>`))
				return
			}

			name := c.PostForm("name")
			phone := c.PostForm("phone")

			logAction("PROFILE UPDATE", fmt.Sprintf("User %s updating profile - Name: %s, Phone: %s", username.(string), name, phone))

			// Demo mode - simulate profile update
			message := "Profile updated successfully (demo mode)"
			if name != "" {
				message += " - Name: " + name
			}
			if phone != "" {
				message += " - Phone: " + phone
			}

			logAction("PROFILE UPDATE SUCCESS", fmt.Sprintf("Profile updated for user: %s", username.(string)))
			c.Data(http.StatusOK, "text/html; charset=utf-8",
				[]byte(`<div class="alert alert-success">`+message+`</div>`))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎭  P2K16 DEMO MODE - Development Testing Server")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📍 Server URL: http://localhost:%s\n", port)
	fmt.Println("🔑 Demo Login: username='demo', password=any")
	fmt.Println("📋 Features Available:")
	fmt.Println("   • User authentication")
	fmt.Println("   • Dashboard with badges")
	fmt.Println("   • Profile management (password change, profile update)")
	fmt.Println("   • Member listing")
	fmt.Println("⚠️  Note: No database - all changes are simulated")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🚀 Starting server...")

	logAction("SERVER STARTUP", fmt.Sprintf("Demo server starting on port %s", port))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
