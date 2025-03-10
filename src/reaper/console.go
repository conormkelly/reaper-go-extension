package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../../sdk
#include "../c/bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

var (
	consoleMutex sync.Mutex
)

// ConsoleLog sends a message to the REAPER console
func ConsoleLog(message string) error {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

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

// ShowConsoleMsg is a direct wrapper for REAPER's ShowConsoleMsg function
// This version is safe to call from any goroutine
func ShowConsoleMsg(message string) error {
	consoleMutex.Lock()
	defer consoleMutex.Unlock()

	if !initialized {
		return fmt.Errorf("REAPER functions not initialized")
	}

	if showConsoleMsgPtr == nil {
		return fmt.Errorf("ShowConsoleMsg function not available")
	}

	// Run on main thread if not already there
	done := make(chan struct{})

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		cMessage := C.CString(message + "\n")
		defer C.free(unsafe.Pointer(cMessage))
		C.plugin_bridge_call_show_console_msg(showConsoleMsgPtr, cMessage)

		close(done)
	}()

	<-done
	return nil
}
