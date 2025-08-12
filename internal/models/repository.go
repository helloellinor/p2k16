package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// FindByID retrieves an account by ID
func (r *AccountRepository) FindByID(id int) (*Account, error) {
	query := `
		SELECT id, username, email, password, name, phone, reset_token, 
		       reset_token_validity, system, created_at, updated_at, created_by, updated_by
		FROM account WHERE id = $1`

	account := &Account{}
	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.Username, &account.Email, &account.Password,
		&account.Name, &account.Phone, &account.ResetToken, &account.ResetTokenValidity,
		&account.System, &account.CreatedAt, &account.UpdatedAt, &account.CreatedBy, &account.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// FindByUsername retrieves an account by username
func (r *AccountRepository) FindByUsername(username string) (*Account, error) {
	query := `
		SELECT id, username, email, password, name, phone, reset_token, 
		       reset_token_validity, system, created_at, updated_at, created_by, updated_by
		FROM account WHERE username = $1`

	account := &Account{}
	err := r.db.QueryRow(query, username).Scan(
		&account.ID, &account.Username, &account.Email, &account.Password,
		&account.Name, &account.Phone, &account.ResetToken, &account.ResetTokenValidity,
		&account.System, &account.CreatedAt, &account.UpdatedAt, &account.CreatedBy, &account.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// ValidatePassword checks if the provided password matches the account's password
func (a *Account) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CircleRepository handles database operations for circles
type CircleRepository struct {
	db *sql.DB
}

func NewCircleRepository(db *sql.DB) *CircleRepository {
	return &CircleRepository{db: db}
}

// GetAll retrieves all circles
func (r *CircleRepository) GetAll() ([]Circle, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, created_by, updated_by
		FROM circle ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var circles []Circle
	for rows.Next() {
		var circle Circle
		err := rows.Scan(
			&circle.ID, &circle.Name, &circle.Description,
			&circle.CreatedAt, &circle.UpdatedAt, &circle.CreatedBy, &circle.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		circles = append(circles, circle)
	}

	return circles, nil
}

// BadgeRepository handles database operations for badges
type BadgeRepository struct {
	db *sql.DB
}

func NewBadgeRepository(db *sql.DB) *BadgeRepository {
	return &BadgeRepository{db: db}
}

// GetAllDescriptions retrieves all badge descriptions
func (r *BadgeRepository) GetAllDescriptions() ([]BadgeDescription, error) {
	query := `
		SELECT id, title, description, certification_circle, slug, icon, color,
		       created_at, updated_at, created_by, updated_by
		FROM badge_description ORDER BY title`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var descriptions []BadgeDescription
	for rows.Next() {
		var desc BadgeDescription
		err := rows.Scan(
			&desc.ID, &desc.Title, &desc.Description, &desc.CertificationCircleID,
			&desc.Slug, &desc.Icon, &desc.Color,
			&desc.CreatedAt, &desc.UpdatedAt, &desc.CreatedBy, &desc.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		descriptions = append(descriptions, desc)
	}

	return descriptions, nil
}

// GetBadgesForAccount retrieves all badges for a specific account
func (r *BadgeRepository) GetBadgesForAccount(accountID int) ([]AccountBadge, error) {
	query := `
		SELECT ab.id, ab.account, ab.badge_description, ab.awarded_by,
		       ab.created_at, ab.updated_at, ab.created_by, ab.updated_by,
		       bd.title, bd.description, bd.icon, bd.color
		FROM account_badge ab
		JOIN badge_description bd ON ab.badge_description = bd.id
		WHERE ab.account = $1
		ORDER BY ab.created_at DESC`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []AccountBadge
	for rows.Next() {
		var badge AccountBadge
		var desc BadgeDescription
		err := rows.Scan(
			&badge.ID, &badge.AccountID, &badge.BadgeDescriptionID, &badge.AwardedByID,
			&badge.CreatedAt, &badge.UpdatedAt, &badge.CreatedBy, &badge.UpdatedBy,
			&desc.Title, &desc.Description, &desc.Icon, &desc.Color,
		)
		if err != nil {
			return nil, err
		}
		desc.ID = badge.BadgeDescriptionID
		badge.BadgeDescription = &desc
		badges = append(badges, badge)
	}

	return badges, nil
}

// CreateBadgeDescription creates a new badge description
func (r *BadgeRepository) CreateBadgeDescription(title string, createdBy int) (*BadgeDescription, error) {
	query := `
		INSERT INTO badge_description (title, created_at, updated_at, created_by, updated_by)
		VALUES ($1, NOW(), NOW(), $2, $2)
		RETURNING id, created_at, updated_at`

	var desc BadgeDescription
	desc.Title = title
	desc.CreatedBy = sql.NullInt64{Int64: int64(createdBy), Valid: true}
	desc.UpdatedBy = sql.NullInt64{Int64: int64(createdBy), Valid: true}

	err := r.db.QueryRow(query, title, createdBy).Scan(
		&desc.ID, &desc.CreatedAt, &desc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &desc, nil
}

// AwardBadge awards a badge to an account
func (r *BadgeRepository) AwardBadge(accountID int, badgeDescriptionID int, awardedBy int) (*AccountBadge, error) {
	query := `
		INSERT INTO account_badge (account, badge_description, awarded_by, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, NOW(), NOW(), $3, $3)
		RETURNING id, created_at, updated_at`

	var badge AccountBadge
	badge.AccountID = accountID
	badge.BadgeDescriptionID = badgeDescriptionID
	badge.AwardedByID = sql.NullInt64{Int64: int64(awardedBy), Valid: true}
	badge.CreatedBy = sql.NullInt64{Int64: int64(awardedBy), Valid: true}
	badge.UpdatedBy = sql.NullInt64{Int64: int64(awardedBy), Valid: true}

	err := r.db.QueryRow(query, accountID, badgeDescriptionID, awardedBy).Scan(
		&badge.ID, &badge.CreatedAt, &badge.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &badge, nil
}

// FindBadgeDescriptionByTitle finds a badge description by title
func (r *BadgeRepository) FindBadgeDescriptionByTitle(title string) (*BadgeDescription, error) {
	query := `
		SELECT id, title, description, certification_circle, slug, icon, color,
		       created_at, updated_at, created_by, updated_by
		FROM badge_description WHERE title = $1`

	var desc BadgeDescription
	err := r.db.QueryRow(query, title).Scan(
		&desc.ID, &desc.Title, &desc.Description, &desc.CertificationCircleID,
		&desc.Slug, &desc.Icon, &desc.Color,
		&desc.CreatedAt, &desc.UpdatedAt, &desc.CreatedBy, &desc.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &desc, nil
}
// ToolRepository handles database operations for tools
type ToolRepository struct {
	db *sql.DB
}

func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

// GetAllTools retrieves all tool descriptions
func (r *ToolRepository) GetAllTools() ([]ToolDescription, error) {
	query := `
		SELECT id, name, type, created_at, updated_at, created_by, updated_by
		FROM tool_description ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []ToolDescription
	for rows.Next() {
		var tool ToolDescription
		err := rows.Scan(
			&tool.ID, &tool.Name, &tool.Type,
			&tool.CreatedAt, &tool.UpdatedAt, &tool.CreatedBy, &tool.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, nil
}

// FindToolByID retrieves a tool by ID
func (r *ToolRepository) FindToolByID(id int) (*ToolDescription, error) {
	query := `
		SELECT id, name, type, created_at, updated_at, created_by, updated_by
		FROM tool_description WHERE id = $1`

	var tool ToolDescription
	err := r.db.QueryRow(query, id).Scan(
		&tool.ID, &tool.Name, &tool.Type,
		&tool.CreatedAt, &tool.UpdatedAt, &tool.CreatedBy, &tool.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &tool, nil
}

// CheckoutTool creates a tool checkout record
func (r *ToolRepository) CheckoutTool(toolID int, accountID int) (*ToolCheckout, error) {
	query := `
		INSERT INTO tool_checkout (tool, account, checkout_at, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, NOW(), NOW(), NOW(), $2, $2)
		RETURNING id, checkout_at, created_at, updated_at`

	var checkout ToolCheckout
	checkout.ToolID = toolID
	checkout.AccountID = accountID
	checkout.CreatedBy = sql.NullInt64{Int64: int64(accountID), Valid: true}
	checkout.UpdatedBy = sql.NullInt64{Int64: int64(accountID), Valid: true}

	err := r.db.QueryRow(query, toolID, accountID).Scan(
		&checkout.ID, &checkout.CheckoutAt, &checkout.CreatedAt, &checkout.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &checkout, nil
}

// CheckinTool updates a tool checkout record with checkin time
func (r *ToolRepository) CheckinTool(checkoutID int) error {
	query := `
		UPDATE tool_checkout 
		SET checkin_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND checkin_at IS NULL`

	result, err := r.db.Exec(query, checkoutID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tool checkout not found or already checked in")
	}

	return nil
}

// GetActiveCheckouts retrieves all currently checked out tools
func (r *ToolRepository) GetActiveCheckouts() ([]ToolCheckout, error) {
	query := `
		SELECT tc.id, tc.tool, tc.account, tc.checkout_at, tc.checkin_at,
		       tc.created_at, tc.updated_at, tc.created_by, tc.updated_by,
		       td.name, td.type,
		       a.username, a.name
		FROM tool_checkout tc
		JOIN tool_description td ON tc.tool = td.id
		JOIN account a ON tc.account = a.id
		WHERE tc.checkin_at IS NULL
		ORDER BY tc.checkout_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkouts []ToolCheckout
	for rows.Next() {
		var checkout ToolCheckout
		var tool ToolDescription
		var account Account
		err := rows.Scan(
			&checkout.ID, &checkout.ToolID, &checkout.AccountID, &checkout.CheckoutAt, &checkout.CheckinAt,
			&checkout.CreatedAt, &checkout.UpdatedAt, &checkout.CreatedBy, &checkout.UpdatedBy,
			&tool.Name, &tool.Type,
			&account.Username, &account.Name,
		)
		if err != nil {
			return nil, err
		}
		tool.ID = checkout.ToolID
		account.ID = checkout.AccountID
		checkout.Tool = &tool
		checkout.Account = &account
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// EventRepository handles database operations for events
type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// CreateEvent creates a new event record
func (r *EventRepository) CreateEvent(domain, key string, createdBy int) (*Event, error) {
	query := `
		INSERT INTO event (domain, key, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, NOW(), NOW(), $3, $3)
		RETURNING id, created_at, updated_at`

	var event Event
	event.Domain = domain
	event.Key = key
	event.CreatedBy = sql.NullInt64{Int64: int64(createdBy), Valid: true}
	event.UpdatedBy = sql.NullInt64{Int64: int64(createdBy), Valid: true}

	err := r.db.QueryRow(query, domain, key, createdBy).Scan(
		&event.ID, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetRecentEvents retrieves recent events for a domain
func (r *EventRepository) GetRecentEvents(domain string, limit int) ([]Event, error) {
	query := `
		SELECT id, domain, key, text1, text2, text3, int1, int2, int3,
		       created_at, updated_at, created_by, updated_by
		FROM event 
		WHERE domain = $1 
		ORDER BY created_at DESC 
		LIMIT $2`

	rows, err := r.db.Query(query, domain, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID, &event.Domain, &event.Key, &event.Text1, &event.Text2, &event.Text3,
			&event.Int1, &event.Int2, &event.Int3,
			&event.CreatedAt, &event.UpdatedAt, &event.CreatedBy, &event.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
