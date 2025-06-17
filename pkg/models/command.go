package models

import (
	"time"
)

// Command represents a single terminal command execution
type Command struct {
	ID         int64     `json:"id" db:"id"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	SessionID  int64     `json:"session_id" db:"session_id"`
	Command    string    `json:"command" db:"command"`
	ExitCode   int       `json:"exit_code" db:"exit_code"`
	CWD        string    `json:"cwd" db:"cwd"`
	DurationMS int64     `json:"duration_ms" db:"duration_ms"`
}

// Session represents a terminal session
type Session struct {
	ID            int64      `json:"id" db:"id"`
	StartTime     time.Time  `json:"start_time" db:"start_time"`
	EndTime       *time.Time `json:"end_time,omitempty" db:"end_time"`
	TerminalPID   int        `json:"terminal_pid" db:"terminal_pid"`
	ShellType     string     `json:"shell_type" db:"shell_type"`
	TotalCommands int        `json:"total_commands" db:"total_commands"`
}

// UserProgress represents gamification progress
type UserProgress struct {
	ID                  int        `json:"id" db:"id"`
	TotalXP             int        `json:"total_xp" db:"total_xp"`
	CurrentLevel        int        `json:"current_level" db:"current_level"`
	CommandsCount       int        `json:"commands_count" db:"commands_count"`
	UniqueCommandsCount int        `json:"unique_commands_count" db:"unique_commands_count"`
	LongestStreak       int        `json:"longest_streak" db:"longest_streak"`
	CurrentStreak       int        `json:"current_streak" db:"current_streak"`
	LastActivityDate    *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// Achievement represents an earned achievement/badge
type Achievement struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	EarnedAt    time.Time `json:"earned_at" db:"earned_at"`
	XPBonus     int       `json:"xp_bonus" db:"xp_bonus"`
}

// DailyStats represents cached daily statistics for performance
type DailyStats struct {
	Date                time.Time `json:"date" db:"date"`
	CommandsCount       int       `json:"commands_count" db:"commands_count"`
	UniqueCommandsCount int       `json:"unique_commands_count" db:"unique_commands_count"`
	SessionCount        int       `json:"session_count" db:"session_count"`
	ActiveTimeMinutes   int       `json:"active_time_minutes" db:"active_time_minutes"`
	XPEarned            int       `json:"xp_earned" db:"xp_earned"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
