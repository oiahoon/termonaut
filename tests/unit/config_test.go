package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oiahoon/termonaut/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	// Test default values
	if cfg.DisplayMode != "enter" {
		t.Errorf("Expected default display mode 'enter', got '%s'", cfg.DisplayMode)
	}

	if cfg.Theme != "emoji" {
		t.Errorf("Expected default theme 'emoji', got '%s'", cfg.Theme)
	}

	if !cfg.ShowGamification {
		t.Error("Expected show_gamification to be true by default")
	}

	if cfg.IdleTimeoutMinutes != 10 {
		t.Errorf("Expected idle timeout 10 minutes, got %d", cfg.IdleTimeoutMinutes)
	}

	// Test UI config defaults
	if cfg.UI.DefaultMode != "smart" {
		t.Errorf("Expected UI default mode 'smart', got '%s'", cfg.UI.DefaultMode)
	}

	if cfg.UI.Theme != "space" {
		t.Errorf("Expected UI theme 'space', got '%s'", cfg.UI.Theme)
	}

	if !cfg.UI.ShowAvatar {
		t.Error("Expected UI show_avatar to be true by default")
	}

	// Test privacy defaults
	if cfg.AnonymousMode {
		t.Error("Expected anonymous_mode to be false by default")
	}

	if !cfg.PrivacySanitizer {
		t.Error("Expected privacy_sanitizer to be true by default")
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create custom config
	cfg := config.DefaultConfig()
	cfg.DisplayMode = "ps1"
	cfg.Theme = "minimal"
	cfg.ShowGamification = false
	cfg.IdleTimeoutMinutes = 15
	cfg.UI.DefaultMode = "compact"
	cfg.UI.Theme = "cyberpunk"
	cfg.AnonymousMode = true

	// Save config
	err := config.Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify config file exists
	configPath := filepath.Join(tempDir, config.DefaultConfigDir, "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load config
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded values
	if loadedCfg.DisplayMode != cfg.DisplayMode {
		t.Errorf("Display mode mismatch: expected '%s', got '%s'", cfg.DisplayMode, loadedCfg.DisplayMode)
	}

	if loadedCfg.Theme != cfg.Theme {
		t.Errorf("Theme mismatch: expected '%s', got '%s'", cfg.Theme, loadedCfg.Theme)
	}

	if loadedCfg.ShowGamification != cfg.ShowGamification {
		t.Errorf("ShowGamification mismatch: expected %v, got %v", cfg.ShowGamification, loadedCfg.ShowGamification)
	}

	if loadedCfg.IdleTimeoutMinutes != cfg.IdleTimeoutMinutes {
		t.Errorf("IdleTimeoutMinutes mismatch: expected %d, got %d", cfg.IdleTimeoutMinutes, loadedCfg.IdleTimeoutMinutes)
	}

	if loadedCfg.UI.DefaultMode != cfg.UI.DefaultMode {
		t.Errorf("UI DefaultMode mismatch: expected '%s', got '%s'", cfg.UI.DefaultMode, loadedCfg.UI.DefaultMode)
	}

	if loadedCfg.AnonymousMode != cfg.AnonymousMode {
		t.Errorf("AnonymousMode mismatch: expected %v, got %v", cfg.AnonymousMode, loadedCfg.AnonymousMode)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		modifyFunc  func(*config.Config)
		expectError bool
	}{
		{
			name: "valid config",
			modifyFunc: func(cfg *config.Config) {
				// No modifications - should be valid
			},
			expectError: false,
		},
		{
			name: "invalid display mode",
			modifyFunc: func(cfg *config.Config) {
				cfg.DisplayMode = "invalid_mode"
			},
			expectError: true,
		},
		{
			name: "invalid theme",
			modifyFunc: func(cfg *config.Config) {
				cfg.Theme = "invalid_theme"
			},
			expectError: true,
		},
		{
			name: "negative idle timeout",
			modifyFunc: func(cfg *config.Config) {
				cfg.IdleTimeoutMinutes = -1
			},
			expectError: true,
		},
		{
			name: "invalid UI mode",
			modifyFunc: func(cfg *config.Config) {
				cfg.UI.DefaultMode = "invalid_ui_mode"
			},
			expectError: true,
		},
		{
			name: "invalid avatar style",
			modifyFunc: func(cfg *config.Config) {
				cfg.AvatarStyle = "invalid_style"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			tt.modifyFunc(cfg)

			err := config.Validate(cfg)
			if tt.expectError && err == nil {
				t.Error("Expected validation error, but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestConfigGetSet(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test setting values
	testCases := []struct {
		key      string
		value    interface{}
		expected interface{}
	}{
		{"display_mode", "ps1", "ps1"},
		{"theme", "minimal", "minimal"},
		{"show_gamification", false, false},
		{"idle_timeout_minutes", 20, 20},
		{"anonymous_mode", true, true},
	}

	for _, tc := range testCases {
		t.Run("set_"+tc.key, func(t *testing.T) {
			err := config.Set(tc.key, tc.value)
			if err != nil {
				t.Fatalf("Failed to set %s: %v", tc.key, err)
			}

			// Get the value back
			value, err := config.Get(tc.key)
			if err != nil {
				t.Fatalf("Failed to get %s: %v", tc.key, err)
			}

			if value != tc.expected {
				t.Errorf("Expected %s = %v, got %v", tc.key, tc.expected, value)
			}
		})
	}
}

func TestConfigReset(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Modify config
	err := config.Set("theme", "minimal")
	if err != nil {
		t.Fatalf("Failed to set theme: %v", err)
	}

	err = config.Set("show_gamification", false)
	if err != nil {
		t.Fatalf("Failed to set show_gamification: %v", err)
	}

	// Reset config
	err = config.Reset()
	if err != nil {
		t.Fatalf("Failed to reset config: %v", err)
	}

	// Verify values are back to defaults
	theme, err := config.Get("theme")
	if err != nil {
		t.Fatalf("Failed to get theme: %v", err)
	}
	if theme != "emoji" {
		t.Errorf("Expected theme to be reset to 'emoji', got '%v'", theme)
	}

	showGamification, err := config.Get("show_gamification")
	if err != nil {
		t.Fatalf("Failed to get show_gamification: %v", err)
	}
	if showGamification != true {
		t.Errorf("Expected show_gamification to be reset to true, got %v", showGamification)
	}
}

func TestConfigFilePermissions(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save config
	cfg := config.DefaultConfig()
	err := config.Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check file permissions
	configPath := filepath.Join(tempDir, config.DefaultConfigDir, "config.toml")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Config file should be readable by owner and group, but not world-readable
	mode := info.Mode()
	if mode.Perm() != 0644 {
		t.Errorf("Expected config file permissions 0644, got %o", mode.Perm())
	}
}

func TestConfigDirectoryCreation(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Config directory should not exist initially
	configDir := filepath.Join(tempDir, config.DefaultConfigDir)
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Error("Config directory should not exist initially")
	}

	// Save config should create directory
	cfg := config.DefaultConfig()
	err := config.Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Config directory should now exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory should be created when saving config")
	}

	// Check directory permissions
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("Failed to stat config directory: %v", err)
	}

	if !info.IsDir() {
		t.Error("Config path should be a directory")
	}

	// Directory should have appropriate permissions
	mode := info.Mode()
	if mode.Perm() != 0755 {
		t.Errorf("Expected config directory permissions 0755, got %o", mode.Perm())
	}
}
