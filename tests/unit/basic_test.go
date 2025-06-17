package unit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
)

func TestBasicDatabaseOperations(t *testing.T) {
	// Create temporary database
	tempDir := t.TempDir()

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Initialize database
	db, err := database.New(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test session creation
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == 0 {
		t.Error("Session ID should not be 0")
	}

	if session.TerminalPID != 12345 {
		t.Errorf("Expected PID 12345, got %d", session.TerminalPID)
	}

	// Test command storage
	cmd := &models.Command{
		Timestamp: time.Now(),
		SessionID: session.ID,
		Command:   "ls -la",
		ExitCode:  0,
		CWD:       "/tmp",
	}

	err = db.StoreCommand(cmd)
	if err != nil {
		t.Fatalf("Failed to store command: %v", err)
	}

	if cmd.ID == 0 {
		t.Error("Command ID should not be 0")
	}

	// Test stats retrieval
	basicStats, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get basic stats: %v", err)
	}

	if basicStats["total_commands"].(int) != 1 {
		t.Errorf("Expected 1 command, got %d", basicStats["total_commands"])
	}
}

func TestConfig(t *testing.T) {
	// Create temporary config directory
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create default config
	cfg := config.DefaultConfig()

	if cfg.DisplayMode != "enter" {
		t.Errorf("Expected default display mode 'enter', got '%s'", cfg.DisplayMode)
	}

	if cfg.Theme != "emoji" {
		t.Errorf("Expected default theme 'emoji', got '%s'", cfg.Theme)
	}

	if !cfg.ShowGamification {
		t.Error("Expected show_gamification to be true by default")
	}

	// Test config save and load
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

	if loadedCfg.DisplayMode != cfg.DisplayMode {
		t.Errorf("Loaded config display mode mismatch: expected '%s', got '%s'",
			cfg.DisplayMode, loadedCfg.DisplayMode)
	}
}
