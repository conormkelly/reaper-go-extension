// Package fx provides FX and parameter management for REAPER tracks
package fx

import (
	"unsafe"
)

// TrackCollection represents a collection of tracks with their FX
// This is the top-level structure for working with multiple tracks
type TrackCollection struct {
	Tracks []TrackWithFX
}

// TrackWithFX represents a track with its FX list
// This allows operations on all FX in a track
type TrackWithFX struct {
	TrackIndex int            // REAPER track index (0-based)
	TrackName  string         // Display name of the track
	MediaTrack unsafe.Pointer // REAPER's internal track pointer
	FXList     []FXWithParams // List of FX on this track
}

// FXWithParams represents an FX with its parameters
// This allows operations on all parameters of an FX
type FXWithParams struct {
	FXIndex    int              // FX index on the track (0-based)
	FXName     string           // Display name of the FX
	Parameters []ParameterState // List of parameters for this FX
}

// ParameterState represents the current state of a parameter
// This includes the current value, name, and metadata
type ParameterState struct {
	FXIndex        int     // FX index this parameter belongs to
	ParamIndex     int     // Parameter index within the FX (0-based)
	ParamName      string  // Display name of the parameter
	Value          float64 // Current normalized value (0.0-1.0)
	FormattedValue string  // Human-readable formatted value
	Min            float64 // Minimum normalized value (typically 0.0)
	Max            float64 // Maximum normalized value (typically 1.0)
	MinFormatted   string  // Human-readable minimum value
	MaxFormatted   string  // Human-readable maximum value
}

// ParameterChange represents a parameter value change to be applied
// Used for batch operations to modify parameters
type ParameterChange struct {
	TrackIndex int     // Track index (0-based)
	FXIndex    int     // FX index within the track (0-based)
	ParamIndex int     // Parameter index within the FX (0-based)
	Value      float64 // New normalized value to set (0.0-1.0)
}

// ParameterModification represents a suggested parameter change
// This includes before/after values and explanations
type ParameterModification struct {
	TrackIndex        int     // Track index (0-based)
	TrackName         string  // Display name of the track
	FXIndex           int     // FX index (0-based)
	FXName            string  // Display name of the FX
	ParamIndex        int     // Parameter index (0-based)
	ParamName         string  // Display name of the parameter
	OriginalValue     float64 // Original normalized value (0.0-1.0)
	NewValue          float64 // New normalized value to set (0.0-1.0)
	OriginalFormatted string  // Human-readable original value
	NewFormatted      string  // Human-readable new value
	Explanation       string  // Explanation for the change
}

// ParameterFormatRequest defines a request to format a parameter value
// Used for batch operations to get formatted values
type ParameterFormatRequest struct {
	TrackIndex int     // Track index (0-based)
	FXIndex    int     // FX index (0-based)
	ParamIndex int     // Parameter index (0-based)
	Value      float64 // Value to format (0.0-1.0)
}
