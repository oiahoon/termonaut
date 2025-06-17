package github

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// BadgeConfig represents configuration for generating badges
type BadgeConfig struct {
	Style     string // flat, flat-square, plastic, for-the-badge, social
	Color     string // Color for the badge
	Logo      string // Logo to include
	LogoColor string // Logo color
}

// DefaultBadgeConfig returns the default badge configuration
func DefaultBadgeConfig() BadgeConfig {
	return BadgeConfig{
		Style:     "flat-square",
		Color:     "brightgreen",
		Logo:      "terminal",
		LogoColor: "white",
	}
}

// BadgeGenerator handles generation of various GitHub badges
type BadgeGenerator struct {
	config BadgeConfig
	stats  *stats.StatsCalculator
}

// NewBadgeGenerator creates a new badge generator
func NewBadgeGenerator(statsCalculator *stats.StatsCalculator, config BadgeConfig) *BadgeGenerator {
	return &BadgeGenerator{
		config: config,
		stats:  statsCalculator,
	}
}

// GenerateXPBadge generates an XP progress badge
func (bg *BadgeGenerator) GenerateXPBadge(userProgress *models.UserProgress) string {
	level := userProgress.CurrentLevel
	currentXP := userProgress.TotalXP
	nextLevelXP := bg.calculateXPForLevel(level + 1)

	label := "XP"
	message := fmt.Sprintf("Level %d (%d/%d)", level, currentXP, nextLevelXP)

	return bg.generateShieldsURL(label, message, bg.GetXPColor(level))
}

// calculateXPForLevel calculates XP required for a given level
func (bg *BadgeGenerator) calculateXPForLevel(level int) int {
	// Simple XP calculation: level^2 * 100
	return level * level * 100
}

// GenerateCommandsBadge generates a total commands badge
func (bg *BadgeGenerator) GenerateCommandsBadge(totalCommands int) string {
	label := "Commands"
	message := fmt.Sprintf("%d", totalCommands)

	color := bg.GetCommandsColor(totalCommands)
	return bg.generateShieldsURL(label, message, color)
}

// GenerateStreakBadge generates a streak badge
func (bg *BadgeGenerator) GenerateStreakBadge(streakDays int) string {
	label := "Streak"
	message := fmt.Sprintf("%d days", streakDays)

	color := bg.GetStreakColor(streakDays)
	return bg.generateShieldsURL(label, message, color)
}

// GenerateProductivityBadge generates a productivity score badge
func (bg *BadgeGenerator) GenerateProductivityBadge(score float64) string {
	label := "Productivity"
	message := fmt.Sprintf("%.1f%%", score*100)

	color := bg.GetProductivityColor(score)
	return bg.generateShieldsURL(label, message, color)
}

// GenerateAchievementsBadge generates an achievements progress badge
func (bg *BadgeGenerator) GenerateAchievementsBadge(completed, total int) string {
	label := "Achievements"
	message := fmt.Sprintf("%d/%d", completed, total)

	percentage := float64(completed) / float64(total)
	color := bg.GetAchievementColor(percentage)
	return bg.generateShieldsURL(label, message, color)
}

// GenerateLastActiveBadge generates a last active badge
func (bg *BadgeGenerator) GenerateLastActiveBadge(lastActive time.Time) string {
	label := "Last Active"

	now := time.Now()
	diff := now.Sub(lastActive)

	var message string
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		message = fmt.Sprintf("%dm ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		message = fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		message = fmt.Sprintf("%dd ago", days)
	}

	color := bg.GetActivityColor(diff)
	return bg.generateShieldsURL(label, message, color)
}

// GenerateCustomBadge generates a custom badge with provided label and message
func (bg *BadgeGenerator) GenerateCustomBadge(label, message, color string) string {
	return bg.generateShieldsURL(label, message, color)
}

// generateShieldsURL generates a Shields.io badge URL
func (bg *BadgeGenerator) generateShieldsURL(label, message, color string) string {
	baseURL := "https://img.shields.io/badge"

	// URL encode the label and message
	encodedLabel := url.QueryEscape(label)
	encodedMessage := url.QueryEscape(message)
	encodedColor := url.QueryEscape(color)

	// Build the badge URL
	badgeURL := fmt.Sprintf("%s/%s-%s-%s", baseURL, encodedLabel, encodedMessage, encodedColor)

	// Add style parameter
	if bg.config.Style != "" {
		badgeURL += "?style=" + bg.config.Style
	}

	// Add logo if specified
	if bg.config.Logo != "" {
		separator := "?"
		if strings.Contains(badgeURL, "?") {
			separator = "&"
		}
		badgeURL += separator + "logo=" + url.QueryEscape(bg.config.Logo)

		if bg.config.LogoColor != "" {
			badgeURL += "&logoColor=" + url.QueryEscape(bg.config.LogoColor)
		}
	}

	return badgeURL
}

// Color determination methods
func (bg *BadgeGenerator) GetXPColor(level int) string {
	switch {
	case level >= 50:
		return "gold"
	case level >= 25:
		return "orange"
	case level >= 10:
		return "blue"
	case level >= 5:
		return "green"
	default:
		return "lightgrey"
	}
}

func (bg *BadgeGenerator) GetCommandsColor(commands int) string {
	switch {
	case commands >= 10000:
		return "gold"
	case commands >= 5000:
		return "orange"
	case commands >= 1000:
		return "blue"
	case commands >= 100:
		return "green"
	default:
		return "lightgrey"
	}
}

func (bg *BadgeGenerator) GetStreakColor(days int) string {
	switch {
	case days >= 100:
		return "gold"
	case days >= 30:
		return "orange"
	case days >= 7:
		return "blue"
	case days >= 3:
		return "green"
	default:
		return "red"
	}
}

func (bg *BadgeGenerator) GetProductivityColor(score float64) string {
	switch {
	case score >= 0.9:
		return "brightgreen"
	case score >= 0.7:
		return "green"
	case score >= 0.5:
		return "yellowgreen"
	case score >= 0.3:
		return "orange"
	default:
		return "red"
	}
}

func (bg *BadgeGenerator) GetAchievementColor(percentage float64) string {
	switch {
	case percentage >= 0.9:
		return "gold"
	case percentage >= 0.7:
		return "orange"
	case percentage >= 0.5:
		return "blue"
	case percentage >= 0.3:
		return "green"
	default:
		return "lightgrey"
	}
}

func (bg *BadgeGenerator) GetActivityColor(lastActiveDiff time.Duration) string {
	switch {
	case lastActiveDiff < time.Hour:
		return "brightgreen"
	case lastActiveDiff < 6*time.Hour:
		return "green"
	case lastActiveDiff < 24*time.Hour:
		return "yellow"
	case lastActiveDiff < 7*24*time.Hour:
		return "orange"
	default:
		return "red"
	}
}
