package gamification

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/pkg/models"
)

// GameEnhancements provides enhanced gaming features
type GameEnhancements struct {
	rand *rand.Rand
}

// NewGameEnhancements creates a new game enhancements manager
func NewGameEnhancements() *GameEnhancements {
	return &GameEnhancements{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// XPMultiplier represents different XP multiplier events
type XPMultiplier struct {
	Name        string
	Multiplier  float64
	Duration    time.Duration
	Icon        string
	Description string
}

// GetActiveMultipliers returns currently active XP multipliers
func (ge *GameEnhancements) GetActiveMultipliers(userProgress *models.UserProgress) []XPMultiplier {
	var multipliers []XPMultiplier

	// Time-based multipliers
	hour := time.Now().Hour()
	weekday := time.Now().Weekday()

	// Early bird bonus (6-9 AM)
	if hour >= 6 && hour <= 9 {
		multipliers = append(multipliers, XPMultiplier{
			Name:        "Early Bird",
			Multiplier:  1.5,
			Duration:    3 * time.Hour,
			Icon:        "🌅",
			Description: "Morning productivity boost!",
		})
	}

	// Night owl bonus (9 PM - 1 AM)
	if hour >= 21 || hour <= 1 {
		multipliers = append(multipliers, XPMultiplier{
			Name:        "Night Owl",
			Multiplier:  1.3,
			Duration:    4 * time.Hour,
			Icon:        "🦉",
			Description: "Late night coding session!",
		})
	}

	// Weekend warrior (Saturday/Sunday)
	if weekday == time.Saturday || weekday == time.Sunday {
		multipliers = append(multipliers, XPMultiplier{
			Name:        "Weekend Warrior",
			Multiplier:  1.2,
			Duration:    48 * time.Hour,
			Icon:        "⚔️",
			Description: "Weekend dedication!",
		})
	}

	// Streak-based multipliers
	if userProgress.CurrentStreak >= 7 {
		streakMultiplier := 1.0 + (float64(userProgress.CurrentStreak) * 0.1)
		if streakMultiplier > 3.0 {
			streakMultiplier = 3.0 // Cap at 3x
		}

		multipliers = append(multipliers, XPMultiplier{
			Name:        fmt.Sprintf("%d-Day Streak", userProgress.CurrentStreak),
			Multiplier:  streakMultiplier,
			Duration:    24 * time.Hour,
			Icon:        "🔥",
			Description: fmt.Sprintf("On fire! %d days strong!", userProgress.CurrentStreak),
		})
	}

	return multipliers
}

// CalculateXPWithMultipliers calculates XP with active multipliers
func (ge *GameEnhancements) CalculateXPWithMultipliers(baseXP int, userProgress *models.UserProgress) (int, []XPMultiplier) {
	multipliers := ge.GetActiveMultipliers(userProgress)

	totalMultiplier := 1.0
	for _, multiplier := range multipliers {
		totalMultiplier *= multiplier.Multiplier
	}

	finalXP := int(float64(baseXP) * totalMultiplier)
	return finalXP, multipliers
}

// CommandCombo represents a command combination chain
type CommandCombo struct {
	Commands    []string
	Multiplier  float64
	Name        string
	Icon        string
	Description string
}

// GetCommandCombos returns available command combinations
func (ge *GameEnhancements) GetCommandCombos() []CommandCombo {
	return []CommandCombo{
		{
			Commands:    []string{"git", "add", "commit", "push"},
			Multiplier:  2.0,
			Name:        "Git Master",
			Icon:        "🚀",
			Description: "Complete git workflow!",
		},
		{
			Commands:    []string{"docker", "build", "run"},
			Multiplier:  1.8,
			Name:        "Container Captain",
			Icon:        "🐳",
			Description: "Docker deployment combo!",
		},
		{
			Commands:    []string{"npm", "install", "npm", "run"},
			Multiplier:  1.6,
			Name:        "Node Ninja",
			Icon:        "📦",
			Description: "NPM workflow mastery!",
		},
		{
			Commands:    []string{"make", "test", "make", "build"},
			Multiplier:  1.7,
			Name:        "Build Master",
			Icon:        "🔨",
			Description: "Build & test combination!",
		},
		{
			Commands:    []string{"vim", "git", "commit"},
			Multiplier:  2.2,
			Name:        "Vim Warrior",
			Icon:        "⚔️",
			Description: "Edit and commit like a pro!",
		},
	}
}

// PowerUp represents temporary power-ups
type PowerUp struct {
	ID          string
	Name        string
	Icon        string
	Description string
	Effect      string
	Duration    time.Duration
	Rarity      float64 // 0.0 to 1.0, lower = rarer
}

// GetAvailablePowerUps returns all available power-ups
func (ge *GameEnhancements) GetAvailablePowerUps() []PowerUp {
	return []PowerUp{
		{
			ID:          "double_xp",
			Name:        "Double XP",
			Icon:        "⚡",
			Description: "Double XP for all commands",
			Effect:      "2x XP multiplier",
			Duration:    30 * time.Minute,
			Rarity:      0.1,
		},
		{
			ID:          "command_frenzy",
			Name:        "Command Frenzy",
			Icon:        "💨",
			Description: "Extra XP for quick commands",
			Effect:      "Bonus XP for rapid commands",
			Duration:    15 * time.Minute,
			Rarity:      0.15,
		},
		{
			ID:          "wisdom_boost",
			Name:        "Wisdom Boost",
			Icon:        "🧠",
			Description: "Learn new commands faster",
			Effect:      "Bonus XP for new commands",
			Duration:    1 * time.Hour,
			Rarity:      0.2,
		},
		{
			ID:          "streak_shield",
			Name:        "Streak Shield",
			Icon:        "🛡️",
			Description: "Protect your streak for 24h",
			Effect:      "Streak protection",
			Duration:    24 * time.Hour,
			Rarity:      0.05,
		},
		{
			ID:          "productivity_surge",
			Name:        "Productivity Surge",
			Icon:        "🌟",
			Description: "Increased XP for all activities",
			Effect:      "1.5x XP for 1 hour",
			Duration:    1 * time.Hour,
			Rarity:      0.12,
		},
	}
}

// TriggerRandomPowerUp has a chance to activate a random power-up
func (ge *GameEnhancements) TriggerRandomPowerUp() *PowerUp {
	powerUps := ge.GetAvailablePowerUps()

	for _, powerUp := range powerUps {
		if ge.rand.Float64() < powerUp.Rarity {
			return &powerUp
		}
	}

	return nil
}

// DailyQuest represents daily challenges
type DailyQuest struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Target      int
	Progress    int
	Reward      int
	XPReward    int
	IsCompleted bool
}

// GenerateDailyQuests creates random daily quests
func (ge *GameEnhancements) GenerateDailyQuests(userProgress *models.UserProgress) []DailyQuest {
	today := time.Now().Format("2006-01-02")
	ge.rand.Seed(int64(userProgress.ID) + int64(time.Now().Unix()/86400)) // Deterministic per day

	quests := []DailyQuest{
		{
			ID:          fmt.Sprintf("commands_%s", today),
			Name:        "Command Explorer",
			Description: fmt.Sprintf("Execute %d commands today", 20+ge.rand.Intn(30)),
			Icon:        "🎯",
			Target:      20 + ge.rand.Intn(30),
			Reward:      50,
			XPReward:    100,
		},
		{
			ID:          fmt.Sprintf("unique_%s", today),
			Name:        "Variety Master",
			Description: fmt.Sprintf("Use %d different commands", 5+ge.rand.Intn(10)),
			Icon:        "🌟",
			Target:      5 + ge.rand.Intn(10),
			Reward:      75,
			XPReward:    150,
		},
		{
			ID:          fmt.Sprintf("git_%s", today),
			Name:        "Git Guru",
			Description: "Perform 5 git operations",
			Icon:        "🚀",
			Target:      5,
			Reward:      60,
			XPReward:    120,
		},
	}

	// Add level-specific quests
	if userProgress.CurrentLevel >= 5 {
		quests = append(quests, DailyQuest{
			ID:          fmt.Sprintf("advanced_%s", today),
			Name:        "Advanced User",
			Description: "Use advanced commands (grep, sed, awk)",
			Icon:        "🧠",
			Target:      3,
			Reward:      100,
			XPReward:    200,
		})
	}

	return quests
}

// WeeklyChallenge represents weekly challenges
type WeeklyChallenge struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Target      int
	Progress    int
	Reward      int
	XPReward    int
	IsCompleted bool
	StartDate   time.Time
	EndDate     time.Time
}

// GenerateWeeklyChallenge creates a weekly challenge
func (ge *GameEnhancements) GenerateWeeklyChallenge(userProgress *models.UserProgress) WeeklyChallenge {
	// Get Monday of current week
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	sunday := monday.AddDate(0, 0, 6)

	challenges := []WeeklyChallenge{
		{
			ID:          fmt.Sprintf("streak_week_%s", monday.Format("2006-01-02")),
			Name:        "Streak Keeper",
			Description: "Maintain a 7-day command streak",
			Icon:        "🔥",
			Target:      7,
			Reward:      500,
			XPReward:    1000,
			StartDate:   monday,
			EndDate:     sunday,
		},
		{
			ID:          fmt.Sprintf("explorer_week_%s", monday.Format("2006-01-02")),
			Name:        "Command Explorer",
			Description: "Execute 200 commands this week",
			Icon:        "🗺️",
			Target:      200,
			Reward:      300,
			XPReward:    600,
			StartDate:   monday,
			EndDate:     sunday,
		},
		{
			ID:          fmt.Sprintf("variety_week_%s", monday.Format("2006-01-02")),
			Name:        "Variety Seeker",
			Description: "Use 50 different commands this week",
			Icon:        "🎨",
			Target:      50,
			Reward:      400,
			XPReward:    800,
			StartDate:   monday,
			EndDate:     sunday,
		},
	}

	// Return random challenge based on user level
	index := (userProgress.CurrentLevel + int(monday.Unix())) % len(challenges)
	return challenges[index]
}

// LevelUpReward represents rewards for leveling up
type LevelUpReward struct {
	Level       int
	Title       string
	Icon        string
	Description string
	Rewards     []string
}

// GetLevelUpReward returns rewards for reaching a specific level
func (ge *GameEnhancements) GetLevelUpReward(level int) LevelUpReward {
	rewards := make([]string, 0)

	// Base rewards
	rewards = append(rewards, fmt.Sprintf("%d XP bonus", level*50))

	// Special milestone rewards
	switch {
	case level == 5:
		return LevelUpReward{
			Level:       5,
			Title:       "Terminal Apprentice",
			Icon:        "🎓",
			Description: "You're getting the hang of this!",
			Rewards:     []string{"100 XP bonus", "New daily quests unlocked", "Command history expanded"},
		}
	case level == 10:
		return LevelUpReward{
			Level:       10,
			Title:       "Command Line Warrior",
			Icon:        "⚔️",
			Description: "A formidable terminal user!",
			Rewards:     []string{"200 XP bonus", "Advanced achievements unlocked", "Weekly challenges available"},
		}
	case level == 25:
		return LevelUpReward{
			Level:       25,
			Title:       "Shell Sage",
			Icon:        "🧙‍♂️",
			Description: "Wisdom flows through your commands!",
			Rewards:     []string{"500 XP bonus", "Master-tier challenges", "Custom badge colors"},
		}
	case level == 50:
		return LevelUpReward{
			Level:       50,
			Title:       "Terminal Master",
			Icon:        "👑",
			Description: "You have mastered the art of the command line!",
			Rewards:     []string{"1000 XP bonus", "Legendary achievements", "Elite status badge"},
		}
	case level == 100:
		return LevelUpReward{
			Level:       100,
			Title:       "Command Line Legend",
			Icon:        "🌟",
			Description: "Your terminal skills are legendary!",
			Rewards:     []string{"2000 XP bonus", "Mythical achievements", "Hall of Fame entry"},
		}
	case level%10 == 0: // Every 10 levels
		return LevelUpReward{
			Level:       level,
			Title:       fmt.Sprintf("Level %d Elite", level),
			Icon:        "💎",
			Description: "Another milestone reached!",
			Rewards:     append(rewards, "Special milestone badge", "Increased XP multiplier"),
		}
	default:
		return LevelUpReward{
			Level:       level,
			Title:       fmt.Sprintf("Level %d Achiever", level),
			Icon:        "⭐",
			Description: "Keep up the great work!",
			Rewards:     rewards,
		}
	}
}

// CommandRarity determines the rarity of commands
type CommandRarity string

const (
	RarityCommon    CommandRarity = "common"
	RarityUncommon  CommandRarity = "uncommon"
	RarityRare      CommandRarity = "rare"
	RarityEpic      CommandRarity = "epic"
	RarityLegendary CommandRarity = "legendary"
)

// GetCommandRarity determines rarity based on usage frequency
func (ge *GameEnhancements) GetCommandRarity(command string, totalUsage int) CommandRarity {
	commonCommands := []string{"ls", "cd", "pwd", "git", "cat", "echo", "mkdir", "rm", "cp", "mv"}
	rareCommands := []string{"awk", "sed", "grep", "find", "xargs", "curl", "wget"}
	epicCommands := []string{"ffmpeg", "pandoc", "docker", "kubectl", "terraform"}
	legendaryCommands := []string{"emacs", "vim", "nano", "gdb", "strace", "tcpdump"}

	// Check predefined rarities first
	for _, cmd := range legendaryCommands {
		if strings.Contains(strings.ToLower(command), cmd) {
			return RarityLegendary
		}
	}

	for _, cmd := range epicCommands {
		if strings.Contains(strings.ToLower(command), cmd) {
			return RarityEpic
		}
	}

	for _, cmd := range rareCommands {
		if strings.Contains(strings.ToLower(command), cmd) {
			return RarityRare
		}
	}

	for _, cmd := range commonCommands {
		if strings.Contains(strings.ToLower(command), cmd) {
			return RarityCommon
		}
	}

	// Determine by usage frequency
	switch {
	case totalUsage > 100:
		return RarityCommon
	case totalUsage > 50:
		return RarityUncommon
	case totalUsage > 10:
		return RarityRare
	case totalUsage > 1:
		return RarityEpic
	default:
		return RarityLegendary
	}
}

// GetRarityInfo returns information about command rarity
func (ge *GameEnhancements) GetRarityInfo(rarity CommandRarity) (string, string, float64) {
	switch rarity {
	case RarityCommon:
		return "⚪", "Common", 1.0
	case RarityUncommon:
		return "🟢", "Uncommon", 1.2
	case RarityRare:
		return "🔵", "Rare", 1.5
	case RarityEpic:
		return "🟣", "Epic", 2.0
	case RarityLegendary:
		return "🟡", "Legendary", 3.0
	default:
		return "⚪", "Common", 1.0
	}
}

// FormatEnhancedXPGain formats XP gain with enhancements
func (ge *GameEnhancements) FormatEnhancedXPGain(baseXP, finalXP int, multipliers []XPMultiplier, rarity CommandRarity) string {
	var parts []string

	// Base XP
	parts = append(parts, fmt.Sprintf("+%d XP", baseXP))

	// Rarity bonus
	rarityIcon, rarityName, rarityMultiplier := ge.GetRarityInfo(rarity)
	if rarityMultiplier > 1.0 {
		parts = append(parts, fmt.Sprintf("%s %s (%.1fx)", rarityIcon, rarityName, rarityMultiplier))
	}

	// Active multipliers
	for _, multiplier := range multipliers {
		parts = append(parts, fmt.Sprintf("%s %s (%.1fx)", multiplier.Icon, multiplier.Name, multiplier.Multiplier))
	}

	// Total
	if finalXP != baseXP {
		parts = append(parts, fmt.Sprintf("= %d XP total", finalXP))
	}

	return strings.Join(parts, " ")
}

// ProgressBar creates an ASCII progress bar
func (ge *GameEnhancements) ProgressBar(current, target int, width int) string {
	if target == 0 {
		return strings.Repeat("█", width)
	}

	progress := float64(current) / float64(target)
	if progress > 1.0 {
		progress = 1.0
	}

	filled := int(progress * float64(width))
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentage := int(progress * 100)

	return fmt.Sprintf("[%s] %d%% (%d/%d)", bar, percentage, current, target)
}
