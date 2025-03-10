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

var (
	// Mutex for undo operations
	undoMutex sync.Mutex
)

// BeginUndoBlock starts a new undo block with the specified description
// Uses Undo_BeginBlock2 with NULL for the active project
func BeginUndoBlock(description string) error {
	undoMutex.Lock() // Will be released in EndUndoBlock

	if !initialized {
		undoMutex.Unlock() // Release lock if we're returning early
		return fmt.Errorf("REAPER functions not initialized")
	}

	// Get the function pointer for Undo_BeginBlock2
	cFuncName := C.CString("Undo_BeginBlock2")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		undoMutex.Unlock() // Release lock if we're returning early
		return fmt.Errorf("could not get Undo_BeginBlock2 function pointer - REAPER version may be too old")
	}

	// Call through our bridge with NULL for the active project
	C.plugin_bridge_call_undo_begin_block2(getFuncPtr, nil)
	logger.Debug("Started undo block: %s", description)

	return nil
}

// EndUndoBlock ends the current undo block with the specified description
// Uses Undo_EndBlock2 with NULL for the active project
func EndUndoBlock(description string, flags int) error {
	defer undoMutex.Unlock() // Release lock from BeginUndoBlock

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	// Convert description to C string
	cDesc := C.CString(description)
	defer C.free(unsafe.Pointer(cDesc))

	// Get the function pointer for Undo_EndBlock2
	cFuncName := C.CString("Undo_EndBlock2")
	defer C.free(unsafe.Pointer(cFuncName))

	getFuncPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getFuncPtr == nil {
		return fmt.Errorf("could not get Undo_EndBlock2 function pointer - REAPER version may be too old")
	}

	// Call through our bridge with NULL for the active project
	C.plugin_bridge_call_undo_end_block2(getFuncPtr, nil, cDesc, C.int(flags))
	logger.Debug("Ended undo block: %s (flags: %d)", description, flags)

	return nil
}
