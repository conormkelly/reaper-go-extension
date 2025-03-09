package core

import (
	"go-reaper/pkg/logger"
)

// Log level constants for backward compatibility
const (
	LogLevelError   = logger.LevelError
	LogLevelWarning = logger.LevelWarning
	LogLevelInfo    = logger.LevelInfo
	LogLevelDebug   = logger.LevelDebug
	LogLevelTrace   = logger.LevelTrace
)

// InitLogging initializes the logging system
func InitLogging() {
	logger.Initialize()
}

// CleanupLogging shuts down the logging system
func CleanupLogging() {
	logger.Cleanup()
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	logger.Error(format, args...)
}

// LogWarning logs a warning message
func LogWarning(format string, args ...interface{}) {
	logger.Warning(format, args...)
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	logger.Info(format, args...)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}

// LogTrace logs a trace message
func LogTrace(format string, args ...interface{}) {
	logger.Trace(format, args...)
}

// IsLoggingEnabled returns true if logging is enabled
func IsLoggingEnabled() bool {
	return logger.IsLoggingEnabled()
}

// SetLoggingEnabled enables or disables logging
func SetLoggingEnabled(enabled bool) {
	logger.SetLoggingEnabled(enabled)
}

// GetLogLevel returns the current log level
func GetLogLevel() int {
	return logger.GetLogLevel()
}

// SetLogLevel sets the log level
func SetLogLevel(level int) {
	logger.SetLogLevel(level)
}

// SetLogPath sets a custom log file path
func SetLogPath(path string) {
	logger.SetLogPath(path)
}
