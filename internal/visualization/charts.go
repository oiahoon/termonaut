package visualization

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/pkg/models"
)

// ChartRenderer creates ASCII charts and visualizations
type ChartRenderer struct {
	Width  int
	Height int
}

// NewChartRenderer creates a new chart renderer
func NewChartRenderer(width, height int) *ChartRenderer {
	return &ChartRenderer{
		Width:  width,
		Height: height,
	}
}

// BarChart represents a bar chart configuration
type BarChart struct {
	Title       string
	Data        map[string]float64
	MaxBarWidth int
	ShowValues  bool
	ShowPercent bool
	SortDesc    bool
}

// LineChart represents a line chart configuration
type LineChart struct {
	Title  string
	Data   []float64
	Labels []string
	YAxis  YAxisConfig
}

// YAxisConfig configures the Y-axis
type YAxisConfig struct {
	Min   float64
	Max   float64
	Steps int
}

// RenderBarChart creates an ASCII bar chart
func (cr *ChartRenderer) RenderBarChart(chart *BarChart) string {
	if len(chart.Data) == 0 {
		return fmt.Sprintf("📊 %s\n(No data available)\n", chart.Title)
	}

	// Sort data
	type item struct {
		label string
		value float64
	}

	var items []item
	for label, value := range chart.Data {
		items = append(items, item{label, value})
	}

	if chart.SortDesc {
		sort.Slice(items, func(i, j int) bool {
			return items[i].value > items[j].value
		})
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].value < items[j].value
		})
	}

	// Find max value for scaling
	maxValue := 0.0
	totalValue := 0.0
	for _, item := range items {
		if item.value > maxValue {
			maxValue = item.value
		}
		totalValue += item.value
	}

	if maxValue == 0 {
		return fmt.Sprintf("📊 %s\n(No data to display)\n", chart.Title)
	}

	// Build chart
	result := fmt.Sprintf("📊 %s\n", chart.Title)
	result += strings.Repeat("━", 50) + "\n\n"

	maxLabelWidth := 0
	for _, item := range items {
		if len(item.label) > maxLabelWidth {
			maxLabelWidth = len(item.label)
		}
	}

	for _, item := range items {
		// Calculate bar width
		barWidth := int((item.value / maxValue) * float64(chart.MaxBarWidth))
		if barWidth < 1 && item.value > 0 {
			barWidth = 1
		}

		// Create bar
		bar := strings.Repeat("█", barWidth)
		if barWidth < chart.MaxBarWidth {
			bar += strings.Repeat("░", chart.MaxBarWidth-barWidth)
		}

		// Format label
		label := fmt.Sprintf("%-*s", maxLabelWidth, item.label)

		// Add values if requested
		valueStr := ""
		if chart.ShowValues && chart.ShowPercent {
			percentage := (item.value / totalValue) * 100
			valueStr = fmt.Sprintf(" %.1f (%.1f%%)", item.value, percentage)
		} else if chart.ShowValues {
			valueStr = fmt.Sprintf(" %.1f", item.value)
		} else if chart.ShowPercent {
			percentage := (item.value / totalValue) * 100
			valueStr = fmt.Sprintf(" %.1f%%", percentage)
		}

		result += fmt.Sprintf("%s [%s]%s\n", label, bar, valueStr)
	}

	return result + "\n"
}

// RenderLineChart creates an ASCII line chart
func (cr *ChartRenderer) RenderLineChart(chart *LineChart) string {
	if len(chart.Data) == 0 {
		return fmt.Sprintf("📈 %s\n(No data available)\n", chart.Title)
	}

	result := fmt.Sprintf("📈 %s\n", chart.Title)
	result += strings.Repeat("━", 50) + "\n\n"

	// Find min/max values
	minVal := chart.Data[0]
	maxVal := chart.Data[0]
	for _, val := range chart.Data {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	// Use provided Y-axis config or auto-scale
	yMin := chart.YAxis.Min
	yMax := chart.YAxis.Max
	if yMin == 0 && yMax == 0 {
		yMin = minVal
		yMax = maxVal
	}

	height := cr.Height
	if height == 0 {
		height = 10
	}

	// Create the chart grid
	for row := height - 1; row >= 0; row-- {
		// Y-axis label
		yValue := yMin + (yMax-yMin)*float64(row)/float64(height-1)
		result += fmt.Sprintf("%6.1f │", yValue)

		// Plot points
		for i, val := range chart.Data {
			// Normalize value to chart height
			normalizedVal := (val - yMin) / (yMax - yMin)
			plotRow := int(normalizedVal * float64(height-1))

			if plotRow == row {
				result += "●"
			} else if plotRow == row-1 || plotRow == row+1 {
				result += "│"
			} else {
				result += " "
			}

			// Add connecting line
			if i < len(chart.Data)-1 {
				nextVal := chart.Data[i+1]
				nextNormalizedVal := (nextVal - yMin) / (yMax - yMin)
				nextPlotRow := int(nextNormalizedVal * float64(height-1))

				if (plotRow < row && nextPlotRow >= row) || (plotRow >= row && nextPlotRow < row) {
					result += "─"
				} else {
					result += " "
				}
			}
		}
		result += "\n"
	}

	// X-axis
	result += "       └"
	for i := 0; i < len(chart.Data)*2-1; i++ {
		result += "─"
	}
	result += "\n"

	// X-axis labels
	if len(chart.Labels) > 0 {
		result += "        "
		for i, label := range chart.Labels {
			if i < len(chart.Data) {
				if len(label) > 3 {
					label = label[:3]
				}
				result += fmt.Sprintf("%-2s", label)
			}
		}
	}

	return result + "\n\n"
}

// RenderCommandFrequencyChart creates a command frequency visualization
func (cr *ChartRenderer) RenderCommandFrequencyChart(commands []*models.Command, limit int) string {
	if len(commands) == 0 {
		return "📊 Command Frequency\n(No commands to analyze)\n"
	}

	// Count command frequencies
	cmdFreq := make(map[string]int)
	for _, cmd := range commands {
		// Get base command (first word)
		parts := strings.Fields(cmd.Command)
		if len(parts) > 0 {
			baseCmd := parts[0]
			cmdFreq[baseCmd]++
		}
	}

	// Convert to float64 for chart
	chartData := make(map[string]float64)
	type cmdCount struct {
		cmd   string
		count int
	}

	var sortedCmds []cmdCount
	for cmd, count := range cmdFreq {
		sortedCmds = append(sortedCmds, cmdCount{cmd, count})
	}

	sort.Slice(sortedCmds, func(i, j int) bool {
		return sortedCmds[i].count > sortedCmds[j].count
	})

	// Limit results
	displayLimit := limit
	if displayLimit == 0 || displayLimit > len(sortedCmds) {
		displayLimit = len(sortedCmds)
	}

	for i := 0; i < displayLimit; i++ {
		chartData[sortedCmds[i].cmd] = float64(sortedCmds[i].count)
	}

	chart := &BarChart{
		Title:       "Top Command Usage",
		Data:        chartData,
		MaxBarWidth: 30,
		ShowValues:  true,
		ShowPercent: true,
		SortDesc:    true,
	}

	return cr.RenderBarChart(chart)
}

// RenderCategoryDistribution creates a category distribution chart
func (cr *ChartRenderer) RenderCategoryDistribution(commands []*models.Command) string {
	if len(commands) == 0 {
		return "📊 Category Distribution\n(No commands to analyze)\n"
	}

	classifier := categories.NewCommandClassifier()
	categoryCount := make(map[string]float64)

	// Count commands by category
	for _, cmd := range commands {
		category := classifier.ClassifyCommand(cmd.Command)
		info := classifier.GetCategoryInfo(category)
		categoryCount[fmt.Sprintf("%s %s", info.Icon, info.Name)]++
	}

	chart := &BarChart{
		Title:       "Command Categories",
		Data:        categoryCount,
		MaxBarWidth: 25,
		ShowValues:  false,
		ShowPercent: true,
		SortDesc:    true,
	}

	return cr.RenderBarChart(chart)
}

// RenderActivityTimeline creates an activity timeline
func (cr *ChartRenderer) RenderActivityTimeline(commands []*models.Command, days int) string {
	if len(commands) == 0 {
		return "📈 Activity Timeline\n(No activity data)\n"
	}

	// Group commands by day
	dailyActivity := make(map[string]int)
	now := time.Now()

	// Initialize all days with 0
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("01/02")
		dailyActivity[dateStr] = 0
	}

	// Count actual activity
	for _, cmd := range commands {
		// Only include commands from the specified time range
		daysDiff := int(now.Sub(cmd.Timestamp).Hours() / 24)
		if daysDiff < days {
			dateStr := cmd.Timestamp.Format("01/02")
			dailyActivity[dateStr]++
		}
	}

	// Prepare data for line chart
	var data []float64
	var labels []string

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("01/02")
		data = append(data, float64(dailyActivity[dateStr]))
		labels = append(labels, dateStr)
	}

	chart := &LineChart{
		Title:  fmt.Sprintf("Activity Timeline (Last %d Days)", days),
		Data:   data,
		Labels: labels,
	}

	return cr.RenderLineChart(chart)
}

// RenderProgressMeter creates a visual progress meter
func (cr *ChartRenderer) RenderProgressMeter(current, target float64, title string) string {
	if target == 0 {
		return fmt.Sprintf("🎯 %s: No target set\n", title)
	}

	percentage := (current / target) * 100
	if percentage > 100 {
		percentage = 100
	}

	barWidth := 40
	filledWidth := int((percentage / 100) * float64(barWidth))
	emptyWidth := barWidth - filledWidth

	bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", emptyWidth)

	result := fmt.Sprintf("🎯 %s\n", title)
	result += fmt.Sprintf("[%s] %.1f%%\n", bar, percentage)
	result += fmt.Sprintf("Progress: %.0f / %.0f\n", current, target)

	return result
}

// RenderSummaryDashboard creates a comprehensive dashboard
func (cr *ChartRenderer) RenderSummaryDashboard(commands []*models.Command, sessions []*models.Session) string {
	dashboard := fmt.Sprintf("🎮 Termonaut Dashboard\n")
	dashboard += strings.Repeat("═", 60) + "\n\n"

	// Basic stats
	totalCommands := len(commands)
	totalSessions := len(sessions)

	// Calculate unique commands
	uniqueCommands := make(map[string]bool)
	for _, cmd := range commands {
		uniqueCommands[cmd.Command] = true
	}

	dashboard += fmt.Sprintf("📊 Quick Stats:\n")
	dashboard += fmt.Sprintf("   Commands: %d | Sessions: %d | Unique: %d\n\n",
		totalCommands, totalSessions, len(uniqueCommands))

	// Recent activity (last 7 days)
	if totalCommands > 0 {
		dashboard += cr.RenderActivityTimeline(commands, 7)
	}

	// Top commands
	if totalCommands > 5 {
		dashboard += cr.RenderCommandFrequencyChart(commands, 5)
	}

	// Category distribution
	if totalCommands > 0 {
		dashboard += cr.RenderCategoryDistribution(commands)
	}

	return dashboard
}

// CreateProgressVisualization creates rich progress visualization
func CreateProgressVisualization(current, max int, width int, style string) string {
	if max == 0 {
		return strings.Repeat("░", width)
	}

	percentage := float64(current) / float64(max)
	if percentage > 1.0 {
		percentage = 1.0
	}

	filled := int(percentage * float64(width))
	empty := width - filled

	switch style {
	case "blocks":
		return strings.Repeat("█", filled) + strings.Repeat("░", empty)
	case "dots":
		return strings.Repeat("●", filled) + strings.Repeat("○", empty)
	case "arrows":
		return strings.Repeat("▶", filled) + strings.Repeat("▷", empty)
	case "stars":
		return strings.Repeat("★", filled) + strings.Repeat("☆", empty)
	default:
		return strings.Repeat("█", filled) + strings.Repeat("░", empty)
	}
}

// FormatTrendIndicator creates trend visualization
func FormatTrendIndicator(current, previous float64) string {
	if previous == 0 {
		return "🆕 New"
	}

	change := ((current - previous) / previous) * 100

	if change > 10 {
		return fmt.Sprintf("📈 +%.1f%%", change)
	} else if change > 0 {
		return fmt.Sprintf("↗️ +%.1f%%", change)
	} else if change < -10 {
		return fmt.Sprintf("📉 %.1f%%", change)
	} else if change < 0 {
		return fmt.Sprintf("↘️ %.1f%%", change)
	} else {
		return "➡️ Stable"
	}
}
