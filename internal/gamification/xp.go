package gamification

import (
	"math"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/pkg/models"
)

// XPConfig holds XP calculation configuration
type XPConfig struct {
	BaseXPPerCommand    int                `json:"base_xp_per_command"`
	NewCommandBonus     int                `json:"new_command_bonus"`
	StreakMultiplier    float64            `json:"streak_multiplier"`
	CategoryMultipliers map[string]float64 `json:"category_multipliers"`
	TimeDecayEnabled    bool               `json:"time_decay_enabled"`
	MaxDailyXP          int                `json:"max_daily_xp"`
}

// DefaultXPConfig returns the default XP configuration
func DefaultXPConfig() *XPConfig {
	return &XPConfig{
		BaseXPPerCommand: 1,
		NewCommandBonus:  5,
		StreakMultiplier: 1.2,
		CategoryMultipliers: map[string]float64{
			"git":         1.5,
			"development": 1.3,
			"system":      1.0,
			"navigation":  0.8,
			"unknown":     1.0,
		},
		TimeDecayEnabled: false,
		MaxDailyXP:       1000,
	}
}

// XPCalculator handles XP calculations
type XPCalculator struct {
	config *XPConfig
}

// NewXPCalculator creates a new XP calculator
func NewXPCalculator(config *XPConfig) *XPCalculator {
	if config == nil {
		config = DefaultXPConfig()
	}
	return &XPCalculator{
		config: config,
	}
}

// CalculateCommandXP calculates XP for a single command
func (xp *XPCalculator) CalculateCommandXP(cmd *models.Command, isNewCommand bool, streak int, category string) int {
	baseXP := float64(xp.config.BaseXPPerCommand)

	// Apply failure penalty first
	if cmd.ExitCode != 0 {
		failurePenalty := xp.calculateFailurePenalty(cmd.ExitCode, category)
		baseXP *= failurePenalty
	}

	// New command bonus (only for successful commands)
	if isNewCommand && cmd.ExitCode == 0 {
		baseXP += float64(xp.config.NewCommandBonus)
	}

	// Category multiplier
	if multiplier, exists := xp.config.CategoryMultipliers[category]; exists {
		baseXP *= multiplier
	}

	// Command complexity bonus
	complexityBonus := xp.calculateComplexityBonus(cmd.Command, category)
	baseXP *= complexityBonus

	// Streak multiplier
	if streak > 1 {
		streakBonus := math.Pow(xp.config.StreakMultiplier, float64(streak-1))
		baseXP *= streakBonus
	}

	// Time-based bonus (morning/evening productivity)
	if xp.config.TimeDecayEnabled {
		timeBonus := xp.calculateTimeBonus(cmd.Timestamp)
		baseXP *= timeBonus
	}

	// Round and ensure minimum XP (even failures give some XP for learning)
	result := int(math.Round(baseXP))
	if result < 1 && cmd.ExitCode == 0 {
		result = 1
	} else if result < 0 && cmd.ExitCode != 0 {
		result = 0 // Failed commands can give 0 XP but not negative
	}

	return result
}

// calculateFailurePenalty calculates XP penalty for failed commands
func (xp *XPCalculator) calculateFailurePenalty(exitCode int, category string) float64 {
	// Different penalties based on exit code severity
	var basePenalty float64
	switch {
	case exitCode == 1:
		basePenalty = 0.7 // Minor failure - 30% penalty
	case exitCode == 2:
		basePenalty = 0.5 // Misuse of command - 50% penalty  
	case exitCode >= 126 && exitCode <= 128:
		basePenalty = 0.3 // Command not found/executable - 70% penalty
	case exitCode == 130:
		basePenalty = 0.9 // Ctrl+C (user interrupted) - minor penalty
	default:
		basePenalty = 0.6 // Other failures - 40% penalty
	}

	// Less penalty for development/learning categories
	if category == "development" || category == "git" || category == "test" {
		basePenalty += 0.2 // Reduce penalty for learning environments
		if basePenalty > 1.0 {
			basePenalty = 1.0
		}
	}

	return basePenalty
}

// calculateComplexityBonus calculates bonus based on command complexity
func (xp *XPCalculator) calculateComplexityBonus(command string, category string) float64 {
	parts := strings.Fields(command)
	argCount := len(parts) - 1 // Exclude the base command

	// Base complexity from argument count
	complexityBonus := 1.0 + float64(argCount)*0.05 // 5% per argument

	// Pipe complexity bonus
	pipeCount := strings.Count(command, "|")
	if pipeCount > 0 {
		complexityBonus += float64(pipeCount) * 0.15 // 15% per pipe
	}

	// Redirection complexity bonus
	redirectCount := strings.Count(command, ">") + strings.Count(command, "<")
	if redirectCount > 0 {
		complexityBonus += float64(redirectCount) * 0.1 // 10% per redirection
	}

	// Advanced syntax bonuses
	if strings.Contains(command, "&&") || strings.Contains(command, "||") {
		complexityBonus += 0.2 // 20% for logical operators
	}

	if strings.Contains(command, "$(") || strings.Contains(command, "`") {
		complexityBonus += 0.25 // 25% for command substitution
	}

	// Regex/pattern bonuses
	if strings.Contains(command, "*") || strings.Contains(command, "?") || 
	   strings.Contains(command, "[") || strings.Contains(command, "{") {
		complexityBonus += 0.15 // 15% for wildcards and patterns
	}

	// Category-specific complexity adjustments
	switch category {
	case "development", "security", "database":
		complexityBonus *= 1.1 // 10% bonus for complex categories
	case "navigation", "system":
		complexityBonus *= 0.95 // 5% reduction for basic categories
	}

	// Cap the bonus to prevent excessive XP
	if complexityBonus > 2.0 {
		complexityBonus = 2.0
	}

	return complexityBonus
}

// calculateTimeBonus calculates time-based XP bonus
func (xp *XPCalculator) calculateTimeBonus(timestamp time.Time) float64 {
	hour := timestamp.Hour()

	// Morning productivity bonus (6-10 AM)
	if hour >= 6 && hour <= 10 {
		return 1.2
	}

	// Evening focus bonus (7-11 PM)
	if hour >= 19 && hour <= 23 {
		return 1.1
	}

	// Late night penalty (12-5 AM)
	if hour >= 0 && hour <= 5 {
		return 0.8
	}

	// Normal hours
	return 1.0
}

// LevelCalculator handles level calculations
type LevelCalculator struct{}

// NewLevelCalculator creates a new level calculator
func NewLevelCalculator() *LevelCalculator {
	return &LevelCalculator{}
}

// CalculateLevel calculates level from total XP
func (lc *LevelCalculator) CalculateLevel(totalXP int) int {
	if totalXP <= 0 {
		return 1
	}

	// Level formula: level = sqrt(totalXP / 100) + 1
	// This creates a gentle curve where higher levels require more XP
	level := int(math.Sqrt(float64(totalXP)/100.0)) + 1

	// Cap at reasonable level
	if level > 100 {
		level = 100
	}

	return level
}

// CalculateXPForLevel calculates XP required for a specific level
func (lc *LevelCalculator) CalculateXPForLevel(level int) int {
	if level <= 1 {
		return 0
	}

	// Inverse of level formula: XP = (level - 1)^2 * 100
	return (level - 1) * (level - 1) * 100
}

// CalculateXPToNextLevel calculates XP needed to reach next level
func (lc *LevelCalculator) CalculateXPToNextLevel(currentXP int) (int, int, int) {
	currentLevel := lc.CalculateLevel(currentXP)
	nextLevel := currentLevel + 1

	xpForCurrentLevel := lc.CalculateXPForLevel(currentLevel)
	xpForNextLevel := lc.CalculateXPForLevel(nextLevel)

	xpInCurrentLevel := currentXP - xpForCurrentLevel
	xpNeeded := xpForNextLevel - currentXP

	return xpInCurrentLevel, xpNeeded, xpForNextLevel - xpForCurrentLevel
}

// GetLevelTitle returns a themed title for the level
func (lc *LevelCalculator) GetLevelTitle(level int) string {
	titles := map[int]string{
		1:   "🌱 Rookie Explorer",
		5:   "🚀 Space Cadet",
		10:  "🌟 Star Navigator",
		15:  "🛸 Cosmic Pilot",
		20:  "🌌 Galaxy Ranger",
		25:  "⭐ Stellar Commander",
		30:  "🌠 Nebula Master",
		35:  "🪐 Planet Walker",
		40:  "☄️ Comet Rider",
		45:  "🌙 Lunar Guardian",
		50:  "☀️ Solar Champion",
		60:  "🌈 Aurora Seeker",
		70:  "⚡ Quantum Navigator",
		80:  "🔮 Cosmic Sage",
		90:  "👑 Universal Master",
		100: "🎆 Legendary Termonaut",
	}

	// Find the highest title that applies
	var title string = "🌱 Rookie Explorer"
	for levelReq, levelTitle := range titles {
		if level >= levelReq {
			title = levelTitle
		}
	}

	return title
}

// ProgressCalculator handles progress tracking
type ProgressCalculator struct {
	xpCalc    *XPCalculator
	levelCalc *LevelCalculator
}

// NewProgressCalculator creates a new progress calculator
func NewProgressCalculator(xpConfig *XPConfig) *ProgressCalculator {
	return &ProgressCalculator{
		xpCalc:    NewXPCalculator(xpConfig),
		levelCalc: NewLevelCalculator(),
	}
}

// CalculateProgress calculates comprehensive progress information
func (pc *ProgressCalculator) CalculateProgress(totalXP int, commandsToday int, streak int) *ProgressInfo {
	currentLevel := pc.levelCalc.CalculateLevel(totalXP)
	xpInLevel, xpToNext, xpForLevel := pc.levelCalc.CalculateXPToNextLevel(totalXP)

	progressPercent := float64(xpInLevel) / float64(xpForLevel) * 100
	if progressPercent > 100 {
		progressPercent = 100
	}

	return &ProgressInfo{
		CurrentLevel:      currentLevel,
		LevelTitle:        pc.levelCalc.GetLevelTitle(currentLevel),
		TotalXP:           totalXP,
		XPInCurrentLevel:  xpInLevel,
		XPToNextLevel:     xpToNext,
		XPForCurrentLevel: xpForLevel,
		ProgressPercent:   progressPercent,
		CommandsToday:     commandsToday,
		CurrentStreak:     streak,
	}
}

// ProgressInfo holds comprehensive progress information
type ProgressInfo struct {
	CurrentLevel      int     `json:"current_level"`
	LevelTitle        string  `json:"level_title"`
	TotalXP           int     `json:"total_xp"`
	XPInCurrentLevel  int     `json:"xp_in_current_level"`
	XPToNextLevel     int     `json:"xp_to_next_level"`
	XPForCurrentLevel int     `json:"xp_for_current_level"`
	ProgressPercent   float64 `json:"progress_percent"`
	CommandsToday     int     `json:"commands_today"`
	CurrentStreak     int     `json:"current_streak"`
}
