package reaper

// FXParameter represents a single parameter of an FX
type FXParameter struct {
	Index          int     `json:"index"`
	Name           string  `json:"name"`
	Value          float64 `json:"value"`          // Normalized value (0.0-1.0)
	FormattedValue string  `json:"formattedValue"` // Human-readable value
	Min            float64 `json:"min"`            // Minimum value
	Max            float64 `json:"max"`            // Maximum value
	MinFormatted   string  `json:"minFormatted"`   // Human-readable minimum value
	MaxFormatted   string  `json:"maxFormatted"`   // Human-readable maximum value
}

// ParameterChange represents a parameter value change to be applied
type ParameterChange struct {
	FXIndex    int     `json:"fxIndex"`    // FX index
	ParamIndex int     `json:"paramIndex"` // Parameter index
	Value      float64 `json:"value"`      // New value to set
}

// FXInfo represents an FX and its parameters
type FXInfo struct {
	Index      int           `json:"index"`
	Name       string        `json:"name"`
	Parameters []FXParameter `json:"parameters"`
}

// ActionHandler defines a function type for handling actions
type ActionHandler func()

// Section ID constants
const (
	SectionMain          = 0
	SectionMainAlt       = 100
	SectionMIDIEditor    = 32060
	SectionMIDIEventList = 32061
	SectionMIDIInline    = 32062
	SectionMediaExplorer = 32063
)
