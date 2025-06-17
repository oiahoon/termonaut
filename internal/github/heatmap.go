package github

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/stats"
)

// HeatmapGenerator generates GitHub-style activity heatmaps
type HeatmapGenerator struct {
	stats *stats.StatsCalculator
}

// NewHeatmapGenerator creates a new heatmap generator
func NewHeatmapGenerator(statsCalculator *stats.StatsCalculator) *HeatmapGenerator {
	return &HeatmapGenerator{
		stats: statsCalculator,
	}
}

// DayActivity represents activity data for a single day
type DayActivity struct {
	Date         time.Time `json:"date"`
	CommandCount int       `json:"command_count"`
	XPEarned     int       `json:"xp_earned"`
	Level        int       `json:"level"`
}

// HeatmapData represents a year's worth of activity data
type HeatmapData struct {
	Year       int           `json:"year"`
	StartDate  time.Time     `json:"start_date"`
	EndDate    time.Time     `json:"end_date"`
	Days       []DayActivity `json:"days"`
	TotalDays  int           `json:"total_days"`
	ActiveDays int           `json:"active_days"`
	MaxDaily   int           `json:"max_daily"`
	MinDaily   int           `json:"min_daily"`
	AvgDaily   float64       `json:"avg_daily"`
}

// GenerateYearHeatmap generates heatmap data for a specific year
func (hg *HeatmapGenerator) GenerateYearHeatmap(year int) (*HeatmapData, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	// For now, generate sample data since we don't have historical data yet
	// In a real implementation, this would query the database for historical stats
	days := hg.generateSampleData(startDate, endDate)

	// Calculate statistics
	totalCommands := 0
	activeDays := 0
	maxDaily := 0
	minDaily := -1

	for _, day := range days {
		totalCommands += day.CommandCount
		if day.CommandCount > 0 {
			activeDays++
		}
		if day.CommandCount > maxDaily {
			maxDaily = day.CommandCount
		}
		if minDaily == -1 || (day.CommandCount < minDaily && day.CommandCount > 0) {
			minDaily = day.CommandCount
		}
	}

	avgDaily := 0.0
	if activeDays > 0 {
		avgDaily = float64(totalCommands) / float64(activeDays)
	}

	if minDaily == -1 {
		minDaily = 0
	}

	return &HeatmapData{
		Year:       year,
		StartDate:  startDate,
		EndDate:    endDate,
		Days:       days,
		TotalDays:  len(days),
		ActiveDays: activeDays,
		MaxDaily:   maxDaily,
		MinDaily:   minDaily,
		AvgDaily:   avgDaily,
	}, nil
}

// generateSampleData creates sample activity data for demonstration
func (hg *HeatmapGenerator) generateSampleData(startDate, endDate time.Time) []DayActivity {
	var days []DayActivity

	current := startDate
	for current.Before(endDate) || current.Equal(endDate.Truncate(24*time.Hour)) {
		// Generate realistic activity patterns
		activity := hg.generateDayActivity(current)
		days = append(days, activity)
		current = current.AddDate(0, 0, 1)
	}

	return days
}

// generateDayActivity generates realistic activity for a single day
func (hg *HeatmapGenerator) generateDayActivity(date time.Time) DayActivity {
	// Create realistic patterns based on day of week and month
	weekday := date.Weekday()

	// Base activity levels (Monday = higher, Weekend = lower)
	baseActivity := 20
	switch weekday {
	case time.Saturday, time.Sunday:
		baseActivity = 5
	case time.Monday, time.Tuesday:
		baseActivity = 30
	case time.Wednesday, time.Thursday, time.Friday:
		baseActivity = 25
	}

	// Add some randomness and seasonal variation
	variation := (date.YearDay() % 7) - 3 // -3 to +3
	commands := baseActivity + variation

	// Some days have no activity
	if date.YearDay()%11 == 0 {
		commands = 0
	}

	// Spike days (high productivity)
	if date.YearDay()%23 == 0 {
		commands = baseActivity * 2
	}

	if commands < 0 {
		commands = 0
	}

	// Calculate XP based on commands (simplified)
	xp := commands * 10

	// Estimate level based on total XP (simplified)
	level := 1 + date.YearDay()/50

	return DayActivity{
		Date:         date,
		CommandCount: commands,
		XPEarned:     xp,
		Level:        level,
	}
}

// GenerateHTMLHeatmap creates an HTML representation of the heatmap
func (hg *HeatmapGenerator) GenerateHTMLHeatmap(data *HeatmapData) string {
	var html strings.Builder

	html.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Termonaut Activity Heatmap</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #0d1117;
            color: #c9d1d9;
        }
        .heatmap-container {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 6px;
            padding: 20px;
            margin: 20px 0;
        }
        .heatmap-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }
        .heatmap-title {
            font-size: 20px;
            font-weight: 600;
            color: #f0f6fc;
        }
        .heatmap-stats {
            font-size: 14px;
            color: #8b949e;
        }
        .heatmap-grid {
            display: grid;
            grid-template-columns: repeat(53, 1fr);
            gap: 3px;
            max-width: 800px;
        }
        .day-cell {
            width: 11px;
            height: 11px;
            border-radius: 2px;
            cursor: pointer;
        }
        .day-empty { background-color: #161b22; border: 1px solid #21262d; }
        .day-level-1 { background-color: #0e4429; }
        .day-level-2 { background-color: #006d32; }
        .day-level-3 { background-color: #26a641; }
        .day-level-4 { background-color: #39d353; }
        .legend {
            display: flex;
            align-items: center;
            gap: 5px;
            margin-top: 15px;
            font-size: 12px;
            color: #8b949e;
        }
        .legend-item {
            width: 11px;
            height: 11px;
            border-radius: 2px;
        }
        .month-labels {
            display: grid;
            grid-template-columns: repeat(12, 1fr);
            gap: 10px;
            margin-bottom: 10px;
            font-size: 12px;
            color: #8b949e;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }
        .stat-card {
            background: #21262d;
            border: 1px solid #30363d;
            border-radius: 6px;
            padding: 15px;
        }
        .stat-value {
            font-size: 24px;
            font-weight: 600;
            color: #f0f6fc;
        }
        .stat-label {
            font-size: 14px;
            color: #8b949e;
            margin-top: 5px;
        }
    </style>
</head>
 <body>
     <div class="heatmap-container">
         <div class="heatmap-header">
             <h1 class="heatmap-title">🚀 Terminal Activity Heatmap %d</h1>
             <div class="heatmap-stats">%d commands in %d active days</div>
         </div>`, data.Year, hg.sumCommands(data.Days), data.ActiveDays))

	// Month labels
	html.WriteString(`<div class="month-labels">`)
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	for _, month := range months {
		html.WriteString(fmt.Sprintf(`<div>%s</div>`, month))
	}
	html.WriteString(`</div>`)

	// Heatmap grid
	html.WriteString(`<div class="heatmap-grid">`)
	for _, day := range data.Days {
		level := hg.getActivityLevel(day.CommandCount, data.MaxDaily)
		class := fmt.Sprintf("day-cell day-level-%d", level)
		if day.CommandCount == 0 {
			class = "day-cell day-empty"
		}

		html.WriteString(fmt.Sprintf(
			`<div class="%s" title="%s: %d commands"></div>`,
			class,
			day.Date.Format("Jan 2, 2006"),
			day.CommandCount,
		))
	}
	html.WriteString(`</div>`)

	// Legend
	html.WriteString(`<div class="legend">
        <span>Less</span>
        <div class="legend-item day-empty"></div>
        <div class="legend-item day-level-1"></div>
        <div class="legend-item day-level-2"></div>
        <div class="legend-item day-level-3"></div>
        <div class="legend-item day-level-4"></div>
        <span>More</span>
    </div>`)

	// Statistics cards
	html.WriteString(`<div class="stats-grid">`)

	stats := []struct {
		value string
		label string
	}{
		{fmt.Sprintf("%d", hg.sumCommands(data.Days)), "Total Commands"},
		{fmt.Sprintf("%d", data.ActiveDays), "Active Days"},
		{fmt.Sprintf("%.1f", data.AvgDaily), "Average per Active Day"},
		{fmt.Sprintf("%d", data.MaxDaily), "Best Day"},
		{fmt.Sprintf("%.1f%%", float64(data.ActiveDays)/float64(data.TotalDays)*100), "Activity Rate"},
		{fmt.Sprintf("%d", hg.getCurrentStreak(data.Days)), "Current Streak"},
	}

	for _, stat := range stats {
		html.WriteString(fmt.Sprintf(`
        <div class="stat-card">
            <div class="stat-value">%s</div>
            <div class="stat-label">%s</div>
        </div>`, stat.value, stat.label))
	}

	html.WriteString(`</div>
    </div>
</body>
</html>`)

	return html.String()
}

// GenerateSVGHeatmap creates an SVG representation of the heatmap
func (hg *HeatmapGenerator) GenerateSVGHeatmap(data *HeatmapData) string {
	var svg strings.Builder

	width := 800
	height := 150
	cellSize := 11
	gap := 3

	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, width, height))
	svg.WriteString(`<style>
        .day-empty { fill: #ebedf0; }
        .day-level-1 { fill: #9be9a8; }
        .day-level-2 { fill: #40c463; }
        .day-level-3 { fill: #30a14e; }
        .day-level-4 { fill: #216e39; }
        .month-label { font: 12px sans-serif; fill: #586069; }
    </style>`)

	// Generate grid
	x, y := 0, 30
	week := 0

	for i, day := range data.Days {
		weekday := int(day.Date.Weekday())
		if weekday == 0 && i > 0 { // Sunday, start new week
			week++
			x = week * (cellSize + gap)
			y = 30
		} else {
			y = 30 + weekday*(cellSize+gap)
		}

		level := hg.getActivityLevel(day.CommandCount, data.MaxDaily)
		class := fmt.Sprintf("day-level-%d", level)
		if day.CommandCount == 0 {
			class = "day-empty"
		}

		svg.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" class="%s">`,
			x, y, cellSize, cellSize, class,
		))
		svg.WriteString(fmt.Sprintf(
			`<title>%s: %d commands</title></rect>`,
			day.Date.Format("Jan 2, 2006"), day.CommandCount,
		))
	}

	// Add month labels
	monthX := 0
	for month := 1; month <= 12; month++ {
		monthName := time.Month(month).String()[:3]
		svg.WriteString(fmt.Sprintf(
			`<text x="%d" y="20" class="month-label">%s</text>`,
			monthX, monthName,
		))
		monthX += 65 // Approximate spacing
	}

	svg.WriteString(`</svg>`)
	return svg.String()
}

// GenerateMarkdownHeatmap creates a markdown representation
func (hg *HeatmapGenerator) GenerateMarkdownHeatmap(data *HeatmapData) string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# 🔥 Terminal Activity Heatmap %d\n\n", data.Year))

	// Statistics table
	md.WriteString("## 📊 Statistics\n\n")
	md.WriteString("| Metric | Value |\n")
	md.WriteString("|--------|-------|\n")
	md.WriteString(fmt.Sprintf("| Total Commands | %d |\n", hg.sumCommands(data.Days)))
	md.WriteString(fmt.Sprintf("| Active Days | %d |\n", data.ActiveDays))
	md.WriteString(fmt.Sprintf("| Average per Day | %.1f |\n", data.AvgDaily))
	md.WriteString(fmt.Sprintf("| Best Day | %d commands |\n", data.MaxDaily))
	md.WriteString(fmt.Sprintf("| Activity Rate | %.1f%% |\n", float64(data.ActiveDays)/float64(data.TotalDays)*100))
	md.WriteString(fmt.Sprintf("| Current Streak | %d days |\n", hg.getCurrentStreak(data.Days)))

	// Monthly breakdown
	md.WriteString("\n## 📅 Monthly Breakdown\n\n")
	monthlyStats := hg.getMonthlyStats(data.Days)

	md.WriteString("| Month | Commands | Active Days | Best Day |\n")
	md.WriteString("|-------|----------|-------------|----------|\n")

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	for i, month := range months {
		if stats, ok := monthlyStats[i+1]; ok {
			md.WriteString(fmt.Sprintf("| %s | %d | %d | %d |\n",
				month, stats.TotalCommands, stats.ActiveDays, stats.BestDay))
		} else {
			md.WriteString(fmt.Sprintf("| %s | 0 | 0 | 0 |\n", month))
		}
	}

	// ASCII heatmap representation
	md.WriteString("\n## 🎨 Activity Heatmap\n\n")
	md.WriteString("```\n")
	md.WriteString(hg.generateASCIIHeatmap(data))
	md.WriteString("```\n\n")

	md.WriteString("Legend: `░` None, `▓` Low, `█` High\n\n")

	md.WriteString("*Generated by [Termonaut](https://github.com/oiahoon/termonaut) - Terminal productivity tracker*\n")

	return md.String()
}

// Helper functions

func (hg *HeatmapGenerator) getActivityLevel(commands, maxDaily int) int {
	if commands == 0 {
		return 0
	}

	if maxDaily == 0 {
		return 1
	}

	ratio := float64(commands) / float64(maxDaily)
	switch {
	case ratio >= 0.75:
		return 4
	case ratio >= 0.5:
		return 3
	case ratio >= 0.25:
		return 2
	default:
		return 1
	}
}

func (hg *HeatmapGenerator) sumCommands(days []DayActivity) int {
	total := 0
	for _, day := range days {
		total += day.CommandCount
	}
	return total
}

func (hg *HeatmapGenerator) getCurrentStreak(days []DayActivity) int {
	// Sort days by date descending
	sorted := make([]DayActivity, len(days))
	copy(sorted, days)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	streak := 0
	for _, day := range sorted {
		if day.CommandCount > 0 {
			streak++
		} else {
			break
		}
	}
	return streak
}

type MonthlyStats struct {
	TotalCommands int
	ActiveDays    int
	BestDay       int
}

func (hg *HeatmapGenerator) getMonthlyStats(days []DayActivity) map[int]MonthlyStats {
	monthlyStats := make(map[int]MonthlyStats)

	for _, day := range days {
		month := int(day.Date.Month())
		stats := monthlyStats[month]

		stats.TotalCommands += day.CommandCount
		if day.CommandCount > 0 {
			stats.ActiveDays++
		}
		if day.CommandCount > stats.BestDay {
			stats.BestDay = day.CommandCount
		}

		monthlyStats[month] = stats
	}

	return monthlyStats
}

func (hg *HeatmapGenerator) generateASCIIHeatmap(data *HeatmapData) string {
	var ascii strings.Builder

	// Create a simplified weekly view
	weeks := make([][]DayActivity, 0)
	currentWeek := make([]DayActivity, 0)

	for _, day := range data.Days {
		currentWeek = append(currentWeek, day)
		if day.Date.Weekday() == time.Saturday || len(currentWeek) == 7 {
			weeks = append(weeks, currentWeek)
			currentWeek = make([]DayActivity, 0)
		}
	}

	// Add remaining days
	if len(currentWeek) > 0 {
		weeks = append(weeks, currentWeek)
	}

	// Generate ASCII representation
	for _, week := range weeks {
		for _, day := range week {
			level := hg.getActivityLevel(day.CommandCount, data.MaxDaily)
			switch level {
			case 0:
				ascii.WriteString("░")
			case 1, 2:
				ascii.WriteString("▓")
			case 3, 4:
				ascii.WriteString("█")
			}
		}
		ascii.WriteString("\n")
	}

	return ascii.String()
}
