package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/helloellinor/p2k16/internal/database"
	"github.com/helloellinor/p2k16/internal/handlers"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
	"github.com/helloellinor/p2k16/internal/session"
)

func main() {
	// Database configuration - matches the existing setup
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 2016),
		User:     getEnv("DB_USER", "p2k16-web"),
		Password: getEnv("DB_PASSWORD", "p2k16-web"),
		DBName:   getEnv("DB_NAME", "p2k16"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database - exit if connection fails
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Printf("💡 Please ensure PostgreSQL is running and credentials are correct")
		log.Printf("🔧 Database config: %s@%s:%d/%s", dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		log.Fatalf("Cannot start server without database connection")
	}

	log.Printf("✅ Database connection successful")
	defer db.Close()

	// Initialize repositories
	accountRepo := models.NewAccountRepository(db.DB)
	circleRepo := models.NewCircleRepository(db.DB)
	badgeRepo := models.NewBadgeRepository(db.DB)
	toolRepo := models.NewToolRepository(db.DB)
	eventRepo := models.NewEventRepository(db.DB)
	membershipRepo := models.NewMembershipRepository(db.DB)

	// Initialize session manager
	sessionManager := session.NewChiSessionManager()

	// Initialize handlers
	handler := handlers.NewChiHandler(accountRepo, circleRepo, badgeRepo, toolRepo, eventRepo, membershipRepo, sessionManager)

	// Set up Chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.ChiLogger())
	r.Use(middleware.ChiRecovery())
	r.Use(middleware.ChiCORS())
	r.Use(sessionManager.LoadAndSave)

	// Serve static files
	r.Handle("/styles/*", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.ChiOptionalAuth(sessionManager, handler.GetAccountRepo()))
		r.Get("/", handler.ChiHome)
		r.Get("/login", handler.ChiLogin)
	})

	// Logout route (no auth required, just clears session)
	r.Post("/logout", handler.ChiLogout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.ChiRequireAuth(sessionManager, handler.GetAccountRepo()))
		
		r.Get("/dashboard", handler.ChiDashboard)
		r.Get("/profile", handler.ChiProfile)
		r.Get("/admin", handler.ChiAdmin)

		// Admin routes
		r.Get("/admin/users", handler.ChiAdminUsers)
		r.Get("/admin/tools", handler.ChiAdminTools)
		r.Get("/admin/companies", handler.ChiAdminCompanies)
		r.Get("/admin/circles", handler.ChiAdminCircles)
		r.Get("/admin/logs", handler.ChiAdminLogs)
		r.Get("/admin/config", handler.ChiAdminConfig)

		// Profile management endpoints
		r.Post("/profile/change-password", handler.ChiChangePassword)
		r.Post("/profile/update", handler.ChiUpdateProfile)
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/members/active", handler.ChiGetActiveMembers)
		r.Post("/auth/login", handler.ChiAuthLogin)

		// Protected API routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.ChiRequireAuth(sessionManager, handler.GetAccountRepo()))
			
			// Account management endpoints
			r.Get("/accounts", handler.ChiGetAccounts)
			r.Get("/accounts/{id}", handler.ChiGetAccount)

			// Badge management endpoints
			r.Get("/badges", handler.ChiGetBadges)
			r.Get("/user/badges", handler.ChiGetUserBadges)
			r.Get("/badges/available", handler.ChiGetAvailableBadges)
			r.Post("/badges/create", handler.ChiCreateBadge)
			r.Post("/badges/award", handler.ChiAwardBadge)

			// Membership endpoints
			r.Get("/memberships", handler.ChiGetMembershipStatusAPI)
			r.Get("/membership/status", handler.ChiGetMembershipStatus)
			r.Get("/membership/active", handler.ChiGetActiveMembersDetailed)

			// Tool management routes
			r.Get("/tools", handler.ChiGetTools)
			r.Get("/tools/checkouts", handler.ChiGetActiveCheckouts)
			r.Post("/tools/checkout", handler.ChiCheckoutTool)
			r.Post("/tools/checkin", handler.ChiCheckinTool)

			// Profile card flip endpoints for HTMX
			r.Get("/profile/card/front", handler.ChiProfileCardFront)
			r.Get("/profile/card/back", handler.ChiProfileCardBack)
		})
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Printf("Application will be available at http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}