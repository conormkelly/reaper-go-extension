package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

// Dialog button constants
const (
	MB_OK               = 0
	MB_OKCANCEL         = 1
	MB_ABORTRETRYIGNORE = 2
	MB_YESNOCANCEL      = 3
	MB_YESNO            = 4
	MB_RETRYCANCEL      = 5

	// Button IDs
	IDOK     = 1
	IDCANCEL = 2
	IDABORT  = 3
	IDRETRY  = 4
	IDIGNORE = 5
	IDYES    = 6
	IDNO     = 7
)

// GetUserInputs shows a dialog with fields for user input
// title: the dialog title
// fields: array of field labels
// defaults: array of default values for fields (same length as fields)
// Returns: array of user input values, error if dialog was cancelled
func GetUserInputs(title string, fields []string, defaults []string) ([]string, error) {
	if !initialized {
		return nil, fmt.Errorf("REAPER functions not initialized")
	}

	// Ensure defaults has the same length as fields
	if len(defaults) != len(fields) {
		defaults = make([]string, len(fields))
	}

	// Get the GetUserInputs function
	cFuncName := C.CString("GetUserInputs")
	defer C.free(unsafe.Pointer(cFuncName))

	getUserInputsPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getUserInputsPtr == nil {
		return nil, fmt.Errorf("could not get GetUserInputs function pointer")
	}

	// Prepare parameters
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Join field names with commas
	captions := strings.Join(fields, ",")
	cCaptions := C.CString(captions)
	defer C.free(unsafe.Pointer(cCaptions))

	// Join default values with commas
	defaultValues := strings.Join(defaults, ",")

	// Allocate a buffer for the values
	// Add extra space for safety
	bufferSize := 1024
	if len(defaultValues)*2 > bufferSize {
		bufferSize = len(defaultValues) * 2
	}

	cValues := C.CString(defaultValues)
	defer C.free(unsafe.Pointer(cValues))

	// Call GetUserInputs
	result := C.plugin_bridge_call_get_user_inputs(
		getUserInputsPtr,
		cTitle,
		C.int(len(fields)),
		cCaptions,
		cValues,
		C.int(bufferSize),
	)

	// If user cancelled, return an error
	if !bool(result) {
		return nil, fmt.Errorf("user cancelled the dialog")
	}

	// Parse result
	goValues := C.GoString(cValues)
	return strings.Split(goValues, ","), nil
}

// ShowMessageBox displays a standard message box and returns the button clicked
func ShowMessageBox(text string, title string, messageType int) (int, error) {
	if !initialized {
		return 0, fmt.Errorf("REAPER functions not initialized")
	}

	// Get the ShowMessageBox function
	cFuncName := C.CString("ShowMessageBox")
	defer C.free(unsafe.Pointer(cFuncName))

	showMessageBoxPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if showMessageBoxPtr == nil {
		return 0, fmt.Errorf("could not get ShowMessageBox function pointer")
	}

	// Prepare parameters
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Call ShowMessageBox
	result := C.plugin_bridge_call_show_message_box(
		showMessageBoxPtr,
		cText,
		cTitle,
		C.int(messageType),
	)

	return int(result), nil
}

// MessageBox is a convenience function that displays a message box with an OK button
func MessageBox(text string, title string) error {
	_, err := ShowMessageBox(text, title, MB_OK)
	return err
}

// YesNoBox is a convenience function that displays a message box with Yes/No buttons
// Returns true if Yes was clicked, false if No was clicked
func YesNoBox(text string, title string) (bool, error) {
	result, err := ShowMessageBox(text, title, MB_YESNO)
	if err != nil {
		return false, err
	}
	return result == IDYES, nil
}
