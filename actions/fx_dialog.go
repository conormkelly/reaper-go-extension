package actions

import (
	"fmt"
	"go-reaper/reaper"
	"strconv"
	"strings"
	"unsafe"
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
		// Handle the case where there's no track selected
		if strings.Contains(err.Error(), "no track selected") {
			reaper.MessageBox("Please select a track before using the LLM FX Assistant.", "LLM FX Assistant")
			reaper.ConsoleLog("No track selected. Please select a track before using the LLM FX Assistant.")
		} else {
			// Handle other errors
			reaper.MessageBox(fmt.Sprintf("Error: %v", err), "LLM FX Assistant")
			reaper.ConsoleLog(fmt.Sprintf("Error getting track info: %v", err))
		}
		return
	}

	// Check if track has FX
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to the track before using the LLM FX Assistant.", "LLM FX Assistant")
		reaper.ConsoleLog("Selected track has no FX. Please add FX to the track before using the LLM FX Assistant.")
		return
	}

	// Get FX list
	fxList, err := reaper.GetTrackFXList(trackInfo.MediaTrack)
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Error getting FX list: %v", err))
		return
	}

	reaper.ConsoleLog(fmt.Sprintf("Found %d FX on track.", len(fxList)))

	// Build FX selection dialog
	fxOptions := buildFXSelectionList(fxList)
	reaper.ConsoleLog(fxOptions)

	// Show FX selection dialog
	fields := []string{
		"FX to adjust (comma-separated numbers)",
		"Your request (e.g., 'make vocals clearer')",
	}

	defaults := []string{
		"1", // Default to first FX
		"",  // Empty prompt
	}

	// Show the dialog
	results, err := reaper.GetUserInputs("LLM FX Assistant", fields, defaults)
	if err != nil {
		reaper.ConsoleLog("User cancelled the dialog")
		return
	}

	// Parse the results
	selectedFXIndices, err := parseFXSelection(results[0], len(fxList))
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Invalid FX selection: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Invalid FX selection: %v", err))
		return
	}

	userPrompt := results[1]
	if userPrompt == "" {
		reaper.MessageBox("Please provide a request for the LLM FX Assistant.", "LLM FX Assistant")
		reaper.ConsoleLog("Empty prompt provided")
		return
	}

	// Log the selections
	reaper.ConsoleLog(fmt.Sprintf("Selected FX indices: %v", selectedFXIndices))
	reaper.ConsoleLog(fmt.Sprintf("User prompt: %s", userPrompt))

	// Collect FX parameters for the selected FX
	fxParameters := collectFXParameters(trackInfo.MediaTrack, selectedFXIndices, fxList)

	// Format the FX parameters for display
	parametersText := formatFXParametersText(fxParameters)

	// For Phase 1, display the collected parameters
	displayText := fmt.Sprintf("Track: %s\nPrompt: %s\n\n%s",
		trackInfo.Name,
		userPrompt,
		parametersText)

	res, err := reaper.ShowMessageBox(displayText, "LLM FX Assistant - Parameters", 0)
	reaper.ConsoleLog(fmt.Sprintf("Got %d res from ShowMessageBox", res))

	reaper.ConsoleLog("Phase 1 complete. Parameters collected successfully.")
	reaper.ConsoleLog("In Phase 2, we will integrate with an LLM API to get suggestions.")
}

// buildFXSelectionList creates a formatted list of FX for the console log
func buildFXSelectionList(fxList []reaper.FXInfo) string {
	var builder strings.Builder
	builder.WriteString("\nAvailable FX:\n")

	for i, fx := range fxList {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, fx.Name))
	}

	return builder.String()
}

// parseFXSelection parses a comma-separated list of FX indices
func parseFXSelection(input string, maxFX int) ([]int, error) {
	if input == "" {
		return nil, fmt.Errorf("no FX selected")
	}

	// Split by comma
	parts := strings.Split(input, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Parse the number
		idx, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid FX number: %s", part)
		}

		// Adjust for 1-based indexing in the UI to 0-based indexing internally
		idx--

		// Check range
		if idx < 0 || idx >= maxFX {
			return nil, fmt.Errorf("FX number out of range: %d", idx+1)
		}

		result = append(result, idx)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid FX selected")
	}

	return result, nil
}

// collectFXParameters collects all parameters for the selected FX
func collectFXParameters(track unsafe.Pointer, indices []int, fxList []reaper.FXInfo) []reaper.FXInfo {
	result := make([]reaper.FXInfo, 0, len(indices))

	for _, fxIndex := range indices {
		// Get full FX parameters
		fxInfo, err := reaper.GetFXParameters(track, fxIndex)
		if err != nil {
			reaper.ConsoleLog(fmt.Sprintf("Error getting FX parameters for %s: %v",
				fxList[fxIndex].Name, err))
			continue
		}

		result = append(result, fxInfo)
	}

	return result
}

// formatFXParametersText formats the FX parameters for display
func formatFXParametersText(fxParameters []reaper.FXInfo) string {
	var builder strings.Builder

	for _, fx := range fxParameters {
		builder.WriteString(fmt.Sprintf("FX: %s\n", fx.Name))
		builder.WriteString("Parameters:\n")

		for _, param := range fx.Parameters {
			builder.WriteString(fmt.Sprintf("  %s: %.2f (%s)\n",
				param.Name, param.Value, param.FormattedValue))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
