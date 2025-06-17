package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal dashboard",
	Long: `Launch a beautiful, interactive terminal user interface (TUI) dashboard.
Navigate through different tabs to view your productivity stats, analytics,
heatmaps, and achievements in real-time.

Features:
  • 📊 Overview - Quick stats and recent commands
  • 📈 Analytics - Deep productivity insights
  • 🔥 Heatmap - Time-based activity visualization
  • 🏆 Achievements - Gamification progress
  • ⚙️ Settings - Customize your experience

Use arrow keys or h/l to navigate tabs, q to quit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUICommand(cmd, args)
	},
}

func runTUICommand(cmd *cobra.Command, args []string) error {
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

	// Run the simple interactive dashboard
	if err := tui.RunSimpleDashboard(db); err != nil {
		return fmt.Errorf("failed to run TUI dashboard: %w", err)
	}

	return nil
}
