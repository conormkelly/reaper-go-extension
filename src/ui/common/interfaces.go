package common

// UI component interfaces

// Window represents a UI window
type Window interface {
	// Show the window
	Show() error

	// Hide the window
	Hide() error

	// Close and dispose of window resources
	Close() error

	// IsVisible returns true if window is visible
	IsVisible() bool

	// AddLabel adds a text label
	AddLabel(text string, x, y, width, height int, options *TextOptions) error

	// AddButton adds a button
	AddButton(text string, x, y, width, height int, callback ActionCallback) error

	// AddSlider adds a horizontal slider
	AddSlider(min, max, value float64, x, y, width, height int, callback ValueChangeCallback) error

	// AddTextField adds a text field
	AddTextField(placeholder string, x, y, width, height int) error

	// SetTitle changes the window title
	SetTitle(title string) error
}

// ParameterView represents a UI component that displays a parameter with slider
type ParameterView interface {
	// SetValue updates the displayed value
	SetValue(value float64) error

	// GetValue returns the current value
	GetValue() float64

	// SetFormattedValue updates the displayed formatted value
	SetFormattedValue(formatted string) error

	// SetExplanation updates the explanation text
	SetExplanation(text string) error

	// SetOriginalValue sets the original value (for comparison)
	SetOriginalValue(value float64, formatted string) error

	// OnValueChanged sets the callback for value changes
	OnValueChanged(callback ValueChangeCallback) error

	// Show the parameter view
	Show() error

	// Hide the parameter view
	Hide() error
}

// ParamViewFactory creates parameter visualization components
type ParamViewFactory interface {
	// CreateParamView creates a parameter view inside the given window
	CreateParamView(window Window, param ParamState, x, y, width, height int) (ParameterView, error)

	// CreateWindow creates a window for parameter visualization
	CreateWindow(options WindowOptions) (Window, error)
}

// UISystem provides access to the platform's UI capabilities
type UISystem interface {
	// GetParamViewFactory returns a factory for parameter visualization components
	GetParamViewFactory() ParamViewFactory

	// RunOnMainThread runs the given function on the UI thread
	RunOnMainThread(func()) error

	// IsMainThread returns true if called from the UI thread
	IsMainThread() bool

	// ShowMessageBox shows a simple message box
	ShowMessageBox(title, message string) error

	// ShowConfirmDialog shows a Yes/No dialog
	ShowConfirmDialog(title, message string) (bool, error)

	// ShowInputDialog shows a dialog with input fields
	ShowInputDialog(title string, fields []string, defaults []string) ([]string, error)
}
