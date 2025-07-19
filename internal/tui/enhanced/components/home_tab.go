package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/oiahoon/termonaut/internal/avatar"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// HomeTabComponent handles the home tab rendering and logic
type HomeTabComponent struct {
	db        *database.DB
	statsCalc *stats.StatsCalculator
	avatarMgr *avatar.AvatarManager
	theme     *Theme
}

// NewHomeTabComponent creates a new home tab component
func NewHomeTabComponent(db *database.DB, statsCalc *stats.StatsCalculator, avatarMgr *avatar.AvatarManager, theme *Theme) *HomeTabComponent {
	return &HomeTabComponent{
		db:        db,
		statsCalc: statsCalc,
		avatarMgr: avatarMgr,
		theme:     theme,
	}
}

// Render renders the home tab content
func (h *HomeTabComponent) Render(width, height int) string {
	// Calculate available space
	contentHeight := height - 4 // Account for header and footer
	contentWidth := width - 4   // Account for padding

	// Get user progress and stats
	progress, err := h.db.GetUserProgress()
	if err != nil {
		return h.renderError("Failed to load user progress", width, height)
	}

	basicStats, err := h.db.GetBasicStats()
	if err != nil {
		return h.renderError("Failed to load basic stats", width, height)
	}

	// Layout sections
	sections := h.layoutSections(contentWidth, contentHeight, progress, basicStats)
	
	// Combine sections
	var content strings.Builder
	for _, section := range sections {
		content.WriteString(section)
		content.WriteString("\n")
	}

	return h.theme.ContentBox.Render(content.String())
}

// layoutSections determines the layout of home tab sections
func (h *HomeTabComponent) layoutSections(width, height int, progress *models.UserProgress, basicStats map[string]interface{}) []string {
	var sections []string

	// Avatar and user info section
	if height > 20 && width > 80 {
		sections = append(sections, h.renderUserSection(width/2, progress))
	}

	// Quick stats section
	sections = append(sections, h.renderQuickStats(width, basicStats))

	// Recent activity section
	if height > 15 {
		sections = append(sections, h.renderRecentActivity(width))
	}

	// Level progress section
	sections = append(sections, h.renderLevelProgress(width, progress))

	return sections
}

// renderUserSection renders the user avatar and basic info
func (h *HomeTabComponent) renderUserSection(width int, progress *models.UserProgress) string {
	var content strings.Builder

	// Get avatar
	avatarStr := ""
	if h.avatarMgr != nil {
		avatarStr = h.avatarMgr.GetAvatar(width/3, 8)
	}

	// User info
	userInfo := fmt.Sprintf(
		"Level %d %s\n%d XP • %d Commands\nStreak: %d days",
		progress.Level,
		h.getLevelTitle(progress.Level),
		progress.TotalXP,
		progress.TotalCommands,
		progress.CurrentStreak,
	)

	// Layout side by side if width allows
	if width > 60 && avatarStr != "" {
		avatarLines := strings.Split(avatarStr, "\n")
		infoLines := strings.Split(userInfo, "\n")
		
		maxLines := len(avatarLines)
		if len(infoLines) > maxLines {
			maxLines = len(infoLines)
		}

		for i := 0; i < maxLines; i++ {
			var line strings.Builder
			
			// Avatar column
			if i < len(avatarLines) {
				line.WriteString(avatarLines[i])
			} else {
				line.WriteString(strings.Repeat(" ", len(avatarLines[0])))
			}
			
			line.WriteString("  ") // Spacing
			
			// Info column
			if i < len(infoLines) {
				line.WriteString(h.theme.UserInfo.Render(infoLines[i]))
			}
			
			content.WriteString(line.String())
			if i < maxLines-1 {
				content.WriteString("\n")
			}
		}
	} else {
		// Vertical layout for narrow screens
		if avatarStr != "" {
			content.WriteString(avatarStr)
			content.WriteString("\n\n")
		}
		content.WriteString(h.theme.UserInfo.Render(userInfo))
	}

	return h.theme.SectionBox.Render(content.String())
}

// renderQuickStats renders quick statistics
func (h *HomeTabComponent) renderQuickStats(width int, basicStats map[string]interface{}) string {
	stats := []struct {
		label string
		value interface{}
		icon  string
	}{
		{"Commands Today", basicStats["commands_today"], "🎯"},
		{"Active Time", basicStats["active_time_today"], "⏱️"},
		{"Sessions", basicStats["sessions_today"], "📱"},
		{"New Commands", basicStats["new_commands_today"], "⭐"},
	}

	// Calculate columns based on width
	cols := 2
	if width > 100 {
		cols = 4
	} else if width > 60 {
		cols = 3
	}

	var content strings.Builder
	for i, stat := range stats {
		if i > 0 && i%cols == 0 {
			content.WriteString("\n")
		}

		statStr := fmt.Sprintf("%s %s\n%v", stat.icon, stat.label, stat.value)
		statBox := h.theme.StatBox.Render(statStr)
		
		content.WriteString(statBox)
		if (i+1)%cols != 0 && i < len(stats)-1 {
			content.WriteString(" ")
		}
	}

	return h.theme.SectionBox.Render(content.String())
}

// renderRecentActivity renders recent command activity
func (h *HomeTabComponent) renderRecentActivity(width int) string {
	recentCommands, err := h.db.GetRecentCommands(5)
	if err != nil {
		return h.theme.ErrorBox.Render("Failed to load recent activity")
	}

	var content strings.Builder
	content.WriteString(h.theme.SectionTitle.Render("🔥 Recent Activity"))
	content.WriteString("\n\n")

	if len(recentCommands) == 0 {
		content.WriteString(h.theme.EmptyState.Render("No recent activity"))
	} else {
		for i, cmd := range recentCommands {
			// Truncate long commands
			command := cmd.Command
			maxLen := width - 20
			if len(command) > maxLen {
				command = command[:maxLen-3] + "..."
			}

			timeStr := cmd.Timestamp.Format("15:04")
			exitIcon := "✅"
			if cmd.ExitCode != 0 {
				exitIcon = "❌"
			}

			activityLine := fmt.Sprintf("%s %s %s", exitIcon, timeStr, command)
			content.WriteString(h.theme.ActivityItem.Render(activityLine))
			
			if i < len(recentCommands)-1 {
				content.WriteString("\n")
			}
		}
	}

	return h.theme.SectionBox.Render(content.String())
}

// renderLevelProgress renders level progression
func (h *HomeTabComponent) renderLevelProgress(width int, progress *models.UserProgress) string {
	var content strings.Builder
	content.WriteString(h.theme.SectionTitle.Render("🚀 Level Progress"))
	content.WriteString("\n\n")

	// Calculate progress to next level
	currentLevelXP := h.getXPForLevel(progress.Level)
	nextLevelXP := h.getXPForLevel(progress.Level + 1)
	progressXP := progress.TotalXP - currentLevelXP
	neededXP := nextLevelXP - currentLevelXP

	progressPercent := float64(progressXP) / float64(neededXP) * 100
	if progressPercent > 100 {
		progressPercent = 100
	}

	// Progress bar
	barWidth := width - 20
	if barWidth > 50 {
		barWidth = 50
	}
	if barWidth < 10 {
		barWidth = 10
	}

	filledWidth := int(float64(barWidth) * progressPercent / 100)
	progressBar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	content.WriteString(fmt.Sprintf("Level %d → %d\n", progress.Level, progress.Level+1))
	content.WriteString(h.theme.ProgressBar.Render(progressBar))
	content.WriteString(fmt.Sprintf("\n%d / %d XP (%.1f%%)", progressXP, neededXP, progressPercent))

	return h.theme.SectionBox.Render(content.String())
}

// renderError renders an error message
func (h *HomeTabComponent) renderError(message string, width, height int) string {
	return h.theme.ErrorBox.Render(fmt.Sprintf("❌ %s", message))
}

// getLevelTitle returns the title for a given level
func (h *HomeTabComponent) getLevelTitle(level int) string {
	titles := []string{
		"Rookie", "Explorer", "Navigator", "Commander", "Captain",
		"Major", "Colonel", "General", "Admiral", "Legend",
	}
	
	if level < len(titles) {
		return titles[level]
	}
	return "Master"
}

// getXPForLevel calculates XP required for a given level
func (h *HomeTabComponent) getXPForLevel(level int) int {
	if level <= 0 {
		return 0
	}
	// Simple exponential progression
	return level * level * 100
}
