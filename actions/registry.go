package actions

import (
	"go-reaper/core"
)

// RegisterAll registers all actions
func RegisterAll() error {
	core.LogDebug("----------------------------------------------------------")
	core.LogDebug("Registering Go REAPER extension actions...")

	// Register FX Assistant (LLM FX Assistant)
	if err := RegisterFXAssistant(); err != nil {
		return err
	}

	// Register Native UI action
	if err := RegisterNativeWindow(); err != nil {
		return err
	}

	if err := RegisterKeyringTest(); err != nil {
		return err
	}

	// Register other actions here as they are implemented

	core.LogDebug("----------------------------------------------------------")
	core.LogDebug("Go plugin actions registered successfully!")
	core.LogDebug("- Main section: Look for actions starting with 'Go:'")
	core.LogDebug("----------------------------------------------------------")

	return nil
}
