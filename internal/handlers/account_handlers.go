package handlers

import (
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

// GetAccount returns details for a specific account (API endpoint: GET /api/accounts/:id)
func (h *Handler) GetAccount(c *gin.Context) {
	logging.LogHandlerAction("API REQUEST", "Account details requested")

	// Parse account ID
	idStr := c.Param("id")
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		logging.LogError("VALIDATION ERROR", "Invalid account ID: "+idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid account ID",
		})
		return
	}

	// Get account from database
	account, err := h.accountRepo.FindByID(accountID)
	if err != nil {
		logging.LogError("DATABASE ERROR", "Failed to find account: "+err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Account not found",
		})
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