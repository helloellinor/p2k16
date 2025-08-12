package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("⚠️  P2K16 SERVER - FALLBACK TO DEMO MODE")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("❌ Database connection failed: %v\n", err)
		fmt.Println("🎭 Falling back to DEMO MODE - no database required")
		fmt.Println("🔑 Demo logins available:")
		fmt.Println("   • demo/password")
		fmt.Println("   • super/super") 
		fmt.Println("   • foo/foo")
		fmt.Println("⚠️  Note: All data operations will be simulated")
		fmt.Println(strings.Repeat("=", 60))
		
		// Initialize demo handlers with nil repositories
		handler = handlers.NewHandler(nil, nil, nil, nil, nil, nil)
		demoMode = true
	} else {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("🚀 P2K16 SERVER - PRODUCTION MODE")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("✅ Database connection successful")
		fmt.Printf("🗄️  Connected to: %s@%s:%d/%s\n", dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		fmt.Println("💾 All data operations will be persisted to database")
		fmt.Println(strings.Repeat("=", 60))
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
	
	// Session validation middleware
	r.Use(middleware.SessionValidationMiddleware())

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
			
			// Profile management routes (Phase 2)
			apiProtected.POST("/profile/change-password", handler.ChangePassword)
			apiProtected.POST("/profile/update", handler.UpdateProfile)
		}
	}

	// Set demo mode in handler
	handler.SetDemoMode(demoMode)

	// Start server
	port := getEnv("PORT", "8080")
	
	if demoMode {
		fmt.Printf("🌐 Demo server starting on http://localhost:%s\n", port)
		fmt.Println("📋 Available features: Dashboard, Profile Management, Badge System")
	} else {
		fmt.Printf("🌐 Production server starting on http://localhost:%s\n", port)
		fmt.Println("📋 Full feature set available with database persistence")
	}
	fmt.Println("🚀 Server starting...")

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
