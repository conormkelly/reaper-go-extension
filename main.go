package main

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/sdk
#include "reaper_plugin_bridge.h"
*/
import "C"
import (
	"unsafe"

	"go-reaper/actions"
	"go-reaper/core"
)

//export GoReaperPluginEntry
func GoReaperPluginEntry(hInstance unsafe.Pointer, rec unsafe.Pointer) C.int {
	// If rec is null, REAPER is unloading the plugin
	if rec == nil {
		// Perform cleanup tasks including logging shutdown
		core.CleanupLogging()
		return 0
	}

	// Initialize logging system
	core.InitLogging()

	// Initialize core functionality
	if err := core.Initialize(hInstance, rec); err != nil {
		core.LogError("Failed to initialize REAPER: %s", err.Error())
		return 0
	}

	// Register all actions
	if err := actions.RegisterAll(); err != nil {
		core.LogError("Failed to register actions: %s", err.Error())
		return 0
	}

	core.LogInfo("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
