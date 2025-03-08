package actions

import (
	"fmt"
	"go-reaper/reaper"
)

// RegisterFXDialogTest registers a test action for the dialog UI
func RegisterFXDialogTest() error {
	actionID, err := reaper.RegisterMainAction("GO_FX_DIALOG_TEST", "Go: Test UI Functions")
	if err != nil {
		return fmt.Errorf("failed to register UI test action: %v", err)
	}

	reaper.ConsoleLog(fmt.Sprintf("UI Test action registered with ID: %d", actionID))
	reaper.SetActionHandler("GO_FX_DIALOG_TEST", handleFXDialogTest)
	return nil
}

// handleFXDialogTest tests all UI functions
func handleFXDialogTest() {
	reaper.ConsoleLog("----- UI Function Tests -----")

	// Test basic message box
	reaper.ConsoleLog("Testing MessageBox...")
	err := reaper.MessageBox("This is a test message box.", "UI Test")
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("MessageBox error: %v", err))
	}

	// Test Yes/No dialog
	reaper.ConsoleLog("Testing YesNoBox...")
	result, err := reaper.YesNoBox("Would you like to continue testing?", "UI Test")
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("YesNoBox error: %v", err))
	} else {
		if result {
			reaper.ConsoleLog("User clicked Yes")
		} else {
			reaper.ConsoleLog("User clicked No")
			return
		}
	}

	// Test GetUserInputs
	reaper.ConsoleLog("Testing GetUserInputs...")
	fields := []string{
		"Name",
		"Email",
		"Favorite plugin",
	}

	defaults := []string{
		"User",
		"user@example.com",
		"ReaComp",
	}

	values, err := reaper.GetUserInputs("Test Form", fields, defaults)
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("GetUserInputs error: %v", err))
	} else {
		reaper.ConsoleLog("GetUserInputs results:")
		for i, value := range values {
			reaper.ConsoleLog(fmt.Sprintf("  %s: %s", fields[i], value))
		}
	}

	// Also test writing to console explicitly
	reaper.ConsoleLog("-------------------------------------")
	reaper.ConsoleLog("Console Output Test")
	reaper.ConsoleLog("-------------------------------------")

	reaper.ConsoleLog("----- UI Function Tests Complete -----")
}
