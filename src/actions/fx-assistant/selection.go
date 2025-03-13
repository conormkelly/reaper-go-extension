package fxassistant

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"go-reaper/src/reaper/fx"
	"strconv"
	"strings"
	"unsafe"
)

// SelectionResult represents the result of the FX selection process
type SelectionResult struct {
	Collection  fx.TrackCollection // The full track collection
	SelectedFX  map[int][]int      // Map of track index to list of selected FX indices
	Description string             // User-friendly description of what was selected
}

// SelectFXForProcessing presents a UI to the user to select tracks and FX for processing
// Returns a collection containing only the selected FX, and a map of track index to selected FX indices
func SelectFXForProcessing() (*SelectionResult, error) {
	logger.Debug("Starting FX selection process")

	// Step 1: Get all selected tracks
	tracks, err := getSelectedTracks()
	if err != nil {
		return nil, fmt.Errorf("error getting selected tracks: %v", err)
	}

	// Step 2: Build track → FX hierarchy for all tracks
	collection, err := buildTrackFXHierarchy(tracks)
	if err != nil {
		return nil, fmt.Errorf("error building track/FX hierarchy: %v", err)
	}

	if len(collection.Tracks) == 0 {
		return nil, fmt.Errorf("no tracks with FX found")
	}

	// Step 3: Present hierarchy to user and get selection
	selectedFX, description, err := presentFXSelectionUI(collection)
	if err != nil {
		return nil, fmt.Errorf("error during FX selection: %v", err)
	}

	// Step 4: Create a result with the selected FX
	result := &SelectionResult{
		Collection:  collection,
		SelectedFX:  selectedFX,
		Description: description,
	}

	return result, nil
}

// getSelectedTracks gets all currently selected tracks in REAPER
func getSelectedTracks() ([]unsafe.Pointer, error) {
	// First check if we have a valid track selection
	selectedTrack, err := reaper.GetSelectedTrack()
	if err != nil {
		return nil, fmt.Errorf("no track selected: %v", err)
	}

	// For now, we'll just use the first selected track
	// TODO: Enhance to support multiple track selection
	return []unsafe.Pointer{selectedTrack}, nil
}

// buildTrackFXHierarchy creates a complete hierarchy of all tracks and their FX
func buildTrackFXHierarchy(tracks []unsafe.Pointer) (fx.TrackCollection, error) {
	logger.Debug("Building track/FX hierarchy for %d tracks", len(tracks))

	// Use our domain model to build a complete collection
	collection, err := fx.BatchGetMultiTrackFXParameters(tracks, nil)
	if err != nil {
		return fx.TrackCollection{}, fmt.Errorf("failed to get FX parameters: %v", err)
	}

	// Log information about what we found
	totalFX := 0
	totalParams := 0
	for _, track := range collection.Tracks {
		totalFX += len(track.FXList)
		for _, fx := range track.FXList {
			totalParams += len(fx.Parameters)
		}
	}

	logger.Debug("Built hierarchy with %d tracks, %d FX, and %d parameters",
		len(collection.Tracks), totalFX, totalParams)

	return collection, nil
}

// presentFXSelectionUI shows a dialog for selecting FX and returns the selection
// Returns a map of track index to slice of selected FX indices
func presentFXSelectionUI(collection fx.TrackCollection) (map[int][]int, string, error) {
	logger.Debug("Presenting FX selection UI")

	// Build a user-friendly representation of the track/FX hierarchy
	var builder strings.Builder
	builder.WriteString("Select FX to include (comma-separated numbers)\n\n")

	// For each track, list its FX
	for i, track := range collection.Tracks {
		builder.WriteString(fmt.Sprintf("Track %d: %s\n", i+1, track.TrackName))

		// List FX for this track
		for j, fx := range track.FXList {
			builder.WriteString(fmt.Sprintf("  %d.%d: %s\n", i+1, j+1, fx.FXName))
		}
		builder.WriteString("\n")
	}

	// Add instructions
	builder.WriteString("Format: <track>.<fx>, <track>.<fx>, ...\n")
	builder.WriteString("Example: 1.1, 1.3 (Track 1, FX 1 and 3)\n")
	builder.WriteString("Or: all (to select all FX)")

	// Show the dialog to the user
	fields := []string{"Selection:"}
	defaults := []string{"1.1"} // Default to first FX on first track

	results, err := reaper.GetUserInputs("FX Selection", fields, defaults)
	if err != nil {
		logger.Info("User cancelled FX selection dialog")
		return nil, "", fmt.Errorf("selection cancelled")
	}

	// Parse the selection
	selection := results[0]
	logger.Debug("User FX selection: %s", selection)

	// Process the selection
	selectedFX, description, err := parseSelection(selection, collection)
	if err != nil {
		return nil, "", err
	}

	return selectedFX, description, nil
}

// parseSelection parses the user's selection string into a map of track index to FX indices
func parseSelection(selection string, collection fx.TrackCollection) (map[int][]int, string, error) {
	selectedFX := make(map[int][]int)
	selection = strings.TrimSpace(selection)

	// Special case for "all" - select all FX from all tracks
	if strings.ToLower(selection) == "all" {
		for _, track := range collection.Tracks {
			fxIndices := make([]int, len(track.FXList))
			for j := range track.FXList {
				fxIndices[j] = track.FXList[j].FXIndex
			}
			selectedFX[track.TrackIndex] = fxIndices
		}
		return selectedFX, "All FX on all tracks", nil
	}

	// Split by commas
	parts := strings.Split(selection, ",")

	// Process each part (track.fx pair)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by dot to get track and FX indices
		elements := strings.Split(part, ".")
		if len(elements) != 2 {
			return nil, "", fmt.Errorf("invalid format: %s (should be track.fx)", part)
		}

		// Parse track index (1-based)
		trackNum, err := strconv.Atoi(elements[0])
		if err != nil {
			return nil, "", fmt.Errorf("invalid track number: %s", elements[0])
		}
		trackIdx := trackNum - 1 // Convert to 0-based

		// Parse FX index (1-based)
		fxNum, err := strconv.Atoi(elements[1])
		if err != nil {
			return nil, "", fmt.Errorf("invalid FX number: %s", elements[1])
		}
		fxIdx := fxNum - 1 // Convert to 0-based

		// Validate track index
		if trackIdx < 0 || trackIdx >= len(collection.Tracks) {
			return nil, "", fmt.Errorf("track number out of range: %d", trackNum)
		}

		// Get the actual track
		track := collection.Tracks[trackIdx]

		// Find the corresponding FX index in the track's FX list
		validFX := false
		for _, fx := range track.FXList {
			if fx.FXIndex == fxIdx {
				validFX = true
				break
			}
		}

		if !validFX {
			return nil, "", fmt.Errorf("FX number %d not found on track %d", fxNum, trackNum)
		}

		// Add to selected FX map
		if selectedFX[track.TrackIndex] == nil {
			selectedFX[track.TrackIndex] = []int{}
		}
		selectedFX[track.TrackIndex] = append(selectedFX[track.TrackIndex], fxIdx)
	}

	if len(selectedFX) == 0 {
		return nil, "", fmt.Errorf("no valid FX selected")
	}

	// Generate a user-friendly description
	description := generateSelectionDescription(selectedFX, collection)

	return selectedFX, description, nil
}

// generateSelectionDescription creates a user-friendly description of the selected FX
func generateSelectionDescription(selectedFX map[int][]int, collection fx.TrackCollection) string {
	var builder strings.Builder

	totalTracks := len(selectedFX)
	totalFX := 0

	for _, fxList := range selectedFX {
		totalFX += len(fxList)
	}

	builder.WriteString(fmt.Sprintf("%d FX on %d track(s):\n", totalFX, totalTracks))

	for trackIdx, fxIndices := range selectedFX {
		// Find the track in the collection
		var track fx.TrackWithFX
		for _, t := range collection.Tracks {
			if t.TrackIndex == trackIdx {
				track = t
				break
			}
		}

		builder.WriteString(fmt.Sprintf("• %s: ", track.TrackName))

		// Add FX names
		fxNames := make([]string, 0, len(fxIndices))
		for _, fxIdx := range fxIndices {
			for _, fx := range track.FXList {
				if fx.FXIndex == fxIdx {
					fxNames = append(fxNames, fx.FXName)
					break
				}
			}
		}

		builder.WriteString(strings.Join(fxNames, ", "))
		builder.WriteString("\n")
	}

	return builder.String()
}
