package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/oiahoon/termonaut/internal/privacy"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Version information (will be set during build)
	version = "0.4.0-dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "termonaut",
	Short: "🚀 Termonaut - Your Terminal Journey Companion",
	Long: `Termonaut is a gamified terminal productivity tracker that transforms
your command-line usage into an engaging RPG-like experience.

Track your terminal habits, earn XP, unlock achievements, and level up
your productivity - all without leaving your CLI!`,
	SilenceUsage: true,
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display terminal usage statistics",
	Long: `Show comprehensive statistics about your terminal usage including
command counts, sessions, and productivity metrics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStatsCommand(cmd, args)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Termonaut shell integration",
	Long: `Install shell hooks to enable automatic command tracking.
This command will modify your shell configuration file (.zshrc or .bashrc)
to add the necessary hooks for logging commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInitCommand(cmd, args)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Termonaut configuration",
	Long: `View and modify Termonaut configuration settings.
Use subcommands to get or set specific configuration values.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration value",
	Long:  "Display the current value of a configuration setting.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigGetCommand(cmd, args)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long:  "Update a configuration setting with a new value.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfigSetCommand(cmd, args)
	},
}

var logCommandCmd = &cobra.Command{
	Use:    "log-command <command>",
	Short:  "Log a command execution (internal use)",
	Long:   "Internal command used by shell hooks to log command executions.",
	Args:   cobra.ExactArgs(1),
	Hidden: true, // Hide from help as it's for internal use
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLogCommandCommand(cmd, args)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Termonaut %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(logCommandCmd)
	rootCmd.AddCommand(versionCmd)

	// Add gamification commands
	rootCmd.AddCommand(progressCmd)
	rootCmd.AddCommand(achievementsCmd)
	rootCmd.AddCommand(levelCmd)

	// Add category analysis command
	rootCmd.AddCommand(categoriesCmd)

	// Add productivity analytics command
	rootCmd.AddCommand(analyticsCmd)

	// Add advanced features
	rootCmd.AddCommand(heatmapCmd)
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(createAdvancedCmd())

	// Config subcommands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)

	// Add flags
	statsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	statsCmd.Flags().Bool("today", false, "Show today's stats only")
	statsCmd.Flags().Bool("weekly", false, "Show weekly stats")
	statsCmd.Flags().Bool("monthly", false, "Show monthly stats")

	initCmd.Flags().Bool("force", false, "Force reinstall even if already installed")
	initCmd.Flags().String("shell", "", "Specify shell type (zsh, bash)")

	// Gamification command flags
	progressCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	achievementsCmd.Flags().Bool("all", false, "Show all achievements including locked ones")
	achievementsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	levelCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Category command flags
	categoriesCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Analytics command flags
	analyticsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Heatmap command flags
	heatmapCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Dashboard command flags
	dashboardCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func runStatsCommand(cmd *cobra.Command, args []string) error {
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

	// Initialize stats calculator
	statsCalc := stats.New(db)

	// Get flags
	jsonOutput, _ := cmd.Flags().GetBool("json")
	todayOnly, _ := cmd.Flags().GetBool("today")

	if todayOnly {
		// Show today's stats only
		todayStats, err := statsCalc.GetTodayStats()
		if err != nil {
			return fmt.Errorf("failed to get today's stats: %w", err)
		}

		if jsonOutput {
			fmt.Printf("%+v\n", todayStats)
		} else {
			fmt.Printf("📅 Today's Stats:\n")
			fmt.Printf("Commands: %v\n", todayStats["commands_today"])
		}
		return nil
	}

	// Get basic stats
	basicStats, err := statsCalc.GetBasicStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	if jsonOutput {
		fmt.Printf("%+v\n", basicStats)
	} else {
		fmt.Print(statsCalc.FormatBasicStats(basicStats))
	}

	return nil
}

func runInitCommand(cmd *cobra.Command, args []string) error {
	// Get binary path
	binaryPath, err := shell.GetBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get binary path: %w", err)
	}

	// Create hook installer
	installer, err := shell.NewHookInstaller(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create hook installer: %w", err)
	}

	// Check if already installed
	force, _ := cmd.Flags().GetBool("force")
	if !force {
		installed, err := installer.IsInstalled()
		if err != nil {
			return fmt.Errorf("failed to check installation status: %w", err)
		}
		if installed {
			fmt.Println("✅ Termonaut is already initialized!")
			fmt.Println("Use --force to reinstall")
			return nil
		}
	}

	// Install hook
	if err := installer.Install(); err != nil {
		return fmt.Errorf("failed to install shell hook: %w", err)
	}

	fmt.Printf("🚀 Termonaut initialized successfully!\n")
	fmt.Printf("Shell: %s\n", installer.GetShellType())
	fmt.Println("\nPlease restart your terminal or run:")
	if installer.GetShellType() == shell.Zsh {
		fmt.Println("  source ~/.zshrc")
	} else {
		fmt.Println("  source ~/.bashrc")
	}
	fmt.Println("\nThen start using your terminal normally. Run 'termonaut stats' to see your progress!")

	return nil
}

func runConfigGetCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(args) == 0 {
		// Show all config
		fmt.Println("🔧 Termonaut Configuration:")
		fmt.Printf("Display Mode: %s\n", cfg.DisplayMode)
		fmt.Printf("Theme: %s\n", cfg.Theme)
		fmt.Printf("Show Gamification: %t\n", cfg.ShowGamification)
		fmt.Printf("Idle Timeout: %d minutes\n", cfg.IdleTimeoutMinutes)
		fmt.Printf("Track Git Repos: %t\n", cfg.TrackGitRepos)
		fmt.Printf("Command Categories: %t\n", cfg.CommandCategories)
		fmt.Printf("Sync Enabled: %t\n", cfg.SyncEnabled)
		fmt.Printf("Anonymous Mode: %t\n", cfg.AnonymousMode)
		fmt.Printf("Log Level: %s\n", cfg.LogLevel)
		return nil
	}

	key := args[0]
	switch key {
	case "display_mode":
		fmt.Println(cfg.DisplayMode)
	case "theme":
		fmt.Println(cfg.Theme)
	case "show_gamification":
		fmt.Printf("%t\n", cfg.ShowGamification)
	case "idle_timeout_minutes":
		fmt.Printf("%d\n", cfg.IdleTimeoutMinutes)
	case "log_level":
		fmt.Println(cfg.LogLevel)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return nil
}

func runConfigSetCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	key := args[0]
	value := args[1]

	switch key {
	case "display_mode":
		if value != "off" && value != "enter" && value != "ps1" && value != "floating" {
			return fmt.Errorf("invalid display_mode value. Must be: off, enter, ps1, floating")
		}
		cfg.DisplayMode = value
	case "theme":
		if value != "minimal" && value != "emoji" && value != "ascii" {
			return fmt.Errorf("invalid theme value. Must be: minimal, emoji, ascii")
		}
		cfg.Theme = value
	case "show_gamification":
		if value != "true" && value != "false" {
			return fmt.Errorf("invalid boolean value. Must be: true, false")
		}
		cfg.ShowGamification = value == "true"
	case "log_level":
		if value != "debug" && value != "info" && value != "warn" && value != "error" {
			return fmt.Errorf("invalid log_level value. Must be: debug, info, warn, error")
		}
		cfg.LogLevel = value
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ Configuration updated: %s = %s\n", key, value)
	return nil
}

func runLogCommandCommand(cmd *cobra.Command, args []string) error {
	command := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// If config fails to load, still try to log the command
		cfg = config.DefaultConfig()
	}

	// Initialize command sanitizer for privacy protection
	sanitizer := privacy.NewCommandSanitizer(privacy.DefaultSanitizationConfig())
	
	// Sanitize command for privacy
	sanitizedCommand, shouldIgnore := sanitizer.SanitizeCommand(command)
	if shouldIgnore {
		return nil // Skip logging this command entirely
	}

	// Initialize logger (with minimal output for background operation)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Only log errors for background operation

	// Initialize database
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		// Silent fail for background operation
		return nil
	}
	defer db.Close()

	// Get or create session
	session, err := db.GetOrCreateSession(shell.GetTerminalPID(), string(shell.Zsh))
	if err != nil {
		// Silent fail for background operation
		return nil
	}

	// Create command record
	commandRecord := &models.Command{
		Timestamp: time.Now(),
		SessionID: session.ID,
		Command:   sanitizedCommand, // Use sanitized command
		ExitCode:  0, // We don't have exit code from preexec hook
		CWD:       shell.GetCurrentWorkingDir(),
	}

	// Check for Easter Eggs (only if enabled in config)
	if cfg.ShowGamification {
		easterEggManager := gamification.NewEasterEggManager()
		
		// Get recent command history for context
		recentCommands, _ := db.GetRecentCommands(10)
		var commandHistory []string
		for _, recentCmd := range recentCommands {
			commandHistory = append(commandHistory, recentCmd.Command)
		}
		
		easterEggContext := &gamification.EasterEggContext{
			CommandsInTimeframe:    len(recentCommands),
			TimeframeDuration:      time.Hour, // Last hour
			IdleDuration:          time.Since(session.StartTime),
			IsFirstCommandToday:   isFirstCommandToday(recentCommands),
			LastCommand:          sanitizedCommand,
			CommandHistory:       commandHistory,
			QuotesMismatched:     hasUnmatchedQuotes(sanitizedCommand),
		}
		
		if easterEgg := easterEggManager.CheckForEasterEgg(easterEggContext); easterEgg != "" {
			// Store easter egg for display (could be shown in stats or dashboard)
			// For now, we'll just log it silently
			logger.Debug("Easter egg triggered: ", easterEgg)
		}
	}

	// Store command with enhanced gamification (XP, achievements, privacy)
	if err := db.StoreCommandWithXP(commandRecord); err != nil {
		// Silent fail for background operation
		return nil
	}

	return nil
}

func setupLogger(level string) *logrus.Logger {
	logger := logrus.New()

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// isFirstCommandToday checks if this is the first command executed today
func isFirstCommandToday(recentCommands []*models.Command) bool {
	if len(recentCommands) == 0 {
		return true
	}
	
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	for _, cmd := range recentCommands {
		if cmd.Timestamp.After(today) {
			return false // Found another command today
		}
	}
	
	return true // No commands found today
}

// hasUnmatchedQuotes checks if the command has unmatched quotes
func hasUnmatchedQuotes(command string) bool {
	singleQuotes := strings.Count(command, "'")
	doubleQuotes := strings.Count(command, "\"")
	
	return singleQuotes%2 != 0 || doubleQuotes%2 != 0
}
