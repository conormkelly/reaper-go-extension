package actions

import (
	fxassistant "go-reaper/src/actions/fx-assistant"
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

	// Register FX Assistant Settings
	if err := fxassistant.RegisterSettingsAction(); err != nil {
		return err
	}

	// Register Native UI action
	if err := RegisterNativeWindow(); err != nil {
		return err
	}

	// Register Keyring test
	if err := RegisterKeyringTest(); err != nil {
		return err
	}

	// Register Undo Test action
	if err := RegisterUndoDemo(); err != nil {
		return err
	}

	// Register Param Format demo action
	if err := RegisterParamFormatDemo(); err != nil {
		return err
	}

	// Register Batch Parameter Demo action
	if err := RegisterBatchParamDemo(); err != nil {
		return err
	}

	// Register other actions here as they are implemented

	logger.Debug("----------------------------------------------------------")
	logger.Debug("Go plugin actions registered successfully!")
	logger.Debug("- Main section: Look for actions starting with 'Go:'")
	logger.Debug("----------------------------------------------------------")

	return nil
}
