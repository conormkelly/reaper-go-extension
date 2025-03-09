package core

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_ext_logging.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Log level constants
const (
	LogLevelError   = 0
	LogLevelWarning = 1
	LogLevelInfo    = 2
	LogLevelDebug   = 3
	LogLevelTrace   = 4
)

// InitLogging initializes the logging system
func InitLogging() {
	C.log_init()
}

// CleanupLogging shuts down the logging system
func CleanupLogging() {
	C.log_cleanup()
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	logMessage(LogLevelError, format, args...)
}

// LogWarning logs a warning message
func LogWarning(format string, args ...interface{}) {
	logMessage(LogLevelWarning, format, args...)
}

// LogInfo logs an info message
func LogInfo(format string, args ...interface{}) {
	logMessage(LogLevelInfo, format, args...)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	logMessage(LogLevelDebug, format, args...)
}

// LogTrace logs a trace message
func LogTrace(format string, args ...interface{}) {
	logMessage(LogLevelTrace, format, args...)
}

// IsLoggingEnabled returns true if logging is enabled
func IsLoggingEnabled() bool {
	return bool(C.log_is_enabled())
}

// SetLoggingEnabled enables or disables logging
func SetLoggingEnabled(enabled bool) {
	C.log_set_enabled(C.bool(enabled))
}

// GetLogLevel returns the current log level
func GetLogLevel() int {
	return int(C.log_get_level())
}

// SetLogLevel sets the log level
func SetLogLevel(level int) {
	if level >= LogLevelError && level <= LogLevelTrace {
		C.log_set_level(C.LogLevel(level))
	}
}

// SetLogPath sets a custom log file path
func SetLogPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.log_set_path(cPath)
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

	// Convert to C strings for the C logging function
	cFuncName := C.CString(funcName)
	cMessage := C.CString(message)

	// Ensure we free the C strings to avoid memory leaks
	defer C.free(unsafe.Pointer(cFuncName))
	defer C.free(unsafe.Pointer(cMessage))

	// Call the C logging function
	C.log_message(C.LogLevel(level), cFuncName, cMessage)
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
