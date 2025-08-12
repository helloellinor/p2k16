package main

import (
	"log"
	"os"

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
		Port:     5432,
		User:     getEnv("DB_USER", "p2k16-web"),
		Password: getEnv("DB_PASSWORD", "p2k16-web"),
		DBName:   getEnv("DB_NAME", "p2k16"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	accountRepo := models.NewAccountRepository(db.DB)
	circleRepo := models.NewCircleRepository(db.DB)

	// Initialize handlers
	handler := handlers.NewHandler(accountRepo, circleRepo)

	// Set up Gin router
	r := gin.New()

	// Add middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

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
	r.GET("/", middleware.OptionalAuth(accountRepo), handler.Home)
	r.GET("/login", middleware.OptionalAuth(accountRepo), handler.Login)
	r.GET("/logout", handler.Logout)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(accountRepo))
	{
		protected.GET("/dashboard", handler.Dashboard)
	}

	// API routes
	api := r.Group("/api")
	{
		api.GET("/members/active", handler.GetActiveMembers)
		api.POST("/auth/login", handler.AuthLogin)

		// Protected API routes
		apiProtected := api.Group("/")
		apiProtected.Use(middleware.RequireAuth(accountRepo))
		{
			apiProtected.GET("/user/badges", handler.GetUserBadges)
		}
	}

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
