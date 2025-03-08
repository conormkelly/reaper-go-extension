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
