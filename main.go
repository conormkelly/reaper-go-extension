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

// Add the action handler for LLM FX Prototype
func handleFXPrototype() {
	reaper.ConsoleLog("----- LLM FX Prototype Activated -----")

	err := reaper.LogCurrentFX()
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Error: %v", err))
		return
	}

	reaper.ConsoleLog("LLM FX Prototype step 1 complete! The FX parameters have been logged.")
	reaper.ConsoleLog("Future steps: Add user input dialog and LLM integration.")
}

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

	// Register our new LLM FX Prototype action
	FXPrototypeID, err := reaper.RegisterMainAction("GO_FX_PROTOTYPE", "Go: LLM FX Prototype")
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Failed to register LLM FX Prototype: %v", err))
	} else {
		reaper.ConsoleLog(fmt.Sprintf("LLM FX Prototype registered with ID: %d", FXPrototypeID))
		// Register the handler for the action
		reaper.SetActionHandler("GO_FX_PROTOTYPE", handleFXPrototype)
	}

	reaper.ConsoleLog("----------------------------------------------------------")
	reaper.ConsoleLog("Go plugin loaded successfully! Check Actions > Show action list...")
	reaper.ConsoleLog("- Main section: Look for actions starting with 'Go:'")
	reaper.ConsoleLog("----------------------------------------------------------")

	fmt.Println("Go plugin loaded successfully!")
	return 1
}

// Required main function for Go builds
func main() {}
