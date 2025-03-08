package main

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/sdk
#include "reaper_plugin_bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"go-reaper/actions"
	"go-reaper/core"
)

//export GoReaperPluginEntry
func GoReaperPluginEntry(hInstance unsafe.Pointer, rec unsafe.Pointer) C.int {
	fmt.Println("Go REAPER plugin entry called")

	// If rec is null, REAPER is unloading the plugin
	if rec == nil {
		fmt.Println("Go plugin unloading")
		return 0
	}

	// Initialize core functionality
	if err := core.Initialize(hInstance, rec); err != nil {
		fmt.Printf("Failed to initialize REAPER: %v\n", err)
		return 0
	}

	// Register all actions
	if err := actions.RegisterAll(); err != nil {
		fmt.Printf("Failed to register actions: %v\n", err)
		return 0
	}

	fmt.Println("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
