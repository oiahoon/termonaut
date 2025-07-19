package integration

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

// TestFullWorkflow tests the complete user workflow
func TestFullWorkflow(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Step 1: Initialize configuration
	cfg := config.DefaultConfig()
	cfg.ShowGamification = true
	cfg.Theme = "emoji"
	
	err := config.Save(cfg)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Step 2: Initialize database
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Step 3: Create a session (simulating user starting terminal)
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == 0 {
		t.Error("Session ID should not be 0")
	}

	// Step 4: Log some commands (simulating user activity)
	commands := []*models.Command{
		{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "ls -la",
			ExitCode:  0,
			CWD:       "/tmp",
		},
		{
			Timestamp: time.Now().Add(1 * time.Second),
			SessionID: session.ID,
			Command:   "git status",
			ExitCode:  0,
			CWD:       "/tmp/project",
		},
		{
			Timestamp: time.Now().Add(2 * time.Second),
			SessionID: session.ID,
			Command:   "npm install",
			ExitCode:  0,
			CWD:       "/tmp/project",
		},
	}

	for _, cmd := range commands {
		err := db.StoreCommand(cmd)
		if err != nil {
			t.Fatalf("Failed to store command: %v", err)
		}
	}

	// Step 5: Verify statistics are calculated correctly
	basicStats, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get basic stats: %v", err)
	}

	totalCommands := basicStats["total_commands"].(int)
	if totalCommands != len(commands) {
		t.Errorf("Expected %d commands, got %d", len(commands), totalCommands)
	}

	// Step 6: Test gamification features
	userProgress, err := db.GetUserProgress()
	if err != nil {
		t.Fatalf("Failed to get user progress: %v", err)
	}

	if userProgress.TotalCommands != len(commands) {
		t.Errorf("Expected user progress to show %d commands, got %d", 
			len(commands), userProgress.TotalCommands)
	}

	if userProgress.TotalXP <= 0 {
		t.Error("Expected user to have earned some XP")
	}

	// Step 7: Test session management
	sessions, err := db.GetRecentSessions(10)
	if err != nil {
		t.Fatalf("Failed to get recent sessions: %v", err)
	}

	if len(sessions) == 0 {
		t.Error("Expected at least one session")
	}

	// Step 8: Test configuration reload
	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if loadedCfg.Theme != cfg.Theme {
		t.Errorf("Config theme mismatch: expected %s, got %s", cfg.Theme, loadedCfg.Theme)
	}
}

// TestDatabaseIntegration tests database operations integration
func TestDatabaseIntegration(t *testing.T) {
	tempDir := t.TempDir()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test 1: Session creation and retrieval
	session1, err := db.GetOrCreateSession(11111, "bash")
	if err != nil {
		t.Fatalf("Failed to create session 1: %v", err)
	}

	session2, err := db.GetOrCreateSession(22222, "zsh")
	if err != nil {
		t.Fatalf("Failed to create session 2: %v", err)
	}

	if session1.ID == session2.ID {
		t.Error("Different sessions should have different IDs")
	}

	// Test 2: Command storage and retrieval
	testCommands := []*models.Command{
		{
			Timestamp: time.Now(),
			SessionID: session1.ID,
			Command:   "echo 'hello'",
			ExitCode:  0,
			CWD:       "/home/user",
		},
		{
			Timestamp: time.Now(),
			SessionID: session2.ID,
			Command:   "ls -l",
			ExitCode:  0,
			CWD:       "/tmp",
		},
	}

	for _, cmd := range testCommands {
		err := db.StoreCommand(cmd)
		if err != nil {
			t.Fatalf("Failed to store command: %v", err)
		}
		
		if cmd.ID == 0 {
			t.Error("Command ID should be set after storage")
		}
	}

	// Test 3: Statistics calculation
	stats, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get basic stats: %v", err)
	}

	expectedCommands := len(testCommands)
	actualCommands := stats["total_commands"].(int)
	if actualCommands != expectedCommands {
		t.Errorf("Expected %d commands in stats, got %d", expectedCommands, actualCommands)
	}

	expectedSessions := 2
	actualSessions := stats["total_sessions"].(int)
	if actualSessions != expectedSessions {
		t.Errorf("Expected %d sessions in stats, got %d", expectedSessions, actualSessions)
	}

	// Test 4: Gamification integration
	progress, err := db.GetUserProgress()
	if err != nil {
		t.Fatalf("Failed to get user progress: %v", err)
	}

	if progress.TotalCommands != expectedCommands {
		t.Errorf("Expected progress to show %d commands, got %d", 
			expectedCommands, progress.TotalCommands)
	}

	if progress.TotalXP <= 0 {
		t.Error("Expected user to have earned XP")
	}

	// Test 5: Cache functionality (if implemented)
	// First call should hit database
	stats1, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get stats (first call): %v", err)
	}

	// Second call should potentially hit cache
	stats2, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get stats (second call): %v", err)
	}

	// Results should be consistent
	if stats1["total_commands"] != stats2["total_commands"] {
		t.Error("Stats should be consistent between calls")
	}
}

// TestConfigIntegration tests configuration system integration
func TestConfigIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test 1: Default configuration creation
	defaultCfg := config.DefaultConfig()
	if defaultCfg == nil {
		t.Fatal("Default config should not be nil")
	}

	// Test 2: Configuration save and load cycle
	err := config.Save(defaultCfg)
	if err != nil {
		t.Fatalf("Failed to save default config: %v", err)
	}

	loadedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test 3: Configuration values consistency
	if loadedCfg.DisplayMode != defaultCfg.DisplayMode {
		t.Errorf("Display mode mismatch: expected %s, got %s", 
			defaultCfg.DisplayMode, loadedCfg.DisplayMode)
	}

	if loadedCfg.Theme != defaultCfg.Theme {
		t.Errorf("Theme mismatch: expected %s, got %s", 
			defaultCfg.Theme, loadedCfg.Theme)
	}

	// Test 4: Configuration modification
	loadedCfg.Theme = "minimal"
	loadedCfg.ShowGamification = false

	err = config.Save(loadedCfg)
	if err != nil {
		t.Fatalf("Failed to save modified config: %v", err)
	}

	modifiedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load modified config: %v", err)
	}

	if modifiedCfg.Theme != "minimal" {
		t.Errorf("Expected theme 'minimal', got %s", modifiedCfg.Theme)
	}

	if modifiedCfg.ShowGamification != false {
		t.Error("Expected show_gamification to be false")
	}

	// Test 5: Configuration file existence and permissions
	configPath := filepath.Join(tempDir, config.DefaultConfigDir, "config.toml")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Config file should exist: %v", err)
	}

	if info.IsDir() {
		t.Error("Config path should be a file, not a directory")
	}

	// Check basic file permissions (readable)
	if info.Mode().Perm()&0400 == 0 {
		t.Error("Config file should be readable by owner")
	}
}
