package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1).
			MarginRight(2)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B7"))

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	focusedButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(lipgloss.Color("#FF75B7")).
				Padding(0, 3)

	blurredButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#626262")).
				Background(lipgloss.Color("#2A2A2A")).
				Padding(0, 3)
)

// TabID represents different tabs in the dashboard
type TabID int

const (
	StatsTab TabID = iota
	AnalyticsTab
	HeatmapTab
	AchievementsTab
	SettingsTab
)

// Tab represents a dashboard tab
type Tab struct {
	ID    TabID
	Title string
	Icon  string
}

// DashboardModel represents the main dashboard model
type DashboardModel struct {
	// Core data
	db       *database.DB
	commands []*models.Command
	sessions []*models.Session

	// UI state
	currentTab TabID
	tabs       []Tab
	width      int
	height     int
	ready      bool
	err        error

	// Components
	list              list.Model
	table             table.Model
	viewport          viewport.Model
	progress          progress.Model
	gamificationStats *stats.GamificationStats

	// Navigation
	keys keyMap
}

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "previous tab"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "next tab"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch tab"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(db *database.DB) (*DashboardModel, error) {
	// Get data
	commands, err := db.GetAllCommands()
	if err != nil {
		return nil, fmt.Errorf("failed to get commands: %w", err)
	}

	sessions, err := db.GetAllSessions()
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	// Get gamification stats
	statsCalc := stats.New(db)
	gamificationStats, err := statsCalc.GetGamificationStats()
	if err != nil {
		// Non-fatal error, continue with empty stats
		gamificationStats = &stats.GamificationStats{}
	}

	// Initialize tabs
	tabs := []Tab{
		{ID: StatsTab, Title: "Overview", Icon: "📊"},
		{ID: AnalyticsTab, Title: "Analytics", Icon: "📈"},
		{ID: HeatmapTab, Title: "Heatmap", Icon: "🔥"},
		{ID: AchievementsTab, Title: "Achievements", Icon: "🏆"},
		{ID: SettingsTab, Title: "Settings", Icon: "⚙️"},
	}

	// Initialize progress bar
	prog := progress.New(progress.WithDefaultGradient())

	return &DashboardModel{
		db:                db,
		commands:          commands,
		sessions:          sessions,
		currentTab:        StatsTab,
		tabs:              tabs,
		keys:              newKeyMap(),
		progress:          prog,
		gamificationStats: gamificationStats,
	}, nil
}

// Init implements tea.Model
func (m DashboardModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update implements tea.Model
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.ready = true
			m.initializeComponents()
		}

		// Update component sizes
		m.updateComponentSizes()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Left):
			if m.currentTab > 0 {
				m.currentTab--
			} else {
				m.currentTab = TabID(len(m.tabs) - 1)
			}

		case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Tab):
			if m.currentTab < TabID(len(m.tabs)-1) {
				m.currentTab++
			} else {
				m.currentTab = 0
			}

		case key.Matches(msg, m.keys.Enter):
			// Handle enter based on current tab
			switch m.currentTab {
			case StatsTab:
				// Maybe refresh stats
			case AnalyticsTab:
				// Maybe show detailed analytics
			}
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
	}

	// Update active components based on current tab (only if ready)
	if m.ready {
		switch m.currentTab {
		case StatsTab, AnalyticsTab:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case HeatmapTab:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		case AchievementsTab:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m DashboardModel) View() string {
	if !m.ready {
		return "Loading Termonaut Dashboard..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Header
	header := m.renderHeader()

	// Tab navigation
	tabNav := m.renderTabNavigation()

	// Content based on current tab
	content := m.renderCurrentTabContent()

	// Help
	help := m.renderHelp()

	// Combine all sections
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabNav,
		content,
		help,
	)
}

func (m *DashboardModel) initializeComponents() {
	// Initialize list for stats/analytics
	items := []list.Item{}
	if len(m.commands) > 0 {
		// Add recent commands as list items
		for i, cmd := range m.commands {
			if i >= 10 { // Limit to 10 recent commands
				break
			}
			items = append(items, ListItem{
				title: cmd.Command,
				desc:  fmt.Sprintf("Executed at %s", cmd.Timestamp.Format("15:04:05")),
			})
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#FF75B7")).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#FF75B7")).
		Padding(0, 0, 0, 1)

	m.list = list.New(items, delegate, 0, 0)
	m.list.Title = "Recent Commands"
	m.list.Styles.Title = titleStyle
	m.list.SetShowStatusBar(false)

	// Initialize table for achievements
	columns := []table.Column{
		{Title: "Icon", Width: 4},
		{Title: "Achievement", Width: 20},
		{Title: "Progress", Width: 15},
		{Title: "XP", Width: 8},
	}

	rows := []table.Row{}
	if m.gamificationStats != nil && len(m.gamificationStats.NextAchievements) > 0 {
		for _, achievement := range m.gamificationStats.NextAchievements {
			progressBar := m.createProgressBar(achievement.Percentage/100.0, 10)
			rows = append(rows, table.Row{
				achievement.Achievement.Icon,
				achievement.Achievement.Name,
				progressBar + fmt.Sprintf(" %.0f%%", achievement.Percentage),
				fmt.Sprintf("+%d", achievement.Achievement.XPReward),
			})
		}
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		BorderBottom(true).
		Bold(false)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#874BFD")).
		Bold(false)
	m.table.SetStyles(tableStyle)

	// Initialize viewport for heatmap
	m.viewport = viewport.New(m.width-4, m.height-10)
	m.viewport.Style = cardStyle
}

func (m *DashboardModel) updateComponentSizes() {
	contentWidth := m.width - 4
	contentHeight := m.height - 8

	m.list.SetSize(contentWidth/2-2, contentHeight)
	m.viewport.Width = contentWidth
	m.viewport.Height = contentHeight

	m.table.SetWidth(contentWidth)
	m.table.SetHeight(contentHeight)
}

func (m DashboardModel) renderHeader() string {
	var title string
	if m.gamificationStats != nil && m.gamificationStats.ProgressInfo != nil {
		title = fmt.Sprintf("🎮 Termonaut Dashboard - Level %d %s (XP: %d)",
			m.gamificationStats.ProgressInfo.CurrentLevel,
			m.gamificationStats.ProgressInfo.LevelTitle,
			m.gamificationStats.ProgressInfo.TotalXP)
	} else {
		title = "🎮 Termonaut Dashboard"
	}

	return headerStyle.Width(m.width).Render(title)
}

func (m DashboardModel) renderTabNavigation() string {
	var tabs []string

	for i, tab := range m.tabs {
		style := blurredButtonStyle
		if TabID(i) == m.currentTab {
			style = focusedButtonStyle
		}
		tabs = append(tabs, style.Render(fmt.Sprintf("%s %s", tab.Icon, tab.Title)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...) + "\n"
}

func (m DashboardModel) renderCurrentTabContent() string {
	switch m.currentTab {
	case StatsTab:
		return m.renderStatsTab()
	case AnalyticsTab:
		return m.renderAnalyticsTab()
	case HeatmapTab:
		return m.renderHeatmapTab()
	case AchievementsTab:
		return m.renderAchievementsTab()
	case SettingsTab:
		return m.renderSettingsTab()
	default:
		return "Tab content not implemented"
	}
}

func (m DashboardModel) renderStatsTab() string {
	leftPanel := m.renderStatsCards()
	rightPanel := m.list.View()

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)
}

func (m DashboardModel) renderStatsCards() string {
	totalCommands := len(m.commands)
	totalSessions := len(m.sessions)

	// Calculate unique commands
	uniqueCommands := make(map[string]bool)
	for _, cmd := range m.commands {
		uniqueCommands[cmd.Command] = true
	}

	// Create stat cards
	card1 := cardStyle.Render(fmt.Sprintf(
		"📊 Total Commands\n%s\n%d",
		highlightStyle.Render("━━━━━━━━━━━━"),
		totalCommands,
	))

	card2 := cardStyle.Render(fmt.Sprintf(
		"⭐ Unique Commands\n%s\n%d",
		highlightStyle.Render("━━━━━━━━━━━━"),
		len(uniqueCommands),
	))

	card3 := cardStyle.Render(fmt.Sprintf(
		"📱 Sessions\n%s\n%d",
		highlightStyle.Render("━━━━━━━━━━━━"),
		totalSessions,
	))

	// Progress card
	var progressCard string
	if m.gamificationStats != nil && m.gamificationStats.ProgressInfo != nil {
		progress := m.gamificationStats.ProgressInfo.ProgressPercent / 100.0
		progressBar := m.progress.ViewAs(progress)
		progressCard = cardStyle.Render(fmt.Sprintf(
			"🎯 Level Progress\n%s\n%s\n%.1f%% to next level",
			highlightStyle.Render("━━━━━━━━━━━━"),
			progressBar,
			m.gamificationStats.ProgressInfo.ProgressPercent,
		))
	} else {
		progressCard = cardStyle.Render(fmt.Sprintf(
			"🎯 Level Progress\n%s\nNo data available",
			highlightStyle.Render("━━━━━━━━━━━━"),
		))
	}

	// Arrange cards
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, card1, card2)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, card3, progressCard)

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

func (m DashboardModel) renderAnalyticsTab() string {
	if len(m.commands) == 0 {
		return cardStyle.Render("No analytics data available. Start using your terminal!")
	}

	analyzer := analytics.NewProductivityAnalyzer()
	metrics := analyzer.AnalyzeProductivity(m.commands, m.sessions)

	content := fmt.Sprintf(`📈 Productivity Analytics

🎯 Overall Score: %.1f/100
⚡ Efficiency: %.1f%%
🔄 Consistency: %.1f%%
🎪 Specialization: %s

📊 Time Distribution:
• Morning (6-12): %.1f%%
• Afternoon (12-18): %.1f%%
• Evening (18-24): %.1f%%
• Night (0-6): %.1f%%

💪 Peak Performance:
• Most productive day: %s
• Peak hours: %v
• Current streak: %d days`,
		metrics.OverallScore,
		metrics.EfficiencyMetrics.UniqueCommandRatio*100,
		metrics.StreakAnalysis.ConsistencyScore*100,
		metrics.CategoryInsights.SpecializationLevel,
		metrics.DailyPattern.MorningScore,
		metrics.DailyPattern.AfternoonScore,
		metrics.DailyPattern.EveningScore,
		metrics.DailyPattern.NightScore,
		metrics.WeeklyPattern.MostProductiveDay.String(),
		metrics.DailyPattern.PeakHours,
		metrics.StreakAnalysis.CurrentStreak,
	)

	return cardStyle.Width(m.width - 4).Render(content)
}

func (m DashboardModel) renderHeatmapTab() string {
	if len(m.commands) == 0 {
		return cardStyle.Render("No heatmap data available. Start using your terminal!")
	}

	heatmapAnalyzer := analytics.NewHeatmapAnalyzer()
	heatmapData := heatmapAnalyzer.GenerateHeatmap(m.commands)
	content := heatmapAnalyzer.FormatHeatmapVisualization(heatmapData)

	m.viewport.SetContent(content)
	return m.viewport.View()
}

func (m DashboardModel) renderAchievementsTab() string {
	if m.gamificationStats == nil {
		return cardStyle.Render("No achievement data available.")
	}

	header := titleStyle.Render(fmt.Sprintf("🏆 Achievements (%d unlocked)",
		len(m.gamificationStats.Achievements)))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.table.View(),
	)
}

func (m DashboardModel) renderSettingsTab() string {
	content := `⚙️ Settings

🎨 Theme: Default
📊 Stats Display: Enabled
🔔 Notifications: Enabled
💾 Auto-save: Enabled

🎮 Gamification:
• XP Multiplier: 1.0x
• Achievement Notifications: On
• Progress Tracking: On

📈 Analytics:
• Data Collection: On
• Productivity Tracking: On
• Time Pattern Analysis: On

Press Enter to modify settings`

	return cardStyle.Width(m.width - 4).Render(content)
}

func (m DashboardModel) renderHelp() string {
	helpText := "←/→ navigate tabs • q quit • ? help"

	switch m.currentTab {
	case StatsTab:
		helpText += " • ↑/↓ scroll commands"
	case HeatmapTab:
		helpText += " • ↑/↓ scroll heatmap"
	case AchievementsTab:
		helpText += " • ↑/↓ browse achievements"
	}

	return helpStyle.Render(helpText)
}

func (m DashboardModel) createProgressBar(progress float64, width int) string {
	filled := int(progress * float64(width))
	empty := width - filled

	return progressStyle.Render(strings.Repeat("█", filled) + strings.Repeat("░", empty))
}

// ListItem represents an item in the list
type ListItem struct {
	title string
	desc  string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

// RunDashboard starts the interactive dashboard
func RunDashboard(db *database.DB) error {
	model, err := NewDashboardModel(db)
	if err != nil {
		return fmt.Errorf("failed to create dashboard model: %w", err)
	}

	program := tea.NewProgram(
		*model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run dashboard: %w", err)
	}

	return nil
}
