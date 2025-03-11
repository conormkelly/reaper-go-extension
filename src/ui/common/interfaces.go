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

	// AddTextField adds a text field
	AddTextField(placeholder string, x, y, width, height int) error

	// SetTitle changes the window title
	SetTitle(title string) error
}

// UISystem provides access to the platform's UI capabilities
type UISystem interface {
	// CreateWindow creates a window with the given options
	CreateWindow(options WindowOptions) (Window, error)

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
