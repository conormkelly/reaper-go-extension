package logger

/*
#cgo CFLAGS: -I${SRCDIR}/../../c -I${SRCDIR}/../../sdk
#include "../../c/logging.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// cLogMessage sends a log message to the C logging system
func cLogMessage(level int, funcName, message string) {
	// Convert Go strings to C strings
	cFuncName := C.CString(funcName)
	cMessage := C.CString(message)

	// Ensure we free the C strings to avoid memory leaks
	defer C.free(unsafe.Pointer(cFuncName))
	defer C.free(unsafe.Pointer(cMessage))

	// Call the C logging function
	C.log_message(C.LogLevel(level), cFuncName, cMessage)
}

// initLogging initializes the logging system
func initLogging() {
	C.log_init()
}

// cleanupLogging shuts down the logging system
func cleanupLogging() {
	C.log_cleanup()
}

// isEnabled returns true if logging is enabled
func isEnabled() bool {
	return bool(C.log_is_enabled())
}

// setEnabled enables or disables logging
func setEnabled(enabled bool) {
	C.log_set_enabled(C.bool(enabled))
}

// getLevel returns the current log level
func getLevel() int {
	return int(C.log_get_level())
}

// setLevel sets the log level
func setLevel(level int) {
	if level >= LevelError && level <= LevelTrace {
		C.log_set_level(C.LogLevel(level))
	}
}

// setPath sets a custom log file path
func setPath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.log_set_path(cPath)
}
