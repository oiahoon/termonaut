package benchmark

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/privacy"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
)

// BenchmarkCommandLogging benchmarks the command logging performance
func BenchmarkCommandLogging(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Benchmark command logging
	b.ResetTimer()
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

// BenchmarkBatchCommandLogging benchmarks batch command logging
func BenchmarkBatchCommandLogging(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

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
			Command:   "test command",
			ExitCode:  0,
			CWD:       "/tmp",
		}
	}

	// Benchmark batch logging
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.StoreCommandsBatch(commands)
		if err != nil {
			b.Fatalf("Failed to store commands batch: %v", err)
		}
	}
}

// BenchmarkStatsCalculation benchmarks statistics calculation
func BenchmarkStatsCalculation(b *testing.B) {
	// Setup with some data
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add some test data
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

	// Benchmark stats calculation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.GetBasicStats()
		if err != nil {
			b.Fatalf("Failed to get basic stats: %v", err)
		}
	}
}

// BenchmarkStatsCalculationWithCache benchmarks cached stats calculation
func BenchmarkStatsCalculationWithCache(b *testing.B) {
	// Setup with some data
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add some test data
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

	// Prime the cache
	db.GetBasicStats()

	// Benchmark cached stats calculation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.GetBasicStats()
		if err != nil {
			b.Fatalf("Failed to get basic stats: %v", err)
		}
	}
}

// BenchmarkUserProgressCalculation benchmarks user progress calculation
func BenchmarkUserProgressCalculation(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add some test data
	for i := 0; i < 500; i++ {
		cmd := &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "test command",
			ExitCode:  0,
			CWD:       "/tmp",
		}
		db.StoreCommand(cmd)
	}

	// Benchmark user progress calculation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.GetUserProgress()
		if err != nil {
			b.Fatalf("Failed to get user progress: %v", err)
		}
	}
}

// BenchmarkCommandSanitization benchmarks command sanitization
func BenchmarkCommandSanitization(b *testing.B) {
	sanitizer := privacy.NewCommandSanitizer(nil)
	
	testCommands := []string{
		"ls -la /home/user",
		"mysql -u root -p'secret123' -h localhost",
		"git clone https://token@github.com/user/repo.git",
		"curl -u admin:password https://api.example.com",
		"ssh user@192.168.1.100",
		"echo 'hello world'",
		"npm install --save-dev package",
		"docker run -e PASSWORD=secret image",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := testCommands[i%len(testCommands)]
		sanitizer.SanitizeCommand(cmd)
	}
}

// BenchmarkConfigOperations benchmarks configuration operations
func BenchmarkConfigOperations(b *testing.B) {
	// Setup temporary directory
	tempDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cfg := config.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark save and load cycle
		err := config.Save(cfg)
		if err != nil {
			b.Fatalf("Failed to save config: %v", err)
		}

		_, err = config.Load()
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

// BenchmarkSessionManagement benchmarks session management operations
func BenchmarkSessionManagement(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pid := 10000 + i
		_, err := db.GetOrCreateSession(pid, "zsh")
		if err != nil {
			b.Fatalf("Failed to get or create session: %v", err)
		}
	}
}

// BenchmarkDatabaseTransaction benchmarks database transaction operations
func BenchmarkDatabaseTransaction(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.WithTransaction(func(tx *sql.Tx) error {
			// Simulate some database operations within transaction
			_, err := tx.Exec("INSERT INTO commands (timestamp, session_id, command, exit_code, cwd) VALUES (?, ?, ?, ?, ?)",
				time.Now(), session.ID, "test command", 0, "/tmp")
			return err
		})
		if err != nil {
			b.Fatalf("Transaction failed: %v", err)
		}
	}
}

// BenchmarkMemoryUsage benchmarks memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	db, err := database.New(tempDir, logger)
	if err != nil {
		b.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	session, err := db.GetOrCreateSession(12345, "zsh")
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs() // Report memory allocations

	for i := 0; i < b.N; i++ {
		// Create and store command
		cmd := &models.Command{
			Timestamp: time.Now(),
			SessionID: session.ID,
			Command:   "memory test command",
			ExitCode:  0,
			CWD:       "/tmp",
		}

		err := db.StoreCommand(cmd)
		if err != nil {
			b.Fatalf("Failed to store command: %v", err)
		}

		// Get stats to test memory usage
		_, err = db.GetBasicStats()
		if err != nil {
			b.Fatalf("Failed to get stats: %v", err)
		}
	}
}
