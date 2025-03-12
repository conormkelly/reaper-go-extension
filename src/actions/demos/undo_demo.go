package demo

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../../sdk
#include "../../c/bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
)

// RegisterUndoDemo registers the undo test action
func RegisterUndoDemo() error {
	actionID, err := reaper.RegisterMainAction("GO_UNDO_TEST", "Go: Undo Framework Test")
	if err != nil {
		return fmt.Errorf("failed to register undo test action: %v", err)
	}

	logger.Info("Undo Test action registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_UNDO_TEST", handleUndoTest)
	return nil
}

// handleUndoTest demonstrates the undo framework
func handleUndoTest() {
	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("Undo Test action triggered")

	// Get track information
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		reaper.MessageBox("Please select a track first", "Undo Test")
		return
	}

	// Get FX count
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to the track.", "Undo Test")
		return
	}

	// Test by adjusting a parameter in the first FX
	fxIndex := 0 // First FX

	// Get parameter count
	paramCount, err := reaper.GetTrackFXParamCount(trackInfo.MediaTrack, fxIndex)
	if err != nil {
		logger.Error("Failed to get parameter count: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to get parameter count: %v", err), "Undo Test")
		return
	}

	if paramCount == 0 {
		reaper.MessageBox("Selected FX has no parameters to adjust.", "Undo Test")
		return
	}

	paramIndex := 0 // First parameter

	// Get parameter name
	paramName, err := reaper.GetTrackFXParamName(trackInfo.MediaTrack, fxIndex, paramIndex)
	if err != nil {
		logger.Error("Failed to get parameter name: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to get parameter name: %v", err), "Undo Test")
		return
	}

	// Get current value
	value, err := reaper.GetTrackFXParamValue(trackInfo.MediaTrack, fxIndex, paramIndex)
	if err != nil {
		logger.Error("Failed to get parameter value: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to get parameter value: %v", err), "Undo Test")
		return
	}

	// Calculate new value (invert it within 0-1 range)
	newValue := 1.0 - value

	// Begin undo block
	if err = reaper.BeginUndoBlock("Undo Test"); err != nil {
		logger.Error("Failed to begin undo block: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to begin undo block: %v", err), "Undo Test")
		return
	}

	// Make the change
	err = reaper.SetTrackFXParamValue(trackInfo.MediaTrack, fxIndex, paramIndex, newValue)
	if err != nil {
		logger.Error("Failed to set parameter value: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to set parameter value: %v", err), "Undo Test")

		// End the undo block even on error
		if endErr := reaper.EndUndoBlock("Undo Test - Failed", 0); endErr != nil {
			logger.Error("Failed to end undo block: %v", endErr)
		}
		return
	}

	// End undo block
	if err = reaper.EndUndoBlock("Undo Test - Parameter Change", 0); err != nil {
		logger.Error("Failed to end undo block: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to end undo block: %v", err), "Undo Test")
		return
	}

	// Get formatted value
	formatted, _ := reaper.GetTrackFXParamFormatted(trackInfo.MediaTrack, fxIndex, paramIndex)

	// Show confirmation
	reaper.MessageBox(
		fmt.Sprintf("Parameter '%s' changed from %.2f to %.2f (%s).\n\nYou can now use undo/redo to toggle this change.",
			paramName, value, newValue, formatted),
		"Undo Test",
	)
}
