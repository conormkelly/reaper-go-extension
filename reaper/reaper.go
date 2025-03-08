package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>

// Define our own simplified structures for REAPER API
typedef struct {
  int uniqueSectionId;  // Section ID (0=main, 32060=midi editor, etc)
  const char* idStr;    // Unique ID string for the action
  const char* name;     // Display name for the action list
  void *extra;          // Reserved for future use (NULL)
} our_custom_action_t;

// Forward declaration of the Go callback function
extern int goHookCommandProc(int commandId, int flag);
extern int goHookCommandProc2(void* section, int commandId, int val, int valhw, int relmode, void* hwnd, void* proj);
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

	// Track registered command IDs
	registeredCommands map[string]int
)

func init() {
	registeredCommands = make(map[string]int)
}

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

	// Store the Register function pointer
	registerFuncPtr = unsafe.Pointer(pluginInfo.Register)
	if registerFuncPtr == nil {
		return fmt.Errorf("could not get Register function pointer")
	}

	// Clear registered commands map
	registeredCommands = make(map[string]int)

	// Register command hooks - THIS IS THE KEY PART THAT WAS MISSING BEFORE
	cHookCmd2 := C.CString("hookcommand2")
	defer C.free(unsafe.Pointer(cHookCmd2))
	C.plugin_bridge_call_register(registerFuncPtr, cHookCmd2, unsafe.Pointer(C.goHookCommandProc2))

	cHookCmd := C.CString("hookcommand")
	defer C.free(unsafe.Pointer(cHookCmd))
	C.plugin_bridge_call_register(registerFuncPtr, cHookCmd, unsafe.Pointer(C.goHookCommandProc))

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

// Command callback handlers
//
//export goHookCommandProc
func goHookCommandProc(commandId C.int, flag C.int) C.int {
	// Check if this is one of our registered commands
	for _, cmdID := range registeredCommands {
		if int(commandId) == cmdID {
			// Command was triggered - log it to console
			ConsoleLog(fmt.Sprintf("GoReaper action triggered! Command ID: %d", int(commandId)))
			return 1 // Return 1 to indicate we handled it
		}
	}
	return 0 // Not our command, let REAPER handle it
}

//export goHookCommandProc2
func goHookCommandProc2(section unsafe.Pointer, commandId C.int, val C.int, valhw C.int, relmode C.int, hwnd unsafe.Pointer, proj unsafe.Pointer) C.int {
	// Similar to hookCommandProc, check if this is one of our commands
	for _, cmdID := range registeredCommands {
		if int(commandId) == cmdID {
			// Command was triggered - log it to console
			ConsoleLog(fmt.Sprintf("GoReaper action triggered! Command ID: %d (via hookcommand2)", int(commandId)))
			return 1 // Return 1 to indicate we handled it
		}
	}
	return 0 // Not our command, let REAPER handle it
}

// RegisterCustomAction registers an action using the correct SWS-like approach
func RegisterCustomAction(actionID string, description string, sectionID int) (int, error) {
	if !initialized {
		return -1, fmt.Errorf("REAPER functions not initialized")
	}

	// 1. Register the command ID first
	mutex.Lock()
	cCommandID := C.CString("command_id")
	defer C.free(unsafe.Pointer(cCommandID))

	cActionID := C.CString(actionID)
	defer C.free(unsafe.Pointer(cActionID))

	cmdIDResult := C.plugin_bridge_call_register(registerFuncPtr, cCommandID, unsafe.Pointer(cActionID))
	cmdID := int(cmdIDResult)

	if cmdID <= 0 {
		mutex.Unlock()
		return -1, fmt.Errorf("failed to register command ID")
	}

	// Store command ID for lookup in hook handlers
	registeredCommands[actionID] = cmdID

	// 2. Now register the custom action with more details
	cDesc := C.CString(description)
	defer C.free(unsafe.Pointer(cDesc))

	// Create custom_action_register_t struct
	customAction := C.our_custom_action_t{
		uniqueSectionId: C.int(sectionID),
		idStr:           cActionID,
		name:            cDesc,
		extra:           nil,
	}

	// Register the custom action
	cCustomAction := C.CString("custom_action")
	defer C.free(unsafe.Pointer(cCustomAction))

	caResult := C.plugin_bridge_call_register(registerFuncPtr, cCustomAction, unsafe.Pointer(&customAction))

	mutex.Unlock()

	ConsoleLog(fmt.Sprintf("Registered custom action: %s (%s) in section %d",
		actionID, description, sectionID))
	ConsoleLog(fmt.Sprintf("  Command ID result: %d, Custom action result: %d",
		cmdID, int(caResult)))

	return cmdID, nil
}

// Section ID constants
const (
	SectionMain          = 0
	SectionMainAlt       = 100
	SectionMIDIEditor    = 32060
	SectionMIDIEventList = 32061
	SectionMIDIInline    = 32062
	SectionMediaExplorer = 32063
)

// RegisterMainAction registers an action in the main section
func RegisterMainAction(actionID string, description string) (int, error) {
	return RegisterCustomAction(actionID, description, SectionMain)
}
