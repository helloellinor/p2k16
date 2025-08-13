package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTMXUtils provides utilities for HTMX responses
type HTMXUtils struct{}

// NewHTMXUtils creates a new HTMX utilities instance
func NewHTMXUtils() *HTMXUtils {
	return &HTMXUtils{}
}

// FormField represents a form field for the htmx-form component
type FormField struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Label       string            `json:"label"`
	Placeholder string            `json:"placeholder"`
	Value       string            `json:"value"`
	Required    bool              `json:"required"`
	Help        string            `json:"help"`
	Options     []FormFieldOption `json:"options,omitempty"`
	Attrs       map[string]string `json:"attrs,omitempty"`
}

// FormFieldOption represents an option for select fields
type FormFieldOption struct {
	Value    string `json:"value"`
	Label    string `json:"label"`
	Selected bool   `json:"selected"`
}

// FormData represents data for the htmx-form component
type FormData struct {
	Action     string      `json:"action"`
	Target     string      `json:"target"`
	Swap       string      `json:"swap,omitempty"`
	Indicator  string      `json:"indicator,omitempty"`
	Ext        string      `json:"ext,omitempty"`
	Class      string      `json:"class,omitempty"`
	Fields     []FormField `json:"fields"`
	SubmitText string      `json:"submitText"`
	CancelUrl  string      `json:"cancelUrl,omitempty"`
}

// AlertData represents data for the alert component
type AlertData struct {
	Type        string `json:"type"`
	Title       string `json:"title,omitempty"`
	Message     string `json:"message"`
	Dismissible bool   `json:"dismissible"`
}

// CardData represents data for the card component
type CardData struct {
	Title    string        `json:"title,omitempty"`
	Subtitle string        `json:"subtitle,omitempty"`
	Content  template.HTML `json:"content"`
	Footer   template.HTML `json:"footer,omitempty"`
	Class    string        `json:"class,omitempty"`
}

// ButtonData represents data for the htmx-button component
type ButtonData struct {
	Text      string `json:"text"`
	Variant   string `json:"variant"`
	Class     string `json:"class,omitempty"`
	Get       string `json:"get,omitempty"`
	Post      string `json:"post,omitempty"`
	Put       string `json:"put,omitempty"`
	Delete    string `json:"delete,omitempty"`
	Target    string `json:"target,omitempty"`
	Swap      string `json:"swap,omitempty"`
	Vals      string `json:"vals,omitempty"`
	Confirm   string `json:"confirm,omitempty"`
	Indicator string `json:"indicator,omitempty"`
	Disabled  bool   `json:"disabled"`
}

// BadgeData represents a badge
type BadgeData struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

// BadgeListData represents data for the badge-list component
type BadgeListData struct {
	Badges       []BadgeData `json:"badges"`
	EmptyMessage string      `json:"emptyMessage"`
}

// TableData represents data for the data-table component
type TableData struct {
	Headers      []string   `json:"headers"`
	Rows         [][]string `json:"rows"`
	EmptyMessage string     `json:"emptyMessage"`
}

// ModalData represents data for the modal component
type ModalData struct {
	ID      string        `json:"id"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
	Footer  template.HTML `json:"footer,omitempty"`
}

// RenderAlert renders an alert component
func (h *HTMXUtils) RenderAlert(c *gin.Context, alertType, title, message string, dismissible bool) {
	alert := AlertData{
		Type:        alertType,
		Title:       title,
		Message:     message,
		Dismissible: dismissible,
	}

	html := h.renderTemplate("alert", alert)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RenderError renders an error alert
func (h *HTMXUtils) RenderError(c *gin.Context, message string) {
	h.RenderAlert(c, "error", "Error", message, true)
}

// RenderSuccess renders a success alert
func (h *HTMXUtils) RenderSuccess(c *gin.Context, message string) {
	h.RenderAlert(c, "success", "Success", message, true)
}

// RenderWarning renders a warning alert
func (h *HTMXUtils) RenderWarning(c *gin.Context, message string) {
	h.RenderAlert(c, "warning", "Warning", message, true)
}

// RenderInfo renders an info alert
func (h *HTMXUtils) RenderInfo(c *gin.Context, message string) {
	h.RenderAlert(c, "info", "Info", message, true)
}

// RenderLoading renders a loading spinner
func (h *HTMXUtils) RenderLoading(c *gin.Context, text string) {
	loading := map[string]string{
		"text": text,
	}

	html := h.renderTemplate("loading-spinner", loading)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RenderForm renders a form component
func (h *HTMXUtils) RenderForm(c *gin.Context, formData FormData) {
	html := h.renderTemplate("htmx-form", formData)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RenderCard renders a card component
func (h *HTMXUtils) RenderCard(c *gin.Context, cardData CardData) {
	html := h.renderTemplate("card", cardData)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RenderBadgeList renders a badge list
func (h *HTMXUtils) RenderBadgeList(c *gin.Context, badges []BadgeData, emptyMessage string) {
	badgeList := BadgeListData{
		Badges:       badges,
		EmptyMessage: emptyMessage,
	}

	html := h.renderTemplate("badge-list", badgeList)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// RenderTable renders a data table
func (h *HTMXUtils) RenderTable(c *gin.Context, headers []string, rows [][]string, emptyMessage string) {
	table := TableData{
		Headers:      headers,
		Rows:         rows,
		EmptyMessage: emptyMessage,
	}

	html := h.renderTemplate("data-table", table)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// IsHTMXRequest checks if the request is from HTMX
func IsHTMXRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}

// IsBoosted checks if the request is boosted by HTMX
func IsBoosted(c *gin.Context) bool {
	return c.GetHeader("HX-Boosted") == "true"
}

// GetHTMXTarget gets the HTMX target from request headers
func GetHTMXTarget(c *gin.Context) string {
	return c.GetHeader("HX-Target")
}

// GetHTMXTrigger gets the HTMX trigger from request headers
func GetHTMXTrigger(c *gin.Context) string {
	return c.GetHeader("HX-Trigger")
}

// SetHTMXTrigger sets an HTMX trigger in the response
func SetHTMXTrigger(c *gin.Context, trigger string) {
	c.Header("HX-Trigger", trigger)
}

// SetHTMXRedirect sets an HTMX redirect in the response
func SetHTMXRedirect(c *gin.Context, url string) {
	c.Header("HX-Redirect", url)
}

// SetHTMXRefresh sets an HTMX refresh in the response
func SetHTMXRefresh(c *gin.Context) {
	c.Header("HX-Refresh", "true")
}

// renderTemplate renders a template with data (simplified version)
func (h *HTMXUtils) renderTemplate(templateName string, data interface{}) string {
	// This is a simplified implementation
	// In a real implementation, you would use Go's template system
	// to load and render the templates from files

	switch templateName {
	case "alert":
		alert := data.(AlertData)
		dismissBtn := ""
		if alert.Dismissible {
			dismissBtn = `<button class="p2k16-alert__close" onclick="this.parentElement.remove()">×</button>`
		}
		title := ""
		if alert.Title != "" {
			title = `<strong>` + alert.Title + `</strong> `
		}
		return `<div class="p2k16-alert p2k16-alert--` + alert.Type + `">` + dismissBtn + title + alert.Message + `</div>`

	case "loading-spinner":
		data := data.(map[string]string)
		text := ""
		if data["text"] != "" {
			text = `<span class="p2k16-loading__text">` + data["text"] + `</span>`
		}
		return `<div class="p2k16-loading"><span class="p2k16-spinner"></span>` + text + `</div>`

	default:
		return `<div class="p2k16-alert p2k16-alert--error">Template not found: ` + templateName + `</div>`
	}
}
