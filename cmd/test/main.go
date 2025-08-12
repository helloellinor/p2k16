package main

import (
	"fmt"

	"github.com/helloellinor/p2k16/internal/database"
	"github.com/helloellinor/p2k16/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Test bcrypt compatibility
	testPasswords()

	// Test database connection and user loading (if available)
	testDatabase()
}

func testPasswords() {
	fmt.Println("=== Testing Password Compatibility ===")

	// These are the actual hashes from the migration
	superHash := "$2b$12$B/kxR5O85fN357.fZNUPoOiNblCj7j2lX3/VLajLvuE42OmqsyUTO"
	fooHash := "$2b$12$o764MV/jh0HnsAtsEz53L.GfbLwCqZ5jTf3aV2yUAFFCaTrzGCcQm"

	// Test super user
	err := bcrypt.CompareHashAndPassword([]byte(superHash), []byte("super"))
	if err == nil {
		fmt.Println("✓ Super user password verification works")
	} else {
		fmt.Printf("✗ Super user password failed: %v\n", err)
	}

	// Test foo user
	err = bcrypt.CompareHashAndPassword([]byte(fooHash), []byte("foo"))
	if err == nil {
		fmt.Println("✓ Foo user password verification works")
	} else {
		fmt.Printf("✗ Foo user password failed: %v\n", err)
	}

	// Test generating new hash
	newHash, err := models.HashPassword("test123")
	if err != nil {
		fmt.Printf("✗ Failed to generate new hash: %v\n", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(newHash), []byte("test123"))
	if err == nil {
		fmt.Println("✓ New password generation and verification works")
	} else {
		fmt.Printf("✗ New password verification failed: %v\n", err)
	}
}

func testDatabase() {
	fmt.Println("\n=== Testing Database Connection ===")

	dbConfig := database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "p2k16-web",
		Password: "p2k16-web",
		DBName:   "p2k16",
		SSLMode:  "disable",
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		fmt.Printf("✗ Database connection failed: %v\n", err)
		fmt.Println("  (This is expected if the database is not running)")
		return
	}
	defer db.Close()

	fmt.Println("✓ Database connection successful")

	// Test loading users
	accountRepo := models.NewAccountRepository(db.DB)

	superUser, err := accountRepo.FindByUsername("super")
	if err != nil {
		fmt.Printf("✗ Failed to load super user: %v\n", err)
		return
	}

	fmt.Printf("✓ Loaded super user: ID=%d, Username=%s, Email=%s\n",
		superUser.ID, superUser.Username, superUser.Email)

	// Test password validation
	if superUser.ValidatePassword("super") {
		fmt.Println("✓ Super user password validation works with database")
	} else {
		fmt.Println("✗ Super user password validation failed")
	}

	// Test foo user
	fooUser, err := accountRepo.FindByUsername("foo")
	if err != nil {
		fmt.Printf("✗ Failed to load foo user: %v\n", err)
		return
	}

	fmt.Printf("✓ Loaded foo user: ID=%d, Username=%s, Email=%s\n",
		fooUser.ID, fooUser.Username, fooUser.Email)

	if fooUser.ValidatePassword("foo") {
		fmt.Println("✓ Foo user password validation works with database")
	} else {
		fmt.Println("✗ Foo user password validation failed")
	}
}
