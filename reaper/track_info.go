package reaper

/*
#cgo CFLAGS: -I${SRCDIR}/../c -I${SRCDIR}/../sdk
#include "../c/bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// TrackInfo represents information about a REAPER track
type TrackInfo struct {
	MediaTrack unsafe.Pointer
	Index      int
	Name       string
	NumFX      int
}

// GetSelectedTrackInfo gets detailed information about the selected track
func GetSelectedTrackInfo() (*TrackInfo, error) {
	// Get selected track pointer
	track, err := GetSelectedTrack()
	if err != nil {
		return nil, err
	}

	// Create track info object
	trackInfo := &TrackInfo{
		MediaTrack: track,
	}

	// Get track index
	cFuncName := C.CString("GetMediaTrackInfo_Value")
	defer C.free(unsafe.Pointer(cFuncName))

	getTrackInfoPtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getTrackInfoPtr == nil {
		return nil, fmt.Errorf("could not get GetMediaTrackInfo_Value function pointer")
	}

	// Get track index (IP_TRACKNUMBER)
	cParam := C.CString("IP_TRACKNUMBER")
	defer C.free(unsafe.Pointer(cParam))

	// Define function pointer type for GetMediaTrackInfo_Value
	getTrackInfoFunc := unsafe.Pointer(getTrackInfoPtr)
	index := C.plugin_bridge_call_track_get_info_value(getTrackInfoFunc, track, cParam)
	trackInfo.Index = int(index)

	// Get track name
	trackName, err := GetTrackName(track)
	if err == nil {
		trackInfo.Name = trackName
	} else {
		trackInfo.Name = fmt.Sprintf("Track %d", trackInfo.Index)
	}

	// Get FX count
	fxCount, err := GetTrackFXCount(track)
	if err == nil {
		trackInfo.NumFX = fxCount
	}

	return trackInfo, nil
}

// GetTrackName gets the name of a track
func GetTrackName(track unsafe.Pointer) (string, error) {
	if !initialized {
		return "", fmt.Errorf("REAPER functions not initialized")
	}

	cFuncName := C.CString("GetTrackName")
	defer C.free(unsafe.Pointer(cFuncName))

	getTrackNamePtr := C.plugin_bridge_call_get_func(C.plugin_bridge_get_get_func(), cFuncName)
	if getTrackNamePtr == nil {
		return "", fmt.Errorf("could not get GetTrackName function pointer")
	}

	// Allocate buffers for track name and flags
	nameBuf := (*C.char)(C.malloc(C.size_t(256)))
	defer C.free(unsafe.Pointer(nameBuf))

	flagsBuf := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0)))))
	defer C.free(unsafe.Pointer(flagsBuf))

	// Call GetTrackName
	result := C.plugin_bridge_call_get_track_name(getTrackNamePtr, track, nameBuf, 256, flagsBuf)
	if !bool(result) {
		return "", fmt.Errorf("failed to get track name")
	}

	return C.GoString(nameBuf), nil
}

// GetTrackFXList gets a list of all FX on a track
func GetTrackFXList(track unsafe.Pointer) ([]FXInfo, error) {
	// Get FX count
	fxCount, err := GetTrackFXCount(track)
	if err != nil {
		return nil, fmt.Errorf("failed to get FX count: %v", err)
	}

	// Gather info for all FX
	result := make([]FXInfo, 0, fxCount)
	for i := 0; i < fxCount; i++ {
		// Get just FX name for the list (optimization)
		fxName, err := GetTrackFXName(track, i)
		if err != nil {
			return nil, fmt.Errorf("failed to get FX name: %v", err)
		}

		// Create minimal FX info (don't load all parameters yet for performance)
		fxInfo := FXInfo{
			Index: i,
			Name:  fxName,
		}

		result = append(result, fxInfo)
	}

	return result, nil
}
