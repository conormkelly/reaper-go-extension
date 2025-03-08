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
	"encoding/json"
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

// Store a map of action handlers
var actionHandlers map[string]func()

// ActionHandler defines a function type for handling actions
type ActionHandler func()

// Initialize action handlers map
func initActionHandlers() {
	actionHandlers = make(map[string]func())
}

// SetActionHandler associates a function with an action ID
func SetActionHandler(actionID string, handler func()) {
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
			ConsoleLog(fmt.Sprintf("GoReaper action triggered! Command ID: %d (%s)", int(commandId), actionID))

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
			ConsoleLog(fmt.Sprintf("GoReaper action triggered! Command ID: %d (%s) (via hookcommand2)", int(commandId), actionID))

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

// TrackFX related functions
func GetSelectedTrack() (unsafe.Pointer, error) {
	if !initialized {
		return nil, fmt.Errorf("REAPER functions not initialized")
	}

	// Get the GetFunc pointer using our bridge function
	getFuncPtr := C.plugin_bridge_get_get_func()

	// Get the GetSetObjectState function pointer
	cFuncName := C.CString("GetSelectedTrack")
	defer C.free(unsafe.Pointer(cFuncName))

	trackFuncPtr := C.plugin_bridge_call_get_func(getFuncPtr, cFuncName)
	if trackFuncPtr == nil {
		return nil, fmt.Errorf("could not get GetSelectedTrack function pointer")
	}

	// Call GetSelectedTrack(0, 0) - first project, first selected track
	track := C.plugin_bridge_call_get_selected_track(trackFuncPtr, 0, 0)
	if track == nil {
		return nil, fmt.Errorf("no track selected")
	}

	return track, nil
}

// GetTrackFXCount gets the number of FX on a track
func GetTrackFXCount(track unsafe.Pointer) (int, error) {
	if !initialized {
		return 0, fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetCount")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return 0, fmt.Errorf("could not get TrackFX_GetCount function pointer")
	}

	count := C.plugin_bridge_call_track_fx_get_count(getFuncPtr, track)
	return int(count), nil
}

// GetTrackFXName gets the name of an FX
func GetTrackFXName(track unsafe.Pointer, fxIndex int) (string, error) {
	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetFXName")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return "", fmt.Errorf("could not get TrackFX_GetFXName function pointer")
	}

	// Allocate buffer for the name
	buf := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(buf))

	C.plugin_bridge_call_track_fx_get_name(getFuncPtr, track, C.int(fxIndex), buf, C.int(256))

	return C.GoString(buf), nil
}

// GetTrackFXParamCount gets the number of parameters for an FX
func GetTrackFXParamCount(track unsafe.Pointer, fxIndex int) (int, error) {
	if !initialized {
		return 0, fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetNumParams")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return 0, fmt.Errorf("could not get TrackFX_GetNumParams function pointer")
	}

	count := C.plugin_bridge_call_track_fx_get_param_count(getFuncPtr, track, C.int(fxIndex))
	return int(count), nil
}

// GetTrackFXParamName gets the name of a parameter
func GetTrackFXParamName(track unsafe.Pointer, fxIndex int, paramIndex int) (string, error) {
	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetParamName")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return "", fmt.Errorf("could not get TrackFX_GetParamName function pointer")
	}

	// Allocate buffer for the name
	buf := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(buf))

	C.plugin_bridge_call_track_fx_get_param_name(getFuncPtr, track, C.int(fxIndex), C.int(paramIndex), buf, C.int(256))

	return C.GoString(buf), nil
}

// GetTrackFXParamValue gets the normalized value (0.0-1.0) of a parameter
func GetTrackFXParamValue(track unsafe.Pointer, fxIndex int, paramIndex int) (float64, error) {
	if !initialized {
		return 0, fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetParam")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return 0, fmt.Errorf("could not get TrackFX_GetParam function pointer")
	}

	value := C.plugin_bridge_call_track_fx_get_param(getFuncPtr, track, C.int(fxIndex), C.int(paramIndex), nil, nil)
	return float64(value), nil
}

// GetTrackFXParamFormatted gets the formatted value of a parameter as a string
func GetTrackFXParamFormatted(track unsafe.Pointer, fxIndex int, paramIndex int) (string, error) {
	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetFormattedParamValue")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return "", fmt.Errorf("could not get TrackFX_GetFormattedParamValue function pointer")
	}

	// Allocate buffer for the formatted value
	buf := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(buf))

	C.plugin_bridge_call_track_fx_get_param_formatted(getFuncPtr, track, C.int(fxIndex), C.int(paramIndex), buf, C.int(256))

	return C.GoString(buf), nil
}

// SetTrackFXParamValue sets the value of a parameter
func SetTrackFXParamValue(track unsafe.Pointer, fxIndex int, paramIndex int, value float64) error {
	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_SetParam")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return fmt.Errorf("could not get TrackFX_SetParam function pointer")
	}

	C.plugin_bridge_call_track_fx_set_param(getFuncPtr, track, C.int(fxIndex), C.int(paramIndex), C.double(value))

	return nil
}

// LogFXParameters logs all parameters of an FX to the REAPER console
func LogFXParameters(track unsafe.Pointer, fxIndex int) error {
	// Get FX name
	fxName, err := GetTrackFXName(track, fxIndex)
	if err != nil {
		return fmt.Errorf("failed to get FX name: %v", err)
	}

	ConsoleLog(fmt.Sprintf("FX: %s", fxName))

	// Get parameter count
	paramCount, err := GetTrackFXParamCount(track, fxIndex)
	if err != nil {
		return fmt.Errorf("failed to get parameter count: %v", err)
	}

	ConsoleLog(fmt.Sprintf("Parameter count: %d", paramCount))

	// Log each parameter
	for i := 0; i < paramCount; i++ {
		paramName, err := GetTrackFXParamName(track, fxIndex, i)
		if err != nil {
			return fmt.Errorf("failed to get parameter name: %v", err)
		}

		paramValue, err := GetTrackFXParamValue(track, fxIndex, i)
		if err != nil {
			return fmt.Errorf("failed to get parameter value: %v", err)
		}

		paramFormatted, err := GetTrackFXParamFormatted(track, fxIndex, i)
		if err != nil {
			return fmt.Errorf("failed to get formatted parameter value: %v", err)
		}

		ConsoleLog(fmt.Sprintf("  Param #%d: %s = %.4f (%s)", i, paramName, paramValue, paramFormatted))
	}

	return nil
}

// LogCurrentFX logs parameters of the currently selected FX
func LogCurrentFX() error {
	// Get selected track
	track, err := GetSelectedTrack()
	if err != nil {
		return fmt.Errorf("failed to get selected track: %v", err)
	}

	// For now, just use the first FX on the track
	// In a more advanced version, we'd get the currently focused FX
	err = LogFXParameters(track, 0)
	if err != nil {
		return fmt.Errorf("failed to log FX parameters: %v", err)
	}

	return nil
}

// FXParameter represents a single parameter of an FX
type FXParameter struct {
	Index          int     `json:"index"`
	Name           string  `json:"name"`
	Value          float64 `json:"value"`          // Normalized value (0.0-1.0)
	FormattedValue string  `json:"formattedValue"` // Human-readable value
	Min            float64 `json:"min"`            // Minimum value
	Max            float64 `json:"max"`            // Maximum value
}

// FXInfo represents an FX and its parameters
type FXInfo struct {
	Index      int           `json:"index"`
	Name       string        `json:"name"`
	Parameters []FXParameter `json:"parameters"`
}

// GetTrackFXParamValueWithRange gets the normalized value and range of a parameter
func GetTrackFXParamValueWithRange(track unsafe.Pointer, fxIndex int, paramIndex int) (value, min, max float64, err error) {
	if !initialized {
		return 0, 0, 0, fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_GetParam")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return 0, 0, 0, fmt.Errorf("could not get TrackFX_GetParam function pointer")
	}

	// Allocate memory for min and max values
	minPtr := (*C.double)(C.malloc(C.size_t(unsafe.Sizeof(C.double(0)))))
	maxPtr := (*C.double)(C.malloc(C.size_t(unsafe.Sizeof(C.double(0)))))
	defer C.free(unsafe.Pointer(minPtr))
	defer C.free(unsafe.Pointer(maxPtr))

	value = float64(C.plugin_bridge_call_track_fx_get_param(getFuncPtr, track, C.int(fxIndex), C.int(paramIndex), minPtr, maxPtr))
	min = float64(*minPtr)
	max = float64(*maxPtr)

	return value, min, max, nil
}

// GetFXParameters retrieves all parameters for a specific FX
func GetFXParameters(track unsafe.Pointer, fxIndex int) (FXInfo, error) {
	result := FXInfo{
		Index:      fxIndex,
		Parameters: []FXParameter{},
	}

	// Get FX name
	fxName, err := GetTrackFXName(track, fxIndex)
	if err != nil {
		return result, fmt.Errorf("failed to get FX name: %v", err)
	}
	result.Name = fxName

	// Get parameter count
	paramCount, err := GetTrackFXParamCount(track, fxIndex)
	if err != nil {
		return result, fmt.Errorf("failed to get parameter count: %v", err)
	}

	// Collect each parameter
	for i := 0; i < paramCount; i++ {
		paramName, err := GetTrackFXParamName(track, fxIndex, i)
		if err != nil {
			return result, fmt.Errorf("failed to get parameter name: %v", err)
		}

		paramValue, min, max, err := GetTrackFXParamValueWithRange(track, fxIndex, i)
		if err != nil {
			return result, fmt.Errorf("failed to get parameter value: %v", err)
		}

		paramFormatted, err := GetTrackFXParamFormatted(track, fxIndex, i)
		if err != nil {
			return result, fmt.Errorf("failed to get formatted parameter value: %v", err)
		}

		param := FXParameter{
			Index:          i,
			Name:           paramName,
			Value:          paramValue,
			FormattedValue: paramFormatted,
			Min:            min,
			Max:            max,
		}

		result.Parameters = append(result.Parameters, param)
	}

	return result, nil
}

// GetCurrentFXInfo gets information about the FX on the currently selected track
func GetCurrentFXInfo() ([]FXInfo, error) {
	// Get selected track
	track, err := GetSelectedTrack()
	if err != nil {
		return nil, fmt.Errorf("failed to get selected track: %v", err)
	}

	// Get FX count
	fxCount, err := GetTrackFXCount(track)
	if err != nil {
		return nil, fmt.Errorf("failed to get FX count: %v", err)
	}

	// Gather info for all FX
	result := make([]FXInfo, 0, fxCount)
	for i := 0; i < fxCount; i++ {
		fxInfo, err := GetFXParameters(track, i)
		if err != nil {
			return nil, fmt.Errorf("failed to get FX parameters: %v", err)
		}
		result = append(result, fxInfo)
	}

	return result, nil
}

// GetCurrentFXInfoJSON returns the FX information as a JSON string
func GetCurrentFXInfoJSON() (string, error) {
	fxInfos, err := GetCurrentFXInfo()
	if err != nil {
		return "", err
	}

	// Import encoding/json package for JSON serialization
	// Make sure to add: import "encoding/json" at the top of reaper.go
	jsonData, err := json.Marshal(fxInfos)
	if err != nil {
		return "", fmt.Errorf("failed to marshal FX info to JSON: %v", err)
	}

	return string(jsonData), nil
}
