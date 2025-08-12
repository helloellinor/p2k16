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
