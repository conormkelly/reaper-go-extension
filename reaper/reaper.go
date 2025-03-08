package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var (
	// Store REAPER function pointers
	showConsoleMsgPtr unsafe.Pointer
	registerFuncPtr   unsafe.Pointer
	mutex             sync.RWMutex
	initialized       bool
)

// Initialize stores the necessary function pointers from REAPER
func Initialize(info unsafe.Pointer) error {
	mutex.Lock()
	defer mutex.Unlock()

	if info == nil {
		return fmt.Errorf("null plugin info pointer")
	}

	pluginInfo := (*C.reaper_plugin_info_t)(info)

	// Check API version
	if pluginInfo.caller_version != 0x20E {
		return fmt.Errorf("wrong REAPER plugin API version. Expected %d, got %d",
			0x20E, pluginInfo.caller_version)
	}

	// Get the GetFunc function from REAPER
	getFuncPtr := pluginInfo.GetFunc
	if getFuncPtr == nil {
		return fmt.Errorf("could not get GetFunc function pointer")
	}

	// Get the ShowConsoleMsg function
	cFuncName := C.CString("ShowConsoleMsg")
	defer C.free(unsafe.Pointer(cFuncName))
	showConsoleMsgPtr = C.plugin_bridge_call_get_func(unsafe.Pointer(getFuncPtr), cFuncName)

	if showConsoleMsgPtr == nil {
		return fmt.Errorf("could not get ShowConsoleMsg function pointer")
	}

	// Store the Register function
	registerFuncPtr = unsafe.Pointer(pluginInfo.Register)
	if registerFuncPtr == nil {
		return fmt.Errorf("could not get Register function pointer")
	}

	initialized = true
	return nil
}

// ConsoleLog sends a message to the REAPER console
func ConsoleLog(message string) error {
	mutex.RLock()
	defer mutex.RUnlock()

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	if showConsoleMsgPtr == nil {
		return fmt.Errorf("ShowConsoleMsg function not available")
	}

	cMessage := C.CString(message + "\n")
	defer C.free(unsafe.Pointer(cMessage))
	C.plugin_bridge_call_show_console_msg(showConsoleMsgPtr, cMessage)
	return nil
}

// RegisterAction registers a custom action with REAPER
func RegisterAction(commandID string, actionName string) (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if !initialized {
		return -1, fmt.Errorf("REAPER functions not initialized")
	}

	if registerFuncPtr == nil {
		return -1, fmt.Errorf("Register function not available")
	}

	cCmdID := C.CString(commandID)
	defer C.free(unsafe.Pointer(cCmdID))

	cActionName := C.CString(actionName)
	defer C.free(unsafe.Pointer(cActionName))

	result := C.plugin_bridge_call_register(registerFuncPtr, cCmdID, unsafe.Pointer(cActionName))
	return int(result), nil
}
