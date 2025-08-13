package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin shows the admin interface (requires authentication)
func (h *Handler) Admin(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Admin - P2K16</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin") + `
    
	<main>
		<div>
			<h1>Admin Console</h1>
			<p>Workspace for privileged tasks: manage users, tools, companies, circles, logs, and configuration.</p>
		</div>
		<section>
			<nav aria-label="Admin sections">
				<ul>
					<li><a href="/admin/users">Users</a></li>
					<li><a href="/admin/tools">Tools</a></li>
					<li><a href="/admin/companies">Companies</a></li>
					<li><a href="/admin/circles">Circles</a></li>
					<li><a href="/admin/logs">Logs</a></li>
					<li><a href="/admin/config">Config</a></li>
				</ul>
			</nav>
		</section>

    </main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminUsers shows the users admin page
func (h *Handler) AdminUsers(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Users - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Users") + `
	<main>
		<h1>Users</h1>
		<p>Manage user accounts and permissions.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminTools shows the tools admin page
func (h *Handler) AdminTools(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Tools - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Tools") + `
	<main>
		<h1>Tools</h1>
		<p>Manage tools and equipment.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminCompanies shows the companies admin page
func (h *Handler) AdminCompanies(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Companies - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Companies") + `
	<main>
		<h1>Companies</h1>
		<p>Manage corporate memberships.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminCircles shows the circles admin page
func (h *Handler) AdminCircles(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Circles - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Circles") + `
	<main>
		<h1>Circles</h1>
		<p>Manage user groups and permissions.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminLogs shows the logs admin page
func (h *Handler) AdminLogs(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Logs - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Logs") + `
	<main>
		<h1>Logs</h1>
		<p>View system logs.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AdminConfig shows the config admin page
func (h *Handler) AdminConfig(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Admin / Config - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Admin / Config") + `
	<main>
		<h1>Config</h1>
		<p>Configure system parameters.</p>
	</main>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}