package gamification

import (
	"fmt"
	"time"
)

// AchievementType represents different types of achievements
type AchievementType string

const (
	CommandCount    AchievementType = "command_count"
	UniqueCommands  AchievementType = "unique_commands"
	Streak          AchievementType = "streak"
	SessionCount    AchievementType = "session_count"
	TimeBasedUsage  AchievementType = "time_based"
	CategoryMastery AchievementType = "category_mastery"
	XPMilestone     AchievementType = "xp_milestone"
	LevelMilestone  AchievementType = "level_milestone"
)

// Achievement represents an achievement definition
type Achievement struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	Type        AchievementType `json:"type"`
	Target      int             `json:"target"`
	XPReward    int             `json:"xp_reward"`
	Hidden      bool            `json:"hidden"`
	Rare        bool            `json:"rare"`
}

// UserAchievement represents an earned achievement
type UserAchievement struct {
	Achievement *Achievement `json:"achievement"`
	EarnedAt    time.Time    `json:"earned_at"`
	Progress    int          `json:"progress"`
	Completed   bool         `json:"completed"`
}

// AchievementManager manages achievements
type AchievementManager struct {
	achievements map[string]*Achievement
}

// NewAchievementManager creates a new achievement manager
func NewAchievementManager() *AchievementManager {
	am := &AchievementManager{
		achievements: make(map[string]*Achievement),
	}
	am.loadDefaultAchievements()
	return am
}

// loadDefaultAchievements loads the default achievement set
func (am *AchievementManager) loadDefaultAchievements() {
	achievements := []*Achievement{
		// Command Count Achievements
		{
			ID:          "first_launch",
			Name:        "🚀 First Launch",
			Description: "Execute your first command with Termonaut",
			Icon:        "🚀",
			Type:        CommandCount,
			Target:      1,
			XPReward:    10,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "century",
			Name:        "🏆 Century",
			Description: "Execute 100 commands",
			Icon:        "🏆",
			Type:        CommandCount,
			Target:      100,
			XPReward:    50,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "millennium",
			Name:        "👑 Millennium",
			Description: "Execute 1,000 commands",
			Icon:        "👑",
			Type:        CommandCount,
			Target:      1000,
			XPReward:    200,
			Hidden:      false,
			Rare:        true,
		},
		{
			ID:          "commander",
			Name:        "⭐ Commander",
			Description: "Execute 10,000 commands",
			Icon:        "⭐",
			Type:        CommandCount,
			Target:      10000,
			XPReward:    500,
			Hidden:      false,
			Rare:        true,
		},

		// Unique Commands Achievements
		{
			ID:          "explorer",
			Name:        "🌟 Explorer",
			Description: "Use 10 different commands",
			Icon:        "🌟",
			Type:        UniqueCommands,
			Target:      10,
			XPReward:    25,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "versatile",
			Name:        "🎭 Versatile",
			Description: "Use 50 different commands",
			Icon:        "🎭",
			Type:        UniqueCommands,
			Target:      50,
			XPReward:    100,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "polyglot",
			Name:        "🗣️ Polyglot",
			Description: "Use 100 different commands",
			Icon:        "🗣️",
			Type:        UniqueCommands,
			Target:      100,
			XPReward:    250,
			Hidden:      false,
			Rare:        true,
		},

		// Streak Achievements
		{
			ID:          "consistent",
			Name:        "📈 Consistent",
			Description: "Maintain a 3-day streak",
			Icon:        "📈",
			Type:        Streak,
			Target:      3,
			XPReward:    30,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "dedicated",
			Name:        "🔥 Dedicated",
			Description: "Maintain a 7-day streak",
			Icon:        "🔥",
			Type:        Streak,
			Target:      7,
			XPReward:    75,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "unstoppable",
			Name:        "⚡ Unstoppable",
			Description: "Maintain a 30-day streak",
			Icon:        "⚡",
			Type:        Streak,
			Target:      30,
			XPReward:    300,
			Hidden:      false,
			Rare:        true,
		},

		// XP Milestones
		{
			ID:          "novice",
			Name:        "🌱 Novice",
			Description: "Earn 100 XP",
			Icon:        "🌱",
			Type:        XPMilestone,
			Target:      100,
			XPReward:    25,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "apprentice",
			Name:        "📚 Apprentice",
			Description: "Earn 500 XP",
			Icon:        "📚",
			Type:        XPMilestone,
			Target:      500,
			XPReward:    75,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "shell_sprinter",
			Name:        "🏃‍♂️ Shell Sprinter",
			Description: "Execute 100 commands in a single day",
			Icon:        "🏃‍♂️",
			Type:        CommandCount,
			Target:      100,
			XPReward:    150,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "config_whisperer",
			Name:        "🧙‍♂️ Config Whisperer",
			Description: "Edit configuration files 10 times",
			Icon:        "🧙‍♂️",
			Type:        CategoryMastery,
			Target:      10,
			XPReward:    75,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "night_coder",
			Name:        "🌙 Night Coder",
			Description: "Use shell between 0:00-5:00 AM for 50 commands",
			Icon:        "🌙",
			Type:        TimeBasedUsage,
			Target:      50,
			XPReward:    100,
			Hidden:      false,
			Rare:        true,
		},
		{
			ID:          "git_commander",
			Name:        "🧬 Git Commander",
			Description: "Use git commands 100 times",
			Icon:        "🧬",
			Type:        CategoryMastery,
			Target:      100,
			XPReward:    125,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "pro_streaker",
			Name:        "🔥 Pro Streaker",
			Description: "Maintain a 7-day terminal usage streak",
			Icon:        "🔥",
			Type:        Streak,
			Target:      7,
			XPReward:    200,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "sudo_smasher",
			Name:        "🛡️ Sudo Smasher",
			Description: "Use sudo commands 50 times",
			Icon:        "🛡️",
			Type:        CategoryMastery,
			Target:      50,
			XPReward:    100,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "docker_whale",
			Name:        "🐳 Docker Whale",
			Description: "Execute 25 Docker commands",
			Icon:        "🐳",
			Type:        CategoryMastery,
			Target:      25,
			XPReward:    75,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "vim_escape_artist",
			Name:        "🎭 Vim Escape Artist",
			Description: "Open Vim/Neovim 20 times",
			Icon:        "🎭",
			Type:        CategoryMastery,
			Target:      20,
			XPReward:    90,
			Hidden:      false,
			Rare:        true,
		},
		{
			ID:          "error_survivor",
			Name:        "💪 Error Survivor",
			Description: "Encounter 50 command failures and keep going",
			Icon:        "💪",
			Type:        CommandCount,
			Target:      50,
			XPReward:    80,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "early_bird",
			Name:        "🐦 Early Bird",
			Description: "Execute 30 commands between 5:00-8:00 AM",
			Icon:        "🐦",
			Type:        TimeBasedUsage,
			Target:      30,
			XPReward:    75,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "weekend_warrior",
			Name:        "⚔️ Weekend Warrior",
			Description: "Stay active on weekends for 4 weeks",
			Icon:        "⚔️",
			Type:        TimeBasedUsage,
			Target:      4,
			XPReward:    120,
			Hidden:      false,
			Rare:        true,
		},
		{
			ID:          "pipe_master",
			Name:        "🔗 Pipe Master",
			Description: "Use 25 commands with pipes (|)",
			Icon:        "🔗",
			Type:        CategoryMastery,
			Target:      25,
			XPReward:    85,
			Hidden:      false,
			Rare:        false,
		},

		{
			ID:          "expert",
			Name:        "🎓 Expert",
			Description: "Earn 2,000 XP",
			Icon:        "🎓",
			Type:        XPMilestone,
			Target:      2000,
			XPReward:    200,
			Hidden:      false,
			Rare:        true,
		},

		// Level Milestones
		{
			ID:          "cadet",
			Name:        "🚀 Space Cadet",
			Description: "Reach level 5",
			Icon:        "🚀",
			Type:        LevelMilestone,
			Target:      5,
			XPReward:    50,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "navigator",
			Name:        "🌟 Star Navigator",
			Description: "Reach level 10",
			Icon:        "🌟",
			Type:        LevelMilestone,
			Target:      10,
			XPReward:    100,
			Hidden:      false,
			Rare:        false,
		},
		{
			ID:          "commander_level",
			Name:        "⭐ Stellar Commander",
			Description: "Reach level 25",
			Icon:        "⭐",
			Type:        LevelMilestone,
			Target:      25,
			XPReward:    500,
			Hidden:      false,
			Rare:        true,
		},

		// Time-based Achievements
		{
			ID:          "early_bird",
			Name:        "🌅 Early Bird",
			Description: "Execute 50 commands before 8 AM",
			Icon:        "🌅",
			Type:        TimeBasedUsage,
			Target:      50,
			XPReward:    75,
			Hidden:      true,
			Rare:        false,
		},
		{
			ID:          "night_owl",
			Name:        "🦉 Night Owl",
			Description: "Execute 50 commands after 10 PM",
			Icon:        "🦉",
			Type:        TimeBasedUsage,
			Target:      50,
			XPReward:    75,
			Hidden:      true,
			Rare:        false,
		},
	}

	for _, achievement := range achievements {
		am.achievements[achievement.ID] = achievement
	}
}

// GetAllAchievements returns all available achievements
func (am *AchievementManager) GetAllAchievements() map[string]*Achievement {
	return am.achievements
}

// GetAchievement returns a specific achievement by ID
func (am *AchievementManager) GetAchievement(id string) (*Achievement, bool) {
	achievement, exists := am.achievements[id]
	return achievement, exists
}

// CheckAchievements checks which achievements should be unlocked based on user stats
func (am *AchievementManager) CheckAchievements(stats *UserStats, earnedAchievements map[string]*UserAchievement) []*UserAchievement {
	var newAchievements []*UserAchievement

	for _, achievement := range am.achievements {
		// Skip if already earned
		if _, earned := earnedAchievements[achievement.ID]; earned {
			continue
		}

		// Check if achievement should be unlocked
		progress := am.calculateProgress(achievement, stats)
		if progress >= achievement.Target {
			newAchievements = append(newAchievements, &UserAchievement{
				Achievement: achievement,
				EarnedAt:    time.Now(),
				Progress:    progress,
				Completed:   true,
			})
		}
	}

	return newAchievements
}

// calculateProgress calculates progress towards an achievement
func (am *AchievementManager) calculateProgress(achievement *Achievement, stats *UserStats) int {
	switch achievement.Type {
	case CommandCount:
		return stats.TotalCommands
	case UniqueCommands:
		return stats.UniqueCommands
	case Streak:
		return stats.LongestStreak
	case SessionCount:
		return stats.TotalSessions
	case XPMilestone:
		return stats.TotalXP
	case LevelMilestone:
		return stats.CurrentLevel
	case TimeBasedUsage:
		// This would require more complex tracking
		// For now, return 0 (will be implemented in later versions)
		return 0
	default:
		return 0
	}
}

// GetProgressToAchievement returns progress information for an achievement
func (am *AchievementManager) GetProgressToAchievement(achievementID string, stats *UserStats) (int, int, float64) {
	achievement, exists := am.achievements[achievementID]
	if !exists {
		return 0, 0, 0
	}

	progress := am.calculateProgress(achievement, stats)
	target := achievement.Target
	percentage := float64(progress) / float64(target) * 100

	if percentage > 100 {
		percentage = 100
	}

	return progress, target, percentage
}

// UserStats represents user statistics for achievement checking
type UserStats struct {
	TotalCommands     int `json:"total_commands"`
	UniqueCommands    int `json:"unique_commands"`
	TotalSessions     int `json:"total_sessions"`
	CurrentStreak     int `json:"current_streak"`
	LongestStreak     int `json:"longest_streak"`
	TotalXP           int `json:"total_xp"`
	CurrentLevel      int `json:"current_level"`
	CommandsToday     int `json:"commands_today"`
	EarlyBirdCommands int `json:"early_bird_commands"`
	NightOwlCommands  int `json:"night_owl_commands"`
}

// FormatAchievement formats an achievement for display
func FormatAchievement(achievement *Achievement, userAchievement *UserAchievement) string {
	if userAchievement != nil && userAchievement.Completed {
		return fmt.Sprintf("%s %s - Earned %s (+%d XP)",
			achievement.Icon, achievement.Name,
			userAchievement.EarnedAt.Format("Jan 2"),
			achievement.XPReward)
	}

	return fmt.Sprintf("%s %s - %s",
		achievement.Icon, achievement.Name, achievement.Description)
}
