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

	// Log to the REAPER console
	reaper.ConsoleLog("----------------------------------------------------------")
	reaper.ConsoleLog("Hello from Go REAPER extension!")
	reaper.ConsoleLog("----------------------------------------------------------")

	// Register some distinctive actions in the Main section
	reaper.RegisterMainAction("GO_PURPLE_DRAGON", "Go: Purple Dragon Attack")
	reaper.RegisterMainAction("GO_RAINBOW_LASER", "Go: Rainbow Laser Beam")

	reaper.ConsoleLog("----------------------------------------------------------")
	reaper.ConsoleLog("Go plugin loaded successfully! Check Actions > Show action list...")
	reaper.ConsoleLog("- Main section: Look for actions starting with 'Go:'")
	reaper.ConsoleLog("----------------------------------------------------------")

	fmt.Println("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
