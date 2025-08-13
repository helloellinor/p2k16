package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// ChiAdmin shows the admin interface (requires authentication)
func (h *ChiHandler) ChiAdmin(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin") + `
	<main class="container mt-4">
		<h1>Admin</h1>
		<section class="mt-4">
			<h2>Admin Operations</h2>
			<nav>
				<ul class="list-group">
					<li class="list-group-item"><a href="/admin/users">Users</a></li>
					<li class="list-group-item"><a href="/admin/tools">Tools</a></li>
					<li class="list-group-item"><a href="/admin/companies">Companies</a></li>
					<li class="list-group-item"><a href="/admin/circles">Circles</a></li>
					<li class="list-group-item"><a href="/admin/logs">Logs</a></li>
					<li class="list-group-item"><a href="/admin/config">Config</a></li>
				</ul>
			</nav>
		</section>
	</main>
</body>
</html>`

	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminUsers shows the users admin page
func (h *ChiHandler) ChiAdminUsers(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Users - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Users") + `
	<main class="container mt-4">
		<h1>Users</h1>
		<p>Manage user accounts and permissions.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminTools shows the tools admin page
func (h *ChiHandler) ChiAdminTools(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Tools - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Tools") + `
	<main class="container mt-4">
		<h1>Tools</h1>
		<p>Manage tools and equipment.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminCompanies shows the companies admin page
func (h *ChiHandler) ChiAdminCompanies(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Companies - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Companies") + `
	<main class="container mt-4">
		<h1>Companies</h1>
		<p>Manage company accounts and memberships.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminCircles shows the circles admin page
func (h *ChiHandler) ChiAdminCircles(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Circles - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Circles") + `
	<main class="container mt-4">
		<h1>Circles</h1>
		<p>Manage circles and permissions.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminLogs shows the logs admin page
func (h *ChiHandler) ChiAdminLogs(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Logs - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Logs") + `
	<main class="container mt-4">
		<h1>Logs</h1>
		<p>View system logs.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiAdminConfig shows the config admin page
func (h *ChiHandler) ChiAdminConfig(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Config - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	` + h.renderChiNavbarWithTrail(r, "Admin / Config") + `
	<main class="container mt-4">
		<h1>Config</h1>
		<p>System configuration settings.</p>
	</main>
</body>
</html>`
	h.writeHTML(w, http.StatusOK, html)
}

// ChiGetActiveMembers returns a list of active members
func (h *ChiHandler) ChiGetActiveMembers(w http.ResponseWriter, r *http.Request) {
	members, err := h.membershipRepo.GetActiveMembers()
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, "<p>Failed to load active members</p>")
		return
	}

	html := "<section aria-labelledby=\"active-members-title\">" +
		"<h2 id=\"active-members-title\">Active Members</h2>" +
		"<div class=\"row\">"

	for _, member := range members {
		html += fmt.Sprintf("<div class=\"col-md-4 mb-2\">"+
			"<div class=\"card border-success\">"+
			"<div class=\"card-body\">"+
			"<h6 class=\"card-title\">%s</h6>"+
			"<p class=\"card-text\">Status: Active</p>"+
			"</div>"+
			"</div>"+
			"</div>", member.Username)
	}

	html += "</div></section>"
	h.writeHTML(w, http.StatusOK, html)
}

// ChiGetAccounts returns a list of all accounts
func (h *ChiHandler) ChiGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountRepo.GetAll()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load accounts"})
		return
	}

	h.writeJSON(w, http.StatusOK, accounts)
}

// ChiGetAccount returns a specific account by ID
func (h *ChiHandler) ChiGetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
		return
	}

	account, err := h.accountRepo.FindByID(id)
	if err != nil {
		h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "Account not found"})
		return
	}

	h.writeJSON(w, http.StatusOK, account)
}

// ChiGetBadges returns a list of all badges
func (h *ChiHandler) ChiGetBadges(w http.ResponseWriter, r *http.Request) {
	badges, err := h.badgeRepo.GetAll()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load badges"})
		return
	}

	h.writeJSON(w, http.StatusOK, badges)
}

// ChiGetUserBadges returns badges for the current user
func (h *ChiHandler) ChiGetUserBadges(w http.ResponseWriter, r *http.Request) {
	user := middleware.ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, "<p>Authentication required</p>")
		return
	}

	badges, err := h.badgeRepo.GetUserBadges(user.ID)
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, "<p>Failed to load badges</p>")
		return
	}

	html := "<div class=\"d-flex flex-wrap gap-2\">"
	for _, badge := range badges {
		html += fmt.Sprintf("<span class=\"badge bg-primary\">%s</span>", badge.Title)
	}
	if len(badges) == 0 {
		html += "<span class=\"text-muted\">No badges earned yet</span>"
	}
	html += "</div>"

	h.writeHTML(w, http.StatusOK, html)
}

// ChiGetAvailableBadges returns available badges for awarding
func (h *ChiHandler) ChiGetAvailableBadges(w http.ResponseWriter, r *http.Request) {
	badges, err := h.badgeRepo.GetAll()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load badges"})
		return
	}

	h.writeJSON(w, http.StatusOK, badges)
}

// ChiCreateBadge creates a new badge
func (h *ChiHandler) ChiCreateBadge(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented yet"})
}

// ChiAwardBadge awards a badge to a user
func (h *ChiHandler) ChiAwardBadge(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented yet"})
}

// ChiGetMembershipStatusAPI returns membership status for API
func (h *ChiHandler) ChiGetMembershipStatusAPI(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "Not implemented yet"})
}

// ChiGetMembershipStatus returns membership status
func (h *ChiHandler) ChiGetMembershipStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, "<p>Authentication required</p>")
		return
	}

	html := `<div class="card border-success">
		<div class="card-header bg-success text-white">
			<h5 class="card-title mb-0">Membership Status</h5>
		</div>
		<div class="card-body">
			<p class="card-text">Status: <strong>Active</strong></p>
			<p class="card-text">Member since: <strong>2024-01-01</strong></p>
		</div>
	</div>`

	h.writeHTML(w, http.StatusOK, html)
}

// ChiGetActiveMembersDetailed returns detailed active members info
func (h *ChiHandler) ChiGetActiveMembersDetailed(w http.ResponseWriter, r *http.Request) {
	members, err := h.membershipRepo.GetActiveMembers()
	if err != nil {
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load active members"})
		return
	}

	h.writeJSON(w, http.StatusOK, members)
}