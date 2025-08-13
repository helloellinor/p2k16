package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set up Gin router
	r := gin.New()

	// Add basic middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Session middleware
	store := cookie.NewStore([]byte("demo-secret-key"))
	store.Options(sessions.Options{
		MaxAge:   86400 * 7, // 7 days
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})
	r.Use(sessions.Sessions("p2k16_session", store))

	// Home page route
	r.GET("/", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html>
<head>
	<title>P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	<header class="navbar navbar-dark bg-dark">
		<div class="container">
			<a class="navbar-brand" href="/">P2K16</a>
			<div class="navbar-nav ms-auto">
				<a class="nav-link" href="/demo-admin">Admin Demo</a>
			</div>
		</div>
	</header>
	
	<main class="container mt-4">
		<div class="text-center mb-4">
			<h1>Welcome to P2K16</h1>
			<p class="lead">Hackerspace Management System</p>
		</div>
		
		<div class="row justify-content-center">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						<h2 class="card-title mb-0">Login</h2>
					</div>
					<div class="card-body">
						<form>
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
					</div>
				</div>
			</div>
		</div>
		
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
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// Admin demo page
	r.GET("/demo-admin", func(c *gin.Context) {
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
	<header class="navbar navbar-dark bg-dark">
		<div class="container">
			<a class="navbar-brand" href="/">P2K16</a>
			<nav class="navbar-nav">
				<span class="navbar-text">/ Admin</span>
			</nav>
		</div>
	</header>
	
	<main class="container mt-4">
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h1>Admin Console</h1>
		</div>
		
		<div class="row">
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Users</h5>
						<p class="card-text">Manage user accounts and permissions</p>
						<a href="/demo-admin-users" class="btn btn-primary">Manage Users</a>
					</div>
				</div>
			</div>
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Tools</h5>
						<p class="card-text">Manage tools and equipment</p>
						<a href="#" class="btn btn-primary">Manage Tools</a>
					</div>
				</div>
			</div>
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Badges</h5>
						<p class="card-text">Manage badges and awards</p>
						<a href="#" class="btn btn-primary">Manage Badges</a>
					</div>
				</div>
			</div>
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Companies</h5>
						<p class="card-text">Manage corporate memberships</p>
						<a href="#" class="btn btn-primary">Manage Companies</a>
					</div>
				</div>
			</div>
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Circles</h5>
						<p class="card-text">Manage user groups and permissions</p>
						<a href="#" class="btn btn-primary">Manage Circles</a>
					</div>
				</div>
			</div>
			<div class="col-md-4 mb-3">
				<div class="card h-100">
					<div class="card-body text-center">
						<h5 class="card-title">Logs</h5>
						<p class="card-text">View system logs</p>
						<a href="#" class="btn btn-primary">View Logs</a>
					</div>
				</div>
			</div>
		</div>
	</main>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// Admin users demo page
	r.GET("/demo-admin-users", func(c *gin.Context) {
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
	<header class="navbar navbar-dark bg-dark">
		<div class="container">
			<a class="navbar-brand" href="/">P2K16</a>
			<nav class="navbar-nav">
				<span class="navbar-text">/ <a href="/demo-admin" class="text-light">Admin</a> / Users</span>
			</nav>
		</div>
	</header>
	
	<main class="container mt-4">
		<div class="d-flex justify-content-between align-items-center mb-4">
			<h1>User Management</h1>
			<nav>
				<a href="/demo-admin" class="btn btn-outline-secondary">← Back to Admin</a>
			</nav>
		</div>
		
		<div class="card">
			<div class="card-header">
				<h5 class="card-title mb-0">User Accounts</h5>
			</div>
			<div class="card-body">
				<div class="table-responsive">
					<table class="table table-hover">
						<thead>
							<tr>
								<th>ID</th>
								<th>Username</th>
								<th>Name</th>
								<th>Email</th>
								<th>System</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td>1</td>
								<td><strong>admin</strong></td>
								<td>Administrator</td>
								<td>admin@p2k16.local</td>
								<td><span class="badge bg-warning">System</span></td>
								<td>
									<button class="btn btn-sm btn-outline-primary">View</button>
								</td>
							</tr>
							<tr>
								<td>2</td>
								<td><strong>john_doe</strong></td>
								<td>John Doe</td>
								<td>john@example.com</td>
								<td></td>
								<td>
									<button class="btn btn-sm btn-outline-primary">View</button>
								</td>
							</tr>
							<tr>
								<td>3</td>
								<td><strong>jane_smith</strong></td>
								<td>Jane Smith</td>
								<td>jane@example.com</td>
								<td></td>
								<td>
									<button class="btn btn-sm btn-outline-primary">View</button>
								</td>
							</tr>
							<tr>
								<td>4</td>
								<td><strong>alice_johnson</strong></td>
								<td>Alice Johnson</td>
								<td>alice@example.com</td>
								<td></td>
								<td>
									<button class="btn btn-sm btn-outline-primary">View</button>
								</td>
							</tr>
							<tr>
								<td>5</td>
								<td><strong>bob_wilson</strong></td>
								<td>Bob Wilson</td>
								<td>bob@example.com</td>
								<td></td>
								<td>
									<button class="btn btn-sm btn-outline-primary">View</button>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div class="d-flex justify-content-between align-items-center mt-3">
					<span class="text-muted">Showing 1-5 of 5 users</span>
					<div>
						<button class="btn btn-outline-secondary btn-sm" disabled>Previous</button>
						<button class="btn btn-outline-secondary btn-sm" disabled>Next</button>
					</div>
				</div>
			</div>
		</div>
	</main>
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// Start server
	port := "8080"
	log.Printf("Starting demo server on port %s", port)
	log.Printf("Visit http://localhost:%s for the home page", port)
	log.Printf("Visit http://localhost:%s/demo-admin for the admin interface", port)
	log.Printf("Visit http://localhost:%s/demo-admin-users for the user management", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}