// Package fx provides FX and parameter management for REAPER tracks
package fx

/*
#cgo CFLAGS: -I${SRCDIR}/../../c -I${SRCDIR}/../../../sdk
#include "../../c/api/fx.h"
#include <stdlib.h>

// Define CGo type mappings for void** to make the conversions cleaner
typedef void** track_ptr_array;
*/
import "C"
import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"unsafe"
)

// BatchGetMultiTrackFXParameters retrieves parameters for multiple tracks and FX in one operation
// This reduces the number of CGo transitions dramatically for multi-track operations
func BatchGetMultiTrackFXParameters(tracks []unsafe.Pointer, fxIndices [][]int) (TrackCollection, error) {
	collection := TrackCollection{
		Tracks: make([]TrackWithFX, 0, len(tracks)),
	}

	trackCount := len(tracks)
	if trackCount == 0 {
		return collection, nil
	}

	// Prepare for the C-API call
	cTracks := make([]unsafe.Pointer, trackCount)
	cFXCounts := make([]C.int, trackCount)
	cFXIndicesPointers := make([]*C.int, trackCount)
	totalFXCount := 0

	// Initialize track data and calculate total FX count
	for i, track := range tracks {
		if track == nil {
			continue
		}

		cTracks[i] = track

		// Get track name for our collection
		trackName, err := reaper.GetTrackName(track)
		if err != nil {
			trackName = fmt.Sprintf("Track %d", i)
		}

		trackWithFX := TrackWithFX{
			TrackIndex: i,
			TrackName:  trackName,
			MediaTrack: track,
			FXList:     make([]FXWithParams, 0),
		}

		// Add track to our collection (we'll fill in FX details later)
		collection.Tracks = append(collection.Tracks, trackWithFX)

		// Determine which FX indices to use for this track
		var trackFXIndices []int
		if i < len(fxIndices) && fxIndices[i] != nil {
			trackFXIndices = fxIndices[i]
		} else {
			// Get all FX on track if no specific indices were provided
			fxCount, err := reaper.GetTrackFXCount(track)
			if err != nil {
				logger.Warning("Could not get FX count for track %d: %v", i, err)
				continue
			}

			trackFXIndices = make([]int, fxCount)
			for j := 0; j < fxCount; j++ {
				trackFXIndices[j] = j
			}
		}

		// Create C array of FX indices
		countFX := len(trackFXIndices)
		if countFX > 0 {
			cFXCounts[i] = C.int(countFX)
			cIndices := make([]C.int, countFX)
			for j, idx := range trackFXIndices {
				cIndices[j] = C.int(idx)
			}

			// Get pointer to the array
			cFXIndicesPointers[i] = &cIndices[0]

			// Update total FX count
			totalFXCount += countFX
		}
	}

	if totalFXCount == 0 {
		return collection, nil // No FX to process
	}

	// For single track operations, fall back to the existing simpler API
	// This makes development easier as we can implement the full multi-track
	// version later while having something working now
	if trackCount == 1 && tracks[0] != nil {
		// Process a single track using existing API
		// Repopulate the collection we started earlier
		trackIdx := 0
		track := tracks[trackIdx]

		// Get specific FX indices or all FX
		var trackFXIndices []int
		if len(fxIndices) > 0 && fxIndices[0] != nil {
			trackFXIndices = fxIndices[0]
		} else {
			// Get all FX on track
			fxCount, err := reaper.GetTrackFXCount(track)
			if err != nil {
				return collection, fmt.Errorf("could not get FX count: %v", err)
			}

			trackFXIndices = make([]int, fxCount)
			for j := 0; j < fxCount; j++ {
				trackFXIndices[j] = j
			}
		}

		// Get FX parameters for each FX
		for _, fxIndex := range trackFXIndices {
			fxInfo, err := reaper.GetFXParameters(track, fxIndex)
			if err != nil {
				logger.Warning("Could not get FX parameters for FX %d: %v", fxIndex, err)
				continue
			}

			// Convert reaper.FXParameter to our ParameterState
			parameters := make([]ParameterState, len(fxInfo.Parameters))
			for j, param := range fxInfo.Parameters {
				parameters[j] = ParameterState{
					FXIndex:        fxIndex,
					ParamIndex:     param.Index,
					ParamName:      param.Name,
					Value:          param.Value,
					FormattedValue: param.FormattedValue,
					Min:            param.Min,
					Max:            param.Max,
					MinFormatted:   param.MinFormatted,
					MaxFormatted:   param.MaxFormatted,
				}
			}

			// Add FX to the track
			fxWithParams := FXWithParams{
				FXIndex:    fxIndex,
				FXName:     fxInfo.Name,
				Parameters: parameters,
			}

			collection.Tracks[0].FXList = append(collection.Tracks[0].FXList, fxWithParams)
		}

		return collection, nil
	}

	// TODO: Implement full multi-track batch operation
	// For now, let's fall back to the existing track-by-track implementation

	// Loop through each track in our collection and fill in FX data
	for i, trackFX := range collection.Tracks {
		track := trackFX.MediaTrack
		if track == nil {
			continue
		}

		// Determine which FX indices to use for this track
		var trackFXIndices []int
		if i < len(fxIndices) && fxIndices[i] != nil {
			trackFXIndices = fxIndices[i]
		} else {
			// Get all FX on track if no specific indices were provided
			fxCount, err := reaper.GetTrackFXCount(track)
			if err != nil {
				logger.Warning("Could not get FX count for track %d: %v", i, err)
				continue
			}

			trackFXIndices = make([]int, fxCount)
			for j := 0; j < fxCount; j++ {
				trackFXIndices[j] = j
			}
		}

		// Get FX parameters for each FX
		for _, fxIndex := range trackFXIndices {
			fxInfo, err := reaper.GetFXParameters(track, fxIndex)
			if err != nil {
				logger.Warning("Could not get FX parameters for track %d, FX %d: %v", i, fxIndex, err)
				continue
			}

			// Convert reaper.FXParameter to our ParameterState
			parameters := make([]ParameterState, len(fxInfo.Parameters))
			for j, param := range fxInfo.Parameters {
				parameters[j] = ParameterState{
					FXIndex:        fxIndex,
					ParamIndex:     param.Index,
					ParamName:      param.Name,
					Value:          param.Value,
					FormattedValue: param.FormattedValue,
					Min:            param.Min,
					Max:            param.Max,
					MinFormatted:   param.MinFormatted,
					MaxFormatted:   param.MaxFormatted,
				}
			}

			// Add FX to the track
			fxWithParams := FXWithParams{
				FXIndex:    fxIndex,
				FXName:     fxInfo.Name,
				Parameters: parameters,
			}

			collection.Tracks[i].FXList = append(collection.Tracks[i].FXList, fxWithParams)
		}
	}

	return collection, nil
}

// BatchFormatParameterValues formats multiple parameter values in a single call
// This is much more efficient than making multiple CGo transitions
func BatchFormatParameterValues(tracks []unsafe.Pointer, requests []ParameterFormatRequest) ([]string, error) {
	if len(tracks) == 0 {
		return nil, fmt.Errorf("no tracks provided")
	}

	if len(requests) == 0 {
		return []string{}, nil
	}

	// Validate all requests
	for i, req := range requests {
		if req.TrackIndex < 0 || req.TrackIndex >= len(tracks) {
			return nil, fmt.Errorf("invalid track index %d in request %d", req.TrackIndex, i)
		}

		if tracks[req.TrackIndex] == nil {
			return nil, fmt.Errorf("nil track pointer for index %d in request %d", req.TrackIndex, i)
		}
	}

	// If only one track is involved, use the simpler API
	if len(tracks) == 1 && tracks[0] != nil {
		// Check if all requests are for this track
		singleTrack := true
		for _, req := range requests {
			if req.TrackIndex != 0 {
				singleTrack = false
				break
			}
		}

		if singleTrack {
			// Convert our requests to reaper.ParameterFormatRequest format
			reaperRequests := make([]reaper.ParameterFormatRequest, len(requests))
			for i, req := range requests {
				reaperRequests[i] = reaper.ParameterFormatRequest{
					FXIndex:    req.FXIndex,
					ParamIndex: req.ParamIndex,
					Value:      req.Value,
				}
			}

			// Call the existing batch format function
			return reaper.BatchFormatFXParameters(tracks[0], reaperRequests)
		}
	}

	// Allocate memory for the format requests
	cFormatRequests := C.malloc(C.size_t(len(requests) * int(unsafe.Sizeof(C.fx_param_format_t{}))))
	if cFormatRequests == nil {
		return nil, fmt.Errorf("failed to allocate memory for format requests")
	}
	defer C.free(cFormatRequests)

	// Fill in the format request data
	formatSlice := (*[1 << 30]C.fx_param_format_t)(cFormatRequests)[:len(requests):len(requests)]

	for i, req := range requests {
		formatSlice[i].track_index = C.int(req.TrackIndex)
		formatSlice[i].fx_index = C.int(req.FXIndex)
		formatSlice[i].param_index = C.int(req.ParamIndex)
		formatSlice[i].value = C.double(req.Value)
	}

	// Create Go slice of tracks
	tracksPtr := make([]unsafe.Pointer, len(tracks))
	for i, track := range tracks {
		tracksPtr[i] = track
	}

	// Use our defined type to properly convert to void**
	cTracks := C.track_ptr_array(unsafe.Pointer(&tracksPtr[0]))

	// Call the batch function with our properly typed tracks array
	result := C.plugin_bridge_batch_format_multi_fx_parameters(
		cTracks,
		(*C.fx_param_format_t)(cFormatRequests),
		C.int(len(requests)),
	)

	if !bool(result) {
		return nil, fmt.Errorf("failed to format parameter values")
	}

	// Extract formatted values
	formatted := make([]string, len(requests))
	for i := 0; i < len(requests); i++ {
		formatted[i] = C.GoString(&formatSlice[i].formatted[0])
	}

	return formatted, nil
}

// BatchSetMultiTrackFXParameters applies parameter changes across multiple tracks in a single operation
// This is much more efficient than making multiple CGo transitions
func BatchSetMultiTrackFXParameters(tracks []unsafe.Pointer, changes []ParameterChange) error {
	if len(changes) == 0 {
		return nil
	}

	// Validate all changes before proceeding
	validChanges := make([]ParameterChange, 0, len(changes))
	for _, change := range changes {
		// Validate track index
		if change.TrackIndex < 0 || change.TrackIndex >= len(tracks) {
			logger.Warning("Invalid track index %d, skipping change", change.TrackIndex)
			continue
		}

		// Skip nil tracks
		if tracks[change.TrackIndex] == nil {
			logger.Warning("Nil track pointer for index %d, skipping change", change.TrackIndex)
			continue
		}

		validChanges = append(validChanges, change)
	}

	// If there's only one track involved, use the simpler API
	if len(tracks) == 1 && tracks[0] != nil {
		// Convert our changes to reaper.ParameterChange
		reaperChanges := make([]reaper.ParameterChange, len(validChanges))
		for i, change := range validChanges {
			reaperChanges[i] = reaper.ParameterChange{
				FXIndex:    change.FXIndex,
				ParamIndex: change.ParamIndex,
				Value:      change.Value,
			}
		}

		// Apply changes using the existing function
		return reaper.BatchSetFXParameters(tracks[0], reaperChanges)
	}

	// Create Go slice of tracks
	tracksPtr := make([]unsafe.Pointer, len(tracks))
	for i, track := range tracks {
		tracksPtr[i] = track
	}

	// Use our defined type to properly convert to void**
	cTracks := C.track_ptr_array(unsafe.Pointer(&tracksPtr[0]))

	// Allocate C memory for the changes
	cChanges := C.malloc(C.size_t(len(validChanges) * int(unsafe.Sizeof(C.fx_param_multi_change_t{}))))
	if cChanges == nil {
		return fmt.Errorf("failed to allocate memory for parameter changes")
	}
	defer C.free(cChanges)

	// Fill in the changes data
	// This creates a slice view of the C array
	changesSlice := (*[1 << 30]C.fx_param_multi_change_t)(cChanges)[:len(validChanges):len(validChanges)]

	for i, change := range validChanges {
		changesSlice[i].track_index = C.int(change.TrackIndex)
		changesSlice[i].fx_index = C.int(change.FXIndex)
		changesSlice[i].param_index = C.int(change.ParamIndex)
		changesSlice[i].value = C.double(change.Value)
	}

	// Call the multi-track batch function
	result := C.plugin_bridge_batch_set_multi_track_fx_parameters(
		cTracks,
		(*C.fx_param_multi_change_t)(cChanges),
		C.int(len(validChanges)),
	)

	if !bool(result) {
		return fmt.Errorf("failed to apply parameter changes")
	}

	logger.Debug("Applied %d parameter changes across %d tracks", len(validChanges), len(tracks))
	return nil
}

// BatchSetMultiTrackFXParametersWithUndo applies parameter changes across multiple tracks
// with undo support to allow for undoing/redoing the entire operation
func BatchSetMultiTrackFXParametersWithUndo(tracks []unsafe.Pointer, changes []ParameterChange, undoLabel string) error {
	// Start undo block
	if err := reaper.BeginUndoBlock(undoLabel); err != nil {
		logger.Warning("Could not start undo block: %v", err)
		// Continue anyway, just without undo support
	}

	// Apply the changes
	err := BatchSetMultiTrackFXParameters(tracks, changes)

	// End undo block (even if there was an error)
	if endErr := reaper.EndUndoBlock(undoLabel, 0); endErr != nil {
		logger.Warning("Could not end undo block: %v", endErr)
	}

	return err
}

// ApplyParameterModifications applies a list of parameter modifications with undo support
// This converts ParameterModification to ParameterChange and applies them
func ApplyParameterModifications(tracks []unsafe.Pointer, modifications []ParameterModification, undoLabel string) error {
	if len(modifications) == 0 {
		return nil
	}

	// Convert modifications to changes
	changes := make([]ParameterChange, len(modifications))
	for i, mod := range modifications {
		changes[i] = ParameterChange{
			TrackIndex: mod.TrackIndex,
			FXIndex:    mod.FXIndex,
			ParamIndex: mod.ParamIndex,
			Value:      mod.NewValue,
		}
	}

	// Apply with undo support
	return BatchSetMultiTrackFXParametersWithUndo(tracks, changes, undoLabel)
}
