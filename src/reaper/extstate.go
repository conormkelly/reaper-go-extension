package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../../sdk
#include "../c/bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"sync"
	"unsafe"
)

// extStateMutex protects access to the ExtState functions
var extStateMutex sync.Mutex

// GetExtState gets an extended state value
func GetExtState(section, key string) (string, error) {
	extStateMutex.Lock()
	defer extStateMutex.Unlock()

	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer
	cFuncName := C.CString("GetExtState")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return "", fmt.Errorf("could not get GetExtState function pointer")
	}

	// Prepare the parameters
	cSection := C.CString(section)
	defer C.free(unsafe.Pointer(cSection))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	// Call GetExtState through our bridge
	result := C.plugin_bridge_call_get_ext_state(getFuncPtr, cSection, cKey)
	if result == nil {
		return "", nil
	}

	// Convert to Go string
	return C.GoString(result), nil
}

// SetExtState sets an extended state value
func SetExtState(section, key, value string, persist bool) error {
	extStateMutex.Lock()
	defer extStateMutex.Unlock()

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer
	cFuncName := C.CString("SetExtState")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return fmt.Errorf("could not get SetExtState function pointer")
	}

	// Prepare the parameters
	cSection := C.CString(section)
	defer C.free(unsafe.Pointer(cSection))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	// Call SetExtState through our bridge
	persistInt := C.int(0)
	if persist {
		persistInt = C.int(1)
	}

	C.plugin_bridge_call_set_ext_state(getFuncPtr, cSection, cKey, cValue, persistInt)

	// Verify with HasExtState for debugging
	// hasKey, err := hasExtStateLocked(section, key)
	// if err != nil {
	// 	logger.Warning("SetExtState: Verification - HasExtState failed: %v", err)
	// } else if !hasKey {
	// 	logger.Warning("SetExtState: Verification - Key not found after setting!")
	// } else {
	// 	logger.Debug("SetExtState: Verification - Key found after setting")
	// }

	return nil
}

// HasExtState checks if an ext state value exists
func HasExtState(section, key string) (bool, error) {
	extStateMutex.Lock()
	defer extStateMutex.Unlock()

	if !initialized {
		return false, fmt.Errorf("REAPER functions not initialized")
	}

	return hasExtStateLocked(section, key)
}

func hasExtStateLocked(section, key string) (bool, error) {
	// Get the function pointer
	cFuncName := C.CString("HasExtState")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return false, fmt.Errorf("could not get HasExtState function pointer")
	}

	// Prepare the parameters
	cSection := C.CString(section)
	defer C.free(unsafe.Pointer(cSection))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	// Call HasExtState through our bridge
	result := C.plugin_bridge_call_has_ext_state(getFuncPtr, cSection, cKey)

	return bool(result), nil
}

// DeleteExtState deletes an ext state value
func DeleteExtState(section, key string) error {
	extStateMutex.Lock()
	defer extStateMutex.Unlock()

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer
	cFuncName := C.CString("DeleteExtState")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return fmt.Errorf("could not get DeleteExtState function pointer")
	}

	// Prepare the parameters
	cSection := C.CString(section)
	defer C.free(unsafe.Pointer(cSection))

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	// Call DeleteExtState through our bridge
	C.plugin_bridge_call_delete_ext_state(getFuncPtr, cSection, cKey)
	logger.Debug("Deleted ext state: [%s]%s", section, key)

	return nil
}
