package fxassistant

import (
	"fmt"
	"go-reaper/src/pkg/config"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"go-reaper/src/reaper/fx"
	"runtime"
)

// RegisterFXAssistantAction registers the LLM FX Assistant action
func RegisterFXAssistantAction() error {
	actionID, err := reaper.RegisterMainAction("GO_FX_ASSISTANT", "Go: LLM FX Assistant")
	if err != nil {
		return fmt.Errorf("failed to register LLM FX Assistant: %v", err)
	}

	logger.Info("LLM FX Assistant registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_FX_ASSISTANT", handleFXAssistant)
	return nil
}

// handleFXAssistant handles the FX Assistant action
func handleFXAssistant() {
	// Lock the current goroutine to the OS thread to ensure thread safety
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Debug("----- LLM FX Assistant Activated -----")

	// STEP 1: Check for first-time setup
	if !isSetupComplete() {
		showFirstTimeSetupMessage()
		return
	}

	// STEP 2: Present FX selection UI and get user selection
	selectionResult, err := SelectFXForProcessing()
	if err != nil {
		handleError("FX Selection Error", err)
		return
	}

	// STEP 3: Show prompt entry dialog
	prompt, err := promptForUserRequest(selectionResult.Description)
	if err != nil {
		handleError("Prompt Entry Error", err)
		return
	}

	// STEP 4: Process the request with the LLM
	modifications, explanation, err := processRequestWithLLM(selectionResult.Collection, selectionResult.SelectedFX, prompt)
	if err != nil {
		handleError("LLM Processing Error", err)
		return
	}

	// STEP 5: Apply the changes if there are any
	if len(modifications) > 0 {
		// Present the changes to the user and ask for confirmation
		confirmed, err := presentChangesForConfirmation(modifications, explanation)
		if err != nil {
			handleError("Confirmation Error", err)
			return
		}

		if confirmed {
			// Apply the changes with undo support
			err = applyParameterModifications(selectionResult.Collection, modifications)
			if err != nil {
				handleError("Parameter Application Error", err)
				return
			}

			// Show success message
			reaper.MessageBox(
				fmt.Sprintf("Successfully applied %d parameter changes.\n\nYou can use Undo to revert these changes if needed.", len(modifications)),
				"LLM FX Assistant")
		} else {
			logger.Info("User chose not to apply the suggested changes")
		}
	} else {
		// No changes suggested - show the explanation
		reaper.MessageBox(
			fmt.Sprintf("No parameter changes were suggested.\n\n%s", explanation),
			"LLM FX Assistant")
	}
}

// isSetupComplete checks if the required setup (API key, etc.) is complete
func isSetupComplete() bool {
	// Get active provider
	provider := config.GetActiveProvider()

	// Check if API key exists
	return config.HasSecureAPIKey(provider)
}

// showFirstTimeSetupMessage shows a message for first-time users
func showFirstTimeSetupMessage() {
	message := `This appears to be your first time using the LLM FX Assistant.

Before you can use this feature, you need to configure your API key.

Please run the "Go: LLM FX Assistant Settings" action to configure your settings.`

	reaper.MessageBox(message, "LLM FX Assistant Setup Required")
}

// promptForUserRequest shows a dialog to get the user's request
func promptForUserRequest(selectionDescription string) (string, error) {
	// Show a dialog with the selection description and prompt for the request
	message := fmt.Sprintf("Selected FX:\n%s\n\nDescribe what you want to do with these FX:", selectionDescription)

	fields := []string{message}
	defaults := []string{"Make the sound warmer"}

	results, err := reaper.GetUserInputs("LLM FX Assistant - Request", fields, defaults)
	if err != nil {
		return "", fmt.Errorf("user cancelled the request dialog")
	}

	userPrompt := results[0]
	if userPrompt == "" {
		return "", fmt.Errorf("empty prompt provided")
	}

	return userPrompt, nil
}

// handleError shows an error message to the user
func handleError(title string, err error) {
	logger.Error("%s: %v", title, err)
	reaper.MessageBox(fmt.Sprintf("Error: %v", err), title)
}

// presentChangesForConfirmation shows the suggested changes to the user and asks for confirmation
func presentChangesForConfirmation(modifications []fx.ParameterModification, explanation string) (bool, error) {
	// Create a user-friendly display of the changes
	message := fmt.Sprintf("The LLM suggests the following changes:\n\n%s\n\n", explanation)

	// Group changes by track and FX
	trackFXChanges := make(map[int]map[int][]fx.ParameterModification)
	for _, mod := range modifications {
		if trackFXChanges[mod.TrackIndex] == nil {
			trackFXChanges[mod.TrackIndex] = make(map[int][]fx.ParameterModification)
		}
		if trackFXChanges[mod.TrackIndex][mod.FXIndex] == nil {
			trackFXChanges[mod.TrackIndex][mod.FXIndex] = []fx.ParameterModification{}
		}
		trackFXChanges[mod.TrackIndex][mod.FXIndex] = append(trackFXChanges[mod.TrackIndex][mod.FXIndex], mod)
	}

	// Format the changes by track and FX
	for trackIdx, fxMap := range trackFXChanges {
		message += fmt.Sprintf("Track %d:\n", trackIdx+1)
		for fxIdx, mods := range fxMap {
			message += fmt.Sprintf("  FX %d (%s):\n", fxIdx+1, mods[0].FXName)
			for _, mod := range mods {
				message += fmt.Sprintf("    • %s: %s → %s\n      %s\n",
					mod.ParamName,
					mod.OriginalFormatted,
					mod.NewFormatted,
					mod.Explanation)
			}
		}
	}

	message += "\nApply these changes?"

	// Ask for confirmation
	confirmed, err := reaper.YesNoBox(message, "LLM FX Assistant - Confirm Changes")
	if err != nil {
		return false, fmt.Errorf("error showing confirmation dialog: %v", err)
	}

	return confirmed, nil
}

// applyParameterModifications applies the suggested parameter modifications
func applyParameterModifications(collection fx.TrackCollection, modifications []fx.ParameterModification) error {
	// Get track pointers from the collection
	tracks := fx.GetTrackMediaTracks(collection)

	// Apply changes using the fx package's helper function
	return fx.ApplyParameterModifications(tracks, modifications, "LLM FX Assistant Changes")
}
