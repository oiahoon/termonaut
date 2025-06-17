package display

import (
	"fmt"
	"strings"

	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// DisplayMode represents different display modes
type DisplayMode string

const (
	ModeMinimal DisplayMode = "minimal"
	ModeRich    DisplayMode = "rich"
	ModeQuiet   DisplayMode = "quiet" // For CI environments
)

// DisplayManager handles different display modes
type DisplayManager struct {
	mode            DisplayMode
	enableColors    bool
	enableEmojis    bool
	enableAnimation bool
}

// NewDisplayManager creates a new display manager
func NewDisplayManager(mode DisplayMode, ciEnvironment bool) *DisplayManager {
	// Auto-detect CI environment and switch to quiet mode
	if ciEnvironment {
		mode = ModeQuiet
	}

	return &DisplayManager{
		mode:            mode,
		enableColors:    mode != ModeQuiet,
		enableEmojis:    mode == ModeRich,
		enableAnimation: mode == ModeRich,
	}
}

// FormatStats formats stats according to the current display mode
func (dm *DisplayManager) FormatStats(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	switch dm.mode {
	case ModeMinimal:
		return dm.formatMinimalStats(basicStats, userProgress)
	case ModeRich:
		return dm.formatRichStats(basicStats, userProgress)
	case ModeQuiet:
		return dm.formatQuietStats(basicStats, userProgress)
	default:
		return dm.formatMinimalStats(basicStats, userProgress)
	}
}

// formatMinimalStats creates a minimal stats display
func (dm *DisplayManager) formatMinimalStats(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	var builder strings.Builder

	builder.WriteString("Termonaut Stats\n")
	builder.WriteString("───────────────\n")
	builder.WriteString(fmt.Sprintf("Level: %d | XP: %d\n", userProgress.CurrentLevel, userProgress.TotalXP))
	builder.WriteString(fmt.Sprintf("Commands: %d | Today: %d\n", basicStats.TotalCommands, basicStats.CommandsToday))
	builder.WriteString(fmt.Sprintf("Streak: %d days\n", userProgress.CurrentStreak))

	if basicStats.MostUsedCommand != "" {
		builder.WriteString(fmt.Sprintf("Top: %s (%d)\n", basicStats.MostUsedCommand, basicStats.MostUsedCount))
	}

	return builder.String()
}

// formatRichStats creates a rich, emoji-filled stats display
func (dm *DisplayManager) formatRichStats(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	var builder strings.Builder

	// Header with emojis and colors
	builder.WriteString("🚀 Termonaut Stats Dashboard\n")
	builder.WriteString("═════════════════════════════════════\n")

	// Level section with progress bar
	levelProgress := dm.calculateLevelProgress(userProgress)
	builder.WriteString(fmt.Sprintf("⚡ Level %d %s (%d XP)\n",
		userProgress.CurrentLevel, levelProgress, userProgress.TotalXP))

	// Commands section
	builder.WriteString(fmt.Sprintf("🎯 Commands: %d total | %d today\n",
		basicStats.TotalCommands, basicStats.CommandsToday))
	builder.WriteString(fmt.Sprintf("⭐ Unique Commands: %d\n", basicStats.UniqueCommands))

	// Streak section with fire emojis
	streakEmoji := dm.getStreakEmoji(userProgress.CurrentStreak)
	builder.WriteString(fmt.Sprintf("%s Streak: %d days (longest: %d)\n",
		streakEmoji, userProgress.CurrentStreak, userProgress.LongestStreak))

	// Most used command
	if basicStats.MostUsedCommand != "" {
		builder.WriteString(fmt.Sprintf("👑 Champion: %s (%d times)\n",
			basicStats.MostUsedCommand, basicStats.MostUsedCount))
	}

	// Top commands with bars
	if len(basicStats.TopCommands) > 0 {
		builder.WriteString("\n🔥 Top Commands:\n")
		for i, cmd := range basicStats.TopCommands {
			if i >= 3 { // Show top 3 in rich mode
				break
			}
			cmdStr := cmd["command"].(string)
			count := cmd["count"].(int)

			// Create visual bar
			barLength := (count * 15) / basicStats.MostUsedCount
			if barLength < 1 {
				barLength = 1
			}
			bar := strings.Repeat("█", barLength)

			builder.WriteString(fmt.Sprintf("%d. %-25s %s (%d)\n",
				i+1, dm.truncateCommand(cmdStr, 25), bar, count))
		}
	}

	// Footer
	builder.WriteString("\n💪 Keep coding, keep growing! 🌱\n")

	return builder.String()
}

// formatQuietStats creates a CI-friendly quiet output
func (dm *DisplayManager) formatQuietStats(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	// Minimal output for CI environments
	return fmt.Sprintf("Termonaut: Level %d, %d commands, %d day streak",
		userProgress.CurrentLevel, basicStats.TotalCommands, userProgress.CurrentStreak)
}

// calculateLevelProgress creates a visual progress bar for level progression
func (dm *DisplayManager) calculateLevelProgress(userProgress *models.UserProgress) string {
	if !dm.enableEmojis {
		return ""
	}

	currentLevelXP := dm.calculateXPForLevel(userProgress.CurrentLevel)
	nextLevelXP := dm.calculateXPForLevel(userProgress.CurrentLevel + 1)
	progressXP := userProgress.TotalXP - currentLevelXP
	neededXP := nextLevelXP - currentLevelXP

	if neededXP <= 0 {
		return "🌟✨🌟✨🌟"
	}

	progress := float64(progressXP) / float64(neededXP)
	barLength := int(progress * 10) // 10 character bar

	var bar strings.Builder
	for i := 0; i < 10; i++ {
		if i < barLength {
			bar.WriteString("█")
		} else {
			bar.WriteString("░")
		}
	}

	return fmt.Sprintf("[%s] %d/%d", bar.String(), progressXP, neededXP)
}

// calculateXPForLevel calculates XP required for a level
func (dm *DisplayManager) calculateXPForLevel(level int) int {
	return level * level * 100
}

// getStreakEmoji returns appropriate emoji for streak length
func (dm *DisplayManager) getStreakEmoji(streak int) string {
	if !dm.enableEmojis {
		return "Streak:"
	}

	switch {
	case streak >= 100:
		return "🔥🔥🔥"
	case streak >= 30:
		return "🔥🔥"
	case streak >= 7:
		return "🔥"
	case streak >= 3:
		return "⚡"
	case streak >= 1:
		return "✨"
	default:
		return "💤"
	}
}

// truncateCommand truncates command to specified length
func (dm *DisplayManager) truncateCommand(cmd string, maxLen int) string {
	if len(cmd) <= maxLen {
		return cmd
	}
	return cmd[:maxLen-3] + "..."
}

// FormatLevelUp formats level up notification
func (dm *DisplayManager) FormatLevelUp(newLevel int) string {
	switch dm.mode {
	case ModeQuiet:
		return fmt.Sprintf("Level up: %d", newLevel)
	case ModeMinimal:
		return fmt.Sprintf("Level Up! Now level %d", newLevel)
	case ModeRich:
		return fmt.Sprintf("🎉 LEVEL UP! 🎉\n⚡ You are now level %d! ⚡\n🌟 New powers unlocked! 🌟", newLevel)
	default:
		return fmt.Sprintf("Level Up! Now level %d", newLevel)
	}
}

// FormatAchievement formats achievement notification
func (dm *DisplayManager) FormatAchievement(name, description string) string {
	switch dm.mode {
	case ModeQuiet:
		return fmt.Sprintf("Achievement: %s", name)
	case ModeMinimal:
		return fmt.Sprintf("Achievement Unlocked: %s", name)
	case ModeRich:
		return fmt.Sprintf("🏆 ACHIEVEMENT UNLOCKED! 🏆\n🎯 %s\n📜 %s\n✨ Keep up the great work! ✨", name, description)
	default:
		return fmt.Sprintf("Achievement Unlocked: %s", name)
	}
}

// ShouldShowEasterEgg determines if easter eggs should be shown
func (dm *DisplayManager) ShouldShowEasterEgg() bool {
	return dm.mode == ModeRich
}

// GetMode returns the current display mode
func (dm *DisplayManager) GetMode() DisplayMode {
	return dm.mode
}

// IsQuietMode checks if in quiet mode (CI environment)
func (dm *DisplayManager) IsQuietMode() bool {
	return dm.mode == ModeQuiet
}
