package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// DefaultConfigDir is the default configuration directory
	DefaultConfigDir = ".termonaut"

	// ConfigFileName is the configuration file name
	ConfigFileName = "config"

	// ConfigFileType is the configuration file type
	ConfigFileType = "toml"
)

// Config represents the application configuration
type Config struct {
	// Display and Theme
	DisplayMode      string `mapstructure:"display_mode"`
	Theme            string `mapstructure:"theme"`
	ShowGamification bool   `mapstructure:"show_gamification"`

	// Tracking Behavior
	IdleTimeoutMinutes int  `mapstructure:"idle_timeout_minutes"`
	TrackGitRepos      bool `mapstructure:"track_git_repos"`
	CommandCategories  bool `mapstructure:"command_categories"`

	// GitHub Integration (Optional)
	SyncEnabled          bool   `mapstructure:"sync_enabled"`
	SyncRepo             string `mapstructure:"sync_repo"`
	BadgeUpdateFrequency string `mapstructure:"badge_update_frequency"`

	// Privacy
	OptOutCommands    []string `mapstructure:"opt_out_commands"`
	AnonymousMode     bool     `mapstructure:"anonymous_mode"`
	PrivacySanitizer  bool     `mapstructure:"privacy_sanitizer"`
	SanitizePasswords bool     `mapstructure:"sanitize_passwords"`
	SanitizeURLs      bool     `mapstructure:"sanitize_urls"`
	SanitizeFilePaths bool     `mapstructure:"sanitize_file_paths"`
	
	// Easter Eggs
	EasterEggsEnabled bool `mapstructure:"easter_eggs_enabled"`

	// Internal
	DataDir  string `mapstructure:"data_dir"`
	LogLevel string `mapstructure:"log_level"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, DefaultConfigDir)

	return &Config{
		// Display and Theme
		DisplayMode:      "enter",
		Theme:            "emoji",
		ShowGamification: true,

		// Tracking Behavior
		IdleTimeoutMinutes: 10,
		TrackGitRepos:      true,
		CommandCategories:  true,

		// GitHub Integration
		SyncEnabled:          false,
		SyncRepo:             "",
		BadgeUpdateFrequency: "daily",

		// Privacy
		OptOutCommands:    []string{"password", "secret", "token"},
		AnonymousMode:     false,
		PrivacySanitizer:  true,
		SanitizePasswords: true,
		SanitizeURLs:      true,
		SanitizeFilePaths: true,

		// Easter Eggs
		EasterEggsEnabled: true,

		// Internal
		DataDir:  dataDir,
		LogLevel: "info",
	}
}

// Load loads configuration from file or creates default config
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)

	// Set up viper
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.AddConfigPath(configDir)

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			config := DefaultConfig()
			if err := Save(config); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure data directory is set
	if config.DataDir == "" {
		config.DataDir = configDir
	}

	return &config, nil
}

// Save saves configuration to file
func Save(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set config values
	viper.Set("display_mode", config.DisplayMode)
	viper.Set("theme", config.Theme)
	viper.Set("show_gamification", config.ShowGamification)
	viper.Set("idle_timeout_minutes", config.IdleTimeoutMinutes)
	viper.Set("track_git_repos", config.TrackGitRepos)
	viper.Set("command_categories", config.CommandCategories)
	viper.Set("sync_enabled", config.SyncEnabled)
	viper.Set("sync_repo", config.SyncRepo)
	viper.Set("badge_update_frequency", config.BadgeUpdateFrequency)
	viper.Set("opt_out_commands", config.OptOutCommands)
	viper.Set("anonymous_mode", config.AnonymousMode)
	viper.Set("data_dir", config.DataDir)
	viper.Set("log_level", config.LogLevel)

	// Write config file
	configPath := filepath.Join(configDir, ConfigFileName+"."+ConfigFileType)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults() {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, DefaultConfigDir)

	viper.SetDefault("display_mode", "enter")
	viper.SetDefault("theme", "emoji")
	viper.SetDefault("show_gamification", true)
	viper.SetDefault("idle_timeout_minutes", 10)
	viper.SetDefault("track_git_repos", true)
	viper.SetDefault("command_categories", true)
	viper.SetDefault("sync_enabled", false)
	viper.SetDefault("sync_repo", "")
	viper.SetDefault("badge_update_frequency", "daily")
	viper.SetDefault("opt_out_commands", []string{"password", "secret", "token"})
	viper.SetDefault("anonymous_mode", false)
	viper.SetDefault("privacy_sanitizer", true)
	viper.SetDefault("sanitize_passwords", true)
	viper.SetDefault("sanitize_urls", true)
	viper.SetDefault("sanitize_file_paths", true)
	viper.SetDefault("easter_eggs_enabled", true)
	viper.SetDefault("data_dir", dataDir)
	viper.SetDefault("log_level", "info")
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DefaultConfigDir)
}

// GetDataDir returns the data directory path from config
func GetDataDir(config *Config) string {
	if config != nil && config.DataDir != "" {
		return config.DataDir
	}
	return GetConfigDir()
}
