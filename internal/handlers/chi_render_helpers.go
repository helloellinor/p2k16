package handlers

import (
	"net/http"

	"github.com/helloellinor/p2k16/internal/middleware"
)

// renderChiNavbar renders the navigation bar for chi
func (h *ChiHandler) renderChiNavbar(r *http.Request) string {
	user := middleware.ChiGetCurrentUser(r)
	
	navbar := `
	<nav class="navbar navbar-expand-lg navbar-dark bg-primary">
		<div class="container">
			<a class="navbar-brand fw-bold" href="/">P2K16</a>
			
			<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
			
			<div class="collapse navbar-collapse" id="navbarNav">
				<ul class="navbar-nav me-auto">
					<li class="nav-item">
						<a class="nav-link" href="/">Home</a>
					</li>`
	
	if user != nil {
		navbar += `
					<li class="nav-item">
						<a class="nav-link" href="/profile">Profile</a>
					</li>
					<li class="nav-item">
						<a class="nav-link" href="/admin">Admin</a>
					</li>`
	}
	
	navbar += `
				</ul>
				<ul class="navbar-nav">`
	
	if user != nil {
		navbar += `
					<li class="nav-item dropdown">
						<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
							` + user.Username + `
						</a>
						<ul class="dropdown-menu" aria-labelledby="navbarDropdown">
							<li><a class="dropdown-item" href="/profile">Profile</a></li>
							<li><hr class="dropdown-divider"></li>
							<li>
								<form method="post" action="/logout" class="d-inline">
									<button type="submit" class="dropdown-item">Logout</button>
								</form>
							</li>
						</ul>
					</li>`
	} else {
		navbar += `
					<li class="nav-item">
						<a class="nav-link" href="/login">Login</a>
					</li>`
	}
	
	navbar += `
				</ul>
			</div>
		</div>
	</nav>`
	
	return navbar
}

// renderChiNavbarWithTrail renders the navigation bar with breadcrumb trail for chi
func (h *ChiHandler) renderChiNavbarWithTrail(r *http.Request, title string) string {
	user := middleware.ChiGetCurrentUser(r)
	
	navbar := h.renderChiNavbar(r)
	
	trail := `
	<div class="container mt-2">
		<nav aria-label="breadcrumb">
			<ol class="breadcrumb">
				<li class="breadcrumb-item"><a href="/">Home</a></li>`
	
	if user != nil {
		trail += `
				<li class="breadcrumb-item active" aria-current="page">` + title + `</li>`
	}
	
	trail += `
			</ol>
		</nav>
	</div>`
	
	return navbar + trail
}

// renderChiProfileCardFrontHTML renders the front of the membership card for chi
func (h *ChiHandler) renderChiProfileCardFrontHTML(user *middleware.ChiAuthenticatedUser) string {
	if user == nil {
		return `<p>No user data available</p>`
	}
	
	return `
	<div class="membership-card-front">
		<div class="row mb-3">
			<div class="col-md-8">
				<h3 class="mb-2">` + user.Username + `</h3>
				<p class="text-muted mb-1">Member #` + string(rune(user.ID)) + `</p>
				<p class="text-muted">` + user.Account.Email + `</p>
			</div>
			<div class="col-md-4 text-end">
				<button class="btn btn-outline-primary btn-sm" 
						hx-get="/api/profile/card/back" 
						hx-target="#membership-card" 
						hx-swap="innerHTML">
					Edit
				</button>
			</div>
		</div>
		
		<div class="row">
			<div class="col-12">
				<h5>Badges</h5>
				<div id="user-badges" hx-get="/api/user/badges" hx-trigger="load" hx-target="this">
					<div class="spinner-border spinner-border-sm" role="status">
						<span class="visually-hidden">Loading badges...</span>
					</div>
				</div>
			</div>
		</div>
	</div>`
}

// renderChiProfileCardBackHTML renders the back (editing) of the membership card for chi
func (h *ChiHandler) renderChiProfileCardBackHTML(user *middleware.ChiAuthenticatedUser) string {
	if user == nil {
		return `<p>No user data available</p>`
	}
	
	return `
	<div class="membership-card-back">
		<div class="row mb-3">
			<div class="col-md-8">
				<h3 class="mb-2">Edit Profile</h3>
			</div>
			<div class="col-md-4 text-end">
				<button class="btn btn-outline-secondary btn-sm" 
						hx-get="/api/profile/card/front" 
						hx-target="#membership-card" 
						hx-swap="innerHTML">
					Cancel
				</button>
			</div>
		</div>
		
		<form hx-post="/profile/update" hx-target="#membership-card" hx-swap="innerHTML">
			<div class="mb-3">
				<label for="username" class="form-label">Username</label>
				<input type="text" class="form-control" id="username" name="username" value="` + user.Username + `" readonly>
			</div>
			<div class="mb-3">
				<label for="email" class="form-label">Email</label>
				<input type="email" class="form-control" id="email" name="email" value="` + user.Account.Email + `">
			</div>
			<div class="mb-3">
				<label for="current_password" class="form-label">Current Password (required to save changes)</label>
				<input type="password" class="form-control" id="current_password" name="current_password">
			</div>
			<div class="mb-3">
				<label for="new_password" class="form-label">New Password (optional)</label>
				<input type="password" class="form-control" id="new_password" name="new_password">
			</div>
			<div class="mb-3">
				<label for="confirm_password" class="form-label">Confirm New Password</label>
				<input type="password" class="form-control" id="confirm_password" name="confirm_password">
			</div>
			<div class="d-grid gap-2 d-md-flex justify-content-md-end">
				<button type="submit" class="btn btn-primary">Save Changes</button>
			</div>
		</form>
	</div>`
}