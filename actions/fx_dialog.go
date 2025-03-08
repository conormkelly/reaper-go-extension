package actions

import (
	"fmt"
	"go-reaper/reaper"
)

// RegisterFXDialog registers the LLM FX Assistant dialog action
func RegisterFXDialog() error {
	FXDialogID, err := reaper.RegisterMainAction("GO_FX_DIALOG", "Go: LLM FX Assistant")
	if err != nil {
		return fmt.Errorf("failed to register LLM FX Assistant: %v", err)
	}

	reaper.ConsoleLog(fmt.Sprintf("LLM FX Assistant registered with ID: %d", FXDialogID))
	reaper.SetActionHandler("GO_FX_DIALOG", handleFXDialog)
	return nil
}

// handleFXDialog opens the LLM FX Assistant dialog
func handleFXDialog() {
	reaper.ConsoleLog("----- LLM FX Assistant Activated -----")

	// Get track info
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Error getting track info: %v", err))
		reaper.ConsoleLog("Please select a track before using the LLM FX Assistant.")
		return
	}

	// Check if track has FX
	if trackInfo.NumFX == 0 {
		reaper.ConsoleLog("Selected track has no FX. Please add FX to the track before using the LLM FX Assistant.")
		return
	}

	// Get FX list
	fxList, err := reaper.GetTrackFXList(trackInfo.MediaTrack)
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Error getting FX list: %v", err))
		return
	}

	reaper.ConsoleLog(fmt.Sprintf("Found %d FX on track.", len(fxList)))

	// For debugging, log the FX information we would show in the dialog
	reaper.ConsoleLog("\nFX on the selected track:")
	for i, fx := range fxList {
		reaper.ConsoleLog(fmt.Sprintf("%d. %s", i+1, fx.Name))
	}

	reaper.ConsoleLog("\nPhase 1 debug complete. UI will be implemented after further research.")

	// Instead of trying to create a dialog, just log the information for now
}
