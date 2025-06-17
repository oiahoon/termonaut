package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AdvancedProgress represents an enhanced progress bar with animations
type AdvancedProgress struct {
	progress.Model
	animated    bool
	value       float64
	target      float64
	animSpeed   float64
	label       string
	showPercent bool
	style       ProgressStyle
}

// ProgressStyle defines the visual style of the progress bar
type ProgressStyle struct {
	FilledColor   lipgloss.Color
	EmptyColor    lipgloss.Color
	PercentColor  lipgloss.Color
	LabelColor    lipgloss.Color
	BorderColor   lipgloss.Color
	GradientStart lipgloss.Color
	GradientEnd   lipgloss.Color
	UseGradient   bool
}

// DefaultProgressStyle returns a default purple gradient style
func DefaultProgressStyle() ProgressStyle {
	return ProgressStyle{
		FilledColor:   lipgloss.Color("#7C3AED"),
		EmptyColor:    lipgloss.Color("#2A2A2A"),
		PercentColor:  lipgloss.Color("#A855F7"),
		LabelColor:    lipgloss.Color("#E5E7EB"),
		BorderColor:   lipgloss.Color("#874BFD"),
		GradientStart: lipgloss.Color("#7C3AED"),
		GradientEnd:   lipgloss.Color("#EC4899"),
		UseGradient:   true,
	}
}

// GamificationProgressStyle returns a gaming-themed style
func GamificationProgressStyle() ProgressStyle {
	return ProgressStyle{
		FilledColor:   lipgloss.Color("#10B981"),
		EmptyColor:    lipgloss.Color("#1F2937"),
		PercentColor:  lipgloss.Color("#34D399"),
		LabelColor:    lipgloss.Color("#F9FAFB"),
		BorderColor:   lipgloss.Color("#059669"),
		GradientStart: lipgloss.Color("#10B981"),
		GradientEnd:   lipgloss.Color("#3B82F6"),
		UseGradient:   true,
	}
}

// XPProgressStyle returns an XP-themed style
func XPProgressStyle() ProgressStyle {
	return ProgressStyle{
		FilledColor:   lipgloss.Color("#F59E0B"),
		EmptyColor:    lipgloss.Color("#374151"),
		PercentColor:  lipgloss.Color("#FBBF24"),
		LabelColor:    lipgloss.Color("#FEF3C7"),
		BorderColor:   lipgloss.Color("#D97706"),
		GradientStart: lipgloss.Color("#F59E0B"),
		GradientEnd:   lipgloss.Color("#EF4444"),
		UseGradient:   true,
	}
}

// NewAdvancedProgress creates a new enhanced progress bar
func NewAdvancedProgress(width int, style ProgressStyle) *AdvancedProgress {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
		progress.WithoutPercentage(),
	)

	if style.UseGradient {
		prog = progress.New(
			progress.WithGradient(string(style.GradientStart), string(style.GradientEnd)),
			progress.WithWidth(width),
			progress.WithoutPercentage(),
		)
	} else {
		prog = progress.New(
			progress.WithSolidFill(string(style.FilledColor)),
			progress.WithWidth(width),
			progress.WithoutPercentage(),
		)
	}

	return &AdvancedProgress{
		Model:       prog,
		animated:    false,
		animSpeed:   0.02,
		showPercent: true,
		style:       style,
	}
}

// SetAnimated enables/disables smooth animations
func (p *AdvancedProgress) SetAnimated(animated bool) {
	p.animated = animated
}

// SetTarget sets the target value for animated progress
func (p *AdvancedProgress) SetTarget(target float64) {
	if target < 0 {
		target = 0
	}
	if target > 1 {
		target = 1
	}
	p.target = target
}

// SetValue sets the current value (immediately if not animated)
func (p *AdvancedProgress) SetValue(value float64) {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}

	if p.animated {
		p.target = value
	} else {
		p.value = value
	}
}

// SetLabel sets the progress bar label
func (p *AdvancedProgress) SetLabel(label string) {
	p.label = label
}

// SetShowPercent controls whether to show percentage
func (p *AdvancedProgress) SetShowPercent(show bool) {
	p.showPercent = show
}

// Update handles progress animations
func (p *AdvancedProgress) Update(msg tea.Msg) (*AdvancedProgress, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if p.animated && p.value != p.target {
		// Smooth animation towards target
		diff := p.target - p.value
		if math.Abs(diff) < p.animSpeed {
			p.value = p.target
		} else {
			if diff > 0 {
				p.value += p.animSpeed
			} else {
				p.value -= p.animSpeed
			}
		}

		// Continue animation
		cmds = append(cmds, p.tickCmd())
	}

	// Update the underlying progress model
	model, cmd := p.Model.Update(msg)
	if progressModel, ok := model.(progress.Model); ok {
		p.Model = progressModel
	}
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

// View renders the enhanced progress bar
func (p *AdvancedProgress) View() string {
	// Create the progress bar view with current value
	progressBar := p.Model.ViewAs(p.value)

	// Add label and percentage if needed
	var result strings.Builder

	if p.label != "" {
		labelStyle := lipgloss.NewStyle().Foreground(p.style.LabelColor)
		result.WriteString(labelStyle.Render(p.label))
		result.WriteString("\n")
	}

	result.WriteString(progressBar)

	if p.showPercent {
		percentStyle := lipgloss.NewStyle().Foreground(p.style.PercentColor)
		percentage := fmt.Sprintf(" %.1f%%", p.value*100)
		result.WriteString(percentStyle.Render(percentage))
	}

	return result.String()
}

// ViewWithStats renders progress with additional statistics
func (p *AdvancedProgress) ViewWithStats(current, total int, timeRemaining string) string {
	var result strings.Builder

	if p.label != "" {
		labelStyle := lipgloss.NewStyle().Foreground(p.style.LabelColor).Bold(true)
		result.WriteString(labelStyle.Render(p.label))
		result.WriteString("\n")
	}

	// Progress bar
	progressBar := p.Model.ViewAs(p.value)
	result.WriteString(progressBar)

	// Stats line
	statsStyle := lipgloss.NewStyle().Foreground(p.style.PercentColor)
	stats := fmt.Sprintf(" %d/%d", current, total)

	if p.showPercent {
		stats += fmt.Sprintf(" (%.1f%%)", p.value*100)
	}

	if timeRemaining != "" {
		stats += fmt.Sprintf(" • %s remaining", timeRemaining)
	}

	result.WriteString(statsStyle.Render(stats))

	return result.String()
}

// tickCmd returns a command for animation ticks
func (p *AdvancedProgress) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return ProgressTickMsg{Time: t}
	})
}

// ProgressTickMsg represents a progress animation tick
type ProgressTickMsg struct {
	Time time.Time
}

// CreateXPProgressBar creates a styled XP progress bar
func CreateXPProgressBar(currentXP, levelXP, totalXP int, level int, title string) string {
	style := XPProgressStyle()

	// Calculate progress within current level
	progress := float64(currentXP) / float64(levelXP)
	if progress > 1 {
		progress = 1
	}

	// Create progress bar
	width := 30
	filled := int(progress * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(style.FilledColor)
	emptyStyle := lipgloss.NewStyle().Foreground(style.EmptyColor)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	// Create the complete XP bar with styling
	labelStyle := lipgloss.NewStyle().Foreground(style.LabelColor).Bold(true)
	percentStyle := lipgloss.NewStyle().Foreground(style.PercentColor)

	result := fmt.Sprintf("%s\n[%s] %s\nXP: %d/%d | Total: %d | Level %d",
		labelStyle.Render(title),
		bar,
		percentStyle.Render(fmt.Sprintf("%.1f%%", progress*100)),
		currentXP,
		levelXP,
		totalXP,
		level,
	)

	return cardStyle.Render(result)
}

// CreateAchievementProgressBar creates a styled achievement progress bar
func CreateAchievementProgressBar(current, target int, achievementName, icon string) string {
	style := GamificationProgressStyle()

	progress := float64(current) / float64(target)
	if progress > 1 {
		progress = 1
	}

	// Create mini progress bar
	width := 15
	filled := int(progress * float64(width))
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(style.FilledColor)
	emptyStyle := lipgloss.NewStyle().Foreground(style.EmptyColor)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	// Format with icon and name
	percentStyle := lipgloss.NewStyle().Foreground(style.PercentColor)

	return fmt.Sprintf("%s %s [%s] %s",
		icon,
		achievementName,
		bar,
		percentStyle.Render(fmt.Sprintf("%.0f%%", progress*100)),
	)
}

// CreateLevelProgressIndicator creates a level-up progress indicator
func CreateLevelProgressIndicator(currentLevel, nextLevel int, progress float64) string {
	style := DefaultProgressStyle()

	width := 25
	filled := int(progress * float64(width))

	// Create gradient effect manually for level progress
	levelBar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			// Gradient from start to end color
			ratio := float64(i) / float64(width)
			if ratio < 0.5 {
				levelBar += lipgloss.NewStyle().Foreground(style.GradientStart).Render("█")
			} else {
				levelBar += lipgloss.NewStyle().Foreground(style.GradientEnd).Render("█")
			}
		} else {
			levelBar += lipgloss.NewStyle().Foreground(style.EmptyColor).Render("░")
		}
	}

	labelStyle := lipgloss.NewStyle().Foreground(style.LabelColor).Bold(true)

	return fmt.Sprintf("%s\n[%s]\nLevel %d → Level %d",
		labelStyle.Render("🌟 Level Progress"),
		levelBar,
		currentLevel,
		nextLevel,
	)
}
