package main

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/sdk
#include "reaper_plugin_bridge.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"go-reaper/reaper"
)

//export GoReaperPluginEntry
func GoReaperPluginEntry(hInstance unsafe.Pointer, rec unsafe.Pointer) C.int {
	fmt.Println("Go REAPER plugin entry called")

	// If rec is null, REAPER is unloading the plugin
	if rec == nil {
		fmt.Println("Go plugin unloading")
		return 0
	}

	// Initialize our REAPER API wrapper
	if err := reaper.Initialize(rec); err != nil {
		fmt.Printf("Failed to initialize REAPER: %v\n", err)
		return 0
	}

	// Now we can use our nice Go functions instead of C calls directly
	if err := reaper.ConsoleLog("Hello from Go REAPER extension!"); err != nil {
		fmt.Printf("Error logging to console: %v\n", err)
	}

	// Register a basic command using our cleaner API
	result, err := reaper.RegisterAction("command_id", "GO_HELLO_WORLD")
	if err != nil {
		fmt.Printf("Error registering action: %v\n", err)
	} else {
		fmt.Printf("Registered command: GO_HELLO_WORLD, result: %d\n", result)
	}

	fmt.Println("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
