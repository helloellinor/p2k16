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
