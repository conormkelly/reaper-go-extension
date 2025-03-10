package actions

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"strings"
)

// RegisterParamFormatDemo registers the parameter formatting demo action
func RegisterParamFormatDemo() error {
	actionID, err := reaper.RegisterMainAction("GO_PARAM_FORMAT_DEMO", "Go: Parameter Formatting Demo")
	if err != nil {
		return fmt.Errorf("failed to register parameter formatting demo action: %v", err)
	}

	logger.Info("Parameter Formatting Demo registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_PARAM_FORMAT_DEMO", handleParamFormatDemo)
	return nil
}

// handleParamFormatDemo demonstrates the parameter formatting functionality
func handleParamFormatDemo() {
	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("Parameter Formatting Demo action triggered")

	// Get track information
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		reaper.MessageBox("Please select a track first", "Parameter Formatting Demo")
		return
	}

	// Get FX count
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to the track.", "Parameter Formatting Demo")
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
		"FX to examine (1-" + fmt.Sprintf("%d", trackInfo.NumFX) + ")",
	}
	defaults := []string{"1"}

	results, err := reaper.GetUserInputs("Parameter Formatting Demo"+fxListBuilder.String(), fields, defaults)
	if err != nil {
		logger.Info("User cancelled the dialog")
		return
	}

	// Parse FX index
	fxNumber := 1
	fmt.Sscanf(results[0], "%d", &fxNumber)
	fxIndex := fxNumber - 1 // Convert to 0-based index

	if fxIndex < 0 || fxIndex >= trackInfo.NumFX {
		reaper.MessageBox(fmt.Sprintf("Invalid FX number. Please enter a number between 1 and %d.", trackInfo.NumFX), "Parameter Formatting Demo")
		return
	}

	// Get FX parameters with min/max formatted values
	fxInfo, err := reaper.GetFXParametersWithMinMax(trackInfo.MediaTrack, fxIndex)
	if err != nil {
		logger.Error("Failed to get parameters: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to get parameters: %v", err), "Parameter Formatting Demo")
		return
	}

	// Build parameter information for display
	var infoBuilder strings.Builder
	infoBuilder.WriteString(fmt.Sprintf("FX: %s\n\n", fxInfo.Name))
	infoBuilder.WriteString("Parameters with Min/Max Values:\n\n")

	for _, param := range fxInfo.Parameters {
		infoBuilder.WriteString(fmt.Sprintf("• %s:\n", param.Name))
		infoBuilder.WriteString(fmt.Sprintf("  Current: %.2f (%s)\n", param.Value, param.FormattedValue))
		infoBuilder.WriteString(fmt.Sprintf("  Range: %.2f to %.2f\n", param.Min, param.Max))
		infoBuilder.WriteString(fmt.Sprintf("  Formatted Range: %s to %s\n\n", param.MinFormatted, param.MaxFormatted))
	}

	// Show results
	reaper.MessageBox(infoBuilder.String(), fmt.Sprintf("Parameter Information for %s", fxInfo.Name))
}
