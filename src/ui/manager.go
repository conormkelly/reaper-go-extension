package ui

// Central UI management

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/ui/common"
	"go-reaper/src/ui/platform"
	"runtime"
	"sync"
)

var (
	// The singleton UI system
	uiSystem common.UISystem

	// Mutex for lazy initialization
	initMutex sync.Mutex

	// Track initialization status
	initialized bool
)

// Initialize the UI system for the current platform
func Initialize() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	logger.Debug("Initializing UI system for platform: %s", runtime.GOOS)

	var err error
	uiSystem, err = platform.GetUISystem()
	if err != nil {
		return fmt.Errorf("failed to initialize UI system: %v", err)
	}

	initialized = true
	logger.Info("UI system initialized successfully")
	return nil
}

// GetUISystem returns the UI system
func GetUISystem() (common.UISystem, error) {
	if !initialized {
		if err := Initialize(); err != nil {
			return nil, err
		}
	}

	return uiSystem, nil
}

// RunOnUIThread executes a function on the UI thread
func RunOnUIThread(fn func()) error {
	sys, err := GetUISystem()
	if err != nil {
		return err
	}

	return sys.RunOnMainThread(fn)
}

// IsUIThread returns true if called from the UI thread
func IsUIThread() (bool, error) {
	sys, err := GetUISystem()
	if err != nil {
		return false, err
	}

	return sys.IsMainThread(), nil
}

// CreateWindow creates a window with the specified options
func CreateWindow(options common.WindowOptions) (common.Window, error) {
	sys, err := GetUISystem()
	if err != nil {
		return nil, err
	}

	return sys.CreateWindow(options)
}

// ShowMessageBox shows a message box
func ShowMessageBox(title, message string) error {
	sys, err := GetUISystem()
	if err != nil {
		return err
	}

	return sys.ShowMessageBox(title, message)
}

// ShowConfirmDialog shows a confirmation dialog
func ShowConfirmDialog(title, message string) (bool, error) {
	sys, err := GetUISystem()
	if err != nil {
		return false, err
	}

	return sys.ShowConfirmDialog(title, message)
}

// ShowInputDialog shows an input dialog
func ShowInputDialog(title string, fields []string, defaults []string) ([]string, error) {
	sys, err := GetUISystem()
	if err != nil {
		return nil, err
	}

	return sys.ShowInputDialog(title, fields, defaults)
}

// Cleanup performs any necessary UI cleanup on shutdown
func Cleanup() {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized && uiSystem != nil {
		logger.Info("Cleaning up UI system")
		// Call any cleanup methods on the UI system

		uiSystem = nil
		initialized = false
	}
}
