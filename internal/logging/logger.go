package logging

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	// Color definitions for different log types
	timestampColor = color.New(color.FgHiBlack).SprintFunc()
	demoColor      = color.New(color.FgMagenta, color.Bold).SprintFunc()
	handlerColor   = color.New(color.FgCyan, color.Bold).SprintFunc()
	serverColor    = color.New(color.FgGreen, color.Bold).SprintFunc()
	errorColor     = color.New(color.FgRed, color.Bold).SprintFunc()
	successColor   = color.New(color.FgGreen).SprintFunc()
	warningColor   = color.New(color.FgYellow).SprintFunc()
	infoColor      = color.New(color.FgBlue).SprintFunc()
)

// LogLevel represents different log levels
type LogLevel string

const (
	LogLevelDemo    LogLevel = "DEMO"
	LogLevelHandler LogLevel = "HANDLER"
	LogLevelServer  LogLevel = "SERVER"
	LogLevelError   LogLevel = "ERROR"
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelSuccess LogLevel = "SUCCESS"
)

// Logger provides enhanced logging functionality with colors and emojis
type Logger struct {
	isDemoMode bool
}

// NewLogger creates a new logger instance
func NewLogger(isDemoMode bool) *Logger {
	return &Logger{
		isDemoMode: isDemoMode,
	}
}

// LogAction logs an action with enhanced formatting
func (l *Logger) LogAction(level LogLevel, action, details string) {
	timestamp := timestampColor(time.Now().Format("15:04:05"))
	
	var emoji, levelStr string
	var colorFunc func(a ...interface{}) string
	
	switch level {
	case LogLevelDemo:
		emoji = "🎯"
		levelStr = demoColor("[DEMO MODE]")
		colorFunc = demoColor
	case LogLevelHandler:
		emoji = "🔧"
		levelStr = handlerColor("[HANDLER]")
		colorFunc = handlerColor
	case LogLevelServer:
		emoji = "🚀"
		levelStr = serverColor("[SERVER]")
		colorFunc = serverColor
	case LogLevelError:
		emoji = "❌"
		levelStr = errorColor("[ERROR]")
		colorFunc = errorColor
	case LogLevelWarning:
		emoji = "⚠️ "
		levelStr = warningColor("[WARNING]")
		colorFunc = warningColor
	case LogLevelSuccess:
		emoji = "✅"
		levelStr = successColor("[SUCCESS]")
		colorFunc = successColor
	case LogLevelInfo:
		emoji = "ℹ️ "
		levelStr = infoColor("[INFO]")
		colorFunc = infoColor
	default:
		emoji = "📝"
		levelStr = "[LOG]"
		colorFunc = func(a ...interface{}) string { return fmt.Sprint(a...) }
	}
	
	fmt.Printf("\n%s %s %s | %s | %s\n", 
		emoji, 
		levelStr, 
		timestamp, 
		colorFunc(action), 
		details)
}

// LogRequest logs HTTP request information
func (l *Logger) LogRequest(method, path, clientIP string, statusCode int, latency time.Duration) {
	var statusEmoji string
	var statusColor func(a ...interface{}) string
	
	switch {
	case statusCode >= 200 && statusCode < 300:
		statusEmoji = "✅"
		statusColor = successColor
	case statusCode >= 300 && statusCode < 400:
		statusEmoji = "➡️ "
		statusColor = infoColor
	case statusCode >= 400 && statusCode < 500:
		statusEmoji = "⚠️ "
		statusColor = warningColor
	default:
		statusEmoji = "❌"
		statusColor = errorColor
	}
	
	timestamp := timestampColor(time.Now().Format("15:04:05"))
	levelStr := l.getLevelString()
	
	fmt.Printf("\n%s %s %s | %s %s %s | %s | %s\n",
		statusEmoji,
		levelStr,
		timestamp,
		statusColor(fmt.Sprintf("%d", statusCode)),
		method,
		path,
		clientIP,
		latency.String())
}

// LogStartup logs server startup information with banner
func (l *Logger) LogStartup(mode, port string, features []string) {
	banner := "============================================================"
	
	if l.isDemoMode {
		fmt.Printf("\n%s\n", banner)
		fmt.Printf("🎭  %s\n", demoColor("P2K16 DEMO MODE - Development Testing Server"))
		fmt.Printf("%s\n", banner)
		fmt.Printf("📍 Server URL: %s\n", infoColor(fmt.Sprintf("http://localhost:%s", port)))
		fmt.Printf("🔑 Demo Login: %s\n", warningColor("username='demo', password=any"))
	} else {
		fmt.Printf("\n%s\n", banner)
		fmt.Printf("🚀 %s\n", serverColor("P2K16 SERVER - Production Mode"))
		fmt.Printf("%s\n", banner)
		fmt.Printf("📍 Server URL: %s\n", infoColor(fmt.Sprintf("http://localhost:%s", port)))
		fmt.Printf("💾 Database: %s\n", successColor("Connected and ready"))
	}
	
	fmt.Printf("📋 Features Available:\n")
	for _, feature := range features {
		fmt.Printf("   • %s\n", feature)
	}
	
	if l.isDemoMode {
		fmt.Printf("⚠️  Note: %s\n", warningColor("No database - all changes are simulated"))
	} else {
		fmt.Printf("💾 Note: %s\n", successColor("All data operations will be persisted"))
	}
	
	fmt.Printf("%s\n", banner)
	fmt.Printf("🚀 Starting server...\n")
}

// LogDatabaseFallback logs when falling back to demo mode due to database issues
func (l *Logger) LogDatabaseFallback(err error) {
	banner := "============================================================"
	fmt.Printf("\n%s\n", banner)
	fmt.Printf("⚠️  %s\n", warningColor("P2K16 SERVER - FALLBACK TO DEMO MODE"))
	fmt.Printf("%s\n", banner)
	fmt.Printf("❌ Database connection failed: %s\n", errorColor(err.Error()))
	fmt.Printf("🎭 Falling back to %s\n", demoColor("DEMO MODE - no database required"))
	fmt.Printf("🔑 Demo logins available:\n")
	fmt.Printf("   • %s\n", infoColor("demo/password"))
	fmt.Printf("   • %s\n", infoColor("super/super"))
	fmt.Printf("   • %s\n", infoColor("foo/foo"))
	fmt.Printf("⚠️  Note: %s\n", warningColor("All data operations will be simulated"))
	fmt.Printf("%s\n", banner)
}

// getLevelString returns the appropriate level string based on mode
func (l *Logger) getLevelString() string {
	if l.isDemoMode {
		return demoColor("[DEMO MODE]")
	}
	return serverColor("[SERVER]")
}

// Global logger instances
var (
	DemoLogger   = NewLogger(true)
	ServerLogger = NewLogger(false)
)

// Convenience functions for direct logging
func LogDemoAction(action, details string) {
	DemoLogger.LogAction(LogLevelDemo, action, details)
}

func LogHandlerAction(action, details string) {
	ServerLogger.LogAction(LogLevelHandler, action, details)
}

func LogServerAction(action, details string) {
	ServerLogger.LogAction(LogLevelServer, action, details)
}

func LogError(action, details string) {
	ServerLogger.LogAction(LogLevelError, action, details)
}

func LogSuccess(action, details string) {
	ServerLogger.LogAction(LogLevelSuccess, action, details)
}

func LogWarning(action, details string) {
	ServerLogger.LogAction(LogLevelWarning, action, details)
}

func LogInfo(action, details string) {
	ServerLogger.LogAction(LogLevelInfo, action, details)
}