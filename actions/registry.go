package actions

import "go-reaper/reaper"

// RegisterAll registers all actions
func RegisterAll() error {
	reaper.ConsoleLog("----------------------------------------------------------")
	reaper.ConsoleLog("Registering Go REAPER extension actions...")

	// Register FX Assistant (LLM FX Assistant)
	if err := RegisterFXAssistant(); err != nil {
		return err
	}

	// Register other actions here as they are implemented

	reaper.ConsoleLog("----------------------------------------------------------")
	reaper.ConsoleLog("Go plugin actions registered successfully!")
	reaper.ConsoleLog("- Main section: Look for actions starting with 'Go:'")
	reaper.ConsoleLog("----------------------------------------------------------")

	return nil
}
