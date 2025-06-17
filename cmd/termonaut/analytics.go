package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/analytics"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/spf13/cobra"
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Show advanced productivity analytics",
	Long: `Display comprehensive productivity analysis including time patterns,
efficiency metrics, category insights, and automation opportunities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAnalyticsCommand(cmd, args)
	},
}

func runAnalyticsCommand(cmd *cobra.Command, args []string) error {
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

	// Get all commands and sessions
	commands, err := db.GetAllCommands()
	if err != nil {
		return fmt.Errorf("failed to get commands: %w", err)
	}

	sessions, err := db.GetAllSessions()
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	if len(commands) == 0 {
		fmt.Println("📊 No analytics data available yet!")
		fmt.Println("Start using your terminal to generate productivity insights.")
		return nil
	}

	// Analyze productivity
	analyzer := analytics.NewProductivityAnalyzer()
	metrics := analyzer.AnalyzeProductivity(commands, sessions)

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		fmt.Printf("%+v\n", metrics)
	} else {
		report := analyzer.FormatProductivityReport(metrics)
		fmt.Print(report)
	}

	return nil
}
