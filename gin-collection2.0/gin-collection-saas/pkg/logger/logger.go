package logger

import (
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a simple structured logger
type Logger struct {
	level LogLevel
}

// LogLevel represents log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var globalLogger *Logger

// Init initializes the global logger
func Init(level string) {
	var logLevel LogLevel
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	default:
		logLevel = INFO
	}

	globalLogger = &Logger{level: logLevel}
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	if globalLogger == nil {
		return
	}
	if globalLogger.level <= DEBUG {
		log("DEBUG", msg, fields...)
	}
}

// Info logs an info message
func Info(msg string, fields ...interface{}) {
	if globalLogger == nil {
		return
	}
	if globalLogger.level <= INFO {
		log("INFO", msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	if globalLogger == nil {
		return
	}
	if globalLogger.level <= WARN {
		log("WARN", msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...interface{}) {
	if globalLogger == nil {
		return
	}
	if globalLogger.level <= ERROR {
		log("ERROR", msg, fields...)
	}
}

func log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	output := os.Stdout
	if level == "ERROR" {
		output = os.Stderr
	}

	_, _ = io.WriteString(output, timestamp+" ["+level+"] "+msg)

	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			_, _ = io.WriteString(output, " "+toString(fields[i])+"="+toString(fields[i+1]))
		}
	}

	_, _ = io.WriteString(output, "\n")
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int64:
		return toString(val)
	default:
		return ""
	}
}

// GinLogger returns a Gin middleware for logging
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		timestamp := time.Now()
		latency := timestamp.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"ip", c.ClientIP(),
		)
	}
}
