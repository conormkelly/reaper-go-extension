package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../sdk
#include "../c/bridge.h"
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
	"unsafe"
)

var (
	// Track registered command IDs
	registeredCommands map[string]int
	// Store a map of action handlers
	actionHandlers map[string]ActionHandler
)

func init() {
	registeredCommands = make(map[string]int)
}

// Initialize action handlers map
func initActionHandlers() {
	actionHandlers = make(map[string]ActionHandler)
}

// SetActionHandler associates a function with an action ID
func SetActionHandler(actionID string, handler ActionHandler) {
	mutex.Lock()
	defer mutex.Unlock()

	actionHandlers[actionID] = handler
}

// Command callback handlers
//
//export goHookCommandProc
func goHookCommandProc(commandId C.int, flag C.int) C.int {
	// Check if this is one of our registered commands
	for actionID, cmdID := range registeredCommands {
		if int(commandId) == cmdID {
			// Log that the command was triggered
			// core.LogInfo("GoReaper action triggered! Command ID: %d (%s)", int(commandId), actionID)

			// Check if we have a handler for this action
			mutex.RLock()
			handler, exists := actionHandlers[actionID]
			mutex.RUnlock()

			if exists {
				// Execute the handler
				handler()
			}

			return 1 // Return 1 to indicate we handled it
		}
	}
	return 0 // Not our command, let REAPER handle it
}

//export goHookCommandProc2
func goHookCommandProc2(section unsafe.Pointer, commandId C.int, val C.int, valhw C.int, relmode C.int, hwnd unsafe.Pointer, proj unsafe.Pointer) C.int {
	// Called by REAPER when an action is triggered. Return 1 if handled, 0 to pass to other plugins.
	// commandId: unique identifier for the action
	// val, valhw: action parameters that may contain state information
	// relmode: relative mouse mode (0=absolute, 1/2=relative from last value)
	// hwnd: window handle
	// proj: project context"

	// Similar to hookCommandProc, check if this is one of our commands
	for actionID, cmdID := range registeredCommands {
		if int(commandId) == cmdID {
			// Command was triggered - log it to console
			// core.LogInfo("GoReaper action triggered! Command ID: %d (%s) (via hookcommand2)", int(commandId), actionID)

			// Check if we have a handler for this action
			mutex.RLock()
			handler, exists := actionHandlers[actionID]
			mutex.RUnlock()

			if exists {
				// Execute the handler
				handler()
			}

			return 1 // Return 1 to indicate we handled it
		}
	}
	return 0 // Not our command, let REAPER handle it
}

// RegisterCustomAction uses a two-step registration process: first register a command ID, then register the custom
// action details. Both must succeed for the action to appear in REAPER's action list.
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

	// caResult var, unused because of commented logs below (temporary)
	_ = C.plugin_bridge_call_register(registerFuncPtr, cCustomAction, unsafe.Pointer(&customAction))

	mutex.Unlock()

	// core.LogInfo("Registered custom action: %s (%s) in section %d", actionID, description, sectionID)
	// core.LogInfo("Command ID result: %d, Custom action result: %d", cmdID, int(caResult))

	return cmdID, nil
}

// RegisterMainAction registers an action in the main section
func RegisterMainAction(actionID string, description string) (int, error) {
	return RegisterCustomAction(actionID, description, SectionMain)
}
