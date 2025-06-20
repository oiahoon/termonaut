package github

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// ProfileGenerator handles generation of shareable profiles and summaries
type ProfileGenerator struct {
	stats  *stats.StatsCalculator
	badges *BadgeGenerator
}

// NewProfileGenerator creates a new profile generator
func NewProfileGenerator(statsCalculator *stats.StatsCalculator, badgeGenerator *BadgeGenerator) *ProfileGenerator {
	return &ProfileGenerator{
		stats:  statsCalculator,
		badges: badgeGenerator,
	}
}

// ProfileData represents a user's complete profile data
type ProfileData struct {
	UserProgress    *models.UserProgress `json:"user_progress"`
	BasicStats      *stats.BasicStats    `json:"basic_stats"`
	Achievements    []AchievementInfo    `json:"achievements"`
	BadgeURLs       map[string]string    `json:"badge_urls"`
	ProfileMarkdown string               `json:"profile_markdown"`
	LastUpdated     time.Time            `json:"last_updated"`
	AvatarURL       string               `json:"avatar_url,omitempty"`
	AvatarASCII     string               `json:"avatar_ascii,omitempty"`
}

// AchievementInfo represents achievement information for sharing
type AchievementInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	EarnedAt    time.Time `json:"earned_at"`
	IsCompleted bool      `json:"is_completed"`
}

// SocialSnippet represents a shareable social media snippet
type SocialSnippet struct {
	Platform string `json:"platform"`
	Content  string `json:"content"`
	Tags     string `json:"tags"`
}

// GenerateProfile generates a complete shareable profile
func (pg *ProfileGenerator) GenerateProfile(userProgress *models.UserProgress) (*ProfileData, error) {
	return pg.GenerateProfileWithAvatar(userProgress, "", "")
}

// GenerateProfileWithAvatar generates a complete shareable profile with avatar
func (pg *ProfileGenerator) GenerateProfileWithAvatar(userProgress *models.UserProgress, avatarURL, avatarASCII string) (*ProfileData, error) {
	basicStats, err := pg.stats.GetBasicStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get basic stats: %w", err)
	}

	// Generate badge URLs
	badgeURLs := map[string]string{
		"xp":           pg.badges.GenerateXPBadge(userProgress),
		"commands":     pg.badges.GenerateCommandsBadge(basicStats.TotalCommands),
		"streak":       pg.badges.GenerateStreakBadge(userProgress.CurrentStreak),
		"productivity": pg.badges.GenerateProductivityBadge(0.8),   // Placeholder calculation
		"achievements": pg.badges.GenerateAchievementsBadge(5, 10), // Placeholder
	}

	if userProgress.LastActivityDate != nil {
		badgeURLs["last_active"] = pg.badges.GenerateLastActiveBadge(*userProgress.LastActivityDate)
	}

	// Generate achievements info (placeholder for now)
	achievements := []AchievementInfo{
		{
			ID:          "first_command",
			Name:        "First Steps",
			Description: "Execute your first command",
			IsCompleted: basicStats.TotalCommands > 0,
		},
		{
			ID:          "command_master",
			Name:        "Command Master",
			Description: "Execute 1000 commands",
			IsCompleted: basicStats.TotalCommands >= 1000,
		},
	}

	// Generate profile markdown
	profileMarkdown := pg.generateProfileMarkdownWithAvatar(userProgress, basicStats, badgeURLs, achievements, avatarURL, avatarASCII)

	return &ProfileData{
		UserProgress:    userProgress,
		BasicStats:      basicStats,
		Achievements:    achievements,
		BadgeURLs:       badgeURLs,
		ProfileMarkdown: profileMarkdown,
		LastUpdated:     time.Now(),
		AvatarURL:       avatarURL,
		AvatarASCII:     avatarASCII,
	}, nil
}

// generateProfileMarkdown creates a markdown representation of the profile
func (pg *ProfileGenerator) generateProfileMarkdown(userProgress *models.UserProgress, basicStats *stats.BasicStats, badgeURLs map[string]string, achievements []AchievementInfo) string {
	return pg.generateProfileMarkdownWithAvatar(userProgress, basicStats, badgeURLs, achievements, "", "")
}

// generateProfileMarkdownWithAvatar creates a markdown representation of the profile with avatar
func (pg *ProfileGenerator) generateProfileMarkdownWithAvatar(userProgress *models.UserProgress, basicStats *stats.BasicStats, badgeURLs map[string]string, achievements []AchievementInfo, avatarURL, avatarASCII string) string {
	var builder strings.Builder

	// Header
	builder.WriteString("# 🚀 My Termonaut Profile\n\n")
	builder.WriteString("*Gamified terminal productivity tracking*\n\n")

	// Badges section at the top
	builder.WriteString("## 📊 Badges\n\n")
	for label, url := range badgeURLs {
		builder.WriteString(fmt.Sprintf("![%s](%s) ", strings.Title(label), url))
	}
	builder.WriteString("\n\n")

	// Avatar and Stats Layout
	if avatarURL != "" {
		builder.WriteString("## 🎨 Profile & Stats\n\n")

		// Use table layout for left-right layout
		builder.WriteString("<table><tr>\n")

		// Avatar column (only SVG, no ASCII for GitHub)
		builder.WriteString("<td width=\"40%\" align=\"center\">\n\n")
		builder.WriteString("### 👤 Avatar\n\n")
		builder.WriteString(fmt.Sprintf("![Avatar](%s)\n\n", avatarURL))
		builder.WriteString("</td>\n")

		// Stats column
		builder.WriteString("<td width=\"60%\">\n\n")
		builder.WriteString("### 📊 Stats Overview\n\n")
		builder.WriteString(fmt.Sprintf("**Level**: %d (XP: %d)  \n", userProgress.CurrentLevel, userProgress.TotalXP))
		builder.WriteString(fmt.Sprintf("**Total Commands**: %d  \n", basicStats.TotalCommands))
		builder.WriteString(fmt.Sprintf("**Unique Commands**: %d  \n", basicStats.UniqueCommands))
		builder.WriteString(fmt.Sprintf("**Current Streak**: %d days  \n", userProgress.CurrentStreak))
		builder.WriteString(fmt.Sprintf("**Longest Streak**: %d days  \n", userProgress.LongestStreak))
		builder.WriteString(fmt.Sprintf("**Commands Today**: %d  \n", basicStats.CommandsToday))

		// Most used command
		if basicStats.MostUsedCommand != "" {
			builder.WriteString(fmt.Sprintf("**Favorite Command**: `%s` (%d times)  \n",
				basicStats.MostUsedCommand, basicStats.MostUsedCount))
		}
		builder.WriteString("\n</td>\n")

		builder.WriteString("</tr></table>\n\n")
	} else {
		// Fallback to original layout without avatar
		builder.WriteString("## 📈 Overview\n\n")
		builder.WriteString(fmt.Sprintf("- **Level**: %d (XP: %d)\n", userProgress.CurrentLevel, userProgress.TotalXP))
		builder.WriteString(fmt.Sprintf("- **Total Commands**: %d\n", basicStats.TotalCommands))
		builder.WriteString(fmt.Sprintf("- **Unique Commands**: %d\n", basicStats.UniqueCommands))
		builder.WriteString(fmt.Sprintf("- **Current Streak**: %d days\n", userProgress.CurrentStreak))
		builder.WriteString(fmt.Sprintf("- **Longest Streak**: %d days\n", userProgress.LongestStreak))
		builder.WriteString(fmt.Sprintf("- **Commands Today**: %d\n", basicStats.CommandsToday))

		// Most used command
		if basicStats.MostUsedCommand != "" {
			builder.WriteString(fmt.Sprintf("- **Favorite Command**: `%s` (%d times)\n",
				basicStats.MostUsedCommand, basicStats.MostUsedCount))
		}
		builder.WriteString("\n")
	}

	// Achievements section
	builder.WriteString("## 🏆 Achievements\n\n")
	completedCount := 0
	for _, achievement := range achievements {
		if achievement.IsCompleted {
			completedCount++
			builder.WriteString(fmt.Sprintf("- ✅ **%s**: %s\n", achievement.Name, achievement.Description))
		} else {
			builder.WriteString(fmt.Sprintf("- ⏳ **%s**: %s\n", achievement.Name, achievement.Description))
		}
	}

	builder.WriteString(fmt.Sprintf("\n*%d/%d achievements unlocked*\n\n", completedCount, len(achievements)))

	// Top commands section
	if len(basicStats.TopCommands) > 0 {
		builder.WriteString("## 🔥 Top Commands\n\n")
		for i, cmd := range basicStats.TopCommands {
			if i >= 5 { // Show top 5
				break
			}
			cmdStr := cmd["command"].(string)
			count := cmd["count"].(int)

			// Create a visual bar
			percentage := float64(count) / float64(basicStats.MostUsedCount) * 100
			barLength := int(percentage / 5) // Scale to 20 chars max
			if barLength < 1 {
				barLength = 1
			}
			bar := strings.Repeat("█", barLength)

			builder.WriteString(fmt.Sprintf("%d. `%s` (%d times) %s\n", i+1, cmdStr, count, bar))
		}
		builder.WriteString("\n")
	}

	// Footer
	builder.WriteString("---\n\n")
	builder.WriteString("*Generated by [Termonaut](https://github.com/oiahoon/termonaut) - Terminal productivity tracker*\n")
	builder.WriteString(fmt.Sprintf("*Last updated: %s*\n", time.Now().Format("January 2, 2006")))

	return builder.String()
}

// GenerateSocialSnippets creates social media snippets for sharing
func (pg *ProfileGenerator) GenerateSocialSnippets(profileData *ProfileData) []SocialSnippet {
	snippets := []SocialSnippet{}

	// Twitter snippet
	twitterContent := fmt.Sprintf(
		"🚀 My terminal stats this week:\n\n"+
			"📊 %d commands executed\n"+
			"⚡ Level %d (XP: %d)\n"+
			"🔥 %d day streak\n"+
			"🎯 %d unique commands\n\n"+
			"#TerminalProductivity #CommandLine #Termonaut",
		profileData.BasicStats.TotalCommands,
		profileData.UserProgress.CurrentLevel,
		profileData.UserProgress.TotalXP,
		profileData.UserProgress.CurrentStreak,
		profileData.BasicStats.UniqueCommands,
	)

	snippets = append(snippets, SocialSnippet{
		Platform: "twitter",
		Content:  twitterContent,
		Tags:     "#TerminalProductivity #CommandLine #Termonaut",
	})

	// LinkedIn snippet
	linkedinContent := fmt.Sprintf(
		"Tracking my terminal productivity with Termonaut! 🚀\n\n"+
			"This week's highlights:\n"+
			"• %d commands executed\n"+
			"• Reached Level %d with %d XP\n"+
			"• Maintained a %d-day streak\n"+
			"• Used %d unique commands\n\n"+
			"Gamification makes even terminal work more engaging! "+
			"What tools do you use to track your productivity?\n\n"+
			"#Productivity #TerminalTools #CommandLine #TechTools",
		profileData.BasicStats.TotalCommands,
		profileData.UserProgress.CurrentLevel,
		profileData.UserProgress.TotalXP,
		profileData.UserProgress.CurrentStreak,
		profileData.BasicStats.UniqueCommands,
	)

	snippets = append(snippets, SocialSnippet{
		Platform: "linkedin",
		Content:  linkedinContent,
		Tags:     "#Productivity #TerminalTools #CommandLine #TechTools",
	})

	// Reddit snippet
	redditContent := fmt.Sprintf(
		"My terminal productivity stats using Termonaut\n\n"+
			"Commands executed: %d\n"+
			"Current level: %d (XP: %d)\n"+
			"Streak: %d days\n"+
			"Unique commands: %d\n\n"+
			"Pretty cool to see my terminal usage gamified. "+
			"Anyone else using productivity tracking tools?",
		profileData.BasicStats.TotalCommands,
		profileData.UserProgress.CurrentLevel,
		profileData.UserProgress.TotalXP,
		profileData.UserProgress.CurrentStreak,
		profileData.BasicStats.UniqueCommands,
	)

	snippets = append(snippets, SocialSnippet{
		Platform: "reddit",
		Content:  redditContent,
		Tags:     "r/commandline r/productivity r/programming",
	})

	return snippets
}

// GenerateStatsSummary creates a brief stats summary for sharing
func (pg *ProfileGenerator) GenerateStatsSummary(userProgress *models.UserProgress, basicStats *stats.BasicStats) string {
	return fmt.Sprintf(
		"🚀 Termonaut Stats: Level %d • %d commands • %d day streak • %d unique commands",
		userProgress.CurrentLevel,
		basicStats.TotalCommands,
		userProgress.CurrentStreak,
		basicStats.UniqueCommands,
	)
}

// ExportProfile exports profile data in various formats
func (pg *ProfileGenerator) ExportProfile(profileData *ProfileData, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(profileData)
	case "markdown":
		return []byte(profileData.ProfileMarkdown), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
