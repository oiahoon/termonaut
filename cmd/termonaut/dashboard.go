package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/internal/visualization"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Show comprehensive productivity dashboard",
	Long: `Display a complete dashboard with analytics, charts, and insights.
Includes activity patterns, command statistics, productivity metrics, and visual charts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDashboardCommand(cmd, args)
	},
}

func runDashboardCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger := setupLogger(cfg.LogLevel)

	// Initialize database
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get all data
	commands, err := db.GetAllCommands()
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	sessions, err := db.GetAllSessions()
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	if len(commands) == 0 {
		fmt.Println("🚀 Welcome to Termonaut!")
		fmt.Println("Start using your terminal to see amazing insights and analytics.")
		fmt.Println()
		fmt.Println("Try these commands to get started:")
		fmt.Println("  • termonaut init      - Setup shell integration")
		fmt.Println("  • termonaut progress  - View your progress")
		fmt.Println("  • termonaut help      - See all available commands")
		return nil
	}

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if jsonOutput {
		// JSON output for programmatic access
		analyzer := analytics.NewProductivityAnalyzer()
		heatmapAnalyzer := analytics.NewHeatmapAnalyzer()

		metrics := analyzer.AnalyzeProductivity(commands, sessions)
		heatmapData := heatmapAnalyzer.GenerateHeatmap(commands)

		dashboardData := map[string]interface{}{
			"productivity_metrics": metrics,
			"heatmap_data":         heatmapData,
			"summary": map[string]interface{}{
				"total_commands": len(commands),
				"total_sessions": len(sessions),
				"data_available": true,
			},
		}

		fmt.Printf("%+v\n", dashboardData)
		return nil
	}

	// Create visual dashboard
	renderer := visualization.NewChartRenderer(60, 12)

	// Header
	fmt.Println("🎮 Termonaut - Complete Productivity Dashboard")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Render main dashboard
	dashboard := renderer.RenderSummaryDashboard(commands, sessions)
	fmt.Print(dashboard)

	// Gamification stats
	statsCalculator := stats.New(db)
	gamificationStats, err := statsCalculator.GetGamificationStats()
	if err == nil {
		gamificationOutput := statsCalculator.FormatGamificationStats(gamificationStats)
		fmt.Print(gamificationOutput)
	}

	// Analytics insights
	analyzer := analytics.NewProductivityAnalyzer()
	metrics := analyzer.AnalyzeProductivity(commands, sessions)

	fmt.Printf("📊 Productivity Analysis\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	fmt.Printf("🎯 Overall Score: %.1f/100\n", metrics.OverallScore)
	fmt.Printf("⚡ Efficiency: %.1f%% | 🔄 Consistency: %.1f%%\n",
		metrics.EfficiencyMetrics.UniqueCommandRatio*100,
		metrics.StreakAnalysis.ConsistencyScore*100)

	fmt.Printf("🏆 Top Category: %s\n", metrics.CategoryInsights.SpecializationLevel)

	// Show peak hours from daily pattern
	if len(metrics.DailyPattern.PeakHours) > 0 {
		fmt.Printf("📈 Peak Hour: %02d:00 | 🎪 Diversity Score: %.1f%%\n",
			metrics.DailyPattern.PeakHours[0],
			metrics.CategoryInsights.DiversityScore*100)
	} else {
		fmt.Printf("📈 Peak Time: Not enough data | 🎪 Diversity Score: %.1f%%\n",
			metrics.CategoryInsights.DiversityScore*100)
	}

	// Recent activity heatmap
	fmt.Printf("\n")
	heatmapAnalyzer := analytics.NewHeatmapAnalyzer()
	heatmapData := heatmapAnalyzer.GenerateHeatmap(commands)
	visualization := heatmapAnalyzer.FormatHeatmapVisualization(heatmapData)
	fmt.Print(visualization)

	// Footer with tips
	fmt.Printf("\n💡 Pro Tips:\n")
	fmt.Printf("  • Use 'termonaut heatmap' for detailed time analysis\n")
	fmt.Printf("  • Use 'termonaut analytics' for deep productivity insights\n")
	fmt.Printf("  • Use 'termonaut progress' to track your achievements\n")
	fmt.Printf("  • Use 'termonaut categories' to see command classification\n")

	return nil
}
