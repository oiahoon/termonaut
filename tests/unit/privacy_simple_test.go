package unit

import (
	"testing"

	"github.com/oiahoon/termonaut/internal/privacy"
)

func TestCommandSanitizerBasic(t *testing.T) {
	// Test with default configuration
	sanitizer := privacy.NewCommandSanitizer(nil)

	tests := []struct {
		name     string
		input    string
		shouldChange bool
	}{
		{
			name:         "simple command",
			input:        "ls -la",
			shouldChange: false,
		},
		{
			name:         "command with potential password",
			input:        "mysql -u root -p secret123",
			shouldChange: true,
		},
		{
			name:         "git command",
			input:        "git status",
			shouldChange: false,
		},
		{
			name:         "empty command",
			input:        "",
			shouldChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			
			if tt.shouldChange && result == tt.input {
				t.Errorf("Expected command to be sanitized, but it remained unchanged: %s", tt.input)
			}
			
			if !tt.shouldChange && result != tt.input {
				t.Errorf("Expected command to remain unchanged, but it was sanitized: %s -> %s", tt.input, result)
			}
		})
	}
}

func TestSanitizationConfig(t *testing.T) {
	// Test default configuration
	config := privacy.DefaultSanitizationConfig()
	
	if !config.Enabled {
		t.Error("Expected sanitization to be enabled by default")
	}
	
	if !config.SanitizePasswords {
		t.Error("Expected password sanitization to be enabled by default")
	}
	
	if !config.SanitizeURLs {
		t.Error("Expected URL sanitization to be enabled by default")
	}
	
	if len(config.IgnoreCommands) == 0 {
		t.Error("Expected some commands to be ignored by default")
	}
}

func TestSanitizerWithCustomConfig(t *testing.T) {
	// Test with custom configuration
	config := &privacy.SanitizationConfig{
		Enabled:           true,
		SanitizePasswords: false, // Disable password sanitization
		SanitizeURLs:      true,
		SanitizeFilePaths: true,
		SanitizeTokens:    true,
		SanitizeEmails:    true,
		IgnoreCommands:    []string{"test"},
		MaxArgLength:      100,
	}
	
	sanitizer := privacy.NewCommandSanitizer(config)
	
	// Test that password sanitization is disabled
	result := sanitizer.SanitizeCommand("mysql -u root -p secret123")
	if !containsString(result, "secret123") {
		t.Error("Expected password to not be sanitized when disabled")
	}
}

func TestIgnoredCommands(t *testing.T) {
	config := &privacy.SanitizationConfig{
		Enabled:           true,
		SanitizePasswords: true,
		IgnoreCommands:    []string{"sudo", "passwd"},
	}
	
	sanitizer := privacy.NewCommandSanitizer(config)
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ignored sudo command",
			input:    "sudo passwd user",
			expected: "sudo passwd user", // Should remain unchanged
		},
		{
			name:     "ignored passwd command",
			input:    "passwd",
			expected: "passwd", // Should remain unchanged
		},
		{
			name:     "non-ignored command",
			input:    "mysql -p secret",
			expected: "", // Should be sanitized (we don't know exact output)
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeCommand(tt.input)
			
			if tt.expected != "" && result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
			
			// For non-ignored commands, just check that they were processed
			if tt.expected == "" && result == tt.input {
				t.Errorf("Expected command to be processed, but remained unchanged: %s", tt.input)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && 
		   (s == substr || len(s) > len(substr) && 
		   (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		   containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
