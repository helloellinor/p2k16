package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/middleware"
)

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

// Login handles user authentication page
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