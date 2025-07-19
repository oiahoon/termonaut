package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oiahoon/termonaut/internal/cache"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
)

const (
	// DatabaseName is the SQLite database filename
	DatabaseName = "termonaut.db"

	// DefaultTimeout for database operations
	DefaultTimeout = 5 * time.Second
	
	// CacheTTL is the time-to-live for cached queries
	CacheTTL = 5 * time.Minute
	
	// CacheCapacity is the maximum number of cached entries
	CacheCapacity = 1000
)

// CacheEntry represents a cached query result
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// DB wraps the SQLite database connection
type DB struct {
	conn   *sql.DB
	logger *logrus.Logger
	
	// Enhanced LRU cache
	lruCache   *cache.LRUCache
	cacheMutex sync.RWMutex
	logger *logrus.Logger
}

// New creates a new database connection
func New(dataDir string, logger *logrus.Logger) (*DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, DatabaseName)

	// Open SQLite database with optimized settings
	conn, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_timeout=10000&_foreign_keys=on&_synchronous=NORMAL&_cache_size=10000&_temp_store=memory")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for better performance
	conn.SetMaxOpenConns(5)  // Allow more concurrent connections for read operations
	conn.SetMaxIdleConns(2)  // Keep more idle connections
	conn.SetConnMaxLifetime(30 * time.Minute) // Shorter lifetime for better resource management

	db := &DB{
		conn:     conn,
		logger:   logger,
		lruCache: cache.NewLRUCache(CacheCapacity),
	}

	// Initialize database schema
	if err := db.initialize(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start cache cleanup timer
	go db.startCacheCleanup()

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// getCachedResult retrieves a cached result if it exists and is not expired
func (db *DB) getCachedResult(key string) (interface{}, bool) {
	if entry, ok := db.cache.Load(key); ok {
		cacheEntry := entry.(*CacheEntry)
		if time.Now().Before(cacheEntry.ExpiresAt) {
			return cacheEntry.Data, true
		}
		// Remove expired entry
		db.cache.Delete(key)
	}
	return nil, false
}

// setCachedResult stores a result in the cache with TTL
func (db *DB) setCachedResult(key string, data interface{}) {
	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(CacheTTL),
	}
	db.cache.Store(key, entry)
}

// clearCache removes all cached entries
func (db *DB) clearCache() {
	db.cache.Range(func(key, value interface{}) bool {
		db.cache.Delete(key)
		return true
	})
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("failed to commit transaction: %w", err)
			}
		}
	}()

	err = fn(tx)
	return err
}

// initialize creates the database schema
func (db *DB) initialize() error {
	schema := `
	-- Commands table: stores each executed command
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		session_id INTEGER NOT NULL,
		command TEXT NOT NULL,
		exit_code INTEGER DEFAULT 0,
		cwd TEXT,
		duration_ms INTEGER,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_commands_timestamp ON commands(timestamp);
	CREATE INDEX IF NOT EXISTS idx_commands_session ON commands(session_id);
	CREATE INDEX IF NOT EXISTS idx_commands_command ON commands(command);

	-- Sessions table: groups commands by terminal session
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		end_time DATETIME,
		terminal_pid INTEGER,
		shell_type TEXT,
		total_commands INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_start ON sessions(start_time);

	-- XP tracking: stores gamification progress
	CREATE TABLE IF NOT EXISTS user_progress (
		id INTEGER PRIMARY KEY CHECK (id = 1), -- Singleton table
		total_xp INTEGER NOT NULL DEFAULT 0,
		current_level INTEGER NOT NULL DEFAULT 1,
		commands_count INTEGER NOT NULL DEFAULT 0,
		unique_commands_count INTEGER NOT NULL DEFAULT 0,
		longest_streak INTEGER NOT NULL DEFAULT 0,
		current_streak INTEGER NOT NULL DEFAULT 0,
		last_activity_date DATE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Achievements: tracks earned badges
	CREATE TABLE IF NOT EXISTS achievements (
		id TEXT PRIMARY KEY,              -- achievement identifier
		name TEXT NOT NULL,               -- display name
		description TEXT,                 -- achievement description
		earned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		xp_bonus INTEGER DEFAULT 0        -- XP awarded for achievement
	);

	-- Daily stats cache: for performance optimization
	CREATE TABLE IF NOT EXISTS daily_stats (
		date DATE PRIMARY KEY,
		commands_count INTEGER NOT NULL DEFAULT 0,
		unique_commands_count INTEGER NOT NULL DEFAULT 0,
		session_count INTEGER NOT NULL DEFAULT 0,
		active_time_minutes INTEGER NOT NULL DEFAULT 0,
		xp_earned INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Initialize user progress if not exists
	INSERT OR IGNORE INTO user_progress (id) VALUES (1);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	db.logger.Info("Database schema initialized successfully")
	return nil
}

// StoreCommand saves a command to the database
func (db *DB) StoreCommand(cmd *models.Command) error {
	query := `
		INSERT INTO commands (timestamp, session_id, command, exit_code, cwd, duration_ms)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		cmd.Timestamp, cmd.SessionID, cmd.Command,
		cmd.ExitCode, cmd.CWD, cmd.DurationMS)
	if err != nil {
		return fmt.Errorf("failed to store command: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get command ID: %w", err)
	}

	cmd.ID = id
	
	// Clear cache when new data is added
	db.clearCache()
	
	return nil
}

// StoreCommandsBatch saves multiple commands to the database in a single transaction
func (db *DB) StoreCommandsBatch(commands []*models.Command) error {
	if len(commands) == 0 {
		return nil
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
		INSERT INTO commands (timestamp, session_id, command, exit_code, cwd, duration_ms)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, cmd := range commands {
		result, execErr := stmt.Exec(
			cmd.Timestamp, cmd.SessionID, cmd.Command,
			cmd.ExitCode, cmd.CWD, cmd.DurationMS)
		if execErr != nil {
			err = fmt.Errorf("failed to execute batch insert: %w", execErr)
			return err
		}

		id, idErr := result.LastInsertId()
		if idErr != nil {
			err = fmt.Errorf("failed to get command ID: %w", idErr)
			return err
		}
		cmd.ID = id
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Clear cache when new data is added
	db.clearCache()
	
	return nil
}

// GetOrCreateSession gets an active session or creates a new one
func (db *DB) GetOrCreateSession(pid int, shellType string) (*models.Session, error) {
	// First, try to find an active session
	query := `
		SELECT id, start_time, end_time, terminal_pid, shell_type, total_commands
		FROM sessions
		WHERE terminal_pid = ? AND end_time IS NULL
		ORDER BY start_time DESC
		LIMIT 1
	`

	var session models.Session
	err := db.conn.QueryRow(query, pid).Scan(
		&session.ID, &session.StartTime, &session.EndTime,
		&session.TerminalPID, &session.ShellType, &session.TotalCommands,
	)

	if err == sql.ErrNoRows {
		// Create new session
		insertQuery := `
			INSERT INTO sessions (terminal_pid, shell_type)
			VALUES (?, ?)
		`
		result, err := db.conn.Exec(insertQuery, pid, shellType)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get session ID: %w", err)
		}

		session = models.Session{
			ID:          id,
			StartTime:   time.Now(),
			TerminalPID: pid,
			ShellType:   shellType,
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	return &session, nil
}

// GetBasicStats returns basic usage statistics with caching
func (db *DB) GetBasicStats() (map[string]interface{}, error) {
	// Check cache first
	cacheKey := "basic_stats"
	if cached, found := db.getCachedResult(cacheKey); found {
		return cached.(map[string]interface{}), nil
	}

	stats := make(map[string]interface{})

	// Use a single query to get multiple stats for better performance
	query := `
		SELECT 
			(SELECT COUNT(*) FROM commands) as total_commands,
			(SELECT COUNT(*) FROM sessions) as total_sessions,
			(SELECT COUNT(DISTINCT command) FROM commands) as unique_commands,
			(SELECT COUNT(*) FROM commands WHERE date(timestamp) = date('now')) as commands_today
	`
	
	var totalCommands, totalSessions, uniqueCommands, commandsToday int
	err := db.conn.QueryRow(query).Scan(&totalCommands, &totalSessions, &uniqueCommands, &commandsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}
	
	stats["total_commands"] = totalCommands
	stats["total_sessions"] = totalSessions
	stats["unique_commands"] = uniqueCommands
	stats["commands_today"] = commandsToday

	// Cache the result
	db.setCachedResult(cacheKey, stats)
	
	return stats, nil

	// Unique commands
	var uniqueCommands int
	err = db.conn.QueryRow("SELECT COUNT(DISTINCT command) FROM commands").Scan(&uniqueCommands)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique commands: %w", err)
	}
	stats["unique_commands"] = uniqueCommands

	// Commands today
	var commandsToday int
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM commands
		WHERE DATE(timestamp) = DATE('now')
	`).Scan(&commandsToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's commands: %w", err)
	}
	stats["commands_today"] = commandsToday

	// Most used command
	var mostUsedCommand string
	var mostUsedCount int
	err = db.conn.QueryRow(`
		SELECT command, COUNT(*) as count
		FROM commands
		GROUP BY command
		ORDER BY count DESC
		LIMIT 1
	`).Scan(&mostUsedCommand, &mostUsedCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get most used command: %w", err)
	}
	if err != sql.ErrNoRows {
		stats["most_used_command"] = mostUsedCommand
		stats["most_used_count"] = mostUsedCount
	}

	return stats, nil
}

// GetTopCommands returns the most frequently used commands
func (db *DB) GetTopCommands(limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT command, COUNT(*) as count
		FROM commands
		GROUP BY command
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top commands: %w", err)
	}
	defer rows.Close()

	var commands []map[string]interface{}
	for rows.Next() {
		var command string
		var count int
		if err := rows.Scan(&command, &count); err != nil {
			return nil, fmt.Errorf("failed to scan command: %w", err)
		}
		commands = append(commands, map[string]interface{}{
			"command": command,
			"count":   count,
		})
	}

	return commands, nil
}

// Health checks database connectivity
func (db *DB) Health() error {
	return db.conn.Ping()
}

// startCacheCleanup starts a background goroutine to clean up expired cache entries
func (db *DB) startCacheCleanup() {
	ticker := time.NewTicker(10 * time.Minute) // Clean up every 10 minutes
	defer ticker.Stop()
	
	for range ticker.C {
		if db.lruCache != nil {
			removed := db.lruCache.CleanupExpired()
			if removed > 0 {
				db.logger.Debugf("Cleaned up %d expired cache entries", removed)
			}
		}
	}
}

// getCachedResult retrieves a cached result
func (db *DB) getCachedResult(key string) (interface{}, bool) {
	if db.lruCache == nil {
		return nil, false
	}
	return db.lruCache.Get(key)
}

// setCachedResult stores a result in cache
func (db *DB) setCachedResult(key string, value interface{}) {
	if db.lruCache != nil {
		db.lruCache.SetWithTTL(key, value, CacheTTL)
	}
}

// invalidateCache removes entries matching a pattern
func (db *DB) invalidateCache(pattern string) {
	if db.lruCache == nil {
		return
	}
	
	// For now, clear all cache on any invalidation
	// In the future, we could implement pattern matching
	db.lruCache.Clear()
	db.logger.Debug("Cache invalidated")
}

// GetCacheStats returns cache statistics
func (db *DB) GetCacheStats() cache.CacheStats {
	if db.lruCache == nil {
		return cache.CacheStats{}
	}
	return db.lruCache.Stats()
}

// StoreCommandsBatch stores multiple commands in a single transaction for better performance
func (db *DB) StoreCommandsBatch(commands []*models.Command) error {
	if len(commands) == 0 {
		return nil
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO commands (timestamp, session_id, command, exit_code, cwd, duration_ms, category)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, cmd := range commands {
		_, err = stmt.Exec(
			cmd.Timestamp,
			cmd.SessionID,
			cmd.Command,
			cmd.ExitCode,
			cmd.CWD,
			cmd.DurationMs,
			cmd.Category,
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Invalidate cache after batch insert
	db.invalidateCache("stats")
	
	return nil
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
