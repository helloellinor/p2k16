package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/helloellinor/p2k16/internal/logging"
)

// Logger middleware for enhanced request logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			// Use our enhanced logger for HTTP requests
			logger := logging.ServerLogger
			if gin.Mode() == gin.DebugMode {
				logger = logging.DemoLogger
			}
			
			// Log the request using our enhanced logger
			logger.LogRequest(
				param.Method,
				param.Path,
				param.ClientIP,
				param.StatusCode,
				param.Latency,
			)
			
			// Return empty string since we handle the output ourselves
			return ""
		},
		Output: nil, // We handle output ourselves
	})
}

// Recovery middleware for panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		}
		c.AbortWithStatus(500)
	})
}

// CORS middleware for handling cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
