package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// ChiGetTools returns a list of all tools
func (h *ChiHandler) ChiGetTools(w http.ResponseWriter, r *http.Request) {
	tools, err := h.toolRepo.GetAllTools()
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, "<p>Failed to load tools</p>")
		return
	}

	html := "<section aria-labelledby=\"tools-title\">" +
		"<h2 id=\"tools-title\">Available Tools</h2>" +
		"<div class=\"row\">"

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

	h.writeHTML(w, http.StatusOK, html)
}

// ChiGetActiveCheckouts returns currently checked out tools
func (h *ChiHandler) ChiGetActiveCheckouts(w http.ResponseWriter, r *http.Request) {
	checkouts, err := h.toolRepo.GetActiveCheckouts()
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, "<p>Failed to load active checkouts</p>")
		return
	}

	html := "<section aria-labelledby=\"checkouts-title\">" +
		"<h2 id=\"checkouts-title\">Active Tool Checkouts</h2>"

	if len(checkouts) == 0 {
		html += "<p>No tools currently checked out.</p>"
	} else {
		html += "<div class=\"row\">"
		for _, checkout := range checkouts {
			html += fmt.Sprintf("<div class=\"col-md-6 mb-3\">"+
				"<div class=\"card border-warning\">"+
				"<div class=\"card-body\">"+
				"<h6 class=\"card-title\">%s</h6>"+
				"<p class=\"card-text\">Checked out to: %s</p>"+
				"<p class=\"card-text\">Since: %s</p>"+
				"<button class=\"btn btn-warning btn-sm\" "+
				"hx-post=\"/api/tools/checkin\" "+
				"hx-vals='{\"checkout_id\":\"%d\"}' "+
				"hx-target=\"#tool-result\" "+
				"hx-swap=\"innerHTML\">"+
				"Check In"+
				"</button>"+
				"</div>"+
				"</div>"+
				"</div>",
				checkout.Tool.Name, checkout.Account.Username, checkout.CheckoutAt.Format("2006-01-02 15:04"), checkout.ID)
		}
		html += "</div>"
	}

	html += "</section>"
	h.writeHTML(w, http.StatusOK, html)
}

// ChiCheckoutTool handles tool checkout
func (h *ChiHandler) ChiCheckoutTool(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.writeHTML(w, http.StatusBadRequest, "<p>Invalid form data</p>")
		return
	}

	toolIDStr := r.FormValue("tool_id")
	toolID, err := strconv.Atoi(toolIDStr)
	if err != nil {
		h.writeHTML(w, http.StatusBadRequest, "<p>Invalid tool ID</p>")
		return
	}

	// Get current user
	user := ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, "<p>Authentication required</p>")
		return
	}

	// Check if tool exists
	tool, err := h.toolRepo.FindToolByID(toolID)
	if err != nil {
		h.writeHTML(w, http.StatusNotFound, "<p>Tool not found</p>")
		return
	}

	// Check if tool is already checked out
	activeCheckouts, err := h.toolRepo.GetActiveCheckouts()
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, "<p>Failed to check tool availability</p>")
		return
	}
	
	// Check if this tool is already checked out
	for _, checkout := range activeCheckouts {
		if checkout.ToolID == toolID {
			h.writeHTML(w, http.StatusConflict, 
				fmt.Sprintf("<p>Tool '%s' is already checked out to %s</p>", tool.Name, checkout.Account.Username))
			return
		}
	}

	// Create checkout
	_, err = h.toolRepo.CheckoutTool(toolID, user.ID)
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, 
			fmt.Sprintf("<p>Failed to checkout tool '%s': %v</p>", tool.Name, err))
		return
	}

	html := fmt.Sprintf("<div class=\"alert alert-success\">"+
		"✅ Successfully checked out tool '%s'"+
		"</div>", tool.Name)

	h.writeHTML(w, http.StatusOK, html)
}

// ChiCheckinTool handles tool checkin
func (h *ChiHandler) ChiCheckinTool(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.writeHTML(w, http.StatusBadRequest, "<p>Invalid form data</p>")
		return
	}

	checkoutIDStr := r.FormValue("checkout_id")
	checkoutID, err := strconv.Atoi(checkoutIDStr)
	if err != nil {
		h.writeHTML(w, http.StatusBadRequest, "<p>Invalid checkout ID</p>")
		return
	}

	// Get current user
	user := ChiGetCurrentUser(r)
	if user == nil {
		h.writeHTML(w, http.StatusUnauthorized, "<p>Authentication required</p>")
		return
	}

	// Get checkout details
	checkout, err := h.toolRepo.GetCheckoutByID(checkoutID)
	if err != nil {
		h.writeHTML(w, http.StatusNotFound, "<p>Checkout not found</p>")
		return
	}

	// Check if user owns this checkout or is admin
	if checkout.AccountID != user.ID {
		// TODO: Add admin check here
		h.writeHTML(w, http.StatusForbidden, "<p>You can only check in tools you checked out</p>")
		return
	}

	// Check in tool
	err = h.toolRepo.CheckinTool(checkoutID)
	if err != nil {
		h.writeHTML(w, http.StatusInternalServerError, 
			fmt.Sprintf("<p>Failed to check in tool: %v</p>", err))
		return
	}

	html := fmt.Sprintf("<div class=\"alert alert-success\">"+
		"✅ Successfully checked in tool '%s'"+
		"</div>", checkout.Tool.Name)

	h.writeHTML(w, http.StatusOK, html)
}