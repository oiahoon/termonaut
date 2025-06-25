package enhanced

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a complete UI theme
type Theme struct {
	Name        string
	Colors      ColorScheme
	Typography  Typography
	Spacing     Spacing
	Animations  AnimationConfig
}

// ColorScheme defines the color palette for a theme
type ColorScheme struct {
	Primary     lipgloss.Color  // 主色调
	Secondary   lipgloss.Color  // 辅助色
	Accent      lipgloss.Color  // 强调色
	Background  lipgloss.Color  // 背景色
	Surface     lipgloss.Color  // 表面色
	Text        lipgloss.Color  // 文本色
	TextMuted   lipgloss.Color  // 次要文本色
	Success     lipgloss.Color  // 成功色
	Warning     lipgloss.Color  // 警告色
	Error       lipgloss.Color  // 错误色
	Info        lipgloss.Color  // 信息色
}

// Typography defines text styling
type Typography struct {
	HeaderSize   int
	BodySize     int
	CaptionSize  int
	LineHeight   float64
}

// Spacing defines layout spacing
type Spacing struct {
	Small   int
	Medium  int
	Large   int
	XLarge  int
}

// AnimationConfig defines animation settings
type AnimationConfig struct {
	Duration    int  // milliseconds
	Enabled     bool
	EasingType  string
}

// DefaultSpaceTheme returns the default space theme
func DefaultSpaceTheme() *Theme {
	return &Theme{
		Name: "Space",
		Colors: ColorScheme{
			Primary:     lipgloss.Color("#7C3AED"),  // 紫色
			Secondary:   lipgloss.Color("#3B82F6"),  // 蓝色
			Accent:      lipgloss.Color("#F59E0B"),  // 金色
			Background:  lipgloss.Color("#0F172A"),  // 深蓝
			Surface:     lipgloss.Color("#1E293B"),  // 灰蓝
			Text:        lipgloss.Color("#F8FAFC"),  // 白色
			TextMuted:   lipgloss.Color("#94A3B8"),  // 灰色
			Success:     lipgloss.Color("#10B981"),  // 绿色
			Warning:     lipgloss.Color("#F59E0B"),  // 橙色
			Error:       lipgloss.Color("#EF4444"),  // 红色
			Info:        lipgloss.Color("#06B6D4"),  // 青色
		},
		Typography: Typography{
			HeaderSize:  16,
			BodySize:    14,
			CaptionSize: 12,
			LineHeight:  1.5,
		},
		Spacing: Spacing{
			Small:  1,
			Medium: 2,
			Large:  4,
			XLarge: 8,
		},
		Animations: AnimationConfig{
			Duration:   300,
			Enabled:    true,
			EasingType: "ease-in-out",
		},
	}
}

// CyberpunkTheme returns a cyberpunk-style theme
func CyberpunkTheme() *Theme {
	return &Theme{
		Name: "Cyberpunk",
		Colors: ColorScheme{
			Primary:     lipgloss.Color("#FF0080"),  // 霓虹粉
			Secondary:   lipgloss.Color("#00FFFF"),  // 青色
			Accent:      lipgloss.Color("#FFFF00"),  // 黄色
			Background:  lipgloss.Color("#000000"),  // 纯黑
			Surface:     lipgloss.Color("#1A1A1A"),  // 深灰
			Text:        lipgloss.Color("#00FF00"),  // 绿色
			TextMuted:   lipgloss.Color("#808080"),  // 灰色
			Success:     lipgloss.Color("#00FF00"),  // 绿色
			Warning:     lipgloss.Color("#FFFF00"),  // 黄色
			Error:       lipgloss.Color("#FF0000"),  // 红色
			Info:        lipgloss.Color("#00FFFF"),  // 青色
		},
		Typography: Typography{
			HeaderSize:  16,
			BodySize:    14,
			CaptionSize: 12,
			LineHeight:  1.4,
		},
		Spacing: Spacing{
			Small:  1,
			Medium: 2,
			Large:  4,
			XLarge: 8,
		},
		Animations: AnimationConfig{
			Duration:   200,
			Enabled:    true,
			EasingType: "ease-out",
		},
	}
}

// MinimalTheme returns a minimal, clean theme
func MinimalTheme() *Theme {
	return &Theme{
		Name: "Minimal",
		Colors: ColorScheme{
			Primary:     lipgloss.Color("#000000"),  // 黑色
			Secondary:   lipgloss.Color("#666666"),  // 深灰
			Accent:      lipgloss.Color("#0066CC"),  // 蓝色
			Background:  lipgloss.Color("#FFFFFF"),  // 白色
			Surface:     lipgloss.Color("#F5F5F5"),  // 浅灰
			Text:        lipgloss.Color("#000000"),  // 黑色
			TextMuted:   lipgloss.Color("#666666"),  // 灰色
			Success:     lipgloss.Color("#00AA00"),  // 绿色
			Warning:     lipgloss.Color("#FF8800"),  // 橙色
			Error:       lipgloss.Color("#CC0000"),  // 红色
			Info:        lipgloss.Color("#0066CC"),  // 蓝色
		},
		Typography: Typography{
			HeaderSize:  16,
			BodySize:    14,
			CaptionSize: 12,
			LineHeight:  1.6,
		},
		Spacing: Spacing{
			Small:  1,
			Medium: 2,
			Large:  4,
			XLarge: 6,
		},
		Animations: AnimationConfig{
			Duration:   400,
			Enabled:    false,
			EasingType: "linear",
		},
	}
}

// RetroTheme returns a retro/vintage theme
func RetroTheme() *Theme {
	return &Theme{
		Name: "Retro",
		Colors: ColorScheme{
			Primary:     lipgloss.Color("#FF6B35"),  // 橙红
			Secondary:   lipgloss.Color("#F7931E"),  // 橙色
			Accent:      lipgloss.Color("#FFD23F"),  // 黄色
			Background:  lipgloss.Color("#2D1B69"),  // 深紫
			Surface:     lipgloss.Color("#3E2A7A"),  // 紫色
			Text:        lipgloss.Color("#FFFFFF"),  // 白色
			TextMuted:   lipgloss.Color("#B8A9D9"),  // 浅紫
			Success:     lipgloss.Color("#4ECDC4"),  // 青绿
			Warning:     lipgloss.Color("#FFE66D"),  // 浅黄
			Error:       lipgloss.Color("#FF6B6B"),  // 浅红
			Info:        lipgloss.Color("#4ECDC4"),  // 青绿
		},
		Typography: Typography{
			HeaderSize:  18,
			BodySize:    14,
			CaptionSize: 12,
			LineHeight:  1.5,
		},
		Spacing: Spacing{
			Small:  2,
			Medium: 3,
			Large:  5,
			XLarge: 8,
		},
		Animations: AnimationConfig{
			Duration:   500,
			Enabled:    true,
			EasingType: "ease-in-out",
		},
	}
}

// NatureTheme returns a nature-inspired theme
func NatureTheme() *Theme {
	return &Theme{
		Name: "Nature",
		Colors: ColorScheme{
			Primary:     lipgloss.Color("#2D5016"),  // 深绿
			Secondary:   lipgloss.Color("#4F7942"),  // 中绿
			Accent:      lipgloss.Color("#8FBC8F"),  // 浅绿
			Background:  lipgloss.Color("#F0F8E8"),  // 浅绿白
			Surface:     lipgloss.Color("#E8F5E8"),  // 极浅绿
			Text:        lipgloss.Color("#2D5016"),  // 深绿
			TextMuted:   lipgloss.Color("#6B8E5A"),  // 中绿
			Success:     lipgloss.Color("#228B22"),  // 森林绿
			Warning:     lipgloss.Color("#DAA520"),  // 金黄
			Error:       lipgloss.Color("#B22222"),  // 火砖红
			Info:        lipgloss.Color("#4682B4"),  // 钢蓝
		},
		Typography: Typography{
			HeaderSize:  16,
			BodySize:    14,
			CaptionSize: 12,
			LineHeight:  1.6,
		},
		Spacing: Spacing{
			Small:  1,
			Medium: 2,
			Large:  4,
			XLarge: 6,
		},
		Animations: AnimationConfig{
			Duration:   350,
			Enabled:    true,
			EasingType: "ease-in-out",
		},
	}
}

// GetAllThemes returns all available themes
func GetAllThemes() map[string]*Theme {
	return map[string]*Theme{
		"space":     DefaultSpaceTheme(),
		"cyberpunk": CyberpunkTheme(),
		"minimal":   MinimalTheme(),
		"retro":     RetroTheme(),
		"nature":    NatureTheme(),
	}
}

// ThemeManager manages theme switching and persistence
type ThemeManager struct {
	currentTheme *Theme
	themes       map[string]*Theme
}

// NewThemeManager creates a new theme manager
func NewThemeManager() *ThemeManager {
	return &ThemeManager{
		currentTheme: DefaultSpaceTheme(),
		themes:       GetAllThemes(),
	}
}

// SetTheme switches to a different theme
func (tm *ThemeManager) SetTheme(name string) error {
	if theme, exists := tm.themes[name]; exists {
		tm.currentTheme = theme
		return nil
	}
	return fmt.Errorf("theme '%s' not found", name)
}

// GetCurrentTheme returns the current active theme
func (tm *ThemeManager) GetCurrentTheme() *Theme {
	return tm.currentTheme
}

// GetAvailableThemes returns a list of available theme names
func (tm *ThemeManager) GetAvailableThemes() []string {
	var names []string
	for name := range tm.themes {
		names = append(names, name)
	}
	return names
}

// Style helpers for common UI elements
func (t *Theme) HeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Text).
		Background(t.Colors.Primary).
		Bold(true).
		Padding(0, t.Spacing.Medium)
}

func (t *Theme) CardStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Colors.Primary).
		Padding(t.Spacing.Medium).
		Background(t.Colors.Surface)
}

func (t *Theme) ButtonStyle(active bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(0, t.Spacing.Medium).
		Margin(0, t.Spacing.Small)
		
	if active {
		return style.
			Foreground(t.Colors.Background).
			Background(t.Colors.Accent).
			Bold(true)
	}
	
	return style.
		Foreground(t.Colors.Text).
		Background(t.Colors.Surface)
}

func (t *Theme) ProgressBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Success)
}

func (t *Theme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Error).
		Bold(true)
}

func (t *Theme) SuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Success).
		Bold(true)
}

func (t *Theme) InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.Info)
}

func (t *Theme) MutedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Colors.TextMuted)
}
