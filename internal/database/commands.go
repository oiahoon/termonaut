package database

import (
	"database/sql"
	"fmt"

	"github.com/oiahoon/termonaut/pkg/models"
)

// GetAllCommands returns all commands from the database
func (db *DB) GetAllCommands() ([]*models.Command, error) {
	query := `
		SELECT id, timestamp, session_id, command, exit_code, cwd, duration_ms
		FROM commands
		ORDER BY timestamp DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query commands: %w", err)
	}
	defer rows.Close()

	var commands []*models.Command
	for rows.Next() {
		var cmd models.Command
		var durationMs sql.NullInt64

		err := rows.Scan(
			&cmd.ID, &cmd.Timestamp, &cmd.SessionID,
			&cmd.Command, &cmd.ExitCode, &cmd.CWD, &durationMs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan command: %w", err)
		}

		if durationMs.Valid {
			cmd.DurationMS = durationMs.Int64
		}

		commands = append(commands, &cmd)
	}

	return commands, nil
}

// GetAllSessions returns all sessions from the database
func (db *DB) GetAllSessions() ([]*models.Session, error) {
	query := `
		SELECT id, start_time, end_time, terminal_pid, shell_type, total_commands
		FROM sessions
		ORDER BY start_time DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		var endTime sql.NullTime

		err := rows.Scan(
			&session.ID, &session.StartTime, &endTime,
			&session.TerminalPID, &session.ShellType, &session.TotalCommands,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		if endTime.Valid {
			session.EndTime = &endTime.Time
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}
