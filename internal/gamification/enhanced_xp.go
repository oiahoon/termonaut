package gamification

import (
	"math"
	"strings"
)

// EnhancedLevelCalculator provides more challenging level progression
type EnhancedLevelCalculator struct {
	baseXP      int     // Base XP for level 2
	exponent    float64 // Exponential growth factor
	linearBonus int     // Linear bonus per level
}

// NewEnhancedLevelCalculator creates a new enhanced level calculator
func NewEnhancedLevelCalculator() *EnhancedLevelCalculator {
	return &EnhancedLevelCalculator{
		baseXP:      150,   // Increased from 100
		exponent:    1.15,  // Exponential growth
		linearBonus: 50,    // Additional linear growth
	}
}

// CalculateLevel calculates level from total XP with enhanced difficulty
func (elc *EnhancedLevelCalculator) CalculateLevel(totalXP int) int {
	if totalXP <= 0 {
		return 1
	}

	// Enhanced formula with exponential + linear growth
	// This makes higher levels significantly more challenging
	level := 1
	
	for level < 200 { // Increased level cap
		nextLevelXP := elc.CalculateXPForLevel(level + 1)
		if totalXP < nextLevelXP {
			break
		}
		level++
	}

	return level
}

// CalculateXPForLevel calculates XP required for a specific level
func (elc *EnhancedLevelCalculator) CalculateXPForLevel(level int) int {
	if level <= 1 {
		return 0
	}

	// Enhanced formula: baseXP * (level-1)^exponent + linearBonus * (level-1)
	// This creates a curve that starts gentle but becomes steep
	exponentialPart := float64(elc.baseXP) * math.Pow(float64(level-1), elc.exponent)
	linearPart := float64(elc.linearBonus) * float64(level-1)
	
	return int(exponentialPart + linearPart)
}

// CalculateXPToNextLevel calculates XP needed to reach next level
func (elc *EnhancedLevelCalculator) CalculateXPToNextLevel(currentXP int) (int, int, int) {
	currentLevel := elc.CalculateLevel(currentXP)
	currentLevelXP := elc.CalculateXPForLevel(currentLevel)
	nextLevelXP := elc.CalculateXPForLevel(currentLevel + 1)
	
	progressXP := currentXP - currentLevelXP
	neededXP := nextLevelXP - currentLevelXP
	remainingXP := nextLevelXP - currentXP
	
	return progressXP, neededXP, remainingXP
}

// GetLevelTitle returns a title based on level with more granular progression
func (elc *EnhancedLevelCalculator) GetLevelTitle(level int) string {
	switch {
	case level >= 100:
		return "🌌 Cosmic Legend"
	case level >= 90:
		return "🌟 Stellar Master"
	case level >= 80:
		return "🚀 Galactic Commander"
	case level >= 70:
		return "🛸 Space Admiral"
	case level >= 60:
		return "👨‍🚀 Veteran Astronaut"
	case level >= 50:
		return "🪐 Planetary Explorer"
	case level >= 40:
		return "🌙 Lunar Specialist"
	case level >= 30:
		return "🛰️ Orbital Engineer"
	case level >= 25:
		return "🚀 Senior Pilot"
	case level >= 20:
		return "👨‍🚀 Space Commander"
	case level >= 15:
		return "🌟 Star Navigator"
	case level >= 10:
		return "🚀 Rocket Pilot"
	case level >= 8:
		return "👨‍🚀 Astronaut"
	case level >= 6:
		return "🛸 Space Cadet"
	case level >= 4:
		return "🌙 Moon Walker"
	case level >= 2:
		return "🚀 Space Trainee"
	default:
		return "🌍 Earth Dweller"
	}
}

// EnhancedXPConfig provides more sophisticated XP calculation
type EnhancedXPConfig struct {
	BaseXPPerCommand    int                `json:"base_xp_per_command"`
	NewCommandBonus     int                `json:"new_command_bonus"`
	StreakMultiplier    float64            `json:"streak_multiplier"`
	CategoryMultipliers map[string]float64 `json:"category_multipliers"`
	TimeDecayEnabled    bool               `json:"time_decay_enabled"`
	MaxDailyXP          int                `json:"max_daily_xp"`
	
	// Enhanced features
	ComplexityBonus     map[string]int     `json:"complexity_bonus"`
	TimeOfDayMultiplier map[int]float64    `json:"time_of_day_multiplier"`
	WeekdayMultiplier   map[int]float64    `json:"weekday_multiplier"`
	ConsistencyBonus    float64            `json:"consistency_bonus"`
}

// DefaultEnhancedXPConfig returns enhanced XP configuration
func DefaultEnhancedXPConfig() *EnhancedXPConfig {
	return &EnhancedXPConfig{
		BaseXPPerCommand: 2, // Increased base XP
		NewCommandBonus:  8, // Increased new command bonus
		StreakMultiplier: 1.3, // Increased streak multiplier
		CategoryMultipliers: map[string]float64{
			"git":         1.8, // Increased git bonus
			"development": 1.6, // Increased development bonus
			"docker":      1.7, // Docker gets good bonus
			"kubernetes":  2.0, // K8s gets highest bonus
			"testing":     1.5, // Testing bonus
			"deployment":  1.8, // Deployment bonus
			"system":      1.2, // Slightly increased system
			"navigation":  0.9, // Slightly decreased navigation
			"unknown":     1.0,
		},
		TimeDecayEnabled: false,
		MaxDailyXP:       2000, // Increased daily cap
		
		// Enhanced features
		ComplexityBonus: map[string]int{
			"pipe":        2, // Commands with pipes
			"redirect":    2, // Commands with redirects
			"background":  3, // Background processes
			"sudo":        3, // Sudo commands
			"ssh":         4, // Remote commands
			"complex":     5, // Very complex commands
		},
		TimeOfDayMultiplier: map[int]float64{
			// Early morning bonus
			5: 1.2, 6: 1.2, 7: 1.1,
			// Late night penalty (encourage healthy habits)
			23: 0.9, 0: 0.8, 1: 0.7, 2: 0.6,
		},
		WeekdayMultiplier: map[int]float64{
			// Monday motivation bonus
			1: 1.1,
			// Weekend coding bonus
			0: 1.2, 6: 1.2,
		},
		ConsistencyBonus: 1.5, // Bonus for consistent daily usage
	}
}

// CalculateComplexityBonus calculates bonus XP based on command complexity
func (config *EnhancedXPConfig) CalculateComplexityBonus(command string) int {
	bonus := 0
	
	// Check for various complexity indicators
	if strings.Contains(command, "|") {
		bonus += config.ComplexityBonus["pipe"]
	}
	if strings.Contains(command, ">") || strings.Contains(command, ">>") {
		bonus += config.ComplexityBonus["redirect"]
	}
	if strings.Contains(command, "&") {
		bonus += config.ComplexityBonus["background"]
	}
	if strings.HasPrefix(command, "sudo ") {
		bonus += config.ComplexityBonus["sudo"]
	}
	if strings.Contains(command, "ssh ") {
		bonus += config.ComplexityBonus["ssh"]
	}
	
	// Complex command detection (multiple operators, long commands, etc.)
	if len(command) > 100 || strings.Count(command, " ") > 10 {
		bonus += config.ComplexityBonus["complex"]
	}
	
	return bonus
}

// GetLevelProgressInfo returns detailed level progress information
func (elc *EnhancedLevelCalculator) GetLevelProgressInfo(currentXP int) map[string]interface{} {
	currentLevel := elc.CalculateLevel(currentXP)
	progressXP, neededXP, remainingXP := elc.CalculateXPToNextLevel(currentXP)
	
	// Calculate progress percentage
	progressPercent := float64(progressXP) / float64(neededXP) * 100.0
	
	// Calculate XP for next few levels
	nextLevels := make([]map[string]interface{}, 0)
	for i := 1; i <= 3; i++ {
		levelXP := elc.CalculateXPForLevel(currentLevel + i)
		nextLevels = append(nextLevels, map[string]interface{}{
			"level": currentLevel + i,
			"title": elc.GetLevelTitle(currentLevel + i),
			"total_xp": levelXP,
			"xp_from_current": levelXP - currentXP,
		})
	}
	
	return map[string]interface{}{
		"current_level":    currentLevel,
		"current_title":    elc.GetLevelTitle(currentLevel),
		"current_xp":       currentXP,
		"progress_xp":      progressXP,
		"needed_xp":        neededXP,
		"remaining_xp":     remainingXP,
		"progress_percent": progressPercent,
		"next_levels":      nextLevels,
	}
}
