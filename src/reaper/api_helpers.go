package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../../sdk
#include "../c/bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"go-reaper/src/pkg/logger"
	"unsafe"
)

// ListAvailableFunctions checks if specific REAPER functions exist and logs the results
func ListAvailableFunctions(functionNames []string) {
	if !initialized {
		logger.Error("REAPER functions not initialized")
		return
	}

	logger.Debug("Checking for available functions:")

	getFuncPtr := C.plugin_bridge_get_get_func()
	if getFuncPtr == nil {
		logger.Error("Error: GetFunc function pointer is nil")
		return
	}

	for _, name := range functionNames {
		cFuncName := C.CString(name)
		funcPtr := C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
		C.free(unsafe.Pointer(cFuncName))

		if funcPtr != nil {
			logger.Debug("- %s: Available", name)
		} else {
			logger.Warning("- %s: Not found", name)
		}
	}
}

// IsFunctionAvailable checks if a specific REAPER function exists
func IsFunctionAvailable(functionName string) bool {
	if !initialized {
		return false
	}

	getFuncPtr := C.plugin_bridge_get_get_func()
	if getFuncPtr == nil {
		return false
	}

	cFuncName := C.CString(functionName)
	defer C.free(unsafe.Pointer(cFuncName))

	funcPtr := C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
	return funcPtr != nil
}

// GetFunctionPointer returns a pointer to a REAPER function if available
func GetFunctionPointer(functionName string) unsafe.Pointer {
	if !initialized {
		return nil
	}

	getFuncPtr := C.plugin_bridge_get_get_func()
	if getFuncPtr == nil {
		return nil
	}

	cFuncName := C.CString(functionName)
	defer C.free(unsafe.Pointer(cFuncName))

	return C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
}

// ReaperConsoleLog sends a message directly to the REAPER console without our package's initialization check
// This is useful for debugging when the main initialization may have failed
func ReaperConsoleLog(message string) {
	cFuncName := C.CString("ShowConsoleMsg")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_get_get_func()
	if getFuncPtr == nil {
		// Can't log as we have no way to access REAPER functions
		return
	}

	showConsoleMsgPtr := C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
	if showConsoleMsgPtr == nil {
		// Can't log as ShowConsoleMsg is not available
		return
	}

	cMessage := C.CString(message + "\n")
	defer C.free(unsafe.Pointer(cMessage))
	C.plugin_bridge_call_show_console_msg(showConsoleMsgPtr, cMessage)
}
