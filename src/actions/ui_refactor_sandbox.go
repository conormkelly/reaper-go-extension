// src/actions/ui_refactor_sandbox.go - Sandbox for testing UI components
package actions

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"go-reaper/src/ui"
	"go-reaper/src/ui/common"
	"runtime"
	"time"
)

// RegisterUISandbox registers the UI sandbox action
func RegisterUISandbox() error {
	logger.Info("Registering UI Sandbox action")

	actionID, err := reaper.RegisterMainAction("GO_UI_SANDBOX", "Go: UI Component Sandbox")
	if err != nil {
		logger.Error("Failed to register UI sandbox action: %v", err)
		return fmt.Errorf("failed to register UI sandbox action: %v", err)
	}

	logger.Info("UI Sandbox action registered with ID: %d", actionID)

	reaper.SetActionHandler("GO_UI_SANDBOX", handleUISandbox)
	return nil
}

// handleUISandbox demonstrates the UI component system
func handleUISandbox() {
	// Lock the current goroutine to the OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("UI Sandbox action triggered")

	// Initialize UI system if needed
	if err := ui.Initialize(); err != nil {
		logger.Error("Failed to initialize UI system: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to initialize UI system: %v", err), "UI Sandbox")
		return
	}

	// Verify we're on the main thread
	isMain, _ := ui.IsUIThread()
	logger.Info("Running on main thread: %v", isMain)

	// Show a test message box through our UI abstraction
	ui.ShowMessageBox("UI Sandbox", "Welcome to the UI Component Sandbox!")

	// Create a window for parameter visualization
	window, err := ui.CreateWindow(common.WindowOptions{
		Title:     "Parameter Visualization Demo",
		X:         100,
		Y:         100,
		Width:     600,
		Height:    400,
		Resizable: true,
	})

	if err != nil {
		logger.Error("Failed to create window: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to create window: %v", err), "UI Sandbox")
		return
	}

	// Add a title label
	window.AddLabel("Parameter Visualization Demo", 20, 350, 560, 30, &common.TextOptions{
		Bold: true,
		Size: 18,
	})

	// Add an instruction label
	window.AddLabel("This is a demonstration of parameter visualization components.", 20, 320, 560, 20, nil)

	// Create some demo parameters
	params := []common.ParamState{
		{
			Name:                   "Gain",
			Value:                  0.75,
			FormattedValue:         "0.75 dB",
			OriginalValue:          0.5,
			OriginalFormattedValue: "0.5 dB",
			Min:                    0.0,
			Max:                    1.0,
			Index:                  0,
			FXIndex:                0,
			Explanation:            "Controls the output volume. Increased to make the signal louder.",
		},
		{
			Name:                   "Threshold",
			Value:                  0.3,
			FormattedValue:         "-12 dB",
			OriginalValue:          0.5,
			OriginalFormattedValue: "-6 dB",
			Min:                    0.0,
			Max:                    1.0,
			Index:                  1,
			FXIndex:                0,
			Explanation:            "Sets the level at which compression begins. Lowered to compress more of the signal.",
		},
		{
			Name:                   "Ratio",
			Value:                  0.6,
			FormattedValue:         "4:1",
			OriginalValue:          0.3,
			OriginalFormattedValue: "2:1",
			Min:                    0.0,
			Max:                    1.0,
			Index:                  2,
			FXIndex:                0,
			Explanation:            "Determines the amount of compression. Increased for stronger compression effect.",
		},
	}

	// Create parameter views
	paramViews := make([]common.ParameterView, len(params))
	for i, param := range params {
		y := 280 - (i * 80)
		paramView, err := ui.CreateParamView(window, param, 20, y, 560, 70)
		if err != nil {
			logger.Error("Failed to create parameter view for %s: %v", param.Name, err)
			continue
		}

		// Set a value change callback
		paramView.OnValueChanged(func(value float64) {
			logger.Info("Parameter %s changed to %.2f", param.Name, value)

			// Get the formatted value (this would normally come from REAPER)
			var formatted string
			switch param.Name {
			case "Gain":
				formatted = fmt.Sprintf("%.2f dB", value)
			case "Threshold":
				// Convert 0-1 to dB range (-60 to 0)
				db := -60.0 + (value * 60.0)
				formatted = fmt.Sprintf("%.1f dB", db)
			case "Ratio":
				// Convert 0-1 to ratio (1:1 to 20:1)
				ratio := 1.0 + (value * 19.0)
				formatted = fmt.Sprintf("%.1f:1", ratio)
			default:
				formatted = fmt.Sprintf("%.2f", value)
			}

			// Update the formatted value
			paramView.SetFormattedValue(formatted)
		})

		paramViews[i] = paramView
	}

	// Add a close button
	window.AddButton("Close", 250, 20, 100, 30, func() {
		logger.Info("Close button clicked")
		window.Close()
	})

	// Show the window
	if err := window.Show(); err != nil {
		logger.Error("Failed to show window: %v", err)
		reaper.MessageBox(fmt.Sprintf("Failed to show window: %v", err), "UI Sandbox")
		return
	}

	// Show parameter views
	for _, paramView := range paramViews {
		if err := paramView.Show(); err != nil {
			logger.Error("Failed to show parameter view: %v", err)
		}
	}

	// Keep the action handler alive briefly to ensure UI operations complete
	// In a real implementation, we'd have proper lifecycle management
	time.Sleep(100 * time.Millisecond)

	logger.Info("UI Sandbox action handler completed")
}

// CloseUISandboxWindow is a helper to close sandbox windows on plugin unload
func CloseUISandboxWindow() {
	// This would close any windows created by the sandbox
	// Implementation would depend on how windows are tracked
}
