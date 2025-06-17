package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/oiahoon/termonaut/pkg/models"
)

// HeatmapAnalyzer generates productivity heatmaps and time insights
type HeatmapAnalyzer struct{}

// NewHeatmapAnalyzer creates a new heatmap analyzer
func NewHeatmapAnalyzer() *HeatmapAnalyzer {
	return &HeatmapAnalyzer{}
}

// HeatmapData represents heatmap visualization data
type HeatmapData struct {
	WeeklyHeatmap    map[time.Weekday]map[int]float64 `json:"weekly_heatmap"`
	MonthlyHeatmap   map[int]map[int]float64          `json:"monthly_heatmap"` // day -> hour -> intensity
	OptimalHours     []int                            `json:"optimal_hours"`
	PeakDays         []time.Weekday                   `json:"peak_days"`
	WorkingPatterns  []WorkingPattern                 `json:"working_patterns"`
	FocusScore       float64                          `json:"focus_score"`
	DistributionType string                           `json:"distribution_type"`
}

// WorkingPattern represents identified working patterns
type WorkingPattern struct {
	Name        string         `json:"name"`
	StartTime   int            `json:"start_time"` // hour
	EndTime     int            `json:"end_time"`   // hour
	Intensity   float64        `json:"intensity"`  // 0-100
	Days        []time.Weekday `json:"days"`
	Confidence  float64        `json:"confidence"` // 0-100
	Description string         `json:"description"`
}

// TimeInsights provides detailed time-based insights
type TimeInsights struct {
	Recommendations  []string             `json:"recommendations"`
	OptimalSchedule  map[string]string    `json:"optimal_schedule"`
	ProductivityTips []ProductivityTip    `json:"productivity_tips"`
	SeasonalTrends   SeasonalTrends       `json:"seasonal_trends"`
	WorkLifeBalance  WorkLifeBalanceScore `json:"work_life_balance"`
}

// ProductivityTip provides actionable productivity advice
type ProductivityTip struct {
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`   // high, medium, low
	Effort      string `json:"effort"`   // high, medium, low
	Priority    int    `json:"priority"` // 1-10
}

// SeasonalTrends analyzes productivity trends over time
type SeasonalTrends struct {
	WeekdayTrend   string  `json:"weekday_trend"`   // increasing, decreasing, stable
	WeekendPattern string  `json:"weekend_pattern"` // active, inactive, balanced
	MonthlyTrend   string  `json:"monthly_trend"`
	Seasonality    float64 `json:"seasonality"` // 0-100
}

// WorkLifeBalanceScore evaluates work-life balance
type WorkLifeBalanceScore struct {
	Score           float64  `json:"score"`            // 0-100
	WeekendActivity float64  `json:"weekend_activity"` // 0-100
	EveningActivity float64  `json:"evening_activity"` // 0-100
	OffHoursRatio   float64  `json:"off_hours_ratio"`  // 0-100
	BalanceLevel    string   `json:"balance_level"`    // excellent, good, needs_improvement, poor
	Recommendations []string `json:"recommendations"`
}

// GenerateHeatmap creates productivity heatmap data
func (ha *HeatmapAnalyzer) GenerateHeatmap(commands []*models.Command) *HeatmapData {
	if len(commands) == 0 {
		return &HeatmapData{}
	}

	// Initialize heatmap grids
	weeklyHeatmap := make(map[time.Weekday]map[int]float64)
	monthlyHeatmap := make(map[int]map[int]float64)

	// Initialize all days and hours
	for d := time.Sunday; d <= time.Saturday; d++ {
		weeklyHeatmap[d] = make(map[int]float64)
		for h := 0; h < 24; h++ {
			weeklyHeatmap[d][h] = 0
		}
	}

	// Process commands
	hourlyActivity := make(map[int]int)
	dayActivity := make(map[time.Weekday]int)
	dailyHourlyActivity := make(map[int]map[int]int) // day of year -> hour -> count

	for _, cmd := range commands {
		hour := cmd.Timestamp.Hour()
		day := cmd.Timestamp.Weekday()
		dayOfYear := cmd.Timestamp.YearDay()

		// Count activities
		hourlyActivity[hour]++
		dayActivity[day]++
		weeklyHeatmap[day][hour]++

		if dailyHourlyActivity[dayOfYear] == nil {
			dailyHourlyActivity[dayOfYear] = make(map[int]int)
		}
		dailyHourlyActivity[dayOfYear][hour]++
	}

	// Normalize weekly heatmap (0-100 scale)
	maxActivity := 0.0
	for d := time.Sunday; d <= time.Saturday; d++ {
		for h := 0; h < 24; h++ {
			if weeklyHeatmap[d][h] > maxActivity {
				maxActivity = weeklyHeatmap[d][h]
			}
		}
	}

	if maxActivity > 0 {
		for d := time.Sunday; d <= time.Saturday; d++ {
			for h := 0; h < 24; h++ {
				weeklyHeatmap[d][h] = (weeklyHeatmap[d][h] / maxActivity) * 100
			}
		}
	}

	// Generate monthly heatmap
	for dayOfYear, hours := range dailyHourlyActivity {
		monthlyHeatmap[dayOfYear] = make(map[int]float64)
		maxDayActivity := 0
		for _, count := range hours {
			if count > maxDayActivity {
				maxDayActivity = count
			}
		}

		for hour, count := range hours {
			if maxDayActivity > 0 {
				monthlyHeatmap[dayOfYear][hour] = float64(count) / float64(maxDayActivity) * 100
			}
		}
	}

	// Find optimal hours (top 25% most active)
	type hourActivity struct {
		hour     int
		activity int
	}

	var sortedHours []hourActivity
	for hour, activity := range hourlyActivity {
		sortedHours = append(sortedHours, hourActivity{hour, activity})
	}

	sort.Slice(sortedHours, func(i, j int) bool {
		return sortedHours[i].activity > sortedHours[j].activity
	})

	optimalCount := int(math.Max(float64(len(sortedHours))/4, 1))
	optimalHours := make([]int, optimalCount)
	for i := 0; i < optimalCount && i < len(sortedHours); i++ {
		optimalHours[i] = sortedHours[i].hour
	}

	// Find peak days
	type dayActivityPair struct {
		day      time.Weekday
		activity int
	}

	var sortedDays []dayActivityPair
	for day, activity := range dayActivity {
		sortedDays = append(sortedDays, dayActivityPair{day, activity})
	}

	sort.Slice(sortedDays, func(i, j int) bool {
		return sortedDays[i].activity > sortedDays[j].activity
	})

	peakDayCount := int(math.Max(float64(len(sortedDays))/2, 1))
	peakDays := make([]time.Weekday, peakDayCount)
	for i := 0; i < peakDayCount && i < len(sortedDays); i++ {
		peakDays[i] = sortedDays[i].day
	}

	// Identify working patterns
	workingPatterns := ha.identifyWorkingPatterns(weeklyHeatmap)

	// Calculate focus score (consistency of activity)
	focusScore := ha.calculateFocusScore(weeklyHeatmap)

	// Determine distribution type
	distributionType := ha.determineDistributionType(weeklyHeatmap)

	return &HeatmapData{
		WeeklyHeatmap:    weeklyHeatmap,
		MonthlyHeatmap:   monthlyHeatmap,
		OptimalHours:     optimalHours,
		PeakDays:         peakDays,
		WorkingPatterns:  workingPatterns,
		FocusScore:       focusScore,
		DistributionType: distributionType,
	}
}

// identifyWorkingPatterns detects common working patterns
func (ha *HeatmapAnalyzer) identifyWorkingPatterns(heatmap map[time.Weekday]map[int]float64) []WorkingPattern {
	patterns := []WorkingPattern{}

	// Morning Worker Pattern (6-12)
	morningScore := ha.calculateTimeRangeScore(heatmap, 6, 12)
	if morningScore > 30 {
		patterns = append(patterns, WorkingPattern{
			Name:        "Early Bird",
			StartTime:   6,
			EndTime:     12,
			Intensity:   morningScore,
			Days:        []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			Confidence:  ha.calculatePatternConfidence(heatmap, 6, 12),
			Description: "You're most productive during morning hours. Peak creativity and focus time.",
		})
	}

	// Standard Work Hours (9-17)
	standardScore := ha.calculateTimeRangeScore(heatmap, 9, 17)
	if standardScore > 25 {
		patterns = append(patterns, WorkingPattern{
			Name:        "Standard Hours",
			StartTime:   9,
			EndTime:     17,
			Intensity:   standardScore,
			Days:        []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			Confidence:  ha.calculatePatternConfidence(heatmap, 9, 17),
			Description: "Traditional 9-to-5 working pattern. Good work-life balance structure.",
		})
	}

	// Night Owl Pattern (18-24)
	eveningScore := ha.calculateTimeRangeScore(heatmap, 18, 24)
	if eveningScore > 25 {
		patterns = append(patterns, WorkingPattern{
			Name:        "Night Owl",
			StartTime:   18,
			EndTime:     24,
			Intensity:   eveningScore,
			Days:        []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
			Confidence:  ha.calculatePatternConfidence(heatmap, 18, 24),
			Description: "You work best during evening hours. Consider optimizing your schedule.",
		})
	}

	// Weekend Warrior
	weekendScore := ha.calculateWeekendScore(heatmap)
	if weekendScore > 20 {
		patterns = append(patterns, WorkingPattern{
			Name:        "Weekend Warrior",
			StartTime:   10,
			EndTime:     22,
			Intensity:   weekendScore,
			Days:        []time.Weekday{time.Saturday, time.Sunday},
			Confidence:  weekendScore,
			Description: "Significant weekend activity. Consider work-life balance optimization.",
		})
	}

	return patterns
}

// calculateTimeRangeScore calculates average activity score for a time range
func (ha *HeatmapAnalyzer) calculateTimeRangeScore(heatmap map[time.Weekday]map[int]float64, startHour, endHour int) float64 {
	total := 0.0
	count := 0

	for d := time.Monday; d <= time.Friday; d++ {
		for h := startHour; h < endHour; h++ {
			total += heatmap[d][h]
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// calculateWeekendScore calculates weekend activity score
func (ha *HeatmapAnalyzer) calculateWeekendScore(heatmap map[time.Weekday]map[int]float64) float64 {
	total := 0.0
	count := 0

	for _, day := range []time.Weekday{time.Saturday, time.Sunday} {
		for h := 0; h < 24; h++ {
			total += heatmap[day][h]
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// calculatePatternConfidence calculates confidence level for a pattern
func (ha *HeatmapAnalyzer) calculatePatternConfidence(heatmap map[time.Weekday]map[int]float64, startHour, endHour int) float64 {
	scores := []float64{}

	for d := time.Monday; d <= time.Friday; d++ {
		dayScore := 0.0
		count := 0
		for h := startHour; h < endHour; h++ {
			dayScore += heatmap[d][h]
			count++
		}
		if count > 0 {
			scores = append(scores, dayScore/float64(count))
		}
	}

	if len(scores) == 0 {
		return 0
	}

	// Calculate coefficient of variation (lower is more consistent)
	mean := 0.0
	for _, score := range scores {
		mean += score
	}
	mean /= float64(len(scores))

	variance := 0.0
	for _, score := range scores {
		variance += math.Pow(score-mean, 2)
	}
	variance /= float64(len(scores))

	stdDev := math.Sqrt(variance)
	cv := stdDev / mean

	// Convert to confidence (0-100, higher is more confident)
	confidence := math.Max(0, 100-cv*100)
	return math.Min(confidence, 100)
}

// calculateFocusScore calculates how focused/consistent the activity is
func (ha *HeatmapAnalyzer) calculateFocusScore(heatmap map[time.Weekday]map[int]float64) float64 {
	// Calculate entropy-like measure
	total := 0.0
	nonZeroSlots := 0

	for d := time.Sunday; d <= time.Saturday; d++ {
		for h := 0; h < 24; h++ {
			if heatmap[d][h] > 0 {
				total += heatmap[d][h]
				nonZeroSlots++
			}
		}
	}

	if nonZeroSlots == 0 {
		return 0
	}

	// More concentrated activity = higher focus score
	maxPossibleSlots := 7 * 24 // 7 days * 24 hours
	concentration := float64(maxPossibleSlots-nonZeroSlots) / float64(maxPossibleSlots) * 100

	// Adjust based on intensity distribution
	entropy := 0.0
	for d := time.Sunday; d <= time.Saturday; d++ {
		for h := 0; h < 24; h++ {
			if heatmap[d][h] > 0 {
				p := heatmap[d][h] / total
				entropy -= p * math.Log2(p)
			}
		}
	}

	// Normalize entropy and combine with concentration
	maxEntropy := math.Log2(float64(nonZeroSlots))
	normalizedEntropy := 0.0
	if maxEntropy > 0 {
		normalizedEntropy = entropy / maxEntropy
	}

	focusScore := (concentration * 0.6) + ((1 - normalizedEntropy) * 40)
	return math.Max(0, math.Min(focusScore, 100))
}

// determineDistributionType categorizes the activity distribution
func (ha *HeatmapAnalyzer) determineDistributionType(heatmap map[time.Weekday]map[int]float64) string {
	morningScore := ha.calculateTimeRangeScore(heatmap, 6, 12)
	afternoonScore := ha.calculateTimeRangeScore(heatmap, 12, 18)
	eveningScore := ha.calculateTimeRangeScore(heatmap, 18, 24)

	maxScore := math.Max(morningScore, math.Max(afternoonScore, eveningScore))

	if maxScore < 10 {
		return "Scattered"
	} else if morningScore == maxScore && morningScore > afternoonScore*1.5 {
		return "Morning-Focused"
	} else if eveningScore == maxScore && eveningScore > afternoonScore*1.5 {
		return "Evening-Focused"
	} else if afternoonScore == maxScore {
		return "Traditional"
	} else {
		return "Balanced"
	}
}

// GenerateTimeInsights creates actionable time management insights
func (ha *HeatmapAnalyzer) GenerateTimeInsights(heatmapData *HeatmapData, commands []*models.Command) *TimeInsights {
	insights := &TimeInsights{
		Recommendations:  []string{},
		OptimalSchedule:  make(map[string]string),
		ProductivityTips: []ProductivityTip{},
	}

	// Generate recommendations based on patterns
	if heatmapData.FocusScore < 50 {
		insights.Recommendations = append(insights.Recommendations,
			"Consider consolidating your work into specific time blocks for better focus")
	}

	// Generate optimal schedule
	if len(heatmapData.OptimalHours) > 0 {
		optimalStart := heatmapData.OptimalHours[0]
		for _, hour := range heatmapData.OptimalHours {
			if hour < optimalStart {
				optimalStart = hour
			}
		}
		insights.OptimalSchedule["peak_start"] = fmt.Sprintf("%02d:00", optimalStart)
		insights.OptimalSchedule["peak_duration"] = fmt.Sprintf("%d hours", len(heatmapData.OptimalHours))
	}

	// Generate productivity tips
	insights.ProductivityTips = append(insights.ProductivityTips, ProductivityTip{
		Category:    "Time Management",
		Title:       "Optimize Your Peak Hours",
		Description: "Schedule your most important tasks during your identified peak productivity hours",
		Impact:      "high",
		Effort:      "low",
		Priority:    8,
	})

	if heatmapData.FocusScore > 70 {
		insights.ProductivityTips = append(insights.ProductivityTips, ProductivityTip{
			Category:    "Focus",
			Title:       "Maintain Your Focus Pattern",
			Description: "Your current time management shows excellent focus. Keep it up!",
			Impact:      "medium",
			Effort:      "low",
			Priority:    5,
		})
	}

	return insights
}

// FormatHeatmapVisualization creates ASCII heatmap visualization
func (ha *HeatmapAnalyzer) FormatHeatmapVisualization(heatmapData *HeatmapData) string {
	if len(heatmapData.WeeklyHeatmap) == 0 {
		return "📊 No heatmap data available yet."
	}

	result := fmt.Sprintf("🔥 Weekly Productivity Heatmap\n")
	result += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Hours header
	result += fmt.Sprintf("       ")
	for h := 0; h < 24; h += 3 {
		result += fmt.Sprintf("%02d  ", h)
	}
	result += fmt.Sprintf("\n")

	// Days and heatmap
	days := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	dayNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	for i, day := range days {
		result += fmt.Sprintf("%s  ", dayNames[i])

		for h := 0; h < 24; h += 3 {
			intensity := heatmapData.WeeklyHeatmap[day][h]
			char := ha.getHeatmapChar(intensity)
			result += fmt.Sprintf(" %s  ", char)
		}
		result += fmt.Sprintf("\n")
	}

	result += fmt.Sprintf("\n")
	result += fmt.Sprintf("Legend: ░ Light  ▓ Medium  █ High  🔥 Peak\n")
	result += fmt.Sprintf("Focus Score: %.1f/100 (%s)\n", heatmapData.FocusScore, ha.getFocusDescription(heatmapData.FocusScore))
	result += fmt.Sprintf("Distribution: %s\n", heatmapData.DistributionType)

	// Optimal hours
	if len(heatmapData.OptimalHours) > 0 {
		result += fmt.Sprintf("\n🎯 Your Peak Hours: ")
		sort.Ints(heatmapData.OptimalHours)
		for i, hour := range heatmapData.OptimalHours {
			if i > 0 {
				result += fmt.Sprintf(", ")
			}
			result += fmt.Sprintf("%02d:00", hour)
		}
		result += fmt.Sprintf("\n")
	}

	// Working patterns
	if len(heatmapData.WorkingPatterns) > 0 {
		result += fmt.Sprintf("\n💼 Identified Patterns:\n")
		for _, pattern := range heatmapData.WorkingPatterns {
			result += fmt.Sprintf("  %s: %02d:00-%02d:00 (%.1f%% intensity, %.1f%% confidence)\n",
				pattern.Name, pattern.StartTime, pattern.EndTime, pattern.Intensity, pattern.Confidence)
		}
	}

	return result
}

// getHeatmapChar returns the appropriate character for intensity level
func (ha *HeatmapAnalyzer) getHeatmapChar(intensity float64) string {
	if intensity > 80 {
		return "🔥"
	} else if intensity > 60 {
		return "█"
	} else if intensity > 30 {
		return "▓"
	} else if intensity > 10 {
		return "░"
	}
	return " "
}

// getFocusDescription returns description for focus score
func (ha *HeatmapAnalyzer) getFocusDescription(score float64) string {
	if score > 80 {
		return "Highly Focused"
	} else if score > 60 {
		return "Well Focused"
	} else if score > 40 {
		return "Moderately Focused"
	} else if score > 20 {
		return "Scattered"
	}
	return "Very Scattered"
}
