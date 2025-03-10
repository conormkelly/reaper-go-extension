// src/ui/platform/macos.go - macOS UI implementation
package platform

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/ui/common"
	"runtime"
	"unsafe"
)

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include "../../c/logging.h"

#import <Cocoa/Cocoa.h>

// Forward declarations
bool macos_is_main_thread(void);
bool macos_run_on_main_thread(void (*func)(void*), void* context);
void* macos_create_window(const char* title, int x, int y, int width, int height, bool resizable);
bool macos_close_window(void* window);
bool macos_show_window(void* window);
bool macos_hide_window(void* window);
bool macos_window_is_visible(void* window);
bool macos_set_window_title(void* window, const char* title);
void* macos_add_label(void* window, const char* text, int x, int y, int width, int height, bool bold, double size);
void* macos_add_button(void* window, const char* text, int x, int y, int width, int height);
void* macos_add_slider(void* window, double min, double max, double value, int x, int y, int width, int height);
void* macos_add_text_field(void* window, const char* placeholder, int x, int y, int width, int height);
int macos_show_alert(const char* title, const char* message, int style);
bool macos_get_user_inputs(const char* title, int num_inputs, const char* captions, char* values, int values_sz);

// Callback typedefs
typedef void (*ButtonCallback)(void* sender);
typedef void (*SliderCallback)(void* sender, double value);

// Register callbacks
bool macos_set_button_callback(void* button, ButtonCallback callback);
bool macos_set_slider_callback(void* slider, SliderCallback callback);

// Parameter view functions
void* macos_create_param_view(void* window, int x, int y, int width, int height, const char* name,
                              double value, double min, double max, const char* formatted);
bool macos_param_view_set_value(void* view, double value);
bool macos_param_view_set_formatted(void* view, const char* formatted);
bool macos_param_view_set_explanation(void* view, const char* explanation);
bool macos_param_view_set_original(void* view, double original, const char* formatted);
bool macos_param_view_show(void* view);
bool macos_param_view_hide(void* view);
double macos_param_view_get_value(void* view);
bool macos_param_view_set_callback(void* view, SliderCallback callback);

// Actual implementations of these functions would be in platform/macos/ui.m
// For now we'll use stub implementations for the interface
*/
import "C"

// Ensure macOS implementation is only used on macOS
func init() {
	if runtime.GOOS != "darwin" {
		logger.Warning("macOS UI implementation loaded on non-macOS platform: %s", runtime.GOOS)
	}
}

// macOSUISystem implements the common.UISystem interface for macOS
type macOSUISystem struct {
	paramViewFactory common.ParamViewFactory
}

// GetUISystem returns the platform-specific UI system
func GetUISystem() (common.UISystem, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("macOS UI implementation not available on %s", runtime.GOOS)
	}

	system := &macOSUISystem{
		paramViewFactory: &macOSParamViewFactory{},
	}

	return system, nil
}

// GetParamViewFactory returns a factory for parameter views
func (s *macOSUISystem) GetParamViewFactory() common.ParamViewFactory {
	return s.paramViewFactory
}

// RunOnMainThread runs the given function on the main thread
func (s *macOSUISystem) RunOnMainThread(fn func()) error {
	// If already on main thread, run directly
	if s.IsMainThread() {
		fn()
		return nil
	}

	// Use a channel to ensure completion
	done := make(chan struct{})

	// Create a closure to execute on main thread and signal completion
	execOnMain := func() {
		defer close(done)
		fn()
	}

	// Execute on main thread via dispatch_async
	// Note: We're using a simplified approach here.
	// In a complete implementation, we would need to create a proper
	// bridge to pass the Go function to Objective-C code.
	dispatch_async(execOnMain)

	// Wait for completion
	<-done
	return nil
}

// dispatch_async is a simplified wrapper around macOS dispatch_async
// In a real implementation, this would use proper CGO bindings
func dispatch_async(fn func()) {
	// Call the C function that dispatches to main thread
	ok := C.macos_run_on_main_thread(nil, nil)
	if !bool(ok) {
		logger.Error("Failed to dispatch to main thread")
		return
	}

	// In a real implementation, the C code would call back to Go
	// and execute our function. For now, we'll just call it directly
	// as a simplification.
	fn()
}

// IsMainThread returns true if called from the main thread
func (s *macOSUISystem) IsMainThread() bool {
	return bool(C.macos_is_main_thread())
}

// ShowMessageBox shows a message box
func (s *macOSUISystem) ShowMessageBox(title, message string) error {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	C.macos_show_alert(cTitle, cMessage, 0) // 0 = OK style
	return nil
}

// ShowConfirmDialog shows a Yes/No dialog
func (s *macOSUISystem) ShowConfirmDialog(title, message string) (bool, error) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))

	result := C.macos_show_alert(cTitle, cMessage, 1) // 1 = Yes/No style
	return result == 1, nil                           // 1 = Yes, 0 = No
}

// ShowInputDialog shows a dialog with input fields
func (s *macOSUISystem) ShowInputDialog(title string, fields []string, defaults []string) ([]string, error) {
	// This is a simplified implementation that would need to be expanded
	return nil, fmt.Errorf("not implemented")
}

// macOSParamViewFactory implements the ParamViewFactory interface for macOS
type macOSParamViewFactory struct{}

// CreateWindow creates a macOS window
func (f *macOSParamViewFactory) CreateWindow(options common.WindowOptions) (common.Window, error) {
	return &macOSWindow{
		options: options,
	}, nil
}

// CreateParamView creates a parameter view
func (f *macOSParamViewFactory) CreateParamView(window common.Window, param common.ParamState, x, y, width, height int) (common.ParameterView, error) {
	macWindow, ok := window.(*macOSWindow)
	if !ok {
		return nil, fmt.Errorf("window is not a macOS window")
	}

	return &macOSParamView{
		window: macWindow,
		param:  param,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}, nil
}

// macOSWindow implements the Window interface for macOS
type macOSWindow struct {
	options common.WindowOptions
	handle  unsafe.Pointer
	visible bool
}

// Show the window
func (w *macOSWindow) Show() error {
	if w.handle == nil {
		// Create the window if it doesn't exist
		cTitle := C.CString(w.options.Title)
		defer C.free(unsafe.Pointer(cTitle))

		w.handle = C.macos_create_window(cTitle,
			C.int(w.options.X),
			C.int(w.options.Y),
			C.int(w.options.Width),
			C.int(w.options.Height),
			C.bool(w.options.Resizable))

		if w.handle == nil {
			return fmt.Errorf("failed to create window")
		}
	}

	if !C.macos_show_window(w.handle) {
		return fmt.Errorf("failed to show window")
	}

	w.visible = true
	return nil
}

// Hide the window
func (w *macOSWindow) Hide() error {
	if w.handle == nil {
		return nil
	}

	if !C.macos_hide_window(w.handle) {
		return fmt.Errorf("failed to hide window")
	}

	w.visible = false
	return nil
}

// Close and dispose of window resources
func (w *macOSWindow) Close() error {
	if w.handle == nil {
		return nil
	}

	if !C.macos_close_window(w.handle) {
		return fmt.Errorf("failed to close window")
	}

	w.handle = nil
	w.visible = false
	return nil
}

// IsVisible returns true if window is visible
func (w *macOSWindow) IsVisible() bool {
	if w.handle == nil {
		return false
	}

	return bool(C.macos_window_is_visible(w.handle))
}

// AddLabel adds a text label
func (w *macOSWindow) AddLabel(text string, x, y, width, height int, options *common.TextOptions) error {
	if w.handle == nil {
		return fmt.Errorf("window not created")
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	bold := false
	size := 12.0
	if options != nil {
		bold = options.Bold
		size = options.Size
	}

	label := C.macos_add_label(w.handle, cText, C.int(x), C.int(y), C.int(width), C.int(height),
		C.bool(bold), C.double(size))

	if label == nil {
		return fmt.Errorf("failed to add label")
	}

	return nil
}

// AddButton adds a button
func (w *macOSWindow) AddButton(text string, x, y, width, height int, callback common.ActionCallback) error {
	if w.handle == nil {
		return fmt.Errorf("window not created")
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	button := C.macos_add_button(w.handle, cText, C.int(x), C.int(y), C.int(width), C.int(height))
	if button == nil {
		return fmt.Errorf("failed to add button")
	}

	// Register callback
	// This is simplified and would need actual implementation

	return nil
}

// AddSlider adds a horizontal slider
func (w *macOSWindow) AddSlider(min, max, value float64, x, y, width, height int, callback common.ValueChangeCallback) error {
	if w.handle == nil {
		return fmt.Errorf("window not created")
	}

	slider := C.macos_add_slider(w.handle, C.double(min), C.double(max), C.double(value),
		C.int(x), C.int(y), C.int(width), C.int(height))

	if slider == nil {
		return fmt.Errorf("failed to add slider")
	}

	// Register callback
	// This is simplified and would need actual implementation

	return nil
}

// AddTextField adds a text field
func (w *macOSWindow) AddTextField(placeholder string, x, y, width, height int) error {
	if w.handle == nil {
		return fmt.Errorf("window not created")
	}

	cPlaceholder := C.CString(placeholder)
	defer C.free(unsafe.Pointer(cPlaceholder))

	textField := C.macos_add_text_field(w.handle, cPlaceholder, C.int(x), C.int(y), C.int(width), C.int(height))
	if textField == nil {
		return fmt.Errorf("failed to add text field")
	}

	return nil
}

// SetTitle changes the window title
func (w *macOSWindow) SetTitle(title string) error {
	if w.handle == nil {
		w.options.Title = title
		return nil
	}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	if !C.macos_set_window_title(w.handle, cTitle) {
		return fmt.Errorf("failed to set window title")
	}

	w.options.Title = title
	return nil
}

// macOSParamView implements the ParameterView interface for macOS
type macOSParamView struct {
	window        *macOSWindow
	param         common.ParamState
	handle        unsafe.Pointer
	x, y          int
	width, height int
	callback      common.ValueChangeCallback
}

// SetValue updates the displayed value
func (p *macOSParamView) SetValue(value float64) error {
	if p.handle == nil {
		p.param.Value = value
		return nil
	}

	if !C.macos_param_view_set_value(p.handle, C.double(value)) {
		return fmt.Errorf("failed to set parameter value")
	}

	p.param.Value = value
	return nil
}

// GetValue returns the current value
func (p *macOSParamView) GetValue() float64 {
	if p.handle == nil {
		return p.param.Value
	}

	return float64(C.macos_param_view_get_value(p.handle))
}

// SetFormattedValue updates the displayed formatted value
func (p *macOSParamView) SetFormattedValue(formatted string) error {
	if p.handle == nil {
		p.param.FormattedValue = formatted
		return nil
	}

	cFormatted := C.CString(formatted)
	defer C.free(unsafe.Pointer(cFormatted))

	if !C.macos_param_view_set_formatted(p.handle, cFormatted) {
		return fmt.Errorf("failed to set formatted value")
	}

	p.param.FormattedValue = formatted
	return nil
}

// SetExplanation updates the explanation text
func (p *macOSParamView) SetExplanation(text string) error {
	if p.handle == nil {
		p.param.Explanation = text
		return nil
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	if !C.macos_param_view_set_explanation(p.handle, cText) {
		return fmt.Errorf("failed to set explanation")
	}

	p.param.Explanation = text
	return nil
}

// SetOriginalValue sets the original value (for comparison)
func (p *macOSParamView) SetOriginalValue(value float64, formatted string) error {
	if p.handle == nil {
		p.param.OriginalValue = value
		p.param.OriginalFormattedValue = formatted
		return nil
	}

	cFormatted := C.CString(formatted)
	defer C.free(unsafe.Pointer(cFormatted))

	if !C.macos_param_view_set_original(p.handle, C.double(value), cFormatted) {
		return fmt.Errorf("failed to set original value")
	}

	p.param.OriginalValue = value
	p.param.OriginalFormattedValue = formatted
	return nil
}

// OnValueChanged sets the callback for value changes
func (p *macOSParamView) OnValueChanged(callback common.ValueChangeCallback) error {
	p.callback = callback

	if p.handle == nil {
		return nil
	}

	// Register callback
	// This is simplified and would need actual implementation

	return nil
}

// Show the parameter view
func (p *macOSParamView) Show() error {
	if p.handle == nil {
		// Create the parameter view
		if p.window.handle == nil {
			if err := p.window.Show(); err != nil {
				return err
			}
		}

		cName := C.CString(p.param.Name)
		defer C.free(unsafe.Pointer(cName))

		cFormatted := C.CString(p.param.FormattedValue)
		defer C.free(unsafe.Pointer(cFormatted))

		p.handle = C.macos_create_param_view(p.window.handle,
			C.int(p.x), C.int(p.y), C.int(p.width), C.int(p.height),
			cName, C.double(p.param.Value), C.double(p.param.Min), C.double(p.param.Max),
			cFormatted)

		if p.handle == nil {
			return fmt.Errorf("failed to create parameter view")
		}

		// Set explanation if available
		if p.param.Explanation != "" {
			cExplanation := C.CString(p.param.Explanation)
			defer C.free(unsafe.Pointer(cExplanation))

			C.macos_param_view_set_explanation(p.handle, cExplanation)
		}

		// Set original value if available
		if p.param.OriginalValue != 0 || p.param.OriginalFormattedValue != "" {
			cOrigFormatted := C.CString(p.param.OriginalFormattedValue)
			defer C.free(unsafe.Pointer(cOrigFormatted))

			C.macos_param_view_set_original(p.handle, C.double(p.param.OriginalValue), cOrigFormatted)
		}

		// Register callback if available
		if p.callback != nil {
			// Register callback
			// This is simplified and would need actual implementation
		}
	}

	if !C.macos_param_view_show(p.handle) {
		return fmt.Errorf("failed to show parameter view")
	}

	return nil
}

// Hide the parameter view
func (p *macOSParamView) Hide() error {
	if p.handle == nil {
		return nil
	}

	if !C.macos_param_view_hide(p.handle) {
		return fmt.Errorf("failed to hide parameter view")
	}

	return nil
}
