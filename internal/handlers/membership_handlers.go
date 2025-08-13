package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/middleware"
)

// GetMembershipStatus returns the membership status for current user
func (h *Handler) GetMembershipStatus(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	// Check if user is an active member
	isActive, err := h.membershipRepo.IsActiveMember(user.ID)
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to check membership status</div>"))
		return
	}

	// Check payment status
	isPaying, _ := h.membershipRepo.IsAccountPayingMember(user.ID)
	isEmployee, _ := h.membershipRepo.IsAccountCompanyEmployee(user.ID)

	// Get membership details
	membership, _ := h.membershipRepo.GetMembershipByAccount(user.ID)

	html := "<section aria-labelledby=\"membership-title\">" +
		"<h2 id=\"membership-title\">Membership Status</h2>" +
		"<div>"

	if isActive {
		html += "<p>" +
			"Active Member"
		if isPaying {
			html += " (Paying Member)"
		}
		if isEmployee {
			html += " (Company Employee)"
		}
		html += "</p>"
	} else {
		html += "<p>Inactive Member</p>"
	}

	if membership != nil {
		html += "<section>" +
			"<h3>Membership Details</h3>" +
			"<p>Member since: " + membership.FirstMembership.Format("2006-01-02") + "</p>" +
			"<p>Current membership start: " + membership.StartMembership.Format("2006-01-02") + "</p>" +
			"<p>Monthly fee: " + fmt.Sprintf("%.2f NOK", float64(membership.Fee)/100) + "</p>"
		if membership.MembershipNumber.Valid {
			html += "<p><strong>Membership number:</strong> " + fmt.Sprintf("%d", membership.MembershipNumber.Int64) + "</p>"
		}
		html += "</section>"
	}

	html += "</div>" +
		"</section>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// GetMembershipStatusAPI returns the membership status for current user as JSON (API endpoint: GET /api/memberships/)
func (h *Handler) GetMembershipStatusAPI(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	// Check if user is an active member
	isActive, err := h.membershipRepo.IsActiveMember(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to check membership status",
		})
		return
	}

	// Check payment status
	isPaying, _ := h.membershipRepo.IsAccountPayingMember(user.ID)
	isEmployee, _ := h.membershipRepo.IsAccountCompanyEmployee(user.ID)

	// Get membership details
	membership, _ := h.membershipRepo.GetMembershipByAccount(user.ID)

	response := gin.H{
		"status": "success",
		"data": gin.H{
			"is_active":  isActive,
			"is_paying":  isPaying,
			"is_employee": isEmployee,
		},
	}

	if membership != nil {
		membershipData := gin.H{
			"first_membership":  membership.FirstMembership.Format("2006-01-02"),
			"start_membership":  membership.StartMembership.Format("2006-01-02"),
			"fee":               membership.Fee,
		}
		if membership.MembershipNumber.Valid {
			membershipData["membership_number"] = membership.MembershipNumber.Int64
		}
		response["data"].(gin.H)["membership"] = membershipData
	}

	c.JSON(http.StatusOK, response)
}

// GetActiveMembers returns a simple list of active members (API endpoint)
func (h *Handler) GetActiveMembers(c *gin.Context) {
	payingMembers, err := h.membershipRepo.GetActivePayingMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load active members"})
		return
	}

	// Return simple JSON response
	var members []map[string]interface{}
	for _, member := range payingMembers {
		displayName := member.Username
		if member.Name.Valid && member.Name.String != "" {
			displayName = member.Name.String
		}
		members = append(members, map[string]interface{}{
			"id":       member.ID,
			"username": member.Username,
			"name":     displayName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// GetActiveMembersDetailed returns a detailed list of active members
func (h *Handler) GetActiveMembersDetailed(c *gin.Context) {
	payingMembers, err := h.membershipRepo.GetActivePayingMembers()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load active members</div>"))
		return
	}

	activeCompanies, err := h.membershipRepo.GetActiveCompanies()
	if err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<div class=\"alert alert-danger\">Failed to load companies</div>"))
		return
	}

	html := "<div class=\"card\">" +
		"<div class=\"card-header\">" +
		"<h6>Active Members</h6>" +
		"</div>" +
		"<div class=\"card-body\">"

	// Show paying members
	if len(payingMembers) > 0 {
		html += "<h6 class=\"text-success\">Paying Members (" + fmt.Sprintf("%d", len(payingMembers)) + ")</h6>" +
			"<div class=\"row\">"
		for _, member := range payingMembers {
			displayName := member.Username
			if member.Name.Valid && member.Name.String != "" {
				displayName = member.Name.String + " (" + member.Username + ")"
			}
			html += "<div class=\"col-md-6 mb-2\">" +
				"<span class=\"badge bg-success\">" + displayName + "</span>" +
				"</div>"
		}
		html += "</div>"
	}

	// Show companies
	if len(activeCompanies) > 0 {
		html += "<h6 class=\"text-primary mt-3\">Active Companies (" + fmt.Sprintf("%d", len(activeCompanies)) + ")</h6>" +
			"<div class=\"list-group\">"
		for _, company := range activeCompanies {
			contactName := company.Contact.Username
			if company.Contact.Name.Valid && company.Contact.Name.String != "" {
				contactName = company.Contact.Name.String
			}
			html += "<div class=\"list-group-item\">" +
				"<h6 class=\"mb-1\">" + company.Name + "</h6>" +
				"<p class=\"mb-1\">Contact: " + contactName + "</p>" +
				"</div>"
		}
		html += "</div>"
	}

	if len(payingMembers) == 0 && len(activeCompanies) == 0 {
		html += "<p class=\"text-muted\">No active members found.</p>"
	}

	html += "</div>" +
		"</div>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}