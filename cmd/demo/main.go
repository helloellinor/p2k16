package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
)

// Demo mode - simple server without database for UI testing
func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Serve static files
	r.Static("/styles", "./styles")

	// Load templates
	r.LoadHTMLGlob("cmd/server/templates/*.html")

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
		logging.LogDemoAction("PAGE REQUEST", "Home page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		var user interface{}
		if username != nil {
			logging.LogDemoAction("USER STATUS", fmt.Sprintf("Authenticated user: %s", username.(string)))
			user = map[string]interface{}{
				"Username": username.(string),
			}
		} else {
			logging.LogDemoAction("USER STATUS", "Anonymous user")
		}

		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":    "Home",
			"user":     user,
			"demoMode": true,
			"content": gin.H{
				"template": "home",
				"data": gin.H{
					"user":       user,
					"demoMode":   true,
					"isLoggedIn": username != nil,
				},
			},
		})
	})

	r.GET("/login", func(c *gin.Context) {
		logging.LogDemoAction("PAGE REQUEST", "Login page visited")
		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":    "Login",
			"narrow":   true,
			"demoMode": true,
			"content": gin.H{
				"template": "login",
				"data": gin.H{
					"demoMode": true,
				},
			},
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		logging.LogDemoAction("PAGE REQUEST", "Dashboard page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		if username == nil {
			logging.LogDemoAction("ACCESS DENIED", "Dashboard access attempted without authentication")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		logging.LogDemoAction("USER ACCESS", fmt.Sprintf("Dashboard accessed by user: %s", username.(string)))

		user := map[string]interface{}{
			"Username": username.(string),
		}

		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":    "Dashboard",
			"user":     user,
			"demoMode": true,
			"content": gin.H{
				"template": "dashboard",
				"data": gin.H{
					"user":     user,
					"demoMode": true,
				},
			},
		})
	})

	r.GET("/profile", func(c *gin.Context) {
		logging.LogDemoAction("PAGE REQUEST", "Profile page visited")
		session := sessions.Default(c)
		username := session.Get("username")

		if username == nil {
			logging.LogDemoAction("ACCESS DENIED", "Profile access attempted without authentication")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		logging.LogDemoAction("USER ACCESS", fmt.Sprintf("Profile page accessed by user: %s", username.(string)))

		user := map[string]interface{}{
			"Username": username.(string),
			"Email":    username.(string) + "@demo.p2k16.io",
		}

		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":    "Profile",
			"user":     user,
			"demoMode": true,
			"content": gin.H{
				"template": "profile",
				"data": gin.H{
					"user":     user,
					"demoMode": true,
				},
			},
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		logging.LogDemoAction("USER ACTION", "User logout requested")
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		logging.LogDemoAction("USER ACTION", "User session cleared - redirecting to home")
		c.Redirect(http.StatusFound, "/")
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/members/active", func(c *gin.Context) {
			logging.LogDemoAction("API REQUEST", "Active members list requested")
			html := `
				<div class="p2k16-grid p2k16-grid--2-col">
					<div class="p2k16-card">
						<div class="p2k16-card__body">
							<h6 class="p2k16-card__title">Demo Admin</h6>
							<p class="p2k16-text--secondary">System Administrator</p>
							<small class="p2k16-text--muted">Last active: 2 hours ago</small>
						</div>
					</div>
					<div class="p2k16-card">
						<div class="p2k16-card__body">
							<h6 class="p2k16-card__title">Demo User</h6>
							<p class="p2k16-text--secondary">Regular Member</p>
							<small class="p2k16-text--muted">Last active: 1 day ago</small>
						</div>
					</div>
				</div>`

			logging.LogDemoAction("API RESPONSE", "Active members list returned (2 demo users)")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		})

		api.POST("/auth/login", func(c *gin.Context) {
			logging.LogDemoAction("API REQUEST", "Login attempt received")
			username := c.PostForm("username")
			password := c.PostForm("password")

			logging.LogDemoAction("LOGIN ATTEMPT", fmt.Sprintf("Username: %s", username))

			if username == "" || password == "" {
				logging.LogDemoAction("LOGIN FAILED", "Missing username or password")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">Username and password are required</div>`))
				return
			}

			// Demo login - accept "demo" with any password
			if username == "demo" {
				logging.LogDemoAction("LOGIN SUCCESS", fmt.Sprintf("Demo user authenticated: %s", username))
				session := sessions.Default(c)
				session.Set("username", username)
				session.Save()

				html := `
					<div class="p2k16-alert p2k16-alert--success">
						Login successful! Welcome to demo mode, ` + username + `
					</div>
					<script>
						setTimeout(function() {
							window.location.href = '/dashboard';
						}, 1000);
					</script>`

				logging.LogDemoAction("SESSION CREATED", fmt.Sprintf("Session created for user: %s", username))
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
				return
			}

			logging.LogDemoAction("LOGIN FAILED", fmt.Sprintf("Invalid credentials for username: %s", username))
			c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
				[]byte(`<div class="p2k16-alert p2k16-alert--danger">Invalid username or password. Use "demo" for demo mode.</div>`))
		})

		api.GET("/user/badges", func(c *gin.Context) {
			logging.LogDemoAction("API REQUEST", "User badges requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logging.LogDemoAction("API ERROR", "Badges request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--warning">Please log in to view badges.</div>`))
				return
			}

			logging.LogDemoAction("BADGES LOADED", fmt.Sprintf("Badges loaded for user: %s", username.(string)))
			html := `
				<div class="badge-list">
					<span class="p2k16-badge p2k16-badge--primary">Go Programming</span>
					<span class="p2k16-badge p2k16-badge--success">HTMX Expert</span>
					<span class="p2k16-badge p2k16-badge--info">Database Admin</span>
					<span class="p2k16-badge p2k16-badge--warning">Demo User</span>
					<p class="p2k16-text--muted p2k16-mt-4">User ` + username.(string) + ` has 4 demo badges.</p>
				</div>`

			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		})

		// Profile management API endpoints (Phase 2)
		api.POST("/profile/change-password", func(c *gin.Context) {
			logging.LogDemoAction("API REQUEST", "Password change requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logging.LogDemoAction("API ERROR", "Password change request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">Please log in to change password.</div>`))
				return
			}

			oldPassword := c.PostForm("oldPassword")
			newPassword := c.PostForm("newPassword")
			confirmPassword := c.PostForm("confirmPassword")

			logging.LogDemoAction("PASSWORD CHANGE", fmt.Sprintf("User %s attempting password change", username.(string)))

			// Validation
			if oldPassword == "" || newPassword == "" || confirmPassword == "" {
				logging.LogDemoAction("PASSWORD CHANGE FAILED", "Missing required fields")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">All fields are required</div>`))
				return
			}

			if newPassword != confirmPassword {
				logging.LogDemoAction("PASSWORD CHANGE FAILED", "New passwords do not match")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">New passwords do not match</div>`))
				return
			}

			if len(newPassword) < 6 {
				logging.LogDemoAction("PASSWORD CHANGE FAILED", "Password too short")
				c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">Password must be at least 6 characters long</div>`))
				return
			}

			// Demo mode - simulate password change
			logging.LogDemoAction("PASSWORD CHANGE SUCCESS", fmt.Sprintf("Password changed successfully for user: %s (demo mode)", username.(string)))
			c.Data(http.StatusOK, "text/html; charset=utf-8",
				[]byte(`<div class="p2k16-alert p2k16-alert--success">Password changed successfully (demo mode)</div>`))
		})

		api.POST("/profile/update", func(c *gin.Context) {
			logging.LogDemoAction("API REQUEST", "Profile update requested")
			session := sessions.Default(c)
			username := session.Get("username")

			if username == nil {
				logging.LogDemoAction("API ERROR", "Profile update request without authentication")
				c.Data(http.StatusUnauthorized, "text/html; charset=utf-8",
					[]byte(`<div class="p2k16-alert p2k16-alert--danger">Please log in to update profile.</div>`))
				return
			}

			name := c.PostForm("name")
			phone := c.PostForm("phone")

			logging.LogDemoAction("PROFILE UPDATE", fmt.Sprintf("User %s updating profile - Name: %s, Phone: %s", username.(string), name, phone))

			// Demo mode - simulate profile update
			message := "Profile updated successfully (demo mode)"
			if name != "" {
				message += " - Name: " + name
			}
			if phone != "" {
				message += " - Phone: " + phone
			}

			logging.LogDemoAction("PROFILE UPDATE SUCCESS", fmt.Sprintf("Profile updated for user: %s", username.(string)))
			c.Data(http.StatusOK, "text/html; charset=utf-8",
				[]byte(`<div class="p2k16-alert p2k16-alert--success">`+message+`</div>`))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Log startup with enhanced formatting
	features := []string{
		"User authentication",
		"Dashboard with badges",
		"Profile management (password change, profile update)",
		"Member listing",
	}
	
	logging.DemoLogger.LogStartup("demo", port, features)

	logging.LogDemoAction("SERVER STARTUP", fmt.Sprintf("Demo server starting on port %s", port))

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
