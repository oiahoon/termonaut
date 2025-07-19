package unit

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
)

func BenchmarkCommandLogging(b *testing.B) {
	// Create temporary database
	tempDir := b.TempDir()

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Initialize database
	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create a test session
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()

	// Benchmark individual command logging
	for i := 0; i < b.N; i++ {
		cmd := &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "ls -la",
			ExitCode:  0,
			CWD:       "/tmp",
		}

		err := db.StoreCommand(cmd)
		if err != nil {
			b.Fatalf("Failed to store command: %v", err)
		}
	}
}

func BenchmarkBatchCommandLogging(b *testing.B) {
	// Create temporary database
	tempDir := b.TempDir()

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Initialize database
	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create a test session
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Prepare batch of commands
	batchSize := 100
	commands := make([]*models.Command, batchSize)
	for i := 0; i < batchSize; i++ {
		commands[i] = &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "ls -la",
			ExitCode:  0,
			CWD:       "/tmp",
		}
	}

	b.ResetTimer()

	// Benchmark batch command logging
	for i := 0; i < b.N; i++ {
		err := db.StoreCommandsBatch(commands)
		if err != nil {
			b.Fatalf("Failed to store commands batch: %v", err)
		}
	}
}

func BenchmarkStatsCalculation(b *testing.B) {
	// Create temporary database
	tempDir := b.TempDir()

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Initialize database
	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Create test data
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add some test commands
	for i := 0; i < 1000; i++ {
		cmd := &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "test command",
			ExitCode:  0,
			CWD:       "/tmp",
		}
		db.StoreCommand(cmd)
	}

	b.ResetTimer()

	// Benchmark stats calculation
	for i := 0; i < b.N; i++ {
		_, err := db.GetBasicStats()
		if err != nil {
			b.Fatalf("Failed to get basic stats: %v", err)
		}
	}
}

func TestCacheEffectiveness(t *testing.T) {
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

	// Create test data
	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Add some test commands
	for i := 0; i < 100; i++ {
		cmd := &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "test command",
			ExitCode:  0,
			CWD:       "/tmp",
		}
		db.StoreCommand(cmd)
	}

	// First call - should hit database
	start1 := time.Now()
	stats1, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get basic stats: %v", err)
	}
	duration1 := time.Since(start1)

	// Second call - should hit cache
	start2 := time.Now()
	stats2, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get basic stats: %v", err)
	}
	duration2 := time.Since(start2)

	// Verify results are the same
	if stats1["total_commands"] != stats2["total_commands"] {
		t.Error("Cached results don't match original results")
	}

	// Cache should be significantly faster
	if duration2 >= duration1 {
		t.Logf("Warning: Cache doesn't seem to be working effectively. First call: %v, Second call: %v", duration1, duration2)
	} else {
		t.Logf("Cache effectiveness: First call: %v, Second call: %v (%.2fx faster)", 
			duration1, duration2, float64(duration1)/float64(duration2))
	}
}

func TestTransactionRollback(t *testing.T) {
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

	// Get initial command count
	initialStats, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get initial stats: %v", err)
	}
	initialCount := initialStats["total_commands"].(int)

	// Test transaction rollback
	err = db.WithTransaction(func(tx *sql.Tx) error {
		// Insert a command within transaction
		_, err := tx.Exec(`
			INSERT INTO commands (timestamp, session_id, command, exit_code, cwd)
			VALUES (?, ?, ?, ?, ?)
		`, time.Now(), 1, "test command", 0, "/tmp")
		if err != nil {
			return err
		}

		// Force an error to trigger rollback
		return fmt.Errorf("intentional error for rollback test")
	})

	// Error should be returned
	if err == nil {
		t.Error("Expected error from transaction, but got nil")
	}

	// Command count should be unchanged due to rollback
	finalStats, err := db.GetBasicStats()
	if err != nil {
		t.Fatalf("Failed to get final stats: %v", err)
	}
	finalCount := finalStats["total_commands"].(int)

	if finalCount != initialCount {
		t.Errorf("Transaction rollback failed: expected %d commands, got %d", initialCount, finalCount)
	}
}
