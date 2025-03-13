package fxassistant

import (
	"encoding/json"
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper/fx"
	"strings"
)

// AssistantResponse represents the structured response from the LLM
type AssistantResponse struct {
	Message string        `json:"message"`
	Changes []ChangeEntry `json:"changes"`
}

// ChangeEntry represents a single parameter change suggestion
type ChangeEntry struct {
	TrackIndex        int     `json:"track_index"`
	TrackName         string  `json:"track_name"`
	FXIndex           int     `json:"fx_index"`
	FXName            string  `json:"fx_name"`
	ParamIndex        int     `json:"param_index"`
	ParamName         string  `json:"param_name"`
	OriginalValue     float64 `json:"original_value"`
	NewValue          float64 `json:"new_value"`
	OriginalFormatted string  `json:"original_formatted"`
	NewFormatted      string  `json:"new_formatted"`
	Explanation       string  `json:"explanation"`
}

// parseAssistantResponse parses the JSON response from the LLM
func parseAssistantResponse(responseText string, collection fx.TrackCollection, selectedFX map[int][]int) ([]fx.ParameterModification, string, error) {
	// Validate input
	if responseText == "" {
		return nil, "", fmt.Errorf("empty response text from LLM")
	}

	logger.Debug("Parsing response text (%d chars)...", len(responseText))

	// Try to extract JSON from the response (it might contain additional text)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		logger.Error("Failed to find valid JSON markers in response: %s", responseText)
		return nil, "", fmt.Errorf("could not find valid JSON in response")
	}

	jsonStr := responseText[jsonStart : jsonEnd+1]
	logger.Debug("Extracted JSON (%d chars)", len(jsonStr))

	// Parse the JSON response
	var response AssistantResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		logger.Error("JSON unmarshal error: %v", err)
		return nil, "", fmt.Errorf("failed to parse LLM response: %v", err)
	}

	// Check for empty changes - not an error, just no suggestions
	if len(response.Changes) == 0 {
		logger.Info("LLM did not suggest any parameter changes")
		return []fx.ParameterModification{}, response.Message, nil
	}

	// Convert changes to parameter modifications
	modifications, err := convertChangesToModifications(response.Changes, collection, selectedFX)
	if err != nil {
		logger.Error("Error converting changes: %v", err)
		return nil, "", fmt.Errorf("failed to process LLM suggestions: %v", err)
	}

	logger.Info("Successfully parsed response with %d suggestions", len(modifications))
	return modifications, response.Message, nil
}

// convertChangesToModifications converts the LLM's change entries to parameter modifications
func convertChangesToModifications(changes []ChangeEntry, collection fx.TrackCollection, selectedFX map[int][]int) ([]fx.ParameterModification, error) {
	modifications := make([]fx.ParameterModification, 0, len(changes))

	// Debug the incoming changes
	for i, change := range changes {
		logger.Debug("Processing change %d: track=%d, fx=%d, param=%d, new_value=%.4f, new_formatted=%s",
			i, change.TrackIndex, change.FXIndex, change.ParamIndex, change.NewValue, change.NewFormatted)
	}

	// Process each change
	for _, change := range changes {
		// Validate track index
		trackFound := false
		for _, track := range collection.Tracks {
			if track.TrackIndex == change.TrackIndex {
				trackFound = true
				break
			}
		}
		if !trackFound {
			logger.Warning("Track index %d not found in collection, skipping suggestion", change.TrackIndex)
			continue
		}

		// Validate track selection
		trackFXIndices, hasTrackSelection := selectedFX[change.TrackIndex]
		if !hasTrackSelection {
			logger.Warning("Track index %d not in selection, skipping suggestion", change.TrackIndex)
			continue
		}

		// Validate FX index
		fxSelected := false
		for _, idx := range trackFXIndices {
			if idx == change.FXIndex {
				fxSelected = true
				break
			}
		}
		if !fxSelected {
			logger.Warning("FX index %d not selected for track %d, skipping suggestion",
				change.FXIndex, change.TrackIndex)
			continue
		}

		// Find FX in collection to verify parameter and get accurate track/fx information
		var trackName string
		var fxName string
		var originalValue float64
		var originalFormatted string

		paramFound := false
		for _, track := range collection.Tracks {
			if track.TrackIndex == change.TrackIndex {
				trackName = track.TrackName
				for _, fx := range track.FXList {
					if fx.FXIndex == change.FXIndex {
						fxName = fx.FXName
						for _, param := range fx.Parameters {
							if param.ParamIndex == change.ParamIndex {
								paramFound = true
								originalValue = param.Value
								originalFormatted = param.FormattedValue
								break
							}
						}
						break
					}
				}
				break
			}
		}
		if !paramFound {
			logger.Warning("Parameter index %d not found for FX %d on track %d, skipping suggestion",
				change.ParamIndex, change.FXIndex, change.TrackIndex)
			continue
		}

		// Validate and clamp new value to 0.0-1.0 range
		newValue := change.NewValue
		if newValue < 0.0 {
			logger.Warning("Value %.4f out of range, clamping to 0.0", newValue)
			newValue = 0.0
		} else if newValue > 1.0 {
			logger.Warning("Value %.4f out of range, clamping to 1.0", newValue)
			newValue = 1.0
		}

		// Create parameter modification
		modification := fx.ParameterModification{
			TrackIndex:        change.TrackIndex,
			TrackName:         trackName, // Use the name from the collection
			FXIndex:           change.FXIndex,
			FXName:            fxName, // Use the name from the collection
			ParamIndex:        change.ParamIndex,
			ParamName:         change.ParamName,
			OriginalValue:     originalValue,     // Use value from the collection
			NewValue:          newValue,          // This comes from the LLM
			OriginalFormatted: originalFormatted, // Use formatted value from the collection
			NewFormatted:      change.NewFormatted,
			Explanation:       change.Explanation,
		}

		logger.Debug("Created modification: track=%d (%s), fx=%d (%s), param=%d (%s), value=%.4f→%.4f (delta: %.6f), formatted=%s→%s",
			modification.TrackIndex, modification.TrackName,
			modification.FXIndex, modification.FXName,
			modification.ParamIndex, modification.ParamName,
			modification.OriginalValue, modification.NewValue,
			modification.NewValue-modification.OriginalValue,
			modification.OriginalFormatted, modification.NewFormatted)

		modifications = append(modifications, modification)
	}

	// Further validate modifications against the collection
	validModifications := fx.ValidateParameterModifications(collection, modifications)

	logger.Debug("Converted %d changes to %d valid modifications",
		len(changes), len(validModifications))

	return validModifications, nil
}
