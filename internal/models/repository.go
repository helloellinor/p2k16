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

// SetPassword hashes and sets a new password for the account
func (a *Account) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
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

// AccountHasBadge checks if an account has a specific badge
func (r *BadgeRepository) AccountHasBadge(accountID int, badgeDescriptionID int) (bool, error) {
	query := `
		SELECT COUNT(*) > 0 
		FROM account_badge 
		WHERE account = $1 AND badge_description = $2`
	
	var has bool
	err := r.db.QueryRow(query, accountID, badgeDescriptionID).Scan(&has)
	if err != nil {
		return false, err
	}
	
	return has, nil
}

// DeleteAccountBadge removes a badge from an account
func (r *BadgeRepository) DeleteAccountBadge(accountBadgeID int, accountID int) error {
	query := `
		DELETE FROM account_badge 
		WHERE id = $1 AND account = $2`
	
	result, err := r.db.Exec(query, accountBadgeID, accountID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no badge found with id %d for account %d", accountBadgeID, accountID)
	}
	
	return nil
}

// GetAll retrieves all badge descriptions (alias for GetAllDescriptions)
func (r *BadgeRepository) GetAll() ([]BadgeDescription, error) {
	return r.GetAllDescriptions()
}

// GetUserBadges retrieves all badges for a specific user
func (r *BadgeRepository) GetUserBadges(accountID int) ([]BadgeDescription, error) {
	query := `
		SELECT bd.id, bd.title, bd.description, bd.certification_circle, bd.slug, bd.icon, bd.color,
		       bd.created_at, bd.updated_at, bd.created_by, bd.updated_by
		FROM badge_description bd
		JOIN account_badge ab ON bd.id = ab.badge_description
		WHERE ab.account = $1
		ORDER BY bd.title`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []BadgeDescription
	for rows.Next() {
		var badge BadgeDescription
		err := rows.Scan(
			&badge.ID, &badge.Title, &badge.Description, &badge.CertificationCircleID, &badge.Slug,
			&badge.Icon, &badge.Color, &badge.CreatedAt, &badge.UpdatedAt, &badge.CreatedBy, &badge.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}

	return badges, nil
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
		SELECT td.id, td.name, td.description, td.circle, td.created_at, td.updated_at, td.created_by, td.updated_by,
		       c.id, c.name, c.description
		FROM tool_description td
		LEFT JOIN circle c ON td.circle = c.id
		ORDER BY td.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []ToolDescription
	for rows.Next() {
		var tool ToolDescription
		var circle Circle
		var circleID sql.NullInt64
		var circleName sql.NullString
		var circleDesc sql.NullString

		err := rows.Scan(
			&tool.ID, &tool.Name, &tool.Description, &tool.CircleID,
			&tool.CreatedAt, &tool.UpdatedAt, &tool.CreatedBy, &tool.UpdatedBy,
			&circleID, &circleName, &circleDesc,
		)
		if err != nil {
			return nil, err
		}

		// Set the circle relationship if it exists
		if circleID.Valid {
			circle.ID = int(circleID.Int64)
			circle.Name = circleName.String
			circle.Description = circleDesc.String
			tool.Circle = &circle
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

// FindToolByID retrieves a tool by ID
func (r *ToolRepository) FindToolByID(id int) (*ToolDescription, error) {
	query := `
		SELECT td.id, td.name, td.description, td.circle, td.created_at, td.updated_at, td.created_by, td.updated_by,
		       c.id, c.name, c.description
		FROM tool_description td
		LEFT JOIN circle c ON td.circle = c.id
		WHERE td.id = $1`

	var tool ToolDescription
	var circle Circle
	var circleID sql.NullInt64
	var circleName sql.NullString
	var circleDesc sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&tool.ID, &tool.Name, &tool.Description, &tool.CircleID,
		&tool.CreatedAt, &tool.UpdatedAt, &tool.CreatedBy, &tool.UpdatedBy,
		&circleID, &circleName, &circleDesc,
	)

	if err != nil {
		return nil, err
	}

	// Set the circle relationship if it exists
	if circleID.Valid {
		circle.ID = int(circleID.Int64)
		circle.Name = circleName.String
		circle.Description = circleDesc.String
		tool.Circle = &circle
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
		       td.name, td.description,
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
			&tool.Name, &tool.Description,
			&account.Username, &account.Name,
		)
		if err != nil {
			return nil, err
		}

		checkout.Tool = &tool
		checkout.Account = &account
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetCheckoutByID retrieves a specific checkout by ID
func (r *ToolRepository) GetCheckoutByID(checkoutID int) (*ToolCheckout, error) {
	query := `
		SELECT tc.id, tc.tool, tc.account, tc.checkout_at, tc.checkin_at,
		       tc.created_at, tc.updated_at, tc.created_by, tc.updated_by,
		       td.name, td.description,
		       a.username, a.name
		FROM tool_checkout tc
		JOIN tool_description td ON tc.tool = td.id
		JOIN account a ON tc.account = a.id
		WHERE tc.id = $1`

	var checkout ToolCheckout
	var tool ToolDescription
	var account Account
	
	err := r.db.QueryRow(query, checkoutID).Scan(
		&checkout.ID, &checkout.ToolID, &checkout.AccountID, &checkout.CheckoutAt, &checkout.CheckinAt,
		&checkout.CreatedAt, &checkout.UpdatedAt, &checkout.CreatedBy, &checkout.UpdatedBy,
		&tool.Name, &tool.Description,
		&account.Username, &account.Name,
	)
	if err != nil {
		return nil, err
	}

	checkout.Tool = &tool
	checkout.Account = &account
	return &checkout, nil
}

// CreateTool creates a new tool description
func (r *ToolRepository) CreateTool(name, description string, circleID *int, userID int) (*ToolDescription, error) {
	query := `
		INSERT INTO tool_description (name, description, circle, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, NOW(), NOW(), $4, $4)
		RETURNING id, created_at, updated_at`

	var tool ToolDescription
	tool.Name = name
	if description != "" {
		tool.Description = sql.NullString{String: description, Valid: true}
	}
	if circleID != nil {
		tool.CircleID = sql.NullInt64{Int64: int64(*circleID), Valid: true}
	}
	tool.CreatedBy = sql.NullInt64{Int64: int64(userID), Valid: true}
	tool.UpdatedBy = sql.NullInt64{Int64: int64(userID), Valid: true}

	err := r.db.QueryRow(query, name, description, circleID, userID).Scan(
		&tool.ID, &tool.CreatedAt, &tool.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tool, nil
}

// UpdateTool updates an existing tool description
func (r *ToolRepository) UpdateTool(id int, name, description string, circleID *int, userID int) (*ToolDescription, error) {
	query := `
		UPDATE tool_description 
		SET name = $2, description = $3, circle = $4, updated_at = NOW(), updated_by = $5
		WHERE id = $1
		RETURNING created_at, updated_at, created_by`

	var tool ToolDescription
	tool.ID = id
	tool.Name = name
	if description != "" {
		tool.Description = sql.NullString{String: description, Valid: true}
	}
	if circleID != nil {
		tool.CircleID = sql.NullInt64{Int64: int64(*circleID), Valid: true}
	}
	tool.UpdatedBy = sql.NullInt64{Int64: int64(userID), Valid: true}

	err := r.db.QueryRow(query, id, name, description, circleID, userID).Scan(
		&tool.CreatedAt, &tool.UpdatedAt, &tool.CreatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &tool, nil
}

// DeleteTool deletes a tool description
func (r *ToolRepository) DeleteTool(id int) error {
	// First check if there are any active checkouts for this tool
	checkQuery := `
		SELECT COUNT(*) FROM tool_checkout 
		WHERE tool = $1 AND checkin_at IS NULL`

	var count int
	err := r.db.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete tool: it has active checkouts")
	}

	deleteQuery := `DELETE FROM tool_description WHERE id = $1`
	result, err := r.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tool not found")
	}

	return nil
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

// MembershipRepository handles database operations for memberships
type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(db *sql.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

// GetMembershipByAccount retrieves membership info for an account
func (r *MembershipRepository) GetMembershipByAccount(accountID int) (*Membership, error) {
	query := `
		SELECT id, account, first_membership, start_membership, fee, membership_number,
		       created_at, updated_at, created_by, updated_by
		FROM membership WHERE account = $1`

	var membership Membership
	err := r.db.QueryRow(query, accountID).Scan(
		&membership.ID, &membership.AccountID, &membership.FirstMembership, &membership.StartMembership,
		&membership.Fee, &membership.MembershipNumber,
		&membership.CreatedAt, &membership.UpdatedAt, &membership.CreatedBy, &membership.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &membership, nil
}

// IsAccountPayingMember checks if an account has an active Stripe payment
func (r *MembershipRepository) IsAccountPayingMember(accountID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM stripe_payment 
		WHERE created_by = $1 AND end_date >= NOW() - INTERVAL '1 day'`

	var count int
	err := r.db.QueryRow(query, accountID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsAccountCompanyEmployee checks if an account is employed by an active company
func (r *MembershipRepository) IsAccountCompanyEmployee(accountID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM company_employee ce
		JOIN company c ON ce.company = c.id
		WHERE ce.account = $1 AND c.active = true`

	var count int
	err := r.db.QueryRow(query, accountID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsActiveMember checks if an account is an active member (paying or company employee)
func (r *MembershipRepository) IsActiveMember(accountID int) (bool, error) {
	// Check if paying member
	isPaying, err := r.IsAccountPayingMember(accountID)
	if err != nil {
		return false, err
	}
	if isPaying {
		return true, nil
	}

	// Check if company employee
	isEmployee, err := r.IsAccountCompanyEmployee(accountID)
	if err != nil {
		return false, err
	}

	return isEmployee, nil
}

// GetActivePayingMembers retrieves all accounts with active payments
func (r *MembershipRepository) GetActivePayingMembers() ([]Account, error) {
	query := `
		SELECT DISTINCT a.id, a.username, a.email, a.password, a.name, a.phone, 
		       a.reset_token, a.reset_token_validity, a.system,
		       a.created_at, a.updated_at, a.created_by, a.updated_by
		FROM account a
		JOIN stripe_payment sp ON sp.created_by = a.id
		WHERE sp.end_date >= NOW() - INTERVAL '1 day'
		ORDER BY a.username`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID, &account.Username, &account.Email, &account.Password,
			&account.Name, &account.Phone, &account.ResetToken, &account.ResetTokenValidity,
			&account.System, &account.CreatedAt, &account.UpdatedAt, &account.CreatedBy, &account.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetActiveMembers retrieves all active members (alias for GetActivePayingMembers)
func (r *MembershipRepository) GetActiveMembers() ([]Account, error) {
	return r.GetActivePayingMembers()
}

// GetActiveCompanies retrieves all active companies
func (r *MembershipRepository) GetActiveCompanies() ([]Company, error) {
	query := `
		SELECT c.id, c.name, c.active, c.contact,
		       c.created_at, c.updated_at, c.created_by, c.updated_by,
		       a.username, a.name
		FROM company c
		JOIN account a ON c.contact = a.id
		WHERE c.active = true
		ORDER BY c.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var company Company
		var contact Account
		err := rows.Scan(
			&company.ID, &company.Name, &company.Active, &company.ContactID,
			&company.CreatedAt, &company.UpdatedAt, &company.CreatedBy, &company.UpdatedBy,
			&contact.Username, &contact.Name,
		)
		if err != nil {
			return nil, err
		}
		contact.ID = company.ContactID
		company.Contact = &contact
		companies = append(companies, company)
	}

	return companies, nil
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

// DoorRepository handles database operations for doors and access
type DoorRepository struct {
	db *sql.DB
}

func NewDoorRepository(db *sql.DB) *DoorRepository {
	return &DoorRepository{db: db}
}

// GetConfiguredDoors returns hardcoded door configurations
// In a real implementation, this would come from a config file or database
func (r *DoorRepository) GetConfiguredDoors() []Door {
	return []Door{
		{
			Key:       "main",
			Name:      "Main Door",
			OpenTime:  5,
			Type:      "mqtt",
			Topic:     "bitraf/door/main",
			CircleIDs: []int{1}, // Admin circle
		},
		{
			Key:       "workshop",
			Name:      "Workshop Door",
			OpenTime:  5,
			Type:      "mqtt",
			Topic:     "bitraf/door/workshop",
			CircleIDs: []int{1, 2}, // Admin and member circles
		},
		{
			Key:       "storage",
			Name:      "Storage Room",
			OpenTime:  3,
			Type:      "dlock",
			URL:       "http://storage-lock.local",
			CircleIDs: []int{1}, // Admin only
		},
	}
}

// CanAccessDoor checks if an account can access a specific door
func (r *DoorRepository) CanAccessDoor(accountID int, door Door, membershipRepo *MembershipRepository) (bool, error) {
	// Check if account is company employee (has access to all doors)
	isEmployee, err := membershipRepo.IsAccountCompanyEmployee(accountID)
	if err != nil {
		return false, err
	}
	if isEmployee {
		return true, nil
	}

	// Check if account is paying member
	isPaying, err := membershipRepo.IsAccountPayingMember(accountID)
	if err != nil {
		return false, err
	}
	if !isPaying {
		return false, nil // Must be paying member for door access
	}

	// Check circle membership if required
	if len(door.CircleIDs) > 0 {
		for _, circleID := range door.CircleIDs {
			isMember, err := r.IsAccountInCircle(accountID, circleID)
			if err != nil {
				return false, err
			}
			if isMember {
				return true, nil
			}
		}
		return false, nil // Not in any required circle
	}

	return true, nil // Paying member with no circle requirements
}

// IsAccountInCircle checks if an account is a member of a circle
func (r *DoorRepository) IsAccountInCircle(accountID int, circleID int) (bool, error) {
	query := `SELECT COUNT(*) FROM circle_member WHERE account = $1 AND circle = $2`

	var count int
	err := r.db.QueryRow(query, accountID, circleID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetAccessibleDoors returns doors that an account can access
func (r *DoorRepository) GetAccessibleDoors(accountID int, membershipRepo *MembershipRepository) ([]Door, error) {
	allDoors := r.GetConfiguredDoors()
	var accessibleDoors []Door

	for _, door := range allDoors {
		canAccess, err := r.CanAccessDoor(accountID, door, membershipRepo)
		if err != nil {
			return nil, err
		}
		if canAccess {
			accessibleDoors = append(accessibleDoors, door)
		}
	}

	return accessibleDoors, nil
}

// LogDoorAccess records a door access event
func (r *DoorRepository) LogDoorAccess(accountID int, doorKey string) (*DoorAccess, error) {
	query := `
		INSERT INTO door_access (account, door_key, opened_at, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, NOW(), NOW(), NOW(), $1, $1)
		RETURNING id, opened_at, created_at, updated_at`

	var access DoorAccess
	access.AccountID = accountID
	access.DoorKey = doorKey
	access.CreatedBy = sql.NullInt64{Int64: int64(accountID), Valid: true}
	access.UpdatedBy = sql.NullInt64{Int64: int64(accountID), Valid: true}

	err := r.db.QueryRow(query, accountID, doorKey).Scan(
		&access.ID, &access.OpenedAt, &access.CreatedAt, &access.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &access, nil
}

// GetRecentDoorAccess returns recent door access events
func (r *DoorRepository) GetRecentDoorAccess(limit int) ([]DoorAccess, error) {
	query := `
		SELECT da.id, da.account, da.door_key, da.opened_at,
		       da.created_at, da.updated_at, da.created_by, da.updated_by,
		       a.username, a.name
		FROM door_access da
		JOIN account a ON da.account = a.id
		ORDER BY da.opened_at DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accesses []DoorAccess
	for rows.Next() {
		var access DoorAccess
		var account Account
		err := rows.Scan(
			&access.ID, &access.AccountID, &access.DoorKey, &access.OpenedAt,
			&access.CreatedAt, &access.UpdatedAt, &access.CreatedBy, &access.UpdatedBy,
			&account.Username, &account.Name,
		)
		if err != nil {
			return nil, err
		}
		account.ID = access.AccountID
		access.Account = &account
		accesses = append(accesses, access)
	}

	return accesses, nil
}

// UpdatePassword updates the password for an account
func (r *AccountRepository) UpdatePassword(accountID int, hashedPassword string) error {
	query := `UPDATE account SET password = $1, updated_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, accountID)
	return err
}

// UpdateEmail updates the email for an account
func (r *AccountRepository) UpdateEmail(accountID int, email string) error {
	query := `UPDATE account SET email = $1, updated_at = now() WHERE id = $2`
	_, err := r.db.Exec(query, email, accountID)
	return err
}

// GetAll retrieves all accounts
func (r *AccountRepository) GetAll() ([]Account, error) {
	query := `
		SELECT id, username, email, phone, name, password, created_at, updated_at, created_by, updated_by
		FROM account ORDER BY username`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID, &account.Username, &account.Email, &account.Phone, &account.Name,
			&account.Password, &account.CreatedAt, &account.UpdatedAt, &account.CreatedBy, &account.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// UpdateProfile updates the profile fields for an account
func (r *AccountRepository) UpdateProfile(account *Account) error {
	query := `
		UPDATE account 
		SET name = $1, phone = $2, updated_at = now() 
		WHERE id = $3`

	var name, phone interface{}
	if account.Name.Valid {
		name = account.Name.String
	}
	if account.Phone.Valid {
		phone = account.Phone.String
	}

	_, err := r.db.Exec(query, name, phone, account.ID)
	return err
}

// GetAllAccounts retrieves all accounts with pagination
func (r *AccountRepository) GetAllAccounts(limit, offset int) ([]Account, error) {
	query := `
		SELECT id, username, email, password, name, phone, reset_token, 
		       reset_token_validity, system, created_at, updated_at, created_by, updated_by
		FROM account 
		ORDER BY username
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID, &account.Username, &account.Email, &account.Password,
			&account.Name, &account.Phone, &account.ResetToken, &account.ResetTokenValidity,
			&account.System, &account.CreatedAt, &account.UpdatedAt, &account.CreatedBy, &account.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// GetAccountCount returns the total number of accounts
func (r *AccountRepository) GetAccountCount() (int, error) {
	query := `SELECT COUNT(*) FROM account`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
