package actions

/*
#cgo darwin CFLAGS: -I${SRCDIR}/../c/platform/macos
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>
#include "../ui/platform/macos/settings_bridge.h"
*/
import "C"
import (
	"fmt"
	"go-reaper/src/pkg/config"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"unsafe"
)

// RegisterFXAssistantSettings registers the LLM FX Assistant Settings action
func RegisterFXAssistantSettings() error {
	actionID, err := reaper.RegisterMainAction("GO_FX_ASSISTANT_SETTINGS", "Go: LLM FX Assistant Settings")
	if err != nil {
		return fmt.Errorf("failed to register LLM FX Assistant Settings: %v", err)
	}

	logger.Info("LLM FX Assistant Settings registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_FX_ASSISTANT_SETTINGS", handleFXAssistantSettings)
	return nil
}

// Export the function for C to call directly
//
//export go_process_settings
func go_process_settings(apiKey *C.char, model *C.char, temperature C.double) {
	// Explicitly lock this goroutine to its OS thread for UI interactions
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Debug information
	threadID := runtime.NumGoroutine()
	isMainThread := runtime.GOMAXPROCS(0) > 1

	logger.Debug("==== go_process_settings START ====")
	logger.Debug("Thread info: goroutine=%d, isMainThread=%v", threadID, isMainThread)

	// Get values as Go strings
	goApiKey := C.GoString(apiKey)
	goModel := C.GoString(model)
	goTemperature := float64(temperature)

	// Log values (without API key for security)
	logger.Debug("Input params: model=%s, temperature=%.1f", goModel, goTemperature)
	logger.Debug("Step 1: Starting provider config retrieval")

	// Get active provider
	activeProvider := config.GetActiveProvider()
	logger.Debug("  Active provider: %s", string(activeProvider))

	// Get existing settings to preserve max tokens
	_, maxTokens, _ := config.GetProviderConfig(activeProvider)
	logger.Debug("  Existing maxTokens: %d", maxTokens)
	logger.Debug("Step 1 completed")

	// Save API key if provided
	var message string
	var success bool

	logger.Debug("Step 2: Starting API key processing")
	if goApiKey != "" {
		logger.Debug("  API key provided (not logging actual key)")
		// Save to keyring
		err := config.StoreSecureAPIKey(activeProvider, goApiKey)
		if err != nil {
			logger.Error("  Failed to save API key to keyring: %v", err)
			message = fmt.Sprintf("Error saving API key: %v", err)
			success = false
		} else {
			logger.Debug("  API key saved successfully")
			success = true
		}
	} else {
		logger.Debug("  No API key provided")
	}
	logger.Debug("Step 2 completed")

	// Use default model if empty
	logger.Debug("Step 3: Processing model")
	if goModel == "" {
		logger.Debug("  Using default model: gpt-3.5-turbo")
		goModel = "gpt-3.5-turbo"
	} else {
		logger.Debug("  Using provided model: %s", goModel)
	}
	logger.Debug("Step 3 completed")

	// Save other settings
	logger.Debug("Step 4: Saving provider config")

	// Save the configuration without checking keyring again
	// This is a key change to avoid potential UI prompts
	err := config.SetProviderConfig(activeProvider, goModel, maxTokens, goTemperature)
	if err != nil {
		logger.Error("  Failed to save provider config: %v", err)
		if success {
			message = fmt.Sprintf("API key saved but failed to save other settings: %v", err)
		} else {
			message = fmt.Sprintf("Failed to save settings: %v", err)
		}
		success = false
	} else {
		logger.Debug("  Provider config saved successfully")
		if success {
			message = fmt.Sprintf("Settings saved successfully!\n\nModel: %s\nTemperature: %.1f",
				goModel, goTemperature)
		} else {
			// We know the API key status from the earlier check
			// No need to check HasSecureAPIKey here
			if goApiKey != "" {
				// We tried to save an API key
				message = fmt.Sprintf("Settings saved successfully!\n\nModel: %s\nTemperature: %.1f",
					goModel, goTemperature)
				success = true
			} else {
				message = "Settings saved but no API key provided.\nYou'll need an API key to use the LLM FX Assistant."
				success = false
			}
		}
	}
	logger.Debug("Step 4 completed")

	// Show confirmation or error message
	logger.Debug("Step 5: Logging messages")
	reaper.ShowConsoleMsg(fmt.Sprintf("LLM FX Assistant Settings: %s\n", message))

	if success {
		logger.Info("Settings saved successfully: model=%s, temperature=%.1f", goModel, goTemperature)
	} else {
		logger.Warning("Settings not fully saved: %s", message)
	}
	logger.Debug("Step 5 completed")

	// Use console message instead of MessageBox to avoid UI blocking
	logger.Debug("Step 6: Showing result message")

	// Use console message for now to avoid potential UI issues
	if success {
		reaper.ShowConsoleMsg("SUCCESS: " + message + "\n")
	} else {
		reaper.ShowConsoleMsg("WARNING: " + message + "\n")
	}

	// Signal to Objective-C that we're done
	logger.Debug("Step 6 completed")
	logger.Debug("==== go_process_settings END ====")
}

// handleFXAssistantSettings handles the settings management for the LLM FX Assistant
func handleFXAssistantSettings() {
	// Only macOS is supported for now
	if runtime.GOOS != "darwin" {
		// Fallback to basic message on non-macOS platforms
		handleFXAssistantSettingsFallback()
		return
	}

	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Debug("----- LLM FX Assistant Settings Activated -----")

	// Get current settings
	activeProvider := config.GetActiveProvider()
	model, _, temperature := config.GetProviderConfig(activeProvider)

	// Try to get existing API key
	apiKey, err := config.GetSecureAPIKey(activeProvider)
	if err != nil {
		apiKey = "" // No key or error getting it
		logger.Debug("No existing API key found or error: %v", err)
	}

	// Convert to C strings
	cTitle := C.CString("REAPER LLM FX Assistant Settings")
	defer C.free(unsafe.Pointer(cTitle))

	cApiKey := C.CString(apiKey)
	defer C.free(unsafe.Pointer(cApiKey))

	cModel := C.CString(model)
	defer C.free(unsafe.Pointer(cModel))

	// Show the settings window
	result := C.settings_show_window(cTitle, cApiKey, cModel, C.double(temperature))

	if bool(result) {
		logger.Info("Settings window created/shown successfully")
	} else {
		logger.Error("Failed to create/show settings window")
		reaper.MessageBox("Failed to create/show settings window. See log for details.", "LLM FX Assistant Settings")
	}

	logger.Info("LLM FX Assistant Settings action handler completed")
}

// CloseSettingsWindow closes the settings window if it exists
func CloseSettingsWindow() {
	if runtime.GOOS == "darwin" {
		logger.Info("Closing settings window...")
		C.settings_close_window()
		logger.Info("Settings window close request completed")
	}
}

// IsSettingsWindowOpen checks if the settings window is currently open
func IsSettingsWindowOpen() bool {
	if runtime.GOOS == "darwin" {
		return bool(C.settings_window_exists())
	}
	return false
}

// handleFXAssistantSettingsFallback provides an error message for non-macOS platforms
func handleFXAssistantSettingsFallback() {
	logger.Debug("----- LLM FX Assistant Settings Fallback Activated -----")

	// Show message that native UI is only available on macOS for now
	reaper.MessageBox(
		"Native settings UI is currently only available on macOS.\n\n"+
			"Support for Windows and Linux will be added in a future update.",
		"LLM FX Assistant Settings")

	logger.Info("Informed user that native settings UI is macOS-only for now")
}
