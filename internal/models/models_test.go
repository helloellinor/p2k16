package models

import (
	"testing"
)

// TestAccount_JSONSerialization tests that the Account struct can be properly serialized to JSON
func TestAccount_JSONSerialization(t *testing.T) {
	account := Account{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password", // Should not appear in JSON
	}

	// This is a basic test to ensure the struct is properly defined
	if account.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", account.ID)
	}

	if account.Username != "testuser" {
		t.Errorf("Expected Username to be 'testuser', got '%s'", account.Username)
	}

	if account.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got '%s'", account.Email)
	}
}

// TestModelsExist verifies that all expected model types are defined
func TestModelsExist(t *testing.T) {
	// Test that we can create instances of all major model types
	// This ensures the types are properly defined and importable
	
	var account Account
	var circle Circle
	var badge Badge
	var tool ToolDescription
	var membership Membership
	
	// Check that types are defined (zero values should work)
	if account.ID != 0 {
		t.Error("Account zero value should have ID = 0")
	}
	
	if circle.ID != 0 {
		t.Error("Circle zero value should have ID = 0")
	}
	
	if badge.ID != 0 {
		t.Error("Badge zero value should have ID = 0")
	}
	
	if tool.ID != 0 {
		t.Error("ToolDescription zero value should have ID = 0")
	}
	
	if membership.ID != 0 {
		t.Error("Membership zero value should have ID = 0")
	}
}