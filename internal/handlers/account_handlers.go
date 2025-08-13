package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
)

// AccountResponse represents the public account information for API responses
type AccountResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	System   bool   `json:"system"`
}

// AccountListResponse represents the response for account listing
type AccountListResponse struct {
	Status string            `json:"status"`
	Data   []AccountResponse `json:"data"`
	Total  int               `json:"total"`
}

// GetAccounts returns a list of all accounts (API endpoint: GET /api/accounts/)
func (h *Handler) GetAccounts(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "Account listing requested")

	// Check if this is an HTMX request for HTML response
	if c.GetHeader("HX-Request") == "true" {
		h.GetAccountsHTML(c)
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 50
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get accounts from database
	accounts, err := h.accountRepo.GetAllAccounts(limit, offset)
	if err != nil {
		logging.LogError("DATABASE ERROR", "Failed to retrieve accounts: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve accounts",
		})
		return
	}

	// Get total count
	total, err := h.accountRepo.GetAccountCount()
	if err != nil {
		logging.LogError("DATABASE ERROR", "Failed to get account count: "+err.Error())
		// Continue with partial data
		total = len(accounts)
	}

	// Convert to response format (exclude password and sensitive fields)
	var accountResponses []AccountResponse
	for _, account := range accounts {
		response := AccountResponse{
			ID:       account.ID,
			Username: account.Username,
			Email:    account.Email,
			System:   account.System.Valid && account.System.Bool,
		}
		
		if account.Name.Valid {
			response.Name = account.Name.String
		}
		if account.Phone.Valid {
			response.Phone = account.Phone.String
		}
		
		accountResponses = append(accountResponses, response)
	}

	logging.LogSuccess("API SUCCESS", "Account list retrieved successfully")
	c.JSON(http.StatusOK, AccountListResponse{
		Status: "success",
		Data:   accountResponses,
		Total:  total,
	})
}

// GetAccountsHTML returns accounts as HTML for HTMX requests
func (h *Handler) GetAccountsHTML(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get accounts from database
	accounts, err := h.accountRepo.GetAllAccounts(limit, offset)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte(`<div class="alert alert-danger">Failed to load users</div>`))
		return
	}

	// Get total count
	total, _ := h.accountRepo.GetAccountCount()

	// Build HTML table
	html := `<div class="table-responsive">
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
			<tbody>`

	for _, account := range accounts {
		name := ""
		if account.Name.Valid {
			name = account.Name.String
		}
		
		systemBadge := ""
		if account.System.Valid && account.System.Bool {
			systemBadge = `<span class="badge bg-warning">System</span>`
		}

		html += fmt.Sprintf(`
				<tr>
					<td>%d</td>
					<td><strong>%s</strong></td>
					<td>%s</td>
					<td>%s</td>
					<td>%s</td>
					<td>
						<button class="btn btn-sm btn-outline-primary" hx-get="/api/accounts/%d" hx-target="#user-details">View</button>
					</td>
				</tr>`, account.ID, account.Username, name, account.Email, systemBadge, account.ID)
	}

	html += `
			</tbody>
		</table>
	</div>
	<div class="d-flex justify-content-between align-items-center mt-3">
		<span class="text-muted">Showing ` + fmt.Sprintf("%d-%d of %d users", offset+1, offset+len(accounts), total) + `</span>
		<div>` 
		
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		html += fmt.Sprintf(`<button class="btn btn-outline-secondary me-2" hx-get="/api/accounts?limit=%d&offset=%d" hx-target="#users-list">Previous</button>`, limit, prevOffset)
	}
	
	if offset+len(accounts) < total {
		nextOffset := offset + limit
		html += fmt.Sprintf(`<button class="btn btn-outline-secondary" hx-get="/api/accounts?limit=%d&offset=%d" hx-target="#users-list">Next</button>`, limit, nextOffset)
	}
	
	html += `</div>
	</div>
	<div id="user-details" class="mt-3"></div>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetAccount returns details for a specific account (API endpoint: GET /api/accounts/:id)
func (h *Handler) GetAccount(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "Account details requested")

	// Parse account ID
	idStr := c.Param("id")
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		logging.LogError("VALIDATION ERROR", "Invalid account ID: "+idStr)
		if c.GetHeader("HX-Request") == "true" {
			c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
				[]byte(`<div class="alert alert-danger">Invalid account ID</div>`))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid account ID",
			})
		}
		return
	}

	// Get account from database
	account, err := h.accountRepo.FindByID(accountID)
	if err != nil {
		logging.LogError("DATABASE ERROR", "Failed to find account: "+err.Error())
		if c.GetHeader("HX-Request") == "true" {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8",
				[]byte(`<div class="alert alert-danger">Account not found</div>`))
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Account not found",
			})
		}
		return
	}

	// Check if this is an HTMX request for HTML response
	if c.GetHeader("HX-Request") == "true" {
		// Return HTML card
		html := `<div class="card">
			<div class="card-header">
				<h6 class="card-title mb-0">User Details</h6>
			</div>
			<div class="card-body">
				<dl class="row">
					<dt class="col-sm-3">ID:</dt>
					<dd class="col-sm-9">` + fmt.Sprintf("%d", account.ID) + `</dd>
					
					<dt class="col-sm-3">Username:</dt>
					<dd class="col-sm-9"><strong>` + account.Username + `</strong></dd>
					
					<dt class="col-sm-3">Email:</dt>
					<dd class="col-sm-9">` + account.Email + `</dd>`
		
		if account.Name.Valid {
			html += `
					<dt class="col-sm-3">Name:</dt>
					<dd class="col-sm-9">` + account.Name.String + `</dd>`
		}
		
		if account.Phone.Valid {
			html += `
					<dt class="col-sm-3">Phone:</dt>
					<dd class="col-sm-9">` + account.Phone.String + `</dd>`
		}
		
		if account.System.Valid && account.System.Bool {
			html += `
					<dt class="col-sm-3">Type:</dt>
					<dd class="col-sm-9"><span class="badge bg-warning">System Account</span></dd>`
		}
		
		html += `
					<dt class="col-sm-3">Created:</dt>
					<dd class="col-sm-9">` + account.CreatedAt.Format("2006-01-02 15:04:05") + `</dd>
				</dl>
			</div>
		</div>`
		
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
		return
	}

	// Convert to response format (exclude password and sensitive fields)
	response := AccountResponse{
		ID:       account.ID,
		Username: account.Username,
		Email:    account.Email,
		System:   account.System.Valid && account.System.Bool,
	}
	
	if account.Name.Valid {
		response.Name = account.Name.String
	}
	if account.Phone.Valid {
		response.Phone = account.Phone.String
	}

	logging.LogSuccess("API SUCCESS", "Account details retrieved successfully")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}