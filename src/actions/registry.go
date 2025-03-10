package actions

import (
	"go-reaper/src/pkg/logger"
)

// RegisterAll registers all actions
func RegisterAll() error {
	logger.Debug("----------------------------------------------------------")
	logger.Debug("Registering Go REAPER extension actions...")

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

	logger.Debug("----------------------------------------------------------")
	logger.Debug("Go plugin actions registered successfully!")
	logger.Debug("- Main section: Look for actions starting with 'Go:'")
	logger.Debug("----------------------------------------------------------")

	return nil
}
