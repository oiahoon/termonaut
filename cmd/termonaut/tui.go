package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/tui"
	"github.com/oiahoon/termonaut/internal/tui/enhanced"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal dashboard",
	Long: `Launch the interactive terminal user interface (TUI) dashboard.
By default uses intelligent mode that adapts to your terminal size.

Mode Options:
  --mode smart     智能模式 (default) - Auto-adapts to terminal size
  --mode compact   普通模式 - Compact layout with small avatars
  --mode full      完整模式 - Full-featured with large avatars  
  --mode classic   经典模式 - Original TUI for compatibility
  --mode minimal   极简模式 - Text-only stats output

Configuration:
  You can set your preferred mode in ~/.termonaut/config.toml:
  [ui]
  default_mode = "smart"  # smart, compact, full, classic, minimal

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
	// Get mode flag
	mode, _ := cmd.Flags().GetString("mode")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If no mode specified, check config, then default to smart
	if mode == "" {
		if cfg.UI.DefaultMode != "" {
			mode = cfg.UI.DefaultMode
		} else {
			mode = "smart"
		}
	}

	// Handle minimal mode (stats output)
	if mode == "minimal" {
		return runMinimalMode(cfg)
	}

	// Initialize logger
	logger := setupLogger(cfg.LogLevel)

	// Initialize database
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Choose TUI implementation based on mode
	switch mode {
	case "classic":
		return runClassicTUI(db)
	case "compact", "full", "smart":
		return runEnhancedTUI(db, mode)
	default:
		return fmt.Errorf("unknown mode: %s. Available modes: smart, compact, full, classic, minimal", mode)
	}
}

func runMinimalMode(cfg *config.Config) error {
	// This will call the stats command functionality
	return runStatsCommand(nil, []string{})
}

func runClassicTUI(db *database.DB) error {
	// Run the original simple interactive dashboard
	if err := tui.RunSimpleDashboard(db); err != nil {
		return fmt.Errorf("failed to run classic TUI dashboard: %w", err)
	}
	return nil
}

func runEnhancedTUI(db *database.DB, mode string) error {
	// Create enhanced dashboard
	dashboard := enhanced.NewEnhancedDashboard(db)
	
	// Set mode preference (the dashboard will adapt accordingly)
	dashboard.SetModePreference(mode)

	// Run the TUI
	program := tea.NewProgram(
		dashboard,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run TUI dashboard: %w", err)
	}

	return nil
}

func init() {
	// Add mode flag
	tuiCmd.Flags().StringP("mode", "m", "", "TUI mode: smart, compact, full, classic, minimal")
}
