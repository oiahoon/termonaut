package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// Simple TUI styles
var (
	simpleHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#7C3AED")).
				Padding(0, 2).
				Bold(true)

	simpleCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Margin(1, 0)

	simpleHighlightStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B9D")).
				Bold(true)

	simpleHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
)

// SimpleDashboard represents a simple dashboard model
type SimpleDashboard struct {
	db       *database.DB
	commands []*models.Command
	sessions []*models.Session
	width    int
	height   int
	view     int // 0=stats, 1=analytics, 2=heatmap, 3=achievements
	ready    bool
	err      error
}

// NewSimpleDashboard creates a new simple dashboard
func NewSimpleDashboard(db *database.DB) (*SimpleDashboard, error) {
	commands, err := db.GetAllCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to get commands: %w", err)
	}

	sessions, err := db.GetAllSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	return &SimpleDashboard{
		db:       db,
		commands: commands,
		sessions: sessions,
		view:     0,
		ready:    true,
	}, nil
}

// Init implements tea.Model
func (m SimpleDashboard) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m SimpleDashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.view = 0
		case "2":
			m.view = 1
		case "3":
			m.view = 2
		case "4":
			m.view = 3
		case "left", "h":
			if m.view > 0 {
				m.view--
			}
		case "right", "l":
			if m.view < 3 {
				m.view++
			}
		case "r":
			// Refresh data
			return m.refreshData()
		}
	}

	return m, nil
}

// View implements tea.Model
func (m SimpleDashboard) View() string {
	if !m.ready {
		return "Loading..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Header
	header := simpleHeaderStyle.Width(m.width).Render("🎮 Termonaut Interactive Dashboard")

	// Navigation
	nav := m.renderNavigation()

	// Content
	var content string
	switch m.view {
	case 0:
		content = m.renderStats()
	case 1:
		content = m.renderAnalytics()
	case 2:
		content = m.renderHeatmap()
	case 3:
		content = m.renderAchievements()
	}

	// Help
	help := simpleHelpStyle.Render("1-4: switch views • ←/→ or h/l: navigate • r: refresh • q: quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		nav,
		content,
		help,
	)
}

func (m SimpleDashboard) renderNavigation() string {
	views := []string{"📊 Stats", "📈 Analytics", "🔥 Heatmap", "🏆 Achievements"}
	var nav []string

	for i, view := range views {
		style := lipgloss.NewStyle().Padding(0, 1)
		if i == m.view {
			style = style.Background(lipgloss.Color("#7C3AED")).Foreground(lipgloss.Color("#FFFFFF"))
		} else {
			style = style.Background(lipgloss.Color("#374151")).Foreground(lipgloss.Color("#9CA3AF"))
		}
		nav = append(nav, style.Render(view))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, nav...) + "\n"
}

func (m SimpleDashboard) renderStats() string {
	totalCommands := len(m.commands)
	totalSessions := len(m.sessions)

	uniqueCommands := make(map[string]bool)
	for _, cmd := range m.commands {
		uniqueCommands[cmd.Command] = true
	}

	// Get gamification stats
	statsCalc := stats.New(m.db)
	gamificationStats, err := statsCalc.GetGamificationStats()
	if err != nil {
		gamificationStats = &stats.GamificationStats{}
	}

	content := fmt.Sprintf(`📊 Overview Statistics

%s: %s
%s: %s
%s: %s

🎮 Gamification Status:`,
		simpleHighlightStyle.Render("Total Commands"),
		fmt.Sprintf("%d", totalCommands),
		simpleHighlightStyle.Render("Unique Commands"),
		fmt.Sprintf("%d", len(uniqueCommands)),
		simpleHighlightStyle.Render("Total Sessions"),
		fmt.Sprintf("%d", totalSessions),
	)

	if gamificationStats.ProgressInfo != nil {
		content += fmt.Sprintf(`
• Level: %s %d (%s)
• XP: %s / %s (%s total)
• Progress: %s`,
			simpleHighlightStyle.Render("Level"),
			gamificationStats.ProgressInfo.CurrentLevel,
			gamificationStats.ProgressInfo.LevelTitle,
			simpleHighlightStyle.Render(fmt.Sprintf("%d", gamificationStats.ProgressInfo.XPInCurrentLevel)),
			fmt.Sprintf("%d", gamificationStats.ProgressInfo.XPForCurrentLevel),
			fmt.Sprintf("%d", gamificationStats.ProgressInfo.TotalXP),
			simpleHighlightStyle.Render(fmt.Sprintf("%.1f%%", gamificationStats.ProgressInfo.ProgressPercent)),
		)
	} else {
		content += "\n• No gamification data available"
	}

	if len(m.commands) > 0 {
		content += "\n\n🕒 Recent Commands:"
		for i, cmd := range m.commands {
			if i >= 5 { // Show only last 5
				break
			}
			content += fmt.Sprintf("\n• %s %s",
				simpleHighlightStyle.Render(cmd.Command),
				cmd.Timestamp.Format("15:04:05"),
			)
		}
	}

	return simpleCardStyle.Render(content)
}

func (m SimpleDashboard) renderAnalytics() string {
	if len(m.commands) == 0 {
		return simpleCardStyle.Render("No analytics data available. Start using your terminal!")
	}

	analyzer := analytics.NewProductivityAnalyzer()
	metrics := analyzer.AnalyzeProductivity(m.commands, m.sessions)

	content := fmt.Sprintf(`📈 Productivity Analytics

%s: %s/100
%s: %s%%
%s: %s%%

📊 Time Distribution:
• Morning (6-12): %s%%
• Afternoon (12-18): %s%%
• Evening (18-24): %s%%
• Night (0-6): %s%%

💪 Performance Insights:
• Most productive day: %s
• Peak hours: %s
• Current streak: %s days
• Specialization: %s`,
		simpleHighlightStyle.Render("Overall Score"),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.OverallScore)),
		simpleHighlightStyle.Render("Efficiency"),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.EfficiencyMetrics.UniqueCommandRatio*100)),
		simpleHighlightStyle.Render("Consistency"),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.StreakAnalysis.ConsistencyScore*100)),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.DailyPattern.MorningScore)),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.DailyPattern.AfternoonScore)),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.DailyPattern.EveningScore)),
		simpleHighlightStyle.Render(fmt.Sprintf("%.1f", metrics.DailyPattern.NightScore)),
		simpleHighlightStyle.Render(metrics.WeeklyPattern.MostProductiveDay.String()),
		simpleHighlightStyle.Render(fmt.Sprintf("%v", metrics.DailyPattern.PeakHours)),
		simpleHighlightStyle.Render(fmt.Sprintf("%d", metrics.StreakAnalysis.CurrentStreak)),
		simpleHighlightStyle.Render(metrics.CategoryInsights.SpecializationLevel),
	)

	return simpleCardStyle.Render(content)
}

func (m SimpleDashboard) renderHeatmap() string {
	if len(m.commands) == 0 {
		return simpleCardStyle.Render("No heatmap data available. Start using your terminal!")
	}

	// Create hour-based activity map
	hourActivity := make(map[int]int)
	for _, cmd := range m.commands {
		hour := cmd.Timestamp.Hour()
		hourActivity[hour]++
	}

	content := "🔥 Activity Heatmap\n\n"

	// Find peak hour
	maxActivity := 0
	peakHour := 0
	for hour, count := range hourActivity {
		if count > maxActivity {
			maxActivity = count
			peakHour = hour
		}
	}

	if maxActivity > 0 {
		content += fmt.Sprintf("Peak Activity: %s at %02d:00\n",
			simpleHighlightStyle.Render(fmt.Sprintf("%d commands", maxActivity)),
			peakHour,
		)
		content += fmt.Sprintf("Total Active Hours: %s\n\n",
			simpleHighlightStyle.Render(fmt.Sprintf("%d", len(hourActivity))),
		)
	}

	// Add simplified heatmap
	content += "Activity by Hour (today):\n"
	hours := []string{"06", "08", "10", "12", "14", "16", "18", "20", "22"}
	for _, hour := range hours {
		hourInt := 0
		fmt.Sscanf(hour, "%d", &hourInt)
		intensity := hourActivity[hourInt]

		bar := ""
		if intensity > 0 {
			level := intensity
			if level > 10 {
				level = 10
			}
			bar = strings.Repeat("█", level) + strings.Repeat("░", 10-level)
		} else {
			bar = strings.Repeat("░", 10)
		}

		content += fmt.Sprintf("%s:00 [%s] %d\n", hour, bar, intensity)
	}

	return simpleCardStyle.Render(content)
}

func (m SimpleDashboard) renderAchievements() string {
	statsCalc := stats.New(m.db)
	gamificationStats, err := statsCalc.GetGamificationStats()
	if err != nil || gamificationStats == nil {
		return simpleCardStyle.Render("No achievement data available.")
	}

	content := fmt.Sprintf("🏆 Achievements (%d unlocked)\n\n", len(gamificationStats.Achievements))

	// Show unlocked achievements
	if len(gamificationStats.Achievements) > 0 {
		content += "✅ Unlocked:\n"
		for _, userAchievement := range gamificationStats.Achievements {
			content += fmt.Sprintf("• %s %s %s (+%d XP)\n",
				userAchievement.Achievement.Icon,
				simpleHighlightStyle.Render(userAchievement.Achievement.Name),
				userAchievement.Achievement.Description,
				userAchievement.Achievement.XPReward,
			)
		}
		content += "\n"
	}

	// Show next achievements
	if len(gamificationStats.NextAchievements) > 0 {
		content += "🎯 In Progress:\n"
		for _, nextAchievement := range gamificationStats.NextAchievements {
			progress := int(nextAchievement.Percentage / 10) // Scale to 10
			bar := strings.Repeat("█", progress) + strings.Repeat("░", 10-progress)
			content += fmt.Sprintf("• %s %s\n  [%s] %s%%\n",
				nextAchievement.Achievement.Icon,
				simpleHighlightStyle.Render(nextAchievement.Achievement.Name),
				bar,
				simpleHighlightStyle.Render(fmt.Sprintf("%.0f", nextAchievement.Percentage)),
			)
		}
	}

	return simpleCardStyle.Render(content)
}

func (m SimpleDashboard) refreshData() (tea.Model, tea.Cmd) {
	// Reload data from database
	commands, err := m.db.GetAllCommands()
	if err != nil {
		m.err = err
		return m, nil
	}

	sessions, err := m.db.GetAllSessions()
	if err != nil {
		m.err = err
		return m, nil
	}

	m.commands = commands
	m.sessions = sessions
	m.err = nil

	return m, nil
}

// RunSimpleDashboard starts the simple interactive dashboard
func RunSimpleDashboard(db *database.DB) error {
	model, err := NewSimpleDashboard(db)
	if err != nil {
		return fmt.Errorf("failed to create dashboard model: %w", err)
	}

	program := tea.NewProgram(
		*model,
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run dashboard: %w", err)
	}

	return nil
}
