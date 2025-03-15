package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../../sdk
#include "../c/bridge.h"
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"go-reaper/src/pkg/logger"
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

	// Use the batch function to get all parameters at once
	parameters, err := BatchGetFXParameters(track, fxIndex)
	if err != nil {
		return result, fmt.Errorf("failed to batch get FX parameters: %v", err)
	}

	// Add parameters to result
	result.Parameters = parameters

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

// BatchGetFXParameters gets all parameters for an FX in a single call
// This reduces the number of C-Go crossings dramatically
func BatchGetFXParameters(track unsafe.Pointer, fxIndex int) ([]FXParameter, error) {
	if !initialized {
		return nil, fmt.Errorf("REAPER functions not initialized")
	}

	// Allocate memory for parameters (we'll allow up to 512 parameters)
	const maxParams = 512
	paramData := (*C.fx_param_t)(C.malloc(C.size_t(maxParams) * C.size_t(unsafe.Sizeof(C.fx_param_t{}))))
	if paramData == nil {
		return nil, fmt.Errorf("failed to allocate memory for parameter data")
	}
	defer C.free(unsafe.Pointer(paramData))

	// Allocate memory for parameter count
	paramCount := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0)))))
	if paramCount == nil {
		return nil, fmt.Errorf("failed to allocate memory for parameter count")
	}
	defer C.free(unsafe.Pointer(paramCount))

	// Call the C function to get all parameters
	result := C.plugin_bridge_batch_get_fx_parameters(
		track,
		C.int(fxIndex),
		paramData,
		C.int(maxParams),
		paramCount,
	)

	if !bool(result) {
		return nil, fmt.Errorf("failed to get FX parameters")
	}

	// Convert C parameter data to Go slice
	count := int(*paramCount)
	parameters := make([]FXParameter, count)

	// Create a slice of paramData
	// This creates a Go slice that points to the C array without copying it
	paramSlice := (*[maxParams]C.fx_param_t)(unsafe.Pointer(paramData))[:count:count]

	// Copy parameter data to Go slice
	for i := 0; i < count; i++ {
		parameters[i] = FXParameter{
			Index:          i,
			Name:           C.GoString(&paramSlice[i].name[0]),
			Value:          float64(paramSlice[i].value),
			FormattedValue: C.GoString(&paramSlice[i].formatted[0]),
			Min:            float64(paramSlice[i].min),
			Max:            float64(paramSlice[i].max),
		}
	}

	return parameters, nil
}

// GetTrackFXParamFormattedValueWithValue gets the formatted string for a specific parameter value
// This is useful to get the formatted display of a value without actually changing the parameter
func GetTrackFXParamFormattedValueWithValue(track unsafe.Pointer, fxIndex int, paramIndex int, value float64) (string, error) {
	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("TrackFX_FormatParamValue")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return "", fmt.Errorf("could not get TrackFX_FormatParamValue function pointer")
	}

	// Allocate buffer for the formatted value
	buf := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(buf))

	C.plugin_bridge_call_track_fx_format_param_value(
		getFuncPtr,
		track,
		C.int(fxIndex),
		C.int(paramIndex),
		C.double(value),
		buf,
		C.int(256),
	)

	return C.GoString(buf), nil
}

// ParameterFormatRequest defines a request to format a parameter value
type ParameterFormatRequest struct {
	FXIndex    int     // FX index
	ParamIndex int     // Parameter index
	Value      float64 // Value to format
}

// BatchFormatFXParameters formats multiple parameter values in a single call
// This is much more efficient than making multiple CGo transitions
func BatchFormatFXParameters(track unsafe.Pointer, requests []ParameterFormatRequest) ([]string, error) {
	if !initialized {
		return nil, fmt.Errorf("REAPER functions not initialized")
	}

	count := len(requests)
	if count == 0 {
		return []string{}, nil
	}

	// Allocate memory for parameter format requests
	paramData := (*C.fx_param_format_t)(C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.fx_param_format_t{}))))
	if paramData == nil {
		return nil, fmt.Errorf("failed to allocate memory for parameter data")
	}
	defer C.free(unsafe.Pointer(paramData))

	// Create a slice view of the C array
	paramSlice := (*[1 << 30]C.fx_param_format_t)(unsafe.Pointer(paramData))[:count:count]

	// Fill in the parameter data
	for i, req := range requests {
		paramSlice[i].fx_index = C.int(req.FXIndex)
		paramSlice[i].param_index = C.int(req.ParamIndex)
		paramSlice[i].value = C.double(req.Value)
	}

	// Call the batch function
	result := C.plugin_bridge_batch_format_fx_parameters(track, paramData, C.int(count))
	if !bool(result) {
		return nil, fmt.Errorf("failed to format parameters")
	}

	// Extract the formatted values
	formattedValues := make([]string, count)
	for i := 0; i < count; i++ {
		formattedValues[i] = C.GoString(&paramSlice[i].formatted[0])
	}

	return formattedValues, nil
}

// GetMinMaxFormatted gets the formatted min and max values for a parameter
func GetMinMaxFormatted(track unsafe.Pointer, fxIndex, paramIndex int) (minFormatted, maxFormatted string, err error) {
	_, min, max, err := GetTrackFXParamValueWithRange(track, fxIndex, paramIndex)
	if err != nil {
		return "", "", fmt.Errorf("failed to get parameter range: %v", err)
	}

	// Create format requests for min and max values
	requests := []ParameterFormatRequest{
		{FXIndex: fxIndex, ParamIndex: paramIndex, Value: min},
		{FXIndex: fxIndex, ParamIndex: paramIndex, Value: max},
	}

	// Format both values in a single call
	formatted, err := BatchFormatFXParameters(track, requests)
	if err != nil {
		return "", "", fmt.Errorf("failed to format min/max values: %v", err)
	}

	if len(formatted) != 2 {
		return "", "", fmt.Errorf("unexpected number of formatted values returned")
	}

	return formatted[0], formatted[1], nil
}

// GetFXParametersWithMinMax retrieves all parameters for a specific FX with min/max formatted values
func GetFXParametersWithMinMax(track unsafe.Pointer, fxIndex int) (FXInfo, error) {
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

	// Use the batch function to get all parameters at once
	parameters, err := BatchGetFXParameters(track, fxIndex)
	if err != nil {
		return result, fmt.Errorf("failed to batch get FX parameters: %v", err)
	}

	// Create batch format requests for min/max values of all parameters
	formatRequests := make([]ParameterFormatRequest, len(parameters)*2)
	for i, param := range parameters {
		// Request for min value
		formatRequests[i*2] = ParameterFormatRequest{
			FXIndex:    fxIndex,
			ParamIndex: param.Index,
			Value:      param.Min,
		}

		// Request for max value
		formatRequests[i*2+1] = ParameterFormatRequest{
			FXIndex:    fxIndex,
			ParamIndex: param.Index,
			Value:      param.Max,
		}
	}

	// Get all formatted min/max values in a single call
	formattedValues, err := BatchFormatFXParameters(track, formatRequests)
	if err != nil {
		// Continue even if formatting fails
		logger.Warning("Failed to format min/max values: %v", err)
		result.Parameters = parameters
		return result, nil
	}

	// Add min/max formatted values to parameters
	enhancedParams := make([]FXParameter, len(parameters))
	for i, param := range parameters {
		enhancedParams[i] = param
		if i*2+1 < len(formattedValues) {
			enhancedParams[i].MinFormatted = formattedValues[i*2]
			enhancedParams[i].MaxFormatted = formattedValues[i*2+1]
		}
	}

	result.Parameters = enhancedParams
	return result, nil
}

// BatchSetFXParameters applies multiple parameter changes in a single call
// This is much more efficient than making multiple CGo transitions
func BatchSetFXParameters(track unsafe.Pointer, changes []ParameterChange) error {
	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	count := len(changes)
	if count == 0 {
		return nil // Nothing to do
	}

	// Allocate memory for parameter changes
	changeData := (*C.fx_param_change_t)(C.malloc(C.size_t(count) * C.size_t(unsafe.Sizeof(C.fx_param_change_t{}))))
	if changeData == nil {
		return fmt.Errorf("failed to allocate memory for parameter changes")
	}
	defer C.free(unsafe.Pointer(changeData))

	// Create a slice view of the C array
	changeSlice := (*[1 << 30]C.fx_param_change_t)(unsafe.Pointer(changeData))[:count:count]

	// Fill in the change data
	for i, change := range changes {
		changeSlice[i].fx_index = C.int(change.FXIndex)
		changeSlice[i].param_index = C.int(change.ParamIndex)
		changeSlice[i].value = C.double(change.Value)
	}

	// Call the batch function
	result := C.plugin_bridge_batch_set_fx_parameters(track, changeData, C.int(count))
	if !bool(result) {
		return fmt.Errorf("failed to apply parameter changes")
	}

	logger.Debug("Applied %d parameter changes successfully", count)
	return nil
}

// BatchSetFXParametersWithUndo applies multiple parameter changes in a single call
// and wraps the changes in an undo block
func BatchSetFXParametersWithUndo(track unsafe.Pointer, changes []ParameterChange, undoLabel string) error {
	// Start undo block
	if err := BeginUndoBlock(undoLabel); err != nil {
		logger.Warning("Could not start undo block: %v", err)
		// Continue anyway, just without undo support
	}

	// Apply the changes
	err := BatchSetFXParameters(track, changes)

	// End undo block (even if there was an error)
	if endErr := EndUndoBlock(undoLabel, 0); endErr != nil {
		logger.Warning("Could not end undo block: %v", endErr)
	}

	return err
}

// TrackFX_GetParameterStepSizes gets the step sizes and toggle status for a parameter
func TrackFX_GetParameterStepSizes(track unsafe.Pointer, fxIndex int, paramIndex int,
	step, smallStep, largeStep *float64, isToggle *bool) bool {
	if !initialized {
		return false
	}

	cFuncName := C.CString("TrackFX_GetParameterStepSizes")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		logger.Warning("Could not get TrackFX_GetParameterStepSizes function pointer")
		return false
	}

	// Convert Go pointers to C pointers
	var cStep, cSmallStep, cLargeStep *C.double
	var cIsToggle *C.bool

	// Only create C pointers if Go pointers are not nil
	if step != nil {
		cStep = (*C.double)(unsafe.Pointer(step))
	}
	if smallStep != nil {
		cSmallStep = (*C.double)(unsafe.Pointer(smallStep))
	}
	if largeStep != nil {
		cLargeStep = (*C.double)(unsafe.Pointer(largeStep))
	}
	if isToggle != nil {
		cIsToggle = (*C.bool)(unsafe.Pointer(isToggle))
	}

	result := C.plugin_bridge_call_track_fx_get_parameter_step_sizes(
		getFuncPtr,
		track,
		C.int(fxIndex),
		C.int(paramIndex),
		cStep,
		cSmallStep,
		cLargeStep,
		cIsToggle,
	)

	return bool(result)
}
