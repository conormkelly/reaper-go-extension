# FX Parameter Analysis & Classification Strategy

## Overview

This document outlines the strategy for analyzing, classifying, and semantically understanding FX parameters to improve the REAPER LLM FX Assistant. The goal is to provide the LLM with human-readable, contextually rich parameter information while handling the technical complexities of parameter normalization behind the scenes.

## Discoveries & Challenges

Our initial investigation revealed several key challenges in working with plugin parameters:

1. **Inconsistent Metadata Across Plugin Formats**:
   - VST3 plugins often provide accurate parameter metadata (toggle flags, step sizes)
   - AU plugins frequently report all parameters as continuous, even binary toggles
   - REAPER's native plugins report all parameters as continuous (zero step sizes)

2. **Parameter Type Ambiguity**:
   - Parameters with binary behavior don't always report as toggles
   - Enum-like parameters (discrete options) often appear as continuous ranges
   - Continuous parameters may use linear, logarithmic, or exponential scaling

3. **Format-Specific Behaviors**:
   - The same plugin in different formats (VST vs AU) may report parameters differently
   - The same parameter may have different behavior across plugin formats

4. **Value Representation Complexity**:
   - Formatted values may use multiple units within the same parameter (Hz to kHz)
   - REAPER's API provides normalized values (0.0-1.0) that need translation
   - The relationship between normalized and formatted values isn't always straightforward

## Classification Strategy

We'll employ a multi-layered classification approach that combines API metadata, value pattern analysis, and semantic examination:

### 1. Parameter Type Detection

The system will classify parameters into one of several fundamental types:

| Type | Description | Detection Method |
|------|-------------|-----------------|
| BINARY | On/Off, Enable/Disable | Toggle flag=1 OR exactly 2 unique formatted values |
| ENUM | Fixed set of discrete options | Small number of unique formatted values (≤10) |
| CONTINUOUS | Smooth range of values | Many unique formatted values with clear min/max |

### 2. Value Pattern Analysis

For continuous parameters, we'll analyze the relationship between normalized and formatted values:

| Scale Type | Characteristics | Detection Method |
|------------|-----------------|------------------|
| LINEAR | Even distribution of values | Consistent difference between adjacent formatted values |
| LOGARITHMIC | Values cluster at lower end | Larger changes at lower normalized values |
| EXPONENTIAL | Values cluster at upper end | Larger changes at higher normalized values |
| STEPPED | Discrete jumps in values | Clustering of formatted values at specific points |

### 3. Semantic Understanding

We'll extract meaningful information from parameter names and formatted values:

| Aspect | Examples | Extraction Method |
|--------|----------|-------------------|
| Units | Hz, kHz, dB, ms, s, % | Regex pattern matching on formatted values |
| Category | Frequency, Time, Gain | Keyword analysis of parameter names |
| Boundaries | Min/max values in human terms | Extraction from formatted values at 0.0 and 1.0 |

## Implementation Strategy

### Caching System

1. **Analysis Pipeline**:
   - When a plugin is first encountered, perform full parameter analysis
   - Generate classification and semantic metadata for each parameter
   - Store results in a persistent SQLite cache indexed by plugin name/type
   - Load relevant data into memory when plugin is used

2. **SQLite Schema Design**:

   ```sql
   CREATE TABLE fx_cache (
       id INTEGER PRIMARY KEY,
       name TEXT NOT NULL,
       format TEXT NOT NULL,  -- VST, VST3, AU
       created_at TIMESTAMP,
       UNIQUE(name, format)
   );

   CREATE TABLE parameter_cache (
       id INTEGER PRIMARY KEY,
       fx_id INTEGER,
       param_index INTEGER,
       name TEXT,
       type TEXT,  -- 'BINARY', 'ENUM', 'CONTINUOUS'
       scale_type TEXT,  -- 'LINEAR', 'LOGARITHMIC', 'EXPONENTIAL', 'STEPPED'
       min_value TEXT,   -- Formatted min value
       max_value TEXT,   -- Formatted max value
       unit TEXT,        -- Primary unit (Hz, dB, etc.)
       value_map TEXT,   -- JSON of {normalized→formatted} mapping points
       enum_values TEXT, -- JSON array of possible values for ENUM type
       confidence REAL,  -- Classification confidence score (0.0-1.0)
       created_at TIMESTAMP,
       FOREIGN KEY (fx_id) REFERENCES fx_cache(id),
       UNIQUE(fx_id, param_index)
   );
   ```

### Value Mapping & Translation

1. **Binary Parameters**:
   - Store the two formatted values (e.g., "on"/"off")
   - Use direct mapping (0.0 → first value, 1.0 → second value)

2. **Enum Parameters**:
   - Store complete list of options with normalized values
   - Use lookup table for bidirectional mapping

3. **Continuous Parameters**:
   - Store representative sample points mapping normalized→formatted
   - Use binary search to find closest match when translating
   - For well-behaved parameters (clear scaling), store interpolation formula

### LLM Integration

**Note: The formats shown below are for illustrative purposes only. The actual input format presented to the LLM will be determined based on empirical testing and optimization.**

When communicating with the LLM, we'll present parameters in human terms without exposing the normalization layer:

```txt
Parameter: Frequency
Current Value: 1.2 kHz
Range: 20 Hz to 20 kHz

Parameter: Mode
Current Value: Hall
Options: Room, Hall, Plate, Ambience

Parameter: Enabled
Current Value: On
Options: Off, On
```

**Note: The response format below is illustrative. The final implementation will provide the LLM with explicit rules about how parameters operate, including allowed value ranges and units.**

When the LLM suggests changes, we'll enforce a structured JSON response format:

```json
{
  "message": "I've adjusted the filter to enhance the bass frequencies",
  "changes": [
    {
      "parameter": "Frequency",
      "value": "800",
      "unit": "Hz"
    },
    {
      "parameter": "Mode",
      "value": "Room"
    }
  ]
}
```

The system will:

1. Validate the JSON structure
2. For each change, identify the parameter and parse the value/unit
3. Look up the corresponding normalized value
4. Apply the changes through REAPER's API

## Technical Implementation Considerations

### Classification Functions

For each type of analysis, create specialized detection functions with confidence scores:

```go
// Type detection with confidence
func IsBinaryParameter(uniqueValues []string, isToggle bool) (bool, float64)
func IsEnumParameter(uniqueValues []string, paramName string) (bool, float64)

// Scale analysis with confidence
func DetectParameterScaling(sampledValues []ParameterSample) (string, float64)

// Unit extraction
func ExtractUnit(formattedValues []string) string

// Overall parameter classification
func ClassifyParameter(param ParameterData) (ParameterClassification, float64)
```

Only parameters with confidence scores above a set threshold (e.g., 0.75) will be included in the LLM context. This ensures that we're only presenting reliable information rather than potentially misleading the LLM with uncertain classifications.

### Lookup Implementation

Efficient value lookup can be implemented with:

```go
// Find normalized value for a given formatted value
func LookupNormalizedValue(param ParameterCache, formattedValue string) float64 {
    // For binary/enum, use direct mapping
    if param.ParamType == "BINARY" || param.ParamType == "ENUM" {
        return directLookup(param, formattedValue)
    }
    
    // For continuous parameters, use binary search
    return binarySearchLookup(param, formattedValue)
}

// Find formatted value for a given normalized value
func LookupFormattedValue(param ParameterCache, normalizedValue float64) string {
    // Implementation varies by parameter type
}
```

## Future Enhancement Possibilities

1. **Plugin-Specific Profiles**:
   - Build specialized profiles for common plugins (Pro-Q, Serum, etc.)
   - Handle unique parameter behaviors for these popular plugins

2. **Multi-format Learning**:
   - If we encounter the same plugin in multiple formats, use the most information-rich format's data

3. **Cache Invalidation Strategy**:
   - Provide a manual "Refresh Analysis" option for users if they notice inconsistencies
   - It should delete the fx name and cascade delete other stuff from cache if triggered
   - Focus on stable cache management rather than complex versioning

## Action Items & Next Steps

1. ✅ Use existing db dump action to collect parameter data from diverse plugins
2. Perform data analysis on the collected parameter data to identify patterns
3. Develop and test classification algorithms with confidence thresholds based on empirical data
4. Implement the caching system and lookup mechanisms
5. Create parameter representation templates for LLM communication (without descriptions)
6. Test the system with various plugins to ensure robustness and accuracy

This approach balances technical precision with human-centric parameter representation, enabling the LLM FX Assistant to make accurate, meaningful parameter adjustments while hiding the complexities of parameter normalization from the user and the LLM.
