package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/database"
)

// StatsCalculator handles statistics computation
type StatsCalculator struct {
	db *database.DB
}

// New creates a new stats calculator
func New(db *database.DB) *StatsCalculator {
	return &StatsCalculator{
		db: db,
	}
}

// BasicStats represents basic usage statistics
type BasicStats struct {
	TotalCommands    int                      `json:"total_commands"`
	TotalSessions    int                      `json:"total_sessions"`
	UniqueCommands   int                      `json:"unique_commands"`
	CommandsToday    int                      `json:"commands_today"`
	MostUsedCommand  string                   `json:"most_used_command,omitempty"`
	MostUsedCount    int                      `json:"most_used_count"`
	TopCommands      []map[string]interface{} `json:"top_commands"`
	FirstCommandTime *time.Time               `json:"first_command_time,omitempty"`
	LastCommandTime  *time.Time               `json:"last_command_time,omitempty"`
}

// GetBasicStats returns basic usage statistics
func (s *StatsCalculator) GetBasicStats() (*BasicStats, error) {
	basicStats, err := s.db.GetBasicStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats from database: %w", err)
	}

	topCommands, err := s.db.GetTopCommands(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top commands: %w", err)
	}

	stats := &BasicStats{
		TotalCommands:  basicStats["total_commands"].(int),
		TotalSessions:  basicStats["total_sessions"].(int),
		UniqueCommands: basicStats["unique_commands"].(int),
		CommandsToday:  basicStats["commands_today"].(int),
		TopCommands:    topCommands,
	}

	// Set most used command if available
	if cmd, ok := basicStats["most_used_command"].(string); ok {
		stats.MostUsedCommand = cmd
		stats.MostUsedCount = basicStats["most_used_count"].(int)
	}

	// Get first and last command times
	firstTime, lastTime, err := s.getFirstAndLastCommandTimes()
	if err != nil {
		return nil, fmt.Errorf("failed to get command time range: %w", err)
	}

	stats.FirstCommandTime = firstTime
	stats.LastCommandTime = lastTime

	return stats, nil
}

// FormatBasicStats returns a formatted string representation of basic stats
func (s *StatsCalculator) FormatBasicStats(stats *BasicStats) string {
	var builder strings.Builder

	builder.WriteString("🚀 Termonaut Stats\n")
	builder.WriteString("─────────────────────────────────────\n")
	builder.WriteString(fmt.Sprintf("Total Commands: %d 🎯\n", stats.TotalCommands))
	builder.WriteString(fmt.Sprintf("Commands Today: %d 📅\n", stats.CommandsToday))
	builder.WriteString(fmt.Sprintf("Unique Commands: %d ⭐\n", stats.UniqueCommands))
	builder.WriteString(fmt.Sprintf("Terminal Sessions: %d 📱\n", stats.TotalSessions))

	if stats.MostUsedCommand != "" {
		builder.WriteString(fmt.Sprintf("Most Used: %s (%d times) 👑\n",
			stats.MostUsedCommand, stats.MostUsedCount))
	}

	if stats.FirstCommandTime != nil {
		daysSince := int(time.Since(*stats.FirstCommandTime).Hours() / 24)
		builder.WriteString(fmt.Sprintf("Using Termonaut for: %d days 📈\n", daysSince))
	}

	if len(stats.TopCommands) > 0 {
		builder.WriteString("\nTop Commands:\n")
		maxCommandLength := 0
		for i, cmd := range stats.TopCommands {
			if i >= 5 { // Show top 5
				break
			}
			cmdStr := cmd["command"].(string)
			if len(cmdStr) > maxCommandLength {
				maxCommandLength = len(cmdStr)
			}
		}

		for i, cmd := range stats.TopCommands {
			if i >= 5 { // Show top 5
				break
			}
			cmdStr := cmd["command"].(string)
			count := cmd["count"].(int)

			// Truncate long commands
			if len(cmdStr) > 30 {
				cmdStr = cmdStr[:27] + "..."
			}

			// Create a simple bar chart
			barLength := (count * 20) / stats.MostUsedCount
			if barLength < 1 {
				barLength = 1
			}
			bar := strings.Repeat("█", barLength)

			builder.WriteString(fmt.Sprintf("%-30s (%3d) %s\n", cmdStr, count, bar))
		}
	}

	return builder.String()
}

// GetTodayStats returns statistics for today
func (s *StatsCalculator) GetTodayStats() (map[string]interface{}, error) {
	// This is a placeholder for more detailed today stats
	// Will be expanded in later versions
	return s.db.GetBasicStats()
}

// getFirstAndLastCommandTimes gets the time range of all commands
func (s *StatsCalculator) getFirstAndLastCommandTimes() (*time.Time, *time.Time, error) {
	// This would query the database for MIN and MAX timestamps
	// For now, return nil to avoid errors
	return nil, nil, nil
}

// GetWeeklyStats returns statistics for the current week
func (s *StatsCalculator) GetWeeklyStats() (map[string]interface{}, error) {
	// Placeholder for weekly stats - will implement in Phase 2
	return map[string]interface{}{
		"commands_this_week": 0,
		"daily_average":      0,
	}, nil
}

// GetMonthlyStats returns statistics for the current month
func (s *StatsCalculator) GetMonthlyStats() (map[string]interface{}, error) {
	// Placeholder for monthly stats - will implement in Phase 2
	return map[string]interface{}{
		"commands_this_month": 0,
		"daily_average":       0,
	}, nil
}

// CalculateProductivityScore calculates a simple productivity score
func (s *StatsCalculator) CalculateProductivityScore(stats *BasicStats) int {
	// Simple scoring algorithm:
	// - 1 point per command today
	// - 2 points per unique command
	// - Bonus for consistency (placeholder)

	score := stats.CommandsToday + (stats.UniqueCommands * 2)

	// Cap at 100 for now
	if score > 100 {
		score = 100
	}

	return score
}
