package analyzer

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Parameter type constants for classification
const (
	ParamTypeBinary      = "BINARY"
	ParamTypeEnumerated  = "ENUMERATED"
	ParamTypeLinear      = "LINEAR"
	ParamTypeLogarithmic = "LOGARITHMIC"
	ParamTypeExponential = "EXPONENTIAL"
	ParamTypeInverted    = "INVERTED"
	ParamTypeUnknown     = "UNKNOWN"

	// Constants for statistics files
	LogPrefix = "Parameter Analysis"
)

// ParameterSample represents a sample of a parameter at a specific normalized value
type ParameterSample struct {
	NormalizedValue float64
	FormattedValue  string
	NumericValue    float64
	IsNumeric       bool
}

// ParameterAnalysis contains the analysis result for a single parameter
type ParameterAnalysis struct {
	FXIndex          int
	FXName           string
	ParamIndex       int
	ParamName        string
	DetectedType     string
	Confidence       float64
	Samples          []ParameterSample
	CurrentValue     float64
	CurrentFormatted string
	Min              float64
	Max              float64
	MinFormatted     string
	MaxFormatted     string
}

// RegisterParameterAnalyzer registers the parameter analyzer action
func RegisterParameterAnalyzer() error {
	actionID, err := reaper.RegisterMainAction("GO_PARAM_ANALYZER", "Go: Parameter Analyzer")
	if err != nil {
		return fmt.Errorf("failed to register parameter analyzer action: %v", err)
	}

	logger.Info("Parameter Analyzer registered with ID: %d", actionID)
	reaper.SetActionHandler("GO_PARAM_ANALYZER", handleParameterAnalyzer)
	return nil
}

// handleParameterAnalyzer runs the parameter analyzer
func handleParameterAnalyzer() {
	// Lock the current goroutine to the OS thread for UI operations
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	logger.Info("Parameter Analyzer action triggered")

	// Get selected track
	track, err := reaper.GetSelectedTrack()
	if err != nil {
		reaper.MessageBox("Please select a track with FX to analyze", "Parameter Analyzer")
		return
	}

	// Get track info
	trackInfo, err := reaper.GetSelectedTrackInfo()
	if err != nil {
		reaper.MessageBox("Error getting track info", "Parameter Analyzer")
		return
	}

	// Check if track has FX
	if trackInfo.NumFX == 0 {
		reaper.MessageBox("Selected track has no FX. Please add FX to analyze.", "Parameter Analyzer")
		return
	}

	// Show confirmation
	confirm, err := reaper.YesNoBox(
		fmt.Sprintf("This will analyze all parameters on %d FX on track '%s'.\n\nResults will be written to the log file and console.\nThis may take some time for complex plugins.\n\nProceed?",
			trackInfo.NumFX, trackInfo.Name),
		"Parameter Analyzer")

	if err != nil || !confirm {
		logger.Info("User cancelled parameter analysis")
		return
	}

	// Start timing
	startTime := time.Now()

	// Create a stats counter
	stats := make(map[string]int)
	var allAnalyses []ParameterAnalysis

	// Log header
	logger.Info("=====================================================")
	logger.Info("%s - TRACK: %s", LogPrefix, trackInfo.Name)
	logger.Info("=====================================================")

	// For each FX on the track
	for fxIndex := 0; fxIndex < trackInfo.NumFX; fxIndex++ {
		// Get FX info
		fxName, err := reaper.GetTrackFXName(track, fxIndex)
		if err != nil {
			logger.Error("Failed to get FX name for index %d: %v", fxIndex, err)
			continue
		}

		// Get parameter count
		paramCount, err := reaper.GetTrackFXParamCount(track, fxIndex)
		if err != nil {
			logger.Error("Failed to get parameter count for FX #%d: %v", fxIndex+1, err)
			continue
		}

		// Log FX info
		logger.Info("FX #%d: %s", fxIndex+1, fxName)
		logger.Info("  Parameter count: %d", paramCount)

		// For each parameter
		for paramIndex := 0; paramIndex < paramCount; paramIndex++ {
			// Get parameter name
			paramName, err := reaper.GetTrackFXParamName(track, fxIndex, paramIndex)
			if err != nil {
				logger.Error("Failed to get parameter name for index %d: %v", paramIndex, err)
				continue
			}

			// Get parameter range and current value
			currentValue, min, max, err := reaper.GetTrackFXParamValueWithRange(track, fxIndex, paramIndex)
			if err != nil {
				logger.Error("Failed to get parameter range: %v", err)
				continue
			}

			// Get current formatted value
			currentFormatted, err := reaper.GetTrackFXParamFormatted(track, fxIndex, paramIndex)
			if err != nil {
				logger.Error("Failed to get current formatted value: %v", err)
				continue
			}

			// Get min/max formatted values - use direct API to avoid batch issues
			minFormatted, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, min)
			if err != nil {
				logger.Error("Failed to get min formatted value: %v", err)
				minFormatted = ""
			}

			maxFormatted, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, max)
			if err != nil {
				logger.Error("Failed to get max formatted value: %v", err)
				maxFormatted = ""
			}

			// Initialize analysis struct
			analysis := ParameterAnalysis{
				FXIndex:          fxIndex,
				FXName:           fxName,
				ParamIndex:       paramIndex,
				ParamName:        paramName,
				DetectedType:     ParamTypeUnknown,
				Confidence:       0.0,
				CurrentValue:     currentValue,
				CurrentFormatted: currentFormatted,
				Min:              min,
				Max:              max,
				MinFormatted:     minFormatted,
				MaxFormatted:     maxFormatted,
			}

			// Define sample points - more at lower end for better logarithmic detection
			samplePoints := []float64{0.0, 0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.95, 0.99, 1.0}

			// Create samples array
			samples := make([]ParameterSample, len(samplePoints))
			for i, point := range samplePoints {
				// Use direct API call instead of batch
				formattedValue, err := reaper.GetTrackFXParamFormattedValueWithValue(track, fxIndex, paramIndex, point)
				if err != nil {
					logger.Warning("Failed to get formatted value for point %.2f: %v", point, err)
					formattedValue = ""
				}

				samples[i] = ParameterSample{
					NormalizedValue: point,
					FormattedValue:  formattedValue,
				}

				// Try to extract numeric value
				numericValue, isNumeric := extractNumericValue(formattedValue)
				samples[i].NumericValue = numericValue
				samples[i].IsNumeric = isNumeric
			}

			analysis.Samples = samples

			// Analyze parameter type
			analysis.DetectedType, analysis.Confidence = classifyParameter(samples)

			// Log the analysis result
			logger.Info("  Parameter #%d: %s", paramIndex+1, paramName)
			logger.Info("    Detected type: %s (confidence: %.2f)", analysis.DetectedType, analysis.Confidence)
			logger.Info("    Current value: %.4f (%s)", analysis.CurrentValue, analysis.CurrentFormatted)
			logger.Info("    Min: %.4f (%s), Max: %.4f (%s)",
				analysis.Min, analysis.MinFormatted, analysis.Max, analysis.MaxFormatted)

			// Log sample points
			logText := "    Sample points: ["
			for i, sample := range analysis.Samples {
				if i > 0 {
					logText += ", "
				}
				logText += fmt.Sprintf("%.2f", sample.NormalizedValue)
			}
			logText += "]"
			logger.Info(logText)

			// Log formatted values
			logText = "    Formatted values: ["
			for i, sample := range analysis.Samples {
				if i > 0 {
					logText += ", "
				}
				logText += fmt.Sprintf("\"%s\"", sample.FormattedValue)
			}
			logText += "]"
			logger.Info(logText)

			// Log numeric values if available
			if numericCount(analysis.Samples) > 0 {
				logText = "    Numeric values: ["
				hasValues := false
				for _, sample := range analysis.Samples {
					if sample.IsNumeric {
						if hasValues {
							logText += ", "
						}
						logText += fmt.Sprintf("%.4f", sample.NumericValue)
						hasValues = true
					} else if hasValues {
						logText += ", null"
					}
				}
				logText += "]"
				logger.Info(logText)
			}

			// Update stats
			stats[analysis.DetectedType]++

			// Add to all analyses
			allAnalyses = append(allAnalyses, analysis)
		}

		logger.Info("-----------------------------------------------------")
	}

	// Calculate duration
	duration := time.Since(startTime)

	// Log stats
	logger.Info("=====================================================")
	logger.Info("%s - STATISTICS", LogPrefix)
	logger.Info("=====================================================")
	logger.Info("Total parameters analyzed: %d", len(allAnalyses))
	logger.Info("Analysis duration: %v", duration)
	logger.Info("Parameter type distribution:")

	// Get sorted stats keys
	var types []string
	for paramType := range stats {
		types = append(types, paramType)
	}
	sort.Slice(types, func(i, j int) bool {
		return stats[types[i]] > stats[types[j]] // Sort by count descending
	})

	// Log stats in order
	for _, paramType := range types {
		count := stats[paramType]
		percent := float64(count) / float64(len(allAnalyses)) * 100
		logger.Info("  %s: %d (%.1f%%)", paramType, count, percent)
	}

	// Prepare console report
	var report strings.Builder
	report.WriteString(fmt.Sprintf("Analyzed %d parameters across %d FX plugins on track '%s'\n\n",
		len(allAnalyses), trackInfo.NumFX, trackInfo.Name))
	report.WriteString("Parameter type distribution:\n")

	for _, paramType := range types {
		count := stats[paramType]
		percent := float64(count) / float64(len(allAnalyses)) * 100
		report.WriteString(fmt.Sprintf("  %s: %d (%.1f%%)\n", paramType, count, percent))
	}

	// Show report in REAPER console
	reaper.ShowConsoleMsg(report.String())

	// Show completion message
	reaper.MessageBox(
		fmt.Sprintf("Analysis complete! Analyzed %d parameters in %v.\n\nSee REAPER console and log file for detailed results.",
			len(allAnalyses), duration.Round(time.Millisecond)),
		"Parameter Analyzer")
}

// classifyParameter determines the parameter type based on samples
func classifyParameter(samples []ParameterSample) (string, float64) {
	// 1. Check for binary parameters (only 2 distinct values)
	uniqueValues := make(map[string]bool)
	for _, sample := range samples {
		uniqueValues[sample.FormattedValue] = true
	}

	if len(uniqueValues) <= 2 {
		return ParamTypeBinary, 0.95
	}

	// 2. Check for enumerated parameters (few distinct values)
	if len(uniqueValues) <= 10 && len(uniqueValues) < len(samples)/2 {
		return ParamTypeEnumerated, 0.90
	}

	// 3. For numeric parameters, analyze relationships
	numericSampleCount := numericCount(samples)
	if numericSampleCount > len(samples)/2 {
		// Calculate confidences for each type
		linearConfidence := detectLinearRelationship(samples)
		logConfidence := detectLogarithmicRelationship(samples)
		expConfidence := detectExponentialRelationship(samples)
		invConfidence := detectInvertedRelationship(samples)

		// Find highest confidence
		bestType := ParamTypeUnknown
		bestConfidence := 0.5 // Minimum threshold

		if linearConfidence > bestConfidence && linearConfidence >= logConfidence &&
			linearConfidence >= expConfidence && linearConfidence >= invConfidence {
			bestType = ParamTypeLinear
			bestConfidence = linearConfidence
		} else if logConfidence > bestConfidence && logConfidence >= linearConfidence &&
			logConfidence >= expConfidence && logConfidence >= invConfidence {
			bestType = ParamTypeLogarithmic
			bestConfidence = logConfidence
		} else if expConfidence > bestConfidence && expConfidence >= linearConfidence &&
			expConfidence >= logConfidence && expConfidence >= invConfidence {
			bestType = ParamTypeExponential
			bestConfidence = expConfidence
		} else if invConfidence > bestConfidence && invConfidence >= linearConfidence &&
			invConfidence >= logConfidence && invConfidence >= expConfidence {
			bestType = ParamTypeInverted
			bestConfidence = invConfidence
		}

		if bestType != ParamTypeUnknown {
			return bestType, bestConfidence
		}
	}

	// 4. Attempt to detect special cases based on parameter names
	lowerName := strings.ToLower(samples[0].FormattedValue)

	// Common frequency parameters
	if strings.Contains(lowerName, "hz") || strings.Contains(lowerName, "khz") {
		return ParamTypeLogarithmic, 0.7
	}

	// Common time parameters
	if strings.Contains(lowerName, "ms") || strings.Contains(lowerName, "sec") {
		return ParamTypeLogarithmic, 0.7
	}

	// Common dB parameters
	if strings.Contains(lowerName, "db") {
		return ParamTypeLinear, 0.7
	}

	return ParamTypeUnknown, 0.0
}

// extractNumericValue attempts to extract a numeric value from a formatted string
func extractNumericValue(formatted string) (float64, bool) {
	if formatted == "" {
		return 0, false
	}

	// Try direct conversion first
	value, err := strconv.ParseFloat(formatted, 64)
	if err == nil {
		return value, true
	}

	// Extract first number from the string
	var numBuilder strings.Builder
	seenDecimalPoint := false
	seenDigit := false

	for _, char := range formatted {
		if unicode.IsDigit(char) {
			numBuilder.WriteRune(char)
			seenDigit = true
		} else if char == '.' && !seenDecimalPoint && seenDigit {
			numBuilder.WriteRune(char)
			seenDecimalPoint = true
		} else if seenDigit && !unicode.IsDigit(char) && char != '.' {
			// Stop at first non-digit after seeing digits
			break
		}
	}

	if seenDigit {
		value, err := strconv.ParseFloat(numBuilder.String(), 64)
		if err == nil {
			return value, true
		}
	}

	return 0, false
}

// numericCount returns the number of samples with valid numeric values
func numericCount(samples []ParameterSample) int {
	count := 0
	for _, sample := range samples {
		if sample.IsNumeric {
			count++
		}
	}
	return count
}

// detectLinearRelationship detects if there's a linear relationship between samples
func detectLinearRelationship(samples []ParameterSample) float64 {
	if len(samples) < 3 {
		return 0.0
	}

	// Count valid pairs for analysis
	var validPairs []struct {
		x, y float64
	}

	for i := 1; i < len(samples); i++ {
		if samples[i].IsNumeric && samples[i-1].IsNumeric {
			validPairs = append(validPairs, struct {
				x, y float64
			}{
				x: samples[i].NormalizedValue - samples[i-1].NormalizedValue,
				y: samples[i].NumericValue - samples[i-1].NumericValue,
			})
		}
	}

	if len(validPairs) < 2 {
		return 0.0
	}

	// Calculate average rate of change
	totalRate := 0.0
	rateCount := 0

	for _, pair := range validPairs {
		if pair.x != 0 {
			totalRate += pair.y / pair.x
			rateCount++
		}
	}

	if rateCount == 0 {
		return 0.0
	}

	avgRate := totalRate / float64(rateCount)

	// Calculate variance of rates from average
	variance := 0.0
	for _, pair := range validPairs {
		if pair.x != 0 {
			rate := pair.y / pair.x
			diff := rate - avgRate
			variance += diff * diff
		}
	}

	variance /= float64(rateCount)

	// Convert variance to confidence score (inverse relationship)
	// Lower variance = higher confidence
	maxVariance := avgRate * avgRate                        // Normalize by avgRate squared
	normalizedVariance := variance / (maxVariance + 0.0001) // Avoid division by zero

	confidence := 1.0 - math.Min(normalizedVariance, 1.0)
	return confidence
}

// detectLogarithmicRelationship detects if there's a logarithmic relationship
func detectLogarithmicRelationship(samples []ParameterSample) float64 {
	// Only analyze if we have enough numeric samples
	numCount := 0
	numericIndices := []int{}

	for i, sample := range samples {
		if sample.IsNumeric {
			numCount++
			numericIndices = append(numericIndices, i)
		}
	}

	if numCount < 4 {
		return 0.0
	}

	// Simple heuristic: Check if changes are larger at lower end of range
	// For logarithmic function, the first 25% of input range might produce
	// 50% or more of the output range

	firstQuarterIdx := -1
	midpointIdx := -1

	// Find indices closest to 25% and 50% of the normalized range
	for _, idx := range numericIndices {
		if firstQuarterIdx == -1 && samples[idx].NormalizedValue >= 0.25 {
			firstQuarterIdx = idx
		}
		if midpointIdx == -1 && samples[idx].NormalizedValue >= 0.5 {
			midpointIdx = idx
			break
		}
	}

	// If we couldn't find appropriate points, return low confidence
	if firstQuarterIdx == -1 || midpointIdx == -1 ||
		!samples[0].IsNumeric || !samples[len(samples)-1].IsNumeric {
		return 0.0
	}

	// Calculate value ranges
	firstQuarterRange := math.Abs(samples[firstQuarterIdx].NumericValue - samples[0].NumericValue)
	totalRange := math.Abs(samples[len(samples)-1].NumericValue - samples[0].NumericValue)

	if totalRange == 0 {
		return 0.0
	}

	// If first quarter produces more than 40% of total change, likely logarithmic
	ratio := firstQuarterRange / totalRange
	if ratio > 0.4 {
		confidence := (ratio - 0.4) * 2.0 // Scale to 0-1 range
		return math.Min(confidence, 0.9)  // Cap at 0.9
	}

	return 0.0
}

// detectExponentialRelationship detects if there's an exponential relationship
func detectExponentialRelationship(samples []ParameterSample) float64 {
	// Only analyze if we have enough numeric samples
	numCount := 0
	numericIndices := []int{}

	for i, sample := range samples {
		if sample.IsNumeric {
			numCount++
			numericIndices = append(numericIndices, i)
		}
	}

	if numCount < 4 {
		return 0.0
	}

	// Simple heuristic: Check if changes are larger at upper end of range
	// For exponential function, the last 25% of input range might produce
	// 50% or more of the output range

	threeQuarterIdx := -1

	// Find indices closest to 75% of the normalized range
	for i := len(numericIndices) - 1; i >= 0; i-- {
		idx := numericIndices[i]
		if samples[idx].NormalizedValue <= 0.75 {
			threeQuarterIdx = idx
			break
		}
	}

	// If we couldn't find appropriate points, return low confidence
	if threeQuarterIdx == -1 ||
		!samples[0].IsNumeric || !samples[len(samples)-1].IsNumeric {
		return 0.0
	}

	// Calculate value ranges
	lastQuarterRange := math.Abs(samples[len(samples)-1].NumericValue - samples[threeQuarterIdx].NumericValue)
	totalRange := math.Abs(samples[len(samples)-1].NumericValue - samples[0].NumericValue)

	if totalRange == 0 {
		return 0.0
	}

	// If last quarter produces more than 40% of total change, likely exponential
	ratio := lastQuarterRange / totalRange
	if ratio > 0.4 {
		confidence := (ratio - 0.4) * 2.0 // Scale to 0-1 range
		return math.Min(confidence, 0.9)  // Cap at 0.9
	}

	return 0.0
}

// detectInvertedRelationship detects if values decrease as normalized values increase
func detectInvertedRelationship(samples []ParameterSample) float64 {
	// Verify we have enough numeric samples
	if !samples[0].IsNumeric || !samples[len(samples)-1].IsNumeric {
		return 0.0
	}

	// Check if the end value is less than the start value
	startValue := samples[0].NumericValue
	endValue := samples[len(samples)-1].NumericValue

	// For an inverted parameter, values should decrease as normalized values increase
	if endValue < startValue {
		// Calculate how significant the inversion is
		diff := math.Abs(startValue - endValue)
		max := math.Max(math.Abs(startValue), math.Abs(endValue))

		if max > 0 {
			ratio := diff / max
			return math.Min(ratio, 0.9) // Cap at 0.9
		}
	}

	return 0.0
}
