package stats

import (
	"fmt"
	"strings"

	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/oiahoon/termonaut/pkg/models"
)

// GamificationStats represents gamification-enhanced statistics
type GamificationStats struct {
	BasicStats         *BasicStats                              `json:"basic_stats"`
	ProgressInfo       *gamification.ProgressInfo               `json:"progress_info"`
	Achievements       map[string]*gamification.UserAchievement `json:"achievements"`
	RecentAchievements []*gamification.UserAchievement          `json:"recent_achievements"`
	NextAchievements   []AchievementProgress                    `json:"next_achievements"`
}

// AchievementProgress represents progress towards an achievement
type AchievementProgress struct {
	Achievement *gamification.Achievement `json:"achievement"`
	Progress    int                       `json:"progress"`
	Target      int                       `json:"target"`
	Percentage  float64                   `json:"percentage"`
}

// GetGamificationStats returns comprehensive gamification statistics
func (s *StatsCalculator) GetGamificationStats() (*GamificationStats, error) {
	// Get basic stats
	basicStats, err := s.GetBasicStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// Get user progress from database
	userProgress, err := s.db.GetUserProgress()
	if err != nil {
		return nil, fmt.Errorf("failed to get user progress: %w", err)
	}

	// Calculate progress info
	progressCalc := gamification.NewProgressCalculator(nil)
	progressInfo := progressCalc.CalculateProgress(
		userProgress.TotalXP,
		basicStats.CommandsToday,
		userProgress.CurrentStreak,
	)

	// Get achievements
	achievements, err := s.db.GetUserAchievements()
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}

	// Get recent achievements (last 5)
	recentAchievements := s.getRecentAchievements(achievements, 5)

	// Get next achievements to unlock
	nextAchievements := s.getNextAchievements(achievements, userProgress)

	return &GamificationStats{
		BasicStats:         basicStats,
		ProgressInfo:       progressInfo,
		Achievements:       achievements,
		RecentAchievements: recentAchievements,
		NextAchievements:   nextAchievements,
	}, nil
}

// getRecentAchievements returns the most recently earned achievements
func (s *StatsCalculator) getRecentAchievements(achievements map[string]*gamification.UserAchievement, limit int) []*gamification.UserAchievement {
	var recent []*gamification.UserAchievement

	for _, achievement := range achievements {
		recent = append(recent, achievement)
	}

	// Sort by earned date (most recent first)
	for i := 0; i < len(recent)-1; i++ {
		for j := i + 1; j < len(recent); j++ {
			if recent[i].EarnedAt.Before(recent[j].EarnedAt) {
				recent[i], recent[j] = recent[j], recent[i]
			}
		}
	}

	// Limit results
	if len(recent) > limit {
		recent = recent[:limit]
	}

	return recent
}

// getNextAchievements returns achievements that are close to being unlocked
func (s *StatsCalculator) getNextAchievements(earned map[string]*gamification.UserAchievement, userProgress *models.UserProgress) []AchievementProgress {
	achievementManager := gamification.NewAchievementManager()
	allAchievements := achievementManager.GetAllAchievements()

	stats := &gamification.UserStats{
		TotalCommands:  userProgress.CommandsCount,
		UniqueCommands: userProgress.UniqueCommandsCount,
		CurrentStreak:  userProgress.CurrentStreak,
		LongestStreak:  userProgress.LongestStreak,
		TotalXP:        userProgress.TotalXP,
		CurrentLevel:   userProgress.CurrentLevel,
	}

	var nextAchievements []AchievementProgress

	for _, achievement := range allAchievements {
		// Skip if already earned
		if _, isEarned := earned[achievement.ID]; isEarned {
			continue
		}

		// Skip hidden achievements for now
		if achievement.Hidden {
			continue
		}

		progress, target, percentage := achievementManager.GetProgressToAchievement(achievement.ID, stats)

		// Include achievements that are at least 20% complete or very close to completion
		if percentage >= 20 || (target-progress) <= 5 {
			nextAchievements = append(nextAchievements, AchievementProgress{
				Achievement: achievement,
				Progress:    progress,
				Target:      target,
				Percentage:  percentage,
			})
		}
	}

	// Sort by completion percentage (highest first)
	for i := 0; i < len(nextAchievements)-1; i++ {
		for j := i + 1; j < len(nextAchievements); j++ {
			if nextAchievements[i].Percentage < nextAchievements[j].Percentage {
				nextAchievements[i], nextAchievements[j] = nextAchievements[j], nextAchievements[i]
			}
		}
	}

	// Limit to top 5
	if len(nextAchievements) > 5 {
		nextAchievements = nextAchievements[:5]
	}

	return nextAchievements
}

// FormatGamificationStats returns a formatted string representation of gamification stats
func (s *StatsCalculator) FormatGamificationStats(stats *GamificationStats) string {
	var builder strings.Builder

	// Header with level and XP
	builder.WriteString("🎮 Termonaut Gamification Dashboard\n")
	builder.WriteString("═══════════════════════════════════════════════════════\n")

	// Level and Progress
	builder.WriteString(fmt.Sprintf("🌟 Level %d - %s\n",
		stats.ProgressInfo.CurrentLevel, stats.ProgressInfo.LevelTitle))

	// XP Progress Bar
	progressBar := s.createProgressBar(stats.ProgressInfo.ProgressPercent, 30)
	builder.WriteString(fmt.Sprintf("XP: %d/%d %s (%.1f%%)\n",
		stats.ProgressInfo.XPInCurrentLevel,
		stats.ProgressInfo.XPForCurrentLevel,
		progressBar,
		stats.ProgressInfo.ProgressPercent))

	builder.WriteString(fmt.Sprintf("Total XP: %d 💎 | Next Level: %d XP needed\n\n",
		stats.ProgressInfo.TotalXP, stats.ProgressInfo.XPToNextLevel))

	// Current Stats Overview
	builder.WriteString("📊 Current Stats:\n")
	builder.WriteString(fmt.Sprintf("Commands: %d 🎯 | Today: %d 📅 | Unique: %d ⭐\n",
		stats.BasicStats.TotalCommands,
		stats.BasicStats.CommandsToday,
		stats.BasicStats.UniqueCommands))

	builder.WriteString(fmt.Sprintf("Current Streak: %d 🔥 | Best Streak: %d 🏆\n\n",
		stats.ProgressInfo.CurrentStreak,
		// Note: We'd need to add LongestStreak to ProgressInfo or get it separately
		stats.ProgressInfo.CurrentStreak)) // Temporary

	// Recent Achievements
	if len(stats.RecentAchievements) > 0 {
		builder.WriteString("🏆 Recent Achievements:\n")
		for _, achievement := range stats.RecentAchievements {
			builder.WriteString(fmt.Sprintf("  %s %s (+%d XP) - %s\n",
				achievement.Achievement.Icon,
				achievement.Achievement.Name,
				achievement.Achievement.XPReward,
				achievement.EarnedAt.Format("Jan 2")))
		}
		builder.WriteString("\n")
	}

	// Next Achievements
	if len(stats.NextAchievements) > 0 {
		builder.WriteString("🎯 Progress Towards Next Achievements:\n")
		for _, next := range stats.NextAchievements {
			progressBar := s.createProgressBar(next.Percentage, 20)
			builder.WriteString(fmt.Sprintf("  %s %s %s (%.0f%%) - %d/%d\n",
				next.Achievement.Icon,
				next.Achievement.Name,
				progressBar,
				next.Percentage,
				next.Progress,
				next.Target))
		}
		builder.WriteString("\n")
	}

	// Achievement Summary
	totalAchievements := len(stats.Achievements)
	achievementManager := gamification.NewAchievementManager()
	allAchievements := achievementManager.GetAllAchievements()
	maxAchievements := 0
	for _, achievement := range allAchievements {
		if !achievement.Hidden {
			maxAchievements++
		}
	}

	builder.WriteString(fmt.Sprintf("🎖️ Achievements Unlocked: %d/%d\n", totalAchievements, maxAchievements))

	return builder.String()
}

// createProgressBar creates a visual progress bar
func (s *StatsCalculator) createProgressBar(percentage float64, width int) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}

	filled := int(percentage * float64(width) / 100)
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s]", bar)
}

// FormatLevelUpNotification formats a level up notification
func FormatLevelUpNotification(oldLevel, newLevel int, newTitle string) string {
	return fmt.Sprintf(`
🎉 LEVEL UP! 🎉
🌟 You've reached Level %d! 🌟
%s

Keep exploring the terminal universe!
`, newLevel, newTitle)
}

// FormatAchievementUnlock formats an achievement unlock notification
func FormatAchievementUnlock(achievement *gamification.UserAchievement) string {
	rarity := ""
	if achievement.Achievement.Rare {
		rarity = " ⭐ RARE"
	}

	return fmt.Sprintf(`
🏆 ACHIEVEMENT UNLOCKED!%s 🏆
%s %s
%s
+%d XP Earned!
`, rarity, achievement.Achievement.Icon, achievement.Achievement.Name,
		achievement.Achievement.Description, achievement.Achievement.XPReward)
}
