package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/spf13/cobra"
)

var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Show productivity heatmap and time patterns",
	Long: `Display a visual heatmap of your productivity patterns across days and hours.
Identify your peak performance times and working patterns.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runHeatmapCommand(cmd, args)
	},
}

func runHeatmapCommand(cmd *cobra.Command, args []string) error {
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

	// Get all commands
	commands, err := db.GetAllCommands()
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	if len(commands) == 0 {
		fmt.Println("📊 No heatmap data available yet!")
		fmt.Println("Start using your terminal to generate time pattern insights.")
		return nil
	}

	// Generate heatmap
	heatmapAnalyzer := analytics.NewHeatmapAnalyzer()
	heatmapData := heatmapAnalyzer.GenerateHeatmap(commands)

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		fmt.Printf("%+v\n", heatmapData)
	} else {
		visualization := heatmapAnalyzer.FormatHeatmapVisualization(heatmapData)
		fmt.Print(visualization)

		// Add insights
		insights := heatmapAnalyzer.GenerateTimeInsights(heatmapData, commands)
		if len(insights.Recommendations) > 0 {
			fmt.Printf("\n💡 Recommendations:\n")
			for i, rec := range insights.Recommendations {
				fmt.Printf("  %d. %s\n", i+1, rec)
			}
		}

		if len(insights.ProductivityTips) > 0 {
			fmt.Printf("\n🎯 Productivity Tips:\n")
			for _, tip := range insights.ProductivityTips {
				fmt.Printf("  • %s: %s (Impact: %s, Effort: %s)\n",
					tip.Title, tip.Description, tip.Impact, tip.Effort)
			}
		}
	}

	return nil
}
