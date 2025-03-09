package main

/*
#cgo CFLAGS: -I${SRCDIR}/c -I${SRCDIR}/sdk
#include "c/bridge.h"
*/
import "C"
import (
	"unsafe"

	"go-reaper/actions"
	"go-reaper/core"
	"go-reaper/pkg/logger"
)

//export GoReaperPluginEntry
func GoReaperPluginEntry(hInstance unsafe.Pointer, rec unsafe.Pointer) C.int {
	// If rec is null, REAPER is unloading the plugin
	if rec == nil {
		// Close any open UI windows
		actions.CloseNativeWindow()
		actions.CloseKeyringWindow()

		// Perform cleanup tasks including logging shutdown
		logger.Cleanup()
		return 0
	}

	// Initialize logging system
	logger.Initialize()

	// Initialize core functionality
	if err := core.Initialize(hInstance, rec); err != nil {
		logger.Error("Failed to initialize REAPER: %v", err)
		return 0
	}

	// Register all actions
	if err := actions.RegisterAll(); err != nil {
		logger.Error("Failed to register actions: %v", err)
		return 0
	}

	logger.Info("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
