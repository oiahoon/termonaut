package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/oiahoon/termonaut/internal/avatar"
	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var avatarCmd = &cobra.Command{
	Use:   "avatar",
	Short: "Manage your terminal avatar",
	Long: `Avatar system provides visual representation through ASCII art that evolves with your level.

Your avatar is generated deterministically based on your username and level,
and changes appearance as you progress through different milestones.`,
}

var avatarShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display your current avatar",
	Long: `Display your current avatar as ASCII art.

The avatar is generated based on your current level and will evolve
as you progress. Different sizes are available for different contexts.`,
	RunE: runAvatarShowCommand,
}

var avatarRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh your avatar cache",
	Long: `Force regeneration of your avatar, bypassing the cache.

This is useful when you want to see changes immediately or if
there are issues with the cached avatar.`,
	RunE: runAvatarRefreshCommand,
}

var avatarConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure avatar settings",
	Long: `Configure avatar settings such as style, size, and display options.

Available styles:
- pixel-art (recommended): Retro 8-bit style
- bottts: Robot-themed avatars  
- adventurer: Fantasy character style
- avataaars: Modern cartoon style`,
	RunE: runAvatarConfigCommand,
}

var avatarPreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview avatar at different levels",
	Long: `Preview how your avatar will look at different levels.

This helps you see the evolution of your avatar as you progress
without actually changing your current level.`,
	RunE: runAvatarPreviewCommand,
}

var avatarStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show avatar system statistics",
	Long: `Display statistics about the avatar system including cache usage,
generation times, and other performance metrics.`,
	RunE: runAvatarStatsCommand,
}

// Command flags
var (
	avatarSize    string
	avatarStyle   string
	avatarLevel   int
	avatarForce   bool
	avatarVerbose bool
)

func init() {
	// Add subcommands
	avatarCmd.AddCommand(avatarShowCmd)
	avatarCmd.AddCommand(avatarRefreshCmd)
	avatarCmd.AddCommand(avatarConfigCmd)
	avatarCmd.AddCommand(avatarPreviewCmd)
	avatarCmd.AddCommand(avatarStatsCmd)

	// Add flags
	avatarShowCmd.Flags().StringVarP(&avatarSize, "size", "s", "small", "Avatar size (mini, small, medium, large)")
	avatarShowCmd.Flags().BoolVarP(&avatarVerbose, "verbose", "v", false, "Show verbose avatar information")

	avatarRefreshCmd.Flags().BoolVarP(&avatarForce, "force", "f", false, "Force refresh even if cache is valid")

	avatarConfigCmd.Flags().StringVar(&avatarStyle, "style", "", "Set avatar style (pixel-art, bottts, adventurer, avataaars)")
	avatarConfigCmd.Flags().StringVar(&avatarSize, "size", "", "Set default avatar size (mini, small, medium, large)")

	avatarPreviewCmd.Flags().IntVarP(&avatarLevel, "level", "l", 0, "Level to preview (required)")
	avatarPreviewCmd.Flags().StringVarP(&avatarSize, "size", "s", "small", "Avatar size for preview")
	avatarPreviewCmd.MarkFlagRequired("level")

	avatarStatsCmd.Flags().BoolVarP(&avatarVerbose, "verbose", "v", false, "Show detailed statistics")

	// Add to root command
	rootCmd.AddCommand(avatarCmd)
}

func runAvatarShowCommand(cmd *cobra.Command, args []string) error {
	// Get avatar manager
	manager, err := getAvatarManager()
	if err != nil {
		return fmt.Errorf("failed to initialize avatar manager: %w", err)
	}

	// Get current user stats to determine level
	username, level, err := getCurrentUserStats()
	if err != nil {
		return fmt.Errorf("failed to get user stats: %w", err)
	}

	// Load configuration to get current style
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse size
	size, err := parseAvatarSize(avatarSize)
	if err != nil {
		return fmt.Errorf("invalid size: %w", err)
	}

	// Create avatar request using configured style
	request := avatar.AvatarRequest{
		Username: username,
		Level:    level,
		Style:    cfg.AvatarStyle, // Use style from configuration
		Size:     size,
	}

	// Generate avatar
	avatarObj, err := manager.Generate(request)
	if err != nil {
		return fmt.Errorf("failed to generate avatar: %w", err)
	}

	// Display avatar
	fmt.Printf("🎨 %s's Avatar - Level %d\n", username, level)
	fmt.Println(strings.Repeat("═", 40))
	fmt.Println()
	fmt.Println(avatarObj.ASCIIArt)
	fmt.Println()

	if avatarVerbose {
		fmt.Printf("📊 Avatar Info:\n")
		fmt.Printf("• Style: %s\n", avatarObj.Style)
		fmt.Printf("• Seed: %s\n", avatarObj.Seed)
		fmt.Printf("• Generated: %s\n", avatarObj.GeneratedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("• Cache Key: %s\n", avatarObj.CacheKey)
		fmt.Printf("• Size: %dx%d\n", avatarObj.Size.ASCIIWidth, avatarObj.Size.ASCIIHeight)
		fmt.Println()
	}

	// Show next evolution info
	nextEvolutionLevel := getNextEvolutionLevel(level)
	if nextEvolutionLevel > 0 {
		fmt.Printf("🚀 Next Evolution: Level %d\n", nextEvolutionLevel)
	} else {
		fmt.Printf("👑 Maximum evolution reached!\n")
	}

	return nil
}

func runAvatarRefreshCommand(cmd *cobra.Command, args []string) error {
	manager, err := getAvatarManager()
	if err != nil {
		return fmt.Errorf("failed to initialize avatar manager: %w", err)
	}

	username, _, err := getCurrentUserStats()
	if err != nil {
		return fmt.Errorf("failed to get user stats: %w", err)
	}

	fmt.Printf("🔄 Refreshing avatar for %s...\n", username)

	err = manager.Refresh(username)
	if err != nil {
		return fmt.Errorf("failed to refresh avatar: %w", err)
	}

	fmt.Println("✅ Avatar refreshed successfully!")
	fmt.Println("💡 Use 'termonaut avatar show' to see your updated avatar")

	return nil
}

func runAvatarConfigCommand(cmd *cobra.Command, args []string) error {
	// Load current configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Handle setting values
	if avatarStyle != "" {
		fmt.Printf("Setting avatar style to: %s\n", avatarStyle)
		cfg.AvatarStyle = avatarStyle
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("✅ Avatar style updated!")
	}

	if avatarSize != "" {
		fmt.Printf("Setting default avatar size to: %s\n", avatarSize)
		cfg.AvatarSize = avatarSize
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("✅ Avatar size updated!")
	}

	// Show current configuration if no changes were made
	if avatarStyle == "" && avatarSize == "" {
		fmt.Println("🎨 Current Avatar Configuration:")
		fmt.Println(strings.Repeat("═", 40))
		fmt.Printf("• Enabled: %t\n", cfg.AvatarEnabled)
		fmt.Printf("• Style: %s\n", cfg.AvatarStyle)
		fmt.Printf("• Size: %s\n", cfg.AvatarSize)
		fmt.Printf("• Color Support: %s\n", cfg.AvatarColorSupport)
		fmt.Printf("• Cache TTL: %s\n", cfg.AvatarCacheTTL)
		fmt.Println()
		
		fmt.Println("📚 Available Options:")
		fmt.Println()
		fmt.Println("Styles:")
		fmt.Println("• pixel-art (recommended): Retro 8-bit style, optimized for terminals")
		fmt.Println("• bottts: Robot-themed avatars with geometric shapes")
		fmt.Println("• adventurer: Fantasy character style with variety")
		fmt.Println("• avataaars: Modern cartoon style with many options")
		fmt.Println()
		fmt.Println("Sizes:")
		fmt.Println("• mini (10x5): Compact size for prompts and inline display")
		fmt.Println("• small (20x10): Good for stats and compact views")
		fmt.Println("• medium (40x20): Standard size for standalone display")
		fmt.Println("• large (60x30): Detailed view with maximum clarity")
		fmt.Println()
		fmt.Println("💡 Usage Examples:")
		fmt.Println("  termonaut avatar config --style pixel-art")
		fmt.Println("  termonaut avatar config --size medium")
		fmt.Println("  termonaut config set avatar_enabled true")
		fmt.Println("  termonaut config set avatar_color_support auto")
	}

	return nil
}

func runAvatarPreviewCommand(cmd *cobra.Command, args []string) error {
	if avatarLevel <= 0 {
		return fmt.Errorf("level must be positive")
	}

	manager, err := getAvatarManager()
	if err != nil {
		return fmt.Errorf("failed to initialize avatar manager: %w", err)
	}

	username, _, err := getCurrentUserStats()
	if err != nil {
		return fmt.Errorf("failed to get user stats: %w", err)
	}

	// Load configuration to get current style
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	size, err := parseAvatarSize(avatarSize)
	if err != nil {
		return fmt.Errorf("invalid size: %w", err)
	}

	request := avatar.AvatarRequest{
		Username: username,
		Level:    avatarLevel,
		Style:    cfg.AvatarStyle, // Use style from configuration
		Size:     size,
	}

	avatarObj, err := manager.Generate(request)
	if err != nil {
		return fmt.Errorf("failed to generate preview avatar: %w", err)
	}

	fmt.Printf("🔮 Preview: %s at Level %d\n", username, avatarLevel)
	fmt.Println(strings.Repeat("═", 40))
	fmt.Println()
	fmt.Println(avatarObj.ASCIIArt)
	fmt.Println()
	fmt.Printf("💡 This is how your avatar will look at level %d\n", avatarLevel)

	return nil
}

func runAvatarStatsCommand(cmd *cobra.Command, args []string) error {
	_, err := getAvatarManager()
	if err != nil {
		return fmt.Errorf("failed to initialize avatar manager: %w", err)
	}

	// Get cache statistics
	// Note: This would need to be implemented in the manager
	fmt.Println("📊 Avatar System Statistics")
	fmt.Println(strings.Repeat("═", 40))
	fmt.Println()
	fmt.Println("Cache Status:")
	fmt.Println("• Total Entries: N/A") // TODO: Implement
	fmt.Println("• Valid Entries: N/A")
	fmt.Println("• Cache Hit Rate: N/A")
	fmt.Println("• Total Size: N/A")
	fmt.Println()
	fmt.Println("Generation Stats:")
	fmt.Println("• Total Generations: N/A")
	fmt.Println("• Average Generation Time: N/A")
	fmt.Println("• API Success Rate: N/A")

	if avatarVerbose {
		fmt.Println()
		fmt.Println("Detailed Statistics:")
		fmt.Println("• Cache Directory: ~/.termonaut/avatars")
		fmt.Println("• Supported Styles: pixel-art, bottts, adventurer, avataaars")
		fmt.Println("• Available Sizes: mini, small, medium, large")
	}

	return nil
}

// Helper functions

func getAvatarManager() (*avatar.AvatarManager, error) {
	config := avatar.GetDefaultConfig()
	return avatar.NewAvatarManager(config)
}

func getCurrentUserStats() (username string, level int, error error) {
	// Get username from environment or system
	username = os.Getenv("USER")
	if username == "" {
		username = "user"
	}
	
	// Load configuration and database to get real level
	cfg, err := config.Load()
	if err != nil {
		// Fallback to level 1 if config fails
		return username, 1, nil
	}

	// Initialize database
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Minimal logging for avatar operations
	db, err := database.New(config.GetDataDir(cfg), logger)
	if err != nil {
		// Fallback to level 1 if database fails
		return username, 1, nil
	}
	defer db.Close()

	// Get user progress
	userProgress, err := db.GetUserProgress()
	if err != nil {
		// Fallback to level 1 if user progress fails
		return username, 1, nil
	}

	return username, userProgress.CurrentLevel, nil
}

func parseAvatarSize(sizeStr string) (avatar.AvatarSize, error) {
	switch strings.ToLower(sizeStr) {
	case "mini":
		return avatar.SizeMini, nil
	case "small":
		return avatar.SizeSmall, nil
	case "medium":
		return avatar.SizeMedium, nil
	case "large":
		return avatar.SizeLarge, nil
	default:
		return avatar.AvatarSize{}, fmt.Errorf("invalid size '%s', must be one of: mini, small, medium, large", sizeStr)
	}
}

func getNextEvolutionLevel(currentLevel int) int {
	evolutionLevels := []int{5, 10, 20, 50, 100}
	
	for _, level := range evolutionLevels {
		if currentLevel < level {
			return level
		}
	}
	
	return 0 // No more evolutions
} 