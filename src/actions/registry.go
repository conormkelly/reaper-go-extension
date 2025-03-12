package actions

import (
	demo "go-reaper/src/actions/demos"
	fxassistant "go-reaper/src/actions/fx-assistant"
	"go-reaper/src/pkg/logger"
)

// RegisterAll registers all actions
func RegisterAll() error {
	logger.Debug("----------------------------------------------------------")
	logger.Debug("Registering Go REAPER extension actions...")

	// Register FX Assistant (LLM FX Assistant)
	if err := demo.RegisterFXAssistant(); err != nil {
		return err
	}

	// Register FX Assistant Settings
	if err := fxassistant.RegisterSettingsAction(); err != nil {
		return err
	}

	// Register Native UI action
	if err := demo.RegisterNativeWindow(); err != nil {
		return err
	}

	// Register Keyring test
	if err := demo.RegisterKeyringTest(); err != nil {
		return err
	}

	// Register Undo Test action
	if err := demo.RegisterUndoDemo(); err != nil {
		return err
	}

	// Register Param Format demo action
	if err := demo.RegisterParamFormatDemo(); err != nil {
		return err
	}

	// Register Batch Parameter Demo action
	if err := demo.RegisterBatchParamDemo(); err != nil {
		return err
	}

	// Register other actions here as they are implemented

	logger.Debug("----------------------------------------------------------")
	logger.Debug("Go plugin actions registered successfully!")
	logger.Debug("- Main section: Look for actions starting with 'Go:'")
	logger.Debug("----------------------------------------------------------")

	return nil
}
