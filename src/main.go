package main

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/../sdk
#include "reaper_plugin_bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Global variables to store REAPER function pointers
var (
	showConsoleMsgPtr unsafe.Pointer
)

//export GoReaperPluginEntry
func GoReaperPluginEntry(hInstance unsafe.Pointer, rec unsafe.Pointer) C.int {
	fmt.Println("Go REAPER plugin entry called")

	// If rec is null, REAPER is unloading the plugin
	if rec == nil {
		fmt.Println("Go plugin unloading")
		return 0
	}

	info := (*C.reaper_plugin_info_t)(rec)

	// Check API version - REAPER_PLUGIN_VERSION is 0x20E as per your definition
	if info.caller_version != 0x20E {
		fmt.Printf("Wrong REAPER plugin API version. Expected %d, got %d\n",
			0x20E, info.caller_version)
		return 0
	}

	// Get the GetFunc function from REAPER
	getFuncPtr := info.GetFunc
	if getFuncPtr == nil {
		fmt.Println("Could not get GetFunc function pointer")
		return 0
	}

	// Get the ShowConsoleMsg function using GetFunc
	cFuncName := C.CString("ShowConsoleMsg")
	defer C.free(unsafe.Pointer(cFuncName))
	showConsoleMsgPtr = C.plugin_bridge_call_get_func(unsafe.Pointer(getFuncPtr), cFuncName)

	if showConsoleMsgPtr == nil {
		fmt.Println("Could not get ShowConsoleMsg function pointer")
		return 0
	}

	// Show message in REAPER console
	message := C.CString("Hello from Go REAPER extension!\n")
	defer C.free(unsafe.Pointer(message))
	C.plugin_bridge_call_show_console_msg(showConsoleMsgPtr, message)

	// Register a basic command
	cmdName := C.CString("command_id")
	defer C.free(unsafe.Pointer(cmdName))

	actionName := C.CString("GO_HELLO_WORLD")
	defer C.free(unsafe.Pointer(actionName))

	// Get the Register function from REAPER
	registerFuncPtr := info.Register
	if registerFuncPtr != nil {
		result := C.plugin_bridge_call_register(unsafe.Pointer(registerFuncPtr), cmdName, unsafe.Pointer(actionName))
		fmt.Printf("Registered command: GO_HELLO_WORLD, result: %d\n", int(result))
	}

	fmt.Println("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
