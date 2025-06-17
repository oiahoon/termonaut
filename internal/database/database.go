package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
)

const (
	// DatabaseName is the SQLite database filename
	DatabaseName = "termonaut.db"

	// DefaultTimeout for database operations
	DefaultTimeout = 5 * time.Second
)

// DB wraps the SQLite database connection
type DB struct {
	conn   *sql.DB
	logger *logrus.Logger
}

// New creates a new database connection
func New(dataDir string, logger *logrus.Logger) (*DB, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, DatabaseName)

	// Open SQLite database
	conn, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(1) // SQLite doesn't benefit from multiple connections
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(time.Hour)

	db := &DB{
		conn:   conn,
		logger: logger,
	}

	// Initialize database schema
	if err := db.initialize(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
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

// GetBasicStats returns basic usage statistics
func (db *DB) GetBasicStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total commands
	var totalCommands int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM commands").Scan(&totalCommands)
	if err != nil {
		return nil, fmt.Errorf("failed to get total commands: %w", err)
	}
	stats["total_commands"] = totalCommands

	// Total sessions
	var totalSessions int
	err = db.conn.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&totalSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get total sessions: %w", err)
	}
	stats["total_sessions"] = totalSessions

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
