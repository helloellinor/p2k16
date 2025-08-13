package handlers

import (
	"fmt"
	"net/http"

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
	` + h.renderNavbar(c) + `
    
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

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Dashboard redirects to the profile page for authenticated users
func (h *Handler) Dashboard(c *gin.Context) {
	logging.LogHandlerAction("PAGE REQUEST", "Dashboard page requested - redirecting to profile")
	c.Redirect(http.StatusFound, "/profile")
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
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Profile") + `
    
	<main>
		<div>
			<h1>Membership Card</h1>
			<p>Front shows your info and badges. Click Edit to flip to the back and manage details.</p>
		</div>
		<section>
			<div id="membership-card">` + h.renderProfileCardFrontHTML(user) + `</div>
		</section>


    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ProfileCardFront returns the front of the membership card (requires auth)
func (h *Handler) ProfileCardFront(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	html := h.renderProfileCardFrontHTML(user)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ProfileCardBack returns the back (editing) of the membership card (requires auth)
func (h *Handler) ProfileCardBack(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	html := h.renderProfileCardBackHTML(user)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
