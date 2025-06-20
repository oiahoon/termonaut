package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/api"
	"github.com/oiahoon/termonaut/internal/categories"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/github"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
			binaryPath, err := shell.GetBinaryPath()
			if err != nil {
				return fmt.Errorf("failed to get binary path: %w", err)
			}

			installer, err := shell.NewHookInstaller(binaryPath)
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			installed, err := installer.IsInstalled()
			if err != nil {
				return fmt.Errorf("failed to check installation status: %w", err)
			}

			fmt.Printf("🐚 Shell Integration Status\n")
			fmt.Printf("===========================\n\n")
			fmt.Printf("Current Shell: %s\n", installer.GetShellType())

			if installed {
				fmt.Printf("Status: ✅ Installed\n\n")
			} else {
				fmt.Printf("Status: ❌ Not Installed\n\n")
			}

			fmt.Printf("Supported Shells:\n")
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
			binaryPath, err := shell.GetBinaryPath()
			if err != nil {
				return fmt.Errorf("failed to get binary path: %w", err)
			}

			installer, err := shell.NewHookInstaller(binaryPath)
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			shellType := installer.GetShellType()
			fmt.Printf("🐚 Installing %s shell integration...\n", shellType)

			force, _ := cmd.Flags().GetBool("force")
			if err := installer.InstallWithForce(force); err != nil {
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
	installCmd.Flags().Bool("force", false, "Force reinstall even if already installed")

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "🔄 Update shell integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			binaryPath, err := shell.GetBinaryPath()
			if err != nil {
				return fmt.Errorf("failed to get binary path: %w", err)
			}

			installer, err := shell.NewHookInstaller(binaryPath)
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

	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "🗑️ Uninstall shell integration",
		RunE: func(cmd *cobra.Command, args []string) error {
			binaryPath, err := shell.GetBinaryPath()
			if err != nil {
				return fmt.Errorf("failed to get binary path: %w", err)
			}

			installer, err := shell.NewHookInstaller(binaryPath)
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			shellType := installer.GetShellType()
			fmt.Printf("🗑️ Uninstalling %s shell integration...\n", shellType)

			if err := installer.Uninstall(); err != nil {
				return fmt.Errorf("failed to uninstall shell hook: %w", err)
			}

			fmt.Printf("✅ Successfully uninstalled %s integration!\n", shellType)
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

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "🔧 Generate shell completion scripts",
		Long: `Generate shell completion scripts for termonaut commands.

This will automatically install completion for the current shell, making it easier
to use termonaut commands with tab completion.

Examples:
  termonaut advanced shell completion          # Auto-detect and install for current shell
  termonaut advanced shell completion bash     # Generate bash completion
  termonaut advanced shell completion zsh      # Generate zsh completion
  termonaut advanced shell completion fish     # Generate fish completion`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var targetShell string

			if len(args) == 0 {
				// Auto-detect current shell
				binaryPath, err := shell.GetBinaryPath()
				if err != nil {
					return fmt.Errorf("failed to get binary path: %w", err)
				}

				installer, err := shell.NewHookInstaller(binaryPath)
				if err != nil {
					return fmt.Errorf("failed to detect shell: %w", err)
				}
				targetShell = string(installer.GetShellType())
			} else {
				targetShell = args[0]
			}

			fmt.Printf("🔧 Setting up shell completion for %s...\n\n", targetShell)

			switch targetShell {
			case "bash":
				return setupBashCompletion()
			case "zsh":
				return setupZshCompletion()
			case "fish":
				return setupFishCompletion()
			case "powershell":
				return setupPowerShellCompletion()
			default:
				return fmt.Errorf("unsupported shell: %s. Supported: bash, zsh, fish, powershell", targetShell)
			}
		},
	}

	cmd.AddCommand(statusCmd)
	cmd.AddCommand(installCmd)
	cmd.AddCommand(updateCmd)
	cmd.AddCommand(uninstallCmd)
	cmd.AddCommand(completionCmd)
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

// setupBashCompletion sets up bash completion
func setupBashCompletion() error {
	fmt.Println("📝 Bash Completion Setup")
	fmt.Println("========================")
	fmt.Println()

	fmt.Println("1. Generate completion script:")
	fmt.Println("   termonaut completion bash > /usr/local/etc/bash_completion.d/termonaut")
	fmt.Println()

	fmt.Println("2. Or add to your ~/.bashrc:")
	fmt.Println("   source <(termonaut completion bash)")
	fmt.Println()

	fmt.Println("3. Reload your shell:")
	fmt.Println("   source ~/.bashrc")
	fmt.Println()

	fmt.Println("✅ After setup, you can use tab completion with termonaut commands!")
	return nil
}

// setupZshCompletion sets up zsh completion
func setupZshCompletion() error {
	fmt.Println("📝 Zsh Completion Setup")
	fmt.Println("=======================")
	fmt.Println()

	fmt.Println("1. Generate completion script:")
	fmt.Println("   termonaut completion zsh > \"${fpath[1]}/_termonaut\"")
	fmt.Println()

	fmt.Println("2. Or add to your ~/.zshrc:")
	fmt.Println("   source <(termonaut completion zsh)")
	fmt.Println()

	fmt.Println("3. For Oh My Zsh users:")
	fmt.Println("   termonaut completion zsh > ~/.oh-my-zsh/completions/_termonaut")
	fmt.Println()

	fmt.Println("4. Reload your shell:")
	fmt.Println("   source ~/.zshrc")
	fmt.Println()

	fmt.Println("✅ After setup, you can use tab completion with termonaut commands!")
	return nil
}

// setupFishCompletion sets up fish completion
func setupFishCompletion() error {
	fmt.Println("📝 Fish Completion Setup")
	fmt.Println("========================")
	fmt.Println()

	fmt.Println("1. Generate completion script:")
	fmt.Println("   termonaut completion fish > ~/.config/fish/completions/termonaut.fish")
	fmt.Println()

	fmt.Println("2. Reload fish:")
	fmt.Println("   source ~/.config/fish/config.fish")
	fmt.Println()

	fmt.Println("✅ After setup, you can use tab completion with termonaut commands!")
	return nil
}

// setupPowerShellCompletion sets up PowerShell completion
func setupPowerShellCompletion() error {
	fmt.Println("📝 PowerShell Completion Setup")
	fmt.Println("==============================")
	fmt.Println()

	fmt.Println("1. Add to your PowerShell profile:")
	fmt.Println("   termonaut completion powershell | Out-String | Invoke-Expression")
	fmt.Println()

	fmt.Println("2. Or save to a file and source it:")
	fmt.Println("   termonaut completion powershell > termonaut.ps1")
	fmt.Println("   . ./termonaut.ps1")
	fmt.Println()

	fmt.Println("3. Reload PowerShell:")
	fmt.Println("   . $PROFILE")
	fmt.Println()

	fmt.Println("✅ After setup, you can use tab completion with termonaut commands!")
	return nil
}

func runGitHubSyncNowCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.SyncEnabled {
		fmt.Println("❌ GitHub sync is disabled")
		fmt.Println("Enable it with: tn config set sync_enabled true")
		return nil
	}

	if cfg.SyncRepo == "" {
		fmt.Println("❌ No sync repository configured")
		fmt.Println("Set it with: tn config set sync_repo username/repository")
		return nil
	}

	fmt.Printf("🚀 Syncing to %s...\n", cfg.SyncRepo)

	// Initialize database and stats
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get user progress
	userProgress, err := db.GetUserProgress()
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// Initialize sync manager
	statsCalc := stats.New(db)
	syncManager := github.NewSyncManager(cfg, statsCalc)

	// Perform sync
	result, err := syncManager.SyncToRepository(userProgress)
	if err != nil {
		fmt.Printf("❌ Sync failed: %v\n", err)
		return err
	}

	if result.Success {
		fmt.Printf("✅ Sync completed successfully!\n")
		fmt.Printf("📁 Files updated: %d\n", len(result.FilesUpdated))
		fmt.Printf("🏷️  Badges updated: %d\n", result.BadgesUpdated)
		fmt.Printf("📄 Profile size: %d bytes\n", result.ProfileSize)
		fmt.Printf("⏱️  Duration: %s\n", result.SyncDuration)
		if result.CommitHash != "" {
			fmt.Printf("🔗 Commit: %s\n", result.CommitHash[:8])
		}

		if len(result.FilesUpdated) > 0 {
			fmt.Println("\n📋 Updated files:")
			for _, file := range result.FilesUpdated {
				fmt.Printf("  • %s\n", file)
			}
		}
	} else {
		fmt.Printf("❌ Sync failed: %s\n", result.ErrorMessage)
	}

	// Save sync result
	if err := saveLastSyncResult(cfg, result); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save sync result: %v\n", err)
	}

	return nil
}

// saveLastSyncResult saves the sync result to a file
func saveLastSyncResult(cfg *config.Config, result *github.SyncResult) error {
	lastSyncFile := filepath.Join(config.GetDataDir(cfg), "last_sync.json")
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(lastSyncFile, data, 0644)
}

func runGitHubSyncStatusCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("📊 GitHub Sync Status")
	fmt.Println("──────────────────────")

	if cfg.SyncEnabled {
		fmt.Println("✅ Sync: Enabled")
	} else {
		fmt.Println("❌ Sync: Disabled")
	}

	if cfg.SyncRepo != "" {
		fmt.Printf("📁 Repository: %s\n", cfg.SyncRepo)
	} else {
		fmt.Println("📁 Repository: Not configured")
	}

	fmt.Printf("⏰ Frequency: %s\n", cfg.BadgeUpdateFrequency)

	// Check last sync
	lastSyncFile := filepath.Join(config.GetDataDir(cfg), "last_sync.json")
	if data, err := os.ReadFile(lastSyncFile); err == nil {
		var lastSync github.SyncResult
		if json.Unmarshal(data, &lastSync) == nil {
			fmt.Printf("🕐 Last sync: %s\n", lastSync.Timestamp.Format("2006-01-02 15:04:05"))
			if lastSync.Success {
				fmt.Printf("✅ Status: Success (%d files, %s)\n", len(lastSync.FilesUpdated), lastSync.SyncDuration)
			} else {
				fmt.Printf("❌ Status: Failed (%s)\n", lastSync.ErrorMessage)
			}
		}
	} else {
		fmt.Println("🕐 Last sync: Never")
	}

	// Show setup instructions if not configured
	if !cfg.SyncEnabled || cfg.SyncRepo == "" {
		fmt.Println()
		fmt.Println("🔧 Setup Instructions:")
		if !cfg.SyncEnabled {
			fmt.Println("1. Enable sync: tn config set sync_enabled true")
		}
		if cfg.SyncRepo == "" {
			fmt.Println("2. Set repository: tn config set sync_repo username/repository")
		}
		fmt.Println("3. Run setup: tn github sync setup")
		fmt.Println("4. Test sync: tn github sync now")
	}

	return nil
}

func runGitHubSyncSetupCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("🔧 GitHub Sync Setup")
	fmt.Println("═══════════════════")
	fmt.Println()

	// Load current configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Interactive setup
	fmt.Println("This will help you set up GitHub synchronization for your Termonaut data.")
	fmt.Println()

	// Step 1: Enable sync
	if !cfg.SyncEnabled {
		fmt.Println("1️⃣  Enabling GitHub sync...")
		cfg.SyncEnabled = true
		fmt.Println("✅ GitHub sync enabled")
	} else {
		fmt.Println("1️⃣  GitHub sync is already enabled")
	}

	// Step 2: Repository configuration
	fmt.Println()
	fmt.Println("2️⃣  Repository Configuration")
	if cfg.SyncRepo == "" {
		fmt.Print("Enter your GitHub repository (username/repository): ")
		var repo string
		fmt.Scanln(&repo)
		if repo != "" {
			cfg.SyncRepo = repo
			fmt.Printf("✅ Repository set to: %s\n", repo)
		}
	} else {
		fmt.Printf("Current repository: %s\n", cfg.SyncRepo)
		fmt.Print("Change repository? (y/N): ")
		var change string
		fmt.Scanln(&change)
		if strings.ToLower(change) == "y" {
			fmt.Print("Enter new repository (username/repository): ")
			var repo string
			fmt.Scanln(&repo)
			if repo != "" {
				cfg.SyncRepo = repo
				fmt.Printf("✅ Repository updated to: %s\n", repo)
			}
		}
	}

	// Step 3: Frequency configuration
	fmt.Println()
	fmt.Println("3️⃣  Sync Frequency")
	fmt.Printf("Current frequency: %s\n", cfg.BadgeUpdateFrequency)
	fmt.Println("Available options: hourly, daily, weekly")
	fmt.Print("Change frequency? (y/N): ")
	var changeFreq string
	fmt.Scanln(&changeFreq)
	if strings.ToLower(changeFreq) == "y" {
		fmt.Print("Enter frequency (hourly/daily/weekly): ")
		var freq string
		fmt.Scanln(&freq)
		if freq == "hourly" || freq == "daily" || freq == "weekly" {
			cfg.BadgeUpdateFrequency = freq
			fmt.Printf("✅ Frequency set to: %s\n", freq)
		}
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("4️⃣  GitHub Actions Setup (Optional)")
	fmt.Println("To automate updates, you can set up GitHub Actions:")
	fmt.Println()
	fmt.Printf("1. Generate workflow: tn github actions generate termonaut-stats-update\n")
	fmt.Printf("2. Commit to your repository: %s\n", cfg.SyncRepo)
	fmt.Printf("3. The workflow will update badges every 6 hours\n")
	fmt.Println()

	fmt.Println("5️⃣  Test Your Setup")
	fmt.Println("Run a test sync to verify everything works:")
	fmt.Println("tn github sync now")
	fmt.Println()

	fmt.Println("✅ Setup complete!")
	fmt.Println()
	fmt.Println("📋 Next Steps:")
	fmt.Println("• Test sync: tn github sync now")
	fmt.Println("• Check status: tn github sync status")
	fmt.Println("• Set up automation: tn github actions generate termonaut-stats-update")

	return nil
}

func runGitHubActionsTriggerCommand(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.SyncRepo == "" {
		fmt.Println("❌ No sync repository configured")
		fmt.Println("Set it with: tn config set sync_repo username/repository")
		return nil
	}

	fmt.Printf("🚀 Triggering workflow '%s' in %s...\n", workflowName, cfg.SyncRepo)
	fmt.Println("✅ GitHub Actions trigger feature coming soon!")
	fmt.Println("📋 For now, manually trigger workflows in GitHub:")
	fmt.Printf("🔗 https://github.com/%s/actions\n", cfg.SyncRepo)

	return nil
}
