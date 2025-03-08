package actions

import (
	"encoding/json"
	"fmt"
	"go-reaper/llm"
	"go-reaper/reaper"
	"strconv"
	"strings"
	"unsafe"
)

// ParameterSuggestion contains a suggestion for a single parameter adjustment
type ParameterSuggestion struct {
	FXIndex     int     `json:"fx_index"`
	ParamIndex  int     `json:"param_index"`
	ParamName   string  `json:"param_name"`
	Value       float64 `json:"value"`
	Explanation string  `json:"explanation"`
}

// AssistantResponse contains the structured response from the LLM
type AssistantResponse struct {
	Suggestions []ParameterSuggestion `json:"suggestions"`
	Reasoning   string                `json:"reasoning"`
}

// RegisterFXAssistant registers the LLM FX Assistant action
func RegisterFXAssistant() error {
	actionID, err := reaper.RegisterMainAction("GO_FX_ASSISTANT", "Go: LLM FX Assistant")
	if err != nil {
		return fmt.Errorf("failed to register LLM FX Assistant: %v", err)
	}

	reaper.ConsoleLog(fmt.Sprintf("LLM FX Assistant registered with ID: %d", actionID))
	reaper.SetActionHandler("GO_FX_ASSISTANT", handleFXAssistant)
	return nil
}

// handleFXAssistant handles the FX Assistant action
func handleFXAssistant() {
	// Lock the current goroutine to the OS thread
	// This ensures we stay on the same thread as REAPER expects
	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()

	reaper.ConsoleLog("----- LLM FX Assistant Activated -----")

	// Get track info
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		// Handle the case where there's no track selected
		if strings.Contains(err.Error(), "no track selected") {
			reaper.MessageBox("Please select a track before using the LLM FX Assistant.", "LLM FX Assistant")
			reaper.ConsoleLog("No track selected. Please select a track before using the LLM FX Assistant.")
		} else {
			// Handle other errors
			reaper.MessageBox(fmt.Sprintf("Error: %v", err), "LLM FX Assistant")
			reaper.ConsoleLog(fmt.Sprintf("Error getting track info: %v", err))
		}
		return
	}

	// Check if track has FX
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to the track before using the LLM FX Assistant.", "LLM FX Assistant")
		reaper.ConsoleLog("Selected track has no FX. Please add FX to the track before using the LLM FX Assistant.")
		return
	}

	// Get FX list
	fxList, err := reaper.GetTrackFXList(trackInfo.MediaTrack)
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Error getting FX list: %v", err))
		return
	}

	reaper.ConsoleLog(fmt.Sprintf("Found %d FX on track.", len(fxList)))

	// Build FX selection dialog
	fxOptions := buildFXSelectionList(fxList)
	reaper.ConsoleLog(fxOptions)

	// Show FX selection dialog
	fields := []string{
		"FX to adjust (comma-separated numbers)",
		"Your request (e.g., 'make vocals clearer')",
	}

	defaults := []string{
		"1", // Default to first FX
		"",  // Empty prompt
	}

	// Show the dialog
	results, err := reaper.GetUserInputs("LLM FX Assistant", fields, defaults)
	if err != nil {
		reaper.ConsoleLog("User cancelled the dialog")
		return
	}

	// Parse the results
	selectedFXIndices, err := parseFXSelection(results[0], len(fxList))
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Invalid FX selection: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Invalid FX selection: %v", err))
		return
	}

	userPrompt := results[1]
	if userPrompt == "" {
		reaper.MessageBox("Please provide a request for the LLM FX Assistant.", "LLM FX Assistant")
		reaper.ConsoleLog("Empty prompt provided")
		return
	}

	// Log the selections
	reaper.ConsoleLog(fmt.Sprintf("Selected FX indices: %v", selectedFXIndices))
	reaper.ConsoleLog(fmt.Sprintf("User prompt: %s", userPrompt))

	// Collect FX parameters for the selected FX
	fxParameters := collectFXParameters(trackInfo.MediaTrack, selectedFXIndices, fxList)

	// Show the parameters (for debugging/validation)
	parametersText := formatFXParametersText(fxParameters)
	reaper.ConsoleLog("Parameters collected:")
	reaper.ConsoleLog(parametersText)

	// Ask if the user wants to proceed with LLM analysis
	confirmMsg := fmt.Sprintf("Track: %s\nFX selected: %d\nRequest: %s\n\nReady to analyze with LLM?\n\nNote: This will require an OpenAI API key.",
		trackInfo.Name, len(selectedFXIndices), userPrompt)

	proceed, err := reaper.YesNoBox(confirmMsg, "LLM FX Assistant")
	if err != nil || !proceed {
		reaper.ConsoleLog("User chose not to proceed with LLM analysis")
		return
	}

	// Get API key
	apiKey, err := getOpenAIKey()
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("API key error: %v", err))
		return
	}

	// Show processing message
	reaper.MessageBox("Analyzing parameters with OpenAI...\nThis might take a few seconds.", "LLM FX Assistant")

	// Create LLM client
	client := llm.NewOpenAIClient(apiKey)

	// Build prompts
	systemPrompt := buildSystemPrompt()
	userPromptText := buildUserPrompt(fxParameters, userPrompt)

	reaper.ConsoleLog("System Prompt:")
	reaper.ConsoleLog(systemPrompt)
	reaper.ConsoleLog("User Prompt:")
	reaper.ConsoleLog(userPromptText)

	// Call LLM
	responseText, err := client.SendPrompt(systemPrompt, userPromptText)
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error calling LLM API: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Error calling LLM API: %v", err))
		return
	}

	reaper.ConsoleLog("LLM Response:")
	reaper.ConsoleLog(responseText)

	// Parse the response
	assistantResponse, err := parseAssistantResponse(responseText)
	if err != nil {
		reaper.MessageBox(fmt.Sprintf("Error parsing LLM response: %v", err), "LLM FX Assistant")
		reaper.ConsoleLog(fmt.Sprintf("Error parsing LLM response: %v", err))
		return
	}

	// Format the results for display
	resultsText := formatAssistantResults(assistantResponse)

	// Ask the user if they want to apply the changes
	applyMsg := fmt.Sprintf("The LLM suggests these parameter changes:\n\n%s\n\nWould you like to apply these changes?", resultsText)
	apply, err := reaper.YesNoBox(applyMsg, "LLM FX Assistant - Apply Changes")
	if err != nil {
		reaper.ConsoleLog(fmt.Sprintf("Dialog error: %v", err))
		return
	}

	if apply {
		// Apply the changes
		err = applyParameterChanges(trackInfo.MediaTrack, assistantResponse.Suggestions)
		if err != nil {
			reaper.MessageBox(fmt.Sprintf("Error applying changes: %v", err), "LLM FX Assistant")
			reaper.ConsoleLog(fmt.Sprintf("Error applying changes: %v", err))
			return
		}

		reaper.MessageBox("Parameter changes applied successfully!", "LLM FX Assistant")
		reaper.ConsoleLog("Parameter changes applied successfully")
	} else {
		reaper.ConsoleLog("User chose not to apply changes")
	}
}

// buildSystemPrompt creates a system prompt for the LLM that explains the task and format
func buildSystemPrompt() string {
	return `You are an audio engineer assistant that helps adjust effects (FX) parameters in a digital audio workstation.
You will be given information about one or more audio effects including their names and parameters.
You will also receive a user request about how they want to adjust the sound.

Your task is to suggest parameter adjustments that will help achieve the user's request.

IMPORTANT RULES:
1. Only suggest adjustments to the parameters provided.
2. Always return values within the normalized range (0.0 to 1.0).
3. Always format your response as valid JSON with this structure:
{
  "suggestions": [
    {
      "fx_index": <integer index of the effect>,
      "param_index": <integer index of the parameter>,
      "param_name": "<name of the parameter>",
      "value": <new value between 0.0 and 1.0>,
      "explanation": "<brief explanation of this adjustment>"
    }
  ],
  "reasoning": "<your overall explanation of the parameter adjustments>"
}

4. Keep explanations concise but technically accurate.
5. Only include parameters you are adjusting in the suggestions array.
6. Focus on achieving the user's sonic goals with the minimum necessary adjustments.
7. The JSON must be valid and complete.`
}

// buildUserPrompt creates a prompt with FX details and the user's request
func buildUserPrompt(fxList []reaper.FXInfo, userRequest string) string {
	var builder strings.Builder

	builder.WriteString("Here are the audio effects and their current parameters:\n\n")

	for _, fx := range fxList {
		builder.WriteString(fmt.Sprintf("FX %d: %s\n", fx.Index, fx.Name))
		builder.WriteString("Parameters:\n")

		for _, param := range fx.Parameters {
			builder.WriteString(fmt.Sprintf("  - %s (index: %d): %.4f (formatted: %s)\n",
				param.Name, param.Index, param.Value, param.FormattedValue))
		}

		builder.WriteString("\n")
	}

	builder.WriteString("User request: " + userRequest + "\n\n")
	builder.WriteString("Please suggest parameter adjustments that will help achieve this request. Remember to format your response as JSON according to the specified structure.")

	return builder.String()
}

// parseAssistantResponse parses the LLM's text response into a structured AssistantResponse
// parseAssistantResponse parses the LLM's text response into a structured AssistantResponse
func parseAssistantResponse(responseText string) (*AssistantResponse, error) {
	// Validate input
	if responseText == "" {
		return nil, fmt.Errorf("empty response text from LLM")
	}

	reaper.ConsoleLog(fmt.Sprintf("Parsing response text (%d chars)...", len(responseText)))

	// Try to extract JSON from the response (it might contain additional text)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		reaper.ConsoleLog("Failed to find valid JSON markers in response")
		reaper.ConsoleLog(fmt.Sprintf("Response text: %s", responseText))
		return nil, fmt.Errorf("could not find valid JSON in response")
	}

	jsonStr := responseText[jsonStart : jsonEnd+1]
	reaper.ConsoleLog(fmt.Sprintf("Extracted JSON (%d chars)", len(jsonStr)))

	// Parse the JSON response
	var response AssistantResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		reaper.ConsoleLog(fmt.Sprintf("JSON unmarshal error: %v", err))
		// Log a truncated version of the JSON for debugging
		if len(jsonStr) > 200 {
			reaper.ConsoleLog(fmt.Sprintf("JSON (truncated): %s...", jsonStr[:200]))
		} else {
			reaper.ConsoleLog(fmt.Sprintf("JSON: %s", jsonStr))
		}
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	// Validate the response structure
	if response.Suggestions == nil {
		return nil, fmt.Errorf("response is missing the suggestions array")
	}

	// Validate the response
	if len(response.Suggestions) == 0 {
		return nil, fmt.Errorf("no parameter suggestions found in response")
	}

	// Validate parameter values
	for i, suggestion := range response.Suggestions {
		// Validate FX index is present
		if suggestion.FXIndex < 0 {
			reaper.ConsoleLog(fmt.Sprintf("Warning: Invalid FX index %d, using 0", suggestion.FXIndex))
			response.Suggestions[i].FXIndex = 0
		}

		// Validate parameter value is in range
		if suggestion.Value < 0 || suggestion.Value > 1 {
			reaper.ConsoleLog(fmt.Sprintf("Warning: Parameter value %f outside 0-1 range, clamping", suggestion.Value))
			if suggestion.Value < 0 {
				response.Suggestions[i].Value = 0
			} else {
				response.Suggestions[i].Value = 1
			}
		}
	}

	reaper.ConsoleLog(fmt.Sprintf("Successfully parsed response with %d suggestions", len(response.Suggestions)))
	return &response, nil
}

// buildFXSelectionList creates a formatted list of FX for display
func buildFXSelectionList(fxList []reaper.FXInfo) string {
	var builder strings.Builder
	builder.WriteString("\nAvailable FX:\n")

	for i, fx := range fxList {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, fx.Name))
	}

	return builder.String()
}

// parseFXSelection parses a comma-separated list of FX indices
func parseFXSelection(input string, maxFX int) ([]int, error) {
	if input == "" {
		return nil, fmt.Errorf("no FX selected")
	}

	// Split by comma
	parts := strings.Split(input, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Parse the number
		idx, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid FX number: %s", part)
		}

		// Adjust for 1-based indexing in the UI to 0-based indexing internally
		idx--

		// Check range
		if idx < 0 || idx >= maxFX {
			return nil, fmt.Errorf("FX number out of range: %d", idx+1)
		}

		result = append(result, idx)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid FX selected")
	}

	return result, nil
}

// collectFXParameters collects all parameters for the selected FX
func collectFXParameters(track unsafe.Pointer, indices []int, fxList []reaper.FXInfo) []reaper.FXInfo {
	result := make([]reaper.FXInfo, 0, len(indices))

	for _, fxIndex := range indices {
		// Get full FX parameters
		fxInfo, err := reaper.GetFXParameters(track, fxIndex)
		if err != nil {
			reaper.ConsoleLog(fmt.Sprintf("Error getting FX parameters for %s: %v",
				fxList[fxIndex].Name, err))
			continue
		}

		result = append(result, fxInfo)
	}

	return result
}

// formatFXParametersText formats the FX parameters for display
func formatFXParametersText(fxParameters []reaper.FXInfo) string {
	var builder strings.Builder

	for _, fx := range fxParameters {
		builder.WriteString(fmt.Sprintf("FX: %s\n", fx.Name))
		builder.WriteString("Parameters:\n")

		for _, param := range fx.Parameters {
			builder.WriteString(fmt.Sprintf("  %s: %.2f (%s)\n",
				param.Name, param.Value, param.FormattedValue))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// formatAssistantResults formats the LLM suggestions for display
func formatAssistantResults(response *AssistantResponse) string {
	var builder strings.Builder

	// Add reasoning if provided
	if response.Reasoning != "" {
		builder.WriteString("Analysis: " + response.Reasoning + "\n\n")
	}

	builder.WriteString("Suggested Changes:\n")

	// Group suggestions by FX
	fxGroups := make(map[int][]ParameterSuggestion)
	for _, suggestion := range response.Suggestions {
		fxGroups[suggestion.FXIndex] = append(fxGroups[suggestion.FXIndex], suggestion)
	}

	// Format each FX group
	for fxIndex, suggestions := range fxGroups {
		builder.WriteString(fmt.Sprintf("\nFX %d:\n", fxIndex))

		for _, suggestion := range suggestions {
			builder.WriteString(fmt.Sprintf("  • %s: %.2f\n    %s\n",
				suggestion.ParamName,
				suggestion.Value,
				suggestion.Explanation))
		}
	}

	return builder.String()
}

// applyParameterChanges applies the parameter changes suggested by the LLM
func applyParameterChanges(track unsafe.Pointer, suggestions []ParameterSuggestion) error {
	for _, suggestion := range suggestions {
		// Apply the parameter change
		err := reaper.SetTrackFXParamValue(track, suggestion.FXIndex, suggestion.ParamIndex, suggestion.Value)
		if err != nil {
			return fmt.Errorf("failed to set parameter value: %v", err)
		}

		// Log the change
		reaper.ConsoleLog(fmt.Sprintf("Applied: FX %d, Parameter %d (%s): %.4f - %s",
			suggestion.FXIndex,
			suggestion.ParamIndex,
			suggestion.ParamName,
			suggestion.Value,
			suggestion.Explanation))
	}

	return nil
}

// getOpenAIKey asks the user for their OpenAI API key
func getOpenAIKey() (string, error) {
	fields := []string{"OpenAI API Key"}
	defaults := []string{""}

	values, err := reaper.GetUserInputs("Enter OpenAI API Key", fields, defaults)
	if err != nil {
		return "", err
	}

	apiKey := values[0]
	if apiKey == "" {
		return "", fmt.Errorf("API key is required")
	}

	return apiKey, nil
}
