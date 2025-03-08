package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

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

	jsonData, err := json.Marshal(fxInfos)
	if err != nil {
		return "", fmt.Errorf("failed to marshal FX info to JSON: %v", err)
	}

	return string(jsonData), nil
}
