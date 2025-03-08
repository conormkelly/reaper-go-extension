package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// GetSelectedTrack returns the first selected track in the current project
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
