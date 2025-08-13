package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// GetTools returns a list of all tools
func (h *Handler) GetTools(c *gin.Context) {
	tools, err := h.toolRepo.GetAllTools()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to load tools</p>"))
		return
	}

	html := "<section aria-labelledby=\"tools-title\">" +
		"<h2 id=\"tools-title\">Available Tools</h2>" +
		"<div>"

	for _, tool := range tools {
		html += "<div class=\"col-md-6 mb-3\">" +
			"<div class=\"card border-primary\">" +
			"<div class=\"card-body\">" +
			"<h6 class=\"card-title\">" + tool.Name + "</h6>" +
			"<p class=\"card-text\">Description: " + tool.Description.String + "</p>" +
			"<button class=\"btn btn-success btn-sm\" " +
			"hx-post=\"/api/tools/checkout\" " +
			"hx-vals='{\"tool_id\":\"" + strconv.Itoa(tool.ID) + "\"}' " +
			"hx-target=\"#tool-result\" " +
			"hx-swap=\"innerHTML\">" +
			"Checkout" +
			"</button>" +
			"</div>" +
			"</div>" +
			"</div>"
	}

	html += "</div>" +
		"<div id=\"tool-result\"></div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetActiveCheckouts returns currently checked out tools
func (h *Handler) GetActiveCheckouts(c *gin.Context) {
	checkouts, err := h.toolRepo.GetActiveCheckouts()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to load active checkouts</p>"))
		return
	}

	html := "<section aria-labelledby=\"checkouts-title\">" +
		"<h2 id=\"checkouts-title\">Currently Checked Out Tools</h2>" +
		"<div>"

	if len(checkouts) == 0 {
		html += "<p>No tools currently checked out.</p>"
	} else {
		html += "<ul>"
		for _, checkout := range checkouts {
			html += "<li>" +
				"<div>" + checkout.Tool.Name + " (" + checkout.Tool.Description.String + ") - Checked out by: " + checkout.Account.Username + " - Since: " + checkout.CheckoutAt.Format("2006-01-02 15:04") + "</div>" +
				"<button " +
				"hx-post=\"/api/tools/checkin\" " +
				"hx-vals='{\"checkout_id\":\"" + strconv.Itoa(checkout.ID) + "\"}' " +
				"hx-target=\"#tool-result\" " +
				"hx-swap=\"innerHTML\">" +
				"Check In" +
				"</button>" +
				"</li>"
		}
		html += "</ul>"
	}

	html += "<div id=\"tool-result\"></div>" +
		"</div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckoutTool handles tool checkout
func (h *Handler) CheckoutTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	toolIDStr := c.PostForm("tool_id")

	toolID, err := strconv.Atoi(toolIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<p>Invalid tool ID</p>"))
		return
	}

	// Check if tool exists
	tool, err := h.toolRepo.FindToolByID(toolID)
	if err != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8",
			[]byte("<p>Tool not found</p>"))
		return
	}

	// Create checkout record
	_, err = h.toolRepo.CheckoutTool(toolID, user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to checkout tool</p>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkout", user.ID)

	html := "<section aria-live=\"polite\">" +
		"<p>Successfully checked out \"" + tool.Name + "\"!</p>" +
		"<button hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// CheckinTool handles tool checkin
func (h *Handler) CheckinTool(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	checkoutIDStr := c.PostForm("checkout_id")

	checkoutID, err := strconv.Atoi(checkoutIDStr)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8",
			[]byte("<p>Invalid checkout ID</p>"))
		return
	}

	// Check in tool
	err = h.toolRepo.CheckinTool(checkoutID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<p>Failed to check in tool: "+err.Error()+"</p>"))
		return
	}

	// Log event
	h.eventRepo.CreateEvent("tool", "checkin", user.ID)

	html := "<section aria-live=\"polite\">" +
		"<p>Tool checked in successfully!</p>" +
		"<button hx-get=\"/api/tools/checkouts\" hx-target=\"#active-checkouts\">" +
		"Refresh Checkouts" +
		"</button>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}