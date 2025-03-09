package core

/*
#cgo CFLAGS: -I${SRCDIR}/.. -I${SRCDIR}/../sdk
#include "../reaper_plugin_bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"go-reaper/reaper"
)

// Initialize initializes the core plugin functionality
func Initialize(hInstance unsafe.Pointer, rec unsafe.Pointer) error {
	// Initialize our REAPER API wrapper
	if err := reaper.Initialize(rec); err != nil {
		return fmt.Errorf("failed to initialize REAPER API: %v", err)
	}

	// Log to the REAPER console
	LogDebug("----------------------------------------------------------")
	LogDebug("Hello from Go REAPER extension!")
	LogDebug("----------------------------------------------------------")

	return nil
}
