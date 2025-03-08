package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
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

	// Store GetFunc for later use using our bridge function
	C.plugin_bridge_set_get_func(unsafe.Pointer(getFuncPtr))

	// Get the ShowConsoleMsg function
	cFuncName := C.CString("ShowConsoleMsg")
	defer C.free(unsafe.Pointer(cFuncName))
	showConsoleMsgPtr = C.plugin_bridge_call_get_func(unsafe.Pointer(getFuncPtr), cFuncName)

	if showConsoleMsgPtr == nil {
		return fmt.Errorf("could not get ShowConsoleMsg function pointer")
	}

	// Store the Register function pointer
	registerFuncPtr = unsafe.Pointer(pluginInfo.Register)
	if registerFuncPtr == nil {
		return fmt.Errorf("could not get Register function pointer")
	}

	// Clear registered commands map
	registeredCommands = make(map[string]int)

	// Initialize action handlers map
	initActionHandlers()

	// Register command hooks
	cHookCmd2 := C.CString("hookcommand2")
	defer C.free(unsafe.Pointer(cHookCmd2))
	C.plugin_bridge_call_register(registerFuncPtr, cHookCmd2, unsafe.Pointer(C.goHookCommandProc2))

	cHookCmd := C.CString("hookcommand")
	defer C.free(unsafe.Pointer(cHookCmd))
	C.plugin_bridge_call_register(registerFuncPtr, cHookCmd, unsafe.Pointer(C.goHookCommandProc))

	initialized = true
	return nil
}
