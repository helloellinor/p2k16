package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/database"
	"github.com/helloellinor/p2k16/internal/handlers"
	"github.com/helloellinor/p2k16/internal/middleware"
	"github.com/helloellinor/p2k16/internal/models"
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

	// Try to connect to database - fall back to demo mode if it fails
	db, err := database.NewConnection(dbConfig)
	var handler *handlers.Handler
	var demoMode bool
	
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		log.Printf("Starting in DEMO MODE - no database required")
		log.Printf("Use 'demo' with any password, 'super/super', or 'foo/foo' to login")
		
		// Initialize demo handlers with nil repositories
		handler = handlers.NewHandler(nil, nil, nil, nil, nil, nil)
		demoMode = true
	} else {
		log.Printf("Database connection successful")
		defer db.Close()

		// Initialize repositories
		accountRepo := models.NewAccountRepository(db.DB)
		circleRepo := models.NewCircleRepository(db.DB)
		badgeRepo := models.NewBadgeRepository(db.DB)
		toolRepo := models.NewToolRepository(db.DB)
		eventRepo := models.NewEventRepository(db.DB)
		membershipRepo := models.NewMembershipRepository(db.DB)

		// Initialize handlers
		handler = handlers.NewHandler(accountRepo, circleRepo, badgeRepo, toolRepo, eventRepo, membershipRepo)
		demoMode = false
	}

	// Set up Gin router
	r := gin.New()

	// Add middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// Serve static files
	r.Static("/styles", "./styles")

	// Session middleware
	sessionSecret := getEnv("SESSION_SECRET", "p2k16-secret-key-change-in-production")
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge:   86400 * 7, // 7 days
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	r.Use(sessions.Sessions(middleware.SessionName, store))

	// Public routes
	r.GET("/", middleware.OptionalAuth(handler.GetAccountRepo()), handler.Home)
	r.GET("/login", middleware.OptionalAuth(handler.GetAccountRepo()), handler.Login)
	r.GET("/logout", handler.Logout)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(handler.GetAccountRepo()))
	{
		protected.GET("/dashboard", handler.Dashboard)
		protected.GET("/profile", handler.Profile)
		protected.GET("/admin", handler.Admin)
	}

	// API routes
	api := r.Group("/api")
	{
		api.GET("/members/active", handler.GetActiveMembers)
		api.POST("/auth/login", handler.AuthLogin)

		// Protected API routes
		apiProtected := api.Group("/")
		apiProtected.Use(middleware.RequireAuth(handler.GetAccountRepo()))
		{
			apiProtected.GET("/user/badges", handler.GetUserBadges)
			apiProtected.GET("/badges/available", handler.GetAvailableBadges)
			apiProtected.POST("/badges/create", handler.CreateBadge)
			apiProtected.POST("/badges/award", handler.AwardBadge)
			
			// Tool management routes
			apiProtected.GET("/tools", handler.GetTools)
			apiProtected.GET("/tools/checkouts", handler.GetActiveCheckouts)
			apiProtected.POST("/tools/checkout", handler.CheckoutTool)
			apiProtected.POST("/tools/checkin", handler.CheckinTool)
			
			// Membership routes
			apiProtected.GET("/membership/status", handler.GetMembershipStatus)
			apiProtected.GET("/membership/active", handler.GetActiveMembersDetailed)
		}
	}

	// Set demo mode in handler
	handler.SetDemoMode(demoMode)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Printf("Application will be available at http://localhost:%s", port)

	if err := r.Run(":" + port); err != nil {
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
