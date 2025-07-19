package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the visual styling for the dashboard
type Theme struct {
	// Colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Muted     lipgloss.Color

	// Component styles
	ContentBox    lipgloss.Style
	SectionBox    lipgloss.Style
	SectionTitle  lipgloss.Style
	StatBox       lipgloss.Style
	UserInfo      lipgloss.Style
	ActivityItem  lipgloss.Style
	ProgressBar   lipgloss.Style
	ErrorBox      lipgloss.Style
	EmptyState    lipgloss.Style
	TabActive     lipgloss.Style
	TabInactive   lipgloss.Style
	Header        lipgloss.Style
	Footer        lipgloss.Style
}

// NewSpaceTheme creates a space-themed dashboard theme
func NewSpaceTheme() *Theme {
	theme := &Theme{
		// Space-inspired color palette
		Primary:   lipgloss.Color("#00D4FF"), // Cyan blue
		Secondary: lipgloss.Color("#7C3AED"), // Purple
		Accent:    lipgloss.Color("#F59E0B"), // Amber
		Success:   lipgloss.Color("#10B981"), // Green
		Warning:   lipgloss.Color("#F59E0B"), // Amber
		Error:     lipgloss.Color("#EF4444"), // Red
		Muted:     lipgloss.Color("#6B7280"), // Gray
	}

	// Base styles
	baseBox := lipgloss.NewStyle().
		Padding(1).
		Margin(0, 1)

	baseBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary)

	// Component styles
	theme.ContentBox = baseBox.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	theme.SectionBox = baseBox.Copy().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Secondary).
		Margin(1, 0)

	theme.SectionTitle = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Underline(true)

	theme.StatBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Accent).
		Padding(1).
		Margin(0, 1).
		Align(lipgloss.Center).
		Width(15)

	theme.UserInfo = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	theme.ActivityItem = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Margin(0, 1)

	theme.ProgressBar = lipgloss.NewStyle().
		Foreground(theme.Success).
		Background(theme.Muted)

	theme.ErrorBox = lipgloss.NewStyle().
		Foreground(theme.Error).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Error).
		Padding(1).
		Margin(1)

	theme.EmptyState = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Italic(true).
		Align(lipgloss.Center)

	theme.TabActive = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Background(theme.Secondary).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	theme.TabInactive = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 2).
		Margin(0, 1)

	theme.Header = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(0, 0, 1, 0)

	theme.Footer = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(1, 0, 0, 0)

	return theme
}

// NewCyberpunkTheme creates a cyberpunk-themed dashboard theme
func NewCyberpunkTheme() *Theme {
	theme := &Theme{
		// Cyberpunk color palette
		Primary:   lipgloss.Color("#00FF41"), // Matrix green
		Secondary: lipgloss.Color("#FF0080"), // Hot pink
		Accent:    lipgloss.Color("#00FFFF"), // Cyan
		Success:   lipgloss.Color("#00FF41"), // Green
		Warning:   lipgloss.Color("#FFFF00"), // Yellow
		Error:     lipgloss.Color("#FF0040"), // Red
		Muted:     lipgloss.Color("#808080"), // Gray
	}

	// Base styles with sharper edges
	baseBox := lipgloss.NewStyle().
		Padding(1).
		Margin(0, 1)

	// Component styles with angular borders
	theme.ContentBox = baseBox.Copy().
		Border(lipgloss.ThickBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	theme.SectionBox = baseBox.Copy().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(theme.Secondary).
		Margin(1, 0)

	theme.SectionTitle = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Underline(true).
		Background(lipgloss.Color("#001100"))

	theme.StatBox = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(theme.Accent).
		Padding(1).
		Margin(0, 1).
		Align(lipgloss.Center).
		Width(15).
		Background(lipgloss.Color("#000011"))

	theme.UserInfo = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Background(lipgloss.Color("#001100"))

	theme.ActivityItem = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Margin(0, 1)

	theme.ProgressBar = lipgloss.NewStyle().
		Foreground(theme.Success).
		Background(lipgloss.Color("#001100"))

	theme.ErrorBox = lipgloss.NewStyle().
		Foreground(theme.Error).
		Border(lipgloss.ThickBorder()).
		BorderForeground(theme.Error).
		Padding(1).
		Margin(1).
		Background(lipgloss.Color("#110000"))

	theme.EmptyState = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Italic(true).
		Align(lipgloss.Center)

	theme.TabActive = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(theme.Primary).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	theme.TabInactive = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 2).
		Margin(0, 1)

	theme.Header = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Border(lipgloss.ThickBorder(), false, false, true, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(0, 0, 1, 0).
		Background(lipgloss.Color("#001100"))

	theme.Footer = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Border(lipgloss.ThickBorder(), true, false, false, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(1, 0, 0, 0)

	return theme
}

// NewMinimalTheme creates a minimal, clean dashboard theme
func NewMinimalTheme() *Theme {
	theme := &Theme{
		// Minimal color palette
		Primary:   lipgloss.Color("#333333"), // Dark gray
		Secondary: lipgloss.Color("#666666"), // Medium gray
		Accent:    lipgloss.Color("#0066CC"), // Blue
		Success:   lipgloss.Color("#00AA00"), // Green
		Warning:   lipgloss.Color("#FF8800"), // Orange
		Error:     lipgloss.Color("#CC0000"), // Red
		Muted:     lipgloss.Color("#999999"), // Light gray
	}

	// Minimal base styles
	baseBox := lipgloss.NewStyle().
		Padding(1).
		Margin(0, 1)

	// Clean, minimal component styles
	theme.ContentBox = baseBox.Copy().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2)

	theme.SectionBox = baseBox.Copy().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Secondary).
		Margin(1, 0)

	theme.SectionTitle = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	theme.StatBox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Accent).
		Padding(1).
		Margin(0, 1).
		Align(lipgloss.Center).
		Width(15)

	theme.UserInfo = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true)

	theme.ActivityItem = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Margin(0, 1)

	theme.ProgressBar = lipgloss.NewStyle().
		Foreground(theme.Success)

	theme.ErrorBox = lipgloss.NewStyle().
		Foreground(theme.Error).
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Error).
		Padding(1).
		Margin(1)

	theme.EmptyState = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Italic(true).
		Align(lipgloss.Center)

	theme.TabActive = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1).
		Underline(true)

	theme.TabInactive = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Padding(0, 2).
		Margin(0, 1)

	theme.Header = lipgloss.NewStyle().
		Foreground(theme.Primary).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(0, 0, 1, 0)

	theme.Footer = lipgloss.NewStyle().
		Foreground(theme.Muted).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(theme.Primary).
		Padding(1).
		Margin(1, 0, 0, 0)

	return theme
}

// GetTheme returns a theme by name
func GetTheme(name string) *Theme {
	switch name {
	case "space":
		return NewSpaceTheme()
	case "cyberpunk":
		return NewCyberpunkTheme()
	case "minimal":
		return NewMinimalTheme()
	default:
		return NewSpaceTheme() // Default theme
	}
}
