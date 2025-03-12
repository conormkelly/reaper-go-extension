package demo

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"unsafe"

	"github.com/zalando/go-keyring"
)

// This file implements a keyring test with native macOS UI

/*
#cgo darwin CFLAGS: -I${SRCDIR}/../../c/platform/macos
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdlib.h>
#include "../../ui/platform/macos/keyring_bridge.h"
*/
import "C"

// Constants for keyring access
const (
	KeyringServiceName = "GoReaperExtension"
	KeyringKeyName     = "APIKey"
)

// Export the function for C to call directly
//
//export go_process_keyring_key
func go_process_keyring_key(keyValue *C.char) {
	// Get key as Go string
	key := C.GoString(keyValue)

	// Log the key length (for debugging)
	logger.Debug("Processing key for keyring (length: %d)", len(key))

	var message string
	var success bool

	// Save to keyring
	err := keyring.Set(KeyringServiceName, KeyringKeyName, key)
	if err != nil {
		logger.Error("Failed to save key to keyring: %v", err)
		message = fmt.Sprintf("Error saving to keyring: %v", err)
		success = false
	} else {
		message = "Success! You've added the key to the keyring!"
		logger.Info("Key saved to keyring successfully")
		success = true
	}

	// Update UI
	updateMessage(success, message)
}

// RegisterKeyringTest registers the keyring test action
func RegisterKeyringTest() error {
	logger.Info("Registering Keyring Test action")

	actionID, err := reaper.RegisterMainAction("GO_KEYRING_TEST", "Go: Keyring Test")
	if err != nil {
		logger.Error("Failed to register keyring test action: %v", err)
		return fmt.Errorf("failed to register keyring test action: %v", err)
	}

	logger.Info("Keyring Test action registered with ID: %d", actionID)

	reaper.SetActionHandler("GO_KEYRING_TEST", handleKeyringTest)
	return nil
}

// handleKeyringTest executes the keyring test action
func handleKeyringTest() {
	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		reaper.MessageBox("This keyring test is currently only implemented for macOS", "Keyring Test")
		return
	}

	// Log action triggered
	logger.Info("Keyring Test action triggered")

	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Check if the key exists in keyring
	key, err := keyring.Get(KeyringServiceName, KeyringKeyName)
	keyExists := (err == nil && key != "")

	var message string
	if keyExists {
		message = "Success! The key is in the keyring."
		logger.Info("API key found in keyring")
	} else {
		message = "No key found. Please enter your API key."
		logger.Info("No API key found in keyring")
	}

	// Show the keyring window
	title := C.CString("REAPER Keyring Test")
	defer C.free(unsafe.Pointer(title))

	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	result := C.kr_show_window(title, C.bool(keyExists), cMessage)

	if bool(result) {
		logger.Info("Keyring window created/shown successfully")
	} else {
		logger.Error("Failed to create/show keyring window")
		reaper.MessageBox("Failed to create/show keyring window. See log for details.", "Keyring Test")
	}

	logger.Info("Keyring Test action handler completed")
}

// updateMessage updates the message in the keyring window
func updateMessage(keyExists bool, message string) {
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	result := C.kr_update_message(C.bool(keyExists), cMessage)

	if bool(result) {
		logger.Info("Keyring message updated successfully")
	} else {
		logger.Error("Failed to update keyring message")
	}
}

// CloseKeyringWindow closes the keyring window if it exists
func CloseKeyringWindow() {
	if runtime.GOOS == "darwin" {
		logger.Info("Closing keyring window...")

		C.kr_close_window()

		logger.Info("Keyring window close request completed")
	}
}

// IsKeyringWindowOpen checks if the keyring window is currently open
func IsKeyringWindowOpen() bool {
	if runtime.GOOS == "darwin" {
		return bool(C.kr_window_exists())
	}
	return false
}
