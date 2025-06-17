package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/pkg/models"
)

// ProductivityAnalyzer analyzes user productivity patterns
type ProductivityAnalyzer struct {
	classifier *categories.CommandClassifier
}

// NewProductivityAnalyzer creates a new productivity analyzer
func NewProductivityAnalyzer() *ProductivityAnalyzer {
	return &ProductivityAnalyzer{
		classifier: categories.NewCommandClassifier(),
	}
}

// ProductivityMetrics holds comprehensive productivity analysis
type ProductivityMetrics struct {
	OverallScore      float64                   `json:"overall_score"`
	DailyPattern      DailyProductivityPattern  `json:"daily_pattern"`
	WeeklyPattern     WeeklyProductivityPattern `json:"weekly_pattern"`
	CategoryInsights  CategoryProductivityStats `json:"category_insights"`
	StreakAnalysis    StreakAnalysis            `json:"streak_analysis"`
	EfficiencyMetrics EfficiencyMetrics         `json:"efficiency_metrics"`
	TimeDistribution  map[string]float64        `json:"time_distribution"`
}

// DailyProductivityPattern represents hourly productivity
type DailyProductivityPattern struct {
	PeakHours      []int             `json:"peak_hours"`
	PeakScore      float64           `json:"peak_score"`
	LowHours       []int             `json:"low_hours"`
	HourlyStats    map[int]*HourStat `json:"hourly_stats"`
	MorningScore   float64           `json:"morning_score"`   // 6-12
	AfternoonScore float64           `json:"afternoon_score"` // 12-18
	EveningScore   float64           `json:"evening_score"`   // 18-24
	NightScore     float64           `json:"night_score"`     // 0-6
}

// HourStat represents statistics for a specific hour
type HourStat struct {
	Hour         int     `json:"hour"`
	CommandCount int     `json:"command_count"`
	Productivity float64 `json:"productivity"`
	Diversity    float64 `json:"diversity"`
}

// WeeklyProductivityPattern represents daily productivity across the week
type WeeklyProductivityPattern struct {
	MostProductiveDay  time.Weekday              `json:"most_productive_day"`
	LeastProductiveDay time.Weekday              `json:"least_productive_day"`
	WeekdayScore       float64                   `json:"weekday_score"`
	WeekendScore       float64                   `json:"weekend_score"`
	DailyStats         map[time.Weekday]*DayStat `json:"daily_stats"`
}

// DayStat represents statistics for a specific day of the week
type DayStat struct {
	Day            time.Weekday `json:"day"`
	CommandCount   int          `json:"command_count"`
	Productivity   float64      `json:"productivity"`
	CategorySpread int          `json:"category_spread"`
}

// CategoryProductivityStats analyzes productivity by command category
type CategoryProductivityStats struct {
	TopCategory         categories.Category             `json:"top_category"`
	MostImproved        categories.Category             `json:"most_improved"`
	CategoryScores      map[categories.Category]float64 `json:"category_scores"`
	DiversityScore      float64                         `json:"diversity_score"`
	SpecializationLevel string                          `json:"specialization_level"`
}

// StreakAnalysis analyzes usage streaks and consistency
type StreakAnalysis struct {
	CurrentStreak    int     `json:"current_streak"`
	LongestStreak    int     `json:"longest_streak"`
	ConsistencyScore float64 `json:"consistency_score"`
	DaysActive       int     `json:"days_active"`
	WeeksActive      int     `json:"weeks_active"`
	AverageGap       float64 `json:"average_gap"`
}

// EfficiencyMetrics measures command efficiency and patterns
type EfficiencyMetrics struct {
	CommandsPerSession   float64 `json:"commands_per_session"`
	AverageSessionLength float64 `json:"average_session_length"`
	RepetitiveCommands   float64 `json:"repetitive_commands"`
	UniqueCommandRatio   float64 `json:"unique_command_ratio"`
	ComplexityScore      float64 `json:"complexity_score"`
	AutomationPotential  float64 `json:"automation_potential"`
}

// AnalyzeProductivity performs comprehensive productivity analysis
func (pa *ProductivityAnalyzer) AnalyzeProductivity(commands []*models.Command, sessions []*models.Session) *ProductivityMetrics {
	if len(commands) == 0 {
		return &ProductivityMetrics{}
	}

	metrics := &ProductivityMetrics{
		DailyPattern:      pa.analyzeDailyPattern(commands),
		WeeklyPattern:     pa.analyzeWeeklyPattern(commands),
		CategoryInsights:  pa.analyzeCategoryProductivity(commands),
		StreakAnalysis:    pa.analyzeStreaks(commands),
		EfficiencyMetrics: pa.analyzeEfficiency(commands, sessions),
		TimeDistribution:  pa.analyzeTimeDistribution(commands),
	}

	// Calculate overall productivity score
	metrics.OverallScore = pa.calculateOverallScore(metrics)

	return metrics
}

// analyzeDailyPattern analyzes hourly productivity patterns
func (pa *ProductivityAnalyzer) analyzeDailyPattern(commands []*models.Command) DailyProductivityPattern {
	hourlyStats := make(map[int]*HourStat)
	hourlyCommands := make(map[int][]*models.Command)

	// Group commands by hour
	for _, cmd := range commands {
		hour := cmd.Timestamp.Hour()
		if _, exists := hourlyStats[hour]; !exists {
			hourlyStats[hour] = &HourStat{Hour: hour}
		}
		hourlyStats[hour].CommandCount++
		hourlyCommands[hour] = append(hourlyCommands[hour], cmd)
	}

	// Calculate productivity and diversity for each hour
	var peakHours []int
	var lowHours []int
	maxProductivity := 0.0
	minProductivity := math.Inf(1)

	morningCommands, afternoonCommands, eveningCommands, nightCommands := 0, 0, 0, 0

	for hour, stat := range hourlyStats {
		// Calculate diversity (unique commands / total commands)
		uniqueCommands := make(map[string]bool)
		for _, cmd := range hourlyCommands[hour] {
			uniqueCommands[cmd.Command] = true
		}
		stat.Diversity = float64(len(uniqueCommands)) / float64(stat.CommandCount)

		// Calculate productivity score (commands * diversity)
		stat.Productivity = float64(stat.CommandCount) * stat.Diversity

		// Track peak and low hours
		if stat.Productivity > maxProductivity {
			maxProductivity = stat.Productivity
			peakHours = []int{hour}
		} else if stat.Productivity == maxProductivity {
			peakHours = append(peakHours, hour)
		}

		if stat.Productivity < minProductivity {
			minProductivity = stat.Productivity
			lowHours = []int{hour}
		} else if stat.Productivity == minProductivity {
			lowHours = append(lowHours, hour)
		}

		// Count commands for time periods
		switch {
		case hour >= 6 && hour < 12:
			morningCommands += stat.CommandCount
		case hour >= 12 && hour < 18:
			afternoonCommands += stat.CommandCount
		case hour >= 18 && hour < 24:
			eveningCommands += stat.CommandCount
		default:
			nightCommands += stat.CommandCount
		}
	}

	totalCommands := float64(len(commands))

	return DailyProductivityPattern{
		PeakHours:      peakHours,
		PeakScore:      maxProductivity,
		LowHours:       lowHours,
		HourlyStats:    hourlyStats,
		MorningScore:   float64(morningCommands) / totalCommands * 100,
		AfternoonScore: float64(afternoonCommands) / totalCommands * 100,
		EveningScore:   float64(eveningCommands) / totalCommands * 100,
		NightScore:     float64(nightCommands) / totalCommands * 100,
	}
}

// analyzeWeeklyPattern analyzes daily productivity patterns across the week
func (pa *ProductivityAnalyzer) analyzeWeeklyPattern(commands []*models.Command) WeeklyProductivityPattern {
	dailyStats := make(map[time.Weekday]*DayStat)
	dailyCommands := make(map[time.Weekday][]*models.Command)

	// Initialize all weekdays
	for d := time.Sunday; d <= time.Saturday; d++ {
		dailyStats[d] = &DayStat{Day: d}
	}

	// Group commands by weekday
	for _, cmd := range commands {
		day := cmd.Timestamp.Weekday()
		dailyStats[day].CommandCount++
		dailyCommands[day] = append(dailyCommands[day], cmd)
	}

	// Calculate productivity for each day
	mostProductiveDay := time.Sunday
	leastProductiveDay := time.Sunday
	maxProductivity := 0.0
	minProductivity := math.Inf(1)

	weekdayCommands, weekendCommands := 0, 0

	for day, stat := range dailyStats {
		if stat.CommandCount == 0 {
			continue
		}

		// Calculate category spread
		categories := make(map[categories.Category]bool)
		for _, cmd := range dailyCommands[day] {
			category := pa.classifier.ClassifyCommand(cmd.Command)
			categories[category] = true
		}
		stat.CategorySpread = len(categories)

		// Calculate productivity (commands * category diversity)
		stat.Productivity = float64(stat.CommandCount) * float64(stat.CategorySpread)

		// Track most/least productive days
		if stat.Productivity > maxProductivity {
			maxProductivity = stat.Productivity
			mostProductiveDay = day
		}
		if stat.Productivity < minProductivity && stat.CommandCount > 0 {
			minProductivity = stat.Productivity
			leastProductiveDay = day
		}

		// Count weekday vs weekend commands
		if day == time.Saturday || day == time.Sunday {
			weekendCommands += stat.CommandCount
		} else {
			weekdayCommands += stat.CommandCount
		}
	}

	totalCommands := float64(len(commands))
	weekdayScore := float64(weekdayCommands) / totalCommands * 100
	weekendScore := float64(weekendCommands) / totalCommands * 100

	return WeeklyProductivityPattern{
		MostProductiveDay:  mostProductiveDay,
		LeastProductiveDay: leastProductiveDay,
		WeekdayScore:       weekdayScore,
		WeekendScore:       weekendScore,
		DailyStats:         dailyStats,
	}
}

// analyzeCategoryProductivity analyzes productivity by command categories
func (pa *ProductivityAnalyzer) analyzeCategoryProductivity(commands []*models.Command) CategoryProductivityStats {
	categoryCount := make(map[categories.Category]int)
	categoryCommands := make(map[categories.Category][]*models.Command)

	// Group commands by category
	for _, cmd := range commands {
		category := pa.classifier.ClassifyCommand(cmd.Command)
		categoryCount[category]++
		categoryCommands[category] = append(categoryCommands[category], cmd)
	}

	// Calculate scores for each category
	categoryScores := make(map[categories.Category]float64)
	topCategory := categories.Unknown
	maxScore := 0.0

	for category, count := range categoryCount {
		// Calculate unique commands in category
		uniqueCommands := make(map[string]bool)
		for _, cmd := range categoryCommands[category] {
			uniqueCommands[cmd.Command] = true
		}

		// Score = frequency * diversity * category XP bonus
		diversity := float64(len(uniqueCommands)) / float64(count)
		xpBonus := pa.classifier.GetXPMultiplier(category)
		score := float64(count) * diversity * xpBonus

		categoryScores[category] = score

		if score > maxScore {
			maxScore = score
			topCategory = category
		}
	}

	// Calculate diversity score (number of categories used)
	diversityScore := float64(len(categoryCount)) / 17.0 * 100 // 17 total categories

	// Determine specialization level
	specializationLevel := "Generalist"
	if diversityScore < 25 {
		specializationLevel = "Specialist"
	} else if diversityScore < 50 {
		specializationLevel = "Focused"
	} else if diversityScore > 75 {
		specializationLevel = "Polymath"
	}

	return CategoryProductivityStats{
		TopCategory:         topCategory,
		MostImproved:        categories.Unknown, // TODO: Implement trend analysis
		CategoryScores:      categoryScores,
		DiversityScore:      diversityScore,
		SpecializationLevel: specializationLevel,
	}
}

// analyzeStreaks analyzes usage consistency and streaks
func (pa *ProductivityAnalyzer) analyzeStreaks(commands []*models.Command) StreakAnalysis {
	if len(commands) == 0 {
		return StreakAnalysis{}
	}

	// Group commands by date
	dailyActivity := make(map[string]bool)
	dates := make([]time.Time, 0)

	for _, cmd := range commands {
		dateStr := cmd.Timestamp.Format("2006-01-02")
		if !dailyActivity[dateStr] {
			dailyActivity[dateStr] = true
			dates = append(dates, cmd.Timestamp.Truncate(24*time.Hour))
		}
	}

	// Sort dates
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Calculate streaks
	currentStreak := 0
	longestStreak := 0
	tempStreak := 0
	gaps := make([]int, 0)

	today := time.Now().Truncate(24 * time.Hour)

	for i, date := range dates {
		if i == 0 {
			tempStreak = 1
			continue
		}

		daysDiff := int(date.Sub(dates[i-1]).Hours() / 24)

		if daysDiff == 1 {
			tempStreak++
		} else {
			if tempStreak > longestStreak {
				longestStreak = tempStreak
			}
			gaps = append(gaps, daysDiff-1)
			tempStreak = 1
		}
	}

	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	// Calculate current streak (from most recent activity to today)
	if len(dates) > 0 {
		lastActivity := dates[len(dates)-1]
		daysSinceLastActivity := int(today.Sub(lastActivity).Hours() / 24)

		if daysSinceLastActivity <= 1 {
			// Find current streak by going backwards
			for i := len(dates) - 1; i >= 0; i-- {
				if i == len(dates)-1 {
					currentStreak = 1
					continue
				}

				daysDiff := int(dates[i+1].Sub(dates[i]).Hours() / 24)
				if daysDiff == 1 {
					currentStreak++
				} else {
					break
				}
			}
		}
	}

	// Calculate consistency score
	daysActive := len(dates)
	totalPeriod := 1
	if len(dates) > 1 {
		totalPeriod = int(dates[len(dates)-1].Sub(dates[0]).Hours()/24) + 1
	}
	consistencyScore := float64(daysActive) / float64(totalPeriod) * 100

	// Calculate average gap
	averageGap := 0.0
	if len(gaps) > 0 {
		sum := 0
		for _, gap := range gaps {
			sum += gap
		}
		averageGap = float64(sum) / float64(len(gaps))
	}

	return StreakAnalysis{
		CurrentStreak:    currentStreak,
		LongestStreak:    longestStreak,
		ConsistencyScore: consistencyScore,
		DaysActive:       daysActive,
		WeeksActive:      int(math.Ceil(float64(totalPeriod) / 7)),
		AverageGap:       averageGap,
	}
}

// analyzeEfficiency analyzes command efficiency and patterns
func (pa *ProductivityAnalyzer) analyzeEfficiency(commands []*models.Command, sessions []*models.Session) EfficiencyMetrics {
	if len(commands) == 0 {
		return EfficiencyMetrics{}
	}

	// Calculate commands per session
	commandsPerSession := float64(len(commands))
	if len(sessions) > 0 {
		commandsPerSession = float64(len(commands)) / float64(len(sessions))
	}

	// Calculate average session length
	totalDuration := int64(0)
	validSessions := 0
	for _, session := range sessions {
		if session.EndTime != nil {
			duration := session.EndTime.Sub(session.StartTime)
			totalDuration += duration.Nanoseconds()
			validSessions++
		}
	}

	averageSessionLength := 0.0
	if validSessions > 0 {
		averageSessionLength = float64(totalDuration) / float64(validSessions) / float64(time.Minute)
	}

	// Calculate command uniqueness
	uniqueCommands := make(map[string]bool)
	commandFrequency := make(map[string]int)

	for _, cmd := range commands {
		uniqueCommands[cmd.Command] = true
		commandFrequency[cmd.Command]++
	}

	uniqueCommandRatio := float64(len(uniqueCommands)) / float64(len(commands)) * 100

	// Calculate repetitive commands percentage
	repetitiveCount := 0
	for _, freq := range commandFrequency {
		if freq > 3 { // Commands used more than 3 times
			repetitiveCount += freq
		}
	}
	repetitiveCommands := float64(repetitiveCount) / float64(len(commands)) * 100

	// Calculate complexity score based on command categories
	complexCategories := []categories.Category{
		categories.Development, categories.Docker, categories.Kubernetes,
		categories.Cloud, categories.Security, categories.Database,
	}

	complexCommandCount := 0
	for _, cmd := range commands {
		category := pa.classifier.ClassifyCommand(cmd.Command)
		for _, complex := range complexCategories {
			if category == complex {
				complexCommandCount++
				break
			}
		}
	}

	complexityScore := float64(complexCommandCount) / float64(len(commands)) * 100

	// Calculate automation potential (high repetition + low complexity = high automation potential)
	automationPotential := (repetitiveCommands * 0.7) + ((100 - complexityScore) * 0.3)

	return EfficiencyMetrics{
		CommandsPerSession:   commandsPerSession,
		AverageSessionLength: averageSessionLength,
		RepetitiveCommands:   repetitiveCommands,
		UniqueCommandRatio:   uniqueCommandRatio,
		ComplexityScore:      complexityScore,
		AutomationPotential:  automationPotential,
	}
}

// analyzeTimeDistribution analyzes how time is distributed across activities
func (pa *ProductivityAnalyzer) analyzeTimeDistribution(commands []*models.Command) map[string]float64 {
	categoryTime := make(map[string]int)
	totalCommands := len(commands)

	for _, cmd := range commands {
		category := pa.classifier.ClassifyCommand(cmd.Command)
		info := pa.classifier.GetCategoryInfo(category)
		categoryTime[info.Name]++
	}

	distribution := make(map[string]float64)
	for category, count := range categoryTime {
		distribution[category] = float64(count) / float64(totalCommands) * 100
	}

	return distribution
}

// calculateOverallScore calculates a comprehensive productivity score
func (pa *ProductivityAnalyzer) calculateOverallScore(metrics *ProductivityMetrics) float64 {
	// Weighted scoring system
	streakWeight := 0.3
	diversityWeight := 0.2
	efficiencyWeight := 0.2
	consistencyWeight := 0.15
	complexityWeight := 0.15

	// Normalize scores to 0-100 range
	streakScore := math.Min(float64(metrics.StreakAnalysis.LongestStreak)*5, 100)
	diversityScore := metrics.CategoryInsights.DiversityScore
	efficiencyScore := math.Min(metrics.EfficiencyMetrics.CommandsPerSession*10, 100)
	consistencyScore := metrics.StreakAnalysis.ConsistencyScore
	complexityScore := metrics.EfficiencyMetrics.ComplexityScore

	overallScore := (streakScore * streakWeight) +
		(diversityScore * diversityWeight) +
		(efficiencyScore * efficiencyWeight) +
		(consistencyScore * consistencyWeight) +
		(complexityScore * complexityWeight)

	return math.Min(overallScore, 100)
}

// FormatProductivityReport generates a formatted productivity report
func (pa *ProductivityAnalyzer) FormatProductivityReport(metrics *ProductivityMetrics) string {
	if metrics.OverallScore == 0 {
		return "📊 No productivity data available yet. Start using your terminal!"
	}

	report := fmt.Sprintf("📊 Productivity Analysis Report\n")
	report += fmt.Sprintf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Overall Score
	scoreEmoji := "🔥"
	if metrics.OverallScore < 50 {
		scoreEmoji = "📈"
	} else if metrics.OverallScore < 75 {
		scoreEmoji = "⭐"
	}

	report += fmt.Sprintf("%s Overall Productivity Score: %.1f/100\n\n", scoreEmoji, metrics.OverallScore)

	// Time Patterns
	report += fmt.Sprintf("⏰ Time Patterns:\n")
	report += fmt.Sprintf("  🌅 Morning (6-12): %.1f%%\n", metrics.DailyPattern.MorningScore)
	report += fmt.Sprintf("  ☀️ Afternoon (12-18): %.1f%%\n", metrics.DailyPattern.AfternoonScore)
	report += fmt.Sprintf("  🌆 Evening (18-24): %.1f%%\n", metrics.DailyPattern.EveningScore)
	report += fmt.Sprintf("  🌙 Night (0-6): %.1f%%\n\n", metrics.DailyPattern.NightScore)

	// Peak Performance
	if len(metrics.DailyPattern.PeakHours) > 0 {
		report += fmt.Sprintf("🎯 Peak Performance: %02d:00 (Score: %.1f)\n\n",
			metrics.DailyPattern.PeakHours[0], metrics.DailyPattern.PeakScore)
	}

	// Category Insights
	report += fmt.Sprintf("📂 Category Focus:\n")
	topCategoryInfo := pa.classifier.GetCategoryInfo(metrics.CategoryInsights.TopCategory)
	report += fmt.Sprintf("  %s Top Category: %s\n", topCategoryInfo.Icon, topCategoryInfo.Name)
	report += fmt.Sprintf("  🎨 Diversity: %.1f%% (%s)\n",
		metrics.CategoryInsights.DiversityScore, metrics.CategoryInsights.SpecializationLevel)

	// Consistency
	report += fmt.Sprintf("\n🔥 Consistency:\n")
	report += fmt.Sprintf("  Current Streak: %d days\n", metrics.StreakAnalysis.CurrentStreak)
	report += fmt.Sprintf("  Longest Streak: %d days\n", metrics.StreakAnalysis.LongestStreak)
	report += fmt.Sprintf("  Consistency Score: %.1f%%\n", metrics.StreakAnalysis.ConsistencyScore)

	// Efficiency
	report += fmt.Sprintf("\n⚡ Efficiency:\n")
	report += fmt.Sprintf("  Commands/Session: %.1f\n", metrics.EfficiencyMetrics.CommandsPerSession)
	report += fmt.Sprintf("  Command Uniqueness: %.1f%%\n", metrics.EfficiencyMetrics.UniqueCommandRatio)
	report += fmt.Sprintf("  Complexity Score: %.1f%%\n", metrics.EfficiencyMetrics.ComplexityScore)

	// Automation Potential
	if metrics.EfficiencyMetrics.AutomationPotential > 60 {
		report += fmt.Sprintf("\n🤖 Automation Opportunity: %.1f%% - Consider scripting repetitive tasks!\n",
			metrics.EfficiencyMetrics.AutomationPotential)
	}

	return report
}
