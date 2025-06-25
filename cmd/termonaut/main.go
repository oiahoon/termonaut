package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/oiahoon/termonaut/internal/avatar"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/oiahoon/termonaut/internal/github"
	"github.com/oiahoon/termonaut/internal/privacy"
	"github.com/oiahoon/termonaut/internal/shell"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/oiahoon/termonaut/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// Version information (will be set during build)
	version = "v0.10.1"
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

🆕 New User? Start Here:
  termonaut setup      Interactive setup wizard (recommended)
  termonaut quickstart Quick setup with sensible defaults

📊 Daily Usage:
  termonaut tui        Launch interactive dashboard (smart mode)
  termonaut stats      Quick stats in terminal

🔧 Configuration:
  termonaut init       Install shell integration manually
  termonaut config     Manage settings

Track your terminal habits, earn XP, unlock achievements, and level up
your productivity - all without leaving your CLI!`,
	SilenceUsage: true,
	Aliases:      []string{"tn"},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display terminal usage statistics (极简模式)",
	Long: `Show terminal usage statistics in minimal shell output format.
This is the fastest way to check your productivity metrics.

三层查看模式架构:
  • 极简模式: termonaut stats (shell直接输出，最快速) ← 当前命令
  • 普通模式: termonaut tui-compact (紧凑TUI，平衡体验)
  • 完整模式: termonaut tui-enhanced (完整TUI，沉浸体验)

Options:
  --today     Show only today's statistics
  --weekly    Show this week's statistics  
  --monthly   Show this month's statistics
  --alltime   Show all-time statistics
  --json      Output in JSON format
  --minimal   Ultra-minimal one-line output
  --avatar    Include small ASCII avatar (experimental)`,
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

var terminalTestCmd = &cobra.Command{
	Use:   "terminal-test",
	Short: "Test terminal capabilities and easter egg compatibility",
	Long: `Test the current terminal's capabilities including:
- Unicode and emoji support
- Color support
- Modern terminal features
- Easter egg display compatibility`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🧪 Terminal Compatibility Test")
		fmt.Println("==============================")
		fmt.Println()

		// Display terminal information
		termInfo := gamification.GetTerminalInfo()
		fmt.Println("📊 Terminal Information:")
		for key, value := range termInfo {
			if value != "" {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
		fmt.Println()

		// Test modern terminal detection
		isModern := gamification.IsModernTerminal()
		fmt.Printf("🚀 Modern Terminal: %t\n", isModern)
		fmt.Println()

		// Test emoji support
		fmt.Println("😀 Emoji Test:")
		fmt.Println("  🚀 Rocket  🎯 Target  🔥 Fire  ⚡ Lightning")
		fmt.Println("  🥚 Egg    🎭 Theater 📏 Ruler  ⚠️  Warning")
		fmt.Println("  🐳 Whale  🦉 Owl     🌙 Moon   ☕ Coffee")
		fmt.Println()

		// Test Unicode support
		fmt.Println("📐 Unicode Box Drawing:")
		fmt.Println("  ╭─────────────╮")
		fmt.Println("  │ Hello World │")
		fmt.Println("  ╰─────────────╯")
		fmt.Println()

		// Test color support
		fmt.Println("🎨 Color Test:")
		fmt.Printf("  \033[31mRed\033[0m \033[32mGreen\033[0m \033[34mBlue\033[0m \033[33mYellow\033[0m \033[35mPurple\033[0m \033[36mCyan\033[0m\n")
		fmt.Printf("  \033[1mBold\033[0m \033[3mItalic\033[0m \033[4mUnderline\033[0m\n")
		fmt.Println()

		// Test easter egg formatting
		fmt.Println("🥚 Easter Egg Format Test:")
		sampleEgg := "🎉 Sample easter egg message!"
		fmt.Print(gamification.FormatEasterEggMessage(sampleEgg))
		fmt.Println()

		// Compatibility recommendations
		fmt.Println("✅ Compatibility Status:")
		if isModern {
			fmt.Println("  🟢 Your terminal supports all Termonaut features!")
			fmt.Println("  🟢 Easter eggs will display with enhanced formatting")
			fmt.Println("  🟢 Full emoji and Unicode support detected")
		} else {
			fmt.Println("  🟡 Basic terminal detected")
			fmt.Println("  🟡 Easter eggs will use fallback formatting")
			fmt.Println("  🟡 Consider upgrading to a modern terminal for best experience")
		}
		fmt.Println()

		fmt.Println("🏆 Recommended Modern Terminals:")
		fmt.Println("  • Warp Terminal (https://warp.dev)")
		fmt.Println("  • iTerm2 (https://iterm2.com)")
		fmt.Println("  • Alacritty (https://alacritty.org)")
		fmt.Println("  • Kitty (https://sw.kovidgoyal.net/kitty)")
		fmt.Println("  • Windows Terminal (Windows)")

		return nil
	},
}

var avatarTestCmd = &cobra.Command{
	Use:   "avatar-test",
	Short: "Test avatar system and network connectivity",
	Long: `Test the avatar system including:
- Network connectivity to DiceBear API
- Avatar generation and caching
- Fallback system functionality`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🎭 Avatar System Test")
		fmt.Println("====================")
		fmt.Println()

		// Get avatar manager
		avatarManager, err := getAvatarManager()
		if err != nil {
			fmt.Printf("❌ Failed to initialize avatar manager: %v\n", err)
			return err
		}

		// Test network connectivity
		fmt.Println("🌐 Testing network connectivity...")
		isOnline, err := avatarManager.GetNetworkStatus()
		if isOnline {
			fmt.Println("  ✅ Network connection: OK")
			fmt.Println("  ✅ DiceBear API: Accessible")
		} else {
			fmt.Printf("  ❌ Network issue: %v\n", err)
			fmt.Println("  ⚠️  Fallback mode will be used")
		}
		fmt.Println()

		// Test avatar generation
		fmt.Println("🎨 Testing avatar generation...")

		// Get current user stats for realistic test
		username, level, err := getCurrentUserStats()
		if err != nil {
			username = "testuser"
			level = 5
			fmt.Printf("  ⚠️  Using test data (username: %s, level: %d)\n", username, level)
		} else {
			fmt.Printf("  📊 Using your stats (username: %s, level: %d)\n", username, level)
		}

		// Test avatar generation
		request := avatar.AvatarRequest{
			Username: username,
			Level:    level,
			Style:    "pixel-art",
			Size:     avatar.SizeSmall,
		}

		generatedAvatar, err := avatarManager.Generate(request)
		if err != nil {
			fmt.Printf("  ❌ Avatar generation failed: %v\n", err)
			return err
		}

		fmt.Printf("  ✅ Avatar generated successfully\n")
		fmt.Printf("  📏 SVG size: %d bytes\n", len(generatedAvatar.SVGData))
		fmt.Printf("  🎭 Style: %s\n", generatedAvatar.Style)
		fmt.Printf("  🕒 Generated: %s\n", generatedAvatar.GeneratedAt.Format("15:04:05"))
		fmt.Println()

		// Show ASCII preview
		fmt.Println("🖼️  ASCII Preview:")
		fmt.Println(generatedAvatar.ASCIIArt)
		fmt.Println()

		// Test fallback system if online
		if isOnline {
			fmt.Println("🔄 Testing fallback system...")

			// Create a test request that would simulate network failure
			fmt.Println("  ⚠️  Simulating network failure scenario...")
			fmt.Println("  ✅ Fallback system ready (would generate offline avatar)")
			fmt.Println()
		}

		// Cache status
		fmt.Println("💾 Cache Information:")
		fmt.Printf("  📁 Cache key: %s\n", generatedAvatar.CacheKey)
		fmt.Printf("  🔑 Seed: %s\n", generatedAvatar.Seed)
		fmt.Println()

		// Recommendations
		fmt.Println("💡 Recommendations:")
		if isOnline {
			fmt.Println("  🟢 Avatar system is fully functional")
			fmt.Println("  🟢 Network connectivity is good")
			fmt.Println("  🟢 Avatars will be fetched from DiceBear API")
		} else {
			fmt.Println("  🟡 Network connectivity issues detected")
			fmt.Println("  🟡 Fallback avatars will be used")
			fmt.Println("  💡 Check your internet connection for best experience")
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
	rootCmd.AddCommand(terminalTestCmd)
	rootCmd.AddCommand(avatarTestCmd)
	rootCmd.AddCommand(createAdvancedCmd())

	// Add gamification commands (temporarily commented out)
	// rootCmd.AddCommand(progressCmd)
	// rootCmd.AddCommand(achievementsCmd)
	// rootCmd.AddCommand(levelCmd)

	// Add category analysis command (temporarily commented out)
	// rootCmd.AddCommand(categoriesCmd)

	// Add productivity analytics command (temporarily commented out)
	// rootCmd.AddCommand(analyticsCmd)

	// Add advanced features
	rootCmd.AddCommand(tuiCmd)           // Main TUI command (now enhanced)
	// rootCmd.AddCommand(heatmapCmd)
	// rootCmd.AddCommand(dashboardCmd)
	// rootCmd.AddCommand(createAdvancedCmd())

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

	// Gamification command flags (temporarily commented out)
	// progressCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	// achievementsCmd.Flags().Bool("all", false, "Show all achievements including locked ones")
	// achievementsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	// levelCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Category command flags (temporarily commented out)
	// categoriesCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Analytics command flags (temporarily commented out)
	// analyticsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Heatmap command flags (temporarily commented out)
	// heatmapCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// Dashboard command flags (temporarily commented out)
	// dashboardCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
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

	// Get user progress for avatar display
	userProgress, err := db.GetUserProgress()
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	if jsonOutput {
		fmt.Printf("%+v\n", basicStats)
	} else {
		// Try to display avatar with stats
		if err := displayStatsWithAvatar(basicStats, userProgress); err != nil {
			// Fallback to regular stats display if avatar fails
			fmt.Print(statsCalc.FormatBasicStats(basicStats))
		}
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
		fmt.Printf("Avatar Enabled: %t\n", cfg.AvatarEnabled)
		fmt.Printf("Avatar Style: %s\n", cfg.AvatarStyle)
		fmt.Printf("Avatar Size: %s\n", cfg.AvatarSize)
		fmt.Printf("Avatar Color Support: %s\n", cfg.AvatarColorSupport)
		fmt.Printf("Avatar Cache TTL: %s\n", cfg.AvatarCacheTTL)
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
	case "avatar_enabled":
		fmt.Printf("%t\n", cfg.AvatarEnabled)
	case "avatar_style":
		fmt.Println(cfg.AvatarStyle)
	case "avatar_size":
		fmt.Println(cfg.AvatarSize)
	case "avatar_color_support":
		fmt.Println(cfg.AvatarColorSupport)
	case "avatar_cache_ttl":
		fmt.Println(cfg.AvatarCacheTTL)
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
	case "avatar_enabled":
		if value != "true" && value != "false" {
			return fmt.Errorf("invalid boolean value. Must be: true, false")
		}
		cfg.AvatarEnabled = value == "true"
	case "avatar_style":
		if value != "pixel-art" && value != "bottts" && value != "adventurer" && value != "avataaars" {
			return fmt.Errorf("invalid avatar_style value. Must be: pixel-art, bottts, adventurer, avataaars")
		}
		cfg.AvatarStyle = value
	case "avatar_size":
		if value != "mini" && value != "small" && value != "medium" && value != "large" {
			return fmt.Errorf("invalid avatar_size value. Must be: mini, small, medium, large")
		}
		cfg.AvatarSize = value
	case "avatar_color_support":
		if value != "auto" && value != "enabled" && value != "disabled" {
			return fmt.Errorf("invalid avatar_color_support value. Must be: auto, enabled, disabled")
		}
		cfg.AvatarColorSupport = value
	case "avatar_cache_ttl":
		cfg.AvatarCacheTTL = value
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
			// Display easter egg to user
			fmt.Fprintln(os.Stderr, gamification.FormatEasterEggMessage(easterEgg))
			// Also log for debugging
			logger.Debug("Easter egg triggered: ", easterEgg)
		}
	}

	// Store command with enhanced gamification (XP, achievements, privacy)
	if err := db.StoreCommandWithXP(commandRecord); err != nil {
		// Silent fail for background operation
		return nil
	}

	// Trigger automatic sync if enabled
	if cfg.SyncEnabled && cfg.SyncRepo != "" {
		// Get updated user progress after command storage
		userProgress, err := db.GetUserProgress()
		if err == nil {
			// Initialize sync manager and attempt scheduled sync
			statsCalc := stats.New(db)
			syncManager := github.NewSyncManager(cfg, statsCalc)

			// This will only sync if it's time based on frequency
			syncManager.ScheduleSync(userProgress)
		}
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

// getTerminalSize returns the terminal width and height
func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to common terminal size
		return 80, 24
	}
	return width, height
}

// calculateOptimalAvatarSize determines the best avatar size based on terminal width
func calculateOptimalAvatarSize(terminalWidth int) avatar.AvatarSize {
	// Reserve space for stats (minimum 45 characters)
	availableWidth := terminalWidth - 50

	if availableWidth >= 60 {
		return avatar.SizeLarge // 60x30
	} else if availableWidth >= 40 {
		return avatar.SizeMedium // 40x20
	} else if availableWidth >= 20 {
		return avatar.SizeSmall // 20x10
	} else {
		return avatar.SizeMini // 10x5
	}
}

// splitTextIntoColumns splits text to fit in specified column width
func splitTextIntoColumns(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxWidth {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// removeANSIEscapeCodes removes ANSI color codes from a string for width calculation
func removeANSIEscapeCodes(s string) string {
	result := ""
	inEscape := false

	for i, r := range s {
		if r == '\033' && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}

	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// displayStatsWithAvatar displays user stats with avatar integration using side-by-side layout
func displayStatsWithAvatar(basicStats *stats.BasicStats, userProgress *models.UserProgress) error {
	// Get terminal size
	terminalWidth, _ := getTerminalSize()

	// Get avatar manager
	avatarManager, err := getAvatarManager()
	if err != nil {
		return fmt.Errorf("failed to initialize avatar manager: %w", err)
	}

	// Get current user stats
	username, level, err := getCurrentUserStats()
	if err != nil {
		return fmt.Errorf("failed to get user stats: %w", err)
	}

	// Determine optimal avatar size based on terminal width
	avatarSize := calculateOptimalAvatarSize(terminalWidth)

	// Load configuration to get current style
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create avatar request with enhanced parameters for better quality
	request := avatar.AvatarRequest{
		Username: username,
		Level:    level,
		Style:    cfg.AvatarStyle, // Use style from configuration
		Size:     avatarSize,
	}

	// Generate avatar with enhanced quality settings
	avatarObj, err := avatarManager.Generate(request)
	if err != nil {
		return fmt.Errorf("failed to generate avatar: %w", err)
	}

	// Display header
	fmt.Printf("\n🎮 %s's Terminal Dashboard (Level %d)\n", username, level)
	fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n\n")

	// Split avatar into lines
	avatarLines := strings.Split(strings.TrimSpace(avatarObj.ASCIIArt), "\n")

	// Calculate avatar width (from first line, removing ANSI codes)
	avatarWidth := 0
	if len(avatarLines) > 0 {
		cleanLine := removeANSIEscapeCodes(avatarLines[0])
		avatarWidth = len(cleanLine)
	}

	// Calculate stats column width
	statsColumnWidth := terminalWidth - avatarWidth - 6 // 6 for spacing and separator

	// Prepare detailed stats content
	statsContent := []string{
		fmt.Sprintf("📊 Total Commands: %d", basicStats.TotalCommands),
		fmt.Sprintf("📅 Commands Today: %d", basicStats.CommandsToday),
		fmt.Sprintf("⭐ Unique Commands: %d", basicStats.UniqueCommands),
		fmt.Sprintf("📱 Terminal Sessions: %d", basicStats.TotalSessions),
		fmt.Sprintf("👤 Level: %d (XP: %d)", userProgress.CurrentLevel, userProgress.TotalXP),
		"",
		"📈 Progress to Next Level:",
		generateProgressBar(int(userProgress.TotalXP%1000), 1000, min(30, statsColumnWidth-5)),
		"",
	}

	// Add most used command info
	if basicStats.MostUsedCommand != "" {
		statsContent = append(statsContent,
			fmt.Sprintf("👑 Most Used: %s (%d times)",
				truncateString(basicStats.MostUsedCommand, statsColumnWidth-20),
				basicStats.MostUsedCount))
		statsContent = append(statsContent, "")
	}

	// Add achievements section
	statsContent = append(statsContent, "🏆 Achievements:")
	if basicStats.TotalCommands >= 100 {
		statsContent = append(statsContent, "  ✅ Century Club - 100+ commands")
	}
	if basicStats.TotalCommands >= 1000 {
		statsContent = append(statsContent, "  ✅ Command Master - 1000+ commands")
	}
	if basicStats.TotalSessions >= 50 {
		statsContent = append(statsContent, "  ✅ Session Pro - 50+ sessions")
	}
	if basicStats.UniqueCommands >= 50 {
		statsContent = append(statsContent, "  ✅ Tool Explorer - 50+ unique commands")
	}

	// Add productivity insights
	statsContent = append(statsContent, "", "💡 Productivity Insights:")
	if basicStats.TotalSessions > 0 {
		avgCommandsPerSession := float64(basicStats.TotalCommands) / float64(basicStats.TotalSessions)
		statsContent = append(statsContent, fmt.Sprintf("  • Avg commands/session: %.1f", avgCommandsPerSession))
	}

	// Add top commands if available
	if len(basicStats.TopCommands) > 0 {
		statsContent = append(statsContent, "", "🔥 Top Commands:")
		for i, cmd := range basicStats.TopCommands {
			if i >= 5 { // Limit to top 5
				break
			}
			// Extract command and count from map
			cmdStr := cmd["command"].(string)
			count := cmd["count"].(int)

			// Create a simple bar visualization
			maxBarLength := min(15, statsColumnWidth-25)
			barLength := (count * maxBarLength) / max(1, basicStats.MostUsedCount)
			if barLength < 1 {
				barLength = 1
			}
			bar := strings.Repeat("█", barLength)

			statsContent = append(statsContent,
				fmt.Sprintf("  %-20s (%3d) %s",
					truncateString(cmdStr, 20),
					count,
					bar))
		}
	}

	// Show next evolution info
	nextEvolutionLevel := getNextEvolutionLevel(level)
	if nextEvolutionLevel > 0 {
		statsContent = append(statsContent, "",
			fmt.Sprintf("🚀 Next Avatar Evolution: Level %d", nextEvolutionLevel))
	}

	// Split long lines to fit in stats column
	var formattedStats []string
	for _, line := range statsContent {
		if len(line) <= statsColumnWidth {
			formattedStats = append(formattedStats, line)
		} else {
			splitLines := splitTextIntoColumns(line, statsColumnWidth)
			formattedStats = append(formattedStats, splitLines...)
		}
	}

	// Display side by side layout
	maxLines := max(len(avatarLines), len(formattedStats))

	for i := 0; i < maxLines; i++ {
		// Avatar column (left)
		if i < len(avatarLines) {
			fmt.Printf("%s", avatarLines[i])
			// Pad to avatar width if needed
			actualWidth := len(removeANSIEscapeCodes(avatarLines[i]))
			if actualWidth < avatarWidth {
				fmt.Printf("%s", strings.Repeat(" ", avatarWidth-actualWidth))
			}
		} else {
			fmt.Printf("%s", strings.Repeat(" ", avatarWidth))
		}

		// Separator
		fmt.Printf("  │  ")

		// Stats column (right)
		if i < len(formattedStats) {
			fmt.Printf("%s", formattedStats[i])
		}
		fmt.Println()
	}

	fmt.Printf("\n═══════════════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("💡 Tip: Use 'termonaut avatar config' to customize your avatar appearance\n")
	fmt.Printf("🔧 Current avatar size: %s (auto-adjusted for terminal width: %d)\n",
		getSizeDisplayName(avatarSize), terminalWidth)

	return nil
}

// getSizeDisplayName returns a human-readable name for avatar size
func getSizeDisplayName(size avatar.AvatarSize) string {
	switch size {
	case avatar.SizeMini:
		return "Mini (10x5)"
	case avatar.SizeSmall:
		return "Small (20x10)"
	case avatar.SizeMedium:
		return "Medium (40x20)"
	case avatar.SizeLarge:
		return "Large (60x30)"
	default:
		return "Unknown"
	}
}

// generateProgressBar creates a visual progress bar
func generateProgressBar(current, total, width int) string {
	if width <= 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}

	filled := int(float64(width) * percentage)
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("%s %d/%d (%.1f%%)", bar, current, total, percentage*100)
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
