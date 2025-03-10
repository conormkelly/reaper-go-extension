package config

import (
	"encoding/json"
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"sync"

	"github.com/zalando/go-keyring"
)

// A unified configuration system for the REAPER Go extension.

// VERSION indicates the settings schema version
// Increment this when making incompatible changes to settings structure
const VERSION = 1

// Constants for keyring access
const (
	KeyringServiceName = "GoReaperExtension"
)

// KeyringKeys represents different API keys we might store
const (
	KeyringOpenAI = "OpenAIAPIKey"
	// KeyringClaude = "ClaudeAPIKey"
)

// Provider represents supported LLM providers
type Provider string

// Provider constants
const (
	ProviderOpenAI Provider = "openai"
	// ProviderClaude   Provider = "claude"
	// ProviderOllama   Provider = "ollama"
	// ProviderLMStudio Provider = "lmstudio"
	// Add more providers as needed
)

// Settings defines the structure of our application settings
type Settings struct {
	// Schema version for migration support
	Version int `json:"version"`

	// Active provider configuration
	ActiveProvider Provider `json:"active_provider"`

	// Provider-specific configurations
	Providers struct {
		// OpenAI specific settings
		OpenAI struct {
			Model       string  `json:"model"`
			MaxTokens   int     `json:"max_tokens"`
			Temperature float64 `json:"temperature"`
		} `json:"openai"`

		// TODO: add other providers
	} `json:"providers"`

	// Prompt settings
	Prompt struct {
		DefaultPrompt string `json:"default_prompt"`
	} `json:"prompt"`

	// General plugin settings
	General struct {
		AutoApplyChanges bool `json:"auto_apply_changes"`
		// Add more general settings as needed
	} `json:"general"`
}

// DefaultSettings provides the default configuration
var DefaultSettings = Settings{
	Version:        VERSION,
	ActiveProvider: ProviderOpenAI,
	Providers: struct {
		OpenAI struct {
			Model       string  `json:"model"`
			MaxTokens   int     `json:"max_tokens"`
			Temperature float64 `json:"temperature"`
		} `json:"openai"`
	}{
		OpenAI: struct {
			Model       string  `json:"model"`
			MaxTokens   int     `json:"max_tokens"`
			Temperature float64 `json:"temperature"`
		}{
			Model:       "gpt-3.5-turbo",
			MaxTokens:   1024,
			Temperature: 0.7,
		},
	},
	Prompt: struct {
		DefaultPrompt string `json:"default_prompt"`
	}{
		DefaultPrompt: "", // TODO: centralize this
	},
	General: struct {
		AutoApplyChanges bool `json:"auto_apply_changes"`
	}{
		AutoApplyChanges: false,
	},
}

// ExtState keys - note we use a consistent key, versioning is handled within the JSON
const (
	ExtStateSection = "GoReaperExtension"
	ExtStateKey     = "Settings"
)

// configMutex protects access to the settings
var configMutex sync.RWMutex

// GetSecureAPIKey retrieves an API key from the system keyring
func GetSecureAPIKey(provider Provider) (string, error) {
	keyName := providerToKeyringKey(provider)
	return keyring.Get(KeyringServiceName, keyName)
}

// StoreSecureAPIKey stores an API key in the system keyring
func StoreSecureAPIKey(provider Provider, apiKey string) error {
	keyName := providerToKeyringKey(provider)
	return keyring.Set(KeyringServiceName, keyName, apiKey)
}

// HasSecureAPIKey checks if an API key exists in the keyring
func HasSecureAPIKey(provider Provider) bool {
	key, err := GetSecureAPIKey(provider)
	return err == nil && key != ""
}

// providerToKeyringKey converts a provider to its keyring key
func providerToKeyringKey(provider Provider) string {
	switch provider {
	case ProviderOpenAI:
		return KeyringOpenAI
	default:
		return KeyringOpenAI
	}
}

// GetSettings returns the current settings
func GetSettings() Settings {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Always load from storage to ensure fresh data
	return loadSettings()
}

// SaveSettings saves the settings
func SaveSettings(settings Settings) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Ensure version is current before saving
	settings.Version = VERSION

	// Convert settings to JSON
	jsonData, err := json.Marshal(settings)
	if err != nil {
		logger.Error("Failed to marshal settings: %v", err)
		return err
	}

	// Save to REAPER's ExtState
	err = reaper.SetExtState(ExtStateSection, ExtStateKey, string(jsonData), true)
	if err != nil {
		logger.Error("Failed to save settings: %v", err)
		return err
	}

	logger.Debug("Settings saved successfully")
	return nil
}

// loadSettings loads settings from REAPER's ExtState
func loadSettings() Settings {
	// Start with defaults
	settings := DefaultSettings

	// Try to get from REAPER
	jsonData, err := reaper.GetExtState(ExtStateSection, ExtStateKey)
	if err != nil || jsonData == "" {
		logger.Debug("No saved settings found, using defaults")
		return settings
	}

	// Parse JSON
	err = json.Unmarshal([]byte(jsonData), &settings)
	if err != nil {
		logger.Warning("Failed to parse settings JSON, using defaults: %v", err)
		return DefaultSettings
	}

	// Handle version migrations if needed
	if settings.Version < VERSION {
		logger.Info("Migrating settings from version %d to %d", settings.Version, VERSION)
		migratedSettings, err := migrateSettingsStepByStep(settings)
		if err != nil {
			logger.Warning("Failed to migrate settings: %v - using defaults", err)
			return DefaultSettings
		}

		// Save the migrated settings
		settings = migratedSettings

		// Update storage with migrated settings
		jsonData, err := json.Marshal(settings)
		if err == nil {
			reaper.SetExtState(ExtStateSection, ExtStateKey, string(jsonData), true)
			logger.Info("Saved migrated settings")
		}
	}

	logger.Debug("Settings loaded successfully")
	return settings
}

// migrateSettingsStepByStep applies migrations sequentially
func migrateSettingsStepByStep(oldSettings Settings) (Settings, error) {
	// Start with the settings as they are
	settings := oldSettings

	// Apply migrations sequentially for each version jump
	// We must apply them in order (v1→v2, then v2→v3, etc.)
	fromVersion := settings.Version

	// When VERSION is incremented, add case statements for the migrations
	for v := fromVersion; v < VERSION; v++ {
		logger.Debug("Applying migration from v%d to v%d", v, v+1)

		switch v {
		case 1:
			// v1 to v2 migration (when we add v2)
			// settings = migrateV1toV2(settings)

		case 2:
			// v2 to v3 migration (when we add v3)
			// settings = migrateV2toV3(settings)

		// Add more cases for future versions
		default:
			logger.Warning("Unknown version to migrate from: %d", v)
		}
	}

	// Update version number to current
	settings.Version = VERSION

	return settings, nil
}

// Future migration helpers would be defined below:

// migrateV1toV2 handles migration from v1 to v2
// func migrateV1toV2(settings Settings) Settings {
//     // Example migration logic
//     newSettings := settings
//
//     // Add new field with default value
//     // newSettings.NewFieldInV2 = defaultValue
//
//     return newSettings
// }

// GetActiveProvider returns the currently active LLM provider
func GetActiveProvider() Provider {
	return GetSettings().ActiveProvider
}

// SetActiveProvider sets the active LLM provider
func SetActiveProvider(provider Provider) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	settings := loadSettings()
	settings.ActiveProvider = provider
	return SaveSettings(settings)
}

// GetProviderConfig returns the configuration for the specified provider
func GetProviderConfig(provider Provider) (model string, maxTokens int, temperature float64) {
	settings := GetSettings()

	switch provider {
	case ProviderOpenAI:
		return settings.Providers.OpenAI.Model,
			settings.Providers.OpenAI.MaxTokens,
			settings.Providers.OpenAI.Temperature
	default:
		// Fallback to OpenAI config
		logger.Warning("Unknown provider %s, using OpenAI configuration", provider)
		return settings.Providers.OpenAI.Model,
			settings.Providers.OpenAI.MaxTokens,
			settings.Providers.OpenAI.Temperature
	}
}

// GetActiveProviderConfig returns the configuration for the active provider
func GetActiveProviderConfig() (model string, maxTokens int, temperature float64) {
	provider := GetActiveProvider()
	return GetProviderConfig(provider)
}

// SetProviderConfig sets the configuration for the specified provider
func SetProviderConfig(provider Provider, model string, maxTokens int, temperature float64) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	settings := loadSettings()

	switch provider {
	case ProviderOpenAI:
		settings.Providers.OpenAI.Model = model
		settings.Providers.OpenAI.MaxTokens = maxTokens
		settings.Providers.OpenAI.Temperature = temperature
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	return SaveSettings(settings)
}

// GetPromptConfig returns the prompt configuration
func GetPromptConfig() (defaultPrompt string) {
	return GetSettings().Prompt.DefaultPrompt
}

// SetPromptConfig sets the prompt configuration
func SetPromptConfig(defaultPrompt string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	settings := loadSettings()
	settings.Prompt.DefaultPrompt = defaultPrompt

	return SaveSettings(settings)
}

// GetGeneralConfig returns the general configuration
func GetGeneralConfig() (autoApplyChanges bool) {
	return GetSettings().General.AutoApplyChanges
}

// SetGeneralConfig sets the general configuration
func SetGeneralConfig(autoApplyChanges bool) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	settings := loadSettings()
	settings.General.AutoApplyChanges = autoApplyChanges

	return SaveSettings(settings)
}

// ResetToDefaults resets all settings to defaults
func ResetToDefaults() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Save default settings
	return SaveSettings(DefaultSettings)
}
