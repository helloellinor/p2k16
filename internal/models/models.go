package models

import (
	"database/sql"
	"time"
)

// Account represents a user account in the system
type Account struct {
	ID                 int            `json:"id"`
	Username           string         `json:"username"`
	Email              string         `json:"email"`
	Password           string         `json:"-"` // Never expose password in JSON
	Name               sql.NullString `json:"name"`
	Phone              sql.NullString `json:"phone"`
	ResetToken         sql.NullString `json:"-"`
	ResetTokenValidity sql.NullTime   `json:"-"`
	System             sql.NullBool   `json:"system"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedBy          sql.NullInt64  `json:"created_by"`
	UpdatedBy          sql.NullInt64  `json:"updated_by"`
}

// Circle represents a group/circle in the system
type Circle struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   sql.NullInt64 `json:"created_by"`
	UpdatedBy   sql.NullInt64 `json:"updated_by"`
}

// CircleMember represents membership in a circle
type CircleMember struct {
	ID        int            `json:"id"`
	CircleID  int            `json:"circle_id"`
	AccountID int            `json:"account_id"`
	IssuerID  int            `json:"issuer_id"`
	Comment   sql.NullString `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy sql.NullInt64  `json:"created_by"`
	UpdatedBy sql.NullInt64  `json:"updated_by"`
}

// BadgeDescription represents a badge type/template
type BadgeDescription struct {
	ID                    int            `json:"id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	CertificationCircleID sql.NullInt64  `json:"certification_circle_id"`
	Slug                  sql.NullString `json:"slug"`
	Icon                  sql.NullString `json:"icon"`
	Color                 sql.NullString `json:"color"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	CreatedBy             sql.NullInt64  `json:"created_by"`
	UpdatedBy             sql.NullInt64  `json:"updated_by"`
}

// AccountBadge represents a badge awarded to an account
type AccountBadge struct {
	ID                 int           `json:"id"`
	AccountID          int           `json:"account_id"`
	BadgeDescriptionID int           `json:"badge_description_id"`
	AwardedByID        sql.NullInt64 `json:"awarded_by_id"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	CreatedBy          sql.NullInt64 `json:"created_by"`
	UpdatedBy          sql.NullInt64 `json:"updated_by"`

	// Relationships (populated when needed)
	Account          *Account          `json:"account,omitempty"`
	BadgeDescription *BadgeDescription `json:"badge_description,omitempty"`
	AwardedBy        *Account          `json:"awarded_by,omitempty"`
}

// Badge represents a competency badge
type Badge struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   sql.NullInt64  `json:"created_by"`
	UpdatedBy   sql.NullInt64  `json:"updated_by"`
}

// BadgeRequest represents a badge given to an account
type BadgeRequest struct {
	ID        int           `json:"id"`
	AccountID int           `json:"account_id"`
	BadgeID   int           `json:"badge_id"`
	AwarderID int           `json:"awarder_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`
}

// Membership represents a membership record
type Membership struct {
	ID               int           `json:"id"`
	AccountID        int           `json:"account_id"`
	FirstMembership  time.Time     `json:"first_membership"`
	StartMembership  time.Time     `json:"start_membership"`
	Fee              int           `json:"fee"`
	MembershipNumber sql.NullInt64 `json:"membership_number"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	CreatedBy        sql.NullInt64 `json:"created_by"`
	UpdatedBy        sql.NullInt64 `json:"updated_by"`
}

// StripePayment represents a payment made through Stripe
type StripePayment struct {
	ID          int           `json:"id"`
	StripeID    string        `json:"stripe_id"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	Amount      float64       `json:"amount"`
	PaymentDate time.Time     `json:"payment_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   sql.NullInt64 `json:"created_by"`
	UpdatedBy   sql.NullInt64 `json:"updated_by"`
}

// ToolDescription represents a tool in the hackerspace
type ToolDescription struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CircleID    sql.NullInt64  `json:"circle_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedBy   sql.NullInt64  `json:"created_by"`
	UpdatedBy   sql.NullInt64  `json:"updated_by"`
	
	// Relationships
	Circle      *Circle        `json:"circle,omitempty"`
}

// ToolCheckout represents a tool checkout record
type ToolCheckout struct {
	ID         int            `json:"id"`
	ToolID     int            `json:"tool_id"`
	AccountID  int            `json:"account_id"`
	CheckoutAt time.Time      `json:"checkout_at"`
	CheckinAt  sql.NullTime   `json:"checkin_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	CreatedBy  sql.NullInt64  `json:"created_by"`
	UpdatedBy  sql.NullInt64  `json:"updated_by"`

	// Relationships
	Tool    *ToolDescription `json:"tool,omitempty"`
	Account *Account         `json:"account,omitempty"`
}

// Event represents a system event
type Event struct {
	ID        int            `json:"id"`
	Domain    string         `json:"domain"`
	Key       string         `json:"key"`
	Text1     sql.NullString `json:"text1"`
	Text2     sql.NullString `json:"text2"`
	Text3     sql.NullString `json:"text3"`
	Int1      sql.NullInt64  `json:"int1"`
	Int2      sql.NullInt64  `json:"int2"`
	Int3      sql.NullInt64  `json:"int3"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy sql.NullInt64  `json:"created_by"`
	UpdatedBy sql.NullInt64  `json:"updated_by"`
}

// Company represents a company in the system
type Company struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Active    bool          `json:"active"`
	ContactID int           `json:"contact_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`

	// Relationships
	Contact   *Account           `json:"contact,omitempty"`
	Employees []CompanyEmployee  `json:"employees,omitempty"`
}

// CompanyEmployee represents an employee relationship
type CompanyEmployee struct {
	ID        int           `json:"id"`
	CompanyID int           `json:"company_id"`
	AccountID int           `json:"account_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`

	// Relationships
	Company *Company `json:"company,omitempty"`
	Account *Account `json:"account,omitempty"`
}

// StripeCustomer represents a Stripe customer record
type StripeCustomer struct {
	ID           int           `json:"id"`
	AccountID    int           `json:"account_id"`
	StripeID     string        `json:"stripe_id"`
	CustomerData string        `json:"customer_data"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	CreatedBy    sql.NullInt64 `json:"created_by"`
	UpdatedBy    sql.NullInt64 `json:"updated_by"`
}

// Door represents a door configuration
type Door struct {
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	OpenTime   int      `json:"open_time"`   // Duration in seconds
	Type       string   `json:"type"`        // "mqtt" or "dlock"
	Topic      string   `json:"topic"`       // For MQTT doors
	URL        string   `json:"url"`         // For dlock doors
	CircleIDs  []int    `json:"circle_ids"`  // Required circles for access
	Circles    []Circle `json:"circles,omitempty"`
}

// DoorAccess represents a door access log entry
type DoorAccess struct {
	ID        int           `json:"id"`
	AccountID int           `json:"account_id"`
	DoorKey   string        `json:"door_key"`
	OpenedAt  time.Time     `json:"opened_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedBy sql.NullInt64 `json:"created_by"`
	UpdatedBy sql.NullInt64 `json:"updated_by"`

	// Relationships
	Account *Account `json:"account,omitempty"`
}
