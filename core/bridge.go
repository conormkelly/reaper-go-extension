package core

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../sdk
#include "../c/bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"go-reaper/pkg/logger"
	"go-reaper/reaper"
)

// Initialize initializes the core plugin functionality
func Initialize(hInstance unsafe.Pointer, rec unsafe.Pointer) error {
	// Initialize our REAPER API wrapper
	if err := reaper.Initialize(rec); err != nil {
		return fmt.Errorf("failed to initialize REAPER API: %v", err)
	}

	// Log to the REAPER console
	logger.Debug("----------------------------------------------------------")
	logger.Debug("Hello from Go REAPER extension!")
	logger.Debug("----------------------------------------------------------")

	return nil
}
