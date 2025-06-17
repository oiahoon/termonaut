package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/oiahoon/termonaut/pkg/models"
)

// UpdateUserProgress updates user progress with XP and achievements
func (db *DB) UpdateUserProgress(xpGained int, newAchievements []*gamification.UserAchievement) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user progress
	if xpGained > 0 {
		err = db.addXPWithTransaction(tx, xpGained)
		if err != nil {
			return fmt.Errorf("failed to add XP: %w", err)
		}
	}

	// Add new achievements
	for _, achievement := range newAchievements {
		err = db.storeAchievementWithTransaction(tx, achievement)
		if err != nil {
			return fmt.Errorf("failed to store achievement: %w", err)
		}
	}

	return tx.Commit()
}

// addXPWithTransaction adds XP to user progress within a transaction
func (db *DB) addXPWithTransaction(tx *sql.Tx, xpGained int) error {
	query := `
		UPDATE user_progress
		SET total_xp = total_xp + ?,
		    current_level = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`

	// Calculate new level (we'll get the current XP first)
	var currentXP int
	err := tx.QueryRow("SELECT total_xp FROM user_progress WHERE id = 1").Scan(&currentXP)
	if err != nil {
		return fmt.Errorf("failed to get current XP: %w", err)
	}

	levelCalc := gamification.NewLevelCalculator()
	newLevel := levelCalc.CalculateLevel(currentXP + xpGained)

	_, err = tx.Exec(query, xpGained, newLevel)
	return err
}

// storeAchievementWithTransaction stores an achievement within a transaction
func (db *DB) storeAchievementWithTransaction(tx *sql.Tx, userAchievement *gamification.UserAchievement) error {
	query := `
		INSERT INTO achievements (id, name, description, earned_at, xp_bonus)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := tx.Exec(query,
		userAchievement.Achievement.ID,
		userAchievement.Achievement.Name,
		userAchievement.Achievement.Description,
		userAchievement.EarnedAt,
		userAchievement.Achievement.XPReward,
	)

	return err
}

// GetUserProgress returns current user progress
func (db *DB) GetUserProgress() (*models.UserProgress, error) {
	query := `
		SELECT id, total_xp, current_level, commands_count, unique_commands_count,
		       longest_streak, current_streak, last_activity_date, created_at, updated_at
		FROM user_progress
		WHERE id = 1
	`

	var progress models.UserProgress
	err := db.conn.QueryRow(query).Scan(
		&progress.ID, &progress.TotalXP, &progress.CurrentLevel,
		&progress.CommandsCount, &progress.UniqueCommandsCount,
		&progress.LongestStreak, &progress.CurrentStreak,
		&progress.LastActivityDate, &progress.CreatedAt, &progress.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	return &progress, nil
}

// GetUserAchievements returns all earned achievements
func (db *DB) GetUserAchievements() (map[string]*gamification.UserAchievement, error) {
	query := `
		SELECT id, name, description, earned_at, xp_bonus
		FROM achievements
		ORDER BY earned_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %w", err)
	}
	defer rows.Close()

	achievements := make(map[string]*gamification.UserAchievement)
	for rows.Next() {
		var id, name, description string
		var earnedAt time.Time
		var xpBonus int

		err := rows.Scan(&id, &name, &description, &earnedAt, &xpBonus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}

		achievement := &gamification.Achievement{
			ID:          id,
			Name:        name,
			Description: description,
			XPReward:    xpBonus,
		}

		userAchievement := &gamification.UserAchievement{
			Achievement: achievement,
			EarnedAt:    earnedAt,
			Completed:   true,
		}

		achievements[id] = userAchievement
	}

	return achievements, nil
}

// UpdateStreakAndCommands updates streak and command counts
func (db *DB) UpdateStreakAndCommands() error {
	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update commands count
	var totalCommands, uniqueCommands int
	err = tx.QueryRow("SELECT COUNT(*), COUNT(DISTINCT command) FROM commands").Scan(&totalCommands, &uniqueCommands)
	if err != nil {
		return fmt.Errorf("failed to get command counts: %w", err)
	}

	// Calculate streak
	currentStreak, longestStreak := db.calculateStreaks(tx)

	// Update user progress
	updateQuery := `
		UPDATE user_progress
		SET commands_count = ?,
		    unique_commands_count = ?,
		    current_streak = ?,
		    longest_streak = CASE WHEN ? > longest_streak THEN ? ELSE longest_streak END,
		    last_activity_date = DATE('now'),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`

	_, err = tx.Exec(updateQuery, totalCommands, uniqueCommands, currentStreak, longestStreak, longestStreak)
	if err != nil {
		return fmt.Errorf("failed to update user progress: %w", err)
	}

	return tx.Commit()
}

// calculateStreaks calculates current and longest streaks
func (db *DB) calculateStreaks(tx *sql.Tx) (int, int) {
	// Get command dates in descending order
	query := `
		SELECT DISTINCT DATE(timestamp) as cmd_date
		FROM commands
		ORDER BY cmd_date DESC
		LIMIT 365
	`

	rows, err := tx.Query(query)
	if err != nil {
		db.logger.Errorf("Failed to query command dates: %v", err)
		return 0, 0
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			continue
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0, 0
	}

	// Calculate current streak
	currentStreak := 1
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Check if user has activity today or yesterday
	if len(dates) > 0 && dates[0] != today && dates[0] != yesterday {
		currentStreak = 0
	} else {
		for i := 1; i < len(dates); i++ {
			prevDate, _ := time.Parse("2006-01-02", dates[i-1])
			currDate, _ := time.Parse("2006-01-02", dates[i])

			// Check if dates are consecutive
			if prevDate.AddDate(0, 0, -1).Format("2006-01-02") == currDate.Format("2006-01-02") {
				currentStreak++
			} else {
				break
			}
		}
	}

	// Calculate longest streak (simplified version)
	longestStreak := currentStreak
	tempStreak := 1

	for i := 1; i < len(dates); i++ {
		prevDate, _ := time.Parse("2006-01-02", dates[i-1])
		currDate, _ := time.Parse("2006-01-02", dates[i])

		if prevDate.AddDate(0, 0, -1).Format("2006-01-02") == currDate.Format("2006-01-02") {
			tempStreak++
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			tempStreak = 1
		}
	}

	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	return currentStreak, longestStreak
}

// GetGamificationStats returns stats needed for gamification calculations
func (db *DB) GetGamificationStats() (*gamification.UserStats, error) {
	progress, err := db.GetUserProgress()
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	basicStats, err := db.GetBasicStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	return &gamification.UserStats{
		TotalCommands:     progress.CommandsCount,
		UniqueCommands:    progress.UniqueCommandsCount,
		TotalSessions:     basicStats["total_sessions"].(int),
		CurrentStreak:     progress.CurrentStreak,
		LongestStreak:     progress.LongestStreak,
		TotalXP:           progress.TotalXP,
		CurrentLevel:      progress.CurrentLevel,
		CommandsToday:     basicStats["commands_today"].(int),
		EarlyBirdCommands: 0, // TODO: Implement time-based tracking
		NightOwlCommands:  0, // TODO: Implement time-based tracking
	}, nil
}

// StoreCommandWithXP stores a command and calculates XP
func (db *DB) StoreCommandWithXP(cmd *models.Command) error {
	// First store the command normally
	err := db.StoreCommand(cmd)
	if err != nil {
		return err
	}

	// Check if this is a new command
	isNewCommand, err := db.isNewCommand(cmd.Command)
	if err != nil {
		db.logger.Warnf("Failed to check if command is new: %v", err)
		isNewCommand = false
	}

	// Update streak and command counts
	err = db.UpdateStreakAndCommands()
	if err != nil {
		db.logger.Warnf("Failed to update streaks: %v", err)
	}

	// Classify command to get category and XP multiplier
	classifier := categories.NewCommandClassifier()
	category := classifier.ClassifyCommand(cmd.Command)
	categoryStr := string(category)

	// Calculate XP for this command
	xpCalc := gamification.NewXPCalculator(nil)
	stats, err := db.GetGamificationStats()
	if err != nil {
		db.logger.Warnf("Failed to get gamification stats: %v", err)
		return nil // Don't fail the whole operation
	}

	xpGained := xpCalc.CalculateCommandXP(cmd, isNewCommand, stats.CurrentStreak, categoryStr)

	// Check for new achievements
	achievementManager := gamification.NewAchievementManager()
	earnedAchievements, err := db.GetUserAchievements()
	if err != nil {
		db.logger.Warnf("Failed to get user achievements: %v", err)
		earnedAchievements = make(map[string]*gamification.UserAchievement)
	}

	newAchievements := achievementManager.CheckAchievements(stats, earnedAchievements)

	// Add XP from achievements
	totalXPGained := xpGained
	for _, achievement := range newAchievements {
		totalXPGained += achievement.Achievement.XPReward
	}

	// Update progress with XP and achievements
	err = db.UpdateUserProgress(totalXPGained, newAchievements)
	if err != nil {
		db.logger.Warnf("Failed to update user progress: %v", err)
	}

	return nil
}

// isNewCommand checks if this is the first time we've seen this command
func (db *DB) isNewCommand(command string) (bool, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM commands WHERE command = ?", command).Scan(&count)
	if err != nil {
		return false, err
	}
	return count <= 1, nil // <= 1 because we just inserted it
}
