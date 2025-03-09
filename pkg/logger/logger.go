// Package logger provides centralized logging functionality for the REAPER Go extension.
package logger

import (
	"fmt"
	"runtime"
)

// Log level constants
const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

// Error logs an error message
func Error(format string, args ...interface{}) {
	logMessage(LevelError, format, args...)
}

// Warning logs a warning message
func Warning(format string, args ...interface{}) {
	logMessage(LevelWarning, format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	logMessage(LevelInfo, format, args...)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	logMessage(LevelDebug, format, args...)
}

// Trace logs a trace message
func Trace(format string, args ...interface{}) {
	logMessage(LevelTrace, format, args...)
}

// IsLoggingEnabled returns true if logging is enabled
func IsLoggingEnabled() bool {
	return isEnabled()
}

// SetLoggingEnabled enables or disables logging
func SetLoggingEnabled(enabled bool) {
	setEnabled(enabled)
}

// GetLogLevel returns the current log level
func GetLogLevel() int {
	return getLevel()
}

// SetLogLevel sets the log level
func SetLogLevel(level int) {
	if level >= LevelError && level <= LevelTrace {
		setLevel(level)
	}
}

// SetLogPath sets a custom log file path
func SetLogPath(path string) {
	setPath(path)
}

// Initialize initializes the logging system
func Initialize() {
	initLogging()
}

// Cleanup shuts down the logging system
func Cleanup() {
	cleanupLogging()
}

// logMessage is the internal function for all logging levels
func logMessage(level int, format string, args ...interface{}) {
	// Skip logging if disabled or level is too verbose
	if !IsLoggingEnabled() || GetLogLevel() < level {
		return
	}

	// Get caller function name for better logging
	pc, _, _, ok := runtime.Caller(2) // Skip logMessage and the specific level function
	funcName := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			funcName = fn.Name()

			// Get just the function name without package path
			if lastDot := last(funcName, '.'); lastDot >= 0 {
				funcName = funcName[lastDot+1:]
			}
		}
	}

	// Format the message with arguments if provided
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format
	}

	// Send to the C logging system
	cLogMessage(level, funcName, message)
}

// last finds the last occurrence of a character in a string
func last(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
