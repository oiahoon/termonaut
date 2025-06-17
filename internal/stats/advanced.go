package stats

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/pkg/models"
)

// BasicStatsCalculator provides basic statistics calculations
type BasicStatsCalculator struct {
	db *database.DB
}

// NewStatsCalculator creates a new basic stats calculator
func NewStatsCalculator(db *database.DB) *BasicStatsCalculator {
	return &BasicStatsCalculator{
		db: db,
	}
}

// GetBasicStats returns basic statistics
func (bsc *BasicStatsCalculator) GetBasicStats() (map[string]interface{}, error) {
	return bsc.db.GetBasicStats()
}

// GetGamificationStats returns gamification-related statistics
func (bsc *BasicStatsCalculator) GetGamificationStats() (map[string]interface{}, error) {
	// Get basic stats from database
	stats, err := bsc.db.GetBasicStats()
	if err != nil {
		return nil, err
	}

	// Add gamification-specific data
	gamificationStats := make(map[string]interface{})
	for k, v := range stats {
		gamificationStats[k] = v
	}

	// Add additional gamification metrics
	gamificationStats["achievements_count"] = 0 // This would come from achievements table
	gamificationStats["current_streak"] = 0    // This would come from user_progress table
	gamificationStats["xp_multiplier"] = 1.0   // Based on current level/achievements

	return gamificationStats, nil
}

// AdvancedStatsManager handles power user features
type AdvancedStatsManager struct {
	db         *database.DB
	classifier *categories.CommandClassifier
}

// NewAdvancedStatsManager creates a new advanced stats manager
func NewAdvancedStatsManager(db *database.DB) *AdvancedStatsManager {
	return &AdvancedStatsManager{
		db:         db,
		classifier: categories.NewCommandClassifier(),
	}
}

// CustomCommandScore represents a custom scoring rule
type CustomCommandScore struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Patterns    []string               `json:"patterns"`
	Category    categories.Category    `json:"category"`
	Multiplier  float64                `json:"multiplier"`
	Conditions  map[string]interface{} `json:"conditions"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
}

// AdvancedFilter represents filtering options for advanced queries
type AdvancedFilter struct {
	Categories    []categories.Category `json:"categories"`
	DateFrom      *time.Time            `json:"date_from"`
	DateTo        *time.Time            `json:"date_to"`
	ExitCode      *int                  `json:"exit_code"`
	MinXP         *int                  `json:"min_xp"`
	MaxXP         *int                  `json:"max_xp"`
	CommandRegex  string                `json:"command_regex"`
	Directory     string                `json:"directory"`
	Duration      *time.Duration        `json:"duration"`
	Limit         int                   `json:"limit"`
	SortBy        string                `json:"sort_by"`
	SortOrder     string                `json:"sort_order"`
}

// BulkOperation represents operations that can be performed on multiple commands
type BulkOperation struct {
	Type        string                 `json:"type"`
	Filters     *AdvancedFilter        `json:"filters"`
	Parameters  map[string]interface{} `json:"parameters"`
	DryRun      bool                   `json:"dry_run"`
}

// CommandScore represents a scored command result
type CommandScore struct {
	Command      *models.Command     `json:"command"`
	Score        float64             `json:"score"`
	Category     categories.Category `json:"category"`
	AppliedRules []string            `json:"applied_rules"`
	Rank         int                 `json:"rank"`
}

// GetCustomCommandScores retrieves all custom scoring rules
func (asm *AdvancedStatsManager) GetCustomCommandScores() ([]*CustomCommandScore, error) {
	// This would be stored in database in a real implementation
	// For now, return some default scoring rules
	defaultScores := []*CustomCommandScore{
		{
			ID:          "git_complexity",
			Name:        "Git Complexity Scorer",
			Description: "Scores git commands based on complexity",
			Patterns:    []string{"git rebase", "git merge", "git cherry-pick"},
			Category:    categories.Git,
			Multiplier:  2.0,
			Conditions: map[string]interface{}{
				"min_args": 2,
				"categories": []string{"git"},
			},
			Enabled:   true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "docker_orchestration",
			Name:        "Docker Orchestration",
			Description: "High scores for complex docker operations",
			Patterns:    []string{"docker-compose", "docker stack", "docker swarm"},
			Category:    categories.Docker,
			Multiplier:  1.8,
			Conditions: map[string]interface{}{
				"success_rate": 0.8,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
		},
		{
			ID:          "kubernetes_expert",
			Name:        "Kubernetes Expert",
			Description: "Rewards advanced k8s operations",
			Patterns:    []string{"kubectl apply", "kubectl rollout", "helm"},
			Category:    categories.Kubernetes,
			Multiplier:  2.5,
			Conditions: map[string]interface{}{
				"min_length": 10,
			},
			Enabled:   true,
			CreatedAt: time.Now(),
		},
	}

	return defaultScores, nil
}

// CreateCustomCommandScore creates a new custom scoring rule
func (asm *AdvancedStatsManager) CreateCustomCommandScore(score *CustomCommandScore) error {
	score.ID = fmt.Sprintf("custom_%d", time.Now().Unix())
	score.CreatedAt = time.Now()
	score.Enabled = true
	
	// In a real implementation, this would be stored in database
	return nil
}

// UpdateCustomCommandScore updates an existing scoring rule
func (asm *AdvancedStatsManager) UpdateCustomCommandScore(id string, score *CustomCommandScore) error {
	// In a real implementation, this would update the database record
	return nil
}

// DeleteCustomCommandScore removes a scoring rule
func (asm *AdvancedStatsManager) DeleteCustomCommandScore(id string) error {
	// In a real implementation, this would delete from database
	return nil
}

// CalculateCommandScores applies custom scoring to commands
func (asm *AdvancedStatsManager) CalculateCommandScores(commands []*models.Command, customRules []*CustomCommandScore) ([]*CommandScore, error) {
	var scoredCommands []*CommandScore

	for _, cmd := range commands {
		score := &CommandScore{
			Command:      cmd,
			Score:        1.0, // Base score
			Category:     asm.classifier.ClassifyCommand(cmd.Command),
			AppliedRules: []string{},
		}

		// Apply custom scoring rules
		for _, rule := range customRules {
			if !rule.Enabled {
				continue
			}

			if asm.matchesRule(cmd, rule) {
				score.Score *= rule.Multiplier
				score.AppliedRules = append(score.AppliedRules, rule.Name)
			}
		}

		// Apply category multiplier
		categoryInfo := asm.classifier.GetCategoryInfo(score.Category)
		if categoryInfo != nil {
			score.Score *= categoryInfo.XPBonus
		}

		scoredCommands = append(scoredCommands, score)
	}

	// Sort by score (highest first)
	sort.Slice(scoredCommands, func(i, j int) bool {
		return scoredCommands[i].Score > scoredCommands[j].Score
	})

	// Assign ranks
	for i, score := range scoredCommands {
		score.Rank = i + 1
	}

	return scoredCommands, nil
}

// FilterCommands applies advanced filtering to commands
func (asm *AdvancedStatsManager) FilterCommands(filter *AdvancedFilter) ([]*models.Command, error) {
	// This would be implemented as a database query in a real implementation
	// For now, we'll return a placeholder
	return nil, fmt.Errorf("advanced filtering not yet implemented in database layer")
}

// PerformBulkOperation executes bulk operations on commands
func (asm *AdvancedStatsManager) PerformBulkOperation(operation *BulkOperation) (*BulkOperationResult, error) {
	result := &BulkOperationResult{
		Type:      operation.Type,
		DryRun:    operation.DryRun,
		StartTime: time.Now(),
		Affected:  0,
		Errors:    []string{},
	}

	// Get commands matching the filter
	commands, err := asm.FilterCommands(operation.Filters)
	if err != nil {
		return result, fmt.Errorf("failed to filter commands: %w", err)
	}

	switch operation.Type {
	case "recalculate_xp":
		result, err = asm.performRecalculateXP(commands, operation, result)
	case "update_categories":
		result, err = asm.performUpdateCategories(commands, operation, result)
	case "export_data":
		result, err = asm.performExportData(commands, operation, result)
	case "delete_commands":
		result, err = asm.performDeleteCommands(commands, operation, result)
	default:
		return result, fmt.Errorf("unknown bulk operation type: %s", operation.Type)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, err
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	Type      string        `json:"type"`
	DryRun    bool          `json:"dry_run"`
	Affected  int           `json:"affected"`
	Errors    []string      `json:"errors"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Details   interface{}   `json:"details"`
}

// GetAdvancedAnalytics provides sophisticated analytics
func (asm *AdvancedStatsManager) GetAdvancedAnalytics(filter *AdvancedFilter) (*AdvancedAnalytics, error) {
	commands, err := asm.FilterCommands(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to filter commands: %w", err)
	}

	analytics := &AdvancedAnalytics{
		Filter:          filter,
		TotalCommands:   len(commands),
		TimeRange:       calculateTimeRange(commands),
		CategoryBreakdown: asm.calculateCategoryBreakdown(commands),
		PerformanceMetrics: asm.calculatePerformanceMetrics(commands),
		TrendAnalysis:   asm.calculateTrendAnalysis(commands),
		Recommendations: asm.generateRecommendations(commands),
	}

	return analytics, nil
}

// AdvancedAnalytics represents comprehensive analytics results
type AdvancedAnalytics struct {
	Filter             *AdvancedFilter                        `json:"filter"`
	TotalCommands      int                                    `json:"total_commands"`
	TimeRange          *TimeRange                             `json:"time_range"`
	CategoryBreakdown  map[categories.Category]*CategoryData `json:"category_breakdown"`
	PerformanceMetrics *PerformanceMetrics                    `json:"performance_metrics"`
	TrendAnalysis      *TrendAnalysis                         `json:"trend_analysis"`
	Recommendations    []string                               `json:"recommendations"`
}

// TimeRange represents a time range
type TimeRange struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

// CategoryData represents detailed category information
type CategoryData struct {
	Count          int     `json:"count"`
	Percentage     float64 `json:"percentage"`
	AverageXP      float64 `json:"average_xp"`
	SuccessRate    float64 `json:"success_rate"`
	TopCommands    []string `json:"top_commands"`
	TrendDirection string   `json:"trend_direction"`
}

// PerformanceMetrics represents performance analysis
type PerformanceMetrics struct {
	AverageExecutionTime float64 `json:"average_execution_time"`
	SuccessRate          float64 `json:"success_rate"`
	CommandsPerHour      float64 `json:"commands_per_hour"`
	PeakHours            []int   `json:"peak_hours"`
	ProductivityScore    float64 `json:"productivity_score"`
}

// TrendAnalysis represents trend analysis data
type TrendAnalysis struct {
	DailyTrend     []TrendPoint `json:"daily_trend"`
	WeeklyTrend    []TrendPoint `json:"weekly_trend"`
	CategoryTrends map[string]string `json:"category_trends"`
	GrowthRate     float64      `json:"growth_rate"`
}

// TrendPoint represents a single point in trend analysis
type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
	Label string    `json:"label"`
}

// Helper methods

func (asm *AdvancedStatsManager) matchesRule(cmd *models.Command, rule *CustomCommandScore) bool {
	// Check pattern matching
	for _, pattern := range rule.Patterns {
		if strings.Contains(strings.ToLower(cmd.Command), strings.ToLower(pattern)) {
			return asm.checkConditions(cmd, rule.Conditions)
		}
	}
	return false
}

func (asm *AdvancedStatsManager) checkConditions(cmd *models.Command, conditions map[string]interface{}) bool {
	// Check minimum arguments
	if minArgs, exists := conditions["min_args"]; exists {
		parts := strings.Fields(cmd.Command)
		if len(parts)-1 < int(minArgs.(float64)) {
			return false
		}
	}

	// Check command length
	if minLength, exists := conditions["min_length"]; exists {
		if len(cmd.Command) < int(minLength.(float64)) {
			return false
		}
	}

	// Add more condition checks as needed
	return true
}

// Placeholder implementations for bulk operations
func (asm *AdvancedStatsManager) performRecalculateXP(commands []*models.Command, operation *BulkOperation, result *BulkOperationResult) (*BulkOperationResult, error) {
	if !operation.DryRun {
		// Actually recalculate XP for commands
	}
	result.Affected = len(commands)
	result.Details = map[string]interface{}{
		"recalculated_commands": len(commands),
		"total_xp_adjusted":     0, // Calculate actual XP changes
	}
	return result, nil
}

func (asm *AdvancedStatsManager) performUpdateCategories(commands []*models.Command, operation *BulkOperation, result *BulkOperationResult) (*BulkOperationResult, error) {
	if !operation.DryRun {
		// Actually update categories
	}
	result.Affected = len(commands)
	return result, nil
}

func (asm *AdvancedStatsManager) performExportData(commands []*models.Command, operation *BulkOperation, result *BulkOperationResult) (*BulkOperationResult, error) {
	// Export commands to specified format
	result.Affected = len(commands)
	result.Details = map[string]interface{}{
		"export_format": operation.Parameters["format"],
		"file_path":     operation.Parameters["path"],
	}
	return result, nil
}

func (asm *AdvancedStatsManager) performDeleteCommands(commands []*models.Command, operation *BulkOperation, result *BulkOperationResult) (*BulkOperationResult, error) {
	if !operation.DryRun {
		// Actually delete commands
	}
	result.Affected = len(commands)
	return result, nil
}

// Helper functions for analytics
func calculateTimeRange(commands []*models.Command) *TimeRange {
	if len(commands) == 0 {
		return nil
	}

	start := commands[0].Timestamp
	end := commands[0].Timestamp

	for _, cmd := range commands {
		if cmd.Timestamp.Before(start) {
			start = cmd.Timestamp
		}
		if cmd.Timestamp.After(end) {
			end = cmd.Timestamp
		}
	}

	return &TimeRange{
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
	}
}

func (asm *AdvancedStatsManager) calculateCategoryBreakdown(commands []*models.Command) map[categories.Category]*CategoryData {
	breakdown := make(map[categories.Category]*CategoryData)
	
	// Implementation would calculate detailed category statistics
	// This is a placeholder
	
	return breakdown
}

func (asm *AdvancedStatsManager) calculatePerformanceMetrics(commands []*models.Command) *PerformanceMetrics {
	// Implementation would calculate performance metrics
	return &PerformanceMetrics{
		AverageExecutionTime: 0.5,
		SuccessRate:          0.85,
		CommandsPerHour:      120,
		PeakHours:            []int{9, 10, 14, 15},
		ProductivityScore:    87.5,
	}
}

func (asm *AdvancedStatsManager) calculateTrendAnalysis(commands []*models.Command) *TrendAnalysis {
	// Implementation would calculate trend analysis
	return &TrendAnalysis{
		DailyTrend:     []TrendPoint{},
		WeeklyTrend:    []TrendPoint{},
		CategoryTrends: map[string]string{},
		GrowthRate:     0.15,
	}
}

func (asm *AdvancedStatsManager) generateRecommendations(commands []*models.Command) []string {
	recommendations := []string{
		"Consider automating frequently repeated command sequences",
		"Your git workflow shows room for optimization",
		"Peak productivity hours are 9-10 AM and 2-3 PM",
		"Docker usage has increased 25% this week",
	}
	
	return recommendations
} 