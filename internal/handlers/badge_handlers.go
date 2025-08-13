package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// BadgeResponse represents the public badge information for API responses
type BadgeResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}

// GetBadges returns a list of all badge descriptions (API endpoint: GET /api/badges/)
func (h *Handler) GetBadges(c *gin.Context) {
	// Get badge descriptions from database
	descriptions, err := h.badgeRepo.GetAllDescriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve badges",
		})
		return
	}

	// Convert to response format
	var badgeResponses []BadgeResponse
	for _, desc := range descriptions {
		response := BadgeResponse{
			ID:    desc.ID,
			Title: desc.Title,
		}
		
		if desc.Description.Valid {
			response.Description = desc.Description.String
		}
		if desc.Slug.Valid {
			response.Slug = desc.Slug.String
		}
		if desc.Icon.Valid {
			response.Icon = desc.Icon.String
		}
		if desc.Color.Valid {
			response.Color = desc.Color.String
		}
		
		badgeResponses = append(badgeResponses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   badgeResponses,
	})
}

// GetUserBadges returns user badges (for HTMX, requires authentication)
func (h *Handler) GetUserBadges(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	// Render centralized section
	html := h.renderUserBadgesSectionHTML(user.ID)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetAvailableBadges returns a list of available badge descriptions
func (h *Handler) GetAvailableBadges(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	descriptions, err := h.badgeRepo.GetAllDescriptions()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to load available badges</p>`))
		return
	}

	html := `
<section aria-labelledby="available-badges-title">
    <h2 id="available-badges-title">Explore Badges</h2>
    <ul>`

	for _, desc := range descriptions {
		// Check if user already has this badge
		has, _ := h.badgeRepo.AccountHasBadge(user.ID, desc.ID)
		if has {
			html += `
		<li>
			<span>` + desc.Title + `</span>
			<span>(already added)</span>
		</li>`
		} else {
			html += `
		<li>
			<span>` + desc.Title + `</span>
			<button 
				hx-post="/api/badges/award" 
				hx-vals='{"badge_title":"` + desc.Title + `"}'
				hx-target="#badge-feedback"
				hx-swap="innerHTML">Add to My Badges</button>
		</li>`
		}
	}

	html += `
	</ul>
	<div id="badge-feedback" aria-live="polite"></div>
	<div>
		<p>Want something new? Use the dedicated page to create a badge.</p>
		<p><a href="/badges/new">Create a new badge</a></p>
	</div>
</section>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CreateBadge creates a new badge description
func (h *Handler) CreateBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	title := c.PostForm("title")

	if title == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge title is required</p>`))
		return
	}

	// Check if badge already exists
	existing, _ := h.badgeRepo.FindBadgeDescriptionByTitle(title)
	if existing != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge "`+title+`" already exists</p>`))
		return
	}

	// Create badge description
	desc, err := h.badgeRepo.CreateBadgeDescription(title, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to create badge</p>`))
		return
	}

	// Award to self
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Badge created but failed to award</p>`))
		return
	}

	// Build success feedback and update the user badges via OOB swap
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := `
<section aria-live="polite">
	<p>Badge "` + title + `" created and added to your badges.</p>
</section>
` + `<div id="user-badges" hx-swap-oob="true">` + updated + `</div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// AwardBadge awards an existing badge to the current user
func (h *Handler) AwardBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	badgeTitle := c.PostForm("badge_title")

	if badgeTitle == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Badge title is required</p>`))
		return
	}

	// Find badge description
	desc, err := h.badgeRepo.FindBadgeDescriptionByTitle(badgeTitle)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte(`<p>Badge "`+badgeTitle+`" not found</p>`))
		return
	}

	// Prevent duplicates
	has, _ := h.badgeRepo.AccountHasBadge(user.ID, desc.ID)
	if has {
		c.Data(http.StatusOK, "text/html; charset=utf-8",
			[]byte(`<section aria-live="polite"><p>You already have '`+badgeTitle+`'.</p></section>`))
		return
	}

	// Award badge
	_, err = h.badgeRepo.AwardBadge(user.ID, desc.ID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to award badge</p>"))
		return
	}

	// Return feedback and out-of-band update of the user badges section
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := "<section aria-live=\"polite\"><p>Added '" + badgeTitle + "' to your badges.</p></section>" +
		"<div id=\"user-badges\" hx-swap-oob=\"true\">" + updated + "</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RemoveBadge removes an awarded badge from the current user
func (h *Handler) RemoveBadge(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	idStr := c.PostForm("account_badge_id")
	accountBadgeID, err := strconv.Atoi(idStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte(`<p>Invalid badge id</p>`))
		return
	}
	if err := h.badgeRepo.DeleteAccountBadge(accountBadgeID, user.ID); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<p>Failed to remove badge</p>`))
		return
	}
	updated := h.renderUserBadgesSectionHTML(user.ID)
	html := `<section aria-live="polite"><p>Badge removed.</p></section>` +
		`<div id="user-badges" hx-swap-oob="true">` + updated + `</div>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// BadgeCreatePage shows a dedicated page to create a new badge (requires authentication)
func (h *Handler) BadgeCreatePage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
	<title>Create Badge - P2K16</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
    
</head>
<body>
	` + h.renderNavbarWithTrail(c, "Create Badge") + `
    
	<main>
		<section aria-labelledby="create-badge-title">
			<header>
				<h1 id="create-badge-title">Create a New Badge</h1>
			</header>
			<div>
				<form hx-post="/api/badges/create" hx-target="#create-result">
					<div>
						<label for="title">Title</label>
						<input id="title" type="text" name="title" required>
						<div>Choose a short, clear name (e.g., Laser Cutter Trained)</div>
					</div>
					<div>
						<button type="submit">Create Badge</button>
					</div>
				</form>
				<div id="create-result" aria-live="polite"></div>
			</div>
		</section>

		<section>
			<p>After creating, you can award it to yourself or others from the Admin Console.</p>
		</section>
	</main>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}