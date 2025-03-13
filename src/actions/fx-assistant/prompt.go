package fxassistant

import (
	"fmt"
	"go-reaper/src/llm"
	"go-reaper/src/pkg/config"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper/fx"
	"strings"
)

// buildSystemPrompt creates the system prompt for the LLM
func buildSystemPrompt() string {
	return `You are an audio engineer assistant that helps adjust effects (FX) parameters in a digital audio workstation.
You will be given information about one or more audio effects including their names and parameters.
You will also receive a user request about how they want to adjust the sound.

Your task is to suggest parameter adjustments that will help achieve the user's request.

IMPORTANT RULES:
1. Only suggest adjustments to the parameters provided.
2. Always return values within the normalized range (0.0 to 1.0).
3. BE EXTREMELY PRECISE with normalized values. For example, a 1dB change may only require a 0.0033 change in normalized value.
4. Always format your response as valid JSON with this structure:
{
  "message": "<your overall explanation of what you're doing>",
  "changes": [
    {
      "track_index": <integer index of the track>,
      "track_name": "<name of the track>",
      "fx_index": <integer index of the effect>,
      "fx_name": "<name of the effect>",
      "param_index": <integer index of the parameter>,
      "param_name": "<name of the parameter>",
      "original_value": <original normalized value between 0.0 and 1.0>,
      "new_value": <new normalized value between 0.0 and 1.0>,
      "original_formatted": "<human-readable original value>",
      "new_formatted": "<human-readable new value>",
      "explanation": "<brief explanation of this adjustment including the exact normalized value change>"
    }
  ]
}

5. Keep explanations concise but technically accurate. ALWAYS specify the exact normalized value change in the explanation.
6. Parameter values are ALWAYS provided to you in the format: "value: 0.752019 (formatted: 6.8)" where 0.752019 is the normalized value and 6.8 is the display value.
7. Only include parameters you are adjusting in the changes array.
8. Focus on achieving the user's sonic goals with the minimum necessary adjustments.
9. If no parameter changes are needed, return an empty changes array with a helpful message explaining why.
10. The JSON must be valid and complete.`
}

// buildUserPrompt creates the user prompt with FX details and the user's request
func buildUserPrompt(collection fx.TrackCollection, selectedFX map[int][]int, userRequest string) string {
	var builder strings.Builder

	builder.WriteString("Here are the audio effects and their current parameters:\n\n")

	// Include only the selected tracks and FX
	for _, track := range collection.Tracks {
		// Check if this track has any selected FX
		trackFXIndices, hasSelection := selectedFX[track.TrackIndex]
		if !hasSelection {
			continue
		}

		builder.WriteString(fmt.Sprintf("Track %d: %s\n", track.TrackIndex, track.TrackName))

		// Include only the selected FX for this track
		for _, fx := range track.FXList {
			// Check if this FX is selected
			isSelected := false
			for _, selectedFXIndex := range trackFXIndices {
				if fx.FXIndex == selectedFXIndex {
					isSelected = true
					break
				}
			}

			if !isSelected {
				continue
			}

			builder.WriteString(fmt.Sprintf("  FX %d: %s\n", fx.FXIndex, fx.FXName))
			builder.WriteString("  Parameters:\n")

			for _, param := range fx.Parameters {
				builder.WriteString(fmt.Sprintf("    - %s (index: %d): %.4f (formatted: %s)\n",
					param.ParamName, param.ParamIndex, param.Value, param.FormattedValue))

				// Include parameter range information to help the LLM understand the parameter better
				builder.WriteString(fmt.Sprintf("      Range: %.4f to %.4f (formatted: %s to %s)\n",
					param.Min, param.Max, param.MinFormatted, param.MaxFormatted))

				// Add an explicit note about normalized vs formatted values
				builder.WriteString("      NOTE: You must work with both types of values - the normalized values (0.0-1.0) when setting parameters,\n")
				builder.WriteString("      and the formatted values when displaying to the user. Make sure your suggestions include both.\n")
			}

			builder.WriteString("\n")
		}

		builder.WriteString("\n")
	}

	builder.WriteString("User request: " + userRequest + "\n\n")
	builder.WriteString("Please suggest parameter adjustments that will help achieve this request. Remember to format your response as JSON according to the specified structure.")

	return builder.String()
}

// sendPromptToLLM sends the prompts to the LLM service and returns the response
func sendPromptToLLM(systemPrompt, userPrompt string) (string, error) {
	// Get API key from keyring
	provider := config.GetActiveProvider()
	apiKey, err := config.GetSecureAPIKey(provider)
	if err != nil {
		return "", fmt.Errorf("failed to get API key: %v", err)
	}

	// Get provider configuration
	model, maxTokens, temperature := config.GetProviderConfig(provider)
	logger.Debug("Using LLM model: %s with temperature: %.2f", model, temperature)

	// Create LLM client
	client := llm.NewOpenAIClient(apiKey)
	client.Model = model
	client.MaxTokens = maxTokens
	client.Temp = temperature

	// Log the prompts
	logger.Debug("System prompt length: %d characters", len(systemPrompt))
	logger.Debug("User prompt length: %d characters", len(userPrompt))

	// Send the prompt to the LLM
	logger.Info("Sending prompt to LLM...")
	responseText, err := client.SendPrompt(systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("error calling LLM API: %v", err)
	}

	logger.Debug("Received response from LLM (%d characters)", len(responseText))
	return responseText, nil
}

// processRequestWithLLM processes the user's request using the LLM
// This function is called from assistant.go - it coordinates the prompting and response parsing
func processRequestWithLLM(collection fx.TrackCollection, selectedFX map[int][]int, userRequest string) ([]fx.ParameterModification, string, error) {
	logger.Debug("Processing request with LLM: %s", userRequest)

	// Build the prompts
	systemPrompt := buildSystemPrompt()
	userPrompt := buildUserPrompt(collection, selectedFX, userRequest)

	// Send to the LLM
	responseText, err := sendPromptToLLM(systemPrompt, userPrompt)
	if err != nil {
		return nil, "", fmt.Errorf("error sending prompt to LLM: %v", err)
	}

	// Parse the response (this function is defined in response.go)
	return parseAssistantResponse(responseText, collection, selectedFX)
}
