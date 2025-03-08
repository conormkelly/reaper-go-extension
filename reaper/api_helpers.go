package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ListAvailableFunctions checks if specific REAPER functions exist and logs the results
func ListAvailableFunctions(functionNames []string) {
	if !initialized {
		ConsoleLog("REAPER functions not initialized")
		return
	}

	ConsoleLog("Checking for available functions:")

	getFuncPtr := C.plugin_bridge_get_get_func()
	if getFuncPtr == nil {
		ConsoleLog("Error: GetFunc function pointer is nil")
		return
	}

	for _, name := range functionNames {
		cFuncName := C.CString(name)
		funcPtr := C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
		C.free(unsafe.Pointer(cFuncName))

		if funcPtr != nil {
			ConsoleLog(fmt.Sprintf("- %s: Available", name))
		} else {
			ConsoleLog(fmt.Sprintf("- %s: Not found", name))
		}
	}
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
