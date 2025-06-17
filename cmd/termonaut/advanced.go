package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/oiahoon/termonaut/internal/api"
	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/oiahoon/termonaut/internal/stats"
)

// initAdvancedDB initializes the database for advanced commands
func initAdvancedDB() (*database.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger := logrus.New()
	if cfg.LogLevel == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	db, err := database.New(cfg.DataDir, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

// createAdvancedCmd creates the advanced features command group
func createAdvancedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advanced",
		Short: "🚀 Advanced power user features (v0.8.0)",
		Long: `Advanced power user features for terminal productivity analysis.

Includes:
• Custom command scoring and filtering
• Bulk operations on command data  
• Shell integration management
• API server for external integrations
• Advanced analytics and insights`,
	}

	// Add subcommands
	cmd.AddCommand(createScoringCmd())
	cmd.AddCommand(createFilterCmd())
	cmd.AddCommand(createBulkCmd())
	cmd.AddCommand(createShellCmd())
	cmd.AddCommand(createAPICmd())
	cmd.AddCommand(createAnalyticsCmd())

	return cmd
}

// createScoringCmd creates command scoring management
func createScoringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scoring",
		Short: "🎯 Manage custom command scoring rules",
		Long:  "Create, update, and manage custom scoring rules for commands",
	}

	// List scoring rules
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "📋 List all custom scoring rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			advancedStats := stats.NewAdvancedStatsManager(db)
			rules, err := advancedStats.GetCustomCommandScores()
			if err != nil {
				return fmt.Errorf("failed to get scoring rules: %w", err)
			}

			if len(rules) == 0 {
				fmt.Println("📝 No custom scoring rules found")
				fmt.Println("\nCreate one with: termonaut advanced scoring create")
				return nil
			}

			fmt.Printf("🎯 Custom Scoring Rules (%d)\n\n", len(rules))
			for _, rule := range rules {
				status := "✅ Enabled"
				if !rule.Enabled {
					status = "❌ Disabled"
				}

				fmt.Printf("📌 %s\n", rule.Name)
				fmt.Printf("   ID: %s\n", rule.ID)
				fmt.Printf("   Description: %s\n", rule.Description)
				fmt.Printf("   Multiplier: %.1fx\n", rule.Multiplier)
				fmt.Printf("   Category: %s\n", rule.Category)
				fmt.Printf("   Patterns: %s\n", strings.Join(rule.Patterns, ", "))
				fmt.Printf("   Status: %s\n", status)
				fmt.Printf("   Created: %s\n\n", rule.CreatedAt.Format("2006-01-02 15:04"))
			}

			return nil
		},
	}

	// Create scoring rule
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "✨ Create a new custom scoring rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			// Interactive creation (simplified for now)
			fmt.Println("🎯 Create Custom Scoring Rule")
			fmt.Println("===============================")
			
			var name string
			fmt.Print("Rule name: ")
			fmt.Scanln(&name)

			var description string
			fmt.Print("Description: ")
			fmt.Scanln(&description)

			var multiplier float64
			fmt.Print("XP Multiplier (e.g., 1.5): ")
			fmt.Scanln(&multiplier)

			var patterns string
			fmt.Print("Command patterns (comma-separated): ")
			fmt.Scanln(&patterns)

			rule := &stats.CustomCommandScore{
				Name:        name,
				Description: description,
				Multiplier:  multiplier,
				Patterns:    strings.Split(patterns, ","),
				Category:    categories.Development, // Default
				Enabled:     true,
				Conditions:  map[string]interface{}{},
			}

			advancedStats := stats.NewAdvancedStatsManager(db)
			if err := advancedStats.CreateCustomCommandScore(rule); err != nil {
				return fmt.Errorf("failed to create scoring rule: %w", err)
			}

			fmt.Printf("✅ Created scoring rule: %s\n", name)
			return nil
		},
	}

	cmd.AddCommand(listCmd)
	cmd.AddCommand(createCmd)
	return cmd
}

// createFilterCmd creates advanced filtering commands
func createFilterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter",
		Short: "🔍 Advanced command filtering and search",
		Long:  "Filter commands with advanced criteria including date ranges, categories, exit codes, and more",
	}

	var dateFrom, dateTo string
	var categoryStrings []string
	var exitCode int
	var limit int
	var commandRegex string

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "🔍 Search commands with advanced filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			// Build filter
			filter := &stats.AdvancedFilter{
				Limit:        limit,
				CommandRegex: commandRegex,
			}

			if exitCode >= 0 {
				filter.ExitCode = &exitCode
			}

			if dateFrom != "" {
				if from, err := time.Parse("2006-01-02", dateFrom); err == nil {
					filter.DateFrom = &from
				}
			}

			if dateTo != "" {
				if to, err := time.Parse("2006-01-02", dateTo); err == nil {
					filter.DateTo = &to
				}
			}

			// Convert category strings to enum
			if len(categoryStrings) > 0 {
				classifier := categories.NewCommandClassifier()
				allCategories := classifier.GetAllCategories()
				
				for _, catStr := range categoryStrings {
					for cat := range allCategories {
						if strings.EqualFold(string(cat), catStr) {
							filter.Categories = append(filter.Categories, cat)
							break
						}
					}
				}
			}

			advancedStats := stats.NewAdvancedStatsManager(db)
			commands, err := advancedStats.FilterCommands(filter)
			if err != nil {
				fmt.Printf("⚠️ Advanced filtering not yet implemented: %v\n", err)
				fmt.Println("📝 Using basic recent commands instead...")
				
				// Fallback to basic command retrieval
				commands, err = db.GetRecentCommands(limit)
				if err != nil {
					return fmt.Errorf("failed to get commands: %w", err)
				}
			}

			if len(commands) == 0 {
				fmt.Println("🔍 No commands match your filters")
				return nil
			}

			fmt.Printf("🔍 Filtered Commands (%d results)\n\n", len(commands))
			for i, cmd := range commands {
				if i >= limit {
					break
				}
				
				status := "✅"
				if cmd.ExitCode != 0 {
					status = "❌"
				}

				classifier := categories.NewCommandClassifier()
				category := classifier.ClassifyCommand(cmd.Command)
				categoryInfo := classifier.GetCategoryInfo(category)

				fmt.Printf("%s [%s] %s %s\n", 
					status, 
					cmd.Timestamp.Format("15:04:05"),
					categoryInfo.Icon,
					cmd.Command)
			}

			return nil
		},
	}

	searchCmd.Flags().StringVar(&dateFrom, "from", "", "Start date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&dateTo, "to", "", "End date (YYYY-MM-DD)")
	searchCmd.Flags().StringSliceVar(&categoryStrings, "categories", []string{}, "Filter by categories")
	searchCmd.Flags().IntVar(&exitCode, "exit-code", -1, "Filter by exit code")
	searchCmd.Flags().IntVar(&limit, "limit", 50, "Maximum number of results")
	searchCmd.Flags().StringVar(&commandRegex, "regex", "", "Command regex pattern")

	cmd.AddCommand(searchCmd)
	return cmd
}

// createBulkCmd creates bulk operations commands
func createBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk",
		Short: "🔄 Bulk operations on command data",
		Long:  "Perform bulk operations like recalculating XP, updating categories, or exporting data",
	}

	var dryRun bool

	recalcXPCmd := &cobra.Command{
		Use:   "recalc-xp",
		Short: "🔄 Recalculate XP for all commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			operation := &stats.BulkOperation{
				Type:       "recalculate_xp",
				Filters:    &stats.AdvancedFilter{},
				Parameters: map[string]interface{}{},
				DryRun:     dryRun,
			}

			advancedStats := stats.NewAdvancedStatsManager(db)
			result, err := advancedStats.PerformBulkOperation(operation)
			if err != nil {
				return fmt.Errorf("bulk operation failed: %w", err)
			}

			mode := "🔄 Executed"
			if dryRun {
				mode = "🔍 Dry Run"
			}

			fmt.Printf("%s Bulk XP Recalculation\n", mode)
			fmt.Printf("Commands affected: %d\n", result.Affected)
			fmt.Printf("Duration: %v\n", result.Duration)

			if result.Details != nil {
				if details, ok := result.Details.(map[string]interface{}); ok {
					fmt.Printf("XP adjustments: %v\n", details["total_xp_adjusted"])
				}
			}

			return nil
		},
	}

	recalcXPCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")

	cmd.AddCommand(recalcXPCmd)
	return cmd
}

// createShellCmd creates shell integration management
func createShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "🐚 Manage shell integrations",
		Long:  "Install, update, and manage shell hooks for different shells (Zsh, Bash, Fish, PowerShell)",
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "📊 Show shell integration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			installer, err := shell.NewHookInstaller("")
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			shellType := installer.GetShellType()
			installed, err := installer.IsInstalled()
			if err != nil {
				return fmt.Errorf("failed to check installation: %w", err)
			}

			fmt.Printf("🐚 Shell Integration Status\n")
			fmt.Printf("===========================\n\n")
			fmt.Printf("Current Shell: %s\n", shellType)
			
			if installed {
				fmt.Printf("Status: ✅ Installed\n")
			} else {
				fmt.Printf("Status: ❌ Not Installed\n")
			}

			fmt.Printf("\nSupported Shells:\n")
			fmt.Printf("  ✅ Zsh (Z Shell)\n")
			fmt.Printf("  ✅ Bash\n")
			fmt.Printf("  ✅ Fish\n")
			fmt.Printf("  ✅ PowerShell\n")

			if !installed {
				fmt.Printf("\nTo install: termonaut shell install\n")
			}

			return nil
		},
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "⚡ Install shell integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			installer, err := shell.NewHookInstaller("")
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			shellType := installer.GetShellType()
			fmt.Printf("🐚 Installing %s shell integration...\n", shellType)

			if err := installer.Install(); err != nil {
				return fmt.Errorf("failed to install shell hook: %w", err)
			}

			fmt.Printf("✅ Successfully installed %s integration!\n", shellType)
			fmt.Printf("\nRestart your terminal or run:\n")
			
			switch shellType {
			case "zsh":
				fmt.Printf("  source ~/.zshrc\n")
			case "bash":
				fmt.Printf("  source ~/.bashrc\n")
			case "fish":
				fmt.Printf("  source ~/.config/fish/config.fish\n")
			case "powershell":
				fmt.Printf("  . $PROFILE\n")
			}

			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "🔄 Update shell integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			installer, err := shell.NewHookInstaller("")
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			fmt.Printf("🔄 Updating shell integration...\n")

			// Uninstall first
			if err := installer.Uninstall(); err != nil {
				fmt.Printf("⚠️ Warning: Failed to uninstall old hook: %v\n", err)
			}

			// Then reinstall
			if err := installer.Install(); err != nil {
				return fmt.Errorf("failed to reinstall shell hook: %w", err)
			}

			fmt.Printf("✅ Shell integration updated successfully!\n")
			return nil
		},
	}

	cmd.AddCommand(statusCmd)
	cmd.AddCommand(installCmd)
	cmd.AddCommand(updateCmd)
	return cmd
}

// createAPICmd creates API server management
func createAPICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "🌐 API server for external integrations",
		Long:  "Start and manage the REST API server for external tool integrations",
	}

	var port int
	var enableCORS bool

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "🚀 Start the API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			server := api.NewAPIServer(db, port)
			
			fmt.Printf("🚀 Starting Termonaut API Server\n")
			fmt.Printf("Port: %d\n", port)
			fmt.Printf("CORS: %v\n", enableCORS)
			fmt.Printf("\nAvailable endpoints:\n")
			fmt.Printf("  GET  /api/v1/health\n")
			fmt.Printf("  GET  /api/v1/stats\n")
			fmt.Printf("  GET  /api/v1/commands\n")
			fmt.Printf("  GET  /api/v1/categories\n")
			fmt.Printf("  POST /api/v1/bulk/operations\n")
			fmt.Printf("\nPress Ctrl+C to stop\n\n")

			return server.Start()
		},
	}

	startCmd.Flags().IntVar(&port, "port", 8080, "Port to run the API server on")
	startCmd.Flags().BoolVar(&enableCORS, "cors", true, "Enable CORS headers")

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "📋 Show API information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("🌐 Termonaut API Server v0.8.0\n")
			fmt.Printf("===============================\n\n")
			fmt.Printf("Available Endpoints:\n")
			fmt.Printf("  Health:\n")
			fmt.Printf("    GET /api/v1/health\n\n")
			fmt.Printf("  Statistics:\n")
			fmt.Printf("    GET /api/v1/stats\n")
			fmt.Printf("    GET /api/v1/stats/basic\n")
			fmt.Printf("    GET /api/v1/stats/gamification\n")
			fmt.Printf("    GET /api/v1/stats/productivity\n\n")
			fmt.Printf("  Commands:\n")
			fmt.Printf("    GET /api/v1/commands?limit=50\n")
			fmt.Printf("    POST /api/v1/commands/search\n\n")
			fmt.Printf("  Categories:\n")
			fmt.Printf("    GET /api/v1/categories\n\n")
			fmt.Printf("  Scoring Rules:\n")
			fmt.Printf("    GET /api/v1/scoring/rules\n")
			fmt.Printf("    POST /api/v1/scoring/rules\n\n")
			fmt.Printf("  Bulk Operations:\n")
			fmt.Printf("    POST /api/v1/bulk/operations\n\n")
			fmt.Printf("  Export:\n")
			fmt.Printf("    GET /api/v1/export/json\n")
			fmt.Printf("    GET /api/v1/export/csv\n\n")

			return nil
		},
	}

	cmd.AddCommand(startCmd)
	cmd.AddCommand(infoCmd)
	return cmd
}

// createAnalyticsCmd creates advanced analytics commands
func createAnalyticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "📊 Advanced analytics and insights",
		Long:  "Generate sophisticated analytics reports and insights about your terminal usage",
	}

	insightsCmd := &cobra.Command{
		Use:   "insights",
		Short: "💡 Generate usage insights",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := initAdvancedDB()
			if err != nil {
				return err
			}
			defer db.Close()

			advancedStats := stats.NewAdvancedStatsManager(db)
			filter := &stats.AdvancedFilter{
				Limit: 1000, // Analyze last 1000 commands
			}

			analytics, err := advancedStats.GetAdvancedAnalytics(filter)
			if err != nil {
				return fmt.Errorf("failed to get analytics: %w", err)
			}

			fmt.Printf("📊 Advanced Terminal Usage Analytics\n")
			fmt.Printf("====================================\n\n")

			fmt.Printf("📈 Overview:\n")
			fmt.Printf("  Total Commands Analyzed: %d\n", analytics.TotalCommands)
			
			if analytics.TimeRange != nil {
				fmt.Printf("  Time Range: %s to %s\n", 
					analytics.TimeRange.Start.Format("2006-01-02"),
					analytics.TimeRange.End.Format("2006-01-02"))
				fmt.Printf("  Analysis Period: %v\n", analytics.TimeRange.Duration)
			}

			fmt.Printf("\n💡 Recommendations:\n")
			for i, rec := range analytics.Recommendations {
				fmt.Printf("  %d. %s\n", i+1, rec)
			}

			fmt.Printf("\n🔧 Advanced features available:\n")
			fmt.Printf("  • Custom scoring rules: termonaut advanced scoring list\n")
			fmt.Printf("  • Advanced filtering: termonaut advanced filter search\n")
			fmt.Printf("  • Bulk operations: termonaut advanced bulk --help\n")
			fmt.Printf("  • API endpoints: termonaut advanced api start\n")

			return nil
		},
	}

	cmd.AddCommand(insightsCmd)
	return cmd
} 