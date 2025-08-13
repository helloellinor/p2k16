package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/helloellinor/p2k16/internal/logging"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
	"github.com/helloellinor/p2k16/internal/session"
)

// ChiHandler wraps the original handler with chi-compatible methods
type ChiHandler struct {
	*Handler
	sessionManager *session.ChiSessionManager
}

// NewChiHandler creates a new chi-compatible handler
func NewChiHandler(accountRepo *models.AccountRepository, circleRepo *models.CircleRepository, badgeRepo *models.BadgeRepository, toolRepo *models.ToolRepository, eventRepo *models.EventRepository, membershipRepo *models.MembershipRepository, sessionManager *session.ChiSessionManager) *ChiHandler {
	return &ChiHandler{
		Handler:        NewHandler(accountRepo, circleRepo, badgeRepo, toolRepo, eventRepo, membershipRepo),
		sessionManager: sessionManager,
	}
}

// Helper function to write HTML response
func (h *ChiHandler) writeHTML(w http.ResponseWriter, statusCode int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(html))
}

// Helper function to write JSON response
func (h *ChiHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to write plain text response
func (h *ChiHandler) writeText(w http.ResponseWriter, statusCode int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(text))
}

// ChiHome renders the front page for chi
func (h *ChiHandler) ChiHome(w http.ResponseWriter, r *http.Request) {
	logging.LogHandlerAction("PAGE REQUEST", "Home page visited")
	user := middleware.ChiGetCurrentUser(r)

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
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbar(r) + `
    
	<main class="container mt-4">
		<div class="text-center mb-4">
			<h1>Welcome to P2K16</h1>
			<p class="lead">Hackerspace Management System</p>
		</div>
`
	if user == nil {
		html += `
		<div class="row justify-content-center">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h2 class="card-title mb-0">Login</h2>
					</div>
					<div class="card-body">
						<form hx-post="/api/auth/login" hx-target="#login-result" method="post" action="/api/auth/login">
							<div class="mb-3">
								<label for="username" class="form-label">Username</label>
								<input type="text" class="form-control" id="username" name="username" required>
							</div>
							<div class="mb-3">
								<label for="password" class="form-label">Password</label>
								<input type="password" class="form-control" id="password" name="password" required>
							</div>
							<div class="d-grid">
								<button type="submit" class="btn btn-primary">Login</button>
							</div>
						</form>
						<div id="login-result" class="mt-3"></div>
					</div>
				</div>
			</div>
		</div>
		`
	} else {
		// Authenticated user content - similar to Python front-page.html
		html += `
		<div class="row">
			<div class="col-md-8 offset-md-2">
				<div id="membership-status" hx-get="/api/membership/status" hx-trigger="load" hx-target="this">
					<div class="text-center">
						<div class="spinner-border" role="status">
							<span class="visually-hidden">Loading membership status...</span>
						</div>
					</div>
				</div>

				<div class="mt-4">
					<div class="card">
						<div class="card-header">
							<h5 class="card-title mb-0">Tools</h5>
						</div>
						<div class="card-body">
							<div id="tools-section" hx-get="/api/tools" hx-trigger="load" hx-target="this">
								<div class="text-center">
									<div class="spinner-border spinner-border-sm" role="status">
										<span class="visually-hidden">Loading tools...</span>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>

				<div class="mt-4">
					<div class="card">
						<div class="card-header">
							<h5 class="card-title mb-0">Active Checkouts</h5>
						</div>
						<div class="card-body">
							<div id="active-checkouts" hx-get="/api/tools/checkouts" hx-trigger="load" hx-target="this">
								<div class="text-center">
									<div class="spinner-border spinner-border-sm" role="status">
										<span class="visually-hidden">Loading checkouts...</span>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		`
	}
	html += `
		<div class="row mt-4">
			<div class="col-md-8 offset-md-2">
				<div class="card">
					<div class="card-header">
						<h5 class="card-title mb-0">System Status</h5>
					</div>
					<div class="card-body">
						<div class="d-flex align-items-center">
							<span class="badge bg-success me-2">Online</span>
							<span>Database connected - all systems operational</span>
						</div>
					</div>
				</div>
			</div>
		</div>

	</main>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`

	h.writeHTML(w, http.StatusOK, html)
}

// ChiDashboard redirects to the profile page for authenticated users
func (h *ChiHandler) ChiDashboard(w http.ResponseWriter, r *http.Request) {
	logging.LogHandlerAction("PAGE REQUEST", "Dashboard page requested - redirecting to profile")
	http.Redirect(w, r, "/profile", http.StatusFound)
}

// ChiProfile shows the user profile page (requires authentication)
func (h *ChiHandler) ChiProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.ChiGetCurrentUser(r)

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Profile - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Profile") + `
    
	<main class="container mt-4">
		<div class="row justify-content-center">
			<div class="col-md-8">
				<div class="text-center mb-4">
					<h1>Membership Card</h1>
					<p class="lead">Front shows your info and badges. Click Edit to flip to the back and manage details.</p>
				</div>
				<div class="card">
					<div class="card-body">
						<div id="membership-card">` + h.renderChiProfileCardFrontHTML(user) + `</div>
					</div>
				</div>
			</div>
		</div>
    </main>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`

	h.writeHTML(w, http.StatusOK, html)
}

// ChiLogin handles user authentication page
func (h *ChiHandler) ChiLogin(w http.ResponseWriter, r *http.Request) {
	user := middleware.ChiGetCurrentUser(r)
	
	// If already authenticated, redirect to home
	if user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
	<title>Login - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Login") + `
    
	<main class="container mt-4">
		<div class="row justify-content-center">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h2 class="card-title mb-0">Login</h2>
					</div>
					<div class="card-body">
						<form hx-post="/api/auth/login" hx-target="#login-result" method="post" action="/api/auth/login">
							<div class="mb-3">
								<label for="username" class="form-label">Username</label>
								<input type="text" class="form-control" id="username" name="username" required>
							</div>
							<div class="mb-3">
								<label for="password" class="form-label">Password</label>
								<input type="password" class="form-control" id="password" name="password" required>
							</div>
							<div class="d-grid">
								<button type="submit" class="btn btn-primary">Login</button>
							</div>
						</form>
						<div id="login-result" class="mt-3"></div>
					</div>
				</div>
			</div>
		</div>
	</main>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`

	h.writeHTML(w, http.StatusOK, html)
}

// ChiLogout handles user logout
func (h *ChiHandler) ChiLogout(w http.ResponseWriter, r *http.Request) {
	logging.LogHandlerAction("USER ACTION", "User logout requested")
	if err := h.sessionManager.LogoutUser(w, r); err != nil {
		logging.LogError("LOGOUT ERROR", fmt.Sprintf("Failed to logout user: %v", err))
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	logging.LogSuccess("USER ACTION", "User successfully logged out - redirecting to home")
	http.Redirect(w, r, "/", http.StatusFound)
}