package actions

import (
	"fmt"
	"go-reaper/reaper"
)

// RegisterFXPrototype registers the FX Prototype action
func RegisterFXPrototype() error {
	FXPrototypeID, err := reaper.RegisterMainAction("GO_FX_PROTOTYPE", "Go: LLM FX Prototype")
	if err != nil {
		return fmt.Errorf("failed to register LLM FX Prototype: %v", err)
	}

	reaper.ConsoleLog(fmt.Sprintf("LLM FX Prototype registered with ID: %d", FXPrototypeID))
	reaper.SetActionHandler("GO_FX_PROTOTYPE", handleFXPrototype)
	return nil
}

// handleFXPrototype handles the FX Prototype action
func handleFXPrototype() {
	reaper.ConsoleLog("----- LLM FX Prototype Activated -----")

	// Get FX info as structured data
	fxInfosJSON, err := reaper.GetCurrentFXInfoJSON()
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Error getting FX info: %v", err))
		return
	}

	// Log the JSON data
	reaper.ConsoleLog("FX Parameters as JSON:")
	reaper.ConsoleLog(fxInfosJSON)

	// Also log the parameters in a readable format for reference
	err = reaper.LogCurrentFX()
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Error logging FX details: %v", err))
		return
	}

	reaper.ConsoleLog("LLM FX Prototype step 1 complete! The FX parameters have been collected.")
	reaper.ConsoleLog("Future steps: Add user input dialog and LLM integration.")
}
