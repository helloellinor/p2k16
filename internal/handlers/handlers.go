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
    
</head>
<body>
	` + h.renderNavbar(c) + `
    
	<main>
		<div>
			<h1>Welcome to P2K16</h1>
			<p>Hackerspace Management System</p>
		</div>
`
	if user == nil {
		html += `
		<section>
			<header>
				<h2>Login</h2>
			</header>
			<div>
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
		`
	}
	html += `
		<section>
			<article>
				<header>
					<h2>System Status</h2>
				</header>
				<div>
					<p>Online</p>
					<p>Database connected - all systems operational</p>
				</div>
			</article>
		</section>

	</main>
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
