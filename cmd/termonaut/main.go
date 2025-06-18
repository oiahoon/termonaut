package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	version = "0.9.0"
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
	Aliases:      []string{"tn"},
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

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Generate shell prompt integration",
	Long: `Generate a compact status for shell prompt integration.

This command outputs a brief summary suitable for shell prompts.
Add to your shell configuration:

Bash/Zsh:
  export PS1="$(termonaut prompt) $PS1"

Fish:
  function fish_prompt
      echo (termonaut prompt) (fish_prompt)
  end`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			// Silent fail for prompt integration
			return nil
		}

		// Initialize logger (silent)
		logger := setupLogger("error")

		// Initialize database
		db, err := database.New(config.GetDataDir(cfg), logger)
		if err != nil {
			// Silent fail for prompt integration
			return nil
		}
		defer db.Close()

		// Initialize stats calculator
		statsCalc := stats.New(db)

		// Get basic stats
		basicStats, err := statsCalc.GetBasicStats()
		if err != nil {
			// Silent fail for prompt integration
			return nil
		}

		// Get user progress
		userProgress, err := db.GetUserProgress()
		if err != nil {
			// Silent fail for prompt integration
			return nil
		}

		// Generate compact prompt
		if cfg.Theme == "minimal" {
			fmt.Printf("[L%d %dc]", userProgress.CurrentLevel, basicStats.TotalCommands)
		} else {
			fmt.Printf("🚀L%d(%dc)", userProgress.CurrentLevel, basicStats.TotalCommands)
		}

		return nil
	},
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(logCommandCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(promptCmd)

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
	rootCmd.AddCommand(createGitHubCmd())

	// Add completion command
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(termonaut completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ termonaut completion bash > /etc/bash_completion.d/termonaut
  # macOS:
  $ termonaut completion bash > /usr/local/etc/bash_completion.d/termonaut

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ termonaut completion zsh > "${fpath[1]}/_termonaut"

  # You will need to start a new shell for this setup to take effect.

fish:
  $ termonaut completion fish | source

  # To load completions for each session, execute once:
  $ termonaut completion fish > ~/.config/fish/completions/termonaut.fish

PowerShell:
  PS> termonaut completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> termonaut completion powershell > termonaut.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	rootCmd.AddCommand(completionCmd)

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
	if err := installer.InstallWithForce(force); err != nil {
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
		fmt.Printf("Easter Eggs Enabled: %t\n", cfg.EasterEggsEnabled)
		fmt.Printf("Empty Command Stats: %t\n", cfg.EmptyCommandStats)
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
	case "easter_eggs_enabled":
		fmt.Printf("%t\n", cfg.EasterEggsEnabled)
	case "empty_command_stats":
		fmt.Printf("%t\n", cfg.EmptyCommandStats)
	case "idle_timeout_minutes":
		fmt.Printf("%d\n", cfg.IdleTimeoutMinutes)
	case "sync_enabled":
		fmt.Printf("%t\n", cfg.SyncEnabled)
	case "sync_repo":
		fmt.Println(cfg.SyncRepo)
	case "badge_update_frequency":
		fmt.Println(cfg.BadgeUpdateFrequency)
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
	case "easter_eggs_enabled":
		if value != "true" && value != "false" {
			return fmt.Errorf("invalid boolean value. Must be: true, false")
		}
		cfg.EasterEggsEnabled = value == "true"
	case "empty_command_stats":
		if value != "true" && value != "false" {
			return fmt.Errorf("invalid boolean value. Must be: true, false")
		}
		cfg.EmptyCommandStats = value == "true"
	case "sync_enabled":
		if value != "true" && value != "false" {
			return fmt.Errorf("invalid boolean value. Must be: true, false")
		}
		cfg.SyncEnabled = value == "true"
	case "sync_repo":
		cfg.SyncRepo = value
	case "badge_update_frequency":
		if value != "hourly" && value != "daily" && value != "weekly" {
			return fmt.Errorf("invalid badge_update_frequency value. Must be: hourly, daily, weekly")
		}
		cfg.BadgeUpdateFrequency = value
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
	// Debug: Check what we received
	if len(args) == 0 {
		// No command provided, treat as empty command
		cfg, err := config.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}

		if cfg.EmptyCommandStats {
			return showQuickStats()
		} else {
			return nil
		}
	}

	command := args[0]

	// Check for empty command (just Enter was pressed)
	trimmedCommand := strings.TrimSpace(command)
	if trimmedCommand == "" {
		// Load configuration to check if feature is enabled
		cfg, err := config.Load()
		if err != nil {
			cfg = config.DefaultConfig()
		}

		if cfg.EmptyCommandStats {
			return showQuickStats()
		} else {
			return nil // Skip if feature is disabled
		}
	}

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
		ExitCode:  0,                // We don't have exit code from preexec hook
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
			CommandsInTimeframe: len(recentCommands),
			TimeframeDuration:   time.Hour, // Last hour
			IdleDuration:        time.Since(session.StartTime),
			IsFirstCommandToday: isFirstCommandToday(recentCommands),
			LastCommand:         sanitizedCommand,
			CommandHistory:      commandHistory,
			QuotesMismatched:    hasUnmatchedQuotes(sanitizedCommand),
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

// showQuickStats displays a concise stats summary when empty command is executed
func showQuickStats() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Check if feature is enabled
	if !cfg.EmptyCommandStats {
		return nil // Feature disabled, do nothing
	}

	// Initialize logger
	logger := setupLogger(cfg.LogLevel)

	// Initialize database
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return nil // Silent fail for background operation
	}
	defer db.Close()

	// Get user progress
	userProgress, err := db.GetUserProgress()
	if err != nil {
		return nil // Silent fail
	}

	// Get basic stats
	statsCalc := stats.New(db)
	basicStats, err := statsCalc.GetBasicStats()
	if err != nil {
		return nil // Silent fail
	}

	// Choose display format based on configuration
	var output string
	switch cfg.DisplayMode {
	case "off":
		return nil // Don't show anything if display is off
	case "enter", "ps1", "floating":
		// Show different levels of detail based on theme
		if cfg.Theme == "minimal" {
			output = formatMinimalQuickStats(basicStats, userProgress)
		} else {
			output = formatRichQuickStats(basicStats, userProgress, cfg.Theme == "emoji")
		}
	default:
		output = formatMinimalQuickStats(basicStats, userProgress)
	}

	// Output the stats
	fmt.Print(output)
	return nil
}

// formatMinimalQuickStats formats a minimal one-line stats display
func formatMinimalQuickStats(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	return fmt.Sprintf("📊 Lv.%d | %d cmds | %d streak | %d XP\n",
		userProgress.CurrentLevel,
		basicStats.TotalCommands,
		userProgress.CurrentStreak,
		userProgress.TotalXP)
}

// formatRichQuickStats formats a rich multi-line stats display
func formatRichQuickStats(basicStats *stats.BasicStats, userProgress *models.UserProgress, useEmojis bool) string {
	var sb strings.Builder

	if useEmojis {
		sb.WriteString("🚀 ")
	}
	sb.WriteString(fmt.Sprintf("Level %d", userProgress.CurrentLevel))

	// Level progress bar
	if useEmojis {
		currentLevelXP := userProgress.CurrentLevel * userProgress.CurrentLevel * 100
		nextLevelXP := (userProgress.CurrentLevel + 1) * (userProgress.CurrentLevel + 1) * 100
		progressXP := userProgress.TotalXP - currentLevelXP
		neededXP := nextLevelXP - currentLevelXP

		if neededXP > 0 {
			progress := float64(progressXP) / float64(neededXP)
			barLength := int(progress * 8)

			sb.WriteString(" [")
			for i := 0; i < 8; i++ {
				if i < barLength {
					sb.WriteString("█")
				} else {
					sb.WriteString("░")
				}
			}
			sb.WriteString(fmt.Sprintf("] %d XP", userProgress.TotalXP))
		}
	} else {
		sb.WriteString(fmt.Sprintf(" (%d XP)", userProgress.TotalXP))
	}

	sb.WriteString("\n")

	// Commands and streak info
	if useEmojis {
		sb.WriteString(fmt.Sprintf("🎯 %d commands today", basicStats.CommandsToday))
		if userProgress.CurrentStreak > 0 {
			streakEmoji := "✨"
			if userProgress.CurrentStreak >= 7 {
				streakEmoji = "🔥"
			}
			sb.WriteString(fmt.Sprintf(" | %s %d day streak", streakEmoji, userProgress.CurrentStreak))
		}
	} else {
		sb.WriteString(fmt.Sprintf("%d commands today", basicStats.CommandsToday))
		if userProgress.CurrentStreak > 0 {
			sb.WriteString(fmt.Sprintf(" | %d day streak", userProgress.CurrentStreak))
		}
	}

	sb.WriteString("\n")

	// Most used command
	if basicStats.MostUsedCommand != "" && useEmojis {
		sb.WriteString(fmt.Sprintf("👑 %s (%dx)\n", basicStats.MostUsedCommand, basicStats.MostUsedCount))
	}

	return sb.String()
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

// createGitHubCmd creates the GitHub command and its subcommands
func createGitHubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "GitHub integration commands",
		Long:  "GitHub integration for badges, profiles, and social sharing.",
	}

	// Badges subcommand
	githubBadgesCmd := &cobra.Command{
		Use:   "badges",
		Short: "Generate GitHub badges",
		Long:  "Generate dynamic badges for your GitHub profile.",
	}

	githubBadgesGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate badges",
		Long:  "Generate badges showing your terminal stats.",
		RunE:  runGitHubBadgesGenerateCommand,
	}

	// Profile subcommand
	githubProfileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Generate GitHub profile",
		Long:  "Generate a comprehensive profile for your GitHub README.",
	}

	githubProfileGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate profile",
		Long:  "Generate a profile markdown for your GitHub README.",
		RunE:  runGitHubProfileGenerateCommand,
	}

	// Actions subcommand
	githubActionsCmd := &cobra.Command{
		Use:   "actions",
		Short: "GitHub Actions workflows",
		Long:  "Generate and manage GitHub Actions workflows for automation.",
	}

	githubActionsGenerateCmd := &cobra.Command{
		Use:   "generate [workflow-name]",
		Short: "Generate workflow file",
		Long:  "Generate a GitHub Actions workflow file for automation.",
		Args:  cobra.ExactArgs(1),
		RunE:  runGitHubActionsGenerateCommand,
	}

	githubActionsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available workflows",
		Long:  "List all available GitHub Actions workflow templates.",
		RunE:  runGitHubActionsListCommand,
	}

	// Build command hierarchy
	githubBadgesCmd.AddCommand(githubBadgesGenerateCmd)
	githubProfileCmd.AddCommand(githubProfileGenerateCmd)
	githubActionsCmd.AddCommand(githubActionsGenerateCmd)
	githubActionsCmd.AddCommand(githubActionsListCmd)

	githubCmd.AddCommand(githubBadgesCmd)
	githubCmd.AddCommand(githubProfileCmd)
	githubCmd.AddCommand(githubActionsCmd)

	// Add flags
	githubBadgesGenerateCmd.Flags().String("format", "url", "Output format (url, json, markdown)")
	githubBadgesGenerateCmd.Flags().String("output", "", "Output file path")
	githubProfileGenerateCmd.Flags().String("format", "markdown", "Output format (markdown, json)")
	githubProfileGenerateCmd.Flags().String("output", "", "Output file path")

	return githubCmd
}

func runGitHubBadgesGenerateCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize components
	logger := setupLogger(cfg.LogLevel)
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get stats and user progress
	statsCalc := stats.New(db)
	basicStats, err := statsCalc.GetBasicStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	userProgress, err := db.GetUserProgress()
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// Generate badges
	badges := generateBadgeURLs(basicStats, userProgress)

	// Get flags
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	// Format output
	var result string
	switch format {
	case "json":
		result = formatBadgesJSON(badges)
	case "markdown":
		result = formatBadgesMarkdown(badges)
	default:
		result = formatBadgesURL(badges)
	}

	// Output or save
	if output != "" {
		return os.WriteFile(output, []byte(result), 0644)
	}

	fmt.Print(result)
	return nil
}

func runGitHubProfileGenerateCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize components
	logger := setupLogger(cfg.LogLevel)
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get stats and user progress
	statsCalc := stats.New(db)
	basicStats, err := statsCalc.GetBasicStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	userProgress, err := db.GetUserProgress()
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// Generate profile
	profile := generateProfileMarkdown(basicStats, userProgress)

	// Get flags
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	// Format output
	var result string
	switch format {
	case "json":
		result = formatProfileJSON(basicStats, userProgress)
	default:
		result = profile
	}

	// Output or save
	if output != "" {
		return os.WriteFile(output, []byte(result), 0644)
	}

	fmt.Print(result)
	return nil
}

func runGitHubActionsGenerateCommand(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Available workflows
	workflows := map[string]string{
		"termonaut-stats-update":  generateStatsUpdateWorkflow(),
		"termonaut-profile-sync":  generateProfileSyncWorkflow(),
		"termonaut-weekly-report": generateWeeklyReportWorkflow(),
	}

	workflow, exists := workflows[workflowName]
	if !exists {
		return fmt.Errorf("unknown workflow: %s", workflowName)
	}

	// Create .github/workflows directory
	workflowDir := ".github/workflows"
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Write workflow file
	workflowFile := filepath.Join(workflowDir, workflowName+".yml")
	if err := os.WriteFile(workflowFile, []byte(workflow), 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	fmt.Printf("✅ Generated workflow: %s\n", workflowFile)
	fmt.Println()
	fmt.Println("📋 Next steps:")
	fmt.Println("1. Commit the workflow file to your repository")
	fmt.Println("2. Configure any required secrets in GitHub repository settings")
	fmt.Println("3. The workflow will run automatically based on its triggers")

	return nil
}

func runGitHubActionsListCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("📋 Available GitHub Actions Workflows:")
	fmt.Println()

	workflows := []struct {
		Name        string
		Description string
	}{
		{"termonaut-stats-update", "Automatically update Termonaut badges and stats"},
		{"termonaut-profile-sync", "Sync Termonaut profile data to repository"},
		{"termonaut-weekly-report", "Generate weekly productivity reports"},
	}

	for _, workflow := range workflows {
		fmt.Printf("🔧 %s\n", workflow.Name)
		fmt.Printf("   %s\n", workflow.Description)
		fmt.Println()
	}

	fmt.Println("💡 Generate a workflow with:")
	fmt.Println("   tn github actions generate [workflow-name]")

	return nil
}

// Helper functions for badge generation
func generateBadgeURLs(basicStats *stats.BasicStats, userProgress *models.UserProgress) map[string]string {
	badges := make(map[string]string)

	// Commands badge
	commandsColor := "lightgrey"
	if basicStats.TotalCommands >= 1000 {
		commandsColor = "brightgreen"
	} else if basicStats.TotalCommands >= 500 {
		commandsColor = "green"
	} else if basicStats.TotalCommands >= 100 {
		commandsColor = "yellow"
	}
	badges["Commands"] = fmt.Sprintf("https://img.shields.io/badge/Commands-%d-%s?style=flat-square&logo=terminal&logoColor=white",
		basicStats.TotalCommands, commandsColor)

	// Level badge
	levelColor := "lightgrey"
	if userProgress.CurrentLevel >= 25 {
		levelColor = "purple"
	} else if userProgress.CurrentLevel >= 10 {
		levelColor = "blue"
	} else if userProgress.CurrentLevel >= 5 {
		levelColor = "green"
	}
	badges["Level"] = fmt.Sprintf("https://img.shields.io/badge/Level-%d-%s?style=flat-square&logo=terminal&logoColor=white",
		userProgress.CurrentLevel, levelColor)

	// Streak badge
	streakColor := "red"
	if userProgress.CurrentStreak >= 30 {
		streakColor = "purple"
	} else if userProgress.CurrentStreak >= 7 {
		streakColor = "green"
	} else if userProgress.CurrentStreak >= 3 {
		streakColor = "yellow"
	}
	badges["Streak"] = fmt.Sprintf("https://img.shields.io/badge/Streak-%d%%2Bdays-%s?style=flat-square&logo=terminal&logoColor=white",
		userProgress.CurrentStreak, streakColor)

	// XP badge
	badges["XP"] = fmt.Sprintf("https://img.shields.io/badge/XP-Level%%20%d%%20%%28%d%%29-lightgrey?style=flat-square&logo=terminal&logoColor=white",
		userProgress.CurrentLevel, userProgress.TotalXP)

	return badges
}

func formatBadgesURL(badges map[string]string) string {
	var result strings.Builder
	for label, url := range badges {
		result.WriteString(fmt.Sprintf("%s: %s\n", label, url))
	}
	return result.String()
}

func formatBadgesJSON(badges map[string]string) string {
	data, _ := json.MarshalIndent(badges, "", "  ")
	return string(data)
}

func formatBadgesMarkdown(badges map[string]string) string {
	var result strings.Builder
	for label, url := range badges {
		result.WriteString(fmt.Sprintf("![%s](%s) ", label, url))
	}
	result.WriteString("\n")
	return result.String()
}

func generateProfileMarkdown(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	var sb strings.Builder

	sb.WriteString("# 🚀 My Termonaut Profile\n\n")
	sb.WriteString("*Gamified terminal productivity tracking*\n\n")

	// Stats badges
	badges := generateBadgeURLs(basicStats, userProgress)
	sb.WriteString("## 📊 Stats\n\n")
	for label, url := range badges {
		sb.WriteString(fmt.Sprintf("![%s](%s) ", label, url))
	}
	sb.WriteString("\n\n")

	// Overview
	sb.WriteString("## 📈 Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **Level**: %d (XP: %d)\n", userProgress.CurrentLevel, userProgress.TotalXP))
	sb.WriteString(fmt.Sprintf("- **Total Commands**: %d\n", basicStats.TotalCommands))
	sb.WriteString(fmt.Sprintf("- **Unique Commands**: %d\n", basicStats.UniqueCommands))
	sb.WriteString(fmt.Sprintf("- **Current Streak**: %d days\n", userProgress.CurrentStreak))
	sb.WriteString(fmt.Sprintf("- **Commands Today**: %d\n", basicStats.CommandsToday))

	if basicStats.MostUsedCommand != "" {
		sb.WriteString(fmt.Sprintf("- **Favorite Command**: `%s` (%d times)\n",
			basicStats.MostUsedCommand, basicStats.MostUsedCount))
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString("*Generated by [Termonaut](https://github.com/oiahoon/termonaut) - Terminal productivity tracker*\n")
	sb.WriteString(fmt.Sprintf("*Last updated: %s*\n", time.Now().Format("January 2, 2006")))

	return sb.String()
}

func formatProfileJSON(basicStats *stats.BasicStats, userProgress *models.UserProgress) string {
	data := map[string]interface{}{
		"basic_stats":   basicStats,
		"user_progress": userProgress,
		"generated_at":  time.Now(),
	}
	result, _ := json.MarshalIndent(data, "", "  ")
	return string(result)
}

func generateStatsUpdateWorkflow() string {
	return `name: Update Termonaut Stats

on:
  schedule:
    # Run every 6 hours
    - cron: '0 */6 * * *'
  workflow_dispatch:

jobs:
  update-stats:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/oiahoon/termonaut/cmd/termonaut@latest

    - name: Generate Badge Data
      run: |
        mkdir -p badges
        # Generate placeholder badges for now
        echo '{"Commands":"https://img.shields.io/badge/Commands-0-lightgrey?style=flat-square&logo=terminal&logoColor=white"}' > badges/badges.json

    - name: Commit and push changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add badges/ || true
        git diff --staged --quiet || git commit -m "🚀 Update Termonaut stats - $(date)"
        git push
`
}

func generateProfileSyncWorkflow() string {
	return `name: Sync Termonaut Profile

on:
  workflow_dispatch:

jobs:
  sync-profile:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/oiahoon/termonaut/cmd/termonaut@latest

    - name: Generate Profile
      run: |
        mkdir -p profile
        # Generate placeholder profile for now
        echo "# Termonaut Profile" > profile/README.md

    - name: Commit changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add profile/ || true
        git diff --staged --quiet || git commit -m "📊 Sync Termonaut profile - $(date)"
        git push
`
}

func generateWeeklyReportWorkflow() string {
	return `name: Weekly Termonaut Report

on:
  schedule:
    # Run every Monday at 9 AM UTC
    - cron: '0 9 * * 1'
  workflow_dispatch:

jobs:
  weekly-report:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Termonaut
      run: |
        go install github.com/oiahoon/termonaut/cmd/termonaut@latest

    - name: Generate Weekly Report
      run: |
        mkdir -p reports
        WEEK=$(date +'%Y-W%U')
        echo "# Weekly Report $WEEK" > reports/week-$WEEK.md

    - name: Commit reports
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add reports/ || true
        git diff --staged --quiet || git commit -m "📈 Weekly Termonaut report - $(date +'%Y-W%U')"
        git push
`
}
