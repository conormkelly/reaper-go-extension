package common

// Common UI type definitions

// ParamStatus represents the current state of a parameter
type ParamState struct {
	// Name of the parameter
	Name string

	// Current normalized value (0.0-1.0)
	Value float64

	// Human-readable formatted value
	FormattedValue string

	// Original value before modification
	OriginalValue float64

	// Original formatted value
	OriginalFormattedValue string

	// Minimum and maximum values
	Min, Max float64

	// Index of the parameter
	Index int

	// FX index this parameter belongs to
	FXIndex int

	// Optional contextual explanation of parameter (for LLM responses)
	Explanation string
}

// WindowOptions defines common options for UI windows
type WindowOptions struct {
	// Window title
	Title string

	// Position and size
	X, Y, Width, Height int

	// Is this window resizable?
	Resizable bool
}

// TextOptions defines text styling options
type TextOptions struct {
	// Bold text
	Bold bool

	// Font size (points)
	Size float64
}

// ColorRGB defines an RGB color
type ColorRGB struct {
	R, G, B uint8
}

// Common colors
var (
	ColorBlack = ColorRGB{0, 0, 0}
	ColorWhite = ColorRGB{255, 255, 255}
	ColorGray  = ColorRGB{128, 128, 128}
	ColorBlue  = ColorRGB{0, 0, 255}
	ColorRed   = ColorRGB{255, 0, 0}
	ColorGreen = ColorRGB{0, 128, 0}
)

// Callback function types
type ActionCallback func()
type ValueChangeCallback func(value float64)
type FormSubmitCallback func(values map[string]string)
