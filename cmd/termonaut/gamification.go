package main

import (
	"fmt"

	"github.com/oiahoon/termonaut/internal/config"
	"github.com/oiahoon/termonaut/internal/database"
	"github.com/oiahoon/termonaut/internal/gamification"
	"github.com/oiahoon/termonaut/internal/stats"
	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show your Termonaut progress and achievements",
	Long: `Display comprehensive gamification statistics including your current level,
XP progress, earned achievements, and progress towards next achievements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProgressCommand(cmd, args)
	},
}

var achievementsCmd = &cobra.Command{
	Use:   "achievements",
	Short: "Show your achievements and badges",
	Long: `Display all earned achievements and progress towards unlocking new ones.
Use --all to see all available achievements including locked ones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAchievementsCommand(cmd, args)
	},
}

var levelCmd = &cobra.Command{
	Use:   "level",
	Short: "Show your current level and XP information",
	Long: `Display detailed information about your current level, XP progress,
and what's needed to reach the next level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLevelCommand(cmd, args)
	},
}

func runProgressCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if gamification is enabled
	if !cfg.ShowGamification {
		fmt.Println("🎮 Gamification is disabled. Enable it with:")
		fmt.Println("  termonaut config set show_gamification true")
		return nil
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

	// Get gamification stats
	gamificationStats, err := statsCalc.GetGamificationStats()
	if err != nil {
		return fmt.Errorf("failed to get gamification stats: %w", err)
	}

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		fmt.Printf("%+v\n", gamificationStats)
	} else {
		fmt.Print(statsCalc.FormatGamificationStats(gamificationStats))
	}

	return nil
}

func runAchievementsCommand(cmd *cobra.Command, args []string) error {
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

	// Get user achievements
	earnedAchievements, err := db.GetUserAchievements()
	if err != nil {
		return fmt.Errorf("failed to get achievements: %w", err)
	}

	// Get flags
	showAll, _ := cmd.Flags().GetBool("all")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if jsonOutput {
		if showAll {
			// Show all achievements with progress
			achievementManager := gamification.NewAchievementManager()
			allAchievements := achievementManager.GetAllAchievements()
			fmt.Printf("%+v\n", allAchievements)
		} else {
			fmt.Printf("%+v\n", earnedAchievements)
		}
		return nil
	}

	// Format output
	fmt.Println("🏆 Termonaut Achievements")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(earnedAchievements) > 0 {
		fmt.Println("\n✅ Earned Achievements:")
		for _, achievement := range earnedAchievements {
			fmt.Printf("  %s %s - Earned %s (+%d XP)\n",
				achievement.Achievement.Icon,
				achievement.Achievement.Name,
				achievement.EarnedAt.Format("Jan 2, 2006"),
				achievement.Achievement.XPReward)
			if achievement.Achievement.Description != "" {
				fmt.Printf("     %s\n", achievement.Achievement.Description)
			}
		}
	} else {
		fmt.Println("\n🎯 No achievements earned yet. Start using the terminal to unlock your first achievement!")
	}

	if showAll {
		// Show all available achievements
		achievementManager := gamification.NewAchievementManager()
		allAchievements := achievementManager.GetAllAchievements()

		// Get user stats for progress calculation
		gamificationStats, err := db.GetGamificationStats()
		if err == nil {
			fmt.Println("\n🔒 Available Achievements:")
			for _, achievement := range allAchievements {
				// Skip if already earned
				if _, earned := earnedAchievements[achievement.ID]; earned {
					continue
				}

				// Skip hidden achievements unless close to completion
				if achievement.Hidden {
					continue
				}

				progress, target, percentage := achievementManager.GetProgressToAchievement(achievement.ID, gamificationStats)

				rarity := ""
				if achievement.Rare {
					rarity = " ⭐"
				}

				fmt.Printf("  %s %s%s - %d/%d (%.0f%%)\n",
					achievement.Icon,
					achievement.Name,
					rarity,
					progress,
					target,
					percentage)
				fmt.Printf("     %s (+%d XP)\n", achievement.Description, achievement.XPReward)
			}
		}
	}

	fmt.Printf("\n📈 Achievement Progress: %d earned\n", len(earnedAchievements))

	return nil
}

func runLevelCommand(cmd *cobra.Command, args []string) error {
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

	// Get user progress
	userProgress, err := db.GetUserProgress()
	if err != nil {
		return fmt.Errorf("failed to get user progress: %w", err)
	}

	// Calculate progress info
	progressCalc := gamification.NewProgressCalculator(nil)
	progressInfo := progressCalc.CalculateProgress(
		userProgress.TotalXP,
		0, // commands today (not needed for this display)
		userProgress.CurrentStreak,
	)

	// Check output format
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		fmt.Printf("%+v\n", progressInfo)
		return nil
	}

	// Format level information
	fmt.Println("🌟 Termonaut Level Information")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\n🚀 Current Level: %d\n", progressInfo.CurrentLevel)
	fmt.Printf("🎭 Title: %s\n", progressInfo.LevelTitle)
	fmt.Printf("💎 Total XP: %d\n\n", progressInfo.TotalXP)

	// Progress to next level
	fmt.Println("📊 Progress to Next Level:")

	// Create progress bar
	progressBar := createProgressBar(progressInfo.ProgressPercent, 40)

	fmt.Printf("Level %d %s Level %d\n",
		progressInfo.CurrentLevel, progressBar, progressInfo.CurrentLevel+1)
	fmt.Printf("XP: %d/%d (%.1f%% complete)\n",
		progressInfo.XPInCurrentLevel,
		progressInfo.XPForCurrentLevel,
		progressInfo.ProgressPercent)
	fmt.Printf("🎯 XP needed for next level: %d\n\n", progressInfo.XPToNextLevel)

	// Show next few levels
	fmt.Println("🎯 Upcoming Levels:")
	levelCalc := gamification.NewLevelCalculator()
	for i := 1; i <= 3; i++ {
		nextLevel := progressInfo.CurrentLevel + i
		if nextLevel > 100 {
			break
		}
		xpRequired := levelCalc.CalculateXPForLevel(nextLevel)
		title := levelCalc.GetLevelTitle(nextLevel)
		fmt.Printf("  Level %d: %s (Total XP: %d)\n", nextLevel, title, xpRequired)
	}

	return nil
}

// Helper function to create progress bar
func createProgressBar(percentage float64, width int) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}

	filled := int(percentage * float64(width) / 100)

	// Use Unicode block characters for better visualization
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return fmt.Sprintf("[%s]", bar)
}
