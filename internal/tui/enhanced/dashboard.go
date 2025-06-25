package enhanced

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oiahoon/termonaut/internal/avatar"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
)

// TabType represents different dashboard tabs
type TabType int

const (
	HomeTab TabType = iota
	AnalyticsTab
	GamificationTab
	ActivityTab
	ToolsTab
	SettingsTab
)

var tabNames = []string{
	"🏠 Home",
	"📊 Analytics", 
	"🎮 Gamification",
	"🔥 Activity",
	"🛠️ Tools",
	"⚙️ Settings",
}

// EnhancedDashboard represents the new enhanced TUI dashboard
type EnhancedDashboard struct {
	// Core components
	activeTab    TabType
	windowWidth  int
	windowHeight int
	
	// Data managers
	db           *database.DB
	statsCalc    *stats.StatsCalculator
	avatarMgr    *avatar.AvatarManager
	
	// UI components
	spinner      spinner.Model
	loading      bool
	
	// Current data
	userProgress *models.UserProgress
	basicStats   *stats.BasicStats
	avatar       *avatar.Avatar
	
	// Theme and styling
	theme        *Theme
	
	// Key bindings
	keyMap       KeyMap
	
	// Mode preference
	modePreference string // smart, compact, full, classic
}

// KeyMap defines keyboard shortcuts
type KeyMap struct {
	Quit        key.Binding
	Help        key.Binding
	Refresh     key.Binding
	NextTab     key.Binding
	PrevTab     key.Binding
	Settings    key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r", "f5"),
			key.WithHelp("r", "refresh"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab", "l", "right"),
			key.WithHelp("tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab", "h", "left"),
			key.WithHelp("shift+tab", "prev tab"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
	}
}

// NewEnhancedDashboard creates a new enhanced dashboard
func NewEnhancedDashboard(db *database.DB) *EnhancedDashboard {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	// Create avatar manager
	avatarMgr, _ := avatar.NewAvatarManager(&avatar.Config{
		CacheDir:     "~/.termonaut/cache/avatars",
		CacheTTL:     24 * time.Hour,
		APITimeout:   10 * time.Second,
		DefaultStyle: "pixel-art",
		DefaultSize:  avatar.SizeMedium,
	})
	
	return &EnhancedDashboard{
		activeTab:      HomeTab,
		db:             db,
		statsCalc:      stats.New(db),
		avatarMgr:      avatarMgr,
		spinner:        s,
		loading:        true,
		theme:          DefaultSpaceTheme(),
		keyMap:         DefaultKeyMap(),
		modePreference: "smart", // Default to smart mode
	}
}

// SetModePreference sets the preferred display mode
func (d *EnhancedDashboard) SetModePreference(mode string) {
	d.modePreference = mode
}

// Init initializes the dashboard
func (d *EnhancedDashboard) Init() tea.Cmd {
	return tea.Batch(
		d.spinner.Tick,
		d.loadInitialData(),
	)
}

// Update handles messages and updates the model
func (d *EnhancedDashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		oldWidth := d.windowWidth
		oldHeight := d.windowHeight
		
		d.windowWidth = msg.Width
		d.windowHeight = msg.Height
		
		// If window size changed significantly, reload avatar with optimal size
		if abs(oldWidth-msg.Width) > 20 || abs(oldHeight-msg.Height) > 10 {
			cmds = append(cmds, d.loadInitialData())
		}
		
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keyMap.Quit):
			return d, tea.Quit
			
		case key.Matches(msg, d.keyMap.NextTab):
			d.activeTab = (d.activeTab + 1) % TabType(len(tabNames))
			
		case key.Matches(msg, d.keyMap.PrevTab):
			if d.activeTab == 0 {
				d.activeTab = TabType(len(tabNames) - 1)
			} else {
				d.activeTab--
			}
			
		case key.Matches(msg, d.keyMap.Refresh):
			d.loading = true
			cmds = append(cmds, d.loadInitialData())
			
		case key.Matches(msg, d.keyMap.Settings):
			d.activeTab = SettingsTab
		}
		
	case dataLoadedMsg:
		d.loading = false
		d.userProgress = msg.userProgress
		d.basicStats = msg.basicStats
		d.avatar = msg.avatar
		
	case spinner.TickMsg:
		d.spinner, cmd = d.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return d, tea.Batch(cmds...)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// View renders the dashboard
func (d *EnhancedDashboard) View() string {
	if d.loading {
		return d.renderLoading()
	}
	
	// Header
	header := d.renderHeader()
	
	// Tab navigation
	tabNav := d.renderTabNavigation()
	
	// Main content based on active tab
	var content string
	switch d.activeTab {
	case HomeTab:
		content = d.renderHomeTab()
	case AnalyticsTab:
		content = d.renderAnalyticsTab()
	case GamificationTab:
		content = d.renderGamificationTab()
	case ActivityTab:
		content = d.renderActivityTab()
	case ToolsTab:
		content = d.renderToolsTab()
	case SettingsTab:
		content = d.renderSettingsTab()
	}
	
	// Footer
	footer := d.renderFooter()
	
	// Combine all parts
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tabNav,
		content,
		footer,
	)
}

// renderLoading shows loading screen
func (d *EnhancedDashboard) renderLoading() string {
	loadingStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(d.windowWidth).
		Height(d.windowHeight).
		Foreground(d.theme.Colors.Primary)
		
	return loadingStyle.Render(
		fmt.Sprintf("%s Loading Termonaut Dashboard...", d.spinner.View()),
	)
}

// renderHeader renders the dashboard header
func (d *EnhancedDashboard) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(d.theme.Colors.Text).
		Background(d.theme.Colors.Primary).
		Padding(0, 2).
		Width(d.windowWidth).
		Bold(true)
		
	title := "🚀 Termonaut Dashboard"
	if d.userProgress != nil {
		title = fmt.Sprintf("🚀 Termonaut - Level %d Space Commander", d.userProgress.CurrentLevel)
	}
	
	return headerStyle.Render(title)
}

// renderTabNavigation renders tab navigation
func (d *EnhancedDashboard) renderTabNavigation() string {
	var tabs []string
	
	for i, name := range tabNames {
		style := lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 1)
			
		if TabType(i) == d.activeTab {
			style = style.
				Foreground(d.theme.Colors.Background).
				Background(d.theme.Colors.Accent).
				Bold(true)
		} else {
			style = style.
				Foreground(d.theme.Colors.Text).
				Background(d.theme.Colors.Surface)
		}
		
		tabs = append(tabs, style.Render(name))
	}
	
	tabBar := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	
	return lipgloss.NewStyle().
		Padding(1, 0).
		Width(d.windowWidth).
		Render(tabBar)
}

// renderHomeTab renders the home tab content
func (d *EnhancedDashboard) renderHomeTab() string {
	if d.userProgress == nil || d.basicStats == nil {
		return "Loading..."
	}
	
	// Calculate layout based on window width
	useWideLayout := d.windowWidth >= 100
	
	if useWideLayout {
		return d.renderHomeTabWide()
	}
	return d.renderHomeTabNarrow()
}

// renderHomeTabWide renders home tab for wide terminals
func (d *EnhancedDashboard) renderHomeTabWide() string {
	// Left column: Avatar
	avatarSection := d.renderAvatarSection()
	
	// Right column: Stats
	statsSection := d.renderQuickStats()
	
	// Top row
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		avatarSection,
		statsSection,
	)
	
	// Bottom row: Recent commands
	recentCommands := d.renderRecentCommands()
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		"",
		recentCommands,
	)
}

// renderHomeTabNarrow renders home tab for narrow terminals
func (d *EnhancedDashboard) renderHomeTabNarrow() string {
	sections := []string{
		d.renderQuickStats(),
		"",
		d.renderRecentCommands(),
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderAvatarSection renders the avatar display
func (d *EnhancedDashboard) renderAvatarSection() string {
	// Calculate avatar width based on terminal size
	avatarWidth := d.calculateAvatarWidth()
	avatarHeight := d.calculateAvatarHeight()
	
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(d.theme.Colors.Primary).
		Padding(1).
		Width(avatarWidth).
		Height(avatarHeight).
		Align(lipgloss.Center, lipgloss.Center)
	
	content := d.renderAvatarContent()
	
	return cardStyle.Render(content)
}

// calculateAvatarWidth determines the optimal avatar width based on terminal size and mode preference
func (d *EnhancedDashboard) calculateAvatarWidth() int {
	// Apply mode preference override
	switch d.modePreference {
	case "compact":
		// Force compact sizes regardless of terminal size
		if d.windowWidth >= 80 {
			return 25
		} else if d.windowWidth >= 60 {
			return 20
		}
		return 15
	case "full":
		// Force large sizes when possible
		if d.windowWidth >= 100 {
			return 70
		} else if d.windowWidth >= 80 {
			return 55
		}
		return 45
	case "smart":
		// Intelligent sizing based on terminal size
		if d.windowWidth >= 140 {
			return 70 // Extra large avatar for very wide terminals
		} else if d.windowWidth >= 120 {
			return 65 // Large avatar for wide terminals
		} else if d.windowWidth >= 100 {
			return 55 // Medium-large avatar for standard terminals
		} else if d.windowWidth >= 80 {
			return 45 // Medium avatar for narrow terminals
		} else if d.windowWidth >= 60 {
			return 25 // Small avatar for compact terminals
		} else if d.windowWidth >= 40 {
			return 15 // Mini avatar for very small terminals
		}
		return 12 // Ultra-compact avatar for minimal terminals
	default:
		// Default to smart mode behavior
		return d.calculateSmartAvatarWidth()
	}
}

// calculateSmartAvatarWidth implements the smart sizing logic
func (d *EnhancedDashboard) calculateSmartAvatarWidth() int {
	if d.windowWidth >= 140 {
		return 70
	} else if d.windowWidth >= 120 {
		return 65
	} else if d.windowWidth >= 100 {
		return 55
	} else if d.windowWidth >= 80 {
		return 45
	} else if d.windowWidth >= 60 {
		return 25
	} else if d.windowWidth >= 40 {
		return 15
	}
	return 12
}

// calculateAvatarHeight determines the optimal avatar height based on terminal size
func (d *EnhancedDashboard) calculateAvatarHeight() int {
	if d.windowHeight >= 35 {
		return 25 // Extra large avatar for very tall terminals
	} else if d.windowHeight >= 30 {
		return 22 // Large avatar for tall terminals
	} else if d.windowHeight >= 25 {
		return 20 // Medium-large avatar for standard terminals
	} else if d.windowHeight >= 20 {
		return 18 // Medium avatar for short terminals
	} else if d.windowHeight >= 15 {
		return 12 // Small avatar for compact terminals
	} else if d.windowHeight >= 10 {
		return 8  // Mini avatar for very small terminals
	}
	return 6 // Ultra-compact avatar for minimal terminals
}

// renderAvatarContent renders the avatar content with proper sizing
func (d *EnhancedDashboard) renderAvatarContent() string {
	if d.avatar == nil || d.avatar.ASCIIArt == "" {
		return d.renderDefaultAvatar()
	}
	
	// If we have a real avatar, use it directly
	avatarArt := d.avatar.ASCIIArt
	
	// Add level and user info below avatar
	levelInfo := ""
	if d.userProgress != nil {
		levelInfo = fmt.Sprintf("\nLevel %d", d.userProgress.CurrentLevel)
		if d.userProgress.TotalXP > 0 {
			levelInfo += fmt.Sprintf("\n%d XP", d.userProgress.TotalXP)
		}
	}
	
	// Combine avatar art with level info
	content := avatarArt + levelInfo
	
	return content
}

// getOptimalAvatarSize returns the best avatar size for current terminal
func (d *EnhancedDashboard) getOptimalAvatarSize() avatar.AvatarSize {
	if d.windowWidth >= 140 && d.windowHeight >= 35 {
		// Extra large custom size for very wide terminals
		return avatar.AvatarSize{SVGSize: 256, ASCIIWidth: 65, ASCIIHeight: 32}
	} else if d.windowWidth >= 120 && d.windowHeight >= 30 {
		// Large custom size for wide terminals  
		return avatar.AvatarSize{SVGSize: 256, ASCIIWidth: 60, ASCIIHeight: 30}
	} else if d.windowWidth >= 100 && d.windowHeight >= 25 {
		// Medium-large custom size for standard terminals
		return avatar.AvatarSize{SVGSize: 128, ASCIIWidth: 50, ASCIIHeight: 25}
	} else if d.windowWidth >= 80 && d.windowHeight >= 20 {
		return avatar.SizeMedium // 40x20
	} else if d.windowWidth >= 60 && d.windowHeight >= 15 {
		return avatar.SizeSmall  // 20x10
	} else if d.windowWidth >= 40 && d.windowHeight >= 10 {
		return avatar.SizeMini   // 10x5
	}
	// Ultra-compact size for very small terminals
	return avatar.AvatarSize{SVGSize: 32, ASCIIWidth: 8, ASCIIHeight: 4}
}

// renderDefaultAvatar creates a default avatar when none is available
func (d *EnhancedDashboard) renderDefaultAvatar() string {
	// Create different sized default avatars based on terminal size
	if d.windowWidth >= 140 {
		return d.renderExtraLargeDefaultAvatar()
	} else if d.windowWidth >= 120 {
		return d.renderLargeDefaultAvatar()
	} else if d.windowWidth >= 100 {
		return d.renderMediumLargeDefaultAvatar()
	} else if d.windowWidth >= 80 {
		return d.renderMediumDefaultAvatar()
	} else if d.windowWidth >= 60 {
		return d.renderSmallDefaultAvatar()
	} else if d.windowWidth >= 40 {
		return d.renderMiniDefaultAvatar()
	}
	return d.renderUltraCompactDefaultAvatar()
}

// renderMiniDefaultAvatar creates a mini default avatar
func (d *EnhancedDashboard) renderMiniDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`🚀
/|\
T
||

L%d`, level)
}

// renderUltraCompactDefaultAvatar creates an ultra-compact default avatar
func (d *EnhancedDashboard) renderUltraCompactDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`🚀
T
%d`, level)
}

// renderExtraLargeDefaultAvatar creates an extra large default avatar
func (d *EnhancedDashboard) renderExtraLargeDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`                🚀
               /|\
              / | \
             /  |  \
            |   T   |
            |       |
            |       |
            ||     ||
            ||     ||
            /\     /\
           /  \   /  \
          /    \ /    \
         /_____\_/_____\
         
         Level %d
      Space Commander
      
    "To infinity and beyond!"`, level)
}

// renderMediumLargeDefaultAvatar creates a medium-large default avatar
func (d *EnhancedDashboard) renderMediumLargeDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`           🚀
          /|\
         / | \
        /  |  \
       |   T   |
       |       |
       ||     ||
       ||     ||
       /\     /\
      /  \   /  \
     /_____\_____\
     
     Level %d
  Space Commander`, level)
}

// renderLargeDefaultAvatar creates a large default avatar
func (d *EnhancedDashboard) renderLargeDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`         🚀
        /|\
       / | \
      /  |  \
     |   T   |
     |       |
     ||     ||
     /\     /\
    /  \   /  \
   /__________\
   
   Level %d
Space Commander`, level)
}

// renderMediumDefaultAvatar creates a medium default avatar
func (d *EnhancedDashboard) renderMediumDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`    🚀
   /|\
  / | \
 |  T  |
 |     |
 ||   ||
 /\   /\
 
Level %d`, level)
}

// renderSmallDefaultAvatar creates a small default avatar
func (d *EnhancedDashboard) renderSmallDefaultAvatar() string {
	level := 1
	if d.userProgress != nil {
		level = d.userProgress.CurrentLevel
	}
	
	return fmt.Sprintf(`  🚀
 /|\
/ | \
|  T |
||  ||

Lv.%d`, level)
}

// renderQuickStats renders quick statistics
func (d *EnhancedDashboard) renderQuickStats() string {
	// Calculate stats width to complement avatar width
	statsWidth := d.calculateStatsWidth()
	statsHeight := d.calculateAvatarHeight() // Match avatar height
	
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(d.theme.Colors.Secondary).
		Padding(1).
		Width(statsWidth).
		Height(statsHeight)
	
	// Handle loading state
	if d.basicStats == nil || d.userProgress == nil {
		content := lipgloss.JoinVertical(lipgloss.Left, 
			"Loading statistics...",
			"",
			"Please wait while we",
			"fetch your data.",
		)
		return cardStyle.Render(content)
	}
	
	// Safe value extraction with defaults
	commandsToday := 0
	totalCommands := 0
	currentStreak := 0
	totalXP := 0
	
	if d.basicStats != nil {
		commandsToday = d.basicStats.CommandsToday
		totalCommands = d.basicStats.TotalCommands
	}
	
	if d.userProgress != nil {
		currentStreak = d.userProgress.CurrentStreak
		totalXP = d.userProgress.TotalXP
	}
	
	// Create stats content with proper spacing
	stats := []string{
		"📊 Today's Activity",
		"",
		fmt.Sprintf("Commands: %d 🎯", commandsToday),
		fmt.Sprintf("Total: %d", totalCommands),
		fmt.Sprintf("Streak: %d days 🔥", currentStreak),
		"",
		"🎮 Progress",
		fmt.Sprintf("XP: %s", d.formatXP(totalXP)),
		d.renderXPProgress(),
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left, stats...)
	return cardStyle.Render(content)
}

// calculateStatsWidth determines the optimal stats width based on terminal size
func (d *EnhancedDashboard) calculateStatsWidth() int {
	avatarWidth := d.calculateAvatarWidth()
	remainingWidth := d.windowWidth - avatarWidth - 10 // Account for borders, margins, and spacing
	
	// Ensure stats area has reasonable width
	if remainingWidth > 80 {
		return 80 // Maximum stats width for very wide terminals
	} else if remainingWidth > 60 {
		return remainingWidth
	} else if remainingWidth > 45 {
		return remainingWidth
	}
	return 45 // Minimum stats width
}

// formatXP formats XP numbers with proper formatting
func (d *EnhancedDashboard) formatXP(xp int) string {
	if xp >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(xp)/1000000)
	} else if xp >= 1000 {
		return fmt.Sprintf("%.1fK", float64(xp)/1000)
	}
	return fmt.Sprintf("%d", xp)
}

// renderXPProgress renders XP progress bar
func (d *EnhancedDashboard) renderXPProgress() string {
	// Handle nil userProgress
	if d.userProgress == nil {
		return "Progress: Loading..."
	}
	
	// Calculate XP needed for next level
	currentLevel := d.userProgress.CurrentLevel
	currentXP := d.userProgress.TotalXP
	
	// Handle edge cases
	if currentLevel < 0 {
		currentLevel = 0
	}
	if currentXP < 0 {
		currentXP = 0
	}
	
	// Simple XP calculation (can be enhanced)
	xpForCurrentLevel := currentLevel * 1000
	xpForNextLevel := (currentLevel + 1) * 1000
	xpProgress := currentXP - xpForCurrentLevel
	xpNeeded := xpForNextLevel - xpForCurrentLevel
	
	// Ensure positive values
	if xpProgress < 0 {
		xpProgress = 0
	}
	if xpNeeded <= 0 {
		xpNeeded = 1000 // Default to 1000 XP per level
	}
	
	percentage := float64(xpProgress) / float64(xpNeeded)
	if percentage > 1 {
		percentage = 1
	}
	if percentage < 0 {
		percentage = 0
	}
	
	// Progress bar
	barWidth := 30
	filledWidth := int(float64(barWidth) * percentage)
	
	// Ensure valid width values
	if filledWidth < 0 {
		filledWidth = 0
	}
	if filledWidth > barWidth {
		filledWidth = barWidth
	}
	
	emptyWidth := barWidth - filledWidth
	if emptyWidth < 0 {
		emptyWidth = 0
	}
	
	// Create progress bar safely
	var filled, empty string
	if filledWidth > 0 {
		filled = lipgloss.NewStyle().
			Foreground(d.theme.Colors.Success).
			Render(strings.Repeat("█", filledWidth))
	}
	
	if emptyWidth > 0 {
		empty = lipgloss.NewStyle().
			Foreground(d.theme.Colors.Surface).
			Render(strings.Repeat("░", emptyWidth))
	}
	
	progressBar := filled + empty
	
	return fmt.Sprintf("Progress to Level %d: %s %.0f%%", 
		currentLevel+1, progressBar, percentage*100)
}

// renderRecentCommands renders recent command history
func (d *EnhancedDashboard) renderRecentCommands() string {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(d.theme.Colors.Accent).
		Padding(1).
		Width(d.windowWidth - 4)
	
	// Mock recent commands (replace with actual data)
	commands := []string{
		"git commit -m \"feat: enhance TUI dashboard\"     2m ago",
		"go build -o termonaut cmd/termonaut/*.go         5m ago", 
		"./termonaut tui                                  8m ago",
		"git status                                       12m ago",
	}
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"Recent Commands:",
		"",
	)
	
	for _, cmd := range commands {
		content += "\n" + lipgloss.NewStyle().
			Foreground(d.theme.Colors.Text).
			Render("  " + cmd)
	}
	
	return cardStyle.Render(content)
}

// renderFooter renders the help footer
func (d *EnhancedDashboard) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(d.theme.Colors.Text).
		Background(d.theme.Colors.Surface).
		Padding(0, 2).
		Width(d.windowWidth)
	
	help := "[Tab] Next • [Shift+Tab] Prev • [r] Refresh • [q] Quit"
	return footerStyle.Render(help)
}

// Placeholder methods for other tabs
func (d *EnhancedDashboard) renderAnalyticsTab() string {
	return "📊 Analytics Tab - Coming Soon!"
}

func (d *EnhancedDashboard) renderGamificationTab() string {
	return "🎮 Gamification Tab - Coming Soon!"
}

func (d *EnhancedDashboard) renderActivityTab() string {
	return "🔥 Activity Tab - Coming Soon!"
}

func (d *EnhancedDashboard) renderToolsTab() string {
	return "🛠️ Tools Tab - Coming Soon!"
}

func (d *EnhancedDashboard) renderSettingsTab() string {
	return "⚙️ Settings Tab - Coming Soon!"
}

// Data loading
type dataLoadedMsg struct {
	userProgress *models.UserProgress
	basicStats   *stats.BasicStats
	avatar       *avatar.Avatar
}

func (d *EnhancedDashboard) loadInitialData() tea.Cmd {
	return func() tea.Msg {
		// Load user progress with fallback
		progress, err := d.db.GetUserProgress()
		if err != nil || progress == nil {
			// Create default progress if none exists
			progress = &models.UserProgress{
				TotalXP:             0,
				CurrentLevel:        1,
				CommandsCount:       0,
				UniqueCommandsCount: 0,
				LongestStreak:       0,
				CurrentStreak:       0,
			}
		}
		
		// Load basic stats with fallback
		basicStats, err := d.statsCalc.GetBasicStats()
		if err != nil || basicStats == nil {
			// Create default stats if none exist
			basicStats = &stats.BasicStats{
				TotalCommands:    0,
				TotalSessions:    0,
				UniqueCommands:   0,
				CommandsToday:    0,
				MostUsedCommand:  "N/A",
				MostUsedCount:    0,
				TopCommands:      []map[string]interface{}{},
			}
		}
		
		// Generate avatar request with optimal size
		username := "user" // You might want to get this from config
		level := 1
		if progress != nil && progress.CurrentLevel > 0 {
			level = progress.CurrentLevel
		}
		
		// Get optimal avatar size based on current terminal dimensions
		optimalSize := d.getOptimalAvatarSize()
		
		// Load avatar with fallback
		var avatarResult *avatar.Avatar
		if d.avatarMgr != nil {
			avatarReq := avatar.AvatarRequest{
				Username: username,
				Level:    level,
				Style:    "pixel-art",
				Size:     optimalSize,
			}
			avatarResult, _ = d.avatarMgr.Generate(avatarReq)
		}
		
		// If avatar generation fails, we'll use the default avatar in renderAvatarContent
		// No need to create a mock here since renderDefaultAvatar handles it
		
		return dataLoadedMsg{
			userProgress: progress,
			basicStats:   basicStats,
			avatar:       avatarResult,
		}
	}
}

// Helper functions
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func generateMockAvatar() string {
	return `    🚀
   /|\
  / | \
 |  T  |
 |     |
 ||   ||
 /\   /\`
}
