// src/ui/params/param_view.go - Parameter visualization components
package params

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/ui"
	"go-reaper/src/ui/common"
)

// ParamViewOptions configures the appearance and behavior of a parameter view
type ParamViewOptions struct {
	// Show the explanation
	ShowExplanation bool

	// Show the original value
	ShowOriginal bool

	// Show value as numeric text
	ShowNumeric bool

	// Custom colors
	SliderColor *common.ColorRGB
	TextColor   *common.ColorRGB

	// Additional labels
	CustomLabels map[string]string
}

// DefaultParamViewOptions provides sensible defaults
var DefaultParamViewOptions = ParamViewOptions{
	ShowExplanation: true,
	ShowOriginal:    true,
	ShowNumeric:     true,
	SliderColor:     nil, // Use default
	TextColor:       nil, // Use default
	CustomLabels:    nil,
}

// A simplified wrapper around the platform-specific parameter view
type ParamView struct {
	// Parameter data
	Param common.ParamState

	// UI component
	view common.ParameterView

	// Parent window
	window common.Window

	// Options
	options ParamViewOptions

	// Change callback
	onChangeCallback func(param *common.ParamState)

	// Position and size
	x, y, width, height int
}

// NewParamView creates a new parameter view in the given window
func NewParamView(window common.Window, param common.ParamState, x, y, width, height int, options *ParamViewOptions) (*ParamView, error) {
	// Use default options if none provided
	opts := DefaultParamViewOptions
	if options != nil {
		opts = *options
	}

	// Create platform-specific view
	view, err := ui.CreateParamView(window, param, x, y, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter view: %v", err)
	}

	// Create wrapper
	paramView := &ParamView{
		Param:   param,
		view:    view,
		window:  window,
		options: opts,
		x:       x,
		y:       y,
		width:   width,
		height:  height,
	}

	// Setup internal value change callback
	view.OnValueChanged(func(value float64) {
		paramView.Param.Value = value
		if paramView.onChangeCallback != nil {
			paramView.onChangeCallback(&paramView.Param)
		}
	})

	return paramView, nil
}

// Show displays the parameter view
func (p *ParamView) Show() error {
	return p.view.Show()
}

// Hide hides the parameter view
func (p *ParamView) Hide() error {
	return p.view.Hide()
}

// SetValue updates the parameter value
func (p *ParamView) SetValue(value float64) error {
	p.Param.Value = value
	return p.view.SetValue(value)
}

// GetValue returns the current parameter value
func (p *ParamView) GetValue() float64 {
	return p.view.GetValue()
}

// SetFormattedValue updates the formatted value display
func (p *ParamView) SetFormattedValue(formatted string) error {
	p.Param.FormattedValue = formatted
	return p.view.SetFormattedValue(formatted)
}

// SetExplanation updates the explanation text
func (p *ParamView) SetExplanation(text string) error {
	p.Param.Explanation = text
	return p.view.SetExplanation(text)
}

// SetOriginalValue sets the original value for comparison
func (p *ParamView) SetOriginalValue(value float64, formatted string) error {
	p.Param.OriginalValue = value
	p.Param.OriginalFormattedValue = formatted
	return p.view.SetOriginalValue(value, formatted)
}

// OnChange sets the callback for value changes
func (p *ParamView) OnChange(callback func(param *common.ParamState)) {
	p.onChangeCallback = callback
}

// Create a group of parameter views from multiple parameters
func CreateParamGroup(window common.Window, params []common.ParamState,
	x, y, width, height int, options *ParamViewOptions) ([]*ParamView, error) {
	if len(params) == 0 {
		return nil, nil
	}

	// Calculate vertical space per param
	paramHeight := height / len(params)

	// Create views
	views := make([]*ParamView, len(params))
	for i, param := range params {
		paramY := y + (i * paramHeight)
		view, err := NewParamView(window, param, x, paramY, width, paramHeight-5, options)
		if err != nil {
			logger.Warning("Failed to create parameter view for %s: %v", param.Name, err)
			continue
		}

		views[i] = view
	}

	return views, nil
}

// Apply parameter changes to a group of parameter views
func ApplyParamChanges(views []*ParamView, values []float64) error {
	if len(views) != len(values) {
		return fmt.Errorf("number of views (%d) doesn't match number of values (%d)",
			len(views), len(values))
	}

	for i, view := range views {
		if err := view.SetValue(values[i]); err != nil {
			return fmt.Errorf("failed to set value for parameter %s: %v",
				view.Param.Name, err)
		}
	}

	return nil
}
