package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

var (
	// Global mutex for UI operations
	uiMutex sync.Mutex
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
func GetUserInputs(title string, fields []string, defaults []string) ([]string, error) {
	// Use a global mutex to ensure only one dialog can be shown at a time
	uiMutex.Lock()
	defer uiMutex.Unlock()

	// Lock the OS thread to ensure UI operations happen on the same thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

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

	// Very important: Use a proper buffer with plenty of space
	bufferSize := 8192 // Much larger buffer to handle any clipboard content

	// Allocate buffer using malloc to ensure it's modifiable
	cValuesBuf := (*C.char)(C.malloc(C.size_t(bufferSize)))
	if cValuesBuf == nil {
		return nil, fmt.Errorf("failed to allocate memory for input buffer")
	}
	defer C.free(unsafe.Pointer(cValuesBuf))

	// Zero the entire buffer first
	C.memset(unsafe.Pointer(cValuesBuf), 0, C.size_t(bufferSize))

	// Copy default values into the buffer
	if len(defaultValues) > 0 {
		C.strncpy(cValuesBuf, C.CString(defaultValues), C.size_t(bufferSize-1))
	}

	// Ensure null termination
	*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cValuesBuf)) + uintptr(bufferSize-1))) = 0

	// Log what we're about to do
	// core.LogDebug("Showing GetUserInputs dialog: %s", title)

	// Ensure we're on the main thread
	// mainThread := runtime.NumGoroutine() // Just a debug helper to verify thread
	// core.LogDebug("Running on goroutine #%d", mainThread)

	// Call GetUserInputs
	result := C.plugin_bridge_call_get_user_inputs(
		getUserInputsPtr,
		cTitle,
		C.int(len(fields)),
		cCaptions,
		cValuesBuf,
		C.int(bufferSize),
	)

	// Check result
	if !bool(result) {
		// core.LogInfo("User cancelled the dialog")
		return nil, fmt.Errorf("user cancelled the dialog")
	}

	// Safely convert the buffer to a Go string
	goValues := C.GoString(cValuesBuf)
	// core.LogInfo("Dialog completed with result: %s", goValues)

	// Split by comma and return
	values := strings.Split(goValues, ",")

	return values, nil
}

// MessageBox is a simplified function that shows a message box with OK button
// We deliberately avoid using channels, goroutines or complex thread handling
func MessageBox(text string, title string) error {
	// Lock the UI mutex to prevent concurrent UI operations
	uiMutex.Lock()
	defer uiMutex.Unlock()

	// core.LogInfo("[MessageBox] %s: %s", title, text)

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer
	cFuncName := C.CString("ShowMessageBox")
	defer C.free(unsafe.Pointer(cFuncName))

	showMessageBoxPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if showMessageBoxPtr == nil {
		return fmt.Errorf("could not get ShowMessageBox function pointer")
	}

	// Prepare the parameters
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Call ShowMessageBox
	C.plugin_bridge_call_show_message_box(
		showMessageBoxPtr,
		cText,
		cTitle,
		C.int(MB_OK),
	)

	// core.LogDebug("Message box %s completed", title)
	return nil
}

// YesNoBox is a simplified function that shows a Yes/No dialog
// Returns true if Yes was clicked, false otherwise
func YesNoBox(text string, title string) (bool, error) {
	// Lock the UI mutex to prevent concurrent UI operations
	uiMutex.Lock()
	defer uiMutex.Unlock()

	// Always log the question
	// core.LogInfo("[QUESTION] %s: %s", title, text)

	if !initialized {
		return false, fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer
	cFuncName := C.CString("ShowMessageBox")
	defer C.free(unsafe.Pointer(cFuncName))

	showMessageBoxPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if showMessageBoxPtr == nil {
		return false, fmt.Errorf("could not get ShowMessageBox function pointer")
	}

	// Prepare the parameters
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Call ShowMessageBox
	result := C.plugin_bridge_call_show_message_box(
		showMessageBoxPtr,
		cText,
		cTitle,
		C.int(MB_YESNO),
	)

	// core.LogDebug("Yes/No box %s completed with result %d", title, int(result))
	return int(result) == IDYES, nil
}
