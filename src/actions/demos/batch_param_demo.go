package demo

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"strings"
)

// RegisterBatchParamDemo registers the batch parameter update demo action
func RegisterBatchParamDemo() error {
	actionID, err := reaper.RegisterMainAction("GO_BATCH_PARAM_DEMO", "Go: Batch Parameter Update Demo")
	if err != nil {
		return fmt.Errorf("failed to register batch parameter demo action: %v", err)
	}

	logger.Info("Batch Parameter Demo registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_BATCH_PARAM_DEMO", handleBatchParamDemo)
	return nil
}

// handleBatchParamDemo demonstrates the batch parameter update functionality
func handleBatchParamDemo() {
	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("Batch Parameter Demo action triggered")

	// Get track information
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		reaper.MessageBox("Please select a track first", "Batch Parameter Demo")
		return
	}

	// Get FX count
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to the track.", "Batch Parameter Demo")
		return
	}

	// Build FX list for user selection
	var fxListBuilder strings.Builder
	fxListBuilder.WriteString("Available FX:\n")
	for i := 0; i < trackInfo.NumFX; i++ {
		fxName, err := reaper.GetTrackFXName(trackInfo.MediaTrack, i)
		if err != nil {
			logger.Error("Failed to get FX name: %v", err)
			continue
		}
		fxListBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, fxName))
	}

	// Get FX selection from user
	fields := []string{
		"FX to adjust (1-" + fmt.Sprintf("%d", trackInfo.NumFX) + ")",
	}
	defaults := []string{"1"}

	results, err := reaper.GetUserInputs("Batch Parameter Demo"+fxListBuilder.String(), fields, defaults)
	if err != nil {
		logger.Info("User cancelled the dialog")
		return
	}

	// Parse FX index
	fxNumber := 1
	fmt.Sscanf(results[0], "%d", &fxNumber)
	fxIndex := fxNumber - 1 // Convert to 0-based index

	if fxIndex < 0 || fxIndex >= trackInfo.NumFX {
		reaper.MessageBox(fmt.Sprintf("Invalid FX number. Please enter a number between 1 and %d.", trackInfo.NumFX), "Batch Parameter Demo")
		return
	}

	// Get FX parameters
	fxInfo, err := reaper.GetFXParameters(trackInfo.MediaTrack, fxIndex)
	if err != nil {
		logger.Error("Failed to get parameters: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to get parameters: %v", err), "Batch Parameter Demo")
		return
	}

	if len(fxInfo.Parameters) == 0 {
		reaper.MessageBox("Selected FX has no parameters to adjust.", "Batch Parameter Demo")
		return
	}

	// Build the list of parameter changes (invert the first 3 parameters, or less if there are fewer)
	numParams := len(fxInfo.Parameters)
	if numParams > 3 {
		numParams = 3
	}

	// Create a list of parameter changes
	changes := make([]reaper.ParameterChange, numParams)
	changeDescriptions := make([]string, numParams)

	for i := 0; i < numParams; i++ {
		param := fxInfo.Parameters[i]
		// Invert the parameter value within its range
		newValue := 1.0 - param.Value

		changes[i] = reaper.ParameterChange{
			FXIndex:    fxIndex,
			ParamIndex: param.Index,
			Value:      newValue,
		}

		changeDescriptions[i] = fmt.Sprintf("• %s: %.2f → %.2f",
			param.Name, param.Value, newValue)
	}

	// Show confirmation message
	message := fmt.Sprintf("The following parameters will be adjusted:\n\n%s\n\nApply changes?",
		strings.Join(changeDescriptions, "\n"))

	confirmation, err := reaper.YesNoBox(message, "Batch Parameter Demo")
	if err != nil || !confirmation {
		logger.Info("User cancelled the changes")
		return
	}

	// Apply the changes with undo support
	err = reaper.BatchSetFXParametersWithUndo(
		trackInfo.MediaTrack,
		changes,
		"Batch Parameter Demo Changes")

	if err != nil {
		logger.Error("Failed to apply changes: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to apply changes: %v", err), "Batch Parameter Demo")
		return
	}

	reaper.MessageBox(
		fmt.Sprintf("Successfully applied %d parameter changes.\nYou can use Undo/Redo to toggle these changes.", len(changes)),
		"Batch Parameter Demo")
}
