// Package fx provides FX and parameter management for REAPER tracks
package fx

/*
#cgo CFLAGS: -I${SRCDIR}/../../c -I${SRCDIR}/../../../sdk
#include "../../c/api/fx.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"unsafe"
)

// GetTrackCollection retrieves all FX and parameters for the selected tracks
// This is a high-level function that creates a complete track collection
func GetTrackCollection() (TrackCollection, error) {
	// Get all selected tracks
	var tracks []unsafe.Pointer
	var err error

	// For now, we'll just use the first selected track
	// TODO: Enhance to support multiple track selection
	track, err := reaper.GetSelectedTrack()
	if err != nil {
		return TrackCollection{}, fmt.Errorf("no track selected: %v", err)
	}

	tracks = append(tracks, track)

	// Get all FX for each track
	collection, err := BatchGetMultiTrackFXParameters(tracks, nil)
	if err != nil {
		return TrackCollection{}, fmt.Errorf("failed to get FX parameters: %v", err)
	}

	return collection, nil
}

// GetSelectedFXFromCollection extracts a subset of FX from a track collection based on user selection
// This allows the user to focus on specific FX they want to work with
func GetSelectedFXFromCollection(collection TrackCollection, selectedFX map[int][]int) (TrackCollection, error) {
	result := TrackCollection{
		Tracks: make([]TrackWithFX, 0),
	}

	// Process each track
	for _, track := range collection.Tracks {
		// Check if we have selected FX for this track
		fxIndices, hasSelection := selectedFX[track.TrackIndex]
		if !hasSelection {
			// No selection for this track, skip it
			continue
		}

		// Create a new track with only the selected FX
		newTrack := TrackWithFX{
			TrackIndex: track.TrackIndex,
			TrackName:  track.TrackName,
			MediaTrack: track.MediaTrack,
			FXList:     make([]FXWithParams, 0),
		}

		// Add only the selected FX
		for _, fxIndex := range fxIndices {
			// Find the FX in the original track
			for _, fx := range track.FXList {
				if fx.FXIndex == fxIndex {
					newTrack.FXList = append(newTrack.FXList, fx)
					break
				}
			}
		}

		// Add the track to the result if it has any FX
		if len(newTrack.FXList) > 0 {
			result.Tracks = append(result.Tracks, newTrack)
		}
	}

	return result, nil
}

// FormatCollectionForDisplay creates a human-readable representation of a track collection
// This is useful for displaying the current state to the user
func FormatCollectionForDisplay(collection TrackCollection) string {
	if len(collection.Tracks) == 0 {
		return "No tracks or FX selected."
	}

	var result string

	// Format each track
	for _, track := range collection.Tracks {
		result += fmt.Sprintf("Track: %s\n", track.TrackName)

		// Format each FX
		for _, fx := range track.FXList {
			result += fmt.Sprintf("  FX: %s\n", fx.FXName)

			// Format each parameter
			for _, param := range fx.Parameters {
				result += fmt.Sprintf("    %s: %s\n", param.ParamName, param.FormattedValue)
			}

			result += "\n"
		}

		result += "\n"
	}

	return result
}

// GetTrackMediaTracks extracts the MediaTrack pointers from a collection
// This is useful for operations that need the raw REAPER track pointers
func GetTrackMediaTracks(collection TrackCollection) []unsafe.Pointer {
	tracks := make([]unsafe.Pointer, len(collection.Tracks))

	for i, track := range collection.Tracks {
		tracks[i] = track.MediaTrack
	}

	return tracks
}

// GetFXCount returns the total number of FX in a collection
func GetFXCount(collection TrackCollection) int {
	count := 0

	for _, track := range collection.Tracks {
		count += len(track.FXList)
	}

	return count
}

// GetParameterCount returns the total number of parameters in a collection
func GetParameterCount(collection TrackCollection) int {
	count := 0

	for _, track := range collection.Tracks {
		for _, fx := range track.FXList {
			count += len(fx.Parameters)
		}
	}

	return count
}

// ValidateParameterModifications ensures that parameter modifications
// are valid for the given track collection
func ValidateParameterModifications(collection TrackCollection, modifications []ParameterModification) []ParameterModification {
	valid := make([]ParameterModification, 0, len(modifications))

	// Create a map for quick lookup of tracks, FX, and parameters
	trackMap := make(map[int]int) // TrackIndex -> Index in collection.Tracks
	for i, track := range collection.Tracks {
		trackMap[track.TrackIndex] = i

		// Also create a map for FX
		fxMap := make(map[int]bool) // FXIndex -> exists
		for _, fx := range track.FXList {
			fxMap[fx.FXIndex] = true

			// Add parameters to a map too
			paramMap := make(map[int]bool) // ParamIndex -> exists
			for _, param := range fx.Parameters {
				paramMap[param.ParamIndex] = true
			}
		}
	}

	// Check each modification
	for _, mod := range modifications {
		// Check track
		trackIdx, hasTrack := trackMap[mod.TrackIndex]
		if !hasTrack {
			logger.Warning("Track index %d not found in collection, skipping modification", mod.TrackIndex)
			continue
		}

		// Check FX
		hasFX := false
		for _, fx := range collection.Tracks[trackIdx].FXList {
			if fx.FXIndex == mod.FXIndex {
				hasFX = true

				// Check parameter
				hasParam := false
				for _, param := range fx.Parameters {
					if param.ParamIndex == mod.ParamIndex {
						hasParam = true
						break
					}
				}

				if !hasParam {
					logger.Warning("Parameter index %d not found in FX %d on track %d, skipping modification",
						mod.ParamIndex, mod.FXIndex, mod.TrackIndex)
					continue
				}

				break
			}
		}

		if !hasFX {
			logger.Warning("FX index %d not found on track %d, skipping modification",
				mod.FXIndex, mod.TrackIndex)
			continue
		}

		// If we got here, the modification is valid
		valid = append(valid, mod)
	}

	return valid
}
